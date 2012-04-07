package glutils

import (
	"github.com/pzsz/gl"
	"github.com/pzsz/glu"
	v "github.com/pzsz/lin3dmath"
	"math"

//	"reflect"
//	"unsafe"
)

func CreateFrustrumMatrix(xmin, xmax, ymin, ymax, zNear, zFar float32) *v.Matrix4 {
	A := -(xmax + xmin) / (xmax - xmin)
	B := -(ymax + ymin) / (ymax - ymin)
	C := -(zFar + zNear) / (zFar - zNear)
	D := -2 * zFar * zNear / (zFar - zNear)

	return &v.Matrix4{
		(2 * zNear) / (xmax - xmin), 0, 0, 0,
		0, (2 * zNear) / (ymax - ymin), 0, 0,
		A, B, C, -1,
		0, 0, D, 0}
}

func CreateOrthoMatrix(xmin, xmax, ymin, ymax, zNear, zFar float32) *v.Matrix4 {
	tx := -(xmax + xmin) / (xmax - xmin)
	ty := -(ymax + ymin) / (ymax - ymin)
	tz := -(zFar + zNear) / (zFar - zNear)

	return &v.Matrix4{
		2 / (xmax - xmin), 0, 0, 0,
		0, 2 / (ymax - ymin), 0, 0,
		0, 0, -2 / (zFar - zNear), 0,
		tx, ty, tz, 1}
}

func CreateLookAtMatrix(eyex, eyey, eyez, centerx, centery, centerz, upx, upy, upz float32) *v.Matrix4 {

	var x, y, z [3]float32
	var mag float32

	/* Make rotation matrix */

	/* Z vector */
	z[0] = centerx - eyex
	z[1] = centery - eyey
	z[2] = centerz - eyez
	mag = float32(math.Sqrt(float64(z[0]*z[0] + z[1]*z[1] + z[2]*z[2])))
	if mag != 0 { /* mpichler, 19950515 */
		z[0] /= mag
		z[1] /= mag
		z[2] /= mag
	}

	/* Y vector */
	y[0] = upx
	y[1] = upy
	y[2] = upz

	/* X vector = Z cross Y */
	x[0] = z[1]*y[2] - z[2]*y[1]
	x[1] = -z[0]*y[2] + z[2]*y[0]
	x[2] = z[0]*y[1] - z[1]*y[0]

	mag = float32(math.Sqrt(float64(x[0]*x[0] + x[1]*x[1] + x[2]*x[2])))
	if mag != 0 {
		x[0] /= mag
		x[1] /= mag
		x[2] /= mag
	}

	/* Recompute Y = X cross Z */
	y[0] = x[1]*z[2] - x[2]*z[1]
	y[1] = -x[0]*z[2] + x[2]*z[0]
	y[2] = x[0]*z[1] - x[1]*z[0]

	// TODO inline this shit
	a := (v.Matrix4{x[0], x[1], -x[2], 0,
		y[0], y[1], -y[2], 0,
		z[0], z[1], -z[2], 0,
		0, 0, 0, 1.0})
	r := a.Mul(v.MatrixTranslate(-eyex, -eyey, -eyez))
	return &r
}

type Camera struct {
	Fov      float32
	NearZ    float32
	FarZ     float32
	Viewport *Viewport

	ModelviewMatrix  v.Matrix4
	ProjectionMatrix v.Matrix4

	EyePos  v.Vector3f
	ViewPos v.Vector3f
}

func NewCamera(viewport *Viewport) *Camera {
	return &Camera{Viewport: viewport}
}

func (self *Camera) SetModelviewOne() {
	self.ModelviewMatrix = *v.MatrixOne()
}

func (self *Camera) SetCustomModelview(posx, posy, posz float32, ModelviewMatrix *v.Matrix4) {
	self.EyePos = v.Vector3f{posx, posy, posz}

	self.ModelviewMatrix = *ModelviewMatrix
}


func (self *Camera) SetModelview(posx, posy, posz, lookx, looky, lookz, upx, upy, upz float32) {
	self.EyePos = v.Vector3f{posx, posy, posz}
	self.ViewPos = v.Vector3f{lookx, looky, lookz}

	self.ModelviewMatrix = *CreateLookAtMatrix(posx, posy, posz,
		lookx, looky, lookz,
		upx, upy, upz)
}

func (self *Camera) SetFrustrumProjection(fov, nearz, farz float32) {
	self.Fov = fov
	self.NearZ = nearz
	self.FarZ = farz

	ymax := self.NearZ * float32(math.Tan(float64(self.Fov*math.Pi/360)))
	ymin := -ymax
	xmin := ymin * self.Viewport.Aspect
	xmax := ymax * self.Viewport.Aspect

	self.ProjectionMatrix = *CreateFrustrumMatrix(xmin, xmax, ymin, ymax, self.NearZ, self.FarZ)
}

func (self *Camera) SetOrthoProjection(nearz, farz float32) {
	self.Fov = 0
	self.NearZ = nearz
	self.FarZ = farz

	self.ProjectionMatrix = *CreateOrthoMatrix(0, self.Viewport.Width,
		0, self.Viewport.Height,
		nearz, farz)
}

func (self *Camera) GetViewRay(x, y float32) v.Vector3f {
	viewport := [4]int32{0, 0, int32(self.Viewport.Width), int32(self.Viewport.Height)}

	mx, my, mz := glu.UnProject(
		float64(x), float64(self.Viewport.Height)-float64(y), float64(self.NearZ),
		self.ModelviewMatrix.ToArray64(),
		self.ProjectionMatrix.ToArray64(), &viewport)

	return v.Vector3f{self.EyePos.X - float32(mx),
		self.EyePos.Y - float32(my),
		self.EyePos.Z - float32(mz)}
}

func (self *Camera) ScreenToPlaneXY(x, y, z float32) v.Vector2f {
	viewRay := self.GetViewRay(x, y)

	if viewRay.Z == 0 {
		return v.Vector2f{0, 0}
	}

	t := (z - self.EyePos.Z) / viewRay.Z
	retx := self.EyePos.X + t*viewRay.X
	rety := self.EyePos.Y + t*viewRay.Y

	return v.Vector2f{retx, rety}
}

func (self *Camera) LoadProjection() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadMatrixf(self.ProjectionMatrix.ToArray32())
}

func (self *Camera) LoadModelview(m *v.Matrix4) {
	gl.MatrixMode(gl.MODELVIEW)

	fu := self.ModelviewMatrix.Mul(m)
	gl.LoadMatrixf(fu.ToArray32())
}
