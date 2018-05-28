package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

func getTransformMatrix(bounds Bounds) mgl32.Mat4 {
	// Get size and position
	size := bounds.Size
	pos := bounds.Position

	// Change pivot to top-left instead of center
	pos = Position{pos.X + size.Width/2, pos.Y + size.Height/2}

	// Set position, taking into account OpenGL's weird coordinate system
	posMtx := mgl32.Translate3D((pos.X*2)-1, -(pos.Y*2)+1, 0.0)

	// Scale to size
	sizeMtx := mgl32.Scale3D(size.Width, size.Height, 1.0)

	// Multiply everything into a transform matrix
	result := posMtx.Mul4(sizeMtx)

	return result
}
