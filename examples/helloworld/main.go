package main

import (
	"image"
	"image/draw"
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

	window, err := opengl.CreateWindow(640, 480, "Hello world", nil, nil, opengl.WindowOptions{
		BackgroundColor: youi.HexColor(0x101020ff),
		Resizable:       true,
	})
	if err != nil {
		panic(err)
	}

	//label := components.Label{}
	//label.SetFontSize(12)
	//label.SetText("Hello world")

	imgdata := image.NewRGBA(image.Rect(0, 0, 10, 10))
	draw.Draw(imgdata, imgdata.Bounds(), &image.Uniform{youi.HexColor(0x00ff00ff)}, image.ZP, draw.Src)

	img := components.Image{}
	img.SetImage(imgdata)

	canvas := components.Canvas{}
	//canvas.AppendChild(&label)
	canvas.AppendChild(&img)
	canvas.SetBounds(image.Rect(10, 10, 100, 100))

	form := youi.MakeForm(window)
	form.Root.AppendChild(&canvas)

	for window.IsOpen() {
		if form.Root.ShouldDraw() {
			form.Draw()
		}
		opengl.Poll()
	}
}
