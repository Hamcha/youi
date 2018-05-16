package components

import (
	"errors"
)

// Container is a component that can hold (and partition) other components
type Container interface {
	Component

	Children() ComponentList
}

// ComponentList is a modifiable, ordered list of components
type ComponentList []Component

// containerBase is the common parent of all container components
type containerBase struct {
	children ComponentList

	dirtyChildren bool
}

func (c *containerBase) Children() ComponentList {
	return c.children
}

func (c *containerBase) drawChildren() {
	for _, child := range c.children {
		if child.ShouldDraw() {
			child.Draw()
		}
	}
}

func (c *containerBase) isDirty() bool {
	return c.dirtyChildren
}

// Append adds a component at the end of the list
func (c *containerBase) AppendChild(component Component) {
	c.children = append(c.children, component)
	c.dirtyChildren = true
}

// Insert inserts a component in one index of the list, moving all other components forward
func (c *containerBase) InsertChild(component Component, index int) error {
	if index < 0 || index >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:index], append(ComponentList{component}, c.children[index:]...)...)
	c.dirtyChildren = true
	return nil
}

// PrependChildren inserts a component at the beginning of the list
func (c *containerBase) PrependChild(component Component) {
	c.children = append(ComponentList{component}, c.children...)
	c.dirtyChildren = true
}

// RemoveChildren removes a component from the list
func (c *containerBase) RemoveChild(component Component) error {
	id := c.FindChildIndex(component)
	if id < 0 {
		return ErrComponentNotFound
	}
	return c.RemoveChildByIndex(id)
}

// FindChildIndex finds a component's index in the list
func (c *containerBase) FindChildIndex(component Component) int {
	for i, cmp := range c.children {
		if cmp == component {
			return i
		}
	}
	return -1
}

// RemoveChildByIndex removes the ith component from the list
func (c *containerBase) RemoveChildByIndex(i int) error {
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
