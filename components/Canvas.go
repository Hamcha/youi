package components

import (
	"image"
)

// Canvas is a container that has absolute positioning and sizing regardless of parent/siblings.
// While Canvas can hold multiple childrens, those childrens will share the same position/size
// so they will all overlap each other. It is therefore recommended to only put one component
// inside a Canvas.
type Canvas struct {
	ContainerBase

	canvasBounds image.Rectangle
}

func (c *Canvas) SetPosition(position image.Point) {
	c.canvasBounds.Min = position
	c.SetBounds(c.bounds)
}

func (c *Canvas) SetSize(size image.Point) {
	c.canvasBounds.Max = c.canvasBounds.Min.Add(size)
	c.SetBounds(c.bounds)
}

func (c *Canvas) SetRect(rect image.Rectangle) {
	c.canvasBounds = rect
	c.SetBounds(c.bounds)
}

func (c *Canvas) SetBounds(bounds Bounds) {
	c.dirtyBounds = true
}

func (c *Canvas) ShouldDraw() bool {
	return c.dirtyBounds || c.ContainerBase.ShouldDraw()
}

func (c *Canvas) Draw() {
	if c.dirtyBounds {
		c.resizeChildren()
	}

	c.ContainerBase.Draw()
	c.ClearFlags()
}

func (c *Canvas) ClearFlags() {
	c.dirtyBounds = false
}

func (c *Canvas) resizeChildren() {
	// Recalculate size in relation to the root

	// Get resolution
	res := c.Root().Bounds().Size
	// Convert from absolute to relative bounds
	relbounds := BoundsFromRect(c.canvasBounds).Scale(res.Inverse())

	// Apply to each children
	for _, child := range c.children {
		child.SetBounds(relbounds)
	}
}
