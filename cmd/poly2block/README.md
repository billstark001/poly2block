# poly2block CLI

Command-line tool for converting polygon meshes to voxels and Minecraft schematics.

## Installation

```bash
go install github.com/billstark001/poly2block/cmd/poly2block@latest
```

Or build from source:

```bash
cd cmd/poly2block
go build
```

## Commands

### mesh-to-vox

Convert a polygon mesh to MagicaVoxel VOX format.

```bash
poly2block mesh-to-vox input.gltf output.vox --resolution 128
```

Options:
- `-r, --resolution`: Voxel resolution (default: 128)
- `--conservative`: Use conservative voxelization (default: true)

### mesh-to-schematic

Convert a polygon mesh directly to Minecraft schematic.

```bash
poly2block mesh-to-schematic input.gltf output.schem \
  --resolution 128 \
  --dither \
  --palette vanilla.msgpack
```

Options:
- `-r, --resolution`: Voxel resolution (default: 128)
- `--conservative`: Use conservative voxelization (default: true)
- `--dither`: Enable error diffusion dithering
- `--dither-algorithm`: Dithering algorithm (default: floyd-steinberg)
- `-p, --palette`: Palette file path (msgpack format)

### vox-to-schematic

Convert a VOX file to Minecraft schematic.

```bash
poly2block vox-to-schematic input.vox output.schem \
  --dither \
  --palette vanilla.msgpack
```

Options:
- `--dither`: Enable error diffusion dithering
- `--dither-algorithm`: Dithering algorithm (default: floyd-steinberg)
- `-p, --palette`: Palette file path (msgpack format)

### generate-palette

Generate a CIELAB color palette for Minecraft blocks.

```bash
poly2block generate-palette --output vanilla.msgpack
```

Options:
- `-o, --output`: Output file path (default: palette.msgpack)
- `--vanilla`: Include vanilla Minecraft blocks (default: true)
- `--custom`: Custom blocks definition file (JSON)

### extract-palette

Extract block colors from Minecraft resource pack or jar file by analyzing textures.

```bash
# Extract from resource pack (zip or directory)
poly2block extract-palette \
  --resource-pack path/to/resourcepack.zip \
  --output custom.msgpack \
  --export-json blocks.json

# Extract from Minecraft jar file
poly2block extract-palette \
  --jar path/to/minecraft.jar \
  --output vanilla.msgpack
```

Options:
- `-o, --output`: Output palette file (default: palette.msgpack)
- `--resource-pack`: Path to resource pack (zip or directory)
- `--jar`: Path to Minecraft jar file
- `--export-json`: Also export blocks as JSON file

### convert

Alias for `mesh-to-schematic`.

```bash
poly2block convert input.gltf output.schem --resolution 128 --dither
```

## Examples

### Basic Conversion

```bash
# Convert glTF to schematic with default settings
poly2block convert model.gltf building.schem

# Higher resolution
poly2block convert model.gltf building.schem --resolution 256
```

### With Custom Palette

```bash
# Generate custom palette
poly2block generate-palette --output my_palette.msgpack

# Use custom palette
poly2block convert model.gltf building.schem \
  --palette my_palette.msgpack \
  --dither
```

### Two-Stage Workflow

```bash
# Step 1: Convert to VOX
poly2block mesh-to-vox model.gltf model.vox --resolution 128

# Step 2: Convert VOX to schematic
poly2block vox-to-schematic model.vox building.schem \
  --palette vanilla.msgpack \
  --dither
```

## Supported Formats

### Input Formats
- glTF (.gltf, .glb)
- OBJ (.obj) - Coming soon

### Output Formats
- VOX (.vox) - MagicaVoxel format
- Schematic (.schem, .schematic) - Minecraft Sponge format

## Performance Tips

1. Start with lower resolutions (64-128) for testing
2. Enable dithering for better color reproduction
3. Use custom palettes for specific resource packs
4. Conservative voxelization is slower but more accurate

## License

MIT
