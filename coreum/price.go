package coreum

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

func ParsePrice(input string) (types.Price, error) {
	// Trim the input and check for valid characters
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return types.Price{}, fmt.Errorf("empty input")
	}

	if input[0] == '-' {
		return types.Price{}, fmt.Errorf("price should be positive: %s", input)
	}

	if input[0] == '+' {
		input = input[1:]
	}

	// Split the input into integer and fractional parts
	parts := strings.Split(input, ".")
	if len(parts) > 2 {
		return types.Price{}, fmt.Errorf("invalid number format: %s", input)
	}

	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) == 2 {
		fractionalPart = parts[1]
	}

	// Remove leading zeros from integer part
	integerPart = strings.TrimLeft(integerPart, "0")
	if integerPart == "" {
		integerPart = "0"
	}

	// Remove trailing zeros from fractional part
	fractionalPart = strings.TrimRight(fractionalPart, "0")

	exp := 0
	if integerPart == "0" {
		for i, ch := range fractionalPart {
			if ch != '0' {
				integerPart = string(ch)
				fractionalPart = fractionalPart[i+1:]
				exp = -i - 1
				break
			}
		}
		if integerPart == "0" {
			return types.Price{}, fmt.Errorf("price should not be zero: %s", input)
		}
	} else {
		if fractionalPart == "" {
			originalLen := len(integerPart)
			integerPart = strings.TrimRight(integerPart, "0")
			exp = originalLen - len(integerPart)
		}
	}

	// Combine integer and fractional parts
	numberStr := integerPart + fractionalPart

	// Convert to uint64
	var mantissa uint64
	var err error
	if len(numberStr) > 20 {
		numberStr = numberStr[:20] // Truncate to avoid exceeding uint64
	}
	mantissa, err = strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		return types.Price{}, fmt.Errorf("mantissa exceeds uint64 range: %s", numberStr)
	}

	// Adjust exponent to account for fractional digits removed
	exp -= len(fractionalPart)

	num := ""
	if exp != 0 {
		num = fmt.Sprintf("%de%d", mantissa, exp)
	} else {
		num = fmt.Sprintf("%d", mantissa)
	}

	// Return the result in dex price format
	return types.NewPriceFromString(num)
}
