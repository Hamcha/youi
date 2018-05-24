package components

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/hamcha/youi/font"
	"github.com/hamcha/youi/opengl"
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

// ComponentBase is the common parent of all components
type ComponentBase struct {
	parent Component

	bounds      Bounds
	dirtyBounds bool

	children      ComponentList
	dirtyChildren bool
}

// ComponentList is a modifiable, ordered list of components
type ComponentList []Component

func (c *ComponentBase) Parent() Component {
	return c.parent
}

func (c *ComponentBase) setParent(container Component) {
	c.parent = container
}

func (c *ComponentBase) SetBounds(bounds Bounds) {
	c.bounds = bounds
	c.dirtyBounds = true
}

func (c *ComponentBase) Bounds() Bounds {
	return c.bounds
}

func (c *ComponentBase) Draw() {
	c.drawChildren()
	c.ClearFlags()
}

func (c *ComponentBase) ShouldDraw() bool {
	return c.dirtyBounds || c.dirtyChildren
}

func (c *ComponentBase) SetRedraw() {
	c.dirtyBounds = true
}

func (c *ComponentBase) ClearFlags() {
	c.dirtyChildren = false
	c.dirtyBounds = false
}

// Children returns the list of all children components
func (c *ComponentBase) Children() ComponentList {
	return c.children
}

func (c *ComponentBase) drawChildren() {
	for _, child := range c.children {
		//if child.ShouldDraw() {
		child.Draw()
		//}
	}
}

func (c *ComponentBase) Root() Component {
	if c.parent == nil {
		return c
	}
	return c.parent.Root()
}

func (c *ComponentBase) isRoot() bool {
	return c.parent == nil
}

// AppendChild adds a component at the end of the list
func (c *ComponentBase) AppendChild(component Component) {
	c.children = append(c.children, component)
	component.setParent(c)
	c.dirtyChildren = true
}

// InsertChild inserts a component in one index of the list, moving all other components forward
func (c *ComponentBase) InsertChild(component Component, index int) error {
	if index < 0 || index >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:index], append(ComponentList{component}, c.children[index:]...)...)
	component.setParent(c)
	c.dirtyChildren = true
	return nil
}

// PrependChild inserts a component at the beginning of the list
func (c *ComponentBase) PrependChild(component Component) {
	c.children = append(ComponentList{component}, c.children...)
	component.setParent(c)
	c.dirtyChildren = true
}

// RemoveChild removes a component from the list
func (c *ComponentBase) RemoveChild(component Component) error {
	id := c.FindChildIndex(component)
	if id < 0 {
		return ErrComponentNotFound
	}
	return c.RemoveChildByIndex(id)
}

// FindChildIndex finds a component's index in the list
func (c *ComponentBase) FindChildIndex(component Component) int {
	for i, cmp := range c.children {
		if cmp == component {
			return i
		}
	}
	return -1
}

// RemoveChildByIndex removes the ith component from the list
func (c *ComponentBase) RemoveChildByIndex(i int) error {
	if i < 0 || i >= len(c.children) {
		return ErrIndexOutOfBounds
	}
	c.children = append(c.children[:i], c.children[i+1:]...)
	c.dirtyChildren = true
	return nil
}

// ChildrenStr calls String() on each children and indents the results
func (c *ComponentBase) ChildrenStr() (out string) {
	for _, child := range c.children {
		out += utils.IndentStrings(fmt.Sprint(child), 1) + "\n"
	}
	return
}

// ComponentText is a common parent of all text-based components
type ComponentText struct {
	fontFace string
	fontSize float64

	_font     *font.Font
	dirtyFont bool
}

func (c *ComponentText) SetFontFace(name string) {
	c.fontFace = name
	c.dirtyFont = true
}

func (c *ComponentText) SetFontSize(size float64) {
	c.fontSize = size
}

func (c *ComponentText) makeFace() {
	// If no font is provided, use Go Regolar
	if c.fontFace == "" {
		c._font = font.DefaultFont()
	} else {
		var err error
		c._font, err = font.LoadFont(c.fontFace)
		if err != nil {
			//TODO Proper error reporting
			panic(err)
		}
	}
}

func (c *ComponentText) ShouldDraw() bool {
	return c.dirtyFont
}

func (c *ComponentText) ClearFlags() {
	c.dirtyFont = false
}

func (c *ComponentText) Draw() {
	if c.dirtyFont {
		c.makeFace()
		c.ClearFlags()
	}
}

type ComponentDrawable struct {
	ComponentBase
	Quad   *opengl.Quad
	Shader *opengl.Shader
}

func (c *ComponentDrawable) Draw() {
	if c.Shader == nil {
		c.Shader = opengl.DefaultShader()
	}
	if c.Quad == nil {
		c.Quad = opengl.MakeQuad(c.Shader)
	}

	// Check if bounds have changed
	if c.dirtyBounds {
		// Update transform matrix
		c.updateTransformMatrix()
	}

	c.Quad.Draw()
	c.ComponentBase.Draw()
}

func (c *ComponentDrawable) updateTransformMatrix() {
	// Get size and position
	size := c.bounds.Size
	pos := c.bounds.Position

	// Change pivot to top-left instead of center
	pos = Position{pos.X + size.Width/2, pos.Y + size.Height/2}

	// Set position, taking into account OpenGL's weird coordinate system
	posMtx := mgl32.Translate3D((pos.X*2)-1, -(pos.Y*2)+1, 0.0)

	// Scale to size
	sizeMtx := mgl32.Scale3D(size.Width, size.Height, 1.0)

	// Multiply everything into a transform matrix
	result := posMtx.Mul4(sizeMtx)

	// Set result matrix as uniform value
	c.Shader.GetUniform("transform").Set(result)
}

// Component handling errors
var (
	ErrIndexOutOfBounds  = errors.New("index out of bounds")
	ErrComponentNotFound = errors.New("component not found in list")
)
