package glutils

import (
	"errors"
	"github.com/pzsz/gl"
	"io/ioutil"
	"os"
	"strings"
)

type Shader struct {
	ShaderObject gl.Shader
	Filename     string
}

func newShader(filename string) (*Shader, error) {
	var shaderType gl.GLenum

	if strings.Index(filename, ".fragment") != -1 {
		shaderType = gl.FRAGMENT_SHADER
	} else {
		shaderType = gl.VERTEX_SHADER
	}

	shader := &Shader{gl.CreateShader(shaderType), filename}

	source, err := readShaderSource(filename)
	if err != nil {
		return nil, err
	}

	shader.ShaderObject.Source(source)
	shader.ShaderObject.Compile()

	if shader.ShaderObject.Get(gl.COMPILE_STATUS) == gl.FALSE {
		info := shader.ShaderObject.GetInfoLog()
		if info != "" {
			return nil, errors.New("Error while compiling GLSL " + filename + ":\n" + info)
		}
	}
	return shader, nil
}

func readShaderSource(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ret, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

type ShaderProgram struct {
	Vertex        *Shader
	Fragment      *Shader
	ProgramObject gl.Program
}

func newShaderProgram(vertex, fragment *Shader) (*ShaderProgram, error) {
	ret := &ShaderProgram{vertex, fragment, gl.CreateProgram()}

	ret.ProgramObject.AttachShader(vertex.ShaderObject)
	ret.ProgramObject.AttachShader(fragment.ShaderObject)
	ret.ProgramObject.Link()
	if ret.ProgramObject.Get(gl.LINK_STATUS) == gl.FALSE {
		info := ret.ProgramObject.GetInfoLog()
		if info != "" {
			return nil, errors.New("Error while linking GLSL " + vertex.Filename + " with " + fragment.Filename + ": " + info)
		}
	}
	ret.ProgramObject.Validate()
	if ret.ProgramObject.Get(gl.VALIDATE_STATUS) == gl.FALSE {
		info2 := ret.ProgramObject.GetInfoLog()
		if info2 != "" {
			return nil, errors.New("Error while linking GLSL " + vertex.Filename + " with " + fragment.Filename + ": " + info2)
		}
	}

	return ret, nil
}

func (self *ShaderProgram) Use() {
	self.ProgramObject.Use()
}

func (self *ShaderProgram) Unuse() {
	gl.ProgramUnuse()
}

func (self *ShaderProgram) GetUniform(name string) gl.UniformLocation {
	return self.ProgramObject.GetUniformLocation(name)
}
