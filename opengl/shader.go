package opengl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader errors
var (
	ErrMissingVertexShader   = errors.New("vertex shader is missing")
	ErrMissingFragmentShader = errors.New("fragment shader is missing")
	ErrUniformInvalidType    = errors.New("uniform set value is not of a valid/supported type")
)

// Shader is a OpenGL program with one vertex and one fragment shader
type Shader struct {
	fragID, vertID uint32
	programID      uint32

	vertAttrib     uint32
	texCoordAttrib uint32

	uniforms map[string]*Uniform

	output *uint8
}

// MakeShader creates a shader object
func MakeShader() *Shader {
	return &Shader{
		uniforms: make(map[string]*Uniform),
	}
}

// SetFragmentSource sets (and updates) the fragment portion of the shader
func (s *Shader) SetFragmentSource(src string) (err error) {
	// Do some cleanup if we're replacing an existing shader
	s.destroyFragmentShader()

	// Set new shader
	s.fragID, err = setShader(src, gl.FRAGMENT_SHADER)

	// If successful and applicable, attach and compile
	if err == nil && s.programID != 0 {
		gl.AttachShader(s.programID, s.fragID)
		err = s.compileProgram()
	}

	return
}

// SetVertexSource sets (and updates) the vertex portion of the shader
func (s *Shader) SetVertexSource(src string) (err error) {
	// Do some cleanup if we're replacing an existing shader
	s.destroyVertexShader()

	// Set new shader
	s.vertID, err = setShader(src, gl.VERTEX_SHADER)

	// If successful and applicable, attach and compile
	if err == nil && s.programID != 0 {
		gl.AttachShader(s.programID, s.vertID)
		err = s.compileProgram()
	}

	return
}

// GetProgram compiles (if needed) and returns an OpenGL program with the specified shaders
func (s *Shader) GetProgram() (uint32, error) {
	// Program exists already? Just return it
	if s.programID != 0 {
		return s.programID, nil
	}

	// Check that both shaders are present
	if s.vertID == 0 {
		return 0, ErrMissingVertexShader
	}
	if s.fragID == 0 {
		return 0, ErrMissingFragmentShader
	}

	// Create program
	s.programID = gl.CreateProgram()

	// Attach shaders
	gl.AttachShader(s.programID, s.vertID)
	gl.AttachShader(s.programID, s.fragID)

	return s.programID, s.compileProgram()
}

// MustGetProgram calls GetProgram and panics if it gets an error
func (s *Shader) MustGetProgram() uint32 {
	out, err := s.GetProgram()
	if err != nil {
		panic(err)
	}
	return out
}

// Use sets the shader program as active
func (s *Shader) Use() {
	gl.UseProgram(s.programID)
}

// BindUniforms sets all the uniforms to their current set value
func (s *Shader) BindUniforms() {
	for _, uniform := range s.uniforms {
		uniform.Bind()
	}

	// Bind output
	gl.BindFragDataLocation(s.programID, 0, s.output)
}

// SetOutput sets the shader color output variable
func (s *Shader) SetOutput(out string) {
	s.output = glString(out)
}

// Destroy frees up the resources used by the shader and makes it unusable
func (s *Shader) Destroy() {
	s.destroyVertexShader()
	s.destroyFragmentShader()
	if s.programID != 0 {
		gl.DeleteProgram(s.programID)
		s.programID = 0
	}
}

func (s *Shader) destroyVertexShader() {
	if s.vertID != 0 {
		// Already in a program? Detach
		if s.programID != 0 {
			gl.DetachShader(s.programID, s.vertID)
		}
		// Delete old shader
		gl.DeleteShader(s.vertID)
		s.vertID = 0
	}
}

func (s *Shader) destroyFragmentShader() {
	if s.fragID != 0 {
		// Already in a program? Detach
		if s.programID != 0 {
			gl.DetachShader(s.programID, s.fragID)
		}
		// Delete old shader
		gl.DeleteShader(s.fragID)
		s.fragID = 0
	}
}

func (s *Shader) compileProgram() error {
	// Link and check status
	gl.LinkProgram(s.programID)
	var success int32
	gl.GetProgramiv(s.programID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		// Get error
		var logLength int32
		gl.GetProgramiv(s.programID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.programID, logLength, nil, gl.Str(log))

		// Cleanup
		gl.DetachShader(s.programID, s.vertID)
		gl.DetachShader(s.programID, s.fragID)
		gl.DeleteProgram(s.programID)
		s.programID = 0

		return fmt.Errorf("link failed: %v", log)
	}

	// Reset attributes
	if s.vertAttrib != 0 {
		gl.DisableVertexAttribArray(s.vertAttrib)
	}
	s.vertAttrib = uint32(gl.GetAttribLocation(s.programID, glString("vert")))
	gl.EnableVertexAttribArray(s.vertAttrib)
	gl.VertexAttribPointer(s.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	if s.texCoordAttrib != 0 {
		gl.DisableVertexAttribArray(s.texCoordAttrib)
	}
	s.texCoordAttrib = uint32(gl.GetAttribLocation(s.programID, glString("vertTexCoord")))
	gl.EnableVertexAttribArray(s.texCoordAttrib)
	gl.VertexAttribPointer(s.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Reset all uniforms
	for uname := range s.uniforms {
		s.uniforms[uname].id = -1
		s.uniforms[uname].program = s.programID
	}

	return nil
}

func setShader(src string, shaderType uint32) (uint32, error) {
	// Create shader
	shaderid := gl.CreateShader(shaderType)

	// Load shader source
	fragSrcs, free := gl.Strs(src)
	gl.ShaderSource(shaderid, 1, fragSrcs, nil)
	free()

	// Compile and check result
	gl.CompileShader(shaderid)
	var success int32
	gl.GetShaderiv(shaderid, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		// Get error
		var logLength int32
		gl.GetShaderiv(shaderid, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderid, logLength, nil, gl.Str(log))
		gl.DeleteShader(shaderid)

		return 0, fmt.Errorf("compile failed: %v", log)
	}

	return shaderid, nil
}

// DefaultShader returns a plain, kinda useless, shader
func DefaultShader() *Shader {
	shader := MakeShader()

	err := shader.SetVertexSource(DefaultVertexShader)
	if err != nil {
		// Default vertex shader should always be fine
		panic(err)
	}

	err = shader.SetFragmentSource(DefaultFragmentShader)
	if err != nil {
		// Default fragment shader should always be fine
		panic(err)
	}

	// Init transform uniform to identity
	transform := mgl32.Ident4()
	shader.GetUniform("transform").Set(transform)

	// Set output color
	shader.SetOutput("color")

	return shader
}

// DefaultVertexShader is the vertex portion of the default shader
const DefaultVertexShader = `
#version 330 core
uniform mat4 transform;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
	fragTexCoord = vertTexCoord;
	gl_Position = transform * vec4(vert, 1);
}
` + "\x00"

// DefaultFragmentShader is the fragment portion of the default shader
const DefaultFragmentShader = `
#version 330 core
in vec2 fragTexCoord;
out vec4 color;
void main() {
    color = vec4(fragTexCoord.xy,1.0,1.0);
}
` + "\x00"

// GetUniform returns the uniform value for a given shader uniform variable
func (s *Shader) GetUniform(str string) *Uniform {
	if uid, ok := s.uniforms[str]; ok {
		return uid
	}
	s.uniforms[str] = &Uniform{
		name:    str,
		program: s.programID,
	}
	return s.uniforms[str]
}

// Uniform is a modifiable input for shaders
type Uniform struct {
	name    string
	id      int32
	program uint32
	value   interface{}
}

// Set sets the value for the uniform that will applied when bound
func (u *Uniform) Set(value interface{}) {
	u.value = value
}

// Bind binds the current value to the current program
func (u *Uniform) Bind() {
	if u.id < 0 {
		u.id = gl.GetUniformLocation(u.program, glString(u.name))
		if u.id < 0 {
			panic(fmt.Errorf("glGetUniformLocation returned -1, GL error: %d", gl.GetError()))
		}
	}
	switch value := u.value.(type) {
	case uint32:
		gl.Uniform1ui(u.id, value)
	case []uint32:
		switch len(value) {
		case 1:
			gl.Uniform1uiv(u.id, 1, &value[0])
		case 2:
			gl.Uniform2uiv(u.id, 2, &value[0])
		case 3:
			gl.Uniform3uiv(u.id, 3, &value[0])
		case 4:
			gl.Uniform4uiv(u.id, 4, &value[0])
		default:
			panic(ErrUniformInvalidType)
		}
	case int32:
		gl.Uniform1i(u.id, value)
	case []int32:
		switch len(value) {
		case 1:
			gl.Uniform1iv(u.id, 1, &value[0])
		case 2:
			gl.Uniform2iv(u.id, 2, &value[0])
		case 3:
			gl.Uniform3iv(u.id, 3, &value[0])
		case 4:
			gl.Uniform4iv(u.id, 4, &value[0])
		default:
			panic(ErrUniformInvalidType)
		}
	case float32:
		gl.Uniform1f(u.id, value)
	case []float32:
		switch len(value) {
		case 1:
			gl.Uniform1fv(u.id, 1, &value[0])
		case 2:
			gl.Uniform2fv(u.id, 2, &value[0])
		case 3:
			gl.Uniform3fv(u.id, 3, &value[0])
		case 4:
			gl.Uniform4fv(u.id, 4, &value[0])
		default:
			panic(ErrUniformInvalidType)
		}
	case float64:
		gl.Uniform1d(u.id, value)
	case []float64:
		switch len(value) {
		case 1:
			gl.Uniform1dv(u.id, 1, &value[0])
		case 2:
			gl.Uniform2dv(u.id, 2, &value[0])
		case 3:
			gl.Uniform3dv(u.id, 3, &value[0])
		case 4:
			gl.Uniform4dv(u.id, 4, &value[0])
		default:
			panic(ErrUniformInvalidType)
		}
	case mgl32.Mat2:
		gl.UniformMatrix2fv(u.id, 1, false, &value[0])
	case mgl32.Mat3:
		gl.UniformMatrix3fv(u.id, 1, false, &value[0])
	case mgl32.Mat4:
		gl.UniformMatrix4fv(u.id, 1, false, &value[0])
	case *Texture:
		value.Bind(0) //TODO support multiple textures per-shader
		err := value.SetUniform(u.id)
		if err != nil {
			panic(err)
		}
	default:
		panic(ErrUniformInvalidType)
	}
}
