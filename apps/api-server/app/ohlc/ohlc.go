package ohlc

import (
	"context"
	"strconv"
	"time"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ohlcgrpclient "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc/client"
)

type Application struct {
	client ohlcgrpc.OHLCServiceClient
}

func NewApplication() *Application {
	return &Application{
		client: ohlcgrpclient.Client(),
	}
}

func (s *Application) GetOHLC(ctx context.Context, ohlcOpt *ohlcgrpc.OHLCFilter) (*ohlcgrpc.OHLCs, error) {
	ohlcOpt.Backfill = true
	ohlcOpt.AllowCache = true
	return s.client.Get(ohlcgrpclient.AuthCtx(ctx), ohlcOpt)
}

func (app *Application) Get(ctx context.Context, ohlcOpt *ohlcgrpc.OHLCFilter) ([][6]interface{}, error) {
	// To get a better response pattern, round the from and to to the period time periods.
	d, err := app.GetOHLC(ctx, ohlcOpt)
	if err != nil {
		return nil, err
	}
	// Data needs to be transforment to expected rather anonymous json array without any labels:
	// [[1676937600,"0.4041000412483307","0.4159999999103348","0.39","0.3981088236566072","1442749.5076381033"]]
	// Which is an array of arrays:
	// timestamp (seconds), open, high, low, close, volume
	retvals := make([][6]interface{}, 0, len(d.OHLCs))
	// Data also needs to be filled where there are blanks: The FE graph sometimes does fill lthe blanks, sometimes does not: So better that the BE fills the blanks.
	// The fill interval can be expressed in time.Duration Minutes:
	deltaT := int64(ohlcOpt.Period.ToMinute().Duration) * int64(time.Minute) // The interval in nanoseconds.

	from := ohlcOpt.From // Used for filling the start of the return value array if no data is present
	to := ohlcOpt.To     // Used for filling the end of the return value array if no data is present
	// minTs is used to fill in the blanks in the period. The algorithm is not allowed to fill timestamps before the first known datapoint (to prevent inventing new data before a coin was initialized)
	if len(d.OHLCs) > 0 {
		minTs := d.OHLCs[0].Timestamp.Seconds
		for index, v := range d.OHLCs {
			// Smooth the outliers first: That way when we backfill the data we do not have to take the actual backfill into account.
			v = dmn.SmoothOutliers(d.OHLCs, index)
			// t := time.Unix(0, v.TimeStamp).Unix()
			for minTs < v.Timestamp.Seconds {
				if minTs >= from.Seconds-deltaT { // We want to be on the edge or 1 period in front of the requested edge
					retvals = append(retvals, dmn.OHLCPointResponse{
						minTs,
						strconv.FormatFloat(v.Close, 'f', -1, 64),
						strconv.FormatFloat(v.Close, 'f', -1, 64),
						strconv.FormatFloat(v.Close, 'f', -1, 64),
						strconv.FormatFloat(v.Close, 'f', -1, 64),
						"0.0",
					})
				}
				// Bit brutal to just iterate like this: The FROM is not aligned to the period, and exact math would be quicker/nicer, but also takes more time to write.
				minTs += deltaT
			}
			retvals = append(retvals, dmn.OHLCPointResponse{
				v.Timestamp.Seconds,
				strconv.FormatFloat(v.Open, 'f', -1, 64),
				strconv.FormatFloat(v.High, 'f', -1, 64),
				strconv.FormatFloat(v.Low, 'f', -1, 64),
				strconv.FormatFloat(v.Close, 'f', -1, 64),
				strconv.FormatFloat(v.Volume, 'f', -1, 64),
			})
			// Move the min else we would still set a timestamp before the originalFrom
			minTs += deltaT
		}
		// Fill the last periods if there is less data than expected:
		for minTs < to.Seconds {
			if minTs >= from.Seconds-deltaT { // Edge case: There is literally no data for the whole period, so this is a flat line fill from the last trade, which is now way in the past
				retvals = append(retvals, dmn.OHLCPointResponse{
					minTs,
					strconv.FormatFloat(d.OHLCs[len(d.OHLCs)-1].Close, 'f', -1, 64),
					strconv.FormatFloat(d.OHLCs[len(d.OHLCs)-1].Close, 'f', -1, 64),
					strconv.FormatFloat(d.OHLCs[len(d.OHLCs)-1].Close, 'f', -1, 64),
					strconv.FormatFloat(d.OHLCs[len(d.OHLCs)-1].Close, 'f', -1, 64),
					"0.0",
				})
			}
			minTs += deltaT
		}
	}

	return retvals, nil
}
