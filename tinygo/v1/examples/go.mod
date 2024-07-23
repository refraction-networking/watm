module github.com/refraction-networking/watm/tinygo/v1/examples

go 1.21

replace golang.org/x/sys => ../../replaced/golang.org/x/sys@v0.19.0

replace github.com/refraction-networking/watm => ../../../

require (
	github.com/CosmWasm/tinyjson v0.9.0
	github.com/refraction-networking/utls v1.6.7-wasm
	github.com/refraction-networking/watm v0.6.5
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/cloudflare/circl v1.3.9 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)
