package components

import (
	"image"

	"github.com/hamcha/youi/opengl"
)

const imageFragShader = `
#version 330 core
uniform sampler2D imgdata;
in vec2 fragTexCoord;
out vec4 color;
void main() {
	color = texture(imgdata, fragTexCoord);
}
` + "\x00"

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
	return i.dirtyContent || i.ComponentBase.ShouldDraw()
}

func (i *Image) Draw() {
	if i.shader == nil {
		i.shader = opengl.DefaultShader()
		err := i.shader.SetFragmentSource(imageFragShader)
		if err != nil {
			panic(err)
		}
	}
	if i.dirtyContent {
		i.texture = opengl.MakeTexture(i.content, opengl.TextureOptions{
			WrapS:     opengl.TextureWrapClamp,
			WrapR:     opengl.TextureWrapClamp,
			MinFilter: opengl.TextureFilterLinear,
			MagFilter: opengl.TextureFilterLinear,
		})
		i.shader.GetUniform("imgdata").Set(i.texture)
	}

	i.ComponentDrawable.Draw()

	i.ClearFlags()
}

func (i *Image) ClearFlags() {
	i.dirtyContent = false
}
