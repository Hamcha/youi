package components

import (
	"image"

	"github.com/hamcha/youi/font"
)

// Component is a renderable UI component
type Component interface {
	Draw()
	ShouldDraw() bool
	Parent() Container
	Bounds() image.Rectangle

	setParent(Container)
	setBounds(image.Rectangle)
}

// componentBase is the common parent of all components
type componentBase struct {
	parent Container
	bounds image.Rectangle

	dirtyBounds bool
}

func (c *componentBase) Parent() Container {
	return c.parent
}

func (c *componentBase) setParent(container Container) {
	c.parent = container
}

func (c *componentBase) setBounds(rect image.Rectangle) {
	c.bounds = rect
	c.dirtyBounds = true
}

func (c *componentBase) Bounds() image.Rectangle {
	return c.bounds
}

func (c *componentBase) isDirty() bool {
	return c.dirtyBounds
}

// componentText is a common parent of all text-based components
type componentText struct {
	fontFace string
	fontSize float64

	_font     *font.Font
	dirtyFont bool
}

func (c *componentText) SetFontFace(name string) {
	c.fontFace = name
	c.dirtyFont = true
}

func (c *componentText) SetFontSize(size float64) {
	c.fontSize = size
}

func (c *componentText) makeFace() {
	if c.dirtyFont {
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
		c.dirtyFont = false
	}
}

func (c *componentText) isDirty() bool {
	return c.dirtyFont
}
