package opengl

import "image/color"

// Converts from color.Color to 4 [0.0,1.0] values (for GL colors)
func toGLColor(col color.Color) (float32, float32, float32, float32) {
	const rgbaDivider = float32(0xffff)
	r, g, b, a := col.RGBA()
	return float32(r) / rgbaDivider, float32(g) / rgbaDivider, float32(b) / rgbaDivider, float32(a) / rgbaDivider
}
