package builtin

import (
	"fmt"
	"image"

	"github.com/hamcha/youi/components"
)

// Page is a special canvas container with less checks
type Page struct {
	components.Base

	pixelBounds image.Rectangle
}

func (r *Page) SetSize(size image.Point) {
	r.pixelBounds.Max = r.pixelBounds.Min.Add(size)
	r.SetBounds(components.BoundsFromRect(r.pixelBounds))
}

func (r *Page) SetRect(rect image.Rectangle) {
	r.pixelBounds = rect
	r.SetBounds(components.BoundsFromRect(r.pixelBounds))
}

func (r *Page) ShouldDraw() bool {
	return r.Base.ShouldDraw()
}

func (r *Page) Draw() {
	r.resizeChildren()
	r.Base.Draw()
}

func (r *Page) resizeChildren() {
	relbounds := components.Bounds{
		components.Position{0, 0},
		components.Size{1, 1},
	}

	// Apply to each children
	for _, child := range r.Children() {
		child.SetBounds(relbounds)
	}
}

func (r *Page) String() string {
	return fmt.Sprintf("<Page>\n%s</Page>", r.ChildrenStr())
}

func makePage(attr components.AttributeList) (components.Component, error) {
	return &Page{}, nil
}
