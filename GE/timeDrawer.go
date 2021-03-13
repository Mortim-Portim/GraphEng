package GE

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

/**
TimeDrawer represents a background, foreground, sun and moon image that can be arranged and drawn according to the time of the day
**/
func TimeToRel(t time.Time) float64 {
	return float64(HMStoS(t.Clock())) / 43200.0
}
func GetTimeDrawer(back, front, sun, moon *ebiten.Image, X, Y, W, H float64) (td *TimeDrawer) {
	td = &TimeDrawer{
		Back:  EbitenImgToImgObj(back),
		Front: EbitenImgToImgObj(front),
		Sun:   EbitenImgToImgObj(sun),
		Moon:  EbitenImgToImgObj(moon),
		X:     X, Y: Y, W: W, H: H,
	}
	td.Update(0)
	return
}

type TimeDrawer struct {
	Back, Front, Sun, Moon *ImageObj
	showSun, showMoon      bool
	lastPercent, Percent   float64
	X, Y, W, H             float64
}

func (td *TimeDrawer) Update(int) {
	if td.lastPercent != td.Percent {
		td.lastPercent = td.Percent
		td.Back.X, td.Back.Y, td.Back.W, td.Back.H = td.X, td.Y, td.W, td.H
		td.Front.X, td.Front.Y, td.Front.W, td.Front.H = td.X, td.Y, td.W, td.H
		td.Sun.W = td.H / 2
		td.Sun.H = td.H / 2
		td.Moon.W = td.H / 2
		td.Moon.H = td.H / 2
		padding := td.Sun.W / 2
		WMinusPadding := td.W - td.Sun.W
		td.showSun, td.showMoon = false, false
		if td.Percent >= 0 && td.Percent < 0.5 {
			td.showMoon = true
			yR := 1 - math.Sin((td.Percent+0.5)*math.Pi)
			xM := td.X + padding + WMinusPadding/2 + WMinusPadding*td.Percent
			yM := td.Y + padding + yR*td.Moon.H
			td.Moon.SetMiddle(xM, yM)
		} else if td.Percent >= 0.5 && td.Percent < 1 {
			td.showSun = true
			yR := 1 - math.Sin((td.Percent-0.5)*math.Pi)
			xM := td.X + padding + WMinusPadding*(td.Percent-0.5)
			yM := td.Y + padding + yR*td.Sun.H
			td.Sun.SetMiddle(xM, yM)
		} else if td.Percent >= 1 && td.Percent < 1.5 {
			td.showSun = true
			yR := 1 - math.Sin((td.Percent-0.5)*math.Pi)
			xM := td.X + padding + WMinusPadding/2 + WMinusPadding*(td.Percent-1)
			yM := td.Y + padding + yR*td.Sun.H
			td.Sun.SetMiddle(xM, yM)
		} else if td.Percent >= 1.5 && td.Percent < 2 {
			td.showMoon = true
			yR := 1 - math.Sin((td.Percent-1.5)*math.Pi)
			xM := td.X + padding + WMinusPadding*(td.Percent-1.5)
			yM := td.Y + padding + yR*td.Moon.H
			td.Moon.SetMiddle(xM, yM)
		}
	}
}
func (td *TimeDrawer) Draw(screen *ebiten.Image) {
	td.Back.Draw(screen)
	if td.showSun {
		td.Sun.Draw(screen)
	}
	if td.showMoon {
		td.Moon.Draw(screen)
	}
	td.Front.Draw(screen)
}
func (td *TimeDrawer) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return td.Update, td.Draw
}
func (td *TimeDrawer) Start(screen *ebiten.Image, data interface{}) {}
func (td *TimeDrawer) Stop(screen *ebiten.Image, data interface{})  {}
