package GE

import (
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

/**
StructureObj is an Object that can be displayed in the world
It always has a hitbox, but it may not collide

The embedded type Animation does not necessarily need consist out of multiple frames
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
	o.Update(0)
	o.SetToXY(x, y)
	return
}
func (o *StructureObj) Clone() *StructureObj {
	return &StructureObj{o.Structure, o.HitBox.Clone(), o.DrawBox.Clone(), o.frame}
}
//Sets the top left corner of the hitbox to a coordinate on the map
func (o *StructureObj) SetToXY(x,y float64) {
	o.HitBox.MoveTo(&Point{x,y})
	w,h := o.DayNightAnim.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}

//Draws the objects hitbox if it can collide
func (o *StructureObj) DrawCollisionMatrix(mat *Matrix, value int16) {
	//fmt.Println(o.HitBox.Print())
	if !o.Collides {
		value = -value
	}
	mat.FillAbs(int(o.HitBox.Min().X), int(o.HitBox.Min().Y), int(o.HitBox.Max().X), int(o.HitBox.Max().Y), value)
	//fmt.Println(mat.Print())
}

//Draws the StructureObj
func (o *StructureObj) DrawStructObj(screen *ebiten.Image, leftTop *Point, sqSize, xStart, yStart float64, lightlevel int16) {
	o.Update(o.frame)
	relPx, relPy := float64(o.DrawBox.Min().X-leftTop.X), float64(o.DrawBox.Min().Y-leftTop.Y)
	o.SetParams(relPx*sqSize+xStart, relPy*sqSize+yStart, float64(o.DrawBox.Bounds().X)*sqSize, float64(o.DrawBox.Bounds().Y)*sqSize)
	o.LightLevel = lightlevel
	o.DrawAnim(screen)
	o.frame ++
}