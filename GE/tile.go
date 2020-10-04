package GE

import (
	"github.com/hajimehoshi/ebiten"
)

type Tile struct {
	Img *ebiten.Image
}

func (t *Tile) Draw(screen *ebiten.Image, x, y, w, h float64, myLayer, drawLayer int, drawer, frame *ImageObj) {
	drawer.Img = t.Img
	drawer.X = x
	drawer.Y = y
	if myLayer == drawLayer {
		drawer.DrawImageObj(screen)
	} else if myLayer < drawLayer {
		dif := 1.0 / float64(drawLayer-myLayer+1)
		drawer.DrawImageObjAlpha(screen, dif)
	} else {
		box := float64(1 + myLayer - drawLayer)
		sq := box*2 + 1
		drawer.DrawImageBlured(screen, int(box), 1.0/((sq*sq)*0.35))
	}
	if frame != nil {
		frame.X = drawer.X
		frame.Y = drawer.Y
		frame.DrawImageObj(screen)
	}
}