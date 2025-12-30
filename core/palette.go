package core

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// PaletteData represents serializable palette data for msgpack.
type PaletteData struct {
	Version string                   `msgpack:"version"`
	Colors  []PaletteColorData       `msgpack:"colors"`
}

// PaletteColorData represents serializable color data.
type PaletteColorData struct {
	Name     string                 `msgpack:"name"`
	RGB      [3]uint8               `msgpack:"rgb"`
	LAB      [3]float64             `msgpack:"lab"`
	Metadata map[string]interface{} `msgpack:"metadata,omitempty"`
}

// ExportPalette exports a palette to msgpack format.
func ExportPalette(palette *Palette, w io.Writer) error {
	data := PaletteData{
		Version: "1.0",
		Colors:  make([]PaletteColorData, len(palette.Colors)),
	}
	
	for i, color := range palette.Colors {
		data.Colors[i] = PaletteColorData{
			Name:     color.Name,
			RGB:      color.RGB,
			LAB:      [3]float64{color.LAB.L, color.LAB.A, color.LAB.B},
			Metadata: color.Metadata,
		}
	}
	
	encoder := msgpack.NewEncoder(w)
	return encoder.Encode(&data)
}

// ImportPalette imports a palette from msgpack format.
func ImportPalette(r io.Reader) (*Palette, error) {
	var data PaletteData
	decoder := msgpack.NewDecoder(r)
	
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	
	palette := &Palette{
		Colors: make([]PaletteColor, len(data.Colors)),
	}
	
	for i, colorData := range data.Colors {
		palette.Colors[i] = PaletteColor{
			Name:     colorData.Name,
			RGB:      colorData.RGB,
			LAB:      LABColor{L: colorData.LAB[0], A: colorData.LAB[1], B: colorData.LAB[2]},
			Metadata: colorData.Metadata,
		}
	}
	
	return palette, nil
}

// GenerateMinecraftPalette creates a palette from Minecraft block definitions.
func GenerateMinecraftPalette(blocks []MinecraftBlock) *Palette {
	palette := &Palette{
		Colors: make([]PaletteColor, len(blocks)),
	}
	
	for i, block := range blocks {
		palette.Colors[i] = PaletteColor{
			Name: block.ID,
			RGB:  block.RGB,
			LAB:  RGBToLAB(block.RGB),
			Metadata: map[string]interface{}{
				"block_id":   block.ID,
				"properties": block.Properties,
			},
		}
	}
	
	return palette
}

// GetVanillaMinecraftBlocks returns a list of common vanilla Minecraft blocks with colors.
// This is a basic set; users can extend or customize this.
func GetVanillaMinecraftBlocks() []MinecraftBlock {
	return []MinecraftBlock{
		{ID: "minecraft:white_wool", RGB: [3]uint8{233, 236, 236}, Properties: map[string]string{}},
		{ID: "minecraft:orange_wool", RGB: [3]uint8{240, 118, 19}, Properties: map[string]string{}},
		{ID: "minecraft:magenta_wool", RGB: [3]uint8{189, 68, 179}, Properties: map[string]string{}},
		{ID: "minecraft:light_blue_wool", RGB: [3]uint8{58, 175, 217}, Properties: map[string]string{}},
		{ID: "minecraft:yellow_wool", RGB: [3]uint8{253, 221, 70}, Properties: map[string]string{}},
		{ID: "minecraft:lime_wool", RGB: [3]uint8{112, 185, 25}, Properties: map[string]string{}},
		{ID: "minecraft:pink_wool", RGB: [3]uint8{237, 141, 172}, Properties: map[string]string{}},
		{ID: "minecraft:gray_wool", RGB: [3]uint8{62, 68, 71}, Properties: map[string]string{}},
		{ID: "minecraft:light_gray_wool", RGB: [3]uint8{142, 142, 134}, Properties: map[string]string{}},
		{ID: "minecraft:cyan_wool", RGB: [3]uint8{21, 137, 145}, Properties: map[string]string{}},
		{ID: "minecraft:purple_wool", RGB: [3]uint8{121, 42, 172}, Properties: map[string]string{}},
		{ID: "minecraft:blue_wool", RGB: [3]uint8{53, 57, 157}, Properties: map[string]string{}},
		{ID: "minecraft:brown_wool", RGB: [3]uint8{114, 71, 40}, Properties: map[string]string{}},
		{ID: "minecraft:green_wool", RGB: [3]uint8{85, 109, 27}, Properties: map[string]string{}},
		{ID: "minecraft:red_wool", RGB: [3]uint8{160, 39, 34}, Properties: map[string]string{}},
		{ID: "minecraft:black_wool", RGB: [3]uint8{20, 21, 25}, Properties: map[string]string{}},
		// Concrete blocks
		{ID: "minecraft:white_concrete", RGB: [3]uint8{207, 213, 214}, Properties: map[string]string{}},
		{ID: "minecraft:orange_concrete", RGB: [3]uint8{224, 97, 1}, Properties: map[string]string{}},
		{ID: "minecraft:magenta_concrete", RGB: [3]uint8{169, 48, 159}, Properties: map[string]string{}},
		{ID: "minecraft:light_blue_concrete", RGB: [3]uint8{36, 137, 199}, Properties: map[string]string{}},
		{ID: "minecraft:yellow_concrete", RGB: [3]uint8{240, 175, 21}, Properties: map[string]string{}},
		{ID: "minecraft:lime_concrete", RGB: [3]uint8{94, 168, 24}, Properties: map[string]string{}},
		{ID: "minecraft:pink_concrete", RGB: [3]uint8{213, 101, 143}, Properties: map[string]string{}},
		{ID: "minecraft:gray_concrete", RGB: [3]uint8{54, 57, 61}, Properties: map[string]string{}},
		{ID: "minecraft:light_gray_concrete", RGB: [3]uint8{125, 125, 115}, Properties: map[string]string{}},
		{ID: "minecraft:cyan_concrete", RGB: [3]uint8{21, 119, 136}, Properties: map[string]string{}},
		{ID: "minecraft:purple_concrete", RGB: [3]uint8{100, 32, 156}, Properties: map[string]string{}},
		{ID: "minecraft:blue_concrete", RGB: [3]uint8{44, 46, 143}, Properties: map[string]string{}},
		{ID: "minecraft:brown_concrete", RGB: [3]uint8{96, 59, 31}, Properties: map[string]string{}},
		{ID: "minecraft:green_concrete", RGB: [3]uint8{73, 91, 36}, Properties: map[string]string{}},
		{ID: "minecraft:red_concrete", RGB: [3]uint8{142, 32, 32}, Properties: map[string]string{}},
		{ID: "minecraft:black_concrete", RGB: [3]uint8{8, 10, 15}, Properties: map[string]string{}},
	}
}
