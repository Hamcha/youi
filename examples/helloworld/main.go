package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
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
		DebugContext:    true,
	})
	if err != nil {
		panic(err)
	}

	//label := components.Label{}
	//label.SetFontSize(12)
	//label.SetText("Hello world")

	imgfile, _ := os.Open("out.png")
	imghelo, _ := png.Decode(imgfile)

	imgdata := image.NewRGBA(imghelo.Bounds())
	if imgdata.Stride != imgdata.Rect.Size().X*4 {
		panic(fmt.Errorf("unsupported stride"))
	}
	draw.Draw(imgdata, imgdata.Bounds(), imghelo, image.ZP, draw.Src)

	img := components.Image{}
	img.SetImage(imgdata)

	canvas := components.Canvas{}
	//canvas.AppendChild(&label)
	canvas.AppendChild(&img)
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
