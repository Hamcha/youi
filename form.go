package youi

import (
	"image"
	"io"

	"github.com/kataras/go-errors"

	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/components/builtin"
	"github.com/hamcha/youi/opengl"
	"github.com/hamcha/youi/yuml"
)

var (
	ErrYUMLRootMustBePage  = errors.New("YUML root must be <Page>")
	ErrCouldNotMakeElement = errors.New("Could not make element \"%s\"")
)

type Form struct {
	Root *builtin.Page

	window *opengl.Window
}

func MakeForm(window *opengl.Window) *Form {
	form := &Form{
		Root:   new(builtin.Page),
		window: window,
	}
	form.Root.SetSize(window.GetSize())

	// Set resize callback
	window.SetResizeCallback(form.onResize)

	return form
}

func (f *Form) Draw() {
	f.window.Clear()
	f.Root.Draw()
	f.window.DrawDone()
}

func (f *Form) onResize(width, height int) {
	f.Root.SetSize(image.Point{width, height})
}

func (f *Form) LoadYUML(reader io.Reader) error {
	// Parse code
	yumlElem, err := yuml.ParseYUML(reader)
	if err != nil {
		return err
	}

	// Create tree
	elem, err := makeYUMLcomponentTree(yumlElem)
	if err != nil {
		return err
	}

	// Add to root and return
	var ok bool
	f.Root, ok = elem.(*builtin.Page)
	if !ok {
		return ErrYUMLRootMustBePage
	}
	return nil
}

func makeYUMLcomponentTree(element *yuml.Element) (components.Component, error) {
	elem, err := makeComponent(element.Name.Space, element.Name.Local, toAttributeList(element.Attributes))
	if err != nil {
		return nil, err
	}

	// Check for children
	for _, child := range element.Children {
		childelem, err := makeYUMLcomponentTree(child.Element)
		if err != nil {
			return nil, ErrCouldNotMakeElement.Format(element.Name.Local).AppendErr(err)
		}
		elem.AppendChild(childelem)
	}

	return elem, nil
}

func toAttributeList(y yuml.Attributes) components.AttributeList {
	out := make(components.AttributeList)
	for _, attr := range y {
		out[attr.Name.Local] = components.Attribute(attr.Value)
	}
	return out
}
