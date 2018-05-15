# youi

youi is an experimental UI library for Go projects using a OpenGL backend for a consistent look across platformers rather than aiming for a native look'n'feel.

## Install

```
go get -u github.com/hamcha/youi/...
```

## Compatibility

youi is currently targeting OpenGL 3.3 core, which should work on most hardware from 2008 onwards:

- [All Apple hardware from 2007/2008 and newer](https://support.apple.com/en-us/HT202823)
- Sixth generation Intel HD Graphics and newer (except on Windows, where you need at least 7th)
- NVIDIA Graphics cards G80 (GeForce 8 series) and newer
- AMD Graphics cards R600 (Radeon HD 2000 series) and newer

Actual support may be affected by which OS drivers are installed.

## License

The code for this project is licensed under the ISC license, full text provided in the `LICENSE` file.

This project depends on third party libraries, you can see their licenses in their relative files:

- `LICENSE-go-gl` for the [go-gl project](https://github.com/go-gl)
- `LICENSE-ftl` for [freetype2](http://freetype.sourceforge.net/index2.html) and [Freetype-Go](https://github.com/golang/freetype)
- `LICENSE-texpack` for [texpack](https://github.com/adinfinit/texpack)
- `LICENSE-go-rice` for [go.rice](https://github.com/GeertJohan/go.rice)