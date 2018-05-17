package components

import (
	"image"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/hamcha/youi/font"
	"github.com/hamcha/youi/opengl"
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

func (c *ComponentBase) Draw() {
	c.clearFlags()
}

func (c *ComponentBase) ShouldDraw() bool {
	return c.dirtyBounds
}

func (c *ComponentBase) clearFlags() {
	c.dirtyBounds = false
}

func (c *ComponentBase) root() Container {
	return c.parent.Root()
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

func (c *componentText) ShouldDraw() bool {
	return c.dirtyFont
}

func (c *componentText) clearFlags() {
	c.dirtyFont = false
}

func (c *componentText) Draw() {
	if c.dirtyFont {
		c.makeFace()
		c.clearFlags()
	}
}

type ComponentDrawable struct {
	ComponentBase
	quad   *opengl.Quad
	shader *opengl.Shader
}

func (c *ComponentDrawable) Draw() {
	if c.quad == nil {
		c.shader = opengl.DefaultShader()
		c.quad = opengl.MakeQuad(c.shader)
	}

	// Check if bounds have changed
	if c.dirtyBounds {
		// Update transform matrix
		c.updateTransformMatrix()
	}

	c.quad.Draw()
	c.ComponentBase.Draw()
}

func (c *ComponentDrawable) updateTransformMatrix() {
	// Get window resolution
	res := c.root().Bounds().Size()

	// Start from identity
	id := mgl32.Ident4()

	// Move to coordinate
	position := mgl32.Translate3D(float32(c.bounds.Min.X), float32(c.bounds.Min.Y), 1.0)

	// Scale to size
	size := mgl32.Scale3D(float32(c.bounds.Max.X), float32(c.bounds.Max.Y), 1.0)

	// Divide by resolution
	resolution := mgl32.Scale3D(1.0/float32(res.X), 1.0/float32(res.Y), 1.0)

	// Multiply everything into a transform matrix
	result := id.Mul4(position).Mul4(size).Mul4(resolution)

	// Set result matrix as uniform value
	c.shader.GetUniform("transform").Set(result)
}
