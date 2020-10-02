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
	//"github.com/golang/freetype/truetype"
	//"github.com/hajimehoshi/ebiten/inpututil"
	//"golang.org/x/image/font"
	//"github.com/nfnt/resize"
	//"image/color"
	//"fmt"
	//"math"
)

type ScrollBar struct {
	ImageObj
	pointer, value *ImageObj

	min, max, current int
	stepsize          float64
}

func GetImageScrollbar(X, Y, W, H float64, bar, pointer *ebiten.Image, min, max, current int) (b *ScrollBar) {
	b = &ScrollBar{min: min, max: max, current: current}
	b.Img = bar
	b.X = X
	b.Y = Y
	b.W = W
	b.H = H
	b.pointer = &ImageObj{W: H, H: H}
	b.value = &ImageObj{W: H, H: H}
	b.pointer.Img = pointer
	b.stepsize = W / float64(max-min)
	return
}
