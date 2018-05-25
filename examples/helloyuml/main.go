package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hamcha/youi"
	"github.com/hamcha/youi/opengl"
	"github.com/hamcha/youi/utils"

	// Import example data
	_ "github.com/hamcha/youi/examples/example_data"
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

	const src = `<Page xmlns="https://yuml.ovo.ovh/schema/components/1.0">
	<Canvas X="10" Y="10" Width="100" Height="100">
		<Image Path="images/hello.png" />
	</Canvas>
</Page>`

	form := youi.MakeForm(window)
	err = form.LoadYUML(strings.NewReader(src))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Source:\n%s\n\nResult:\n%s\n", src, form.Root)

	for window.IsOpen() {
		if form.Root.ShouldDraw() {
			form.Draw()
		}
		opengl.Poll()
	}
}
