package components

import (
	"fmt"

	"github.com/hamcha/youi/utils"
)

// Component is a renderable UI component that can optionally hold children
type Component interface {
	Draw()
	ShouldDraw() bool

	Bounds() Bounds
	SetBounds(Bounds)

	Parent() Component
	Children() ComponentList
	Root() Component

	AppendChild(Component)
	InsertChild(Component, int) error
	PrependChild(Component)
	RemoveChild(Component) error
	FindChildIndex(Component) int
	RemoveChildByIndex(int) error

	setParent(Component)
}

// Base is the common parent of all components
type Base struct {
	parent Component

	bounds      Bounds
	dirtyBounds bool

	children      ComponentList
	dirtyChildren bool
}

// ComponentList is a modifiable, ordered list of components
type ComponentList []Component

func (c *Base) Parent() Component {
	return c.parent
}

func (c *Base) setParent(container Component) {
	c.parent = container
}

func (c *Base) SetBounds(bounds Bounds) {
	c.bounds = bounds
	c.dirtyBounds = true
}

func (c *Base) Bounds() Bounds {
	return c.bounds
}

func (c *Base) Draw() {
	c.drawChildren()
	c.ClearFlags()
}

func (c *Base) ShouldDraw() bool {
	for _, child := range c.children {
		if child.ShouldDraw() {
			return true
		}
	}
	return c.dirtyBounds || c.dirtyChildren
}

func (c *Base) SetRedraw() {
	c.dirtyBounds = true
}

func (c *Base) ClearFlags() {
	c.dirtyChildren = false
	c.dirtyBounds = false
}

// Children returns the list of all children components
func (c *Base) Children() ComponentList {
	return c.children
}

func (c *Base) drawChildren() {
	for _, child := range c.children {
		//if child.ShouldDraw() {
		child.Draw()
		//}
	}
}

func (c *Base) Root() Component {
	if c.parent == nil {
		return c
	}
	return c.parent.Root()
}

func (c *Base) isRoot() bool {
	return c.parent == nil
}

// AppendChild adds a component at the end of the list
func (c *Base) AppendChild(component Component) {
	c.children = append(c.children, component)
	component.setParent(c)
	c.dirtyChildren = true
}

// InsertChild inserts a component in one index of the list, moving all other components forward
func (c *Base) InsertChild(component Component, index int) error {
	if index < 0 || index >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:index], append(ComponentList{component}, c.children[index:]...)...)
	component.setParent(c)
	c.dirtyChildren = true
	return nil
}

// PrependChild inserts a component at the beginning of the list
func (c *Base) PrependChild(component Component) {
	c.children = append(ComponentList{component}, c.children...)
	component.setParent(c)
	c.dirtyChildren = true
}

// RemoveChild removes a component from the list
func (c *Base) RemoveChild(component Component) error {
	id := c.FindChildIndex(component)
	if id < 0 {
		return ErrComponentNotFound
	}
	return c.RemoveChildByIndex(id)
}

// FindChildIndex finds a component's index in the list
func (c *Base) FindChildIndex(component Component) int {
	for i, cmp := range c.children {
		if cmp == component {
			return i
		}
	}
	return -1
}

// RemoveChildByIndex removes the ith component from the list
func (c *Base) RemoveChildByIndex(i int) error {
	if i < 0 || i >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:i], c.children[i+1:]...)
	c.dirtyChildren = true
	return nil
}

// ChildrenStr calls String() on each children and indents the results
func (c *Base) ChildrenStr() (out string) {
	for _, child := range c.children {
		out += utils.IndentStrings(fmt.Sprint(child), 1) + "\n"
	}
	return
}
