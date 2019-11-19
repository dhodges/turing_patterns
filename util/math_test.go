package util

import (
	"testing"
)

func testRound(t *testing.T, f float64, i int) {
	if Round(f) != i {
		t.Errorf("Round(%.4f) is %v, but it should be %v", f, Round(f), i)
	}
}

func testFixed(t *testing.T, num float64, precision int, result float64) {
	if ToFixed(num, precision) != result {
		t.Errorf("ToFixed(%.4f, %v) is %v, but it should be %v", num, precision, ToFixed(num, precision), result)
	}
}

func testConstrain(t *testing.T, min, n, max float64) {
	result := Constrain(min, n, max)
	if result < min || max < result {
		t.Errorf("Constrain(%.2f, %.2f, %.2f) is %.2f which is outside the constraint", min, n, max, result)
	}
}

func TestRoundFunction(t *testing.T) {
	testRound(t, 0.6, 1)
	testRound(t, 0.5, 1)
	testRound(t, 0.4, 0)
	testRound(t, 10.6, 11)
	testRound(t, 10.5, 11)
	testRound(t, 10.4, 10)
}

func TestToFixedFunction(t *testing.T) {
	testFixed(t, 10.1234, 2, 10.12)
	testFixed(t, 15.9876, 3, 15.988)
	testFixed(t, 100.2468, 4, 100.2468)
}

func TestConstrain(t *testing.T) {
	testConstrain(t, 0, 0, 5)
	testConstrain(t, 0, 5, 5)
	testConstrain(t, 0, 3, 5)
	testConstrain(t, 10, 15, 20)
}
