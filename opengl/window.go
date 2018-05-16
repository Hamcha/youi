package opengl

import (
	"fmt"
	"image"
	"image/color"

	"github.com/go-gl/gl/v3.3-core/gl"
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

	// Initialize OpenGL on the window
	err = gl.Init()
	if err != nil {
		return nil, err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Setup global properties
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	return &Window{
		handle:  window,
		options: options,
	}, err
}

func (w *Window) Clear() {
	if w.options.BackgroundColor != nil {
		SetBGColor(w.options.BackgroundColor)
	}
	Clear()
}

func (w *Window) DrawDone() {
	w.handle.SwapBuffers()
}

func (w *Window) IsOpen() bool {
	return !w.handle.ShouldClose()
}

func (w *Window) GetSize() image.Point {
	width, height := w.handle.GetSize()
	return image.Point{width, height}
}

type WindowResizeCallback func(width, height int)

func (w *Window) SetResizeCallback(fn WindowResizeCallback) {
	w.handle.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		fn(width, height)
	})
}
