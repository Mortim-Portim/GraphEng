package GE

import (
	"fmt"
	"math"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

type Point struct {
	X, Y float64
}

func (p *Point) ToVec() *Vector {
	return &Vector{p.X, p.Y, 0}
}
func (p *Point) Equals(p2 *Point) bool {
	if p.X == p2.X && p.Y == p2.Y {
		return true
	}
	return false
}
func (p *Point) Print() string {
	return fmt.Sprintf("X:%v, Y:%v", p.X, p.Y)
}
func (p *Point) Copy() *Point {
	return &Point{p.X, p.Y}
}
func (p *Point) DistanceTo(p2 *Point) float64 {
	return math.Pow(math.Pow(p.X-p2.X, 2)+math.Pow(p.Y-p2.Y, 2), 1.0/2.0)
}
func (p *Point) InBounds(r *Rectangle) bool {
	if p.X >= r.Min().X && p.X <= r.Max().X && p.Y >= r.Min().Y && p.Y <= r.Max().Y {
		return true
	}
	return false
}
func (p *Point) ToBytes() []byte {
	return append(cmp.Float64ToBytes(p.X), cmp.Float64ToBytes(p.Y)...)
}
func PointFromBytes(bs []byte) (p *Point) {
	p = &Point{}
	p.X = cmp.BytesToFloat64(bs[:8])
	p.Y = cmp.BytesToFloat64(bs[8:])
	return
}

/**
Rectangle represents two Points spanning a rectangle:

P1
#----------------------+
|					   |
|					   |
|					   |
|					   |
+----------------------#
					   P2

the width and height of the Rectangle can be accessed by calling Bounds()
**/
func GetRectangle(x1, y1, x2, y2 float64) (r *Rectangle) {
	r = &Rectangle{&Point{x1, y1}, &Point{x2, y2}, nil}
	r.updateBounds()
	return
}

type Rectangle struct {
	min, max, bounds *Point
}

func (r *Rectangle) Print() string {
	return fmt.Sprintf("Min: %s, Max: %s, Bounds: %s", r.Min().Print(), r.Max().Print(), r.Bounds().Print())
}
func (r *Rectangle) Copy() *Rectangle {
	return &Rectangle{r.min.Copy(), r.max.Copy(), r.bounds.Copy()}
}
func (r *Rectangle) MoveTo(pnt *Point) {
	w, h := r.Bounds().X, r.Bounds().Y
	r.SetMin(pnt)
	r.SetBounds(&Point{w, h})
}
func (r *Rectangle) MoveBy(dx, dy float64) {
	r.min.X += dx
	r.min.Y += dy
	r.max.X += dx
	r.max.Y += dy
}
func (r *Rectangle) SetMin(min *Point) {
	r.min = min
	r.updateBounds()
}
func (r *Rectangle) SetMax(max *Point) {
	r.max = max
	r.updateBounds()
}
func (r *Rectangle) SetBounds(bounds *Point) {
	r.SetMax(&Point{r.min.X + bounds.X, r.min.Y + bounds.Y})
}
func (r *Rectangle) GetMiddle() *Point {
	return &Point{r.min.X + r.bounds.X/2, r.min.Y + r.bounds.Y/2}
}
func (r *Rectangle) SetMiddle(pnt *Point) {
	r.SetMin(&Point{pnt.X - r.Bounds().X/2, pnt.Y - r.Bounds().Y/2})
}
func (r *Rectangle) Min() *Point    { return r.min }
func (r *Rectangle) Max() *Point    { return r.max }
func (r *Rectangle) Bounds() *Point { return r.bounds }
func (r *Rectangle) updateBounds() {
	r.bounds = &Point{r.max.X - r.min.X, r.max.Y - r.min.Y}
}
func (r *Rectangle) Inside(r2 *Rectangle) bool {
	if r.Min().X > r2.Min().X && r.Min().Y > r2.Min().Y && 
		r.Max().X < r2.Max().X && r.Max().Y < r2.Max().Y {
		return true
	}
	return false
}
func (r *Rectangle) Overlaps(r2 *Rectangle) bool {
	// If one rectangle is on left side of other
	if r.Min().X >= r2.Max().X || r2.Min().X >= r.Max().X {
		return false
	}
	// If one rectangle is above other
	if r.Min().Y >= r2.Max().Y || r2.Min().Y >= r.Max().Y {
		return false
	}
	return true
}
func (r *Rectangle) DistanceTo(p *Point) float64 {
	return p.DistanceTo(r.GetMiddle())
}

func (r *Rectangle) ToBytes() []byte {
	return append(r.min.ToBytes(), r.max.ToBytes()...)
}
func RectangleFromBytes(bs []byte) (r *Rectangle) {
	r = &Rectangle{}
	r.min = PointFromBytes(bs[:16])
	r.max = PointFromBytes(bs[16:])
	r.updateBounds()
	return
}
func GetOrderedRectangleI(p1, p2 [2]int) *Rectangle {
	return GetOrderedRectangleF([2]float64{float64(p1[0]), float64(p1[1])}, [2]float64{float64(p2[0]), float64(p2[1])})
}

func GetOrderedRectangleF(p1, p2 [2]float64) *Rectangle {
	minX := p1[0];maxX := p2[0]
	minY := p1[1];maxY := p2[1]
	if p1[0] > p2[0] {
		minX = p2[0];maxX = p1[0]
	}
	if p1[1] > p2[1] {
		minY = p2[1];maxY = p1[1]
	}
	return GetRectangle(minX, minY, maxX, maxY)
}

/**
func (r *Rectangle) GetLines() (l []*line) {
	l = make([]*line, 4)
	l[0] = &line{r.min.X, r.min.Y, r.max.X, r.min.Y}
	l[1] = &line{r.max.X, r.min.Y, r.max.X, r.max.Y}
	l[2] = &line{r.max.X, r.max.Y, r.min.X, r.max.Y}
	l[3] = &line{r.min.X, r.max.Y, r.min.X, r.min.Y}
	return
}
**/
