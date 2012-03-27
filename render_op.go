package glutils

import (
	"github.com/pzsz/gl"
	v "github.com/pzsz/lin3dmath"
)

type IRenderOp interface {
	Render(cam *Camera, transform *v.Matrix4)
}

type SimpleRenderOp struct {
	Textures       []*Texture
	Blending       bool
	Buffer         *MeshBuffer
	SProgram       *ShaderProgram
	SProgramConf   func(*ShaderProgram)
}

func NewSimpleRenderOp(blending bool, buffer *MeshBuffer, tex ...*Texture) *SimpleRenderOp {
	return &SimpleRenderOp{tex, blending, buffer, nil, nil}
}

func NewShaderRenderOp(blending bool, sprogram *ShaderProgram, conf func(*ShaderProgram), buffer *MeshBuffer, tex ...*Texture) *SimpleRenderOp {
	return &SimpleRenderOp{tex, blending, buffer, sprogram, conf}
}

func (self *SimpleRenderOp) Render(cam *Camera, m *v.Matrix4) {
	cam.LoadProjection()
	cam.LoadModelview(m)

	if self.Blending {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
		gl.DepthMask(false);
	}

	if len(self.Textures) > 0 {
		gl.Enable(gl.TEXTURE_2D)
		for i:=0;i < len(self.Textures);i++ {
			self.Textures[i].Bind(i)
		}
	}

	if self.SProgram != nil {
		self.SProgram.Use()
		self.SProgramConf(self.SProgram)
	}

	if self.Buffer.HaveVBO() {
		DrawVBO(self.Buffer)
	} else {
		DrawArray(self.Buffer)
	}

	if self.SProgram != nil {
		self.SProgram.Unuse()
	}

	if self.Blending {
		gl.Disable(gl.BLEND)
		gl.DepthMask(true)
	}

	if len(self.Textures) > 0 {
		gl.Disable(gl.TEXTURE_2D)
		for i:=0;i < len(self.Textures);i++ {
			self.Textures[i].Unbind(i)
		}
	}
}

func DrawArray(buffer *MeshBuffer) {
	vertexSize := buffer.CalcVertexSize()

	gl.EnableClientState( gl.VERTEX_ARRAY)
	gl.VertexPointerTyped(3, gl.FLOAT, vertexSize, buffer.vertexArray)

	if ((buffer.Buffers & BUF_NORMAL) != 0) {
		off := buffer.CalcVertexOffset(BUF_NORMAL)
		gl.EnableClientState( gl.NORMAL_ARRAY)
		gl.NormalPointerTyped(gl.FLOAT, vertexSize, buffer.vertexArray[off:])
	}

	if ((buffer.Buffers & BUF_COLOUR) != 0) {
		off := buffer.CalcVertexOffset(BUF_COLOUR)
		gl.EnableClientState( gl.COLOR_ARRAY);
		gl.ColorPointerTyped(4, gl.UNSIGNED_BYTE, vertexSize, buffer.vertexArray[off:])
	}

	if ((buffer.Buffers & BUF_TEX_COORD0) != 0) {
		off := buffer.CalcVertexOffset(BUF_TEX_COORD0)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY);
		gl.TexCoordPointerTyped(2, gl.FLOAT, vertexSize, buffer.vertexArray[off:])
	}

	gl.DrawElementsTyped(gl.TRIANGLES, buffer.IndiceCount, gl.UNSIGNED_SHORT, buffer.indiceArray);

	if ((buffer.Buffers & BUF_NORMAL) != 0) {
		gl.DisableClientState(gl.NORMAL_ARRAY);
	}
	if ((buffer.Buffers & BUF_COLOUR) != 0) {
		gl.DisableClientState(gl.COLOR_ARRAY);
	}
	if ((buffer.Buffers & BUF_TEX_COORD0) != 0) {
		gl.DisableClientState(gl.TEXTURE_COORD_ARRAY);
	}

	gl.DisableClientState( gl.VERTEX_ARRAY)
}

func DrawVBO(buffer *MeshBuffer) {
	vertexSize := buffer.CalcVertexSize()

	buffer.VertexBuffer.Bind(gl.ARRAY_BUFFER)
	buffer.IndiceBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	gl.EnableClientState( gl.VERTEX_ARRAY)
	gl.VertexPointerVBO(3, gl.FLOAT, vertexSize, 0)

	if ((buffer.Buffers & BUF_NORMAL) != 0) {
		off := buffer.CalcVertexOffset(BUF_NORMAL)
		gl.EnableClientState( gl.NORMAL_ARRAY)
		gl.NormalPointerVBO(gl.FLOAT, vertexSize, off)
	}

	if ((buffer.Buffers & BUF_COLOUR) != 0) {
		off := buffer.CalcVertexOffset(BUF_COLOUR)
		gl.EnableClientState( gl.COLOR_ARRAY);
		gl.ColorPointerVBO(4, gl.UNSIGNED_BYTE, vertexSize, off)
	}

	if ((buffer.Buffers & BUF_TEX_COORD0) != 0) {
		off := buffer.CalcVertexOffset(BUF_TEX_COORD0)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY);
		gl.TexCoordPointerVBO(2, gl.FLOAT, vertexSize, off)
	}

	gl.DrawElementsVBO(gl.TRIANGLES, gl.UNSIGNED_SHORT, buffer.IndiceCount);

	if ((buffer.Buffers & BUF_NORMAL) != 0) {
		gl.DisableClientState(gl.NORMAL_ARRAY);
	}
	if ((buffer.Buffers & BUF_COLOUR) != 0) {
		gl.DisableClientState(gl.COLOR_ARRAY);
	}
	if ((buffer.Buffers & BUF_TEX_COORD0) != 0) {
		gl.DisableClientState(gl.TEXTURE_COORD_ARRAY);
	}

	gl.DisableClientState( gl.VERTEX_ARRAY)

	gl.BufferUnbind(gl.ARRAY_BUFFER)
	gl.BufferUnbind(gl.ELEMENT_ARRAY_BUFFER)
}