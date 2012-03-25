package glutils

import (
	"bytes"
	"encoding/binary"
)

var byteOrder binary.ByteOrder = binary.LittleEndian

type MeshBuilder struct {
	VertexCounter int
	IndiceCounter int
	VertexBuf     *bytes.Buffer
	Buffer        *MeshBuffer
	UseVBO        bool
}

func NewMeshBuilder(vertexCount, indiceCount, renderOp, buffers int, useVBO bool,
	attr ... MeshBufferAttribute) *MeshBuilder {

	buf := NewMeshBuffer(indiceCount, vertexCount, renderOp, buffers, attr...)
	buf.AllocArrays()
	vertbuf := bytes.NewBuffer(buf.VertexArray)
	ret := &MeshBuilder{0, 0,  vertbuf, buf, useVBO}
	return ret
}

func ReuseMeshBuilder(buf *MeshBuffer) *MeshBuilder {
	buf.ResetArrays()
	vertbuf := bytes.NewBuffer(buf.VertexArray)
	buf.IndiceArray = buf.IndiceArray[0:0]
	return &MeshBuilder{0, 0, vertbuf, buf, buf.HaveVBO()}
}

func (self *MeshBuilder) StartVertex() (r int) {
	r = self.VertexCounter
	self.VertexCounter++
	return 
}

func (self *MeshBuilder) AddPosition(x,y,z float32) {
	binary.Write(self.VertexBuf, byteOrder, x)
	binary.Write(self.VertexBuf, byteOrder, y)
	binary.Write(self.VertexBuf, byteOrder, z)
}

func (self *MeshBuilder) AddNormal(x,y,z float32) {
	binary.Write(self.VertexBuf, byteOrder, x)
	binary.Write(self.VertexBuf, byteOrder, y)
	binary.Write(self.VertexBuf, byteOrder, z)
}


func (self *MeshBuilder) AddColour(r,g,b,a byte) {
	binary.Write(self.VertexBuf, byteOrder, r)
	binary.Write(self.VertexBuf, byteOrder, g)
	binary.Write(self.VertexBuf, byteOrder, b)
	binary.Write(self.VertexBuf, byteOrder, a)
}

func (self *MeshBuilder) AddTexCoord(u, v float32) {
	binary.Write(self.VertexBuf, byteOrder, u)
	binary.Write(self.VertexBuf, byteOrder, v)
}

func (self *MeshBuilder) AddAttr2F(u, v float32) {
	binary.Write(self.VertexBuf, byteOrder, u)
	binary.Write(self.VertexBuf, byteOrder, v)
}

func (self *MeshBuilder) AddIndice3(a,b,c int) {
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(a))
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(b))
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(c))
	self.IndiceCounter += 3
}

func (self *MeshBuilder) AddIndice4(a,b,c,d int) {
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(a))
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(b))
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(c))
	self.Buffer.IndiceArray = append(self.Buffer.IndiceArray,uint16(d))
	self.IndiceCounter += 4
}

func (self *MeshBuilder) Finalize() *MeshBuffer {
	self.Buffer.VertexCount = self.VertexCounter
	self.Buffer.IndiceCount = self.IndiceCounter

	self.Buffer.VertexArray = self.VertexBuf.Bytes()
	self.Buffer.IndiceArray = self.Buffer.IndiceArray[0:self.IndiceCounter]

	if self.UseVBO {
		self.Buffer.BufferData()
	}
	return self.Buffer
}
