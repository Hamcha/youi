package components

import "github.com/hamcha/youi/font"

// Text is a common parent of all text-based components
type Text struct {
	fontFace string
	fontSize float64

	_font     *font.Font
	dirtyFont bool
}

func (c *Text) SetFontFace(name string) {
	c.fontFace = name
	c.dirtyFont = true
}

func (c *Text) SetFontSize(size float64) {
	c.fontSize = size
}

func (c *Text) makeFace() {
	// If no font is provided, use Go Regolar
	if c.fontFace == "" {
		c._font = font.DefaultFont()
	} else {
		var err error
		c._font, err = font.LoadFont(c.fontFace)
		if err != nil {
			//TODO Proper error reporting
			panic(err)
		}
	}
}

func (c *Text) ShouldDraw() bool {
	return c.dirtyFont
}

func (c *Text) ClearFlags() {
	c.dirtyFont = false
}

func (c *Text) Draw() {
	if c.dirtyFont {
		c.makeFace()
		c.ClearFlags()
	}

	//TODO Draw font
}
