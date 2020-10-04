package GE

import (
	"github.com/hajimehoshi/ebiten"
	//"fmt"
)

//Returns a WorldStructure object
func GetWorldStructure(X, Y, W, H float64, WTiles, HTiles int) (p *WorldStructure) {
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
	
	IdxMat, LayerMat, CollisionMat *Matrix
	xTiles, yTiles, middleX, middleY int
	
	frame  *ImageObj
	drawer *ImageObj

	xStart, yStart, tileS float64
}
//Applies a frame
func (p *WorldStructure) GetFrame(thickness float64, alpha uint8) {
	p.frame = p.drawer.GetFrame(thickness, alpha)
}
//Updates the collision Matrix
func (p *WorldStructure) UpdateCollisionMat() {
	p.CollisionMat = GetMatrix(p.IdxMat.W(),p.IdxMat.H(),COLLIDING_IDX-1)
	for _,obj := range(p.StructureObjs) {
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
//returns the middle of the view
func (p *WorldStructure) Middle() (int, int) {
	return p.middleX, p.middleY
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
	Zlayer := p.LayerMat.GetAbs(p.Middle())
	for x := 0; x < p.IdxMat.W(); x++ {
		for y := 0; y < p.IdxMat.H(); y++ {
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
			//fmt.Println("Current Layer: ", p.LayerMat.GetAbs(int(objPnt.X), int(objPnt.Y)), ", at Point: ", objPnt.Print())
			//fmt.Println(p.LayerMat.Print())
			obj.Draw(screen, p.LayerMat.GetAbs(int(objPnt.X), int(objPnt.Y)), Zlayer, p.IdxMat.Focus().Min(), p.tileS, p.xStart, p.yStart)
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
}
