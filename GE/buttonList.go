package GE

import (
	//"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/**
ButtonList represents a list of buttons, that the user can scroll through

PreRenderPerfect does nothing

ButtonList implements UpdateAble
**/
func GetButtonListFromStrings(X, Y, W, H, BtnH, ySpace, ScrollH float64, textCol, backCol color.Color, strs ...string) *ButtonList {
	btns := make([]*Button, len(strs))
	for i, str := range strs {
		btns[i] = GetTextButton(str, str, StandardFont, 0, 0, BtnH, textCol, backCol)
		btns[i].Data = str
	}
	return GetButtonList(X, Y, W, H, ySpace, ScrollH, btns...)
}
func GetButtonList(X, Y, W, H, ySpace, ScrollH float64, btns ...*Button) *ButtonList {
	return &ButtonList{X, Y, W, H, ySpace, ScrollH, 0.0, btns, false}
}

type ButtonList struct {
	X, Y, W, H, ySpace, ScrollH, currentDelta float64
	btns                                      []*Button

	PreRenderPerfect bool
}

func (l *ButtonList) Content() []*Button {
	return l.btns
}
func (l *ButtonList) Reset() {
	l.currentDelta = 0
	l.UpdateButtonPositions()
}
func (l *ButtonList) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	l.Reset()
	return l.Update, l.Draw
}
func (l *ButtonList) Start(screen *ebiten.Image, data interface{}) {
	l.Reset()
}
func (l *ButtonList) Stop(screen *ebiten.Image, data interface{}) {}
func (l *ButtonList) Update(frame int) {
	x, y := ebiten.CursorPosition()
	_, scroll := ebiten.Wheel()
	hasFocus := int(l.X) <= x && x < int(l.X+l.W) && int(l.Y) <= y && y < int(l.Y+l.H)

	if hasFocus && scroll != 0 {
		l.Scroll(float64(scroll))
		l.UpdateButtonPositions()
		l.SetActiveButtons()
	}
	for _, btn := range l.btns {
		btn.Update(frame)
	}
}
func (l *ButtonList) Draw(screen *ebiten.Image) {
	if !l.PreRenderPerfect {
		for _, btn := range l.btns {
			btn.Draw(screen)
		}
	}
}
func (l *ButtonList) Scroll(dy float64) {
	l.currentDelta -= dy * l.ScrollH
}
func (l *ButtonList) UpdateButtonPositions() {
	Y := l.Y + l.ySpace + l.currentDelta
	for _, btn := range l.btns {
		btn.Img.X = l.X + l.ySpace
		btn.Img.Y = Y
		Y += btn.Img.H + l.ySpace
	}
}
func (l *ButtonList) SetActiveButtons() {
	rect := GetRectangle(l.X, l.Y, l.X+l.W*1000, l.Y+l.H)
	for _, btn := range l.btns {
		btn.Active = btn.Img.Rectangle().Overlaps(rect)
	}
}
