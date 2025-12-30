package core

import "math"

// CIELABMatcher implements ColorMatcher using CIELAB color space.
type CIELABMatcher struct {
	palette *Palette
}

// NewCIELABMatcher creates a new CIELAB color matcher.
func NewCIELABMatcher(palette *Palette) *CIELABMatcher {
	return &CIELABMatcher{palette: palette}
}

// Match finds the best matching palette color for the given RGB color.
func (m *CIELABMatcher) Match(rgb [3]uint8) *PaletteColor {
	if m.palette == nil || len(m.palette.Colors) == 0 {
		return nil
	}
	
	targetLAB := RGBToLAB(rgb)
	
	var bestMatch *PaletteColor
	bestDistance := math.MaxFloat64
	
	for i := range m.palette.Colors {
		distance := DeltaE(targetLAB, m.palette.Colors[i].LAB)
		if distance < bestDistance {
			bestDistance = distance
			bestMatch = &m.palette.Colors[i]
		}
	}
	
	return bestMatch
}

// MatchWithDithering finds the best match considering dithering error.
func (m *CIELABMatcher) MatchWithDithering(rgb [3]uint8, error [3]float64) (*PaletteColor, [3]float64) {
	// Apply accumulated error to the input color
	adjustedRGB := [3]uint8{
		clampUint8(float64(rgb[0]) + error[0]),
		clampUint8(float64(rgb[1]) + error[1]),
		clampUint8(float64(rgb[2]) + error[2]),
	}
	
	// Find best match
	matched := m.Match(adjustedRGB)
	if matched == nil {
		return nil, [3]float64{0, 0, 0}
	}
	
	// Calculate quantization error
	quantError := [3]float64{
		float64(adjustedRGB[0]) - float64(matched.RGB[0]),
		float64(adjustedRGB[1]) - float64(matched.RGB[1]),
		float64(adjustedRGB[2]) - float64(matched.RGB[2]),
	}
	
	return matched, quantError
}

// SetPalette updates the palette used for matching.
func (m *CIELABMatcher) SetPalette(palette *Palette) {
	m.palette = palette
}

// clampUint8 clamps a float64 value to uint8 range [0, 255].
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
