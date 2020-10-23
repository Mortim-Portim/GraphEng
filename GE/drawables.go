package GE

import (
	"github.com/hajimehoshi/ebiten"
	"reflect"
	"errors"
	"sort"
	"fmt"
)
func GetDrawables() *Drawables {
	d := Drawables(make([]*Drawable, 0))
	return &d
}
type GetX func()float64
type GetY func()float64
type Draw func(screen *ebiten.Image, x, y float64, lvMat int16, sqSize float64)

type Drawable struct {
	GetX
	GetY
	Draw
	original interface{}
}
type Drawables []*Drawable
func (d Drawables) Len() int {return len(d)}
func (d Drawables) Swap(i, j int) {d[i], d[j] = d[j], d[i]}
func (d Drawables) Less(i, j int) bool {return d[i].GetY() < d[j].GetY()}
func (d Drawables) Sort() {
	sort.Sort(d)
}
func (dp *Drawables) AddStructureObjOfWorld(obj *StructureObj, wrld *WorldStructure) (*Drawables) {
	pnt := obj.HitBox.GetMiddle()
	leftTop := wrld.ObjMat.Focus().Min()
	pX := pnt.X-leftTop.X+0.5
	pY := pnt.Y-leftTop.Y+1
	gx := func()float64{
		return pX
	}
	gy := func()float64{
		return pY
	}
	dr := func(screen *ebiten.Image, x, y float64, lv int16, sqSize float64){
		obj.Draw(screen, lv, x, y, sqSize)
	}
	dobj := &Drawable{gx,gy,dr, obj}
	d := append(*dp, dobj)
	return &d
}
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
func (d Drawables) Remove(obj interface{}) error {
	objType := reflect.TypeOf(obj)
	i := sort.Search(len(d), func(idx int) bool {
			tp := reflect.TypeOf(d[idx].original)
			if !tp.AssignableTo(objType) {return false}
			if d[idx].original == obj {return true}
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