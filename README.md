# `watm`: WebAssembly Transport Module
![Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-green)
![GPLv3](https://img.shields.io/badge/License-GPL--3.0-red)
[![Test](https://github.com/refraction-networking/watm/actions/workflows/watm.yml/badge.svg?branch=master)](https://github.com/refraction-networking/watm/actions/workflows/watm.yml)
[![Release Status](https://github.com/refraction-networking/watm/actions/workflows/release.yml/badge.svg)](https://github.com/refraction-networking/watm/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/refraction-networking/watm.svg)](https://pkg.go.dev/github.com/refraction-networking/watm)

This repository contains tools for building WebAssembly Transport Modules (WATMs) for [water](https://github.com/refraction-networking/water) project. 

### Cite our work

If you quoted or used our work in your own project/paper/research, please cite our paper [Just add WATER: WebAssembly-based Circumvention Transports](https://www.petsymposium.org/foci/2024/foci-2024-0003.pdf), which is published in the proceedings of Free and Open Communications on the Internet (FOCI) in 2024 issue 1, pages 22-28.

<details>
  <summary>BibTeX</summary>
    
  ```bibtex
    @inproceedings{water-foci24,
        author = {Chi, Erik and Wang, Gaukas and Halderman, J. Alex and Wustrow, Eric and Wampler, Jack},
        year = {2024},
        month = {02},
        number = {1},
        pages = {22-28},
        title = {Just add {WATER}: {WebAssembly}-based Circumvention Transports},
        howpublished = "\url{https://www.petsymposium.org/foci/2024/foci-2024-0003.php}",
        publisher = {PoPETs},
        address = {Virtual Event},
        series = {FOCI '24},
        booktitle = {Free and Open Communications on the Internet},
    }
  ```
</details>

# License

This project is dual-licensed under both the Apache 2.0 license and the GPLv3 license. The license applies differently depending on how this project is used.

- **Apache 2.0**: applies for the project itself, and all of its packages EXCEPT `examples`.
- **GPLv3** applies when your project uses the code from `examples` package, including but not limited to when you modify and redistribute the example code, or even use it for a non-water scenario. However, if you decide to distribute the examples ONLY in a compiled form (i.e., the `.wasm` file), you are free to use the compiled output without a problem.