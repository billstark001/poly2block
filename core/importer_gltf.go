package core

import (
	"fmt"
	"io"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// GLTFImporter implements MeshImporter for glTF format.
type GLTFImporter struct{}

// NewGLTFImporter creates a new glTF importer.
func NewGLTFImporter() *GLTFImporter {
	return &GLTFImporter{}
}

// Import reads and parses a glTF mesh from the given reader.
func (imp *GLTFImporter) Import(r io.Reader) (*Mesh, error) {
	// Parse glTF
	doc := gltf.NewDocument()
	decoder := gltf.NewDecoder(r)
	if err := decoder.Decode(doc); err != nil {
		return nil, fmt.Errorf("failed to parse glTF: %w", err)
	}
	
	mesh := &Mesh{
		Vertices:  []Vertex{},
		Faces:     []Face{},
		Materials: []Material{},
	}
	
	// Extract materials
	for _, mat := range doc.Materials {
		material := Material{
			Name: mat.Name,
		}
		
		if mat.PBRMetallicRoughness != nil {
			pbr := mat.PBRMetallicRoughness
			if len(pbr.BaseColorFactor) >= 3 {
				material.DiffuseColor = [3]float64{
					float64(pbr.BaseColorFactor[0]),
					float64(pbr.BaseColorFactor[1]),
					float64(pbr.BaseColorFactor[2]),
				}
			}
		}
		
		mesh.Materials = append(mesh.Materials, material)
	}
	
	// Extract geometry from all meshes
	for _, gltfMesh := range doc.Meshes {
		for _, primitive := range gltfMesh.Primitives {
			if err := imp.extractPrimitive(doc, primitive, mesh); err != nil {
				return nil, fmt.Errorf("failed to extract primitive: %w", err)
			}
		}
	}
	
	mesh.CalculateBounds()
	return mesh, nil
}

// extractPrimitive extracts geometry from a glTF primitive.
func (imp *GLTFImporter) extractPrimitive(doc *gltf.Document, primitive *gltf.Primitive, mesh *Mesh) error {
	// Get position accessor
	posAccessor, ok := primitive.Attributes[gltf.POSITION]
	if !ok {
		return fmt.Errorf("primitive missing POSITION attribute")
	}
	
	positions, err := modeler.ReadPosition(doc, doc.Accessors[posAccessor], nil)
	if err != nil {
		return fmt.Errorf("failed to read positions: %w", err)
	}
	
	// Read normals if available
	var normals [][3]float32
	if normalAccessor, ok := primitive.Attributes[gltf.NORMAL]; ok {
		normals, err = modeler.ReadNormal(doc, doc.Accessors[normalAccessor], nil)
		if err != nil {
			return fmt.Errorf("failed to read normals: %w", err)
		}
	}
	
	// Read texture coordinates if available
	var texCoords [][2]float32
	if texCoordAccessor, ok := primitive.Attributes[gltf.TEXCOORD_0]; ok {
		texCoords, err = modeler.ReadTextureCoord(doc, doc.Accessors[texCoordAccessor], nil)
		if err != nil {
			return fmt.Errorf("failed to read texture coordinates: %w", err)
		}
	}
	
	// Add vertices
	vertexOffset := len(mesh.Vertices)
	for i, pos := range positions {
		vertex := Vertex{
			Position: [3]float64{float64(pos[0]), float64(pos[1]), float64(pos[2])},
		}
		
		if i < len(normals) {
			vertex.Normal = [3]float64{float64(normals[i][0]), float64(normals[i][1]), float64(normals[i][2])}
		}
		
		if i < len(texCoords) {
			vertex.TexCoord = [2]float64{float64(texCoords[i][0]), float64(texCoords[i][1])}
		}
		
		mesh.Vertices = append(mesh.Vertices, vertex)
	}
	
	// Read indices
	if primitive.Indices != nil {
		indices, err := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], nil)
		if err != nil {
			return fmt.Errorf("failed to read indices: %w", err)
		}
		
		// Create faces (assuming triangles)
		for i := 0; i < len(indices); i += 3 {
			if i+2 < len(indices) {
				materialIndex := -1
				if primitive.Material != nil {
					materialIndex = *primitive.Material
				}
				face := Face{
					VertexIndices: []int{
						vertexOffset + int(indices[i]),
						vertexOffset + int(indices[i+1]),
						vertexOffset + int(indices[i+2]),
					},
					MaterialIndex: materialIndex,
				}
				mesh.Faces = append(mesh.Faces, face)
			}
		}
	} else {
		// No indices, treat as triangle list
		for i := 0; i < len(positions); i += 3 {
			if i+2 < len(positions) {
				materialIndex := -1
				if primitive.Material != nil {
					materialIndex = *primitive.Material
				}
				face := Face{
					VertexIndices: []int{
						vertexOffset + i,
						vertexOffset + i + 1,
						vertexOffset + i + 2,
					},
					MaterialIndex: materialIndex,
				}
				mesh.Faces = append(mesh.Faces, face)
			}
		}
	}
	
	return nil
}

// SupportedFormats returns the list of supported file extensions.
func (imp *GLTFImporter) SupportedFormats() []string {
	return []string{".gltf", ".glb"}
}
