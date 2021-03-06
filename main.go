package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

const (
	TEST_XRES   = 1600
	TEST_YRES   = 900
	TEST_MIDDLE = 50
)

func main() {
	fmt.Println("Testing WorldStructure")
	GE.Init("", 30)
	wrld := GE.GetWorldStructure(0, 0, TEST_XRES, TEST_YRES, 100, 100, 32, 18)
	wrld.LoadTiles("./res/tiles")
	wrld.LoadStructureObjs("./res/structures")

	img, err := GE.LoadImg("./res/testMap.png")
	GE.ShitImDying(err)
	wrld.TileMat.FillFromImage((*img).(*image.Gray))

	wrld.SetLightStats(30, 220)
	wrld.SetMiddleSmooth(TEST_MIDDLE, TEST_MIDDLE)
	wrld.UpdateLIdxMat()

	sprites, err := GE.LoadEbitenImg("./res/anims/arrow.png")
	GE.ShitImDying(err)
	psAnim := GE.GetAnimation(0, 0, 1, 1, 16, 6, sprites)
	pf := GE.GetNewParticleFactory(100, 30, psAnim)
	ps := GE.GetNewParticleSystem(10, pf)

	wrld.Add_Drawables = wrld.Add_Drawables.Add(ps)
	wrld.UpdateObjMat()
	wrld.UpdateObjDrawables()

	game := &TestGame{wrld, ps, pf, 0}
	ebiten.SetWindowSize(TEST_XRES, TEST_YRES)
	ebiten.SetWindowTitle("WorldStructure Test")
	ebiten.SetMaxTPS(30)
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type TestGame struct {
	wrld         *GE.WorldStructure
	ps           *GE.ParticleSystem
	pf           *GE.ParticleFactory
	frameCounter int
}

func (g *TestGame) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		wx, wy := g.wrld.GetTileOfCoordsFP(float64(x), float64(y))
		g.ps.Add(g.pf.GetNewRandom(g.frameCounter, 1.0, wx, wy, 1, 1))
	}

	g.wrld.UpdateAllLightsIfNecassary()
	g.wrld.UpdateTime(time.Second)
	g.ps.Update(g.frameCounter)
	g.frameCounter++
	return nil
}
func (g *TestGame) Draw(screen *ebiten.Image) {
	g.wrld.Draw(screen)
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return TEST_XRES, TEST_YRES
}
