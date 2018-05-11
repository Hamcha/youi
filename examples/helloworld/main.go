package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/hamcha/youi"
	"github.com/hamcha/youi/components"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Hello world", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	youi.Init()
	youi.SetBGColor(youi.HexColor(0xff0000ff))

	label := components.Label{
		x:       320,
		y:       240,
		content: "Hello world",
	}

	for !window.ShouldClose() {
		youi.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
