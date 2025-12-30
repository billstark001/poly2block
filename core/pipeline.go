package core

import "io"

// Pipeline represents the complete conversion pipeline.
type Pipeline struct {
	Importer  MeshImporter
	Voxelizer Voxelizer
	Matcher   ColorMatcher
}

// PipelineConfig holds all configuration for the conversion pipeline.
type PipelineConfig struct {
	Voxelization VoxelizationConfig
	Dithering    DitherConfig
	Palette      *Palette
}

// MeshToVoxelGrid converts a mesh directly to a voxel grid.
func (p *Pipeline) MeshToVoxelGrid(meshReader io.Reader, config PipelineConfig) (*VoxelGrid, error) {
	// Import mesh
	mesh, err := p.Importer.Import(meshReader)
	if err != nil {
		return nil, err
	}
	
	// Voxelize
	voxelGrid, err := p.Voxelizer.Voxelize(mesh, config.Voxelization)
	if err != nil {
		return nil, err
	}
	
	return voxelGrid, nil
}

// MeshToVOX converts a mesh to VOX format.
func (p *Pipeline) MeshToVOX(meshReader io.Reader, voxWriter io.Writer, config PipelineConfig) error {
	voxelGrid, err := p.MeshToVoxelGrid(meshReader, config)
	if err != nil {
		return err
	}
	
	exporter := NewVOXExporter()
	return exporter.Export(voxelGrid, voxWriter)
}

// VoxelGridToSchematic converts a voxel grid to Minecraft schematic.
func (p *Pipeline) VoxelGridToSchematic(vg *VoxelGrid, schematicWriter io.Writer, config PipelineConfig) error {
	// Apply color matching and dithering
	if config.Palette != nil && p.Matcher != nil {
		p.Matcher.SetPalette(config.Palette)
		
		// Apply dithering if enabled
		if config.Dithering.Enabled {
			vg = p.applyDithering(vg, config.Dithering)
		} else {
			// Simple color matching without dithering
			vg = p.applyColorMatching(vg)
		}
	}
	
	// Export to schematic
	exporter := NewSchematicExporter("1.13+")
	return exporter.Export(vg, config.Palette, config.Dithering, schematicWriter)
}

// MeshToSchematic converts a mesh directly to Minecraft schematic.
func (p *Pipeline) MeshToSchematic(meshReader io.Reader, schematicWriter io.Writer, config PipelineConfig) error {
	voxelGrid, err := p.MeshToVoxelGrid(meshReader, config)
	if err != nil {
		return err
	}
	
	return p.VoxelGridToSchematic(voxelGrid, schematicWriter, config)
}

// applyColorMatching applies color matching without dithering.
func (p *Pipeline) applyColorMatching(vg *VoxelGrid) *VoxelGrid {
	result := NewVoxelGrid(vg.SizeX, vg.SizeY, vg.SizeZ)
	result.Scale = vg.Scale
	result.Origin = vg.Origin
	
	for pos, voxel := range vg.Voxels {
		matched := p.Matcher.Match(voxel.Color)
		if matched != nil {
			result.SetVoxel(pos[0], pos[1], pos[2], matched.RGB)
		}
	}
	
	return result
}

// applyDithering applies error diffusion dithering during color matching.
func (p *Pipeline) applyDithering(vg *VoxelGrid, config DitherConfig) *VoxelGrid {
	result := NewVoxelGrid(vg.SizeX, vg.SizeY, vg.SizeZ)
	result.Scale = vg.Scale
	result.Origin = vg.Origin
	
	// Error buffer for dithering
	errorBuffer := make(map[[3]int][3]float64)
	
	// Process voxels in order (for error diffusion)
	for z := 0; z < vg.SizeZ; z++ {
		for y := 0; y < vg.SizeY; y++ {
			for x := 0; x < vg.SizeX; x++ {
				voxel := vg.GetVoxel(x, y, z)
				if voxel == nil {
					continue
				}
				
				pos := [3]int{x, y, z}
				error := errorBuffer[pos]
				
				matched, quantError := p.Matcher.MatchWithDithering(voxel.Color, error)
				if matched != nil {
					result.SetVoxel(x, y, z, matched.RGB)
					
					// Distribute error to neighbors (Floyd-Steinberg pattern)
					p.distributeError(errorBuffer, x, y, z, quantError, config.Algorithm)
				}
			}
		}
	}
	
	return result
}

// distributeError distributes quantization error to neighboring voxels.
func (p *Pipeline) distributeError(buffer map[[3]int][3]float64, x, y, z int, error [3]float64, algorithm string) {
	// Floyd-Steinberg coefficients
	if algorithm == "floyd-steinberg" || algorithm == "" {
		p.addError(buffer, x+1, y, z, error, 7.0/16.0)
		p.addError(buffer, x-1, y+1, z, error, 3.0/16.0)
		p.addError(buffer, x, y+1, z, error, 5.0/16.0)
		p.addError(buffer, x+1, y+1, z, error, 1.0/16.0)
	}
	// Other algorithms can be added here
}

// addError adds error to the buffer at the given position.
func (p *Pipeline) addError(buffer map[[3]int][3]float64, x, y, z int, error [3]float64, weight float64) {
	pos := [3]int{x, y, z}
	current := buffer[pos]
	for i := 0; i < 3; i++ {
		current[i] += error[i] * weight
	}
	buffer[pos] = current
}
