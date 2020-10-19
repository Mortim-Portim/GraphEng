package GE

import (
	"github.com/hajimehoshi/ebiten"
)

type ButtonList struct {
	btns []*Button
	
	displayBtns, current int
	X,Y, W,H, space float64
}
func (l *ButtonList) Reset() {
	l.current = 0
}
func (l *ButtonList) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	l.Reset()
	return l.Update, l.Draw
}
func (l *ButtonList) Start(screen *ebiten.Image, data interface{}) {l.Reset()}
func (l *ButtonList) Stop(screen *ebiten.Image, data interface{}) {}

func (l *ButtonList) Update(_ int) {
	x, y := ebiten.CursorPosition()
	if int(l.X) <= x && x < int(l.X+l.W) && int(l.Y) <= y && y < int(l.Y+l.H) {
		_, dy := ebiten.Wheel()
		l.current -= int(dy)
		if l.current < 0 {
			l.current = 0
		}
		if l.current >= len(l.btns) {
			l.current = len(l.btns)-1
		}
	}
}
func (l *ButtonList) Draw(screen *ebiten.Image) {
	for i := l.current; i < l.displayBtns + l.current; i++ {
		
	}
}
