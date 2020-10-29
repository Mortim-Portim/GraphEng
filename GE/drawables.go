package GE

import (
	"github.com/hajimehoshi/ebiten"
	"reflect"
	"errors"
	"sort"
	"math"
	"fmt"
)

//type GetYPos func()(y float64, isBack bool)
//type Draw func(screen *ebiten.Image, x, y float64, lv int16, xStart, yStart, sqSize float64)

func GetDrawables() *Drawables {
	d := Drawables(make([]Drawable, 0))
	return &d
}


type Drawable interface {
	Height() float64
	GetPos()(x, y float64, layer int8)
	Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64)
}
type Drawables []Drawable
func (d Drawables) Len() int {return len(d)}
func (d Drawables) Swap(i, j int) {d[i], d[j] = d[j], d[i]}
func (d Drawables) Less(i, j int) bool {
	_,y1,layer1 := d[i].GetPos()
	_,y2,layer2 := d[j].GetPos()
	if layer1 < layer2 {
		return true
	}else if layer1 > layer2 {
		return false
	}
	if math.Abs(y1-y2) < 0.5 {
		if d[i].Height() < d[j].Height() {
			return true
		}
		return false
	}
	return y1 < y2
}
func (d Drawables) Sort() {
	sort.Sort(d)
}
func (dp *Drawables) Add(obj Drawable) (*Drawables) {
	d := append(*dp, obj)
	return &d
}
func (d Drawables) Remove(obj interface{}) error {
	objType := reflect.TypeOf(obj)
	i := sort.Search(len(d), func(idx int) bool {
			tp := reflect.TypeOf(d[idx])
			if !tp.AssignableTo(objType) {return false}
			if d[idx] == obj {return true}
			return false
	})
	if i < 0 || i >= len(d) {
		return errors.New(fmt.Sprintf("Cannot remove %s, %v does not exist", objType.String(), obj))
	}
	d.removeIdx(i)
	return nil
}

func (d Drawables) removeIdx(i int) {
	d[i] = d[len(d)-1]
	d = d[:len(d)-1]
}

type WObj struct {
	img *DayNightAnim
	layer int8
	HitBox, DrawBox *Rectangle
	frame, squareSize int
}
func (o *WObj) SetToXY(x,y float64) {
	o.HitBox.MoveTo(&Point{x,y})
	w,h := o.img.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}
func (o *WObj) GetPos() (float64, float64, int8) {
	pnt := o.HitBox.GetMiddle()
	return pnt.X+0.5,pnt.Y+0.5,o.layer
}
func (o *WObj) Height() float64 {
	return o.HitBox.Bounds().Y
}
func (o *WObj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, sqSize, xStart, yStart float64) {
	y := (o.DrawBox.Min().Y-leftTopY)*sqSize
	x := (o.DrawBox.Min().X-leftTopX)*sqSize
	o.img.SetParams(x+xStart,y+yStart, float64(o.DrawBox.Bounds().X)*sqSize, float64(o.DrawBox.Bounds().Y)*sqSize)
	o.img.LightLevel = lv
	o.img.DrawAnim(screen)
	o.frame ++
}
func GetWObjFromParams(img *ebiten.Image, p *Params) (o *WObj) {
	layer := int8(p.Get("Layer"))
	spW := int(p.Get("spriteWidth"))
	uP := int(p.Get("updatePeriod"))
	hitBoxW := p.Get("hitBoxWidth")-1
	hitBoxH := p.Get("hitBoxHeight")-1
	XPos := p.Get("XPos")
	YPos := p.Get("YPos")
	anim := GetDayNightAnim(1,1,1,1, spW, uP, img)
	o = GetWObj(anim, hitBoxW, hitBoxH, XPos, YPos, int(p.Get("squareSize")), layer)
	return
}
//Returns a StructureObj
func GetWObj(img *DayNightAnim, HitboxW,HitboxH, x, y float64, squareSize int, layer int8) (o *WObj) {
	o = &WObj{img:img, layer:layer, squareSize:squareSize}
	o.HitBox = GetRectangle(x,y, x+HitboxW, y+HitboxH)
	if img != nil {
		o.img.Update(0)
	}
	o.SetToXY(x, y)
	return
}



/**
func (dp *Drawables) AddImageObj(drawBox *Rectangle, obj *ImageObj) (*Drawables) {
	pnt := drawBox.GetMiddle()
	bounds := drawBox.Bounds()
	gx := func()float64{
		return pnt.X
	}
	gy := func()float64{
		return pnt.Y
	}
	dr := func(screen *ebiten.Image, x, y float64, _ int16, sqSize float64){
		obj.W = bounds.X*sqSize
		obj.H = bounds.Y*sqSize
		obj.SetMiddle(x,y)
		obj.DrawImageObj(screen)
	}
	dobj := &Drawable{gx,gy,dr, obj}
	d := append(*dp, dobj)
	return &d
}
func (dp *Drawables) AddDayNightImg(drawBox *Rectangle, obj *DayNightImg) (*Drawables) {
	pnt := drawBox.GetMiddle()
	bounds := drawBox.Bounds()
	gx := func()float64{
		return pnt.X
	}
	gy := func()float64{
		return pnt.Y
	}
	dr := func(screen *ebiten.Image, x, y float64, lv int16, sqSize float64){
		obj.SetParams(0,0,bounds.X*sqSize, bounds.Y*sqSize)
		obj.SetMiddle(x,y)
		obj.Draw(screen, float64(lv)/255.0)
	}
	dobj := &Drawable{gx,gy,dr, obj}
	d := append(*dp, dobj)
	return &d
}
**/