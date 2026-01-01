package core

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// LABColor represents a color in CIELAB color space.
type LABColor struct {
	L, A, B float64
}

// Palette represents a collection of colors with their CIELAB values.
type Palette struct {
	Colors []PaletteColor
}

// PaletteColor represents a color entry in a palette.
type PaletteColor struct {
	Name     string
	RGB      [3]uint8
	LAB      LABColor
	Metadata map[string]interface{} // For Minecraft-specific data (block ID, etc.)
}

// ColorMatcher is the interface for finding the closest color match.
type ColorMatcher interface {
	// Match finds the best matching palette color for the given RGB color.
	Match(rgb [3]uint8) *PaletteColor
	
	// MatchWithDithering finds the best match considering dithering error.
	MatchWithDithering(rgb [3]uint8, error [3]float64) (*PaletteColor, [3]float64)
	
	// SetPalette updates the palette used for matching.
	SetPalette(palette *Palette)
}

// DitherConfig holds parameters for error diffusion dithering.
type DitherConfig struct {
	Enabled   bool
	Algorithm string // "floyd-steinberg", "jarvis", "stucki", etc.
}

// RGBToLAB converts an RGB color to CIELAB color space.
func RGBToLAB(rgb [3]uint8) LABColor {
	// Convert uint8 to float64 [0,1]
	r := float64(rgb[0]) / 255.0
	g := float64(rgb[1]) / 255.0
	b := float64(rgb[2]) / 255.0
	
	// Use go-colorful for conversion
	color := colorful.Color{R: r, G: g, B: b}
	l, a, bVal := color.Lab()
	
	return LABColor{L: l, A: a, B: bVal}
}

// LABToRGB converts a CIELAB color to RGB color space.
func LABToRGB(lab LABColor) [3]uint8 {
	// Use go-colorful for conversion
	color := colorful.Lab(lab.L, lab.A, lab.B)
	
	// Clamp values to [0,1]
	r := math.Max(0, math.Min(1, color.R))
	g := math.Max(0, math.Min(1, color.G))
	b := math.Max(0, math.Min(1, color.B))
	
	return [3]uint8{
		uint8(r * 255.0),
		uint8(g * 255.0),
		uint8(b * 255.0),
	}
}

// DeltaE calculates the color difference using CIEDE2000 formula.
func DeltaE(lab1, lab2 LABColor) float64 {
	// Convert to go-colorful colors
	c1 := colorful.Lab(lab1.L, lab1.A, lab1.B)
	c2 := colorful.Lab(lab2.L, lab2.A, lab2.B)
	
	// Use CIEDE2000 distance
	return c1.DistanceCIEDE2000(c2)
}
