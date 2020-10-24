package GE

import (
	cmp "marvin/GraphEng/Compression"
	"github.com/hajimehoshi/ebiten"
	"math"
)

//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles int) (p *WorldStructure) {
	p = &WorldStructure{X:X,Y:Y,W:W,H:H}
	p.TileMat = GetMatrix(WTiles, HTiles, 0)
	p.LIdxMat = GetMatrix(WTiles, HTiles, -1)
	p.ObjMat  = GetMatrix(WTiles, HTiles, 0)
	p.Add_Drawables = GetDrawables()
	p.SO_Drawables = GetDrawables()
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
	
	//Represent basic Drawables, that can be drawn by their location on the map
	Add_Drawables	*Drawables
	SO_Drawables	*Drawables
	
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

//Draws the tiles first and then SO_Drawables
func (p *WorldStructure) Draw(screen *ebiten.Image) {
	p.drawTiles(screen)
	
	for _,dwa := range(*p.SO_Drawables) {
		x := dwa.GetX(); y := dwa.GetY()
		lv, err := p.CurrentLightMat.GetNearest(int(x), int(y-0.5))
		if err != nil {ShitImDying(err)}
		X := (x)*p.tileS + p.xStart
		Y := (y)*p.tileS + p.yStart
		dwa.Draw(screen, X, Y, lv, p.tileS)
		//fmt.Printf("Drawing %v at x:%v, y:%v, X:%v, Y:%v, lv:%v\n", dwa, x, y, X, Y, lv)
	}
}

//!DEPRECATED!
//Draw the top first
//Draws The World Ground Tiles and the Objects form the layer which is currently in the middle of the screen
func (p *WorldStructure) DrawBack(screen *ebiten.Image) {
	p.drawTiles(screen)
	drawnObjs := make([]int,0)
	for y := 0; y < p.ObjMat.H(); y++ {
		for x := 0; x < p.ObjMat.W(); x++ {
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
//!DEPRECATED!
func (p *WorldStructure) DrawFront(screen *ebiten.Image) {
	drawnObjs := make([]int,0)
	for y := 0; y < p.ObjMat.H(); y++ {
		for x := 0; x < p.ObjMat.W(); x++ {
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

func (p *WorldStructure) drawTiles(screen *ebiten.Image) {
	for y := 0; y < p.TileMat.H(); y++ {
		for x := 0; x < p.TileMat.W(); x++ {
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
}

//ONLY use when adding or removing lights
func (p *WorldStructure) UpdateLIdxMat() {
	p.LIdxMat = GetMatrix(p.TileMat.WAbs(), p.TileMat.HAbs(), -1)
	for i,l := range(p.Lights) {
		p.LIdxMat.SetAbs(int(l.Location.X), int(l.Location.Y), int16(i))
	}
	p.TileMat.CopyFocus(p.LIdxMat)
}
//ONLY use when changing objects or lights in order to reaplly raycasting
func (p *WorldStructure) UpdateLights(ls []*Light) {
	for _,l := range(ls) {
		l.ApplyRaycasting(p.ObjMat, 1)
	}
}

//ONLY use when moving the world before drawing tiles or objects
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
//ONLY use when changing an objects location or hitbox
func (p *WorldStructure) UpdateObjMat() {
	p.ObjMat = GetMatrix(p.TileMat.WAbs(),p.TileMat.HAbs(),0)
	for i,obj := range(p.Objects) {
		obj.DrawCollisionMatrix(p.ObjMat, int16(i+1))
	}
	p.TileMat.CopyFocus(p.ObjMat)
}
//ONLY use after moving the world to a different position
func (p *WorldStructure) UpdateObjDrawables() {
	p.SO_Drawables = p.Add_Drawables
	drawnObjs := make([]int,0)
	for y := 0; y < p.ObjMat.H(); y++ {
		for x := 0; x < p.ObjMat.W(); x++ {
			mvi,err := p.ObjMat.Get(x, y)
			idx := int(math.Abs(float64(mvi)))
			if idx != 0 && err == nil {
				idx -= 1
				obj := p.Objects[idx]
				if !containsI(drawnObjs, int(idx)){
					p.SO_Drawables = p.SO_Drawables.AddStructureObjOfWorld(obj, p)
					drawnObjs = append(drawnObjs, int(idx))
				}
			}
		}
	}
	p.SO_Drawables.Sort()
}
//Checks if an object obstructs the point
func (p *WorldStructure) Collides(x,y int) bool {
	v, err := p.ObjMat.GetAbs(x,y)
	if v <= 0 && err == nil {
		return false
	}
	return true
}

func (p *WorldStructure) ObjectsToBytes() (bs []byte) {
	bss := make([][]byte, 0)
	for _,obj := range(p.Objects) {
		bss = append(bss, cmp.CompressAll([][]byte{[]byte(obj.Name)}, cmp.Float64ToBytes(obj.HitBox.Min().X), cmp.Float64ToBytes(obj.HitBox.Min().Y)))
	}
	bs = cmp.CompressAll(bss)
	return
}
func (p *WorldStructure) BytesToObjects(bsss []byte) {
	bss := cmp.DecompressAll(bsss, []int{})
	p.Objects = make([]*StructureObj, 0)
	for _,bs := range(bss) {
		b := cmp.DecompressAll(bs, []int{8,8})
		x := cmp.BytesToFloat64(b[0])
		y := cmp.BytesToFloat64(b[1])
		name := string(b[2])
		strct := p.GetNamedStructure(name)
		if strct != nil {
			obj := GetStructureObj(strct, x, y)
			p.Objects = append(p.Objects, obj)
		}
	}
}
func (p *WorldStructure) LightsToBytes() (bs []byte) {
	bss := make([][]byte, 0)
	for _,l := range(p.Lights) {
		bss = append(bss, l.ToBytes())
	}
	bs = cmp.CompressAll(bss)
	return
}
func (p *WorldStructure) BytesToLights(bs []byte) {
	bss := cmp.DecompressAll(bs, []int{})
	p.Lights = make([]*Light, len(bss))
	for i,b := range(bss) {
		p.Lights[i] = GetLightSourceFromBytes(b)
	}
	p.UpdateLIdxMat()
	p.UpdateLights(p.Lights)
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