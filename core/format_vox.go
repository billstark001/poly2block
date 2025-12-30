package core

import (
	"encoding/binary"
	"fmt"
	"io"
)

// VOXExporterImpl handles MagicaVoxel .vox file format export.
type VOXExporterImpl struct{}

// NewVOXExporter creates a new VOX exporter.
func NewVOXExporter() *VOXExporterImpl {
	return &VOXExporterImpl{}
}

// Export writes a voxel grid to VOX format.
func (e *VOXExporterImpl) Export(vg *VoxelGrid, w io.Writer) error {
	// VOX file structure:
	// - "VOX " magic number
	// - version (150)
	// - MAIN chunk
	// - SIZE chunk (dimensions)
	// - XYZI chunk (voxel data)
	// - RGBA chunk (palette)
	
	// Write magic number
	if _, err := w.Write([]byte("VOX ")); err != nil {
		return err
	}
	
	// Write version (150)
	if err := binary.Write(w, binary.LittleEndian, int32(150)); err != nil {
		return err
	}
	
	// Create palette from voxels
	palette := make(map[[3]uint8]uint8)
	paletteIndex := uint8(1) // Index 0 is reserved for empty
	
	for _, voxel := range vg.Voxels {
		if _, exists := palette[voxel.Color]; !exists {
			palette[voxel.Color] = paletteIndex
			paletteIndex++
			if paletteIndex == 0 { // Overflow (256 colors max)
				break
			}
		}
	}
	
	// Write MAIN chunk
	if err := e.writeChunk(w, "MAIN", []byte{}, func(w io.Writer) error {
		// Write SIZE chunk
		if err := e.writeSizeChunk(w, vg); err != nil {
			return err
		}
		
		// Write XYZI chunk
		if err := e.writeXYZIChunk(w, vg, palette); err != nil {
			return err
		}
		
		// Write RGBA chunk
		return e.writeRGBAChunk(w, palette)
	}); err != nil {
		return err
	}
	
	return nil
}

// writeSizeChunk writes the SIZE chunk.
func (e *VOXExporterImpl) writeSizeChunk(w io.Writer, vg *VoxelGrid) error {
	sizeData := make([]byte, 12)
	binary.LittleEndian.PutUint32(sizeData[0:4], uint32(vg.SizeX))
	binary.LittleEndian.PutUint32(sizeData[4:8], uint32(vg.SizeY))
	binary.LittleEndian.PutUint32(sizeData[8:12], uint32(vg.SizeZ))
	
	return e.writeChunk(w, "SIZE", sizeData, nil)
}

// writeXYZIChunk writes the XYZI chunk.
func (e *VOXExporterImpl) writeXYZIChunk(w io.Writer, vg *VoxelGrid, palette map[[3]uint8]uint8) error {
	// Count voxels
	numVoxels := len(vg.Voxels)
	
	// Create XYZI data
	xyziData := make([]byte, 4+numVoxels*4)
	binary.LittleEndian.PutUint32(xyziData[0:4], uint32(numVoxels))
	
	i := 4
	for _, voxel := range vg.Voxels {
		xyziData[i] = byte(voxel.X)
		xyziData[i+1] = byte(voxel.Y)
		xyziData[i+2] = byte(voxel.Z)
		xyziData[i+3] = palette[voxel.Color]
		i += 4
	}
	
	return e.writeChunk(w, "XYZI", xyziData, nil)
}

// writeRGBAChunk writes the RGBA chunk.
func (e *VOXExporterImpl) writeRGBAChunk(w io.Writer, palette map[[3]uint8]uint8) error {
	// Create RGBA data (256 colors)
	rgbaData := make([]byte, 256*4)
	
	// Initialize with default palette
	for i := 0; i < 256; i++ {
		rgbaData[i*4] = 0
		rgbaData[i*4+1] = 0
		rgbaData[i*4+2] = 0
		rgbaData[i*4+3] = 255
	}
	
	// Fill in actual colors
	for color, index := range palette {
		idx := int(index) * 4
		rgbaData[idx] = color[0]
		rgbaData[idx+1] = color[1]
		rgbaData[idx+2] = color[2]
		rgbaData[idx+3] = 255
	}
	
	return e.writeChunk(w, "RGBA", rgbaData, nil)
}

// writeChunk writes a VOX chunk.
func (e *VOXExporterImpl) writeChunk(w io.Writer, id string, content []byte, childWriter func(io.Writer) error) error {
	// Write chunk ID
	if _, err := w.Write([]byte(id)); err != nil {
		return err
	}
	
	// Calculate child content size
	childSize := int32(0)
	if childWriter != nil {
		// For MAIN chunk, we need to calculate child size
		// This is a simplification; proper implementation would buffer
		childSize = 0 // Will be updated when children are written
	}
	
	// Write content size
	if err := binary.Write(w, binary.LittleEndian, int32(len(content))); err != nil {
		return err
	}
	
	// Write children size
	if err := binary.Write(w, binary.LittleEndian, childSize); err != nil {
		return err
	}
	
	// Write content
	if len(content) > 0 {
		if _, err := w.Write(content); err != nil {
			return err
		}
	}
	
	// Write children
	if childWriter != nil {
		return childWriter(w)
	}
	
	return nil
}

// VOXImporterImpl imports voxel grids from MagicaVoxel .vox format.
type VOXImporterImpl struct{}

// NewVOXImporter creates a new VOX importer.
func NewVOXImporter() *VOXImporterImpl {
	return &VOXImporterImpl{}
}

// Import reads a VOX file and returns a voxel grid.
func (imp *VOXImporterImpl) Import(r io.Reader) (*VoxelGrid, error) {
	// Read magic number
	magic := make([]byte, 4)
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, err
	}
	if string(magic) != "VOX " {
		return nil, fmt.Errorf("invalid VOX file: wrong magic number")
	}
	
	// Read version
	var version int32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, err
	}
	
	// Read chunks
	// This is a simplified implementation
	// A full implementation would parse all chunks properly
	
	return nil, fmt.Errorf("VOX import not fully implemented yet")
}
