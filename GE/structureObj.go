package GE

import (
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
)

/**
StructureObj is an Object that can be displayed in the world
It always has a hitbox, but it may not collide

The embedded type Animation does not necessarily need consist out of multiple frames
**/
					
type StructureObj struct {
	Animation
	frame, squareSize int
	HitBox, DrawBox *Rectangle
	Collides bool
	Name string
}
/**
Params.txt:
spriteWidth: 	[1-NaN]
updatePeriod: 	[0-NaN]
hitBoxWidth:	[1-NaN]
hitBoxHeight:	[1-NaN]
Collides:		[true/false]
squareSize:		[1-NaN]
**/
func GetStructObjFromParams(img *ebiten.Image, p *Params) (s *StructureObj) {
	anim := GetAnimation(1,1,1,1, int(p.Get("spriteWidth")), int(p.Get("updatePeriod")), img)
	hitBox := GetRectangle(0,0,p.Get("hitBoxWidth")-1,p.Get("hitBoxHeight")-1)
	collides := false
	if p.GetS("Collides") != "false" {
		collides = true
	}
	s = GetStructureObj(anim, hitBox, int(p.Get("squareSize")), collides)
	return
}

//Returns a StructureObj
func GetStructureObj(anim *Animation, HitBox *Rectangle, squareSize int, Collides bool) (o *StructureObj) {
	o = &StructureObj{Animation:*anim, frame:0, HitBox:HitBox, squareSize:squareSize, Collides:Collides}
	o.Update(0)
	pnt := HitBox.Min()
	o.SetToXY(int(pnt.X), int(pnt.Y))
	return
}
func (o *StructureObj) Clone() *StructureObj {
	return &StructureObj{o.Animation, o.frame, o.squareSize, o.HitBox.Clone(), o.DrawBox.Clone(), o.Collides, o.Name}
}
//Sets the top left corner of the hitbox to a coordinate on the map
func (o *StructureObj) SetToXY(x,y int) {
	o.HitBox.MoveTo(&Point{float64(x),float64(y)})
	w,h := o.Img.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}
//Draws the objects hitbox if it can collide
func (o *StructureObj) DrawCollisionMatrix(mat *Matrix) {
	if o.Collides {
		mat.Fill(int(o.HitBox.Min().X), int(o.HitBox.Min().Y), int(o.HitBox.Max().X), int(o.HitBox.Max().Y), COLLIDING_IDX)
	}
}
//Draws the StructureObj
func (o *StructureObj) Draw(screen *ebiten.Image, myLayer, drawLayer int, leftTop *Point, sqSize, xStart, yStart float64) {
	o.Update(o.frame)
	
	relPx, relPy := float64(o.DrawBox.Min().X-leftTop.X), float64(o.DrawBox.Min().Y-leftTop.Y)
	o.X = relPx*sqSize+xStart
	o.Y = relPy*sqSize+yStart
	o.W = float64(o.DrawBox.Bounds().X)*sqSize
	o.H = float64(o.DrawBox.Bounds().Y)*sqSize
	
	if myLayer == drawLayer {
		o.DrawImageObj(screen)
	} else if myLayer < drawLayer {
		dif := 1.0 / float64(drawLayer-myLayer+1)
		o.DrawImageObjAlpha(screen, dif)
	} else {
		//box := float64(1 + myLayer - drawLayer)
		//sq := box*2 + 1
		//o.DrawImageBlured(screen, int(box), 1.0/((sq*sq)*0.35))
		o.DrawImageObj(screen)
	}
	o.frame ++
}

/**
Reads a slice of StructureObjs from a folder like this:
folder
----> obj1.png
----> obj1.txt
----> obj2.png
----> obj2.txt
**/
func ReadStructureObj(folderPath string) ([]*StructureObj, error) {
	files, err1 := ioutil.ReadDir(folderPath)
    if err1 != nil {
    	return nil, err1
    }
	ts := make([]*StructureObj, 0)
	names := make([]string, 0)
    for _, f := range files {
            name := f.Name()[:len(f.Name())-4]
	        if !containsS(names, name) {
				names = append(names, name)
				img, err := LoadEbitenImg(folderPath+name+".png")
				if err != nil {
					return nil, err
				}
				p := &Params{}
				err2 := p.LoadFromFile(folderPath+name+".txt")
				if err2 != nil {
					return nil, err2
				}
				obj := GetStructObjFromParams(img, p)
				obj.Name = name
				ts = append(ts, obj)
            }
    }
    return ts, nil
}