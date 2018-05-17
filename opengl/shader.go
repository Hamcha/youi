package opengl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Shader errors
var (
	ErrMissingVertexShader   = errors.New("vertex shader is missing")
	ErrMissingFragmentShader = errors.New("fragment shader is missing")
)

// Shader is a OpenGL program with one vertex and one fragment shader
type Shader struct {
	fragID, vertID uint32
	programID      uint32

	output *uint8
}

// SetFragmentSource sets (and updates) the fragment portion of the shader
func (s *Shader) SetFragmentSource(src string) (err error) {
	// Do some cleanup if we're replacing an existing shader
	if s.fragID != 0 {
		// Already in a program? Detach
		if s.programID != 0 {
			gl.DetachShader(s.programID, s.fragID)
		}
		// Delete old shader
		gl.DeleteShader(s.fragID)
	}

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
	if s.vertID != 0 {
		// Already in a program? Detach
		if s.programID != 0 {
			gl.DetachShader(s.programID, s.vertID)
		}
		// Delete old shader
		gl.DeleteShader(s.vertID)
	}

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

func (s *Shader) Use() {
	gl.UseProgram(s.programID)
}

func (s *Shader) BindUniforms() {
	//TODO

	// Bind output
	gl.BindFragDataLocation(s.programID, 0, s.output)
}

func (s *Shader) SetOutput(out string) {
	s.output = gl.Str(out + "\x00")
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
	shader := new(Shader)

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

	shader.SetOutput("color")

	return shader
}

// DefaultVertexShader is the vertex portion of the default shader
const DefaultVertexShader = `
#version 330 core
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
out vec2 position;
void main() {
	fragTexCoord = vertTexCoord;
	position = vert.xy;
	gl_Position = vec4(vert, 1);
}
` + "\x00"

// DefaultFragmentShader is the fragment portion of the default shader
const DefaultFragmentShader = `
#version 330 core
in vec2 fragTexCoord;
in vec2 position;
out vec4 color;
void main() {
	vec2 pos = position.xy / fragTexCoord.xy;
    color = vec4(pos.x,pos.y,1.0,1.0);
}
` + "\x00"
