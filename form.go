package youi

import (
	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/opengl"
)

type Form struct {
	Root *components.Canvas

	window *opengl.Window
}

func MakeForm(window *opengl.Window) *Form {
	canvas := &Form{
		Root:   &components.Canvas{},
		window: window,
	}
	canvas.Root.SetSize(window.GetSize())

	return canvas
}

func (f *Form) Draw() {
	f.window.Clear()
	f.Root.Draw()
	f.window.DrawDone()
}
