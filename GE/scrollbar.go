package GE

import (
	//"os"
	//"image"
	//_ "image/jpeg"
	//"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/golang/freetype/truetype"
	//"github.com/hajimehoshi/ebiten/inpututil"
	//"golang.org/x/image/font"
	//"github.com/nfnt/resize"
	"image/color"
	"fmt"
	"math"
)

type ScrollBar struct {
	ImageObj
	pointer, value *ImageObj

	min, max, length, current int
	stepsize, relAbsPos float64
	ttf *truetype.Font
}

func (b *ScrollBar) Init(screen *ebiten.Image, data interface{}) {
	
}
func (b *ScrollBar)	Start(screen *ebiten.Image, data interface{}) {}
func (b *ScrollBar)	Stop(screen *ebiten.Image, data interface{}) {}
func (b *ScrollBar)	Update(frame int) {
	x, y := ebiten.CursorPosition()
	if int(b.X) <= x && x < int(b.X+b.W) && int(b.Y) <= y && y < int(b.Y+b.H) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.CheckChange((float64(x)-b.X)/b.W)
		}
	}
}
func (b *ScrollBar)	Draw(screen *ebiten.Image, frame int) {
	b.DrawImageObj(screen)
	b.pointer.DrawImageObj(screen)
	b.value.DrawImageObj(screen)
}
func (b *ScrollBar) UpdatePos() {
	b.relAbsPos = float64(b.current-b.min)
	b.pointer.SetMiddleX(b.X + b.stepsize*(b.relAbsPos))
	b.value = GetTextImage(fmt.Sprintf("%v",b.current), b.pointer.X, b.pointer.Y+b.pointer.H, b.pointer.H/2, b.ttf, &color.RGBA{255,255,255,255}, &color.RGBA{0,0,0,0})
	b.value.SetMiddleX(b.pointer.X+b.pointer.W/2)
}
func (b *ScrollBar) CheckChange(x float64) {
	xC := x-(float64(b.current-b.min)/float64(b.length))
	xCAbs := math.Abs(xC)*b.W
	if xCAbs > b.stepsize*0.7 {
		steps := int(xCAbs/(b.stepsize*0.7))
		if xC < 0 {
			b.current -= steps
		}else{
			b.current += steps
		}
		if b.current < b.min {
			b.current = b.min
		}
		if b.current > b.max {
			b.current = b.max
		}
		b.UpdatePos()
	}
}
