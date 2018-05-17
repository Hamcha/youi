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
}

func (c *Canvas) SetPosition(position image.Point) {
	c.bounds.Min = position
	c.dirtyBounds = true
	c.resizeChildren()
}

func (c *Canvas) SetSize(size image.Point) {
	c.bounds.Max = c.bounds.Min.Add(size)
	c.dirtyBounds = true
	c.resizeChildren()
}

func (c *Canvas) SetBounds(rect image.Rectangle) {
	c.setBounds(rect)
	c.resizeChildren()
}

func (c *Canvas) ShouldDraw() bool {
	return c.dirtyBounds || c.ContainerBase.isDirty()
}

func (c *Canvas) Draw() {
	c.ContainerBase.Draw()
	c.clearFlags()
}

func (c *Canvas) clearFlags() {
	c.dirtyBounds = false
}

func (c *Canvas) resizeChildren() {
	for _, child := range c.children {
		child.setBounds(c.bounds)
	}
}
