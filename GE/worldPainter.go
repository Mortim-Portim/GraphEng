package GE

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE/WObjs"
)

type WorldPainter struct {
	tiles  []*WObjs.Tile
	frame  *ImageObj
	drawer *ImageObj

	xStart, yStart, tileS float64
}

func GetWorldPainter(X, Y, W, H float64, Xtiles, Ytiles int) (p *WorldPainter) {
	p = &WorldPainter{}
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
func (p *WorldPainter) GetFrame(thickness float64, alpha uint8) {
	p.frame = p.drawer.GetFrame(thickness, alpha)
}
func (p *WorldPainter) AddTile(img *WObjs.Tile) {
	if p.tiles == nil {
		p.tiles = make([]*WObjs.Tile, 0)
	}
	p.tiles = append(p.tiles, img)
}

func (p *WorldPainter) Paint(screen *ebiten.Image, idxMat, layerMat *Matrix, Zlayer int) {
	for x := 0; x < idxMat.X; x++ {
		for y := 0; y < idxMat.Y; y++ {
			idx := idxMat.Get(x, y, 0)
			if idx >= 0 || idx < len(p.tiles) {
				p.drawer.Img = p.tiles[idx].Img
				p.drawer.X = float64(x)*p.tileS + p.xStart
				p.drawer.Y = float64(y)*p.tileS + p.yStart
				layerIdx := layerMat.Get(x, y, 0)
				if layerIdx == Zlayer {
					p.drawer.DrawImageObj(screen)
				} else if layerIdx < Zlayer {
					dif := 1.0 / float64(Zlayer-layerIdx+1)
					p.drawer.DrawImageObjAlpha(screen, dif)
				} else {
					box := float64(1 + layerIdx - Zlayer)
					sq := box*2 + 1
					p.drawer.DrawImageBlured(screen, int(box), 1.0/((sq*sq)*0.35))
				}
				if p.frame != nil {
					p.frame.X = p.drawer.X
					p.frame.Y = p.drawer.Y
					p.frame.DrawImageObj(screen)
				}
			}

		}
	}
}
