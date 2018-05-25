package main

import (
	"image"
	"runtime"

	"github.com/hamcha/youi"
	"github.com/hamcha/youi/components/builtin"
	"github.com/hamcha/youi/opengl"
	"github.com/hamcha/youi/utils"
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

	window, err := opengl.CreateWindow(640, 480, "Hello world", nil, nil, opengl.WindowOptions{
		BackgroundColor: utils.HexColor(0x101020ff),
		Resizable:       true,
		DebugContext:    true,
	})
	if err != nil {
		panic(err)
	}

	label := builtin.Label{}
	label.SetFontSize(12)
	label.SetText("Hello world")

	canvas := builtin.Canvas{}
	canvas.AppendChild(&label)
	canvas.SetRect(image.Rect(10, 10, 110, 110))

	form := youi.MakeForm(window)
	form.Root.AppendChild(&canvas)

	for window.IsOpen() {
		if form.Root.ShouldDraw() {
			form.Draw()
		}
		opengl.Poll()
	}
}
