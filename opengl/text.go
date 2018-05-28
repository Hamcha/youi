package opengl

import (
	"image/color"

	"github.com/hamcha/youi/font"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	font    *font.Font
	texture *Texture
	content string
	Shader  *Shader
	Mesh    *Mesh
}

func MakeText(fnt *font.Font, text string) *Text {
	vertices, indices := quadFromText(fnt, text)
	texture := MakeTexture(fnt.Texture, TextureOptions{
		WrapS:     TextureWrapClamp,
		WrapR:     TextureWrapClamp,
		MinFilter: TextureFilterLinear,
		MagFilter: TextureFilterLinear,
	})
	fontShader := getFontShader(texture, fnt)
	mesh := MakeMesh(vertices, indices, fontShader)
	return &Text{
		font:    fnt,
		content: text,
		texture: texture,
		Shader:  fontShader,
		Mesh:    mesh,
	}
}

func (t *Text) SetContent(text string) {
	t.Mesh.Destroy()
	vertices, indices := quadFromText(t.font, text)
	t.Mesh = MakeMesh(vertices, indices, t.Shader)
}

func (t *Text) SetColor(col color.Color) {
	t.Shader.GetUniform("fontColor").Set(col)
}

const fontVertexShader = `
#version 330 core
uniform mat4 transform;
uniform float ttfScale;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
	fragTexCoord = vertTexCoord;
	gl_Position = transform * vec4(vert * ttfScale * 0.2, 1);
}
` + "\x00"

const fontFragmentShader = `
#version 330 core
uniform sampler2D fontTexture;
uniform vec4 fontColor;
in vec2 fragTexCoord;
out vec4 color;

const float tolerance = 0.05;

void main() {
	float distance = texture(fontTexture, fragTexCoord).a;
	float w1 = smoothstep(0.5-tolerance, 0.5+tolerance, distance);
	float w2 = smoothstep(0.5-tolerance*2, 0.5+tolerance*2, distance);
	float alpha = (w1 + w2) / 2.0;
	color = vec4(fontColor.rgb, fontColor.a * alpha);
}
` + "\x00"

func getFontShader(texture *Texture, font *font.Font) *Shader {
	// Make shader
	fontShader := DefaultShader()

	// Set custom shaders
	err := fontShader.SetFragmentSource(fontFragmentShader)
	if err != nil {
		panic(err)
	}
	err = fontShader.SetVertexSource(fontVertexShader)
	if err != nil {
		panic(err)
	}

	fontShader.GetUniform("fontTexture").Set(texture)
	fontShader.GetUniform("fontColor").Set(color.White)
	fontShader.GetUniform("ttfScale").Set(1.0 / float32(font.Size))

	// Return it
	return fontShader
}

func quadFromText(fnt *font.Font, text string) (vertices []float32, indices []uint32) {
	//TODO Word wrapping

	// Each glyph is a quad
	// Each quad is 2 triangles / 6 vertices
	// Each vertex is 5 values
	// X1 Y1 Z1 U1 V1 (top left)
	// X2 Y2 Z2 U2 V2 (top right)
	// X3 Y3 Z3 U3 V3 (bottom left)
	// X4 Y4 Z4 U4 V4 (bottom right)

	// Get texture size
	tsize := fnt.Texture.Bounds().Size()
	twidth, theight := float32(tsize.X), float32(tsize.Y)

	// Make arrays
	vertices = make([]float32, 20*len(text))
	indices = make([]uint32, 6*len(text))

	// Get font scale and other TTF parameters
	fscale := fixed.Int26_6(fnt.Size << 6)
	prevCharIndex, isFirst := fnt.TTF.Index(0), true

	curx := float32(0)
	cury := float32(0)
	for index, chr := range text {
		// Increase space by whatever kerning is
		curCharIndex := fnt.TTF.Index(chr)
		if !isFirst {
			kern := fnt.TTF.Kern(fscale, prevCharIndex, curCharIndex)
			//kern = (kern + 32) &^ 63 // ??
			curx += float32(kern) / 64
		}

		// Get quad width/height
		bounds := fnt.Atlas[chr]
		size := bounds.Size()

		leftu := float32(bounds.Min.X) / twidth
		rightu := float32(bounds.Max.X) / twidth
		topv := float32(bounds.Min.Y) / theight
		bottomv := float32(bounds.Max.Y) / theight

		vertidx := index * 20
		vertices[vertidx+0] = curx                    // X1
		vertices[vertidx+1] = cury                    // Y1
		vertices[vertidx+3] = leftu                   // U1
		vertices[vertidx+4] = bottomv                 // V1
		vertices[vertidx+5] = curx + float32(size.X)  // X2
		vertices[vertidx+6] = cury                    // Y2
		vertices[vertidx+8] = rightu                  // U2
		vertices[vertidx+9] = bottomv                 // V2
		vertices[vertidx+10] = curx                   // X3
		vertices[vertidx+11] = cury + float32(size.Y) // Y3
		vertices[vertidx+13] = leftu                  // U3
		vertices[vertidx+14] = topv                   // V3
		vertices[vertidx+15] = curx + float32(size.X) // X4
		vertices[vertidx+16] = cury + float32(size.Y) // Y4
		vertices[vertidx+18] = rightu                 // U4
		vertices[vertidx+19] = topv                   // V4

		ididx := index * 6
		basevtx := uint32(index * 4)
		indices[ididx+0] = basevtx
		indices[ididx+1] = basevtx + 1
		indices[ididx+2] = basevtx + 2
		indices[ididx+3] = basevtx + 1
		indices[ididx+4] = basevtx + 2
		indices[ididx+5] = basevtx + 3

		// Get font metrics for advancement
		metrics := fnt.TTF.HMetric(fscale, curCharIndex)
		advanceWidth := float32(metrics.AdvanceWidth) / 64

		curx += advanceWidth

		// Set index as previous
		prevCharIndex, isFirst = curCharIndex, false
	}

	return
}
