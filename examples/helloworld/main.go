// +build example

package main

import (
	"runtime"

	"github.com/hamcha/youi"
	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/opengl"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := opengl.Init()
	if err != nil {
		panic(err)
	}
	defer opengl.Terminate()

	window, err := opengl.CreateWindow(640, 480, "Hello world", opengl.WindowOptions{
		BackgroundColor: youi.HexColor(0xff0000ff),
	})
	if err != nil {
		panic(err)
	}

	label := components.Label{}
	label.SetFontSize(12)
	label.SetText("Hello world")

	for window.IsOpen() {
		window.Draw()
		opengl.Poll()
	}
}
