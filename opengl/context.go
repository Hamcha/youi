package opengl

import (
	"image/color"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const SYSDPI = 72 //TODO Support customizable DPI

// Init initializes an OpenGL context
func Init() error {
	// Setup GLFW
	glfw.Init()
	//TODO Should check for GLFW errors

	// Setup GLFW properties
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	return nil
}

// Terminate cleanly closes all the currently opened windows, frees resources etc.
func Terminate() {
	glfw.Terminate()
}

// SetBGColor sets the background (clear) color
func SetBGColor(col color.Color) {
	r, g, b, a := toGLColor(col)
	gl.ClearColor(r, g, b, a)
}

// Clear clears the screen
func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func Poll() {
	glfw.PollEvents()
}
