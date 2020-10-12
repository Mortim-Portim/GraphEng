package GE

import (
	"github.com/hajimehoshi/ebiten"
	"math"
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

/**
Saved:

Objects

Lights

TileMat
middleX, middleY
minLight, maxLight
deltaB
**/
type WorldStructure struct {
	//Tiles and Structures should be the same on all devices
	Tiles  			[]*Tile
	Structures		[]*Structure
	
	//Represent Pointer to structures with a hitbox
	Objects			[]*StructureObj
	
	//Represent Pointer to light sources
	Lights			[]*Light
	
	//The standard light level
	lightLevel, minLight, maxLight int16
	deltaB, currentD float64
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
			tile_idx,err := p.TileMat.Get(x, y)
			if err == nil {
				p.drawer.X, p.drawer.Y = float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
				if int(tile_idx) >= 0 && int(tile_idx) < len(p.Tiles) {
					lv, _ := p.CurrentLightMat.Get(x, y)
					p.Tiles[tile_idx].Draw(screen, p.drawer, p.frame, lv)
				}
			}
		}
	}
	drawnObjs := make([]int,0)
	for x := 0; x < p.ObjMat.W(); x++ {
		for y := 0; y < p.ObjMat.H(); y++ {
			mvi,err := p.ObjMat.Get(x, y)
			idx := int(math.Abs(float64(mvi)))
			if idx != 0 && err == nil {
				idx -= 1
				obj := p.Objects[idx]
				if obj.Background && !containsI(drawnObjs, int(idx)){
					lv,_ := p.CurrentLightMat.GetAbs(x, y)
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, lv)
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
			mvi,err := p.ObjMat.Get(x, y)
			idx := int(math.Abs(float64(mvi)))
			if idx != 0 && err == nil {
				idx -= 1
				obj := p.Objects[idx]
				if !obj.Background && !containsI(drawnObjs, int(idx)){
					lv,_ := p.CurrentLightMat.GetAbs(x,y)
					obj.DrawStructObj(screen, p.ObjMat.Focus().Min(), p.tileS, p.xStart, p.yStart, lv)
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
	}
}
func (p *WorldStructure) UpdateLights(ls []*Light) {
	for _,l := range(ls) {
		l.ApplyRaycasting(p.ObjMat, 1)
	}
}

const LIGHT_COMP_RADIUS = 20
func (p *WorldStructure) DrawLights(update bool) {
	ls := make([]*Light, 0)
	for x := -LIGHT_COMP_RADIUS; x < p.CurrentLightMat.W()+LIGHT_COMP_RADIUS; x++ {
		for y := -LIGHT_COMP_RADIUS; y < p.CurrentLightMat.H()+LIGHT_COMP_RADIUS; y++ {
			idx,err := p.LIdxMat.Get(x,y)
			if int(idx) >= 0 && int(idx) < len(p.Lights) && err == nil {
				ls = append(ls, p.Lights[idx])
			}
		}
	}
	if update {
		p.UpdateLights(ls)
	}
	pnt := p.TileMat.Focus().Min()
	for x := 0; x < p.TileMat.W(); x++ {
		for y := 0; y < p.TileMat.H(); y++ {
			p.CurrentLightMat.SetAbs(x,y, p.GetLightValueForPoint(x+int(pnt.X), y+int(pnt.Y), ls, int16(p.lightLevel)))
		}
	}
}
func (p *WorldStructure) GetLightValueForPoint(x,y int, ls []*Light, standard int16) (v int16) {
	v = standard
	for _,l := range(ls) {
		lv,err := l.GetAtAbs(x,y)
		if err == nil {
			v += lv
		}
		if v > p.maxLight {
			v = p.maxLight
			return
		}
	}
	return
}
func (p *WorldStructure) UpdateObjMat() {
	p.ObjMat = GetMatrix(p.TileMat.WAbs(),p.TileMat.HAbs(),0)
	for i,obj := range(p.Objects) {
		obj.DrawCollisionMatrix(p.ObjMat, int16(i+1))
	}
	p.TileMat.CopyFocus(p.ObjMat)
}
func (p *WorldStructure) Collides(x,y int) bool {
	v, err := p.ObjMat.Get(x,y)
	if v <= 0 && err == nil {
		return false
	}
	return true
}
func (p *WorldStructure) BytesToObjects(bsss []byte) {
	bss := DecompressAll(bsss, []int{})
	p.Objects = make([]*StructureObj, 0)
	for _,bs := range(bss) {
		b := DecompressAll(bs, []int{8,8})
		x := BytesToFloat64(b[0])
		y := BytesToFloat64(b[1])
		name := string(b[2])
		obj := GetStructureObj(p.GetNamedStructure(name), x, y)
		p.Objects = append(p.Objects, obj)
	}
}
func (p *WorldStructure) ObjectsToBytes() (bs []byte) {
	bss := make([][]byte, 0)
	for _,obj := range(p.Objects) {
		bss = append(bss, CompressAll([][]byte{[]byte(obj.Name)}, Float64ToBytes(obj.HitBox.Min().X), Float64ToBytes(obj.HitBox.Min().Y)))
	}
	bs = CompressAll(bss)
	return
}

func (p *WorldStructure) GetNamedStructure(name string) (s *Structure) {
	for _,st := range(p.Structures) {
		if st.Name == name {
			s = st
			break
		}
	}
	return
}