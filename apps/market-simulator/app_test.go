package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNumExpPrice(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.0, "1"},
		{12.0, "12"},
		{1.1, "11e-1"},
		{72.185474, "721855e-4"},
		{75.518003, "75518e-3"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("num=%f", test.input), func(t *testing.T) {
			result := buildNumExpPrice(test.input)
			assert.Equal(t, test.expected, result.String())
		})
	}
}
