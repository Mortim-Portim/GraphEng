package GE

import (
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
)
//TODO Set world xTiles and yTiles


//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles int16) (p *WorldStructure) {
	p = &WorldStructure{xTiles:WTiles, yTiles:HTiles}
	p.xStart = X
	p.tileS = W / float64(WTiles)
	p.yStart = Y - (float64(HTiles)*p.tileS-H)/2
	if W < H {
		p.tileS = H / float64(HTiles)
		p.xStart = X - (float64(WTiles)*p.tileS-W)/2
		p.yStart = Y
	}
	
	p.drawer = &ImageObj{}
	p.drawer.W = p.tileS
	p.drawer.H = p.tileS
	return
}


const COLLIDING_IDX = 1
type WorldStructure struct {
	Tiles  			[]*Tile
	StructureObjs	[]*StructureObj
	TmpObj			map[*StructureObj]int
	
	IdxMat, LayerMat, CollisionMat *Matrix
	xTiles, yTiles, middleX, middleY int16
	
	frame  *ImageObj
	drawer *ImageObj

	xStart, yStart, tileS float64
}

//Updates the collision Matrix
func (p *WorldStructure) UpdateCollisionMat() {
	p.CollisionMat = GetMatrix(p.IdxMat.WAbs(),p.IdxMat.HAbs(),COLLIDING_IDX-1)
	for _,obj := range(p.StructureObjs) {
		obj.DrawCollisionMatrix(p.CollisionMat)
	}
	for obj,_ := range(p.TmpObj) {
		obj.DrawCollisionMatrix(p.CollisionMat)
	}
}
//Checks if point collides with the collision matrix
func (p *WorldStructure) Collides(x,y int16) bool {
	if p.CollisionMat.Get(x,y) == COLLIDING_IDX {
		return true
	}
	return false
}


//sets the Middle of the View
func (p *WorldStructure) SetMiddle(pnt *Point) {
	p.middleX = int16(pnt.X)
	p.middleY = int16(pnt.Y)
	x,y := p.middleX-(p.xTiles-1)/2, p.middleY-(p.yTiles-1)/2
	p.IdxMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	p.LayerMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
}
//moves the view by dx and dy
func (p *WorldStructure) Move(dx,dy int16) {
	p.middleX += dx
	p.middleY += dy
	x,y := p.middleX-(p.xTiles-1)/2, p.middleY-(p.yTiles-1)/2
	p.IdxMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	p.LayerMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
}
//Draws The World Ground Tiles and the Objects form the layer which is currently in the middle of the screen
func (p *WorldStructure) Draw(screen *ebiten.Image) {
	Zlayer := p.LayerMat.GetAbs(p.Middle())
	for x := int16(0); x < p.IdxMat.W(); x++ {
		for y := int16(0); y < p.IdxMat.H(); y++ {
			tile_idx := p.IdxMat.Get(x, y)
			xp, yp := float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
			if tile_idx >= 0 && tile_idx < len(p.Tiles) {
				p.Tiles[tile_idx].Draw(screen, xp, yp, p.tileS, p.tileS, p.LayerMat.Get(x, y), Zlayer, p.drawer, p.frame)
			}
		}
	}
	for _,obj := range(p.StructureObjs) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			objPnt := obj.HitBox.Min()
			obj.Draw(screen, p.LayerMat.GetAbs(int16(objPnt.X), int16(objPnt.Y)), Zlayer, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart)
		}
	}
	for obj,frm := range(p.TmpObj) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			objPnt := obj.HitBox.Min()
			obj.Draw(screen, p.LayerMat.GetAbs(int16(objPnt.X), int16(objPnt.Y)), Zlayer, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart)
		}
		if frm < 1 {
			delete(p.TmpObj, obj)
		}else{
			p.TmpObj[obj] = frm-1
		}
	}
}

func (p *WorldStructure) AddTile(img *Tile) {
	if p.Tiles == nil {
		p.Tiles = make([]*Tile, 0)
	}
	p.Tiles = append(p.Tiles, img)
}
func (p *WorldStructure) AddStructureObj(obj *StructureObj) {
	if p.StructureObjs == nil {
		p.StructureObjs = make([]*StructureObj, 0)
	}
	p.StructureObjs = append(p.StructureObjs, obj)
	p.UpdateCollisionMat()
}
func (p *WorldStructure) AddTempObj(obj *StructureObj, frames int) {
	if p.TmpObj == nil {
		p.TmpObj = make(map[*StructureObj]int)
	}
	p.TmpObj[obj] = frames
	p.UpdateCollisionMat()
}
func (p *WorldStructure) Save(path string) error {
	idxBs, err1 := p.IdxMat.Compress()
	if err1 != nil {return err1}
	layBs, err2 := p.LayerMat.Compress()
	if err2 != nil {return err2}
	mats := append(idxBs, layBs...)
	bs, err3 := CompressBytes(append(mats, AppendInt16ToBytes( int16(len(layBs)), Int16ToBytes(int16(len(idxBs))) )...))
	if err3 != nil {return err3}
	return ioutil.WriteFile(path, bs, 0644)
}
func (p *WorldStructure) Load(path string) error {
	data, err1 := ioutil.ReadFile(path)
   	if err1 != nil {return err1}
   	bs, err2 := DecompressBytes(data)
   	if err2 != nil {return err2}
   	lenIdx := BytesToInt16(bs[len(bs)-4:len(bs)-2])
   	lenLay := BytesToInt16(bs[len(bs)-2:len(bs)])
   	err3 := p.IdxMat.Decompress(bs[:lenIdx])
   	if err3 != nil {return err3}
   	err4 := p.LayerMat.Decompress(bs[lenIdx:lenIdx+lenLay])
   	if err4 != nil {return err4}
   	return nil
}
//returns the middle of the view
func (p *WorldStructure) Middle() (int16, int16) {
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
func (p *WorldStructure) GetTopLeft() *Point {
	return &Point{p.xStart, p.yStart}
}