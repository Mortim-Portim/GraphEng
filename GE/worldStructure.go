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


const NON_COLLIDING_IDX = -1
type WorldStructure struct {
	//Tiles and Structures should be the same on all devices
	Tiles  			[]*Tile
	Structures		[]*Structure
	
	//Represent Pointer to structures with a hitbox
	Objects			[]*StructureObj
	
	//Represent Pointer to light sources
	Lights			[]*Light
	
	//The standard light level
	lightLevel uint8
	//TileMat stores indexes of tiles, LightMat stores the lightlevel, ObjMat stores indexes of Objects
	TileMat, LightMat , ObjMat *Matrix
	
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
	for x := 0; x < p.TileMat.W(); x++ {
		for y := 0; y < p.TileMat.H(); y++ {
			tile_idx := p.TileMat.Get(x, y)
			p.drawer.X, p.drawer.Y = float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
			if int(tile_idx) >= 0 && int(tile_idx) < len(p.Tiles) {
				p.Tiles[tile_idx].Draw(screen, p.drawer, p.frame, p.LightMat.Get(x, y))
			}
		}
	}
	drawnObjs := make([]int,0)
	for x := 0; x < p.ObjMat.W(); x++ {
		for y := 0; y < p.ObjMat.H(); y++ {
			idx := p.ObjMat.Get(x, y)
			if idx >= 0 {
				obj := p.Objects[idx]
				if obj.Background && !containsI(drawnObjs, int(idx)){
					pnt := obj.HitBox.Min()
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, p.LightMat.GetAbs(int(pnt.X), int(pnt.Y)))
					drawnObjs = append(drawnObjs, int(idx))
				}
			}
		}
	}
}
func (p *WorldStructure) DrawFront(screen *ebiten.Image) {
	drawnObjs := make([]int,0)
	for x := 0; x < p.ObjMat.W(); x++ {
		for y := 0; y < p.ObjMat.H(); y++ {
			idx := p.ObjMat.Get(x, y)
			if idx >= 0 {
				obj := p.Objects[idx]
				if !obj.Background && !containsI(drawnObjs, int(idx)){
					pnt := obj.HitBox.Min()
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, p.LightMat.GetAbs(int(pnt.X), int(pnt.Y)))
					drawnObjs = append(drawnObjs, int(idx))
				}
			}
		}
	}
}
func (p *WorldStructure) UpdateLightMat(newLightLevel uint8) {
	if p.LightMat == nil {
		p.LightMat = GetMatrix(p.TileMat.WAbs(), p.TileMat.HAbs(), 0)
	}
	p.LightMat.AddToAllAbs(int16(newLightLevel-p.lightLevel))
	p.lightLevel = newLightLevel
}
func (p *WorldStructure) UpdateObjMat() {
	p.ObjMat = GetMatrix(p.TileMat.WAbs(),p.TileMat.HAbs(),NON_COLLIDING_IDX)
	for i,obj := range(p.Objects) {
		obj.DrawCollisionMatrix(p.ObjMat, int16(i))
	}
	p.TileMat.CopyFocus(p.ObjMat)
}
func (p *WorldStructure) Collides(x,y int) bool {
	if p.ObjMat.Get(x,y) == NON_COLLIDING_IDX {
		return false
	}
	return true
}