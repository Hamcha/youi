package youi

import (
	"image"

	"github.com/hamcha/youi/components"
)

// Root is a special canvas container with less checks
type Root struct {
	components.ContainerBase

	pixelBounds image.Rectangle
}

func (r *Root) SetSize(size image.Point) {
	r.pixelBounds.Max = r.pixelBounds.Min.Add(size)
	r.SetBounds(components.BoundsFromRect(r.pixelBounds))
}

func (r *Root) SetRect(rect image.Rectangle) {
	r.pixelBounds = rect
	r.SetBounds(components.BoundsFromRect(r.pixelBounds))
}

func (r *Root) ShouldDraw() bool {
	return r.ContainerBase.ShouldDraw()
}

func (r *Root) Draw() {
	r.resizeChildren()
	r.ContainerBase.Draw()
}

func (r *Root) resizeChildren() {
	relbounds := components.Bounds{
		components.Position{0, 0},
		components.Size{1, 1},
	}

	// Apply to each children
	for _, child := range r.Children() {
		child.SetBounds(relbounds)
	}
}
