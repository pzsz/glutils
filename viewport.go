package glutils

import (
	"github.com/pzsz/gl"
)

type Viewport struct {
	Width    float32
	Height   float32
	Aspect   float32
}

var viewportInstance *Viewport = &Viewport{}

func GetViewport() *Viewport {
	return viewportInstance
}

func (self *Viewport) SetScreenSize(w,h float32) {
	self.Width = w
	self.Height = h
	self.Aspect = w/h
	gl.Viewport(0, 0, int(w), int(h))
}
