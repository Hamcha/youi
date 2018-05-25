package components

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hamcha/youi/opengl"
)

type Drawable struct {
	Base
	Mesh *opengl.Mesh
}

func (c *Drawable) Draw() {
	if c.Mesh == nil {
		c.Mesh = opengl.MakeQuad(opengl.DefaultShader())
	}

	// Check if bounds have changed
	if c.dirtyBounds {
		// Update transform matrix
		c.updateTransformMatrix()
	}

	c.Mesh.Draw()
	c.Base.Draw()
}

func (c *Drawable) updateTransformMatrix() {
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
	c.Mesh.Shader.GetUniform("transform").Set(result)
}

// Component handling errors
var (
	ErrIndexOutOfBounds  = errors.New("index out of bounds")
	ErrComponentNotFound = errors.New("component not found in list")
)
