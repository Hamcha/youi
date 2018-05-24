package youi

import (
	"image"
	"io"

	"github.com/kataras/go-errors"

	"github.com/hamcha/youi/components/builtin"
	"github.com/hamcha/youi/opengl"
)

var (
	ErrYUMLRootMustBePage = errors.New("YUML root must be <Page>")
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
	yumlElem, err := parseYUML(reader)
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
