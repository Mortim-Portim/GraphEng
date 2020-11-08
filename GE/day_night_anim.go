package GE

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type DayNightAnim struct {
	*DayNightImg
	sprites, current, spriteWidth, spriteHeight, UpdatePeriod int
	LightLevel int16
	spriteSheet *ebiten.Image
}
func (a *DayNightAnim) Copy() (*DayNightAnim) {
	return &DayNightAnim{a.DayNightImg.Clone(), a.sprites, a.current, a.spriteWidth, a.spriteHeight, a.UpdatePeriod, a.LightLevel, a.spriteSheet}
}
func (a *DayNightAnim) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return a.Update, a.DrawAnim
}
func (a *DayNightAnim) Start(screen *ebiten.Image, data interface{}) {}
func (a *DayNightAnim) Stop(screen *ebiten.Image, data interface{}) {}
func (a *DayNightAnim) Update(frame int) {
	if a.UpdatePeriod > 0 && frame%a.UpdatePeriod == 0 {
		a.current ++
		if a.current >= a.sprites {
			a.current = 0
		}
		a.DayNightImg.SetDay(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight/2)).(*ebiten.Image))
		a.DayNightImg.SetNight(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, a.spriteHeight/2, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image))
		
	}else if a.UpdatePeriod == 0 && (a.DayNightImg.day.Img == nil || a.DayNightImg.night.Img == nil) {
		a.DayNightImg.SetDay(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight/2)).(*ebiten.Image))
		a.DayNightImg.SetNight(a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, a.spriteHeight/2, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image))
	}
}
func (a *DayNightAnim) DrawAnim(screen *ebiten.Image) {
	a.Draw(screen, float64(a.LightLevel)/float64(255))
}