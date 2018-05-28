package opengl

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window is a system window with an OpenGL context inside it
type Window struct {
	handle  *glfw.Window
	options WindowOptions
}

// WindowOptions are extra options that can affect how the window looks or behaves
type WindowOptions struct {
	Resizable       bool
	BackgroundColor color.Color
	DebugContext    bool
}

// CreateWindow creates a window with an opengl context and returns it
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

	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Renderer:", renderer)
	fmt.Println("OpenGL version", version)

	// Setup global properties
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if options.DebugContext {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(logstderr, nil)
	}

	return &Window{
		handle:  window,
		options: options,
	}, err
}

// Clear clears the window and its contents
func (w *Window) Clear() {
	if w.options.BackgroundColor != nil {
		SetBGColor(w.options.BackgroundColor)
	}
	Clear()
}

// DrawDone swaps the buffers after drawing
func (w *Window) DrawDone() {
	w.handle.SwapBuffers()
}

// IsOpen returns wether the window is open or should not be closed
func (w *Window) IsOpen() bool {
	return !w.handle.ShouldClose()
}

// GetSize returns the window size
func (w *Window) GetSize() image.Point {
	width, height := w.handle.GetSize()
	return image.Point{width, height}
}

// WindowResizeCallback is a callback that can be called when the window is resized
type WindowResizeCallback func(width, height int)

// SetResizeCallback sets a callback to be called when the window is resized
func (w *Window) SetResizeCallback(fn WindowResizeCallback) {
	w.handle.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		fn(width, height)
	})
}

var debugType = map[uint32]string{
	gl.DEBUG_TYPE_ERROR:               "Error",
	gl.DEBUG_TYPE_MARKER:              "Marker",
	gl.DEBUG_TYPE_PERFORMANCE:         "Performance",
	gl.DEBUG_TYPE_PUSH_GROUP:          "Push group",
	gl.DEBUG_TYPE_POP_GROUP:           "Pop group",
	gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR: "Deprecated behavior",
	gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:  "Undefined behavior",
	gl.DEBUG_TYPE_PORTABILITY:         "Portability",
	gl.DEBUG_TYPE_OTHER:               "Other",
}

var debugSeverity = map[uint32]string{
	gl.DEBUG_SEVERITY_NOTIFICATION: "Notification",
	gl.DEBUG_SEVERITY_LOW:          "Low",
	gl.DEBUG_SEVERITY_MEDIUM:       "Medium",
	gl.DEBUG_SEVERITY_HIGH:         "High",
}

func logstderr(source, mtype, id, severity uint32, length int32, msg string, data unsafe.Pointer) {
	typestr, ok := debugType[mtype]
	if !ok {
		typestr = "Unknown"
	}

	sevstr, ok := debugSeverity[severity]
	if !ok {
		sevstr = "Unknown"
	}

	fmt.Printf("[GL-DEBUG] Type: %s, Severity: %s, Message: %s\n", typestr, sevstr, msg)
}
