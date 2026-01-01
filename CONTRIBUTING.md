# Contributing to poly2block

Thank you for your interest in contributing to poly2block! This document provides guidelines and information for contributors.

## Code of Conduct

This project follows a simple code of conduct: be respectful, professional, and constructive in all interactions.

## How to Contribute

### Reporting Issues

- Check if the issue already exists
- Provide clear steps to reproduce the problem
- Include relevant system information (OS, Go version, etc.)
- Add screenshots or error logs if applicable

### Suggesting Features

- Explain the use case and motivation
- Describe the proposed solution
- Consider performance implications
- Check if similar features exist

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write or update tests
5. Run tests (`make test`)
6. Format code (`make fmt`)
7. Commit changes (`git commit -m 'Add amazing feature'`)
8. Push to branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional, for convenience)

### Getting Started

```bash
# Clone repository
git clone https://github.com/billstark001/poly2block.git
cd poly2block

# Build project
make build

# Run tests
make test

# Build WASM
make wasm
```

## Project Structure

```
poly2block/
├── core/              # Core algorithm library
├── wasm/              # WebAssembly bindings  
├── cmd/poly2block/    # CLI application
└── .github/workflows/ # CI/CD pipelines
```

## Coding Guidelines

### Go Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, self-documenting code
- Add comments for exported functions/types
- Keep functions focused and small

### Testing

- Write unit tests for new functionality
- Maintain or improve test coverage
- Test edge cases and error conditions
- Use table-driven tests when appropriate

### Documentation

- Update README files when adding features
- Document public APIs with godoc comments
- Include usage examples
- Keep documentation concise and clear

## Areas for Contribution

### High Priority

- OBJ/MTL mesh importer
- Additional voxelization algorithms
- Performance optimizations
- More comprehensive tests

### Medium Priority

- Additional dithering algorithms
- Custom palette editor/generator
- More output formats
- Better error messages

### Ideas Welcome

- GPU-accelerated voxelization
- Web-based GUI
- Minecraft resource pack integration
- Additional mesh formats (FBX, etc.)
- Batch processing tools

## Algorithm Implementation

When implementing new algorithms:

1. **Interface First**: Define clear interfaces in `core/`
2. **Generic Design**: Make implementations swappable
3. **Performance**: Consider both speed and memory
4. **Testing**: Add comprehensive tests
5. **Documentation**: Explain algorithm choices

### Example: Adding a New Voxelizer

```go
// In core/voxelizer_custom.go
type CustomVoxelizer struct {
    // Configuration
}

func NewCustomVoxelizer() *CustomVoxelizer {
    return &CustomVoxelizer{}
}

func (v *CustomVoxelizer) Voxelize(mesh *Mesh, config VoxelizationConfig) (*VoxelGrid, error) {
    // Implementation
}

func (v *CustomVoxelizer) Name() string {
    return "custom-voxelizer"
}
```

## Testing Guidelines

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific package
cd core && go test -v ./...
```

### Writing Tests

```go
func TestNewFeature(t *testing.T) {
    // Setup
    input := setupTestData()
    
    // Execute
    result, err := NewFeature(input)
    
    // Verify
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

## Performance Considerations

- Profile before optimizing
- Consider memory allocation patterns
- Use benchmarks to measure improvements
- Document performance characteristics

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. Update version in relevant files
2. Create and push tag: `git tag v1.0.0 && git push origin v1.0.0`
3. GitHub Actions builds and publishes release

## Questions?

Feel free to:
- Open an issue for questions
- Start a discussion on GitHub
- Check existing issues and PRs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
