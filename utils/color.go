package utils

// HexColor is a single 8bpc color from a single hex number (similar to CSS)
type HexColor uint32

// RGBA returns the color as separate alpha-premultiplied color components within [0, 0xffff]
func (hcol HexColor) RGBA() (r, g, b, a uint32) {
	h32 := uint32(hcol)
	r = (h32 >> 16) & 0xff00
	g = (h32 >> 8) & 0xff00
	b = h32 & 0xff00
	a = (h32 & 0xff) << 8
	return
}
