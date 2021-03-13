package GE

import (
	"math"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

/**
StructureObj is an Object that can be displayed in the world
It always has a hitbox, but it may not collide

The Animations do not necessarily need to consist out of multiple frames
**/
type StructureObj struct {
	*Structure
	Hitbox, Drawbox *Rectangle
	frame           int
}

//Returns a StructureObj
func GetStructureObj(stru *Structure, x, y float64) (o *StructureObj) {
	o = &StructureObj{Structure: stru}
	o.Hitbox = GetRectangle(x, y, x+stru.H_W, y+stru.H_H)
	if o.NUA != nil {
		o.NUA.Update(0)
	}
	if o.UA != nil {
		o.UA.Update(0)
	}
	o.SetToXY(x, y)
	return
}
func (o *StructureObj) Clone() *StructureObj {
	return &StructureObj{o.Structure, o.Hitbox.Copy(), o.Drawbox.Copy(), o.frame}
}

//Sets the top left corner of the hitbox to a coordinate on the map
func (o *StructureObj) SetToXY(x, y float64) {
	o.Hitbox.MoveTo(&Point{x, y})
	w, h := o.NUA.Size()
	W := float64(w) / float64(o.squareSize)
	H := float64(h) / float64(o.squareSize)
	o.Drawbox = GetRectangle(o.Hitbox.Min().X-(W-o.Hitbox.Bounds().X)/2, o.Hitbox.Min().Y-(H-o.Hitbox.Bounds().Y), 0, 0)
	o.Drawbox.SetBounds(&Point{W, H})
}
func (o *StructureObj) GetPos() (float64, float64, int8) {
	pnt := o.Hitbox.GetMiddle()
	return pnt.X, pnt.Y, o.layer
}
func (o *StructureObj) GetDrawBox() *Rectangle {
	return o.Drawbox
}

//Draws the objects hitbox if it can collide
func (o *StructureObj) DrawCollisionMatrix(mat *Matrix, value int64) {
	if !o.Collides {
		value = -value
	}
	pnt := o.Hitbox.GetMiddle()
	mat.SetAbs(int(math.Round(pnt.X)), int(math.Round(pnt.Y)), value)
}

func (o *StructureObj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	if !o.understandable || (o.understandable && !o.IsUnderstood) {
		o.drawImg(o.NUA, screen, lv, leftTopX, leftTopY, sqSize, xStart, yStart)
	} else {
		o.drawImg(o.UA, screen, lv, leftTopX, leftTopY, sqSize, xStart, yStart)
	}
}
func (o *StructureObj) drawImg(img *DayNightAnim, screen *ebiten.Image, lv int16, leftTopX, leftTopY, sqSize, xStart, yStart float64) {
	img.Update(o.frame)
	y := (o.Drawbox.Min().Y - leftTopY) * sqSize
	x := (o.Drawbox.Min().X - leftTopX) * sqSize
	img.SetParams(x+xStart, y+yStart, float64(o.Drawbox.Bounds().X)*sqSize, float64(o.Drawbox.Bounds().Y)*sqSize)
	img.LightLevel = lv
	img.DrawAnim(screen)
	o.frame++
}

//!DEPRECATED!
//Draws the StructureObj
func (o *StructureObj) DrawStructObj(screen *ebiten.Image, leftTop *Point, sqSize, xStart, yStart float64, lightlevel int16) {
	if !o.understandable || (o.understandable && !o.IsUnderstood) {
		o.drawDNImg(o.NUA, screen, leftTop, sqSize, xStart, yStart, lightlevel)
	} else {
		o.drawDNImg(o.UA, screen, leftTop, sqSize, xStart, yStart, lightlevel)
	}
}

//!DEPRECATED!
func (o *StructureObj) drawDNImg(img *DayNightAnim, screen *ebiten.Image, leftTop *Point, sqSize, xStart, yStart float64, lightlevel int16) {
	img.Update(o.frame)
	relPx, relPy := float64(o.Drawbox.Min().X-leftTop.X), float64(o.Drawbox.Min().Y-leftTop.Y)
	img.SetParams(relPx*sqSize+xStart, relPy*sqSize+yStart, float64(o.Drawbox.Bounds().X)*sqSize, float64(o.Drawbox.Bounds().Y)*sqSize)
	img.LightLevel = lightlevel
	img.DrawAnim(screen)
	o.frame++
}
