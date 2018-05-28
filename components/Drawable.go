package components

import (
	"errors"

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
		mat := getTransformMatrix(c.bounds)
		c.Mesh.Shader.GetUniform("transform").Set(mat)
	}

	c.Mesh.Draw()
	c.Base.Draw()
}

// Component handling errors
var (
	ErrIndexOutOfBounds  = errors.New("index out of bounds")
	ErrComponentNotFound = errors.New("component not found in list")
)
