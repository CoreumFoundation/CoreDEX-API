package denom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDenom(t *testing.T) {
	d, err := NewDenom("100000000utestcore")
	assert.NoError(t, err)
	assert.Equal(t, "utestcore", d.Currency)
	assert.Equal(t, "", d.Issuer)
	assert.False(t, d.IsIBC)
	assert.Equal(t, "utestcore", d.Denom)

	d, err = NewDenom("1000000000usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc")
	assert.NoError(t, err)
	assert.Equal(t, "usara", d.Currency)
	assert.Equal(t, "devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc", d.Issuer)
	assert.False(t, d.IsIBC)
	assert.Equal(t, "usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc", d.Denom)

	d, err = NewDenom("160500000ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D")
	assert.NoError(t, err)
	assert.Equal(t, "", d.Currency)
	assert.Equal(t, "ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D", d.Issuer)
	assert.True(t, d.IsIBC)
	assert.Equal(t, "ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D", d.Denom)

	d, err = NewDenom("utestcore")
	assert.NoError(t, err)
	assert.Equal(t, "utestcore", d.Currency)
	assert.Equal(t, "", d.Issuer)
	assert.False(t, d.IsIBC)
	assert.Equal(t, "utestcore", d.Denom)

	d, err = NewDenom("usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc")
	assert.NoError(t, err)
	assert.Equal(t, "usara", d.Currency)
	assert.Equal(t, "devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc", d.Issuer)
	assert.False(t, d.IsIBC)
	assert.Equal(t, "usara-devcore1z06se860rwqazs3t7mzmm4ekev7ll0csyee2l24mzaqle6tq7hvsrrwrxc", d.Denom)

	d, err = NewDenom("ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D")
	assert.NoError(t, err)
	assert.Equal(t, "", d.Currency)
	assert.Equal(t, "ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D", d.Issuer)
	assert.True(t, d.IsIBC)
	assert.Equal(t, "ibc/E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D", d.Denom)
}
