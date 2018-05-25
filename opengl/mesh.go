package opengl

import "github.com/go-gl/gl/v3.3-core/gl"

// Mesh is a mesh that can be drawn
type Mesh struct {
	vertices []float32
	vao      uint32
	vbo      uint32
	Shader   *Shader
}

// MakeMesh creates a mesh with given vertices and an optional shader
func MakeMesh(vertices []float32, shader *Shader) *Mesh {
	// Fallback to default shader if not specified
	if shader == nil {
		shader = DefaultShader()
	}

	mesh := new(Mesh)
	mesh.Shader = shader
	mesh.vertices = vertices

	// Generate vertex array object
	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	// Generate vertex buffer object
	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	return mesh
}

// Draw sets the quad's shader and draws it
func (m *Mesh) Draw() {
	// Load shader
	m.Shader.MustGetProgram() // Make sure it's updated
	m.Shader.Use()

	// Setup uniforms
	m.Shader.BindUniforms()

	// Bind VAO
	gl.BindVertexArray(m.vao)

	// Draw vertices
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.vertices)/5))
}

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

// MakeQuad creates a quad with either a provided shader or a default one
func MakeQuad(shader *Shader) *Mesh {
	return MakeMesh(quadVertices, shader)
}
