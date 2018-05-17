package opengl

import (
	"image/color"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Converts from color.Color to 4 [0.0,1.0] values (for GL colors)
func toGLColor(col color.Color) (float32, float32, float32, float32) {
	const rgbaDivider = float32(0xff00)
	r, g, b, a := col.RGBA()
	return float32(r) / rgbaDivider, float32(g) / rgbaDivider, float32(b) / rgbaDivider, float32(a) / rgbaDivider
}

func glString(str string) *uint8 {
	return gl.Str(str + "\x00")
}
