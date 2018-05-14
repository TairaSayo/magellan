package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"time"
)

const resPath = "res/pilot/"
const texPath = "res/textures/"

var (
	WinW int
	WinH int
)

var last time.Time
var Data commons.TData
var dt float64

func mainLoop(window *ebiten.Image) error {
	input.Update()

	if Data.CommonData.PilotData != nil {
		Data.CommonData.PilotData.MsgID++
	}
	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %v\ndt = %.3f\n", fps, dt)
	ebitenutil.DebugPrint(window, msg)

	t := time.Now()
	dt = t.Sub(last).Seconds()
	last = t

	return nil
}

func main() {
	if DEFVAL.DoProf {
		commons.StartProfile(DEFVAL.CpuProfFileName)
		defer commons.StopProfile(DEFVAL.CpuProfFileName, DEFVAL.MemProfFileName)
	}

	WinW = DEFVAL.WinW
	WinH = DEFVAL.WinH

	graph.SetScreenSize(WinW, WinH)

	initClient()
	input.LoadConf(resPath)

	Data = commons.NewData()

	draw.InitTexAtlas(texPath)
	commons.SetGravityConsts(DEFVAL.GravityConst, DEFVAL.WarpGravityConst)

	createScenes()

	Client.Start()
	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)

	last = time.Now()

	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Fatal(err)
	}
}
