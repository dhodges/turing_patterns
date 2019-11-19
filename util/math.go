package util

import "math"

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
