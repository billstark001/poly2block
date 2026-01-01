// +build js,wasm

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/billstark001/poly2block/core"
)

func main() {
	c := make(chan struct{}, 0)
	
	// Register functions to JavaScript
	js.Global().Set("poly2block", js.ValueOf(map[string]interface{}{
		"meshToVox":       js.FuncOf(meshToVox),
		"meshToSchematic": js.FuncOf(meshToSchematic),
		"generatePalette": js.FuncOf(generatePalette),
		"version":         js.ValueOf("0.1.0"),
	}))
	
	fmt.Println("poly2block WASM module loaded")
	<-c
}

// meshToVox converts a mesh to VOX format
// Args: meshData (base64 or Uint8Array), resolution (int), conservative (bool)
// Returns: voxData (base64 string) or error
func meshToVox(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return wrapError("meshToVox requires 3 arguments: meshData, resolution, conservative")
	}
	
	// Get mesh data
	meshData, err := extractBytes(args[0])
	if err != nil {
		return wrapError(fmt.Sprintf("failed to extract mesh data: %v", err))
	}
	
	resolution := args[1].Int()
	conservative := args[2].Bool()
	
	// Create pipeline
	importer := core.NewGLTFImporter()
	voxelizer := core.NewSurfaceVoxelizer()
	
	pipeline := &core.Pipeline{
		Importer:  importer,
		Voxelizer: voxelizer,
	}
	
	config := core.PipelineConfig{
		Voxelization: core.VoxelizationConfig{
			Resolution:   resolution,
			Conservative: conservative,
		},
	}
	
	// Convert
	meshReader := bytes.NewReader(meshData)
	var voxWriter bytes.Buffer
	
	if err := pipeline.MeshToVOX(meshReader, &voxWriter, config); err != nil {
		return wrapError(fmt.Sprintf("conversion failed: %v", err))
	}
	
	// Return as base64
	result := base64.StdEncoding.EncodeToString(voxWriter.Bytes())
	return wrapSuccess(result)
}

// meshToSchematic converts a mesh to Minecraft schematic
// Args: meshData, resolution, conservative, dither, paletteData (optional)
func meshToSchematic(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 {
		return wrapError("meshToSchematic requires at least 4 arguments: meshData, resolution, conservative, dither")
	}
	
	// Get mesh data
	meshData, err := extractBytes(args[0])
	if err != nil {
		return wrapError(fmt.Sprintf("failed to extract mesh data: %v", err))
	}
	
	resolution := args[1].Int()
	conservative := args[2].Bool()
	dither := args[3].Bool()
	
	// Get palette (use vanilla if not provided)
	var palette *core.Palette
	if len(args) >= 5 && !args[4].IsNull() && !args[4].IsUndefined() {
		paletteData, err := extractBytes(args[4])
		if err != nil {
			return wrapError(fmt.Sprintf("failed to extract palette data: %v", err))
		}
		palette, err = core.ImportPalette(bytes.NewReader(paletteData))
		if err != nil {
			return wrapError(fmt.Sprintf("failed to import palette: %v", err))
		}
	} else {
		blocks := core.GetVanillaMinecraftBlocks()
		palette = core.GenerateMinecraftPalette(blocks)
	}
	
	// Create pipeline
	importer := core.NewGLTFImporter()
	voxelizer := core.NewSurfaceVoxelizer()
	matcher := core.NewCIELABMatcher(palette)
	
	pipeline := &core.Pipeline{
		Importer:  importer,
		Voxelizer: voxelizer,
		Matcher:   matcher,
	}
	
	config := core.PipelineConfig{
		Voxelization: core.VoxelizationConfig{
			Resolution:   resolution,
			Conservative: conservative,
		},
		Dithering: core.DitherConfig{
			Enabled:   dither,
			Algorithm: "floyd-steinberg",
		},
		Palette: palette,
	}
	
	// Convert
	meshReader := bytes.NewReader(meshData)
	var schematicWriter bytes.Buffer
	
	if err := pipeline.MeshToSchematic(meshReader, &schematicWriter, config); err != nil {
		return wrapError(fmt.Sprintf("conversion failed: %v", err))
	}
	
	// Return as base64
	result := base64.StdEncoding.EncodeToString(schematicWriter.Bytes())
	return wrapSuccess(result)
}

// generatePalette generates a Minecraft block palette
// Args: none (uses vanilla blocks)
// Returns: paletteData (base64 string) or error
func generatePalette(this js.Value, args []js.Value) interface{} {
	blocks := core.GetVanillaMinecraftBlocks()
	palette := core.GenerateMinecraftPalette(blocks)
	
	var buf bytes.Buffer
	if err := core.ExportPalette(palette, &buf); err != nil {
		return wrapError(fmt.Sprintf("failed to export palette: %v", err))
	}
	
	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	return wrapSuccess(result)
}

// Helper functions

func extractBytes(val js.Value) ([]byte, error) {
	if val.Type() == js.TypeString {
		// Base64 encoded string
		return base64.StdEncoding.DecodeString(val.String())
	} else if val.InstanceOf(js.Global().Get("Uint8Array")) {
		// Uint8Array
		length := val.Get("length").Int()
		data := make([]byte, length)
		js.CopyBytesToGo(data, val)
		return data, nil
	}
	return nil, fmt.Errorf("unsupported data type")
}

func wrapSuccess(data string) interface{} {
	return js.ValueOf(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func wrapError(msg string) interface{} {
	return js.ValueOf(map[string]interface{}{
		"success": false,
		"error":   msg,
	})
}
