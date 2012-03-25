package glutils 

import (
	"github.com/pzsz/gl"
	v "github.com/pzsz/lin3dmath"
)

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
	gl.Vertex3f(-size,-size,0)
	gl.Vertex3f(-size,size,0)

	gl.Vertex3f(-size,size,0)
	gl.Vertex3f(size,size,0)

	gl.Vertex3f(size,size,0)
	gl.Vertex3f(size,-size,0)

	gl.Vertex3f(size,-size,0)
	gl.Vertex3f(-size,-size,0)
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
	gl.Vertex3f(-size_x,-size_y,0)
	gl.Vertex3f(-size_x,size_y,0)

	gl.Vertex3f(-size_x,size_y,0)
	gl.Vertex3f(size_x,size_y,0)

	gl.Vertex3f(size_x,size_y,0)
	gl.Vertex3f(size_x,-size_y,0)

	gl.Vertex3f(size_x,-size_y,0)
	gl.Vertex3f(-size_x,-size_y,0)

	gl.End()

	gl.Color4ub(255, 255, 255, 255)
}

func RenderSprite(cam *Camera, m *v.Matrix4, sizeX, sizeY float32, t *Texture) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
	gl.DepthMask(false);
	//gl.Disable(gl.DEPTH_TEST)	

	RenderTexturedRect(cam, m, sizeX,sizeY, t)

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
	gl.TexCoord2f(0,1)
	gl.Vertex3f(-sizeX,-sizeY,0)

	gl.TexCoord2f(1,1)
	gl.Vertex3f(sizeX,-sizeY,0)

	gl.TexCoord2f(1,0)
	gl.Vertex3f(sizeX,sizeY,0)

	gl.TexCoord2f(0,0)
	gl.Vertex3f(-sizeX,sizeY,0)
	gl.End()

	gl.Disable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0);
}

func RenderRect(cam *Camera, m *v.Matrix4, sizeX, sizeY float32, c Colour) {
	cam.LoadProjection()
	cam.LoadModelview(m)

	gl.Color4ub(c.R, c.G, c.B, c.A)

	gl.Begin(gl.QUADS)
	gl.Vertex3f(-sizeX,-sizeY,0)
	gl.Vertex3f(sizeX,-sizeY,0)
	gl.Vertex3f(sizeX,sizeY,0)
	gl.Vertex3f(-sizeX,sizeY,0)
	gl.End()
}