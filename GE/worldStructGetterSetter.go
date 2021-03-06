package GE

import (
	"fmt"
	"math"
	"time"
)

func (p *WorldStructure) MoveSmooth(dx, dy int, update, force bool) {
	p.MoveMiddleDelta(-dx, -dy)
	mX, mY := p.GetMiddleDelta()
	ts := int(p.tileS)
	if IntAbs(mX) >= ts {
		if mX < 0 {
			p.Move(1, 0, update, force)
			p.middleDx += ts
		} else {
			p.Move(-1, 0, update, force)
			p.middleDx -= ts
		}

	}
	if IntAbs(mY) >= ts {
		if mY < 0 {
			p.Move(0, 1, update, force)
			p.middleDy += ts
		} else {
			p.Move(0, -1, update, force)
			p.middleDy -= ts
		}
	}
}
func (p *WorldStructure) SetMiddleSmooth(x, y float64) {
	xR := math.Remainder(x, 1.0)
	yR := math.Remainder(y, 1.0)
	p.SetMiddleDelta(-int(p.tileS*xR), -int(p.tileS*yR))
	p.SetMiddle(int(x-xR), int(y-yR), false)
}
func (p *WorldStructure) MoveMiddleDelta(dx, dy int) {
	odx, ody := p.GetMiddleDelta()
	p.SetMiddleDelta(odx+dx, ody+dy)
}
func (p *WorldStructure) SetMiddleDelta(dx, dy int) {
	p.middleDx, p.middleDy = dx, dy
}
func (p *WorldStructure) GetMiddleDelta() (int, int) {
	return p.middleDx, p.middleDy
}

//sets the Middle of the View
func (p *WorldStructure) SetMiddle(xP, yP int, force bool) {
	if (xP != p.middleX || yP != p.middleY) || force {
		p.middleX = xP
		p.middleY = yP
		x, y := p.middleX-(p.xTilesS-1)/2, p.middleY-(p.yTilesS-1)/2
		p.TileMat.SetFocus(x, y, x+p.xTilesS, y+p.yTilesS)
		p.TileMat.CopyFocus(p.LIdxMat)
		p.TileMat.CopyFocus(p.LightMat)
		p.TileMat.CopyFocus(p.ObjMat)
		p.TileMat.CopyFocus(p.RegionMat)
		r, err := p.RegionMat.GetAbs(xP, yP)
		if err != nil && int(r) != p.CurrentRegion {
			if p.OnRegionChange != nil {
				p.OnRegionChange(p.CurrentRegion, int(r))
			}
			p.CurrentRegion = int(r)
		}
	}
}

//moves the view by dx and dy
func (p *WorldStructure) Move(dx, dy int, update, force bool) {
	p.SetMiddle(p.middleX+dx, p.middleY+dy, force)
	if update && (dx != 0 || dy != 0) {
		p.UpdateObjDrawables()
	}
}

//Sets the number of tiles to be displayed in X and Y direction
func (p *WorldStructure) SetDisplayWH(x, y int) {
	p.xTilesS = x + 2
	p.yTilesS = y + 2
	p.tileS = p.W / float64(x)
	p.xStart = p.X - p.tileS
	p.yStart = p.Y - (float64(y)*p.tileS-p.H)/2 - p.tileS
	if p.W < p.H {
		p.tileS = p.H / float64(y)
		p.xStart = p.X - (float64(x)*p.tileS-p.W)/2 - p.tileS
		p.yStart = p.Y - p.tileS
	}
	p.drawer = &ImageObj{}
	p.drawer.W = p.tileS
	p.drawer.H = p.tileS
	mx, my := p.Middle()
	p.SetMiddle(mx, my, true)
}
func (p *WorldStructure) SmoothMiddle() (float64, float64) {
	return float64(p.middleX) - float64(p.middleDx)/p.tileS, float64(p.middleY) - float64(p.middleDy)/p.tileS
}

//returns the middle of the view
func (p *WorldStructure) Middle() (int, int) {
	return p.middleX, p.middleY
}

//Applies a frame
func (p *WorldStructure) GetFrame(thickness float64, alpha uint8, scale int) {
	p.frameThickness = thickness
	p.frameAlpha = alpha
	p.frameScale = scale
}
func (p *WorldStructure) HasFrame() bool {
	return p.frameThickness != 0 && p.frameAlpha != 0 && p.frameScale != 0
}

//Returns the width and height of the tiles
func (p *WorldStructure) GetTileS() float64 {
	return p.tileS
}

//Returns the top left corner of the WorldStruct on the screen
func (p *WorldStructure) GetTopLeft() (float64, float64) {
	return p.xStart, p.yStart
}

//Adds a tile to the index list of the worlds tiles
func (p *WorldStructure) AddTile(img *Tile) {
	if p.Tiles == nil {
		p.Tiles = make([]*Tile, 0)
	}
	p.Tiles = append(p.Tiles, img)
}

//Adds a StructureObj to the list of the worlds Structures
func (p *WorldStructure) AddStruct(obj ...*Structure) {
	if p.Structures == nil {
		p.Structures = make([]*Structure, 0)
	}
	p.Structures = append(p.Structures, obj...)
}

//Adds a StructureObj to the list of the worlds Objs
func (p *WorldStructure) AddStructObj(obj ...*StructureObj) {
	if p.Objects == nil {
		p.Objects = make([]*StructureObj, 0)
	}
	p.Objects = append(p.Objects, obj...)
}
func GetStandardTimeToLvFunc(minL, maxL float64) func(secs int) (lv int16) {
	return func(x int) int16 {
		return int16(minL + (maxL-minL)*0.5*(1+math.Sin((math.Pi/43200)*float64(x)-0.5*math.Pi)))
	}
}
func (p *WorldStructure) SetLightLevel(newL int16) {
	p.lightLevel = newL
}
func (p *WorldStructure) GetLightLevel() int16 {
	return p.lightLevel
}

//Updates the background lightlevel
func (p *WorldStructure) UpdateTime(t time.Duration) {
	newT := p.CurrentTime.Add(t)
	p.CurrentTime = &newT
	secs := p.TimeToSec()
	if p.TimeToLV != nil {
		p.SetLightLevel(p.TimeToLV(secs))
	} else {
		fmt.Println("NO Lightlevel function !!!!!!!!!!!!!!!!!!!!!!!")
	}
}
func (p *WorldStructure) SetLightStats(minLightLevel, maxLightLevel int16) {
	p.TimeToLV = GetStandardTimeToLvFunc(float64(minLightLevel), float64(maxLightLevel))
	p.maxLightLevel = maxLightLevel
	p.minLightLevel = minLightLevel
	p.UpdateTime(0)
}
func (p *WorldStructure) TimeHM() string {
	return fmt.Sprintf("%v:%v", p.CurrentTime.Hour(), p.CurrentTime.Minute())
}
func (p *WorldStructure) TimeToSec() int {
	return HMStoS(p.CurrentTime.Clock())
}
func (p *WorldStructure) GetTimeRel() float64 {
	return float64(p.TimeToSec()) / 43200.0
}

func HMStoS(h, m, s int) int {
	return h*60*60 + m*60 + s
}
