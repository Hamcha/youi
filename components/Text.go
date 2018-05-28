package components

import (
	"github.com/hamcha/youi/font"
	"github.com/hamcha/youi/opengl"
)

// Text is a common parent of all text-based components
type Text struct {
	Base

	fontFace string
	fontSize float64

	font      *font.Font
	text      *opengl.Text
	dirtyFont bool

	content      string
	dirtyContent bool
}

// SetText changes the text content of the text control
func (c *Text) SetText(str string) {
	c.content = str
	c.dirtyContent = true
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
		c.font = font.DefaultFont()
	} else {
		var err error
		c.font, err = font.LoadFont(c.fontFace)
		if err != nil {
			//TODO Proper error reporting
			panic(err)
		}
	}
}

func (c *Text) ShouldDraw() bool {
	return c.dirtyFont || c.dirtyContent
}

func (c *Text) ClearFlags() {
	c.dirtyFont = false
	c.dirtyContent = false
}

func (c *Text) Draw() {
	if c.font == nil || c.text == nil || c.dirtyFont {
		c.makeFace()
		c.text = opengl.MakeText(c.font, c.content)
	} else if c.dirtyContent {
		c.text.SetContent(c.content)
	}

	// Check if bounds have changed
	if c.dirtyBounds {
		// Update transform matrix
		mat := getTransformMatrix(c.bounds)
		c.text.Shader.GetUniform("transform").Set(mat)
	}

	c.text.Mesh.Draw()
	c.Base.Draw()
	c.ClearFlags()
}
