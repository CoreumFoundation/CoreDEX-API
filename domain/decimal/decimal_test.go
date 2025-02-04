package decimal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_ToBigInt(t *testing.T) {
	dec, _ := decimal.NewFromString("1234567890123456789012345678901234567890")
	i, e := ToBigInt(dec)
	assert.Equal(t, i, int64(1234567890123456789))
	assert.Equal(t, e, int32(21))

	dec, _ = decimal.NewFromString("92233720368547758070")
	i, e = ToBigInt(dec)
	assert.Equal(t, i, int64(9223372036854775807))
	assert.Equal(t, e, int32(1))

	dec, _ = decimal.NewFromString("9223372036854775807")
	i, e = ToBigInt(dec)
	assert.Equal(t, i, int64(9223372036854775807))
	assert.Equal(t, e, int32(0))

	dec, _ = decimal.NewFromString("223372036854775807")
	i, e = ToBigInt(dec)
	assert.Equal(t, i, int64(223372036854775807))
	assert.Equal(t, e, int32(0))

	dec, _ = decimal.NewFromString("9323372036854775807")
	i, e = ToBigInt(dec)
	assert.Equal(t, i, int64(932337203685477580))
	assert.Equal(t, e, int32(1))

	dec, _ = decimal.NewFromString("9223363036854775807")
	i, e = ToBigInt(dec)
	assert.Equal(t, i, int64(9223363036854775807))
	assert.Equal(t, e, int32(0))
}

func Test_NewDecimal(t *testing.T) {
	d, err := NewDecimal("100000000utestcore")
	assert.NoError(t, err)
	assert.Equal(t, d.Value, int64(100000000))
	assert.Equal(t, d.Exp, int32(0))

	d, err = NewDecimal("1000000000usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc")
	assert.NoError(t, err)
	assert.Equal(t, d.Value, int64(1000000000))
	assert.Equal(t, d.Exp, int32(0))

	d, err = NewDecimal("160500000ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D")
	assert.NoError(t, err)
	assert.Equal(t, d.Value, int64(160500000))
	assert.Equal(t, d.Exp, int32(0))
}
