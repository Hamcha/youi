package youi

import (
	"image"

	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/opengl"
)

type Form struct {
	Root *components.Canvas

	window *opengl.Window
}

func MakeForm(window *opengl.Window) *Form {
	form := &Form{
		Root:   new(components.Canvas),
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
