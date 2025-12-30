package cmd

import (
	"fmt"
	"os"

	"github.com/billstark001/poly2block/core"
	"github.com/spf13/cobra"
)

var (
	vanillaBlocks bool
	customBlocks  string
)

var generatePaletteCmd = &cobra.Command{
	Use:   "generate-palette",
	Short: "Generate CIELAB color palette for Minecraft blocks",
	Long: `Generate a CIELAB color space palette file for Minecraft blocks.
The palette can be used for color matching when converting meshes to schematics.`,
	RunE: runGeneratePalette,
}

func init() {
	generatePaletteCmd.Flags().StringVarP(&outputFile, "output", "o", "palette.msgpack", "Output palette file")
	generatePaletteCmd.Flags().BoolVar(&vanillaBlocks, "vanilla", true, "Include vanilla Minecraft blocks")
	generatePaletteCmd.Flags().StringVar(&customBlocks, "custom", "", "Custom blocks definition file (JSON)")
}

func runGeneratePalette(cmd *cobra.Command, args []string) error {
	fmt.Println("Generating Minecraft block palette...")
	
	var blocks []core.MinecraftBlock
	
	if vanillaBlocks {
		fmt.Println("Including vanilla Minecraft blocks")
		blocks = append(blocks, core.GetVanillaMinecraftBlocks()...)
	}
	
	if customBlocks != "" {
		// TODO: Load custom blocks from JSON file
		fmt.Printf("Loading custom blocks from %s\n", customBlocks)
		return fmt.Errorf("custom blocks not yet implemented")
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
