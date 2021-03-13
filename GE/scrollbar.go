package GE

import (
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
)

/**
A ScrollBar is a horizontal bar that can be used to get input
index of the scrollbar can be between a min and max value

				  ++++
------------------++++-----------------
				  ++++
				   2

ScrollBar implements UpdateAble
**/

type ScrollBar struct {
	ImageObj
	pointer, value            *ImageObj
	focused, hideValue        bool
	min, max, length, current int
	stepsize, relAbsPos       float64
	ttf                       *truetype.Font
	OnChange                  func(b *ScrollBar)
}

//Registers a method to be called when the index changes
func (b *ScrollBar) RegisterOnChange(OnChange func(*ScrollBar)) {
	b.OnChange = OnChange
}

//Returns the current index of the scrollbar
func (b *ScrollBar) Current() int { return b.current }
func (b *ScrollBar) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return b.Update, b.Draw
}
func (b *ScrollBar) Start(screen *ebiten.Image, data interface{}) {}
func (b *ScrollBar) Stop(screen *ebiten.Image, data interface{})  {}
func (b *ScrollBar) Update(frame int) {

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.focused = false
	} else {
		x, y := ebiten.CursorPosition()
		if int(b.X) <= x && x < int(b.X+b.W) && int(b.Y) <= y && y < int(b.Y+b.H) {
			b.focused = true
		}
		if b.focused {
			b.CheckChange((float64(x) - b.X) / b.W)
		}
	}
}
func (b *ScrollBar) Draw(screen *ebiten.Image) {
	b.DrawImageObj(screen)
	b.pointer.DrawImageObj(screen)
	if !b.hideValue {
		b.value.DrawImageObj(screen)
	}
}
func (b *ScrollBar) UpdatePos() {
	b.relAbsPos = float64(b.current - b.min)
	b.pointer.SetMiddleX(b.X + b.stepsize*(b.relAbsPos))
	b.value = GetTextImage(fmt.Sprintf("%v", b.current), b.pointer.X, b.pointer.Y+b.pointer.H, b.pointer.H/2, b.ttf, &color.RGBA{255, 255, 255, 255}, &color.RGBA{0, 0, 0, 0})
	b.value.SetMiddleX(b.pointer.X + b.pointer.W/2)

	if b.OnChange != nil {
		b.OnChange(b)
	}
}
func (b *ScrollBar) CheckChange(x float64) {
	b.current = int(x*float64(b.length)) + b.min
	if b.current < b.min {
		b.current = b.min
	}
	if b.current > b.max {
		b.current = b.max
	}
	b.UpdatePos()
}
func (b *ScrollBar) HideValue() { b.hideValue = true }
func (b *ScrollBar) ShowValue() { b.hideValue = false }
