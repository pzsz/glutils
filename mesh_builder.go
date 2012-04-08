package glutils

import (
	"bytes"
	"encoding/binary"
)

var byteOrder binary.ByteOrder = binary.LittleEndian

type MeshBuilder struct {
	Buffer *MeshBuffer

	VertexCounter int
	IndiceCounter int

	VertexWriter *bytes.Buffer
	IndiceWriter *bytes.Buffer
}

func NewMeshBuilder(vertexCount, indiceCount, renderOp, buffers int,
	attr ...MeshBufferAttribute) *MeshBuilder {

	buf := NewMeshBuffer(indiceCount, vertexCount, renderOp, buffers, attr...)
	buf.AllocArrays()

	vertex_writer := bytes.NewBuffer(buf.vertexArray)
	indice_writer := bytes.NewBuffer(buf.indiceArray)

	ret := &MeshBuilder{buf, 0, 0, vertex_writer, indice_writer}
	return ret
}

func ReuseMeshBuilder(buf *MeshBuffer) *MeshBuilder {
	// Those 3 lines with make arrays empty, but with previous capacity
	vertex_writer := bytes.NewBuffer(buf.vertexArray[0:0])
	indice_writer := bytes.NewBuffer(buf.indiceArray[0:0])

	buf.indiceArray = buf.indiceArray[0:0]
	return &MeshBuilder{buf, 0, 0, vertex_writer, indice_writer}
}

func (self *MeshBuilder) StartVertex() (r int) {
	r = self.VertexCounter
	self.VertexCounter++
	return
}

func (self *MeshBuilder) AddPosition(x, y, z float32) {
	binary.Write(self.VertexWriter, byteOrder, x)
	binary.Write(self.VertexWriter, byteOrder, y)
	binary.Write(self.VertexWriter, byteOrder, z)
}

func (self *MeshBuilder) AddNormal(x, y, z float32) {
	binary.Write(self.VertexWriter, byteOrder, x)
	binary.Write(self.VertexWriter, byteOrder, y)
	binary.Write(self.VertexWriter, byteOrder, z)
}

func (self *MeshBuilder) AddColour(r, g, b, a byte) {
	binary.Write(self.VertexWriter, byteOrder, r)
	binary.Write(self.VertexWriter, byteOrder, g)
	binary.Write(self.VertexWriter, byteOrder, b)
	binary.Write(self.VertexWriter, byteOrder, a)
}

func (self *MeshBuilder) AddTexCoord(u, v float32) {
	binary.Write(self.VertexWriter, byteOrder, u)
	binary.Write(self.VertexWriter, byteOrder, v)
}

func (self *MeshBuilder) AddAttr2F(u, v float32) {
	binary.Write(self.VertexWriter, byteOrder, u)
	binary.Write(self.VertexWriter, byteOrder, v)
}

func (self *MeshBuilder) AddIndice3(a, b, c int) {
	binary.Write(self.IndiceWriter, byteOrder, int16(a))
	binary.Write(self.IndiceWriter, byteOrder, int16(b))
	binary.Write(self.IndiceWriter, byteOrder, int16(c))
	self.IndiceCounter += 3
}

func (self *MeshBuilder) AddIndice4(a, b, c, d int) {
	binary.Write(self.IndiceWriter, byteOrder, uint16(a))
	binary.Write(self.IndiceWriter, byteOrder, uint16(b))
	binary.Write(self.IndiceWriter, byteOrder, uint16(c))
	binary.Write(self.IndiceWriter, byteOrder, uint16(d))
	self.IndiceCounter += 4
}

func (self *MeshBuilder) IsEmpty() bool {
	return self.IndiceCounter == 0
}

func (self *MeshBuilder) Finalize(useVBO bool) *MeshBuffer {
	self.Buffer.VertexCount = self.VertexCounter
	self.Buffer.IndiceCount = self.IndiceCounter

	self.Buffer.vertexArray = self.VertexWriter.Bytes()
	self.Buffer.indiceArray = self.IndiceWriter.Bytes()

	if useVBO {
		self.Buffer.CopyArraysToVBO()
	}
	return self.Buffer
}
