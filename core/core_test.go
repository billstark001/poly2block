package core

import (
	"testing"
)

func TestRGBToLAB(t *testing.T) {
	tests := []struct {
		name string
		rgb  [3]uint8
	}{
		{"White", [3]uint8{255, 255, 255}},
		{"Black", [3]uint8{0, 0, 0}},
		{"Red", [3]uint8{255, 0, 0}},
		{"Green", [3]uint8{0, 255, 0}},
		{"Blue", [3]uint8{0, 0, 255}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lab := RGBToLAB(tt.rgb)
			
			// LAB L should be in range [0, 100] (allow small negative for black due to float precision)
			if lab.L < -0.01 || lab.L > 100 {
				t.Errorf("LAB L out of range: %f", lab.L)
			}
			
			// Convert back to RGB
			rgb := LABToRGB(lab)
			
			// Allow small differences due to rounding
			for i := 0; i < 3; i++ {
				diff := int(tt.rgb[i]) - int(rgb[i])
				if diff < 0 {
					diff = -diff
				}
				if diff > 2 {
					t.Errorf("RGB conversion mismatch: expected %v, got %v", tt.rgb, rgb)
					break
				}
			}
		})
	}
}

func TestDeltaE(t *testing.T) {
	// Same colors should have zero distance
	lab1 := RGBToLAB([3]uint8{128, 128, 128})
	lab2 := RGBToLAB([3]uint8{128, 128, 128})
	
	distance := DeltaE(lab1, lab2)
	if distance > 1.0 {
		t.Errorf("Same colors should have near-zero distance, got %f", distance)
	}
	
	// Different colors should have positive distance
	lab3 := RGBToLAB([3]uint8{255, 255, 255})
	lab4 := RGBToLAB([3]uint8{0, 0, 0})
	
	distance = DeltaE(lab3, lab4)
	if distance <= 0 {
		t.Errorf("Different colors should have positive distance, got %f", distance)
	}
}

func TestPaletteGeneration(t *testing.T) {
	blocks := GetVanillaMinecraftBlocks()
	
	if len(blocks) == 0 {
		t.Fatal("No vanilla blocks returned")
	}
	
	palette := GenerateMinecraftPalette(blocks)
	
	if len(palette.Colors) != len(blocks) {
		t.Errorf("Expected %d colors, got %d", len(blocks), len(palette.Colors))
	}
	
	// Check that LAB values are populated
	for i, color := range palette.Colors {
		if color.LAB.L == 0 && color.LAB.A == 0 && color.LAB.B == 0 {
			// Only valid if RGB is also (0,0,0)
			if color.RGB != [3]uint8{0, 0, 0} {
				t.Errorf("Color %d has zero LAB but non-zero RGB", i)
			}
		}
	}
}

func TestCIELABMatcher(t *testing.T) {
	blocks := GetVanillaMinecraftBlocks()
	palette := GenerateMinecraftPalette(blocks)
	matcher := NewCIELABMatcher(palette)
	
	// Test exact match
	testColor := blocks[0].RGB
	matched := matcher.Match(testColor)
	
	if matched == nil {
		t.Fatal("Matcher returned nil")
	}
	
	// Should match the same or very similar color
	if matched.RGB != testColor {
		distance := DeltaE(RGBToLAB(testColor), matched.LAB)
		if distance > 5.0 {
			t.Errorf("Matched color too different: distance %f", distance)
		}
	}
}

func TestVoxelGrid(t *testing.T) {
	vg := NewVoxelGrid(10, 10, 10)
	
	if vg.SizeX != 10 || vg.SizeY != 10 || vg.SizeZ != 10 {
		t.Errorf("Grid size mismatch")
	}
	
	// Test setting and getting voxels
	color := [3]uint8{255, 0, 0}
	vg.SetVoxel(5, 5, 5, color)
	
	if !vg.HasVoxel(5, 5, 5) {
		t.Error("Voxel should exist at (5,5,5)")
	}
	
	voxel := vg.GetVoxel(5, 5, 5)
	if voxel == nil {
		t.Fatal("GetVoxel returned nil")
	}
	
	if voxel.Color != color {
		t.Errorf("Color mismatch: expected %v, got %v", color, voxel.Color)
	}
	
	if vg.Count() != 1 {
		t.Errorf("Expected 1 voxel, got %d", vg.Count())
	}
}

func TestMeshBounds(t *testing.T) {
	mesh := &Mesh{
		Vertices: []Vertex{
			{Position: [3]float64{0, 0, 0}},
			{Position: [3]float64{1, 1, 1}},
			{Position: [3]float64{-1, 2, 0.5}},
		},
	}
	
	mesh.CalculateBounds()
	
	expected := BoundingBox{
		Min: [3]float64{-1, 0, 0},
		Max: [3]float64{1, 2, 1},
	}
	
	if mesh.Bounds != expected {
		t.Errorf("Bounds mismatch: expected %v, got %v", expected, mesh.Bounds)
	}
}
