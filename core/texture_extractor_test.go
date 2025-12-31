package core

import (
	"image"
	"image/color"
	"testing"
)

func TestCalculateAverageColor(t *testing.T) {
	te := NewTextureExtractor()
	
	// Create a simple 2x2 test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	
	// Set pixels: red, green, blue, white
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img.Set(1, 0, color.RGBA{0, 255, 0, 255})
	img.Set(0, 1, color.RGBA{0, 0, 255, 255})
	img.Set(1, 1, color.RGBA{255, 255, 255, 255})
	
	avgColor := te.calculateAverageColor(img)
	
	// Average should be roughly (127, 127, 127)
	expected := [3]uint8{127, 127, 127}
	
	// Allow some tolerance due to rounding
	for i := 0; i < 3; i++ {
		diff := int(avgColor[i]) - int(expected[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > 2 {
			t.Errorf("Average color component %d: expected ~%d, got %d", i, expected[i], avgColor[i])
		}
	}
}

func TestCalculateAverageColorWithTransparency(t *testing.T) {
	te := NewTextureExtractor()
	
	// Create a 2x2 image with some transparent pixels
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	
	// Set pixels: red, transparent, blue, transparent
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img.Set(1, 0, color.RGBA{0, 0, 0, 0}) // transparent
	img.Set(0, 1, color.RGBA{0, 0, 255, 255})
	img.Set(1, 1, color.RGBA{0, 0, 0, 0}) // transparent
	
	avgColor := te.calculateAverageColor(img)
	
	// Average should be between red and blue (ignoring transparent pixels)
	// Expected: (127, 0, 127)
	if avgColor[0] < 120 || avgColor[0] > 135 {
		t.Errorf("Red component: expected ~127, got %d", avgColor[0])
	}
	if avgColor[1] > 10 {
		t.Errorf("Green component: expected ~0, got %d", avgColor[1])
	}
	if avgColor[2] < 120 || avgColor[2] > 135 {
		t.Errorf("Blue component: expected ~127, got %d", avgColor[2])
	}
}

func TestLoadBlocksFromJSON(t *testing.T) {
	// Create a temporary JSON file
	tmpfile := "/tmp/test_blocks.json"
	
	blocks := []MinecraftBlock{
		{ID: "test:red_block", RGB: [3]uint8{255, 0, 0}, Properties: map[string]string{}},
		{ID: "test:green_block", RGB: [3]uint8{0, 255, 0}, Properties: map[string]string{}},
	}
	
	// Save to JSON
	if err := SaveBlocksToJSON(blocks, tmpfile); err != nil {
		t.Fatalf("Failed to save blocks to JSON: %v", err)
	}
	
	// Load from JSON
	loadedBlocks, err := LoadBlocksFromJSON(tmpfile)
	if err != nil {
		t.Fatalf("Failed to load blocks from JSON: %v", err)
	}
	
	if len(loadedBlocks) != len(blocks) {
		t.Errorf("Expected %d blocks, got %d", len(blocks), len(loadedBlocks))
	}
	
	for i, block := range loadedBlocks {
		if block.ID != blocks[i].ID {
			t.Errorf("Block %d: expected ID %s, got %s", i, blocks[i].ID, block.ID)
		}
		if block.RGB != blocks[i].RGB {
			t.Errorf("Block %d: expected RGB %v, got %v", i, blocks[i].RGB, block.RGB)
		}
	}
}

func TestResolveTexture(t *testing.T) {
	te := NewTextureExtractor()
	
	// Test direct texture reference
	model := BlockModel{
		Textures: map[string]string{
			"all": "block/stone",
		},
	}
	
	texture := te.resolveTexture(model)
	if texture != "block/stone" {
		t.Errorf("Expected 'block/stone', got '%s'", texture)
	}
	
	// Test texture variable reference
	model2 := BlockModel{
		Textures: map[string]string{
			"all":  "#base",
			"base": "block/wood",
		},
	}
	
	texture = te.resolveTexture(model2)
	if texture != "block/wood" {
		t.Errorf("Expected 'block/wood', got '%s'", texture)
	}
}
