package GE

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Animation struct {
	ImageObj
	sprites, current, spriteWidth, spriteHeight, UpdatePeriod int
	
	spriteSheet *ebiten.Image
}

func (a *Animation) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return a.Update, a.Draw
}
func (a *Animation) Start(screen *ebiten.Image, data interface{}) {}
func (a *Animation) Stop(screen *ebiten.Image, data interface{}) {}
func (a *Animation) Update(frame int) {
	if a.UpdatePeriod > 0 && frame%a.UpdatePeriod == 0 {
		a.current ++
		if a.current >= a.sprites {
			a.current = 0
		}
	}
}
func (a *Animation) Draw(screen *ebiten.Image) {
	a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	a.DrawImageObj(screen)
}