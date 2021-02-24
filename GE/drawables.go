package GE

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten"
)

/**
Drawables represents a list of Drawable, that can be drawn in a worldstructure and have a position in it
they can be sorted in order to draw the top ones first

WObj implements Drawable and represents any object, that has a DayNightAnim and a position in the worldstructure
**/
func GetDrawables() *Drawables {
	d := Drawables(make([]Drawable, 0))
	return &d
}

type Drawable interface {
	GetDrawBox() *Rectangle
	GetPos() (x, y float64, layer int8)
	Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64)
}
type Drawables []Drawable

func (d Drawables) Len() int      { return len(d) }
func (d Drawables) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d Drawables) Less(i, j int) bool {
	x1, y1, layer1 := d[i].GetPos()
	x2, y2, layer2 := d[j].GetPos()

	if math.Abs(y1-y2) < 0.1 {
		if layer1 < layer2 {
			return true
		} else if layer1 > layer2 {
			return false
		} else {
			return x1 < x2
		}
	}

	return y1 < y2
}
func (d Drawables) Sort() {
	sort.Sort(d)
}
func (dp *Drawables) Clear() *Drawables {
	*dp = Drawables(make([]Drawable, 0))
	return dp
}
func (dp *Drawables) Add(obj Drawable) *Drawables {
	*dp = append(*dp, obj)
	return dp
}
func (dp *Drawables) Remove(obj Drawable) (error, *Drawables) {
	idx := -1
	for i := 0; i < len(*dp); i++ {
		if AreInterfacesEqual((*dp)[i], obj) {
			idx = i
			break
		}
	}
	if idx < 0 || idx >= len(*dp) {
		return errors.New(fmt.Sprintf("Cannot remove %p from %v does not exist", obj, *dp)), nil
	}
	(*dp)[idx] = (*dp)[len(*dp)-1]
	(*dp) = (*dp)[:len(*dp)-1]
	return nil, dp
}
func AreInterfacesEqual(i1 interface{}, i2 interface{}) bool {
	pnt1 := fmt.Sprintf("%p", i1)
	pnt2 := fmt.Sprintf("%p", i2)
	return pnt1 == pnt2
}

type WObj struct {
	img             *DayNightAnim
	layer           int8
	Hitbox, Drawbox *Rectangle
	squareSize      int
	Name            string
}

//Copys the WObj
func (o *WObj) Copy() (o2 *WObj) {
	o2 = &WObj{o.img.Copy(), o.layer, o.Hitbox.Copy(), o.Drawbox.Copy(), o.squareSize, o.Name}
	return
}

//Updates the animation
func (o *WObj) Update(frame int) {
	o.img.Update(frame)
}

//Moves the WObj by delta
func (o *WObj) MoveBy(dx, dy float64) {
	if dx != 0 || dy != 0 {
		x := o.Hitbox.Min().X + dx
		y := o.Hitbox.Min().Y + dy
		o.SetTopLeft(x, y)
	}
}

//Sets the WObjs top left to x and y
func (o *WObj) SetTopLeft(x, y float64) {
	o.Hitbox.MoveTo(&Point{x, y})
	w, h := o.img.Size()
	W := float64(w) / float64(o.squareSize)
	H := float64(h) / float64(o.squareSize)
	o.Drawbox = GetRectangle(o.Hitbox.Min().X-(W-o.Hitbox.Bounds().X)/2, o.Hitbox.Min().Y-(H-o.Hitbox.Bounds().Y), 0, 0)
	o.Drawbox.SetBounds(&Point{W, H})
}

//Sets the WObjs bottom right to x and y
func (o *WObj) SetBottomRight(x, y float64) {
	w, h := o.Bounds()
	o.SetTopLeft(x-w, y-h)
}

//Sets the WObjs middle to x and y
func (o *WObj) SetMiddle(x, y float64) {
	w, h := o.Bounds()
	o.SetTopLeft(x-w/2, y-h/2)
}

//Returns the Position of the middle of the WObj
func (o *WObj) GetMiddle() (float64, float64, int8) {
	pnt := o.Hitbox.GetMiddle()
	return pnt.X, pnt.Y, o.layer
}
func (o *WObj) GetPos() (float64, float64, int8) {
	return o.GetMiddle()
}

//Returns the Position of the left top of the WObj
func (o *WObj) GetTopLeft() (float64, float64) {
	pnt := o.Hitbox.Min()
	return pnt.X, pnt.Y
}
func (o *WObj) GetBottomRight() (float64, float64) {
	pnt := o.Hitbox.Max()
	return pnt.X, pnt.Y
}

//Sets the layer the WObj is drawn to
func (o *WObj) SetLayer(l int8) {
	o.layer = l
}

//Returns the DrawBox of the WObj
func (o *WObj) GetDrawBox() *Rectangle {
	return o.Drawbox
}

//Returns the real Bounds of the WObj
func (o *WObj) Bounds() (float64, float64) {
	bnds := o.Hitbox.Bounds()
	return bnds.X, bnds.Y
}

//Draws the WObj to the screen
func (o *WObj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	y := (o.Drawbox.Min().Y - leftTopY) * sqSize
	x := (o.Drawbox.Min().X - leftTopX) * sqSize
	o.img.SetParams(x+xStart, y+yStart, float64(o.Drawbox.Bounds().X)*sqSize, float64(o.Drawbox.Bounds().Y)*sqSize)
	o.img.LightLevel = lv
	o.img.DrawAnim(screen)
	//o.frame ++
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
fps: 	        [0-NaN]
hitBoxWidth:	[1-NaN]
hitBoxHeight:	[1-NaN]
XPos:			[0-NaN]
YPos:			[0-NaN]
layer:			[-128-128]
**/
func GetWObjFromParams(img *ebiten.Image, p *Params, name string) (o *WObj) {
	layer := int8(p.Get("layer"))
	spW := int(p.Get("spriteWidth"))
	fps := p.Get("fps")
	hitBoxW := p.Get("hitBoxWidth")
	hitBoxH := p.Get("hitBoxHeight")
	XPos := p.Get("XPos")
	YPos := p.Get("YPos")
	sqSize := int(p.Get("squareSize"))
	if img != nil {
		return GetWObj(GetDayNightAnim(1, 1, 1, 1, spW, FPS/fps, img), hitBoxW, hitBoxH, XPos, YPos, sqSize, layer, name)
	}
	return GetWObj(nil, hitBoxW, hitBoxH, XPos, YPos, sqSize, layer, name)
}
func GetEmptyWObjFromPath(name, path string) (*WObj, error) {
	ps := &Params{}
	err := ps.LoadFromFile(path + ".txt")
	if err != nil {
		return nil, err
	}
	return GetWObjFromParams(nil, ps, name), nil
}

//provide at least on path
func GetWObjFromPathAndAnim(name, path string, anim *DayNightAnim) (*WObj, error) {
	p := &Params{}
	err := p.LoadFromFile(path + ".txt")
	if err != nil {
		return nil, err
	}
	layer := int8(p.Get("layer"))
	hitBoxW := p.Get("hitBoxWidth")
	hitBoxH := p.Get("hitBoxHeight")
	XPos := p.Get("XPos")
	YPos := p.Get("YPos")
	sqSize := int(p.Get("squareSize"))
	return GetWObj(anim, hitBoxW, hitBoxH, XPos, YPos, sqSize, layer, name), nil
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
	} else if len(path) == 2 {
		pathImg = path[0]
		pathParams = path[1]
	}

	img, err := LoadEbitenImg(pathImg + ".png")
	if err != nil {
		return nil, err
	}
	ps := &Params{}
	err = ps.LoadFromFile(pathParams + ".txt")
	if err != nil {
		return nil, err
	}
	return GetWObjFromParams(img, ps, name), nil
}

//Loads all WObjs from a directory
func LoadAllWObjs(folderPath string) (map[string]*WObj, error) {
	if folderPath[len(folderPath)-1:] != "/" {
		folderPath += "/"
	}
	files, err := OSReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var currentError error
	objs := make(map[string]*WObj)
	names := make([]string, 0)
	for _, f := range files {
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
func GetWObj(img *DayNightAnim, HitboxW, HitboxH, x, y float64, squareSize int, layer int8, name string) (o *WObj) {
	o = &WObj{img: img, layer: layer, squareSize: squareSize, Name: name}
	o.Hitbox = GetRectangle(x, y, x+HitboxW, y+HitboxH)
	if img != nil {
		o.img.Update(0)
		o.SetTopLeft(x, y)
	}
	return
}

func IsStringInList(s string, l []string) bool {
	for _, str := range l {
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
