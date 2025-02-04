package decimal

import (
	"fmt"
	"regexp"

	"github.com/shopspring/decimal"
	sdecimal "github.com/shopspring/decimal"
)

// Get the max int64 value
const maxInt64Str = "9223372036854775807"

// Get the length of the max int64 value
const maxInt64StrLen = len(maxInt64Str)

var decimalSplitRegex = regexp.MustCompile(`^(\d+)?([a-zA-Z/].*)$`)

/*
Accepts the following inputs (examples):
* 100000000utestcore
* 1000000000usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc
* 160500000ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D

It will return a Decimal struct with the following values:
* Coefficient: 100000000
* Exponent: 0
*/
func NewDecimal(s string) (*Decimal, error) {
	matches := decimalSplitRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		numberPart := matches[1]
		d, err := sdecimal.NewFromString(numberPart)
		if err != nil {
			return nil, err
		}
		return &Decimal{
			Value: d.CoefficientInt64(),
			Exp:   d.Exponent(),
		}, nil
	}
	return nil, fmt.Errorf("invalid decimal string: %s", s)
}

func FromDec(d sdecimal.Decimal) *Decimal {
	return &Decimal{
		Value: d.CoefficientInt64(),
		Exp:   d.Exponent(),
	}
}

// Non-lossless method to handle values with many decimals
// Returns lossless if the value fits in an int64 and remainder is of 10^x
func ToBigInt(dec decimal.Decimal) (int64, int32) {
	bigValue := dec.BigInt()
	// Convert the bigInt to a string
	bigStr := bigValue.String()
	// Get the length of the string
	strLen := len(bigStr)
	// Get the length of the exponent
	exponent := strLen - maxInt64StrLen
	if exponent < 0 {
		return bigValue.Int64(), 0
	}
	d, _ := decimal.NewFromString(bigStr[:maxInt64StrLen])
	// Check if all the values in the bigStr are less than the max int64 value
	// Where we break if the value of MaxInt64 is larger than the value of bigStr
	if len(bigStr) >= maxInt64StrLen {
		// Create an array with difference of the input value and the max value:
		diffArr := make([]int, maxInt64StrLen)
		for i := 0; i < maxInt64StrLen; i++ {
			diffArr[i] = int(maxInt64Str[i]) - int(bigStr[i])
		}
		// If all the values are 0 or larger than 0, we are ok
		// If there is a negative value, however it occurs after a positive value
		// If there is a negative value occuring before a positive value, we need change the exponent
		hasPos := false
		for i := 0; i < maxInt64StrLen; i++ {
			switch {
			case diffArr[i] > 0:
				hasPos = true
			case diffArr[i] < 0:
				if hasPos {
					break
				}
				exponent++
				d, _ = decimal.NewFromString(bigStr[:maxInt64StrLen-1])
				return d.CoefficientInt64(), int32(exponent)
			}
		}
	}
	// Get the int64 value
	intValue := d.CoefficientInt64()
	return intValue, int32(exponent)
}

func (d *Decimal) Float64() float64 {
	f, _ := sdecimal.New(d.Value, d.Exp).Float64()
	return f
}
