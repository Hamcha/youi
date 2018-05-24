package builtin

import (
	"image"

	"github.com/hamcha/youi/components"
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
	components.ComponentDrawable

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
	if i.Shader == nil {
		i.Shader = opengl.DefaultShader()
		err := i.Shader.SetFragmentSource(imageFragShader)
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
		i.Shader.GetUniform("imgdata").Set(i.texture)
	}

	i.ComponentDrawable.Draw()

	i.ClearFlags()
}

func (i *Image) ClearFlags() {
	i.dirtyContent = false
}

func (i *Image) String() string {
	//TODO Add attributes
	return "<Image />"
}

func makeImage(list components.AttributeList) (components.Component, error) {
	//TODO Parse attributes
	return &Image{}, nil
}
