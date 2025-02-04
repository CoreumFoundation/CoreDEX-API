package denom

import (
	"fmt"
	"regexp"
	"strings"
)

var decimalSplitRegex = regexp.MustCompile(`^(\d+)?([a-zA-Z/].*)$`)

/*
Accepts the following inputs (examples):
* 100000000utestcore
* 1000000000usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc
* 160500000ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D
* utestcore
* usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc
* ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D

There it will return the following outputs:
* utestcore
* currency: usara, issuer: devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc
* currency: "", issuer: ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D, IsIBC: true
* utestcore
* currency: usara, issuer: devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc
* currency: "", issuer: ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D, IsIBC: true
*/
func NewDenom(s string) (*Denom, error) {
	matches := decimalSplitRegex.FindStringSubmatch(s)
	if len(matches) > 2 {
		denom := &Denom{
			Currency: matches[2],
		}
		if len(denom.Currency) > 5 && denom.Currency[:4] == "ibc/" {
			denom.Currency = ""
			denom.Issuer = matches[2]
			denom.IsIBC = true

		}
		if len(matches) == 3 && !denom.IsIBC {
			// Split the issuer and currency:
			s := strings.Split(matches[2], "-")
			denom.Currency = s[0]
			if len(s) > 1 {
				denom.Issuer = s[1]
			}
		}
		// Set the denom:
		denom.Denom = denom.ToString()
		return denom, nil

	}
	return nil, fmt.Errorf("invalid denom string: %s", s)
}

func (s *Denom) ToString() string {
	if s.Issuer != "" && !s.IsIBC {
		return fmt.Sprintf("%s-%s", s.Currency, s.Issuer)
	}
	if s.IsIBC {
		return s.Issuer
	}
	return s.Currency
}
