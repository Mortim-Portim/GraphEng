package GE

import (
	"github.com/hajimehoshi/ebiten"
	"reflect"
	"errors"
	"strings"
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
func (dp *Drawables) Remove(obj interface{}) error {
	objType := reflect.TypeOf(obj)
	i := sort.Search(len(*dp), func(idx int) bool {
			tp := reflect.TypeOf((*dp)[idx])
			if !tp.AssignableTo(objType) {return false}
			if (*dp)[idx] == obj {return true}
			return false
	})
	if i < 0 || i >= len(*dp) {
		return errors.New(fmt.Sprintf("Cannot remove %s, %v does not exist", objType.String(), obj))
	}
	(*dp).removeIdx(i)
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
	Name string
}
//Copys the WObj
func (o *WObj) Copy() (o2 *WObj) {
	o2 = &WObj{o.img.Copy(), o.layer, o.HitBox.Copy(), o.DrawBox.Copy(), o.frame, o.squareSize, o.Name}
	return
}
//Updates the animation
func (o *WObj) Update(frame int) {
	o.img.Update(frame)
}
//Moves the WObj by delta
func (o *WObj) MoveBy(dx,dy float64) {
	if dx != 0 || dy != 0 {
		x := o.HitBox.Min().X+dx;y := o.HitBox.Min().Y+dy
		o.SetToXY(x,y)
	}
}
//Sets the WObjs top left to x and y
func (o *WObj) SetToXY(x,y float64) {
	o.HitBox.MoveTo(&Point{x,y})
	w,h := o.img.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}
//Sets the WObjs middle to x and y
func (o *WObj) SetPos(x,y float64) {
	w,h := o.Bounds()
	o.SetToXY(x+w/2, y+h/2)
}
//Returns the Position of the middle of the WObj
func (o *WObj) GetPos() (float64, float64, int8) {
	pnt := o.HitBox.GetMiddle()
	return pnt.X+0.5,pnt.Y+0.5,o.layer
}
//Sets the layer the WObj is drawn to
func (o *WObj) SetLayer(l int8) {
	o.layer = l
}
//Returns the Height-1 of the WObj
func (o *WObj) Height() float64 {
	return o.HitBox.Bounds().Y
}
//Returns the real Bounds of the WObj
func (o *WObj) Bounds() (float64, float64) {
	bnds := o.HitBox.Bounds()
	return bnds.X+1, bnds.Y+1
}
//Draws the WObj to the screen
func (o *WObj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	y := (o.DrawBox.Min().Y-leftTopY)*sqSize
	x := (o.DrawBox.Min().X-leftTopX)*sqSize
	o.img.SetParams(x+xStart,y+yStart, float64(o.DrawBox.Bounds().X)*sqSize, float64(o.DrawBox.Bounds().Y)*sqSize)
	o.img.LightLevel = lv
	o.img.DrawAnim(screen)
	o.frame ++
}
//Sets the animation the WObj is using
func (o *WObj) SetAnim(anim *DayNightAnim) {
	if anim == nil {
		panic("anim should never be nil")
	}
	o.img = anim
}
/**
Params.txt:
squareSize:		[1-NaN]
spriteWidth: 	[1-NaN]
updatePeriod: 	[0-NaN]
hitBoxWidth:	[1-NaN]
hitBoxHeight:	[1-NaN]
XPos:			[0-NaN]
YPos:			[0-NaN]
layer:			[-128-128]
**/
func GetWObjFromParams(img *ebiten.Image, p *Params, name string) (o *WObj) {
	layer := int8(p.Get("layer"))
	spW := int(p.Get("spriteWidth"))
	uP := int(p.Get("updatePeriod"))
	hitBoxW := p.Get("hitBoxWidth")-1
	hitBoxH := p.Get("hitBoxHeight")-1
	XPos := p.Get("XPos")
	YPos := p.Get("YPos")
	sqSize := int(p.Get("squareSize"))
	anim := GetDayNightAnim(1,1,1,1, spW, uP, img)
	o = GetWObj(anim, hitBoxW, hitBoxH, XPos, YPos, sqSize, layer, name)
	return
}
//provide at least on path
func GetWObjFromPath(name string, path ...string) (*WObj, error) {
	if len(path) == 0 {
		panic("No path given")
	}
	var pathImg, pathParams string
	if len(path) == 1 {
		pathImg = path[0]
		pathParams = path[0]
	}else if len(path) == 2 {
		pathImg = path[0]
		pathParams = path[1]
	}
	
	img, err := LoadEbitenImg(pathImg+".png")
	if err != nil {return nil,err}
	ps := &Params{}; err = ps.LoadFromFile(pathParams+".txt")
	if err != nil {return nil,err}
	return GetWObjFromParams(img, ps, name), nil
}
//Loads all WObjs from a directory
func LoadAllWObjs(folderPath string) (map[string]*WObj, error) {
	if folderPath[len(folderPath)-1:] != "/" {
		folderPath += "/"
	}
	files, err := OSReadDir(folderPath)
	if err != nil {return nil, err}
	
	var currentError error
	objs := make(map[string]*WObj)
	names := make([]string, 0)
	for _,f := range(files) {
		n := strings.Split(f, ".")[0]
		if !IsStringInList(n, names) {
			obj, err := GetWObjFromPath(n, folderPath+n)
			currentError = err
			names = append(names, n)
			objs[n] = obj
		}
	}
	return objs, currentError
}
func GetWObj(img *DayNightAnim, HitboxW,HitboxH, x, y float64, squareSize int, layer int8, name string) (o *WObj) {
	o = &WObj{img:img, layer:layer, squareSize:squareSize, Name:name}
	o.HitBox = GetRectangle(x,y, x+HitboxW, y+HitboxH)
	if img != nil {
		o.img.Update(0)
	}
	o.SetToXY(x, y)
	return
}


func IsStringInList(s string, l []string) bool {
	for _,str := range(l) {
		if str == s {
			return true
		}
	}
	return false
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