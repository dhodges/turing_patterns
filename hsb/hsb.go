package hsb

import (
	"image/color"
	"math"
)

// NHSBA models a non-alpha-premultiplied HSB (== HSV) color with alpha channel
type NHSBA struct {
	H float64 // Hue, within range 0 <= H <= 360°
	S float64 // Saturation
	B float64 // Brightness (or "Value")
	A float64 // Alpha
}

// constrain the given float within the range min <= n <= max
func constrain(min, n, max float64) float64 {
	n = math.Max(min, n)
	n = math.Min(n, max)
	return n
}

// NewNHSBA from four floats, constrained within sensible limits
func NewNHSBA(h, s, b, a float64) *NHSBA {
	h = constrain(0.0, h, 360.0)
	s = constrain(0.0, s, 1.0)
	b = constrain(0.0, b, 1.0)
	a = constrain(0.0, a, 1.0)
	return &NHSBA{H: h, S: s, B: b, A: a}
}

// ToNRGBA converts to a std NRGBA color
// see: http://en.wikipedia.org/wiki/HSV_color_space
// and: https://jsfiddle.net/Lamik/Lr61wqub
func (h *NHSBA) ToNRGBA() *color.NRGBA {
	c := h.B * h.S
	k := h.H / 60.0
	x := c * (1 - math.Abs(math.Mod(k, 2)-1))

	r, g, b := 0.0, 0.0, 0.0

	if 0 <= k && k <= 1 {
		r, g = c, x
	}
	if 1 < k && k <= 2 {
		r, g = x, c
	}
	if 2 < k && k <= 3 {
		g, b = c, x
	}
	if 3 < k && k <= 4 {
		g, b = x, c
	}
	if 4 < k && k <= 5 {
		r, b = x, c
	}
	if 5 < k && k <= 6 {
		r, b = c, x
	}

	m := h.B - c

	// convert from range [0 <= n <= 1.0] to [0 <= n <= 255]

	r = math.Round((r + m) * 255)
	g = math.Round((g + m) * 255)
	b = math.Round((b + m) * 255)
	a := math.Round(h.A * 255)

	return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
