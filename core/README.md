# poly2block/core

Core algorithm library for converting polygon meshes to voxels and Minecraft schematics.

## Features

- **Generic Interfaces**: Pluggable implementations for mesh import, voxelization, and color matching
- **Multiple Input Formats**: Support for OBJ+MTL and glTF
- **Voxelization**: Configurable voxelization with multiple algorithms
- **CIELAB Color Matching**: Perceptually accurate color matching using CIELAB color space
- **Output Formats**: VOX (MagicaVoxel) and Minecraft schematic formats
- **Error Diffusion Dithering**: Optional Floyd-Steinberg and other dithering algorithms
- **Palette Generation**: Generate CIELAB color palettes for Minecraft blocks (msgpack format)
- **Texture Extraction**: Extract block colors from Minecraft resource packs and jar files

## Architecture

The library is built around generic interfaces to allow algorithm swapping:

- `MeshImporter`: Import polygon meshes from various formats
- `Voxelizer`: Convert meshes to voxel grids
- `ColorMatcher`: Match colors to predefined palettes using CIELAB
- `VOXExporter/Importer`: Handle MagicaVoxel format
- `SchematicExporter/Importer`: Handle Minecraft schematic format

## Usage

### Basic Pipeline

```go
import "github.com/billstark001/poly2block/core"

// Create pipeline
pipeline := &core.Pipeline{
    Importer:  core.NewGLTFImporter(),
    Voxelizer: core.NewSurfaceVoxelizer(),
    Matcher:   core.NewCIELABMatcher(palette),
}

// Configure
config := core.PipelineConfig{
    Voxelization: core.VoxelizationConfig{
        Resolution: 128,
        Conservative: true,
    },
    Dithering: core.DitherConfig{
        Enabled: true,
        Algorithm: "floyd-steinberg",
    },
    Palette: myPalette,
}

// Convert mesh to schematic
err := pipeline.MeshToSchematic(meshReader, schematicWriter, config)
```

### Extracting Palettes from Resource Packs

```go
import "github.com/billstark001/poly2block/core"

// Create texture extractor
extractor := core.NewTextureExtractor()

// Extract from resource pack (zip or directory)
blocks, err := extractor.ExtractFromResourcePack("path/to/resourcepack.zip")
if err != nil {
    panic(err)
}

// Generate palette
palette := core.GenerateMinecraftPalette(blocks)

// Export to msgpack
f, _ := os.Create("custom_palette.msgpack")
defer f.Close()
core.ExportPalette(palette, f)
```

### Working with Custom Block Definitions

```go
// Load blocks from JSON
blocks, err := core.LoadBlocksFromJSON("blocks.json")

// Generate palette
palette := core.GenerateMinecraftPalette(blocks)

// Save blocks back to JSON
core.SaveBlocksToJSON(blocks, "modified_blocks.json")
```

## License

MIT
