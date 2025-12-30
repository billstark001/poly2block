package core

// Voxel represents a single voxel with position and color.
type Voxel struct {
	X, Y, Z int
	Color   [3]uint8 // RGB [0,255]
}

// VoxelGrid represents a 3D grid of voxels.
type VoxelGrid struct {
	SizeX, SizeY, SizeZ int
	Voxels              map[[3]int]*Voxel // Sparse representation
	Scale               float64           // Scale factor from mesh units to voxels
	Origin              [3]float64        // Origin in mesh space
}

// VoxelizationConfig holds parameters for voxelization.
type VoxelizationConfig struct {
	Resolution   int     // Target resolution (voxels along longest axis)
	Scale        float64 // Manual scale override (0 = auto)
	Conservative bool    // Use conservative voxelization
}

// Voxelizer is the interface for converting meshes to voxels.
type Voxelizer interface {
	// Voxelize converts a mesh to a voxel grid.
	Voxelize(mesh *Mesh, config VoxelizationConfig) (*VoxelGrid, error)
	
	// Name returns the algorithm name.
	Name() string
}

// NewVoxelGrid creates a new empty voxel grid.
func NewVoxelGrid(sizeX, sizeY, sizeZ int) *VoxelGrid {
	return &VoxelGrid{
		SizeX:  sizeX,
		SizeY:  sizeY,
		SizeZ:  sizeZ,
		Voxels: make(map[[3]int]*Voxel),
		Scale:  1.0,
	}
}

// SetVoxel sets a voxel at the given position.
func (vg *VoxelGrid) SetVoxel(x, y, z int, color [3]uint8) {
	if x >= 0 && x < vg.SizeX && y >= 0 && y < vg.SizeY && z >= 0 && z < vg.SizeZ {
		vg.Voxels[[3]int{x, y, z}] = &Voxel{X: x, Y: y, Z: z, Color: color}
	}
}

// GetVoxel retrieves a voxel at the given position.
func (vg *VoxelGrid) GetVoxel(x, y, z int) *Voxel {
	return vg.Voxels[[3]int{x, y, z}]
}

// HasVoxel checks if a voxel exists at the given position.
func (vg *VoxelGrid) HasVoxel(x, y, z int) bool {
	_, ok := vg.Voxels[[3]int{x, y, z}]
	return ok
}

// Count returns the number of voxels in the grid.
func (vg *VoxelGrid) Count() int {
	return len(vg.Voxels)
}
