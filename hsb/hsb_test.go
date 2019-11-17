package hsb

import (
	"fmt"
	"image/color"
	"testing"
)

func testNHSBAtoNRGBA(t *testing.T, h *NHSBA, rgba *color.NRGBA) {
	result := h.ToNRGBA()

	prefix := fmt.Sprintf("NHSBA(%v, %v, %v, %v) -> RGB, error:", h.H, h.S, h.B, h.A)

	if result.R != rgba.R {
		t.Errorf("%s r is %v, but it should be %v", prefix, result.R, rgba.R)
	}
	if result.G != rgba.G {
		t.Errorf("%s g is %v, but it should be %v", prefix, result.G, rgba.G)
	}
	if result.B != rgba.B {
		t.Errorf("%s b is %v, but it should be %v", prefix, result.B, rgba.B)
	}
	if result.A != rgba.A {
		t.Errorf("%s a is %v, but it should be %v", prefix, result.A, rgba.A)
	}
}

func TestToRGBAValues(t *testing.T) {
	testNHSBAtoNRGBA(t, NewNHSBA(0, 1.0, 1.0, 1.0), &color.NRGBA{255, 0, 0, 255})
	testNHSBAtoNRGBA(t, NewNHSBA(360, 1.0, 1.0, 1.0), &color.NRGBA{255, 0, 0, 255})
	testNHSBAtoNRGBA(t, NewNHSBA(200, 0.75, 0.75, 1.0), &color.NRGBA{48, 143, 191, 255})
	testNHSBAtoNRGBA(t, NewNHSBA(100, 0.5, 0.5, 1.0), &color.NRGBA{85, 128, 64, 255})
}
