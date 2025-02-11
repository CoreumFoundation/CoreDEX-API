package main

import (
	"math"
	"time"
)

const T = float64(24)          // Period of 24 hours
const lowValue = float64(10)   // Value in other hours
const peakValue = float64(250) // Value at peak

// Calculate amplitude and baseline
// The fluctuation should center around the midpoint between peak and low values.
var midpoint = (peakValue + lowValue) / 2
var amplitude = (peakValue - lowValue) / 2

// Set phase shift to make the function peak approximately 9 hours into the cycle
var phaseShift = float64(9)

func TradesCount(t time.Time) int {
	return int(math.Ceil(amplitude*math.Cos(2*math.Pi*(float64(t.Hour())-phaseShift)/T) + midpoint))
}
