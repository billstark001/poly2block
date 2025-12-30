package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/billstark001/poly2block/core"
	"github.com/spf13/cobra"
)

var meshToVoxCmd = &cobra.Command{
	Use:   "mesh-to-vox <input> <output>",
	Short: "Convert mesh to VOX format",
	Long:  `Convert a polygon mesh (OBJ, glTF) to MagicaVoxel VOX format.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runMeshToVox,
}

var voxToSchematicCmd = &cobra.Command{
	Use:   "vox-to-schematic <input> <output>",
	Short: "Convert VOX to Minecraft schematic",
	Long:  `Convert a MagicaVoxel VOX file to Minecraft schematic format.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runVoxToSchematic,
}

var meshToSchematicCmd = &cobra.Command{
	Use:   "mesh-to-schematic <input> <output>",
	Short: "Convert mesh to Minecraft schematic",
	Long:  `Convert a polygon mesh (OBJ, glTF) directly to Minecraft schematic format.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runMeshToSchematic,
}

var convertCmd = &cobra.Command{
	Use:   "convert <input> <output>",
	Short: "Convert mesh to schematic (alias)",
	Long:  `Convert a polygon mesh to Minecraft schematic (same as mesh-to-schematic).`,
	Args:  cobra.ExactArgs(2),
	RunE:  runMeshToSchematic,
}

func init() {
	// mesh-to-vox flags
	addVoxelizationFlags(meshToVoxCmd)
	
	// vox-to-schematic flags
	addDitheringFlags(voxToSchematicCmd)
	addPaletteFlags(voxToSchematicCmd)
	
	// mesh-to-schematic flags
	addVoxelizationFlags(meshToSchematicCmd)
	addDitheringFlags(meshToSchematicCmd)
	addPaletteFlags(meshToSchematicCmd)
	
	// convert flags (same as mesh-to-schematic)
	addVoxelizationFlags(convertCmd)
	addDitheringFlags(convertCmd)
	addPaletteFlags(convertCmd)
}

func runMeshToVox(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	outputFile := args[1]
	
	fmt.Printf("Converting %s to VOX format...\n", inputFile)
	
	// Open input file
	meshReader, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer meshReader.Close()
	
	// Create output file
	voxWriter, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer voxWriter.Close()
	
	// Determine importer based on file extension
	importer, err := getImporter(inputFile)
	if err != nil {
		return err
	}
	
	// Create pipeline
	pipeline := &core.Pipeline{
		Importer:  importer,
		Voxelizer: core.NewSurfaceVoxelizer(),
	}
	
	// Configure
	config := core.PipelineConfig{
		Voxelization: core.VoxelizationConfig{
			Resolution:   resolution,
			Conservative: conservative,
		},
	}
	
	// Convert
	if err := pipeline.MeshToVOX(meshReader, voxWriter, config); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	
	fmt.Printf("Successfully converted to %s\n", outputFile)
	return nil
}

func runVoxToSchematic(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	outputFile := args[1]
	
	fmt.Printf("Converting %s to Minecraft schematic...\n", inputFile)
	
	// Load palette
	palette, err := loadPalette()
	if err != nil {
		return err
	}
	
	// Open input file
	voxReader, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer voxReader.Close()
	
	// Import VOX
	voxImporter := core.NewVOXImporter()
	voxelGrid, err := voxImporter.Import(voxReader)
	if err != nil {
		return fmt.Errorf("failed to import VOX file: %w", err)
	}
	
	// Create output file
	schematicWriter, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer schematicWriter.Close()
	
	// Create pipeline
	pipeline := &core.Pipeline{
		Matcher: core.NewCIELABMatcher(palette),
	}
	
	// Configure
	config := core.PipelineConfig{
		Dithering: core.DitherConfig{
			Enabled:   ditherEnable,
			Algorithm: ditherAlgo,
		},
		Palette: palette,
	}
	
	// Convert
	if err := pipeline.VoxelGridToSchematic(voxelGrid, schematicWriter, config); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	
	fmt.Printf("Successfully converted to %s\n", outputFile)
	return nil
}

func runMeshToSchematic(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	outputFile := args[1]
	
	fmt.Printf("Converting %s to Minecraft schematic...\n", inputFile)
	
	// Load palette
	palette, err := loadPalette()
	if err != nil {
		return err
	}
	
	// Open input file
	meshReader, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer meshReader.Close()
	
	// Create output file
	schematicWriter, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer schematicWriter.Close()
	
	// Determine importer
	importer, err := getImporter(inputFile)
	if err != nil {
		return err
	}
	
	// Create pipeline
	pipeline := &core.Pipeline{
		Importer:  importer,
		Voxelizer: core.NewSurfaceVoxelizer(),
		Matcher:   core.NewCIELABMatcher(palette),
	}
	
	// Configure
	config := core.PipelineConfig{
		Voxelization: core.VoxelizationConfig{
			Resolution:   resolution,
			Conservative: conservative,
		},
		Dithering: core.DitherConfig{
			Enabled:   ditherEnable,
			Algorithm: ditherAlgo,
		},
		Palette: palette,
	}
	
	// Convert
	if err := pipeline.MeshToSchematic(meshReader, schematicWriter, config); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	
	fmt.Printf("Successfully converted to %s\n", outputFile)
	return nil
}

func getImporter(filename string) (core.MeshImporter, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".gltf", ".glb":
		return core.NewGLTFImporter(), nil
	case ".obj":
		return nil, fmt.Errorf("OBJ importer not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func loadPalette() (*core.Palette, error) {
	if paletteFile == "" {
		// Use default vanilla palette
		fmt.Println("Using default vanilla Minecraft palette")
		blocks := core.GetVanillaMinecraftBlocks()
		return core.GenerateMinecraftPalette(blocks), nil
	}
	
	// Load from file
	fmt.Printf("Loading palette from %s\n", paletteFile)
	f, err := os.Open(paletteFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open palette file: %w", err)
	}
	defer f.Close()
	
	palette, err := core.ImportPalette(f)
	if err != nil {
		return nil, fmt.Errorf("failed to import palette: %w", err)
	}
	
	return palette, nil
}
