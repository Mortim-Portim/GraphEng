package GE

import (
	"github.com/hajimehoshi/ebiten"
)

type StructureObj struct {
	Animation
	frame int
	XTile,YTile,TileW,TileH int
}

func (o *StructureObj) DrawCollisionMatrix(mat *Matrix) {
	mat.Fill(o.XTile, o.YTile, o.TileW, o.TileH, 0, 1)
}

func (o *StructureObj) Draw(screen *ebiten.Image, x, y, w, h float64, myLayer, drawLayer int) {
	o.Update(o.frame)
	o.X = x; o.Y = y; o.W = w; o.H = h
	if myLayer == drawLayer {
		o.DrawImageObj(screen)
	} else if myLayer < drawLayer {
		dif := 1.0 / float64(drawLayer-myLayer+1)
		o.DrawImageObjAlpha(screen, dif)
	} else {
		box := float64(1 + myLayer - drawLayer)
		sq := box*2 + 1
		o.DrawImageBlured(screen, int(box), 1.0/((sq*sq)*0.35))
	}
	o.frame ++
}