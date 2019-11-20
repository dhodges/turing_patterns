package util

import (
	"math"
	"math/rand"
)

// Round rounds the given float to an int
// see: https://stackoverflow.com/a/29786394
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// ToFixed rounds the given float to the given precision
// see: https://stackoverflow.com/a/29786394
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// Constrain the given float within the range min <= n <= max
func Constrain(min, n, max float64) float64 {
	n = math.Max(min, n)
	n = math.Min(n, max)
	return n
}

// ConstrainInt the given int within the range min <= n <= max
func ConstrainInt(min, n, max int) int {
	n = int(math.Max(float64(min), float64(n)))
	n = int(math.Min(float64(n), float64(max)))
	return n
}

// RandFloat64 generate a random float64 between the given min and max
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
