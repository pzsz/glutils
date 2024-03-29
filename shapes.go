package glutils

import (
	"github.com/pzsz/gl"
	v "github.com/pzsz/lin3dmath"
)

func BuildCubeBuffer(halfSize v.Vector3f) *MeshBuffer {
	buffer := NewMeshBuffer(8, 6*6, RENDER_POLYGONS, 0)
	build := ReuseMeshBuilder(buffer)

	//   7----6
	//  /|   /|
	// 4----5 |
	// | |  | |
	// | 3--|-2
	// |/   |/
	// 0----1
	
	i0 := build.StartVertex()
	build.AddPosition(-halfSize.X, -halfSize.Y, -halfSize.Z)

	i1 := build.StartVertex()
	build.AddPosition(halfSize.X, -halfSize.Y, -halfSize.Z)

	i2 := build.StartVertex()
	build.AddPosition(halfSize.X, halfSize.Y, -halfSize.Z)

	i3 := build.StartVertex()
	build.AddPosition(-halfSize.X, halfSize.Y, -halfSize.Z)

	i4 := build.StartVertex()
	build.AddPosition(-halfSize.X, -halfSize.Y, halfSize.Z)

	i5 := build.StartVertex()
	build.AddPosition(halfSize.X, -halfSize.Y, halfSize.Z)

	i6 := build.StartVertex()
	build.AddPosition(halfSize.X, halfSize.Y, halfSize.Z)

	i7 := build.StartVertex()
	build.AddPosition(-halfSize.X, halfSize.Y, halfSize.Z)

	build.AddIndice3(i0, i1, i5)
	build.AddIndice3(i5, i4, i0)

	build.AddIndice3(i1, i2, i6)
	build.AddIndice3(i6, i5, i1)

	build.AddIndice3(i6, i2, i3)
	build.AddIndice3(i3, i7, i6)

	build.AddIndice3(i3, i0, i4)
	build.AddIndice3(i4, i7, i3)

	build.AddIndice3(i4, i5, i6)
	build.AddIndice3(i6, i7, i4)

	build.AddIndice3(i0, i3, i2)
	build.AddIndice3(i2, i1, i0)

	build.Finalize(true, buffer)

	return buffer
}


func RenderLine(camera *Camera, m *v.Matrix4, from, to v.Vector3f, colour Colour) {
	camera.LoadProjection()
	camera.LoadModelview(m)

	gl.Color4ub(colour.R, colour.G, colour.B, colour.A)

	gl.Begin(gl.LINES)
	gl.Vertex3f(from.X, from.Y, from.Z)
	gl.Vertex3f(to.X, to.Y, to.Z)
	gl.End()
}

func RenderWireQuad(camera *Camera, m *v.Matrix4, size float32, colour Colour) {
	camera.LoadProjection()
	camera.LoadModelview(m)

	gl.Color4ub(colour.R, colour.G, colour.B, colour.A)

	gl.Begin(gl.LINES)
	gl.Vertex3f(-size, -size, 0)
	gl.Vertex3f(-size, size, 0)

	gl.Vertex3f(-size, size, 0)
	gl.Vertex3f(size, size, 0)

	gl.Vertex3f(size, size, 0)
	gl.Vertex3f(size, -size, 0)

	gl.Vertex3f(size, -size, 0)
	gl.Vertex3f(-size, -size, 0)
	gl.End()
}

func RenderUIStart() {
	gl.Disable(gl.DEPTH_TEST)
	gl.DepthMask(false)
}

func RenderUIEnd() {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
}

func RenderWireRect(cam *Camera, m *v.Matrix4, size_x, size_y float32, colour Colour) {
	cam.LoadProjection()
	cam.LoadModelview(m)

	gl.Color4ub(colour.R, colour.G, colour.B, colour.A)

	gl.Begin(gl.LINES)
	gl.Vertex3f(-size_x, -size_y, 0)
	gl.Vertex3f(-size_x, size_y, 0)

	gl.Vertex3f(-size_x, size_y, 0)
	gl.Vertex3f(size_x, size_y, 0)

	gl.Vertex3f(size_x, size_y, 0)
	gl.Vertex3f(size_x, -size_y, 0)

	gl.Vertex3f(size_x, -size_y, 0)
	gl.Vertex3f(-size_x, -size_y, 0)

	gl.End()

	gl.Color4ub(255, 255, 255, 255)
}

func RenderSprite(cam *Camera, m *v.Matrix4, sizeX, sizeY float32, t *Texture) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(false)
	//gl.Disable(gl.DEPTH_TEST)	

	RenderTexturedRect(cam, m, sizeX, sizeY, t)

	gl.Disable(gl.BLEND)
	gl.DepthMask(true)
	//gl.Enable(gl.DEPTH_TEST)
}

func RenderTexturedRect(cam *Camera, m *v.Matrix4, sizeX, sizeY float32, t *Texture) {
	cam.LoadProjection()
	cam.LoadModelview(m)

	gl.Enable(gl.TEXTURE_2D)
	t.tex.Bind(gl.TEXTURE_2D)

	gl.Color4ub(255, 255, 255, 255)
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-sizeX, -sizeY, 0)

	gl.TexCoord2f(1, 1)
	gl.Vertex3f(sizeX, -sizeY, 0)

	gl.TexCoord2f(1, 0)
	gl.Vertex3f(sizeX, sizeY, 0)

	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-sizeX, sizeY, 0)
	gl.End()

	gl.Disable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func RenderRect(cam *Camera, m *v.Matrix4, sizeX, sizeY float32, c Colour) {
	cam.LoadProjection()
	cam.LoadModelview(m)

	gl.Color4ub(c.R, c.G, c.B, c.A)

	gl.Begin(gl.QUADS)
	gl.Vertex3f(-sizeX, -sizeY, 0)
	gl.Vertex3f(sizeX, -sizeY, 0)
	gl.Vertex3f(sizeX, sizeY, 0)
	gl.Vertex3f(-sizeX, sizeY, 0)
	gl.End()
}
