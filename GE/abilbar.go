package GE

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

//Creates an bar to display stats
//img:				backgrundimage
//x, y, w			size of image
//dx, dy, dw, dh:	size of the bar, x and y relativ to the image pos and will be scaled like the imgage
//color:			color of the bar
func GetAbilbar(img *ebiten.Image, x, y, w, dx, dy, dw, dh float64, barcolor, backcol color.Color) *Abilbar {
	imgobj := EbitenImgToImgObj(img)
	imgobj.X, imgobj.Y = x, y
	scale := w / imgobj.W
	imgobj.ScaleToX(w)

	backimg := GetEmptyImage(2, 2)
	backimg.Fill(backcol)
	back := NewImageObj(nil, backimg, x+dx*scale, y+dy*scale, dw*scale, dh*scale, 0)

	barimg := GetEmptyImage(2, 2)
	barimg.Fill(barcolor)
	bar := NewImageObj(nil, barimg, x+dx*scale, y+dy*scale, dw*scale, dh*scale, 0)
	return &Abilbar{imgobj, bar, back, dw * scale}
}

type Abilbar struct {
	img  *ImageObj
	bar  *ImageObj
	back *ImageObj

	w float64
}

func (ab *Abilbar) Set(percent float64) {
	ab.bar.W = ab.w * percent
}

func (ab *Abilbar) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return ab.Update, ab.Draw
}

func (ab *Abilbar) Start(screen *ebiten.Image, data interface{}) {}

func (ab *Abilbar) Stop(screen *ebiten.Image, data interface{}) {}

func (ab *Abilbar) Update(frame int) {}

func (ab *Abilbar) Draw(screen *ebiten.Image) {
	ab.back.Draw(screen)
	ab.bar.Draw(screen)
	ab.img.Draw(screen)
}
