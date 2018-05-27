package opengl

import "github.com/hamcha/youi/font"

type Text struct {
	font    font.Font
	content string
	mesh    *Mesh
}

func MakeText(fnt font.Font, text string) *Text {
	mesh := MakeMesh(quadFromText(fnt, text), getFontShader())
	return &Text{
		font:    fnt,
		content: text,
		mesh:    mesh,
	}
}

func (t *Text) SetContent(newstr string) {
	t.mesh.Destroy()
	t.mesh = MakeMesh(quadFromText(t.font, newstr), getFontShader())
}

var fontShader *Shader

func getFontShader() *Shader {
	if fontShader != nil {
		return fontShader
	}
	// Make shader
	return nil
}

func quadFromText(fnt font.Font, text string) []float32 {

	return nil
}
