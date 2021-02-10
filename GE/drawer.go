package GE

import (
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/vector"
)

/**
Points are a slice of vectors, that can be drawn as a polygon using a specific color

This is currently only used to draw simple forms such as a Circle
**/

type Points []*Vector
func (ps *Points) Fill(screen *ebiten.Image, col color.Color) {
	var path vector.Path
	for i,p := range(*ps) {
		if i == 0 {
			path.MoveTo(float32(p.X), float32(p.Y))
		}else{
			path.LineTo(float32(p.X), float32(p.Y))
		}
	}
	op := &vector.FillOptions{
		Color: col,
	}
	path.Fill(screen, op)
}
func GetHorizontalBar(x,y,w,h float64) *Points {
	return GetLineOfPoints(x,y+h/2, x+w, y+h/2, h/2)
}
func GetLineOfPoints(x1,y1, x2,y2, thickness float64) *Points {
	ps := make([]*Vector, 4)
	
	loc := &Vector{x1, y1, 0}
	dir := &Vector{x2-x1, y2-y1, 0}
	
	right := dir.CrossProduct(&Vector{0,0,1}).Normalize().Mul(thickness)
	left := dir.CrossProduct(&Vector{0,0,-1}).Normalize().Mul(thickness)
	
	ps[0] = loc.Add(right)
	ps[3] = loc.Add(left)
	ps[1] = loc.Add(right).Add(dir)
	ps[2] = loc.Add(left).Add(dir)
	
	newPs := Points(ps)
	return &newPs
}