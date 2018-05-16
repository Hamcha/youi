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

// ComponentBase is the common parent of all components
type ComponentBase struct {
	parent Container
	bounds image.Rectangle

	dirtyBounds bool
}

func (c *ComponentBase) Parent() Container {
	return c.parent
}

func (c *ComponentBase) setParent(container Container) {
	c.parent = container
}

func (c *ComponentBase) setBounds(rect image.Rectangle) {
	c.bounds = rect
	c.dirtyBounds = true
}

func (c *ComponentBase) Bounds() image.Rectangle {
	return c.bounds
}

func (c *ComponentBase) isDirty() bool {
	return c.dirtyBounds
}

func (c *ComponentBase) ClearFlags() {
	c.dirtyBounds = false
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

func (c *componentText) ClearFlags() {
	c.dirtyFont = false
}
