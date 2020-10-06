package GE

import (
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
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
	StructureObjs	[]*StructureObj
	TmpObj			map[*StructureObj]int
	
	IdxMat, LayerMat, CollisionMat *Matrix
	xTiles, yTiles, middleX, middleY int
	
	frame  *ImageObj
	drawer *ImageObj

	X,Y,W,H, xStart, yStart, tileS float64
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
func (p *WorldStructure) Collides(x,y int) bool {
	if p.CollisionMat.Get(x,y) == COLLIDING_IDX {
		return true
	}
	return false
}


//sets the Middle of the View
func (p *WorldStructure) SetMiddle(pnt *Point) {
	p.middleX = int(pnt.X)
	p.middleY = int(pnt.Y)
	x,y := p.middleX-(p.xTiles-1)/2, p.middleY-(p.yTiles-1)/2
	p.IdxMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	p.LayerMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
}
//moves the view by dx and dy
func (p *WorldStructure) Move(dx,dy int) {
	p.middleX += dx
	p.middleY += dy
	x,y := p.middleX-(p.xTiles-1)/2, p.middleY-(p.yTiles-1)/2
	p.IdxMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
	p.LayerMat.SetFocus(x,y, x+p.xTiles, y+p.yTiles)
}
//Draws The World Ground Tiles and the Objects form the layer which is currently in the middle of the screen
func (p *WorldStructure) Draw(screen *ebiten.Image) {
	Zlayer := int(p.LayerMat.GetAbs(p.Middle()))
	for x := 0; x < p.IdxMat.W(); x++ {
		for y := 0; y < p.IdxMat.H(); y++ {
			tile_idx := p.IdxMat.Get(x, y)
			xp, yp := float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
			if int(tile_idx) >= 0 && int(tile_idx) < len(p.Tiles) {
				p.Tiles[tile_idx].Draw(screen, xp, yp, p.tileS, p.tileS, int(p.LayerMat.Get(x, y)), Zlayer, p.drawer, p.frame)
			}
		}
	}
	for _,obj := range(p.StructureObjs) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			objPnt := obj.HitBox.Min()
			obj.Draw(screen, int(p.LayerMat.GetAbs(int(objPnt.X), int(objPnt.Y))), Zlayer, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart)
		}
	}
	for obj,frm := range(p.TmpObj) {
		if p.IdxMat.Focus().Overlaps(obj.DrawBox) {
			objPnt := obj.HitBox.Min()
			obj.Draw(screen, int(p.LayerMat.GetAbs(int(objPnt.X), int(objPnt.Y))), Zlayer, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart)
		}
		if frm < 1 {
			delete(p.TmpObj, obj)
		}else{
			p.TmpObj[obj] = frm-1
		}
	}
}
//Adds a tile to the index list of the worlds tiles
func (p *WorldStructure) AddTile(img *Tile) {
	if p.Tiles == nil {
		p.Tiles = make([]*Tile, 0)
	}
	p.Tiles = append(p.Tiles, img)
}
//Adds a StructureObj to the list of the worlds Objs
func (p *WorldStructure) AddStructureObj(obj *StructureObj) {
	if p.StructureObjs == nil {
		p.StructureObjs = make([]*StructureObj, 0)
	}
	p.StructureObjs = append(p.StructureObjs, obj)
	p.UpdateCollisionMat()
}
//Adds a temporary Obj to the map of the worlds TmpObj
func (p *WorldStructure) AddTempObj(obj *StructureObj, frames int) {
	if p.TmpObj == nil {
		p.TmpObj = make(map[*StructureObj]int)
	}
	p.TmpObj[obj] = frames
	p.UpdateCollisionMat()
}
//Converts the World into a []byte slice
func (p *WorldStructure) ToBytes() ([]byte, error) {
	idxBs, err1 := p.IdxMat.Compress()
	if err1 != nil {return nil, err1}
	layBs, err2 := p.LayerMat.Compress()
	if err2 != nil {return nil, err2}
	mats := append(idxBs, layBs...)
	bs, err3 := CompressBytes(append(mats, AppendInt16ToBytes( int16(len(layBs)), Int16ToBytes(int16(len(idxBs))) )...))
	if err3 != nil {return nil, err3}
	return bs, nil
}
//Saves the world in a highly compressed way to the file system
func (p *WorldStructure) Save(path string) error {
	bs, err := p.ToBytes()
	if err != nil {return err}
	return ioutil.WriteFile(path, bs, 0644)
}
//Converts a []byte slice into a WorldStructure
func (p *WorldStructure) FromBytes(data []byte) error {
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
//Loads the world from the file system
func (p *WorldStructure) Load(path string) error {
	data, err1 := ioutil.ReadFile(path)
   	if err1 != nil {return err1}
   	err2 := p.FromBytes(data)
   	if err2 != nil {return err2}
   	return nil
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
func (p *WorldStructure) GetTopLeft() *Point {
	return &Point{p.xStart, p.yStart}
}