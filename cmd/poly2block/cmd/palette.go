package cmd

import (
	"fmt"
	"os"

	"github.com/billstark001/poly2block/core"
	"github.com/spf13/cobra"
)

var (
	vanillaBlocks   bool
	customBlocks    string
	resourcePack    string
	jarFile         string
	exportJSON      string
)

var generatePaletteCmd = &cobra.Command{
	Use:   "generate-palette",
	Short: "Generate CIELAB color palette for Minecraft blocks",
	Long: `Generate a CIELAB color space palette file for Minecraft blocks.
The palette can be used for color matching when converting meshes to schematics.`,
	RunE: runGeneratePalette,
}

var extractPaletteCmd = &cobra.Command{
	Use:   "extract-palette",
	Short: "Extract palette from Minecraft resource pack or jar",
	Long: `Extract block colors from Minecraft resource pack (zip or directory) or jar file.
This analyzes textures and generates accurate color information.`,
	RunE: runExtractPalette,
}

func init() {
	generatePaletteCmd.Flags().StringVarP(&outputFile, "output", "o", "palette.msgpack", "Output palette file")
	generatePaletteCmd.Flags().BoolVar(&vanillaBlocks, "vanilla", true, "Include vanilla Minecraft blocks")
	generatePaletteCmd.Flags().StringVar(&customBlocks, "custom", "", "Custom blocks definition file (JSON)")
	
	extractPaletteCmd.Flags().StringVarP(&outputFile, "output", "o", "palette.msgpack", "Output palette file")
	extractPaletteCmd.Flags().StringVar(&resourcePack, "resource-pack", "", "Path to resource pack (zip or directory)")
	extractPaletteCmd.Flags().StringVar(&jarFile, "jar", "", "Path to Minecraft jar file")
	extractPaletteCmd.Flags().StringVar(&exportJSON, "export-json", "", "Also export blocks as JSON")
}

func runGeneratePalette(cmd *cobra.Command, args []string) error {
	fmt.Println("Generating Minecraft block palette...")
	
	var blocks []core.MinecraftBlock
	
	if vanillaBlocks {
		fmt.Println("Including vanilla Minecraft blocks")
		blocks = append(blocks, core.GetVanillaMinecraftBlocks()...)
	}
	
	if customBlocks != "" {
		fmt.Printf("Loading custom blocks from %s\n", customBlocks)
		customBlocksList, err := core.LoadBlocksFromJSON(customBlocks)
		if err != nil {
			return fmt.Errorf("failed to load custom blocks: %w", err)
		}
		blocks = append(blocks, customBlocksList...)
	}
	
	if len(blocks) == 0 {
		return fmt.Errorf("no blocks specified")
	}
	
	// Generate palette
	palette := core.GenerateMinecraftPalette(blocks)
	
	// Export to file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()
	
	if err := core.ExportPalette(palette, outFile); err != nil {
		return fmt.Errorf("failed to export palette: %w", err)
	}
	
	fmt.Printf("Successfully generated palette with %d colors\n", len(palette.Colors))
	fmt.Printf("Saved to %s\n", outputFile)
	
	return nil
}

func runExtractPalette(cmd *cobra.Command, args []string) error {
	if resourcePack == "" && jarFile == "" {
		return fmt.Errorf("must specify either --resource-pack or --jar")
	}
	
	extractor := core.NewTextureExtractor()
	var blocks []core.MinecraftBlock
	var err error
	
	if resourcePack != "" {
		fmt.Printf("Extracting blocks from resource pack: %s\n", resourcePack)
		blocks, err = extractor.ExtractFromResourcePack(resourcePack)
		if err != nil {
			return fmt.Errorf("failed to extract from resource pack: %w", err)
		}
	} else if jarFile != "" {
		fmt.Printf("Extracting blocks from jar file: %s\n", jarFile)
		blocks, err = extractor.ExtractFromJar(jarFile)
		if err != nil {
			return fmt.Errorf("failed to extract from jar: %w", err)
		}
	}
	
	if len(blocks) == 0 {
		return fmt.Errorf("no blocks found in the resource pack/jar")
	}
	
	fmt.Printf("Found %d blocks with textures\n", len(blocks))
	
	// Export as JSON if requested
	if exportJSON != "" {
		fmt.Printf("Exporting blocks to JSON: %s\n", exportJSON)
		if err := core.SaveBlocksToJSON(blocks, exportJSON); err != nil {
			return fmt.Errorf("failed to export JSON: %w", err)
		}
	}
	
	// Generate palette
	palette := core.GenerateMinecraftPalette(blocks)
	
	// Export to file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()
	
	if err := core.ExportPalette(palette, outFile); err != nil {
		return fmt.Errorf("failed to export palette: %w", err)
	}
	
	fmt.Printf("Successfully generated palette with %d colors\n", len(palette.Colors))
	fmt.Printf("Saved to %s\n", outputFile)
	
	return nil
}
