package builtin

import (
	"image"

	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/loader"
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
	components.Drawable

	src          string
	content      *image.RGBA
	dirtyContent bool
	texture      *opengl.Texture
}

func (i *Image) SetPath(src string) error {
	i.src = src

	// Load image, then set it
	img, err := loader.Image(src)
	if err != nil {
		return err
	}

	i.SetImage(img)
	return nil
}

func (i *Image) SetImage(img *image.RGBA) {
	i.content = img
	i.dirtyContent = true
}

func (i *Image) ShouldDraw() bool {
	return i.dirtyContent || i.Drawable.ShouldDraw()
}

func (i *Image) Draw() {
	if i.Mesh == nil {
		shader := opengl.DefaultShader()
		err := shader.SetFragmentSource(imageFragShader)
		if err != nil {
			panic(err)
		}
		i.Mesh = opengl.MakeQuad(shader)
	}
	if i.dirtyContent {
		i.texture = opengl.MakeTexture(i.content, opengl.TextureOptions{
			WrapS:     opengl.TextureWrapClamp,
			WrapR:     opengl.TextureWrapClamp,
			MinFilter: opengl.TextureFilterLinear,
			MagFilter: opengl.TextureFilterLinear,
		})
		i.Mesh.Shader.GetUniform("imgdata").Set(i.texture)
	}

	i.Drawable.Draw()

	i.ClearFlags()
}

func (i *Image) ClearFlags() {
	i.dirtyContent = false
}

func (i *Image) String() string {
	return "<Image Path=\"" + i.src + "\" />"
}

func makeImage(list components.AttributeList) (components.Component, error) {
	img := &Image{}

	src := list.Get("Path", "")
	if src != "" {
		err := img.SetPath(src.String())
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}
