package components

import (
	"github.com/golang/freetype/truetype"
	"github.com/hamcha/youi"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Component is a renderable UI component
type Component interface {
	Draw()
}

// componentPosition add position to a component
type componentPosition struct {
	x, y          float32
	dirtyPosition bool
}

func (c *componentPosition) SetPosition(x, y float32) {
	c.x = x
	c.y = y
	c.dirtyPosition = true
}

// componentText allows components to have text in them
type componentText struct {
	fontFace *truetype.Font
	fontSize float64

	_font     font.Face
	dirtyFont bool
}

func (c *componentText) SetFontFace(fnt *truetype.Font) {
	c.fontFace = fnt
	c.dirtyFont = true
}

func (c *componentText) SetFontSize(size float64) {
	c.fontSize = size
	c.dirtyFont = true
}

func (c *componentText) makeFace() {
	if c.dirtyFont {
		// If no font is provided, use Go Regolar
		if c.fontFace == nil {
			c.fontFace, _ = truetype.Parse(goregular.TTF)
		}
		c._font = truetype.NewFace(c.fontFace, &truetype.Options{
			Size: c.fontSize,
			DPI:  youi.SYSDPI,
		})
	}
}
