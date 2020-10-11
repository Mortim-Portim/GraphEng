package GE

import (
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles int) (p *WorldStructure) {
	p = &WorldStructure{X:X,Y:Y,W:W,H:H}
	p.TileMat = GetMatrix(WTiles, HTiles, 0)
	p.LIdxMat = GetMatrix(WTiles, HTiles, 0)
	p.ObjMat  = GetMatrix(WTiles, HTiles, 0)
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
	LightLevel uint8
	//TileMat stores indexes of tiles, LightMat stores the lightlevel, ObjMat stores indexes of Objects
	TileMat, LIdxMat, CurrentLightMat, ObjMat *Matrix
	
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
				p.Tiles[tile_idx].Draw(screen, p.drawer, p.frame, p.CurrentLightMat.Get(x, y))
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
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, p.CurrentLightMat.GetAbs(int(pnt.X), int(pnt.Y)))
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
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, p.CurrentLightMat.GetAbs(int(pnt.X), int(pnt.Y)))
					drawnObjs = append(drawnObjs, int(idx))
				}
			}
		}
	}
}
func (p *WorldStructure) UpdateLIdxMat() {
	p.LIdxMat = GetMatrix(p.TileMat.WAbs(), p.TileMat.HAbs(), -1)
	for i,l := range(p.Lights) {
		p.LIdxMat.SetAbs(int(l.Location.X), int(l.Location.Y), int16(i))
		//fmt.Println("Setting Light to Point: ", l.Location.Print())
	}
	//fmt.Println(p.LIdxMat.Print())
}
func (p *WorldStructure) UpdateLights() {
	for _,l := range(p.Lights) {
		l.ApplyRaycasting(p.ObjMat, 1)
	}
}

const LIGHT_COMP_RADIUS = 20
func (p *WorldStructure) DrawLights() {
	ls := make([]*Light, 0)
	for x := -LIGHT_COMP_RADIUS; x < p.CurrentLightMat.W()+LIGHT_COMP_RADIUS; x++ {
		for y := -LIGHT_COMP_RADIUS; y < p.CurrentLightMat.H()+LIGHT_COMP_RADIUS; y++ {
			idx := int(p.LIdxMat.Get(x,y))
			if idx >= 0 && idx < len(p.Lights) {
				ls = append(ls, p.Lights[idx])
			}
		}
	}
	//fmt.Println(ls)
	pnt := p.TileMat.Focus().Min()
	//fmt.Println(pnt.Print())
	for x := 0; x < p.TileMat.W(); x++ {
		for y := 0; y < p.TileMat.H(); y++ {
			p.CurrentLightMat.SetAbs(x,y, p.GetLightValueForPoint(x+int(pnt.X), y+int(pnt.Y), ls, int16(p.LightLevel)))
		}
	}
}
func (p *WorldStructure) GetLightValueForPoint(x,y int, ls []*Light, standard int16) (v int16) {
	v = standard
	for _,l := range(ls) {
		lv := l.GetAtAbs(x,y)
		if lv >= 0 {
			v += lv
		}
	}
	return
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