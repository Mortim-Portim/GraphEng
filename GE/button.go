package GE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Button struct {
	Img         *ImageObj
	dark, light *ebiten.Image

	LPressed, RPressed, LastL, LastR, Active, DrawDark, ChangeDrawDarkOnLeft bool
	onPressLeft, onPressRight                                                func(b *Button)
	Data                                                                     interface{}
}

func (b *Button) UpImg() *ebiten.Image {
	return b.light
}
func (b *Button) DownImg() *ebiten.Image {
	return b.dark
}

/**
TODO
Use JustDown and check for other buttons that are pressed
**/

/**
Button represents a struct, which should be updated every frame

Button calls onPressLeft() or onPressRight() if being clicked by the mouse

It can be created from a Image or using Text

Button implements UpdateAble
**/
func (b *Button) Copy() *Button {
	return &Button{b.Img.Copy(), b.dark, b.light, b.LPressed, b.RPressed, b.LastL, b.LastR, b.Active, b.DrawDark, b.ChangeDrawDarkOnLeft, b.onPressLeft, b.onPressRight, b.Data}
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
		x, y := ebiten.CursorPosition()
		hasFocus := int(b.Img.X) <= x && x < int(b.Img.X+b.Img.W) && int(b.Img.Y) <= y && y < int(b.Img.Y+b.Img.H)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && hasFocus {
			b.LPressed = true
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			b.LPressed = false
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && hasFocus {
			b.RPressed = true
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			b.RPressed = false
		}

		if b.LPressed != b.LastL {
			if b.onPressLeft != nil {
				if b.ChangeDrawDarkOnLeft {
					b.DrawDark = !b.DrawDark
				}
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
		oldX := b.Img.X
		oldY := b.Img.Y
		if b.LPressed || b.RPressed {
			b.Img.X = (oldX + b.Img.W*MoveOnButtonDown)
			b.Img.Y = (oldY + b.Img.H*MoveOnButtonDown)
		}
		if b.DrawDark {
			b.Img.Img = b.dark
		} else {
			b.Img.Img = b.light
		}
		b.Img.DrawImageObj(screen)
		if b.LPressed || b.RPressed {
			b.Img.X = (oldX)
			b.Img.Y = (oldY)
		}
	}
}
