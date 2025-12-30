package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "poly2block",
	Short: "Convert polygon meshes to voxels and Minecraft schematics",
	Long: `poly2block is a tool for converting 3D polygon meshes (OBJ, glTF) to voxel formats
and Minecraft schematics using CIELAB color matching for accurate block selection.`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	
	// Add subcommands
	rootCmd.AddCommand(meshToVoxCmd)
	rootCmd.AddCommand(voxToSchematicCmd)
	rootCmd.AddCommand(meshToSchematicCmd)
	rootCmd.AddCommand(generatePaletteCmd)
	rootCmd.AddCommand(convertCmd)
}

// Common flags
var (
	resolution   int
	conservative bool
	ditherEnable bool
	ditherAlgo   string
	paletteFile  string
	outputFile   string
)

func addVoxelizationFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&resolution, "resolution", "r", 128, "Voxel resolution (voxels along longest axis)")
	cmd.Flags().BoolVar(&conservative, "conservative", true, "Use conservative voxelization")
}

func addDitheringFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&ditherEnable, "dither", false, "Enable error diffusion dithering")
	cmd.Flags().StringVar(&ditherAlgo, "dither-algorithm", "floyd-steinberg", "Dithering algorithm (floyd-steinberg)")
}

func addPaletteFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&paletteFile, "palette", "p", "", "Palette file (msgpack format)")
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (required)")
	cmd.MarkFlagRequired("output")
}

// printError prints an error message
func printError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
