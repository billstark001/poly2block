# poly2block WASM

WebAssembly bindings for poly2block core library.

## Building

### Using Standard Go Compiler

```bash
GOOS=js GOARCH=wasm go build -o poly2block.wasm
```

### Using TinyGo (Smaller Binary)

```bash
tinygo build -o poly2block.wasm -target wasm .
```

## Usage

### Loading the Module

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>poly2block WASM Demo</title>
    <script src="wasm_exec.js"></script>
</head>
<body>
    <h1>poly2block WASM Demo</h1>
    <input type="file" id="meshFile" accept=".gltf,.glb">
    <button id="convert">Convert to Schematic</button>
    <div id="output"></div>

    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("poly2block.wasm"), go.importObject)
            .then((result) => {
                go.run(result.instance);
                console.log("poly2block version:", poly2block.version);
            });

        document.getElementById('convert').addEventListener('click', async () => {
            const file = document.getElementById('meshFile').files[0];
            if (!file) {
                alert('Please select a file');
                return;
            }

            const arrayBuffer = await file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);

            const result = poly2block.meshToSchematic(
                uint8Array,
                128,     // resolution
                true,    // conservative
                true,    // dither
                null     // use vanilla palette
            );

            if (result.success) {
                // result.data is base64 encoded schematic
                console.log('Conversion successful!');
                downloadBase64(result.data, 'output.schem');
            } else {
                console.error('Conversion failed:', result.error);
                alert('Error: ' + result.error);
            }
        });

        function downloadBase64(base64Data, filename) {
            const binary = atob(base64Data);
            const bytes = new Uint8Array(binary.length);
            for (let i = 0; i < binary.length; i++) {
                bytes[i] = binary.charCodeAt(i);
            }
            const blob = new Blob([bytes], { type: 'application/octet-stream' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = filename;
            a.click();
            URL.revokeObjectURL(url);
        }
    </script>
</body>
</html>
```

## API

### poly2block.meshToVox(meshData, resolution, conservative)

Convert a mesh to VOX format.

**Parameters:**
- `meshData`: Uint8Array or base64 string containing glTF/GLB data
- `resolution`: Number - voxel resolution (e.g., 128)
- `conservative`: Boolean - use conservative voxelization

**Returns:**
```javascript
{
    success: true,
    data: "base64-encoded-vox-data"
}
// or
{
    success: false,
    error: "error message"
}
```

### poly2block.meshToSchematic(meshData, resolution, conservative, dither, paletteData)

Convert a mesh to Minecraft schematic.

**Parameters:**
- `meshData`: Uint8Array or base64 string containing glTF/GLB data
- `resolution`: Number - voxel resolution
- `conservative`: Boolean - use conservative voxelization
- `dither`: Boolean - enable Floyd-Steinberg dithering
- `paletteData`: Uint8Array, base64 string, or null (uses vanilla blocks)

**Returns:** Same format as `meshToVox`

### poly2block.generatePalette()

Generate a vanilla Minecraft block palette.

**Returns:**
```javascript
{
    success: true,
    data: "base64-encoded-palette-data"
}
```

## Examples

### Convert with Custom Palette

```javascript
// Generate palette
const paletteResult = poly2block.generatePalette();
const paletteData = paletteResult.data;

// Convert mesh with custom palette
const result = poly2block.meshToSchematic(
    meshData,
    256,      // higher resolution
    true,
    true,
    paletteData  // use generated palette
);
```

### Convert to VOX First

```javascript
// Step 1: Convert to VOX
const voxResult = poly2block.meshToVox(meshData, 128, true);

// Step 2: Can save VOX or convert to schematic
// Note: VOX-to-schematic conversion not exposed in WASM yet
```

## Performance Tips

1. Use lower resolutions (64-128) for faster processing
2. Disable dithering for faster but lower quality results
3. Use TinyGo for smaller WASM binaries (~2MB vs ~10MB)
4. Process files in a Web Worker to avoid blocking UI

## License

MIT
