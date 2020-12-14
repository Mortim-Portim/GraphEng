package GE

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

func GetButtonListFromStrings(text []string, X, Y, W, H, space float64, txtcolor, backcolor color.Color) *ButtonList {
	btns := make([]*Button, len(text))
	btnlist := &ButtonList{clapped: true, height: H}

	for i, line := range text {
		btn := GetSizedTextButton(line, StandardFont, X, Y+float64(i)*(H+space), W, H, txtcolor, backcolor)
		btn.Data = i
		btn.RegisterOnLeftEvent(func(btn *Button) {
			if btn.LPressed {
				if !btnlist.clapped {
					btnlist.MoveUpDown(btnlist.current - btn.Data.(int))
				}

				btnlist.clapped = !btnlist.clapped
				btnlist.current = btn.Data.(int)
			}
		})
		btns[i] = btn
	}

	btnlist.btns = btns
	return btnlist
}

type ButtonList struct {
	btns []*Button

	current, max int
	clapped bool
	height  float64
}

func (l *ButtonList) Reset() {
	l.current = 0
}
func (l *ButtonList) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return l.Update, l.Draw
}
func (l *ButtonList) Start(screen *ebiten.Image, data interface{}) { l.Reset() }
func (l *ButtonList) Stop(screen *ebiten.Image, data interface{})  {}

func (l *ButtonList) Update(frame int) {
	if !l.clapped {
		_, dy := ebiten.Wheel()

		if dy > 0 {
			l.MoveUpDown(1)
		}

		if dy < 0 {
			l.MoveUpDown(-1)
		}
	}

	if l.clapped {
		l.btns[l.current].Update(frame)
	} else {
		for i := l.current; i < len(l.btns); i++ {
			l.btns[i].Update(frame)
		}
	}
}

func (l *ButtonList) Draw(screen *ebiten.Image) {
	if l.clapped {
		l.btns[l.current].Draw(screen)
	} else {
		for i := l.current; i < len(l.btns); i++ {
			l.btns[i].Draw(screen)
		}
	}
}

func (l *ButtonList) MoveUpDown(Y int) {
	fmt.Println(l.current, Y)
	if l.current-Y < 0 {
		return
	}

	for _, btn := range l.btns {
		btn.Img.Y += float64(Y) * l.height
	}

	l.current -= Y
}
