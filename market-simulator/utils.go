package main

import "math/rand"

func randIntInRange(rnd *rand.Rand, minRange, maxRange int) int {
	return rnd.Intn(maxRange-minRange+1) + minRange
}

func getAnyItemByIndex[T any](slice []T, ind int) T {
	return slice[ind%len(slice)]
}
