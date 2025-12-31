# Examples

This directory contains examples for using poly2block.

## Extracting Palettes from Minecraft Resources

### 1. Extract from Vanilla Minecraft JAR

```bash
# First, locate your Minecraft JAR file
# On Windows: %APPDATA%\.minecraft\versions\<version>\<version>.jar
# On macOS: ~/Library/Application Support/minecraft/versions/<version>/<version>.jar
# On Linux: ~/.minecraft/versions/<version>/<version>.jar

# Extract palette from the JAR
poly2block extract-palette \
  --jar ~/.minecraft/versions/1.19.2/1.19.2.jar \
  --output vanilla_1.19.2.msgpack \
  --export-json vanilla_1.19.2_blocks.json
```

### 2. Extract from Resource Pack

```bash
# From a ZIP file
poly2block extract-palette \
  --resource-pack path/to/resourcepack.zip \
  --output custom_palette.msgpack

# From an unpacked resource pack directory
poly2block extract-palette \
  --resource-pack path/to/resourcepack/ \
  --output custom_palette.msgpack
```

### 3. Extract from Mod JAR

```bash
# Mods with custom blocks work the same as vanilla
poly2block extract-palette \
  --jar path/to/mod.jar \
  --output mod_palette.msgpack \
  --export-json mod_blocks.json
```

### 4. Combine Multiple Sources

```bash
# Step 1: Extract from each source to JSON
poly2block extract-palette --jar vanilla.jar --export-json vanilla.json
poly2block extract-palette --jar mod1.jar --export-json mod1.json
poly2block extract-palette --resource-pack pack1.zip --export-json pack1.json

# Step 2: Manually merge JSON files (or use a script)
# Combine the arrays from all JSON files into one

# Step 3: Generate palette from combined JSON
poly2block generate-palette --custom combined.json --output complete.msgpack
```

## Using Custom Palettes for Conversion

```bash
# Extract palette from your custom resource pack
poly2block extract-palette \
  --resource-pack mypack.zip \
  --output mypack.msgpack

# Use the custom palette for conversion
poly2block convert model.gltf building.schem \
  --resolution 128 \
  --dither \
  --palette mypack.msgpack
```

## Block Definition JSON Format

The exported JSON has the following structure:

```json
[
  {
    "ID": "minecraft:stone",
    "RGB": [125, 125, 125],
    "Properties": {}
  },
  {
    "ID": "minecraft:oak_planks",
    "RGB": [162, 130, 78],
    "Properties": {}
  }
]
```

You can manually edit this file to:
- Add custom blocks
- Adjust colors
- Remove unwanted blocks
- Add block properties

## Programmatic Usage

### Extract Palette in Go Code

```go
package main

import (
	"os"
	"github.com/billstark001/poly2block/core"
)

func main() {
	// Create extractor
	extractor := core.NewTextureExtractor()
	
	// Extract from resource pack
	blocks, err := extractor.ExtractFromResourcePack("resourcepack.zip")
	if err != nil {
		panic(err)
	}
	
	// Generate palette
	palette := core.GenerateMinecraftPalette(blocks)
	
	// Save to msgpack
	f, _ := os.Create("palette.msgpack")
	defer f.Close()
	core.ExportPalette(palette, f)
}
```

### Customize Block Colors

```go
package main

import (
	"github.com/billstark001/poly2block/core"
)

func main() {
	// Load existing blocks
	blocks, _ := core.LoadBlocksFromJSON("vanilla.json")
	
	// Modify colors
	for i := range blocks {
		if blocks[i].ID == "minecraft:grass_block" {
			// Make grass brighter
			blocks[i].RGB = [3]uint8{100, 200, 100}
		}
	}
	
	// Save modified blocks
	core.SaveBlocksToJSON(blocks, "modified.json")
	
	// Generate palette
	palette := core.GenerateMinecraftPalette(blocks)
	// ... use palette
}
```

## Tips

1. **Finding Minecraft JAR**: The easiest way is to use the Minecraft Launcher:
   - Launch the game once
   - Check the versions directory
   - Each version has its own JAR file

2. **Resource Pack Priority**: If you want accurate colors:
   - Extract from your active resource pack first
   - Fall back to vanilla for missing blocks

3. **Performance**: Extracting from large resource packs may take a few seconds
   - The CLI will show progress
   - Extracted palettes can be reused

4. **Block Selection**: Not all blocks are suitable for building:
   - Some blocks are transparent
   - Some have animations
   - Consider filtering the JSON to include only solid blocks

5. **Color Accuracy**: 
   - The extractor calculates average color from texture
   - Transparent pixels are ignored
   - For best results, use high-resolution texture packs
