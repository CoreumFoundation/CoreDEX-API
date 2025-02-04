package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
)

type OHLCPointResponse [6]interface{}

var allowedValues = map[string]bool{
	"1m":  true,
	"3m":  true,
	"5m":  true,
	"15m": true,
	"30m": true,
	"1h":  true,
	"3h":  true,
	"6h":  true,
	"12h": true,
	"1d":  true,
	"3d":  true,
	"1w":  true,
}

var (
	periodRegex             = regexp.MustCompile(`^(\d+)([a-zA-Z]+)$`)
	ErrIncorrectRequestParm = errors.New("incorrect request parameter")
)

// Input: ["1m","3m","5m","15m","30m","1h","3h","6h","12h","1d","3d","1w"]
// Or invalid input
// Output:
// ohlcgrpc.Period
func HttpPeriodToPeriod(value string) (*ohlcgrpc.Period, error) {
	if !allowedValues[value] {
		return nil, ErrIncorrectRequestParm
	}
	period := &ohlcgrpc.Period{}
	matches := periodRegex.FindStringSubmatch(value)
	if len(matches) != 3 {
		return period, fmt.Errorf("invalid value: %s", value)
	}

	duration, _ := strconv.Atoi(matches[1])
	period.Duration = int32(duration)
	period.PeriodType = mapStringToPeriodType(matches[2])
	return period, nil
}

func mapStringToPeriodType(s string) ohlcgrpc.PeriodType {
	switch s {
	case "m":
		return ohlcgrpc.PeriodType_PERIOD_TYPE_MINUTE
	case "h":
		return ohlcgrpc.PeriodType_PERIOD_TYPE_HOUR
	case "d":
		return ohlcgrpc.PeriodType_PERIOD_TYPE_DAY
	case "w":
		return ohlcgrpc.PeriodType_PERIOD_TYPE_WEEK
	default:
		return ohlcgrpc.PeriodType_PERIOD_TYPE_DO_NOT_USE
	}
}

// Outliers are correct and can occur due to ledger behaviour: A very small trade can occur at a very high price.
// This disturbs the graph and does not represent the reality of the pricing which occurs.
// Since such a transaction can occur as a single transaction in a single minute (so no other data to evaluate it against to be able to identify it as an outlier)
// the choice is made to smooth the outliers on retrieval.
// The smoothing is done by replacing the outlier with the average of the previous and next value if the current value deviates outside the norm and is not correctable within the current time interval
// Correction within the current time interval is done by replacing the outlier with another ohlc value from within the time interval.
// Outlier detection to see if the current minute with all respectable values are within range of the next interval, and for the last value to see if it is within range of the 1 to last interval.
func SmoothOutliers(series []*ohlcgrpc.OHLC, index int) *ohlcgrpc.OHLC {
	data := series[index]
	// We replace inline the values which are deviating.
	// Assumption for the very simple scenario is that the values provided will have a very small divider used, so will show a much too large value for the high price.
	// If the high price deviations outside of the norm, we will inspect open/close too and replace accordingly.
	// The second check this code does is to see if the previous or next value has a high price which is without the range of reasonable values.
	// If uses the previousohlc by default except for the first value, for that it will use the next value.
	// Since this detection only works if there is more than 1 trade in the base data for calculate the OHLC, a single trade in a single ohlc will not show up and will not be corrected
	if data.High/data.Low > 3 {
		o := &ohlcgrpc.OHLC{
			Timestamp: data.Timestamp,
			Low:       data.Low,
		}
		newHigh := data.Low
		// That is a large drop (or increase), assuming calculation error:
		switch {
		case data.High/data.Open > 3:
			// Open is correct
			newHigh = data.Open
			o.Open = data.Open
			fallthrough
		case data.High/data.Close > 3:
			// Close is correct
			o.Close = data.Close
			if data.Close > newHigh {
				newHigh = data.Close
			}
			fallthrough
		case data.High/data.Open < 2:
			// Open is incorrect (the high was the open)
			o.Open = newHigh
			fallthrough
		case data.High/data.Close < 2:
			// Close is incorrect (the high was the close)
			o.Close = newHigh
		}
		o.High = newHigh
		return o
	}
	// Check for major deviations in the data
	// Missing scenario is when the value being retrieved is the current minute and that would have a graphing deviation.
	if len(series) > 1 {
		lookup := index - 1
		if index == 0 {
			lookup = 1
		}
		if data.High/series[lookup].Low > 10 {
			// Severe deviation found: all values are incorrect: Use backfill style to smooth the data
			o := &ohlcgrpc.OHLC{
				Timestamp: series[lookup].Timestamp,
				Open:      series[lookup].Close,
				High:      series[lookup].Close,
				Low:       series[lookup].Close,
				Close:     series[lookup].Close,
				Volume:    0.0,
			}
			return o
		}
	}

	return data
}
