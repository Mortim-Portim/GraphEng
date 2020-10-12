package GE

import (
	"math"
)

//sets the Middle of the View
func (p *WorldStructure) SetMiddle(xP, yP int) {
	p.middleX = xP
	p.middleY = yP
	x,y := p.middleX-(p.xTiles-1)/2, p.middleY-(p.yTiles-1)/2
	p.TileMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	p.LIdxMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	
	//fmt.Println(x,":",y," : ", x+p.xTiles, ":", y+p.yTiles)
	//fmt.Println("Tiles:\n", p.TileMat.Print())
	//fmt.Println("LIdx:\n", p.LIdxMat.Print())
	p.ObjMat.ResetFocus()
	//fmt.Println("ObjMat:\n", p.ObjMat.Print())
	p.ObjMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	//fmt.Println("ObjMat:\n", p.ObjMat.Print())
}
//moves the view by dx and dy
func (p *WorldStructure) Move(dx,dy int) {
	p.middleX += dx
	p.middleY += dy
	p.SetMiddle(p.middleX, p.middleY)
}
//Sets the number of tiles to be displayed in X and Y direction
func (p *WorldStructure) SetDisplayWH(x,y int) {
	p.xTiles = x
	p.yTiles = y
	p.xStart = p.X
	p.tileS = p.W / float64(x)
	p.yStart = p.Y - (float64(y)*p.tileS-p.H)/2
	if p.W < p.H {
		p.tileS = p.H / float64(y)
		p.xStart = p.X - (float64(x)*p.tileS-p.W)/2
		p.yStart = p.Y
	}
	p.drawer = &ImageObj{}
	p.drawer.W = p.tileS
	p.drawer.H = p.tileS
	p.CurrentLightMat = GetMatrix(x,y, 0)
	p.SetMiddle(p.Middle())
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
	p.DrawLights(true)
}
