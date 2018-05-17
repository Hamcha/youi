package components

import (
	"image"

	"github.com/hamcha/youi/opengl"
)

// Image is a simple box that can contain an image or any sort of drawable surface
type Image struct {
	ComponentDrawable

	content      *image.RGBA
	dirtyContent bool
	texture      *opengl.Texture
}

func (i *Image) SetImage(img *image.RGBA) {
	i.content = img
	i.dirtyContent = true
}

func (i *Image) ShouldDraw() bool {
	return i.dirtyContent || i.ComponentBase.isDirty()
}

func (i *Image) Draw() {
	if i.dirtyContent {
		i.texture = opengl.MakeTexture(i.content, opengl.TextureOptions{
			WrapS:     opengl.TextureWrapClamp,
			WrapR:     opengl.TextureWrapClamp,
			MinFilter: opengl.TextureFilterLinear,
			MagFilter: opengl.TextureFilterLinear,
		})
	}

	i.ComponentDrawable.Draw()

	i.ClearFlags()
}

func (i *Image) ClearFlags() {
	i.ComponentBase.ClearFlags()
	i.dirtyContent = false
}
