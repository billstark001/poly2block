# poly2block

Convert polygon meshes to voxels and Minecraft schematics with CIELAB color matching.

[![Build Status](https://github.com/billstark001/poly2block/workflows/Build%20and%20Test/badge.svg)](https://github.com/billstark001/poly2block/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/billstark001/poly2block)](https://goreportcard.com/report/github.com/billstark001/poly2block)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Mesh Import**: glTF format support (OBJ planned)
- **Voxelization**: Surface voxelization algorithm with conservative mode
- **CIELAB Color Matching**: Perceptually accurate color matching using CIEDE2000
- **Output Formats**: VOX (MagicaVoxel) and Minecraft Schematic (Sponge v2)
- **Error Diffusion Dithering**: Floyd-Steinberg dithering for better color reproduction
- **Palette Generation**: Generate and export CIELAB color palettes (msgpack format)
- **Multiple Interfaces**: CLI, Go library, and WebAssembly

## Quick Start

### Installation

```bash
go install github.com/billstark001/poly2block/cmd/poly2block@latest
```

### Basic Usage

```bash
# Convert glTF to Minecraft schematic
poly2block convert model.gltf output.schem --resolution 128 --dither

# Generate vanilla Minecraft palette
poly2block generate-palette --output vanilla.msgpack
```

## Documentation

- [Core Library Documentation](./core/README.md)
- [CLI Documentation](./cmd/poly2block/README.md)
- [WASM Documentation](./wasm/README.md)

## Development

See individual module READMEs for detailed documentation.

## License

MIT License - see [LICENSE](LICENSE) for details.
