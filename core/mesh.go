package core

import "io"

// Mesh represents a 3D polygon mesh with vertices, faces, and optional materials.
type Mesh struct {
	Vertices  []Vertex
	Faces     []Face
	Materials []Material
	Bounds    BoundingBox
}

// Vertex represents a 3D point with optional normal and texture coordinates.
type Vertex struct {
	Position [3]float64
	Normal   [3]float64
	TexCoord [2]float64
}

// Face represents a polygon face with vertex indices and material reference.
type Face struct {
	VertexIndices   []int
	NormalIndices   []int
	TexCoordIndices []int
	MaterialIndex   int
}

// Material represents surface properties including color and texture.
type Material struct {
	Name          string
	DiffuseColor  [3]float64 // RGB [0,1]
	AmbientColor  [3]float64
	SpecularColor [3]float64
	EmissiveColor [3]float64
	Opacity       float64
	TexturePath   string
}

// BoundingBox represents axis-aligned bounding box.
type BoundingBox struct {
	Min [3]float64
	Max [3]float64
}

// MeshImporter is the interface for importing polygon meshes from various formats.
type MeshImporter interface {
	// Import reads and parses a mesh from the given reader.
	Import(r io.Reader) (*Mesh, error)
	
	// SupportedFormats returns the list of supported file extensions.
	SupportedFormats() []string
}

// CalculateBounds computes the bounding box of the mesh.
func (m *Mesh) CalculateBounds() {
	if len(m.Vertices) == 0 {
		return
	}
	
	m.Bounds.Min = m.Vertices[0].Position
	m.Bounds.Max = m.Vertices[0].Position
	
	for _, v := range m.Vertices[1:] {
		for i := 0; i < 3; i++ {
			if v.Position[i] < m.Bounds.Min[i] {
				m.Bounds.Min[i] = v.Position[i]
			}
			if v.Position[i] > m.Bounds.Max[i] {
				m.Bounds.Max[i] = v.Position[i]
			}
		}
	}
}
