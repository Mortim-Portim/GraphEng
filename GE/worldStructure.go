package GE

import (
	"fmt"
	"image/color"
	"math"
	"github.com/hajimehoshi/ebiten"
	cmp "github.com/mortim-portim/GraphEng/Compression"
)

//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles, ScreenWT, ScreenHT int) (p *WorldStructure) {
	p = &WorldStructure{X: X, Y: Y, W: W, H: H, xTilesAbs: WTiles, yTilesAbs: HTiles}
	p.TileMat = 	GetMatrix(WTiles, HTiles, 0)
	p.LIdxMat = 	GetMatrix(WTiles, HTiles, -1)
	p.ObjMat = 		GetMatrix(WTiles, HTiles, 0)
	p.LightMat = 	GetMatrix(WTiles, HTiles, 0)
	p.RegionMat = 	GetMatrix(WTiles, HTiles, 0)
	p.Add_Drawables = GetDrawables()
	p.SO_Drawables = GetDrawables()
	p.SetDisplayWH(ScreenWT, ScreenHT)
	p.UpdateLightValue(p.Lights, true)
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
	Tiles      []*Tile
	Structures []*Structure

	//Represent Pointer to structures with a hitbox
	Objects []*StructureObj

	//Represent Pointer to light sources
	Lights []*Light

	//Represent basic Drawables, that can be drawn by their location on the map
	Add_Drawables *Drawables
	SO_Drawables  *Drawables

	//The standard light level
	lightLevel, minLight, maxLight int16
	deltaB, currentD               float64
	
	//The current region of the Player
	CurrentRegion int
	OnRegionChange func(oldR, newR int)
	
	//TileMat stores indexes of tiles, LightMat stores the lightlevel, ObjMat stores indexes of Objects, RegionMat stores the index of the Region
	TileMat, LIdxMat, LightMat, ObjMat, RegionMat *Matrix

	//SetMiddle/Move
	middleX, middleY int

	//GetFrame
	frame *ImageObj

	frameThickness float64
	frameAlpha     uint8
	frameScale     int

	//Set in GetWorldStructure
	X, Y, W, H float64
	drawer     *ImageObj
	//Set by SetDisplayWH and
	xTilesAbs, yTilesAbs, xTilesS, yTilesS, middleDx, middleDy int
	xStart, yStart, tileS                                      float64
}
func (p *WorldStructure) Size() (int, int) {
	return p.xTilesAbs, p.yTilesAbs
}
func (p *WorldStructure) ScaleTo(w, h int) {
	p.xTilesAbs = w
	p.yTilesAbs = h
	p.TileMat.ScaleTo(w,h, 0)
	p.LIdxMat.ScaleTo(w,h, -1)
	p.ObjMat.ScaleTo(w,h, 0)
	p.LightMat.ScaleTo(w,h, 0)
	p.RegionMat.ScaleTo(w,h, 0)
}
func (p *WorldStructure) Print() (out string) {
	out = fmt.Sprintf("Tiles: %v, Structures: %v, Objects: %v, Lights: %v, Add_Drawables: %v, SO_Drawables: %v\n",
		len(p.Tiles), len(p.Structures), len(p.Objects), len(p.Lights), p.Add_Drawables.Len(), p.SO_Drawables.Len())
	out += fmt.Sprintf("LightLevel: %v, minL: %v, maxL: %v, deltaB: %v, currentD: %v, middleX: %v, middleY: %v\n",
		p.lightLevel, p.minLight, p.maxLight, p.deltaB, p.currentD, p.middleX, p.middleY)
	out += fmt.Sprintf("X:%v, Y:%v, W:%v, H:%v, xTsAbs: %v, yTsAbs: %v, xTs: %v, yTs: %v, middleDx: %v, middleDy: %v, xStart: %v, yStart: %v, tileS: %v\n",
		p.X, p.Y, p.W, p.H, p.xTilesAbs, p.yTilesAbs, p.xTilesS, p.yTilesS, p.middleDx, p.middleDy, p.xStart, p.yStart, p.tileS)
	out += fmt.Sprintf("TilesMat:\n%s\nLIdxMat:\n%s\nLightMat:\n%s\nObjMat:\n%s",
		p.TileMat.Print(), p.LIdxMat.Print(), p.LightMat.Print(), p.ObjMat.Print())
	return
}

//Draws the tiles first and then SO_Drawables
func (p *WorldStructure) Draw(screen *ebiten.Image) {
	p.drawTiles(screen)
	p.DrawFrame(screen)
	lT := p.TileMat.Focus().Min()
	for _, dwa := range *p.SO_Drawables {
		xp, yp, _ := dwa.GetPos()
		lv := p.GetLightValueAtPoint(int(xp-lT.X), int(yp-lT.Y))
		dwa.Draw(screen, int16(lv), lT.X, lT.Y, p.xStart+float64(p.middleDx), p.yStart+float64(p.middleDy), p.tileS)
	}
}

//simply draws the current tiles to the screen
func (p *WorldStructure) drawTiles(screen *ebiten.Image) {
	for y := 0; y < p.TileMat.H(); y++ {
		for x := 0; x < p.TileMat.W(); x++ {
			tile_idx, err := p.TileMat.Get(x, y)
			if err == nil {
				if int(tile_idx) >= 0 && int(tile_idx) < len(p.Tiles) {
					idx := tile_idx
					p.drawer.X, p.drawer.Y = float64(x)*p.tileS+p.xStart+float64(p.middleDx), float64(y)*p.tileS+p.yStart+float64(p.middleDy)
					lv := p.GetLightValueAtPoint(x, y)
					p.Tiles[idx].Draw(screen, p.drawer, int16(lv))
				}
			}
		}
	}
}
func (p *WorldStructure) DrawFrame(screen *ebiten.Image) {
	if !p.HasFrame() {
		return
	}
	col := &color.RGBA{0, 0, 0, p.frameAlpha}
	lT := p.TileMat.Focus().Min()
	sX, sY := p.xStart+float64(p.middleDx), p.yStart+float64(p.middleDy)
	w, h := sX+float64(2+p.xTilesS)*p.tileS, sY+float64(2+p.yTilesS)*p.tileS
	for x := -(int(lT.X) % p.frameScale); x <= p.xTilesS; x += p.frameScale {
		X := sX + float64(x)*p.tileS
		l := GetLineOfPoints(X, sY, X, sY+h, p.frameThickness)
		l.Fill(screen, col)
	}
	for y := -(int(lT.Y) % p.frameScale); y <= p.yTilesS; y += p.frameScale {
		Y := sY + float64(y)*p.tileS
		l := GetLineOfPoints(sX, Y, sX+w, Y, p.frameThickness)
		l.Fill(screen, col)
	}
}

//ONLY use when adding or removing lights
func (p *WorldStructure) UpdateLIdxMat() {
	p.LIdxMat = GetMatrix(p.xTilesAbs, p.yTilesAbs, -1)
	for i, l := range p.Lights {
		p.LIdxMat.SetAbs(int(l.Loc().X), int(l.Loc().Y), int64(i))
	}
	p.TileMat.CopyFocus(p.LIdxMat)
}
func (p *WorldStructure) MakeLightMat() {
	p.LightMat = GetMatrix(p.xTilesAbs, p.yTilesAbs, 0)
	p.UpdateLightValue(p.Lights, true)
}
func (p *WorldStructure) AddLights(ls ...*Light) {
	p.Lights = append(p.Lights, ls...)
	p.UpdateLIdxMat()
	p.UpdateLightValue(ls, true)
}
func (p *WorldStructure) RemoveLight(idx int) {
	if idx >= 0 && idx < len(p.Lights) {
		p.Lights[idx].SetDark()
		p.UpdateAllLightsIfNecassary()
		p.Lights[idx] = p.Lights[len(p.Lights)-1]
		p.Lights = p.Lights[:len(p.Lights)-1]
		p.UpdateLIdxMat()
	}
}

//Updates all lights that somehow changed
func (p *WorldStructure) UpdateAllLightsIfNecassary() int {
	return p.UpdateLightValue(p.Lights, false)
}

//applies raycasting if necassary and updates the LightMat if necassary
func (p *WorldStructure) UpdateLightValue(ls []*Light, forceUpdate bool) (UpdatedLights int) {
	UpdatedLights = 0
	for _, l := range ls {
		l.ApplyRaycasting(p.ObjMat, 1, forceUpdate)
		if l.Changed() || forceUpdate {
			loc := l.Loc()
			r := int(math.Round(l.GetRadius()))
			p.drawLightsToMat(int(loc.X)-r, int(loc.Y)-r, r*2+1, r*2+1)
			l.SetChanged(false)
			UpdatedLights++
		}
	}
	return
}

const LIGHT_COMP_RADIUS = 50

func (p *WorldStructure) drawLightsToMat(xL, yL, w, h int) {
	ls := make([]*Light, 0)
	for x := xL - LIGHT_COMP_RADIUS; x < xL+w+LIGHT_COMP_RADIUS; x++ {
		for y := yL - LIGHT_COMP_RADIUS; y < yL+h+LIGHT_COMP_RADIUS; y++ {
			idx, err := p.LIdxMat.GetAbs(x, y)
			if int(idx) >= 0 && int(idx) < len(p.Lights) && err == nil {
				ls = append(ls, p.Lights[idx])
			}
		}
	}
	for x := xL; x < xL+w; x++ {
		for y := yL; y < yL+h; y++ {
			v := p.calcLightValueForPoint(x, y, ls)
			p.LightMat.SetAbs(x, y, v)
		}
	}
}

//Returns the sum of all the light values of all lights in ls at a relative point and the lightLevel
func (p *WorldStructure) GetLightValueAtPoint(x, y int) int64 {
	v, _ := p.LightMat.GetNearest(x, y)
	return v + int64(p.lightLevel)
}

//Calculates the sum of all the light values of all lights in ls at an absolute point
func (p *WorldStructure) calcLightValueForPoint(x, y int, ls []*Light) (v int64) {
	v = 0
	for _, l := range ls {
		lv, err := l.GetAtAbs(x, y)
		if err == nil {
			v += lv
		}
		if p.maxLight != 0 && v > int64(p.maxLight) {
			v = int64(p.maxLight)
			return
		}
	}
	return
}

//ONLY use when changing an objects location or hitbox
func (p *WorldStructure) UpdateObjMat() {
	p.ObjMat = GetMatrix(p.TileMat.WAbs(), p.TileMat.HAbs(), 0)
	for i, obj := range p.Objects {
		obj.DrawCollisionMatrix(p.ObjMat, int64(i+1))
	}
	p.TileMat.CopyFocus(p.ObjMat)
}

//ONLY use after moving the world to a different position
func (p *WorldStructure) UpdateObjDrawables() {
	*p.SO_Drawables = *p.Add_Drawables
	drawnObjs := make([]int, 0)
	minP := p.ObjMat.Focus().Min()
	maxP := p.ObjMat.Focus().Max()
	for y := int(minP.Y-OBJ_MAX_SIZE); y < int(maxP.Y+OBJ_MAX_SIZE); y++ {
		for x := int(minP.X-OBJ_MAX_SIZE); x < int(maxP.X+OBJ_MAX_SIZE); x++ {
			mvi, err := p.ObjMat.GetAbs(x, y)
			idx := int(math.Abs(float64(mvi)))
			if idx != 0 && err == nil {
				idx -= 1
				obj := p.Objects[idx]
				if !containsI(drawnObjs, int(idx)) {
					p.SO_Drawables = p.SO_Drawables.Add(obj)
					drawnObjs = append(drawnObjs, int(idx))
				}
			}
		}
	}
	p.SO_Drawables.Sort()
}
func (p *WorldStructure) AddDrawable(d Drawable) {
	p.Add_Drawables = p.Add_Drawables.Add(d)
}
func FloatPosToIntPos(fx, fy float64) (int, int) {
	return int(math.Round(fx-0.5)), int(math.Round(fy-0.5))
}
const OBJ_MAX_SIZE = 20
//Checks if an object obstructs the point
func (p *WorldStructure) Collides(x, y, w, h float64) bool {
	idxs := p.GetObjectsInField(int(x)-OBJ_MAX_SIZE, int(y)-OBJ_MAX_SIZE, int(w)+2*OBJ_MAX_SIZE, int(h)+2*OBJ_MAX_SIZE)
	return p.collidesWithObjs(x,y,w,h, idxs...)
}
func (p *WorldStructure) collidesWithObjs(x,y,w,h float64, idxs ...int) bool {
	r := GetRectangle(x,y,x+w,y+h)
	for _,idx := range(idxs) {
		obj := p.Objects[idx]
		if r.Overlaps(obj.Hitbox) {
			return true
		}
	}
	return !r.Overlaps(GetRectangle(0,0,float64(p.xTilesAbs),float64(p.yTilesAbs)))
}
func (p *WorldStructure) GetObjectsInField(X,Y,W,H int) (idxs []int) {
	for x := X; x < X+W; x++ {
		for y := Y; y < Y+H; y++ {
			v, err := p.ObjMat.GetAbs(x, y)
			if v > 0 && err == nil {
				idxs = append(idxs, int(v-1))
			}
		}
	}
	return
}
func (p *WorldStructure) ObjectsToBytes() (bs []byte) {
	bss := make([][]byte, 0)
	for _, obj := range p.Objects {
		bss = append(bss, cmp.CompressAll([][]byte{[]byte(obj.Name)}, cmp.Float64ToBytes(obj.Hitbox.Min().X), cmp.Float64ToBytes(obj.Hitbox.Min().Y)))
	}
	bs = cmp.CompressAll(bss)
	return
}
func (p *WorldStructure) BytesToObjects(bsss []byte) {
	bss := cmp.DecompressAll(bsss, []int{})
	p.Objects = make([]*StructureObj, 0)
	for _, bs := range bss {
		b := cmp.DecompressAll(bs, []int{8, 8})
		x := cmp.BytesToFloat64(b[0])
		y := cmp.BytesToFloat64(b[1])
		name := string(b[2])
		strct := p.GetNamedStructure(name)
		if strct != nil {
			obj := GetStructureObj(strct, x, y)
			p.Objects = append(p.Objects, obj)
		}
	}
	p.UpdateObjMat()
	p.UpdateObjDrawables()
}
func (p *WorldStructure) LightsToBytes() (bs []byte) {
	bss := make([][]byte, 0)
	for _, l := range p.Lights {
		bss = append(bss, l.ToBytes())
	}
	bs = cmp.CompressAll(bss)
	return
}
func (p *WorldStructure) BytesToLights(bs []byte) {
	bss := cmp.DecompressAll(bs, []int{})
	p.Lights = make([]*Light, len(bss))
	for i, b := range bss {
		p.Lights[i] = GetLightSourceFromBytes(b)
	}
	p.UpdateLIdxMat()
	p.MakeLightMat()
	/**
	for _,l := range(p.Lights) {
		fmt.Println(l.Matrix().Print())
	}
	fmt.Println(p.LightMat.Print())
	**/
}
func (p *WorldStructure) GetNamedStructure(name string) (s *Structure) {
	for _, st := range p.Structures {
		if st.Name == name {
			s = st
			break
		}
	}
	return
}
func (p *WorldStructure) GetTileOfCoords(x, y int) (xT, yT int) {
	x -= int(p.xStart)
	y -= int(p.yStart)
	xWithoutDx := float64(x - x%int(p.tileS))
	yWithoutDy := float64(y - y%int(p.tileS))
	tilesDX := xWithoutDx / p.tileS
	tilesDY := yWithoutDy / p.tileS
	loc := p.TileMat.Focus().Min()
	return int(loc.X + tilesDX), int(loc.Y + tilesDY)
}
func (p *WorldStructure) GetTileOfCoordsFP(x, y float64) (xT, yT float64) {
	x -= p.xStart
	y -= p.yStart
	tilesDX := x / p.tileS
	tilesDY := x / p.tileS
	loc := p.TileMat.Focus().Min()
	return loc.X + tilesDX, loc.Y + tilesDY
}
/**
//!DEPRECATED!
//ONLY use when moving the world before drawing tiles or objects
func (p *WorldStructure) DrawLights(update bool) {
	t_LightIdx_s := time.Now()
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
	fmt.Println(time.Now().Sub(t_LightIdx_s))
	pnt := p.TileMat.Focus().Min()
	for x := 0; x < p.TileMat.W(); x++ {
		for y := 0; y < p.TileMat.H(); y++ {
			p.CurrentLightMat.SetAbs(x,y, p.GetLightValueForPoint(x+int(pnt.X), y+int(pnt.Y), ls, int16(p.lightLevel)))
		}
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
**/
