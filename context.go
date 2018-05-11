package youi

import (
	"image/color"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const SYSDPI = 72 //TODO Support customizable DPI

// Init initializes youi and its dependencies
func Init() error {
	// Initialize OpenGL
	err := gl.Init()
	if err != nil {
		return err
	}

	// Setup global properties
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	return nil
}

// SetBGColor sets the background (clear) color
func SetBGColor(col color.Color) {
	r, g, b, a := toGLColor(col)
	gl.ClearColor(r, g, b, a)
}

func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
