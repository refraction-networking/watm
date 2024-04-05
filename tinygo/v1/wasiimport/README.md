# `v1/wasiimport` package

This package exposes the imported functions via WASI interface, which includes both functions from `wasi_snapshot_preview1` module and `env` module.

The idea of this package is to share the imported functions across different packages in this repository. Therefore this package does not do any type/argument transformation, but provide the original functions as they are.