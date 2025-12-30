package core

import "io"

// VOXFormat handles MagicaVoxel .vox file format.
type VOXFormat struct{}

// VOXExporter is the interface for exporting voxel grids to VOX format.
type VOXExporter interface {
	// Export writes a voxel grid to VOX format.
	Export(vg *VoxelGrid, w io.Writer) error
}

// VOXImporter is the interface for importing VOX files.
type VOXImporter interface {
	// Import reads a VOX file and returns a voxel grid.
	Import(r io.Reader) (*VoxelGrid, error)
}

// SchematicFormat handles Minecraft schematic format.
type SchematicFormat struct {
	Version string // "1.13+", "1.12" for different Minecraft versions
}

// MinecraftBlock represents a Minecraft block with its properties.
type MinecraftBlock struct {
	ID         string
	Properties map[string]string
	RGB        [3]uint8
	LAB        LABColor
}

// SchematicExporter is the interface for exporting to Minecraft schematic format.
type SchematicExporter interface {
	// Export writes a voxel grid as a Minecraft schematic.
	Export(vg *VoxelGrid, palette *Palette, config DitherConfig, w io.Writer) error
}

// SchematicImporter is the interface for importing Minecraft schematics.
type SchematicImporter interface {
	// Import reads a schematic file and returns a voxel grid.
	Import(r io.Reader) (*VoxelGrid, error)
}
