package glutils

import (
	"math"
	v "github.com/pzsz/lin3dmath"
	)

type FpsController struct {
	Camera  *Camera
	Pos     v.Vector3f

	HorAxis float32
        VerAxis float32
}

func NewFpsController(camera *Camera) *FpsController {
	ret := &FpsController{
	Camera: camera}
		
	return ret
}

func (s *FpsController) MoveBy(forward, strafe float32) {
	s.Pos.AddIP(s.GetForwardVector().Mul(forward))
	s.Pos.AddIP(s.GetStrafeVector().Mul(strafe))
}

func (s *FpsController) GetForwardVector() v.Vector3f {
	return v.Vector3f{
		float32(math.Sin(float64(-s.HorAxis))),
		0,
		float32(-math.Cos(float64(-s.HorAxis)))}
}

func (s *FpsController) GetViewVector() v.Vector3f {
	
	hor_x := math.Sin(float64(-s.HorAxis))
	hor_z := -math.Cos(float64(-s.HorAxis))

	ver_len := math.Cos(float64(-s.VerAxis))
	ver_hor := math.Sin(float64(-s.VerAxis))

	return v.Vector3f{
		float32(hor_x * ver_len),
		float32(-ver_hor),
		float32(hor_z * ver_len)}
}


func (s *FpsController) GetStrafeVector() v.Vector3f {
	return v.Vector3f{
		float32(math.Cos(float64(-s.HorAxis))),
		0,
		float32(math.Sin(float64(-s.HorAxis)))}
}


func (s *FpsController) RotateBy(deltaHor, deltaVer float32) {
	s.HorAxis += deltaHor
	if s.HorAxis > 2*math.Pi {
		s.HorAxis -= 2*math.Pi
	}
	if s.HorAxis < 0 {
		s.HorAxis += 2*math.Pi
	}

	s.VerAxis += deltaVer
	if s.VerAxis > 0.49*math.Pi {
		s.VerAxis = 0.49*math.Pi
	}
	if s.VerAxis < -0.49*math.Pi {
		s.VerAxis = -0.49*math.Pi
	}
}

func (s *FpsController) SetupCamera() {
	hor := v.MatrixRotate(v.Angle(-s.HorAxis), 0, 1, 0)
	ver := v.MatrixRotate(v.Angle(-s.VerAxis), 1, 0, 0)
	tr := v.MatrixTranslate(-s.Pos.X, -s.Pos.Y, -s.Pos.Z)

	rot := ver.Mul(hor)
	mat := rot.Mul(tr)
	s.Camera.SetCustomModelview(s.Pos.X, s.Pos.Y, s.Pos.Z, &mat)
}
