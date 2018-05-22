package components

import (
	"errors"
)

// Container is a component that can hold (and partition) other components
type Container interface {
	Component

	Children() ComponentList
	Root() Container
}

// ComponentList is a modifiable, ordered list of components
type ComponentList []Component

// ContainerBase is the common parent of all container components
type ContainerBase struct {
	ComponentBase

	children      ComponentList
	dirtyChildren bool
}

// Children returns the list of all children components
func (c *ContainerBase) Children() ComponentList {
	return c.children
}

func (c *ContainerBase) drawChildren() {
	for _, child := range c.children {
		//if child.ShouldDraw() {
		child.Draw()
		//}
	}
}

func (c *ContainerBase) ClearFlags() {
	c.dirtyChildren = false
}

func (c *ContainerBase) Root() Container {
	if c.parent == nil {
		return c
	}
	return c.parent.Root()
}

func (c *ContainerBase) isRoot() bool {
	return c.parent == nil
}

func (c *ContainerBase) ShouldDraw() bool {
	return c.dirtyChildren
}

func (c *ContainerBase) Draw() {
	c.drawChildren()
	c.ClearFlags()
}

// AppendChild adds a component at the end of the list
func (c *ContainerBase) AppendChild(component Component) {
	c.children = append(c.children, component)
	component.setParent(c)
	c.dirtyChildren = true
}

// InsertChild inserts a component in one index of the list, moving all other components forward
func (c *ContainerBase) InsertChild(component Component, index int) error {
	if index < 0 || index >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:index], append(ComponentList{component}, c.children[index:]...)...)
	component.setParent(c)
	c.dirtyChildren = true
	return nil
}

// PrependChild inserts a component at the beginning of the list
func (c *ContainerBase) PrependChild(component Component) {
	c.children = append(ComponentList{component}, c.children...)
	component.setParent(c)
	c.dirtyChildren = true
}

// RemoveChild removes a component from the list
func (c *ContainerBase) RemoveChild(component Component) error {
	id := c.FindChildIndex(component)
	if id < 0 {
		return ErrComponentNotFound
	}
	return c.RemoveChildByIndex(id)
}

// FindChildIndex finds a component's index in the list
func (c *ContainerBase) FindChildIndex(component Component) int {
	for i, cmp := range c.children {
		if cmp == component {
			return i
		}
	}
	return -1
}

// RemoveChildByIndex removes the ith component from the list
func (c *ContainerBase) RemoveChildByIndex(i int) error {
	if i < 0 || i >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:i], c.children[i+1:]...)
	c.dirtyChildren = true
	return nil
}

// Component handling errors
var (
	ErrIndexOutOfBounds  = errors.New("index out of bounds")
	ErrComponentNotFound = errors.New("component not found in list")
)
