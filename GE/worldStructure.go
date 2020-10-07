package GE

import (
	"github.com/hajimehoshi/ebiten"
)

//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles int) (p *WorldStructure) {
	p = &WorldStructure{X:X,Y:Y,W:W,H:H}
	p.SetDisplayWH(WTiles, HTiles)
	return
}


const COLLIDING_IDX = 1
type WorldStructure struct {
	Tiles  			[]*Tile
	BackStructObjs	[]*StructureObj
	FrontStructObjs	[]*StructureObj
	TmpObj			map[*StructureObj]int
	
	LightLevel uint8
	IdxMat, LightMat, CollisionMat *Matrix
	
	//SetMiddle/Move
	middleX, middleY int
	
	//GetFrame
	frame  *ImageObj
	
	//Set in GetWorldStructure
	X,Y,W,H float64
	drawer *ImageObj
	//Set by SetDisplayWH and 
	xTiles, yTiles int
	xStart, yStart, tileS float64
}




//Draws The World Ground Tiles and the Objects form the layer which is currently in the middle of the screen
func (p *WorldStructure) DrawBack(screen *ebiten.Image) {
	for x := 0; x < p.IdxMat.W(); x++ {
		for y := 0; y < p.IdxMat.H(); y++ {
			tile_idx := p.IdxMat.Get(x, y)
			p.drawer.X, p.drawer.Y = float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
			if int(tile_idx) >= 0 && int(tile_idx) < len(p.Tiles) {
				p.Tiles[tile_idx].Draw(screen, p.drawer, p.frame, uint8(p.LightMat.Get(x, y)))
			}
		}
	}
	for _,obj := range(p.BackStructObjs) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			pnt := obj.HitBox.Min()
			obj.DrawStructObj(screen, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart, uint8(p.LightMat.GetAbs(int(pnt.X), int(pnt.Y))))
		}
	}
}
func (p *WorldStructure) DrawFront(screen *ebiten.Image) {
	for _,obj := range(p.FrontStructObjs) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			pnt := obj.HitBox.Min()
			obj.DrawStructObj(screen, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart, uint8(p.LightMat.GetAbs(int(pnt.X), int(pnt.Y))))
		}
	}
	for obj,frm := range(p.TmpObj) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			pnt := obj.HitBox.Min()
			obj.DrawStructObj(screen, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart, uint8(p.LightMat.GetAbs(int(pnt.X), int(pnt.Y))))
		}
		if frm < 1 {
			delete(p.TmpObj, obj)
		}else{
			p.TmpObj[obj] = frm-1
		}
	}
}



//Updates the collision Matrix
func (p *WorldStructure) UpdateCollisionMat() {
	p.CollisionMat = GetMatrix(p.IdxMat.WAbs(),p.IdxMat.HAbs(),COLLIDING_IDX-1)
	for _,obj := range(append(p.FrontStructObjs, p.BackStructObjs...)) {
		obj.DrawCollisionMatrix(p.CollisionMat)
	}
	for obj,_ := range(p.TmpObj) {
		obj.DrawCollisionMatrix(p.CollisionMat)
	}
}
//Checks if point collides with the collision matrix
func (p *WorldStructure) Collides(x,y int) bool {
	if p.CollisionMat.Get(x,y) == COLLIDING_IDX {
		return true
	}
	return false
}