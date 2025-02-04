package coreum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePrice(t *testing.T) {
	tests := []struct {
		strPrice    string
		parsedPrice string
		wantErr     bool
	}{
		{
			// normal price
			strPrice:    "1.231",
			parsedPrice: "1231e-3",
			wantErr:     false,
		},
		{
			// normal price
			strPrice:    "423000",
			parsedPrice: "423e3",
			wantErr:     false,
		},
		{
			// normal price
			strPrice:    "423000 ", // extra space at the end
			parsedPrice: "423e3",   // extra space at the end
			wantErr:     false,
		},
		{
			// normal price
			strPrice:    " 423000", // extra space at the start
			parsedPrice: "423e3",   // extra space at the start
			wantErr:     false,
		},
		{
			// normal price
			strPrice:    "323141245",
			parsedPrice: "323141245",
			wantErr:     false,
		},
		{
			// zero price
			strPrice: "0",
			wantErr:  true,
		},
		{
			// invalid zero price with exponent
			strPrice: "0e1",
			wantErr:  true,
		},
		{
			// invalid price with leading
			strPrice: "01e1",
			wantErr:  true,
		},
		{
			// max uint64 num
			strPrice:    "9999999999999999999",
			parsedPrice: "9999999999999999999",
			wantErr:     false,
		},
		{
			// invalid max uint64 + 1 num
			strPrice: "18446744073709551616",
			wantErr:  true,
		},
		{
			// invalid negative num part
			strPrice: "-1",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.strPrice, func(t *testing.T) {
			got, err := ParsePrice(tt.strPrice)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.parsedPrice, got.String())
		})
	}
}
