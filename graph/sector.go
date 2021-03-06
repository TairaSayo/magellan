package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
)

type Sector struct {
	camParams CamParams

	center   v2.V2
	radius   float64
	startAng float64
	endAng   float64

	color color.Color
	alpha float64

	sprite *Sprite
}

const (
	sectorLen = 1000
	sectorDeg = 1
)

var oneDegreeTex Tex

func init() {
	h := sectorLen
	w := 1 + int(sectorLen*math.Tan(sectorDeg*Deg2Rad))
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)

	p := make([]byte, w*h*4)
	dw := w * 4
	tan := math.Tan(sectorDeg * Deg2Rad)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if x*x+y*y > sectorLen*sectorLen {
				continue
			}
			if x == 0 || (y > 0 && float64(x)/float64(y) < tan) {
				ix := w - 1 - x
				iy := h - 1 - y
				for i := 0; i < 4; i++ {
					p[4*ix+iy*dw+i] = 255
				}
			}
		}
	}

	img.ReplacePixels(p)
	oneDegreeTex = TexFromImage(img, ebiten.FilterDefault, 0, 0, 0, "~oneDegree")
}

func NewSector(params CamParams) *Sector {
	sprite := NewSprite(oneDegreeTex, params)
	sprite.SetPivot(BotRight())
	return &Sector{
		camParams: params,
		sprite:    sprite,
		color:     color.White,
		alpha:     1,
	}
}

func (s *Sector) Draw(dest *ebiten.Image) {
	s.recalcSprite()

	for ang := s.startAng; (ang + sectorDeg) <= s.endAng; ang += sectorDeg {
		s.sprite.SetAng(ang)
		s.sprite.Draw(dest)
	}

	lastPart := s.endAng - sectorDeg
	if lastPart > s.startAng {
		s.sprite.SetAng(lastPart)
		s.sprite.Draw(dest)
	}
}

func (s *Sector) DrawF() (DrawF, string) {
	return s.Draw, "~oneDegree"
}

func (s *Sector) recalcSprite() {
	s.sprite.SetPos(s.center)
	effRadius := s.radius
	camScale := 1.0
	if s.camParams.Cam != nil {
		camScale = s.camParams.Cam.Scale
	}
	maxRadius := winW * 2 / camScale
	if effRadius > maxRadius {
		effRadius = maxRadius
	}

	scale := effRadius / sectorLen
	s.sprite.SetScale(scale, scale)
	s.sprite.SetColor(s.color)
	s.sprite.SetAlpha(s.alpha)
}

func (s *Sector) SetCenter(center v2.V2) {
	s.center = center
}

func (s *Sector) SetRadius(radius float64) {
	s.radius = radius
}

func (s *Sector) SetCenterRadius(center v2.V2, radius float64) {
	s.center = center
	s.radius = radius
}

func (s *Sector) SetAngles(start, end float64) {
	s.startAng, s.endAng = NormAngRange(start, end)
}

func (s *Sector) SetColor(color color.Color) {
	s.color = color
}

func (s *Sector) SetAlpha(alpha float64) {
	s.alpha = alpha
}
