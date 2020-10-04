package GE

import (
	"github.com/hajimehoshi/ebiten"
)

type Button struct {
	Img *ImageObj
	dark, light *ebiten.Image
		
	LPressed, RPressed, LastL, LastR, Active, DrawDark bool
	onPressLeft, onPressRight func(b *Button)
	Data interface{}
}
func (b *Button) Reset() {
	b.LPressed = false
	b.RPressed = false
	b.LastL = false
	b.LastR = false
	b.Active = true
}
func (b *Button) RegisterOnEvent(OnEvent func(*Button)) {
	b.RegisterOnLeftEvent(OnEvent)
	b.RegisterOnRightEvent(OnEvent)
}
func (b *Button) RegisterOnLeftEvent(OnLeftEvent func(*Button)) {
	b.onPressLeft = OnLeftEvent
}
func (b *Button) RegisterOnRightEvent(OnRightEvent func(*Button)) {
	b.onPressRight = OnRightEvent
}
func (b *Button) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	b.Reset()
	return b.Update, b.Draw
}
func (b *Button) Start(screen *ebiten.Image, data interface{}) {
	b.Reset()
}
func (b *Button) Stop(screen *ebiten.Image, data interface{}) {
	b.Reset()
}
func (b *Button) Update(frame int) {
	if b.Active {
		b.LPressed = false
		b.RPressed = false
		x, y := ebiten.CursorPosition()
		if int(b.Img.X) <= x && x < int(b.Img.X+b.Img.W) && int(b.Img.Y) <= y && y < int(b.Img.Y+b.Img.H) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				b.LPressed = true
			}
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
				b.RPressed = true
			}
		}
		if b.LPressed != b.LastL {
			if b.onPressLeft != nil {
				b.onPressLeft(b)
			}
		}
		b.LastL = b.LPressed
		
		if b.RPressed != b.LastR {
			if b.onPressRight != nil {
				b.onPressRight(b)
			}
		}
		b.LastR = b.RPressed
	}
}
func (b *Button) Draw(screen *ebiten.Image) {
	if b.Active {
		oldX := b.Img.X; oldY := b.Img.Y
		if b.LPressed || b.RPressed {
			b.Img.X = (oldX+b.Img.W*MoveOnButtonDown)
			b.Img.Y = (oldY+b.Img.H*MoveOnButtonDown)
		}
		if b.DrawDark {
			b.Img.Img = b.dark
		}else{
			b.Img.Img = b.light
		}
		b.Img.DrawImageObj(screen)
		if b.LPressed || b.RPressed {
			b.Img.X = (oldX)
			b.Img.Y = (oldY)
		}
	}
}
