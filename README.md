# poly2block

Convert polygon meshes to voxels and Minecraft schematics with CIELAB color matching.

## Features

- **Mesh Import**: OBJ+MTL and glTF format support
- **Voxelization**: Convert meshes to voxel grids using surface voxelization
- **CIELAB Color Matching**: Perceptually accurate color matching for block selection
- **Output Formats**: 
  - VOX (MagicaVoxel)
  - Minecraft Schematic (Sponge format)
- **Error Diffusion Dithering**: Optional Floyd-Steinberg dithering for better color reproduction
- **Palette Generation**: Generate and export CIELAB color palettes (msgpack format)

## Architecture

This project uses Go workspaces with three modules:

- **core**: Core algorithm library with generic interfaces
- **wasm**: WebAssembly bindings for web integration
- **cmd/poly2block**: Command-line interface

## Installation

### From Source

```bash
git clone https://github.com/billstark001/poly2block.git
cd poly2block
go build ./cmd/poly2block
```

### Using Go Install

```bash
go install github.com/billstark001/poly2block/cmd/poly2block@latest
```

## Usage

### CLI Examples

```bash
# Convert mesh to VOX format
poly2block mesh-to-vox input.gltf output.vox --resolution 128

# Convert mesh to Minecraft schematic
poly2block mesh-to-schematic input.obj output.schem --palette vanilla.msgpack --dither

# Generate vanilla Minecraft palette
poly2block generate-palette --output vanilla.msgpack

# Direct conversion without intermediate VOX
poly2block convert input.gltf output.schem --resolution 128 --dither
```

### Library Usage

```go
import (
    "github.com/billstark001/poly2block/core"
)

// Create pipeline
pipeline := &core.Pipeline{
    Importer:  core.NewGLTFImporter(),
    Voxelizer: core.NewSurfaceVoxelizer(),
    Matcher:   core.NewCIELABMatcher(palette),
}

// Configure conversion
config := core.PipelineConfig{
    Voxelization: core.VoxelizationConfig{
        Resolution: 128,
        Conservative: true,
    },
    Dithering: core.DitherConfig{
        Enabled: true,
        Algorithm: "floyd-steinberg",
    },
    Palette: palette,
}

// Convert
err := pipeline.MeshToSchematic(meshFile, schematicFile, config)
```

## Documentation

- [Core Library Documentation](./core/README.md)
- [CLI Documentation](./cmd/poly2block/README.md)
- [WASM Documentation](./wasm/README.md)

## Development

### Prerequisites

- Go 1.24 or later
- For WASM: TinyGo (optional, for smaller binaries)

### Building

```bash
# Build CLI
go build ./cmd/poly2block

# Build WASM
cd wasm
GOOS=js GOARCH=wasm go build -o poly2block.wasm

# Run tests
go test ./...
```

### Project Structure

```
poly2block/
├── core/               # Core algorithm library
│   ├── mesh.go        # Mesh data structures
│   ├── voxel.go       # Voxel data structures
│   ├── color.go       # Color space conversions
│   ├── pipeline.go    # Conversion pipeline
│   └── ...
├── wasm/              # WebAssembly bindings
├── cmd/poly2block/    # CLI application
└── go.work            # Go workspace configuration
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Uses [go-colorful](https://github.com/lucasb-eyer/go-colorful) for CIELAB color space conversions
- Uses [gltf](https://github.com/qmuntal/gltf) for glTF parsing
- Uses [go-mc](https://github.com/Tnze/go-mc) for Minecraft NBT format
