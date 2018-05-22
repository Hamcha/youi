package components

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/hamcha/youi/font"
	"github.com/hamcha/youi/opengl"
)

// Component is a renderable UI component
type Component interface {
	Draw()
	ShouldDraw() bool
	Parent() Container
	Bounds() Bounds
	SetBounds(Bounds)

	setParent(Container)
}

// ComponentBase is the common parent of all components
type ComponentBase struct {
	parent Container
	bounds Bounds

	dirtyBounds bool
}

func (c *ComponentBase) Parent() Container {
	return c.parent
}

func (c *ComponentBase) setParent(container Container) {
	c.parent = container
}

func (c *ComponentBase) SetBounds(bounds Bounds) {
	c.bounds = bounds
	c.dirtyBounds = true
}

func (c *ComponentBase) Bounds() Bounds {
	return c.bounds
}

func (c *ComponentBase) Draw() {
	c.ClearFlags()
}

func (c *ComponentBase) ShouldDraw() bool {
	return c.dirtyBounds
}

func (c *ComponentBase) SetRedraw() {
	c.dirtyBounds = true
}

func (c *ComponentBase) ClearFlags() {
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

func (c *componentText) ClearFlags() {
	c.dirtyFont = false
}

func (c *componentText) Draw() {
	if c.dirtyFont {
		c.makeFace()
		c.ClearFlags()
	}
}

type ComponentDrawable struct {
	ComponentBase
	quad   *opengl.Quad
	shader *opengl.Shader
}

func (c *ComponentDrawable) Draw() {
	if c.shader == nil {
		c.shader = opengl.DefaultShader()
	}
	if c.quad == nil {
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
	// Get size and position
	size := c.bounds.Size
	pos := c.bounds.Position

	// Change pivot to top-left instead of center
	pos = Position{pos.X + size.Width/2, pos.Y + size.Height/2}

	// Set position, taking into account OpenGL's weird coordinate system
	posMtx := mgl32.Translate3D((pos.X*2)-1, -(pos.Y*2)+1, 0.0)

	// Scale to size
	sizeMtx := mgl32.Scale3D(size.Width, size.Height, 1.0)

	// Multiply everything into a transform matrix
	result := posMtx.Mul4(sizeMtx)

	// Set result matrix as uniform value
	c.shader.GetUniform("transform").Set(result)
}
