package GE

import (
	"github.com/hajimehoshi/ebiten"
)

/**
StructureObj is an Object that can be displayed in the world
It always has a hitbox, but it may not collide

The Animations do not necessarily need to consist out of multiple frames
**/					
type StructureObj struct {
	*Structure
	HitBox, DrawBox *Rectangle
	frame int
}
//Returns a StructureObj
func GetStructureObj(stru *Structure, x, y float64) (o *StructureObj) {
	o = &StructureObj{Structure:stru}
	o.HitBox = GetRectangle(x,y, x+stru.HitboxW, y+stru.HitboxH)
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
	return &StructureObj{o.Structure, o.HitBox.Copy(), o.DrawBox.Copy(), o.frame}
}

//Sets the top left corner of the hitbox to a coordinate on the map
func (o *StructureObj) SetToXY(x,y float64) {
	o.HitBox.MoveTo(&Point{x,y})
	w,h := o.NUA.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}
func (o *StructureObj) GetPos() (float64, float64, int8) {
	pnt := o.HitBox.GetMiddle()
	return pnt.X+0.5,pnt.Y+0.5,o.layer
}
func (o *StructureObj) GetDrawBox() *Rectangle {
	return o.DrawBox
}
//Draws the objects hitbox if it can collide
func (o *StructureObj) DrawCollisionMatrix(mat *Matrix, value int64) {
	if !o.Collides {
		value = -value
	}
	mat.FillAbs(int(o.HitBox.Min().X), int(o.HitBox.Min().Y), int(o.HitBox.Max().X), int(o.HitBox.Max().Y), value)
}

func (o *StructureObj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	if !o.understandable || (o.understandable && !o.IsUnderstood) {
		o.drawImg(o.NUA, screen, lv, leftTopX, leftTopY, sqSize, xStart, yStart)
	}else{
		o.drawImg(o.UA, screen, lv, leftTopX, leftTopY, sqSize, xStart, yStart)
	}
}
func (o *StructureObj) drawImg(img *DayNightAnim, screen *ebiten.Image, lv int16, leftTopX, leftTopY, sqSize, xStart, yStart float64) {
	y := (o.DrawBox.Min().Y-leftTopY)*sqSize
	x := (o.DrawBox.Min().X-leftTopX)*sqSize
	img.SetParams(x+xStart,y+yStart, float64(o.DrawBox.Bounds().X)*sqSize, float64(o.DrawBox.Bounds().Y)*sqSize)
	img.LightLevel = lv
	img.DrawAnim(screen)
	o.frame ++
}
//!DEPRECATED!
//Draws the StructureObj
func (o *StructureObj) DrawStructObj(screen *ebiten.Image, leftTop *Point, sqSize, xStart, yStart float64, lightlevel int16) {
	if !o.understandable || (o.understandable && !o.IsUnderstood) {
		o.drawDNImg(o.NUA, screen, leftTop, sqSize, xStart, yStart, lightlevel)
	}else{
		o.drawDNImg(o.UA, screen, leftTop, sqSize, xStart, yStart, lightlevel)
	}
}
//!DEPRECATED!
func (o *StructureObj) drawDNImg(img *DayNightAnim, screen *ebiten.Image, leftTop *Point, sqSize, xStart, yStart float64, lightlevel int16) {
	img.Update(o.frame)
	relPx, relPy := float64(o.DrawBox.Min().X-leftTop.X), float64(o.DrawBox.Min().Y-leftTop.Y)
	img.SetParams(relPx*sqSize+xStart, relPy*sqSize+yStart, float64(o.DrawBox.Bounds().X)*sqSize, float64(o.DrawBox.Bounds().Y)*sqSize)
	img.LightLevel = lightlevel
	img.DrawAnim(screen)
	o.frame ++
}