package GE

import (
	"fmt"
)

type Point struct {
	X,Y float64
}
func (p *Point) Print() string {
	return fmt.Sprintf("X:%v, Y:%v", p.X, p.Y)
}
func (p *Point) Clone() (*Point) {
	return &Point{p.X, p.Y}
}

func GetRectangle(x1,y1,x2,y2 float64) (r *Rectangle) {
	r = &Rectangle{&Point{x1,y1}, &Point{x2,y2}, nil}
	r.updateBounds()
	return
}
type Rectangle struct {
	min, max, bounds *Point
}
func (r *Rectangle) Clone() (*Rectangle) {
	return &Rectangle{r.min.Clone(), r.max.Clone(), r.bounds.Clone()}
}
func (r *Rectangle) MoveTo(pnt *Point) {
	w,h := r.Bounds().X, r.Bounds().Y
	r.SetMin(pnt)
	r.SetBounds(&Point{w,h})
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
	r.SetMax(&Point{r.min.X+bounds.X, r.min.Y+bounds.Y})
}
func (r *Rectangle) GetMiddle() (*Point) {
	return &Point{r.min.X+r.bounds.X/2, r.min.Y+r.bounds.Y/2}
}
func (r *Rectangle) SetMiddle(pnt *Point) {
	r.SetMin(&Point{pnt.X-r.Bounds().X/2, pnt.Y-r.Bounds().Y/2})
}
func (r *Rectangle) Min() *Point {return r.min}
func (r *Rectangle) Max() *Point {return r.max}
func (r *Rectangle) Bounds() *Point {return r.bounds}
func (r *Rectangle) updateBounds() {
	r.bounds = &Point{r.max.X-r.min.X, r.max.Y-r.min.Y}
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