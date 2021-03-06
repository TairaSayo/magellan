package main

import (
	. "github.com/Shnifer/magellan/commons"
	hls "github.com/gerow/go-color"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
	"strconv"
)

type pOpts struct {
	parent   string
	t        string
	orbit    float64
	period   float64
	phase    float64
	size     float64
	r10      float64
	maxG     float64
	shps     [15]int
	minerals []int
}

func (o pOpts) gp() *GalaxyPoint {
	okr := func(x float64) float64 {
		const sgn = 100
		return float64(int(x*sgn)) / sgn
	}

	count := 0
	clr := colornames.White
	switch o.t {
	case GPT_STAR:
		count = Opts.StarANCount
		clr = randomRed()
	case GPT_HARDPLANET:
		count = Opts.HardANCount
	case GPT_GASPLANET:
		count = Opts.GasANCount
	case GPT_ASTEROID:
		count = Opts.AsteroidANCount
	default:
		clr = randBright()
	}

	massSizeK := KDev(Opts.SizeMassDevPercent)

	zd := o.r10 / 3 * massSizeK
	maxG := o.maxG * massSizeK
	mass := maxG * zd * zd

	signatures := sphs2sigs(o.shps)

	return &GalaxyPoint{
		ParentID:   o.parent,
		Type:       o.t,
		SpriteAN:   sAN(o.t, count),
		Orbit:      okr(o.orbit),
		Period:     okr(o.period),
		AngPhase:   okr(o.phase),
		Size:       okr(o.size * massSizeK),
		Mass:       okr(mass),
		GDepth:     okr(zd),
		Emissions:  nil,
		Signatures: signatures,
		Minerals:   o.minerals,
		Color:      clr,
	}
}

func randBright() color.RGBA {
	rgb := hls.HSL{
		S: 0.5 + 0.5*rand.Float64(),
		L: 0.8 + 0.2*rand.Float64(),
		H: rand.Float64(),
	}.ToRGB()
	return color.RGBA{
		R: uint8(rgb.R * 255),
		G: uint8(rgb.G * 255),
		B: uint8(rgb.B * 255),
		A: 255,
	}
}

func randomRed() color.RGBA {
	return color.RGBA{
		R: uint8(200 + rand.Intn(56)),
		G: uint8(rand.Intn(100)),
		B: uint8(rand.Intn(100)),
		A: 255,
	}
}

func sAN(t string, count int) string {
	if count == 0 {
		return ""
	} else {
		return t + "-" + strconv.Itoa(rand.Intn(count)+1)
	}
}
