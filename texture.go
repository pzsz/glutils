package glutils

import (
	"github.com/pzsz/gl"
	"image"
	"image/png"
	"os"
	//	"runtime"
	//	"fmt"
)

type TexFilterType int

const (
	NEAREST   = TexFilterType(1)
	LINEAR    = TexFilterType(2)
	TRILINEAR = TexFilterType(3)
)

type TexSetup struct {
	InternalFormat int
	Format         gl.GLenum
	Mipmaps        bool
	Filtering      TexFilterType
}

var DEFAULT_TEXSETUP = TexSetup{gl.RGBA, gl.RGBA, true, TRILINEAR}
var NO_MIPMAP_TEXSETUP = TexSetup{gl.RGBA, gl.RGBA, false, LINEAR}
var ALPHA_TEXSETUP = TexSetup{gl.ALPHA, gl.ALPHA, false, LINEAR}

type Texture struct {
	tex    gl.Texture
	Name   string
	Width  int
	Height int

	Setup TexSetup
}

func (self *Texture) Bind(i int) {
	gl.ActiveTexture(gl.GLenum(gl.TEXTURE0 + i))
	self.tex.Bind(gl.TEXTURE_2D)
}

func (self *Texture) Unbind(i int) {
	gl.ActiveTexture(gl.GLenum(gl.TEXTURE0 + i))
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (self *Texture) Destroy() {
	self.tex.Delete()
}

func (self *Texture) LoadData(data []uint8) {
	self.tex.Bind(gl.TEXTURE_2D)

	if data == nil {
		gl.TexImage2D(gl.TEXTURE_2D, 0, self.Setup.InternalFormat,
			self.Width, self.Height, 0, self.Setup.Format,
			gl.UNSIGNED_BYTE,
			nil)
	} else {
		gl.TexImage2D(gl.TEXTURE_2D, 0, self.Setup.InternalFormat,
			self.Width, self.Height, 0, self.Setup.Format,
			gl.UNSIGNED_BYTE,
			data)
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (self *Texture) setupParams() {
	self.tex.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	if self.Setup.Mipmaps {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	switch self.Setup.Filtering {
	case NEAREST:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		break
	case LINEAR:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		break
	case TRILINEAR:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		break
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

type TextureManager struct {
	Textures map[string]*Texture
}

var tmInstance *TextureManager = NewTextureManager()

func NewTextureManager() *TextureManager {
	return &TextureManager{Textures: map[string]*Texture{}}
}

func GetTexture(name string) (*Texture, error) {
	return tmInstance.GetTexture(name)
}

func GetTextureManager() *TextureManager {
	return tmInstance
}

func (self *TextureManager) GetTexture(name string) (*Texture, error) {
	return self.GetTextureEx(name, DEFAULT_TEXSETUP)
}

func (self *TextureManager) GetTextureEx(name string, newSetup TexSetup) (*Texture, error) {
	cachet := self.Textures[name]
	if cachet != nil {
		return cachet, nil
	}

	tex, er := self.loadTexture(name, newSetup)
	if er != nil {
		return nil, er
	}

	self.Textures[name] = tex

	return tex, nil

}

func (self *TextureManager) CreateEmptyTexture(name string, width, height int, setup TexSetup) *Texture {
	t := gl.GenTexture()

	texture := &Texture{t, name, width, height, setup}
	texture.LoadData(nil)
	texture.setupParams()

	return texture
}

func getByteArray(img image.Image) (ret []byte) {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	size := img.Bounds().Max.X * img.Bounds().Max.Y * 4
	ret = make([]uint8, 0, size)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			ret = append(ret, uint8(r/256))
			ret = append(ret, uint8(g/256))
			ret = append(ret, uint8(b/256))
			ret = append(ret, uint8(a/256))
		}
	}
	return
}

func (self *TextureManager) loadTexture(filename string, setup TexSetup) (*Texture, error) {
	//	runtime.LockOSThread()

	f, er := os.Open(filename)
	if er != nil {
		return nil, er
	}
	defer f.Close()
	img, er := png.Decode(f)
	if er != nil {
		return nil, er
	}

	bytes := getByteArray(img)

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y

	t := gl.GenTexture()

	texture := &Texture{t, filename, width, height, setup}
	texture.LoadData(bytes)
	texture.setupParams()

	return texture, nil
}
