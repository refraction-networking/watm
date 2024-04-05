module github.com/refraction-networking/watm/tinygo/v1/examples

go 1.21

replace golang.org/x/sys v0.16.0 => ./replace/x/sys

replace github.com/refraction-networking/watm => ../../../

require (
	github.com/CosmWasm/tinyjson v0.9.0
	github.com/refraction-networking/utls v1.6.3-wasm
	github.com/refraction-networking/watm v0.0.0-00010101000000-000000000000
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/quic-go/quic-go v0.42.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/sys v0.19.0 // indirect
)
