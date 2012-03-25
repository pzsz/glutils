package glutils

type ShaderManager struct {
	Shaders  map[string]*Shader
	Programs map[string]*ShaderProgram
}

var shaderManager *ShaderManager = newShaderManager()

func GetProgram(vertexFilename, fragFilename string) (*ShaderProgram, error) {
	return shaderManager.GetProgram(vertexFilename, fragFilename)
}

func newShaderManager() *ShaderManager {
	return &ShaderManager{
		make(map[string]*Shader),
		make(map[string]*ShaderProgram)}
}

func (self *ShaderManager) getShader(filename string) (*Shader, error) {
	kept := self.Shaders[filename]
	if kept != nil {
		return kept, nil
	}

	newShader, err := newShader(filename)
	if err != nil {
		return nil, err
	}

	self.Shaders[filename] = newShader
	return newShader, nil
}

func (self *ShaderManager) GetProgram(vertexFilename, fragFilename string) (*ShaderProgram, error) {
	compositeName := vertexFilename + "|" + fragFilename
	kept := self.Programs[compositeName]
	if kept != nil {
		return kept, nil
	}

	vertexShader, err := self.getShader(vertexFilename)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := self.getShader(fragFilename)
	if err != nil {
		return nil, err
	}

	program, err := newShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	self.Programs[compositeName] = program
	return program, nil
}
