package opengl

import (
	"image/color"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	handle  *glfw.Window
	options WindowOptions
}

type WindowOptions struct {
	Resizable       bool
	BackgroundColor color.Color
}

func CreateWindow(width, height int, title string, monitor *glfw.Monitor, parent *Window, options WindowOptions) (*Window, error) {
	// Get parent window handle if a parent was specified
	var parentWnd *glfw.Window
	if parent != nil {
		parentWnd = parent.handle
	}

	// Setup options
	resizable := glfw.False
	if options.Resizable {
		resizable = glfw.True
	}
	glfw.WindowHint(glfw.Resizable, resizable)

	// Create window
	window, err := glfw.CreateWindow(width, height, title, monitor, parentWnd)
	window.MakeContextCurrent()

	return &Window{
		handle:  window,
		options: options,
	}, err
}

func (w *Window) Draw() {
	//TODO Only draw if needed
	if w.options.BackgroundColor != nil {
		SetBGColor(w.options.BackgroundColor)
	}
	Clear()

	//TODO Actually draw something??

	w.handle.SwapBuffers()
}

func (w *Window) IsOpen() bool {
	return !w.handle.ShouldClose()
}
