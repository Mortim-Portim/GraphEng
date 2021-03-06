package GE

import (
	"image"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

/**
DayNightAnim represents an Animtation, that can be drawn using different light levels
see DayNightImg for further infos

DayNightAnim implements UpdateAble
**/
type DayNightAnim struct {
	*DayNightImg
	sprites, current, spriteWidth, spriteHeight, UpdatePeriod, lastFrame int
	LightLevel                                                           int16
	spriteSheet                                                          *ebiten.Image
	OnAnimFinished                                                       func(a *DayNightAnim)
}

func (a *DayNightAnim) Copy() *DayNightAnim {
	return &DayNightAnim{a.DayNightImg.Clone(), a.sprites, a.current, a.spriteWidth, a.spriteHeight, a.UpdatePeriod, a.lastFrame, a.LightLevel, a.spriteSheet, a.OnAnimFinished}
}
func (a *DayNightAnim) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return a.Update, a.DrawAnim
}
func (a *DayNightAnim) Start(screen *ebiten.Image, data interface{}) {}
func (a *DayNightAnim) Stop(screen *ebiten.Image, data interface{})  {}
func (a *DayNightAnim) Update(frame int) {
	if a.UpdatePeriod > 0 && a.lastFrame != frame && frame%a.UpdatePeriod == 0 {
		a.lastFrame = frame
		a.current++
		if a.current >= a.sprites {
			a.current = 0
			if a.OnAnimFinished != nil {
				a.OnAnimFinished(a)
			}
		}
		a.DayNightImg.SetDay(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight/2)).(*ebiten.Image))
		a.DayNightImg.SetNight(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, a.spriteHeight/2, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image))

	} else if a.UpdatePeriod == 0 && (a.DayNightImg.day.Img == nil || a.DayNightImg.night.Img == nil) {
		a.DayNightImg.SetDay(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight/2)).(*ebiten.Image))
		a.DayNightImg.SetNight(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, a.spriteHeight/2, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image))
	}
}
func (a *DayNightAnim) DrawAnim(screen *ebiten.Image) {
	a.Draw(screen, float64(a.LightLevel)/float64(255))
}
