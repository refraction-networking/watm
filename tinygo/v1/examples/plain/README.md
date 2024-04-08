# Example: `plain.wasm` 

This example shows how to build a minimal WATM with TinyGo which passes through the received/sent data without any processing.

## Build

Go 1.20/1.21/1.22 is required to build this example. TinyGo started to support Go 1.22 in v0.31.0.

### Debug

```bash
tinygo build -o plain.wasm -target=wasi -scheduler=none -gc=conservative .
```

### Release

```bash
tinygo build -o plain.wasm -target=wasi -no-debug -scheduler=none -gc=conservative .
```