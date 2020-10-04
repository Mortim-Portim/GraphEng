package GE

import (
	"github.com/hajimehoshi/ebiten"
)
/**
Matrix
	Standard Index = -1
	Layer 0: Ground
**/
type WorldStructure struct {
	tiles  			[]*Tile
	structureObjs	[]*StructureObj
	
	IdxMat, LayerMat, CollisionMat *Matrix
	
	frame  *ImageObj
	drawer *ImageObj

	xStart, yStart, tileS float64
}
func GetWorldStructure(X, Y, W, H float64, Xtiles, Ytiles int) (p *WorldStructure) {
	p = &WorldStructure{}
	p.xStart = X
	p.tileS = W / float64(Xtiles)
	p.yStart = Y - (float64(Ytiles)*p.tileS-H)/2
	if W < H {
		p.tileS = H / float64(Ytiles)
		p.xStart = X - (float64(Xtiles)*p.tileS-W)/2
		p.yStart = Y
	}

	p.drawer = &ImageObj{}
	p.drawer.W = p.tileS
	p.drawer.H = p.tileS
	return
}
func (p *WorldStructure) GetFrame(thickness float64, alpha uint8) {
	p.frame = p.drawer.GetFrame(thickness, alpha)
}


//TODO implement tiles, collObj, nonCollObjs, idxMat, layerMat loader and writer
func (p *WorldStructure) AddTile(img *Tile) {
	if p.tiles == nil {
		p.tiles = make([]*Tile, 0)
	}
	p.tiles = append(p.tiles, img)
}
func (p *WorldStructure) AddStructureObj(obj *StructureObj) {
	if p.structureObjs == nil {
		p.structureObjs = make([]*StructureObj, 0)
	}
	p.structureObjs = append(p.structureObjs, obj)
}

func (p *WorldStructure) Collides(x,y int) bool {
	if p.IdxMat.Get(x,y, 1) >= 0 {
		return true
	}
	return false
}

func (p *WorldStructure) Draw(screen *ebiten.Image, Zlayer int) {
	for x := 0; x < p.IdxMat.X; x++ {
		for y := 0; y < p.IdxMat.Y; y++ {
			tile_idx := p.IdxMat.Get(x, y, 0)
			coll_idx := p.IdxMat.Get(x, y, 1)
			no_coll_idx := p.IdxMat.Get(x, y, 2)
			xp, yp := float64(x)*p.tileS + p.xStart, float64(y)*p.tileS + p.yStart
			if tile_idx >= 0 && tile_idx < len(p.tiles) {
				p.tiles[tile_idx].Draw(screen, xp, yp, p.tileS, p.tileS, p.LayerMat.Get(x, y, 0), Zlayer, p.drawer, p.frame)
			}
			if coll_idx >= 0 && coll_idx < len(p.structureObjs) {
				p.structureObjs[coll_idx].Draw(screen, xp, yp, p.tileS, p.tileS, p.LayerMat.Get(x, y, 0), Zlayer)
			}
			if no_coll_idx >= 0 && no_coll_idx < len(p.structureObjs) {
				p.structureObjs[no_coll_idx].Draw(screen, xp, yp, p.tileS, p.tileS, p.LayerMat.Get(x, y, 0), Zlayer)
			}
		}
	}
}
