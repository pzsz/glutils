package glutils

import "github.com/pzsz/gl"

const (
	BUF_NORMAL     = 1
	BUF_COLOUR     = 2
	BUF_TEX_COORD0 = 4
	BUF_ATTR0      = 8

	RENDER_POLYGONS = 1
	RENDER_STRIPE   = 2

	MAX_ATTR = 2
)

type MeshBufferAttribute struct {
	Name string
	Size int
}

type MeshBuffer struct {
	vertexArray []uint8
	indiceArray []uint8

	VertexBuffer gl.Buffer
	IndiceBuffer gl.Buffer

	Buffers     int
	RenderOp    int
	VertexCount int
	IndiceCount int

	Attributes []MeshBufferAttribute
}

func NewMeshBuffer(indiceCount, vertexCount, renderOp, buffers int, attr ...MeshBufferAttribute) *MeshBuffer {

	ret := &MeshBuffer{
		Buffers:     buffers,
		RenderOp:    renderOp,
		VertexCount: vertexCount,
		IndiceCount: indiceCount,
		Attributes:  attr}
	return ret
}

func (self *MeshBuffer) AllocArrays() {
	vs := self.CalcVertexSize()
	self.vertexArray = make([]uint8, 0, vs*self.VertexCount)
	self.indiceArray = make([]uint8, 0, 2*self.IndiceCount)
}

func (self *MeshBuffer) GetArrays() (vertexArray, indiceArray []uint8) {
	vertexArray, indiceArray = self.vertexArray, self.indiceArray
	return
}

func (self *MeshBuffer) AllocBuffers() {
	if self.VertexBuffer == 0 {
		self.VertexBuffer = gl.GenBuffer()
	}
	if self.IndiceBuffer == 0 {
		self.IndiceBuffer = gl.GenBuffer()
	}
}

func (self *MeshBuffer) HaveVBO() bool {
	return self.VertexBuffer != 0
}

func (self *MeshBuffer) Destroy() {
	if self.VertexBuffer != 0 {
		self.VertexBuffer.Delete()
	}
	if self.IndiceBuffer != 0 {
		self.IndiceBuffer.Delete()
	}

	self.vertexArray = nil
	self.indiceArray = nil
}

func (self *MeshBuffer) CopyArraysToVBO() {
	self.AllocBuffers()

	vs := self.CalcVertexSize()
	self.VertexBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, vs*self.VertexCount,
		self.vertexArray, gl.STATIC_DRAW)
	gl.BufferUnbind(gl.ARRAY_BUFFER)

	self.IndiceBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*self.IndiceCount,
		self.indiceArray, gl.STATIC_DRAW)
	gl.BufferUnbind(gl.ELEMENT_ARRAY_BUFFER)
}

func (self *MeshBuffer) CalcVertexSize() int {
	return self.CalcVertexOffset(0)
}

func (self *MeshBuffer) SetAttribute(i int, attr MeshBufferAttribute) {
	self.Attributes[i] = attr
}

func (self *MeshBuffer) CalcVertexOffset(element int) int {
	ret := 3 * 4 // POS

	if element == BUF_NORMAL {
		return ret
	}
	if (self.Buffers & BUF_NORMAL) != 0 {
		ret += 3 * 4
	}

	if element == BUF_COLOUR {
		return ret
	}
	if (self.Buffers & BUF_COLOUR) != 0 {
		ret += 4
	}

	if element == BUF_TEX_COORD0 {
		return ret
	}
	if (self.Buffers & BUF_TEX_COORD0) != 0 {
		ret += 2 * 4
	}

	if element == BUF_ATTR0 {
		return ret
	}
	if (self.Buffers & BUF_ATTR0) != 0 {
		ret += 4 * self.Attributes[0].Size
	}

	return ret
}
