package core

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// TextureExtractor extracts block textures and calculates average colors.
type TextureExtractor struct {
	blockModels map[string]BlockModel
	textures    map[string]image.Image
}

// BlockModel represents a Minecraft block model.
type BlockModel struct {
	Parent   string                 `json:"parent"`
	Textures map[string]string      `json:"textures"`
	Elements []interface{}          `json:"elements"`
}

// BlockStateDefinition represents a block state definition.
type BlockStateDefinition struct {
	Variants map[string]interface{} `json:"variants"`
}

// NewTextureExtractor creates a new texture extractor.
func NewTextureExtractor() *TextureExtractor {
	return &TextureExtractor{
		blockModels: make(map[string]BlockModel),
		textures:    make(map[string]image.Image),
	}
}

// ExtractFromResourcePack extracts blocks from a resource pack (zip file or directory).
func (te *TextureExtractor) ExtractFromResourcePack(path string) ([]MinecraftBlock, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat resource pack: %w", err)
	}
	
	if info.IsDir() {
		return te.extractFromDirectory(path)
	}
	return te.extractFromZip(path)
}

// ExtractFromJar extracts blocks from a Minecraft jar file.
func (te *TextureExtractor) ExtractFromJar(jarPath string) ([]MinecraftBlock, error) {
	return te.extractFromZip(jarPath)
}

// extractFromZip extracts blocks from a zip file (jar or resource pack).
func (te *TextureExtractor) extractFromZip(zipPath string) ([]MinecraftBlock, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()
	
	// Load textures
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "assets/minecraft/textures/block/") && 
		   (strings.HasSuffix(f.Name, ".png") || strings.HasSuffix(f.Name, ".jpg")) {
			
			rc, err := f.Open()
			if err != nil {
				continue
			}
			
			img, _, err := image.Decode(rc)
			rc.Close()
			
			if err != nil {
				continue
			}
			
			// Extract texture name
			textureName := strings.TrimPrefix(f.Name, "assets/minecraft/textures/")
			textureName = strings.TrimSuffix(textureName, filepath.Ext(textureName))
			te.textures[textureName] = img
		}
	}
	
	// Load block models
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "assets/minecraft/models/block/") && 
		   strings.HasSuffix(f.Name, ".json") {
			
			rc, err := f.Open()
			if err != nil {
				continue
			}
			
			var model BlockModel
			decoder := json.NewDecoder(rc)
			err = decoder.Decode(&model)
			rc.Close()
			
			if err != nil {
				continue
			}
			
			modelName := strings.TrimPrefix(f.Name, "assets/minecraft/models/block/")
			modelName = strings.TrimSuffix(modelName, ".json")
			te.blockModels[modelName] = model
		}
	}
	
	return te.generateBlocksFromModels()
}

// extractFromDirectory extracts blocks from a directory.
func (te *TextureExtractor) extractFromDirectory(dirPath string) ([]MinecraftBlock, error) {
	// Load textures
	texturesDir := filepath.Join(dirPath, "assets", "minecraft", "textures", "block")
	if _, err := os.Stat(texturesDir); err == nil {
		err = filepath.Walk(texturesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if info.IsDir() {
				return nil
			}
			
			if !strings.HasSuffix(path, ".png") && !strings.HasSuffix(path, ".jpg") {
				return nil
			}
			
			f, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()
			
			img, _, err := image.Decode(f)
			if err != nil {
				return nil
			}
			
			// Extract texture name
			relPath, _ := filepath.Rel(filepath.Join(dirPath, "assets", "minecraft", "textures"), path)
			textureName := strings.TrimSuffix(relPath, filepath.Ext(relPath))
			textureName = strings.ReplaceAll(textureName, string(filepath.Separator), "/")
			te.textures[textureName] = img
			
			return nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to walk textures: %w", err)
		}
	}
	
	// Load block models
	modelsDir := filepath.Join(dirPath, "assets", "minecraft", "models", "block")
	if _, err := os.Stat(modelsDir); err == nil {
		err = filepath.Walk(modelsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}
			
			f, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()
			
			var model BlockModel
			decoder := json.NewDecoder(f)
			err = decoder.Decode(&model)
			if err != nil {
				return nil
			}
			
			relPath, _ := filepath.Rel(modelsDir, path)
			modelName := strings.TrimSuffix(relPath, ".json")
			modelName = strings.ReplaceAll(modelName, string(filepath.Separator), "/")
			te.blockModels[modelName] = model
			
			return nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to walk models: %w", err)
		}
	}
	
	return te.generateBlocksFromModels()
}

// generateBlocksFromModels generates MinecraftBlock entries from loaded models and textures.
func (te *TextureExtractor) generateBlocksFromModels() ([]MinecraftBlock, error) {
	var blocks []MinecraftBlock
	
	for modelName, model := range te.blockModels {
		// Get primary texture
		texturePath := te.resolveTexture(model)
		if texturePath == "" {
			continue
		}
		
		img, ok := te.textures[texturePath]
		if !ok {
			continue
		}
		
		// Calculate average color
		avgColor := te.calculateAverageColor(img)
		
		block := MinecraftBlock{
			ID:         "minecraft:" + modelName,
			RGB:        avgColor,
			Properties: make(map[string]string),
		}
		
		blocks = append(blocks, block)
	}
	
	return blocks, nil
}

// resolveTexture resolves the primary texture path from a block model.
func (te *TextureExtractor) resolveTexture(model BlockModel) string {
	// Try common texture keys
	keys := []string{"all", "texture", "particle", "side", "top", "front"}
	
	for _, key := range keys {
		if texture, ok := model.Textures[key]; ok {
			return te.resolveTextureReference(texture, model)
		}
	}
	
	// If no texture found, try parent model
	if model.Parent != "" {
		parentName := strings.TrimPrefix(model.Parent, "minecraft:block/")
		if parent, ok := te.blockModels[parentName]; ok {
			return te.resolveTexture(parent)
		}
	}
	
	return ""
}

// resolveTextureReference resolves a texture reference (which may start with #).
func (te *TextureExtractor) resolveTextureReference(texture string, model BlockModel) string {
	// Remove minecraft: prefix
	texture = strings.TrimPrefix(texture, "minecraft:")
	
	// If it references another texture variable, resolve it recursively
	if strings.HasPrefix(texture, "#") {
		varName := strings.TrimPrefix(texture, "#")
		if resolved, ok := model.Textures[varName]; ok {
			return te.resolveTextureReference(resolved, model)
		}
	}
	
	return texture
}

// calculateAverageColor calculates the average color of an image.
func (te *TextureExtractor) calculateAverageColor(img image.Image) [3]uint8 {
	bounds := img.Bounds()
	var r, g, b uint64
	var count uint64
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			pr, pg, pb, pa := pixel.RGBA()
			
			// Skip fully transparent pixels
			if pa == 0 {
				continue
			}
			
			// Convert from 16-bit to 8-bit
			r += uint64(pr >> 8)
			g += uint64(pg >> 8)
			b += uint64(pb >> 8)
			count++
		}
	}
	
	if count == 0 {
		return [3]uint8{128, 128, 128}
	}
	
	return [3]uint8{
		uint8(r / count),
		uint8(g / count),
		uint8(b / count),
	}
}

// LoadBlocksFromJSON loads block definitions from a JSON file.
func LoadBlocksFromJSON(path string) ([]MinecraftBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer f.Close()
	
	var blocks []MinecraftBlock
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&blocks); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	
	return blocks, nil
}

// SaveBlocksToJSON saves block definitions to a JSON file.
func SaveBlocksToJSON(blocks []MinecraftBlock, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer f.Close()
	
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(blocks); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	
	return nil
}
