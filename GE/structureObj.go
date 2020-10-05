package GE

import (
	"github.com/hajimehoshi/ebiten"
)

type StructureObj struct {
	Animation
	frame, squareSize int
	HitBox, DrawBox *Rectangle
	Collides bool
}

func GetStructureObj(anim *Animation, HitBox *Rectangle, squareSize int, Collides bool) (o *StructureObj) {
	o = &StructureObj{Animation:*anim, frame:0, HitBox:HitBox, squareSize:squareSize, Collides:Collides}
	pnt := HitBox.Min()
	o.SetToXY(int16(pnt.X), int16(pnt.Y))
	return
}
func (o *StructureObj) SetToXY(x,y int16) {
	o.HitBox.MoveTo(&Point{float64(x),float64(y)})
	w,h := o.Img.Size()
	W := float64(w)/float64(o.squareSize); H := float64(h)/float64(o.squareSize)
	o.DrawBox = GetRectangle(o.HitBox.Min().X-(W-o.HitBox.Bounds().X-1)/2, o.HitBox.Min().Y-(H-o.HitBox.Bounds().Y-1), 0,0)
	o.DrawBox.SetBounds(&Point{W,H})
}

func (o *StructureObj) DrawCollisionMatrix(mat *Matrix) {
	if o.Collides {
		mat.Fill(int16(o.HitBox.Min().X), int16(o.HitBox.Min().Y), int16(o.HitBox.Max().X), int16(o.HitBox.Max().Y), COLLIDING_IDX)
	}
}

func (o *StructureObj) Draw(screen *ebiten.Image, myLayer, drawLayer int, leftTop *Point, sqSize, xStart, yStart float64) {
	o.Update(o.frame)
	
	relPx, relPy := float64(o.DrawBox.Min().X-leftTop.X), float64(o.DrawBox.Min().Y-leftTop.Y)
	o.X = relPx*sqSize+xStart
	o.Y = relPy*sqSize+yStart
	o.W = float64(o.DrawBox.Bounds().X)*sqSize
	o.H = float64(o.DrawBox.Bounds().Y)*sqSize
	
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