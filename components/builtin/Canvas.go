package builtin

import (
	"errors"
	"fmt"
	"image"

	"github.com/hamcha/youi/components"
)

// Canvas is a container that has absolute positioning and sizing regardless of parent/siblings.
// While Canvas can hold multiple childrens, those childrens will share the same position/size
// so they will all overlap each other. It is therefore recommended to only put one component
// inside a Canvas.
type Canvas struct {
	components.Base

	canvasBounds image.Rectangle
}

func (c *Canvas) SetPosition(position image.Point) {
	c.canvasBounds.Min = position
	c.SetBounds(c.Bounds())
}

func (c *Canvas) SetSize(size image.Point) {
	c.canvasBounds.Max = c.canvasBounds.Min.Add(size)
	c.SetBounds(c.Bounds())
}

func (c *Canvas) SetRect(rect image.Rectangle) {
	c.canvasBounds = rect
	c.SetBounds(c.Bounds())
}

func (c *Canvas) SetBounds(bounds components.Bounds) {
	c.SetRedraw()
}

func (c *Canvas) Draw() {
	//if c.Base.ShouldDraw() {
	c.resizeChildren()
	//}

	c.Base.Draw()
}

func (c *Canvas) resizeChildren() {
	// Recalculate size in relation to the root

	// Get resolution
	res := c.Root().Bounds().Size
	fmt.Println(res)
	// Convert from absolute to relative bounds
	relbounds := components.BoundsFromRect(c.canvasBounds).Scale(res.Inverse())

	// Apply to each children
	for _, child := range c.Children() {
		child.SetBounds(relbounds)
	}
}

func (c *Canvas) String() string {
	cpos := c.canvasBounds.Min
	csize := c.canvasBounds.Size()
	return fmt.Sprintf(`<Canvas X="%d" Y="%d" Width="%d" Height="%d">
%s</Canvas>`,
		cpos.X, cpos.Y, csize.X, csize.Y, c.ChildrenStr())
}

func makeCanvas(list components.AttributeList) (components.Component, error) {
	x, err := list.Get("X", "0").Int()
	if err != nil {
		return nil, errors.New("X must be an integer number")
	}

	y, err := list.Get("Y", "0").Int()
	if err != nil {
		return nil, errors.New("Y must be an integer number")
	}

	width, err := list.Get("Width", "0").Int()
	if err != nil {
		return nil, errors.New("Width must be an integer number")
	}

	height, err := list.Get("Height", "0").Int()
	if err != nil {
		return nil, errors.New("Height must be an integer number")
	}

	return &Canvas{
		canvasBounds: image.Rect(x, y, x+width, y+height),
	}, nil
}
