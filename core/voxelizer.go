package core

import (
	"fmt"
	"math"
)

// SurfaceVoxelizer implements basic surface voxelization.
type SurfaceVoxelizer struct{}

// NewSurfaceVoxelizer creates a new surface voxelizer.
func NewSurfaceVoxelizer() *SurfaceVoxelizer {
	return &SurfaceVoxelizer{}
}

// Voxelize converts a mesh to a voxel grid using surface voxelization.
func (v *SurfaceVoxelizer) Voxelize(mesh *Mesh, config VoxelizationConfig) (*VoxelGrid, error) {
	if len(mesh.Vertices) == 0 {
		return nil, fmt.Errorf("mesh has no vertices")
	}
	
	// Calculate bounds if not already done
	if mesh.Bounds.Min == [3]float64{} && mesh.Bounds.Max == [3]float64{} {
		mesh.CalculateBounds()
	}
	
	// Calculate dimensions
	dims := [3]float64{
		mesh.Bounds.Max[0] - mesh.Bounds.Min[0],
		mesh.Bounds.Max[1] - mesh.Bounds.Min[1],
		mesh.Bounds.Max[2] - mesh.Bounds.Min[2],
	}
	
	// Find longest dimension
	maxDim := math.Max(dims[0], math.Max(dims[1], dims[2]))
	if maxDim == 0 {
		return nil, fmt.Errorf("mesh has zero size")
	}
	
	// Calculate scale
	scale := float64(config.Resolution) / maxDim
	if config.Scale > 0 {
		scale = config.Scale
	}
	
	// Calculate grid size
	sizeX := int(math.Ceil(dims[0] * scale))
	sizeY := int(math.Ceil(dims[1] * scale))
	sizeZ := int(math.Ceil(dims[2] * scale))
	
	// Create voxel grid
	voxelGrid := NewVoxelGrid(sizeX, sizeY, sizeZ)
	voxelGrid.Scale = scale
	voxelGrid.Origin = mesh.Bounds.Min
	
	// Voxelize each face
	for _, face := range mesh.Faces {
		if len(face.VertexIndices) < 3 {
			continue
		}
		
		// Get triangle vertices
		v0 := mesh.Vertices[face.VertexIndices[0]].Position
		v1 := mesh.Vertices[face.VertexIndices[1]].Position
		v2 := mesh.Vertices[face.VertexIndices[2]].Position
		
		// Get material color
		color := [3]uint8{128, 128, 128} // Default gray
		if face.MaterialIndex >= 0 && face.MaterialIndex < len(mesh.Materials) {
			mat := mesh.Materials[face.MaterialIndex]
			color = [3]uint8{
				uint8(mat.DiffuseColor[0] * 255),
				uint8(mat.DiffuseColor[1] * 255),
				uint8(mat.DiffuseColor[2] * 255),
			}
		}
		
		// Rasterize triangle
		v.rasterizeTriangle(voxelGrid, v0, v1, v2, color, config.Conservative)
	}
	
	return voxelGrid, nil
}

// rasterizeTriangle rasterizes a triangle into the voxel grid.
func (v *SurfaceVoxelizer) rasterizeTriangle(grid *VoxelGrid, v0, v1, v2 [3]float64, color [3]uint8, conservative bool) {
	// Transform vertices to voxel space
	v0Voxel := v.worldToVoxel(v0, grid)
	v1Voxel := v.worldToVoxel(v1, grid)
	v2Voxel := v.worldToVoxel(v2, grid)
	
	// Calculate triangle bounds
	minX := int(math.Floor(math.Min(v0Voxel[0], math.Min(v1Voxel[0], v2Voxel[0]))))
	maxX := int(math.Ceil(math.Max(v0Voxel[0], math.Max(v1Voxel[0], v2Voxel[0]))))
	minY := int(math.Floor(math.Min(v0Voxel[1], math.Min(v1Voxel[1], v2Voxel[1]))))
	maxY := int(math.Ceil(math.Max(v0Voxel[1], math.Max(v1Voxel[1], v2Voxel[1]))))
	minZ := int(math.Floor(math.Min(v0Voxel[2], math.Min(v1Voxel[2], v2Voxel[2]))))
	maxZ := int(math.Ceil(math.Max(v0Voxel[2], math.Max(v1Voxel[2], v2Voxel[2]))))
	
	// Clamp to grid bounds
	minX = max(0, minX)
	maxX = min(grid.SizeX-1, maxX)
	minY = max(0, minY)
	maxY = min(grid.SizeY-1, maxY)
	minZ = max(0, minZ)
	maxZ = min(grid.SizeZ-1, maxZ)
	
	// Scan all voxels in the bounding box
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				voxelCenter := [3]float64{
					float64(x) + 0.5,
					float64(y) + 0.5,
					float64(z) + 0.5,
				}
				
				// Check if voxel intersects triangle
				if v.voxelIntersectsTriangle(voxelCenter, v0Voxel, v1Voxel, v2Voxel, conservative) {
					grid.SetVoxel(x, y, z, color)
				}
			}
		}
	}
}

// worldToVoxel transforms world coordinates to voxel coordinates.
func (v *SurfaceVoxelizer) worldToVoxel(world [3]float64, grid *VoxelGrid) [3]float64 {
	return [3]float64{
		(world[0] - grid.Origin[0]) * grid.Scale,
		(world[1] - grid.Origin[1]) * grid.Scale,
		(world[2] - grid.Origin[2]) * grid.Scale,
	}
}

// voxelIntersectsTriangle checks if a voxel intersects with a triangle.
// This is a simplified check using barycentric coordinates.
func (v *SurfaceVoxelizer) voxelIntersectsTriangle(voxel, v0, v1, v2 [3]float64, conservative bool) bool {
	// Calculate triangle normal
	edge1 := sub3(v1, v0)
	edge2 := sub3(v2, v0)
	normal := cross3(edge1, edge2)
	
	// Calculate distance from voxel to triangle plane
	d := dot3(normal, v0)
	dist := math.Abs(dot3(normal, voxel) - d)
	
	// Check if voxel is close to plane
	threshold := 0.866 // sqrt(3)/2 for voxel diagonal
	if conservative {
		threshold *= 1.5
	}
	
	if dist > threshold {
		return false
	}
	
	// Check if projection is inside triangle using barycentric coordinates
	// Simplified check: test if point is on same side of all edges
	return v.pointInTriangle2D(voxel, v0, v1, v2)
}

// pointInTriangle2D checks if a point is inside a triangle using 2D projection.
func (v *SurfaceVoxelizer) pointInTriangle2D(p, v0, v1, v2 [3]float64) bool {
	// Use XY projection for simplicity
	sign := func(p1, p2, p3 [3]float64) float64 {
		return (p1[0]-p3[0])*(p2[1]-p3[1]) - (p2[0]-p3[0])*(p1[1]-p3[1])
	}
	
	d1 := sign(p, v0, v1)
	d2 := sign(p, v1, v2)
	d3 := sign(p, v2, v0)
	
	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	
	return !(hasNeg && hasPos)
}

// Name returns the algorithm name.
func (v *SurfaceVoxelizer) Name() string {
	return "surface-voxelizer"
}

// Helper functions
func sub3(a, b [3]float64) [3]float64 {
	return [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func cross3(a, b [3]float64) [3]float64 {
	return [3]float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func dot3(a, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
