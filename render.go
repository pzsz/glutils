package glutils

import (
	"github.com/pzsz/gl"
)

func Setup() {
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.LIGHTING)
	gl.Enable(gl.DEPTH_TEST)
}

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
