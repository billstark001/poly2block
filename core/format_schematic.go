package core

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/Tnze/go-mc/nbt"
)

// SchematicExporterImpl implements SchematicExporter for Minecraft schematics.
type SchematicExporterImpl struct {
	Version string
}

// NewSchematicExporter creates a new schematic exporter.
func NewSchematicExporter(version string) *SchematicExporterImpl {
	return &SchematicExporterImpl{Version: version}
}

// Export writes a voxel grid as a Minecraft schematic.
func (e *SchematicExporterImpl) Export(vg *VoxelGrid, palette *Palette, config DitherConfig, w io.Writer) error {
	// Create NBT structure for schematic
	schematic := map[string]interface{}{
		"Version":      int32(2), // Sponge Schematic version 2
		"DataVersion":  int32(2975), // Minecraft 1.19
		"Width":        int16(vg.SizeX),
		"Height":       int16(vg.SizeY),
		"Length":       int16(vg.SizeZ),
		"Offset":       []int32{0, 0, 0},
	}
	
	// Build palette mapping
	blockPalette := make(map[string]int32)
	paletteIndex := int32(0)
	
	// Default air block
	blockPalette["minecraft:air"] = paletteIndex
	paletteIndex++
	
	// Add blocks from palette
	if palette != nil {
		for _, color := range palette.Colors {
			blockID := "minecraft:white_concrete" // Default
			if id, ok := color.Metadata["block_id"].(string); ok {
				blockID = id
			}
			if _, exists := blockPalette[blockID]; !exists {
				blockPalette[blockID] = paletteIndex
				paletteIndex++
			}
		}
	} else {
		// Add a default block if no palette
		blockPalette["minecraft:white_concrete"] = paletteIndex
		paletteIndex++
	}
	
	// Convert palette map to NBT format
	paletteNBT := make(map[string]interface{})
	for blockID, idx := range blockPalette {
		paletteNBT[blockID] = idx
	}
	schematic["Palette"] = paletteNBT
	schematic["PaletteMax"] = paletteIndex
	
	// Build block data array
	blockData := make([]byte, vg.SizeX*vg.SizeY*vg.SizeZ)
	
	// Initialize with air (0)
	for i := range blockData {
		blockData[i] = 0
	}
	
	// Fill voxels
	matcher := NewCIELABMatcher(palette)
	for _, voxel := range vg.Voxels {
		// Calculate index (YZX order for Minecraft)
		index := voxel.Y + voxel.Z*vg.SizeY + voxel.X*vg.SizeY*vg.SizeZ
		
		if palette != nil {
			// Match color to palette
			matched := matcher.Match(voxel.Color)
			if matched != nil {
				if blockID, ok := matched.Metadata["block_id"].(string); ok {
					if idx, exists := blockPalette[blockID]; exists {
						blockData[index] = byte(idx)
					}
				}
			}
		} else {
			// Use default block
			blockData[index] = 1
		}
	}
	
	schematic["BlockData"] = blockData
	
	// Add metadata
	metadata := map[string]interface{}{
		"Name":   "poly2block export",
		"Author": "poly2block",
	}
	schematic["Metadata"] = metadata
	
	// Encode to NBT
	var buf bytes.Buffer
	encoder := nbt.NewEncoder(&buf)
	if err := encoder.Encode(schematic, "Schematic"); err != nil {
		return fmt.Errorf("failed to encode NBT: %w", err)
	}
	
	// Compress with gzip
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()
	
	if _, err := gzipWriter.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to compress schematic: %w", err)
	}
	
	return nil
}

// SchematicImporterImpl implements SchematicImporter for Minecraft schematics.
type SchematicImporterImpl struct{}

// NewSchematicImporter creates a new schematic importer.
func NewSchematicImporter() *SchematicImporterImpl {
	return &SchematicImporterImpl{}
}

// Import reads a schematic file and returns a voxel grid.
func (imp *SchematicImporterImpl) Import(r io.Reader) (*VoxelGrid, error) {
	// Decompress gzip
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()
	
	// Decode NBT
	var schematic map[string]interface{}
	decoder := nbt.NewDecoder(gzipReader)
	_, err = decoder.Decode(&schematic)
	if err != nil {
		return nil, fmt.Errorf("failed to decode NBT: %w", err)
	}
	
	// Extract dimensions
	width := int(schematic["Width"].(int16))
	height := int(schematic["Height"].(int16))
	length := int(schematic["Length"].(int16))
	
	// Create voxel grid
	vg := NewVoxelGrid(width, height, length)
	
	// Extract block data
	blockData := schematic["BlockData"].([]byte)
	palette := schematic["Palette"].(map[string]interface{})
	
	// Build reverse palette
	reversePalette := make(map[int32]string)
	for blockID, idx := range palette {
		reversePalette[idx.(int32)] = blockID
	}
	
	// Fill voxel grid
	for y := 0; y < height; y++ {
		for z := 0; z < length; z++ {
			for x := 0; x < width; x++ {
				index := y + z*height + x*height*length
				blockIndex := int32(blockData[index])
				
				if blockIndex > 0 { // Skip air
					// Get block ID
					if blockID, ok := reversePalette[blockIndex]; ok && blockID != "minecraft:air" {
						// Use a default color for now
						// In a full implementation, we'd look up the actual block color
						vg.SetVoxel(x, y, z, [3]uint8{128, 128, 128})
					}
				}
			}
		}
	}
	
	return vg, nil
}
