package opengl

import (
	"errors"
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Texture is a single OpenGL texture
type Texture struct {
	handle uint32
	unit   uint32
}

// TextureOptions contains extra options for setting up a texture
type TextureOptions struct {
	WrapS     WrapType
	WrapR     WrapType
	MinFilter TextureFilter
	MagFilter TextureFilter
}

// WrapType sets how a texture wraps
type WrapType int32

// Texture wrap types
const (
	TextureWrapRepeat WrapType = gl.REPEAT
	TextureWrapClamp  WrapType = gl.CLAMP_TO_EDGE
	TextureWrapMirror WrapType = gl.MIRRORED_REPEAT
)

// TextureFilter sets what algorithm is used when a texture is scaled
type TextureFilter int32

// Texture filter types
const (
	TextureFilterNearest TextureFilter = gl.NEAREST
	TextureFilterLinear  TextureFilter = gl.LINEAR
)

// Texture errors
var (
	ErrTextureNotBound = errors.New("texture not bound to any texture unit")
)

// MakeTexture creates an OpenGL texture and returns it, if possible
func MakeTexture(img *image.RGBA, options TextureOptions) *Texture {
	texture := new(Texture)

	// Generate texture handle
	gl.GenTextures(1, &texture.handle)

	// Bind texture to set up stuff
	texture.Bind(0)
	defer texture.Unbind()

	// Set wrap types
	if options.WrapR != 0 {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, int32(options.WrapR))
	}
	if options.WrapS != 0 {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, int32(options.WrapS))
	}

	// Set filtering
	if options.MinFilter != 0 {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(options.MinFilter))
	}
	if options.MagFilter != 0 {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(options.MagFilter))
	}

	width := int32(img.Rect.Size().X)
	height := int32(img.Rect.Size().Y)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return texture
}

// Bind binds the texture to a hardware texture unit
func (t *Texture) Bind(unit uint32) {
	t.unit = gl.TEXTURE0 + unit
	gl.ActiveTexture(t.unit)
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

// Unbind removes the texture from its current texture unit
func (t *Texture) Unbind() {
	gl.ActiveTexture(t.unit)
	t.unit = 0
	gl.BindTexture(gl.TEXTURE_2D, t.unit)
}

// SetUniform sets the texture as uniform for a GL program
func (t *Texture) SetUniform(uloc int32) error {
	if t.unit == 0 {
		return ErrTextureNotBound
	}
	gl.Uniform1i(uloc, int32(t.unit-gl.TEXTURE0))
	return nil
}
