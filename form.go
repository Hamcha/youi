package youi

import (
	"image"
	"io"

	"github.com/hamcha/youi/opengl"
)

type Form struct {
	Root *Root

	window *opengl.Window
}

func MakeForm(window *opengl.Window) *Form {
	form := &Form{
		Root:   new(Root),
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
	element, err := parseYUML(reader)
	if err != nil {
		return err
	}

	current := f.Root

	elem, err := makeComponent(element.Name.Space, element.Name.Local)
	if err != nil {
		return err
	}

	return nil
}
