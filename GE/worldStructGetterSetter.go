package GE

import (
	"math"
)

func (p *WorldStructure) MoveSmooth(dx, dy float64, update, force bool) {
	p.MoveMiddleDelta(-dx,-dy)
	mX, mY := p.GetMiddleDelta()
	if math.Abs(mX) >= 1 {
		if mX < 0 {
			p.Move(1, 0, update, force)
			p.middleDx += 1
		}else{
			p.Move(-1, 0, update, force)
			p.middleDx -= 1
		}
		
	}
	if math.Abs(mY) >= 1 {
		if mY < 0 {
			p.Move(0, 1, update, force)
			p.middleDy += 1
		}else{
			p.Move(0, -1, update, force)
			p.middleDy -= 1
		}
	}
}
func (p *WorldStructure) MoveMiddleDelta(dx, dy float64) {
	odx, ody := p.GetMiddleDelta()
	p.SetMiddleDelta(odx+dx, ody+dy)
}
func (p *WorldStructure) SetMiddleDelta(dx, dy float64) {
	p.middleDx, p.middleDy = dx, dy
}
func (p *WorldStructure) GetMiddleDelta() (float64, float64) {
	return p.middleDx, p.middleDy
}
//sets the Middle of the View
func (p *WorldStructure) SetMiddle(xP, yP int, force bool) {
	if (xP != p.middleX || yP != p.middleY) || force {
		p.middleX = xP
		p.middleY = yP
		x,y := p.middleX-(p.xTilesS-1)/2, p.middleY-(p.yTilesS-1)/2
		p.TileMat.SetFocus(x,y, x+p.xTilesS, y+p.yTilesS)
		p.TileMat.CopyFocus(p.LIdxMat)
		p.TileMat.CopyFocus(p.LightMat)
		p.TileMat.CopyFocus(p.ObjMat)
	}
}
//moves the view by dx and dy
func (p *WorldStructure) Move(dx,dy int, update, force bool) {
	p.SetMiddle(p.middleX+dx, p.middleY+dy, force)
	if update && (dx != 0 || dy != 0) {
		p.UpdateObjDrawables()
	}
}
//Sets the number of tiles to be displayed in X and Y direction
func (p *WorldStructure) SetDisplayWH(x,y int) {
	p.xTilesS = x+2
	p.yTilesS = y+2
	p.tileS = p.W / float64(x)
	p.xStart = p.X-p.tileS
	p.yStart = p.Y - (float64(y)*p.tileS-p.H)/2 - p.tileS
	if p.W < p.H {
		p.tileS = p.H / float64(y)
		p.xStart = p.X - (float64(x)*p.tileS-p.W)/2 - p.tileS
		p.yStart = p.Y -p.tileS
	}
	p.drawer = &ImageObj{}
	p.drawer.W = p.tileS
	p.drawer.H = p.tileS
	mx, my := p.Middle()
	p.SetMiddle(mx,my, true)
}
func (p *WorldStructure) SmoothMiddle() (float64, float64) {
	return float64(p.middleX)-p.middleDx, float64(p.middleY)-p.middleDy
}
//returns the middle of the view
func (p *WorldStructure) Middle() (int, int) {
	return p.middleX, p.middleY
}
//Applies a frame
func (p *WorldStructure) GetFrame(thickness float64, alpha uint8) {
	p.frame = p.drawer.GetFrame(thickness, alpha)
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

func (p *WorldStructure) SetLightLevel(newL int16) {
	p.lightLevel = newL
	if newL > p.maxLight {
		p.lightLevel = p.maxLight
		p.deltaB *= -1
	}
	if newL < p.minLight {
		p.lightLevel = p.minLight
		p.deltaB *= -1
	}
}
//Updates the background lightlevel
func (p *WorldStructure) UpdateLightLevel(ticks float64) {
	p.currentD += p.deltaB*ticks
	if math.Abs(p.currentD) >= 1 {
		p.SetLightLevel(p.GetLightLevel()+int16(p.currentD))
		p.currentD = 0
	}
}
func (p *WorldStructure) GetLightLevel() int16 {
	return p.lightLevel
}

func (p *WorldStructure) SetLightStats(min, max int16, lightChange float64) {
	p.minLight = min
	p.maxLight = max
	p.lightLevel = max
	p.deltaB = -lightChange
	p.currentD = 0
}
