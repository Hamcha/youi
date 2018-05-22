package opengl

import "github.com/go-gl/gl/v3.3-core/gl"

var quadVertices = []float32{
	// Triangle 1
	-1, -1, 0, 0, 1,
	-1, 1, 0, 0, 0,
	1, 1, 0, 1, 0,
	// Triangle 2
	-1, -1, 0, 0, 1,
	1, -1, 0, 1, 1,
	1, 1, 0, 1, 0,
}

// Quad is a simple rectangle mesh made of two triangles
type Quad struct {
	vao    uint32
	vbo    uint32
	Shader *Shader
}

// MakeQuad creates a quad with either a provided shader or a default one
func MakeQuad(shader *Shader) *Quad {
	quad := new(Quad)

	// Fallback to default shader if not specified
	if shader == nil {
		shader = DefaultShader()
	}

	program := shader.MustGetProgram()

	quad.Shader = shader

	// Generate vertex array object
	gl.GenVertexArrays(1, &quad.vao)
	gl.BindVertexArray(quad.vao)

	// Generate vertex buffer object
	gl.GenBuffers(1, &quad.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, quad.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, gl.Ptr(quadVertices), gl.STATIC_DRAW)

	// Setup shader uniforms
	vertAttrib := uint32(gl.GetAttribLocation(program, glString("vert")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, glString("vertTexCoord")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return quad
}

// Draw sets the quad's shader and draws it
func (q *Quad) Draw() {
	// Load shader
	q.Shader.MustGetProgram() // Make sure it's updated
	q.Shader.Use()

	// Setup uniforms
	q.Shader.BindUniforms()

	// Bind VAO
	gl.BindVertexArray(q.vao)

	// Draw vertices
	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}
