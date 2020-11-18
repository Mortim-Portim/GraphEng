package GE

import (
	"errors"
	"fmt"
	"math"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

const LIGHT_DIS_FACTOR = 20
const LIGHT_EXTINCTION_LEVEL = 10
const LIGHT_EXTINCTION_POWER = 3

type Light struct {
	location  *Point
	direction *Vector
	angle     float64
	lightMat  *Matrix
	static    bool

	//0-255
	maximumIntesity int16
	//0-0.05
	extinctionRate                                float64
	radius                                        float64
	radiusNeedsUpdate, matrixNeedsUpdate, changed bool
}

func GetLightSource(loc *Point, direction *Vector, angle float64, maxI int16, extRate float64, static bool) (l *Light) {
	l = &Light{location: loc, maximumIntesity: maxI, extinctionRate: extRate, direction: direction, angle: angle, static: static, lightMat: nil}
	l.CalcRadius()
	l.matrixNeedsUpdate = true
	l.changed = true
	return
}

func (l *Light) ToBytes() (b []byte) {
	bs := make([][]byte, 0)
	bs = append(bs, l.location.ToBytes())
	bs = append(bs, l.direction.ToBytes())
	bs = append(bs, cmp.Float64ToBytes(l.angle))
	bs = append(bs, cmp.Float64ToBytes(l.extinctionRate))
	bs = append(bs, cmp.Int16ToBytes(l.maximumIntesity))
	bs = append(bs, []byte{cmp.BoolToByte(l.static)})
	b = cmp.CompressAll([][]byte{}, bs...)
	return
}
func GetLightSourceFromBytes(b []byte) (l *Light) {
	bs := cmp.DecompressAll(b, []int{16, 24, 8, 8, 2, 1})
	loc := PointFromBytes(bs[0])
	dir := VectorFromBytes(bs[1])
	ang := cmp.BytesToFloat64(bs[2])
	//acc := cmp.BytesToFloat64(bs[3])
	ext := cmp.BytesToFloat64(bs[3])
	maI := cmp.BytesToInt16(bs[4])
	stt := cmp.ByteToBool(bs[5][0])
	l = GetLightSource(loc, dir, ang, maI, ext, stt)
	return
}
func (l *Light) Move(dx, dy float64) {
	l.location.X += dx
	l.location.Y += dy
}
func (l *Light) GetAtAbs(x, y int) (int64, error) {
	if l.lightMat == nil {
		return 0, errors.New(fmt.Sprintf("LightMat of %v not yet initialized", l))
	}
	xp := x - int(l.location.X) + int(float64(l.lightMat.WAbs()-1)/2.0)
	yp := y - int(l.location.Y) + int(float64(l.lightMat.HAbs()-1)/2.0)
	return l.lightMat.GetAbs(xp, yp)
}
func (l *Light) applyOnMatrix(mat *Matrix, factor float64) {
	rad := l.GetRadius()
	md := l.location
	for x := int(md.X - rad - 1); x < int(md.X+rad+1); x++ {
		for y := int(md.Y - rad - 1); y < int(md.Y+rad+1); y++ {
			pnt := &Point{float64(x), float64(y)}
			if pnt.X == md.X && pnt.Y == md.Y {
				mat.Add(x, y, int64(float64(l.maximumIntesity)*factor))
			} else {
				mat.Add(x, y, int64(l.GetValueAtXY(x, y)*factor))
			}
		}
	}
}
func (l *Light) GetValueAtXY(x, y int) float64 {
	return l.getValueAtRadius(l.location.DistanceTo(&Point{float64(x), float64(y)}))
}
func (l *Light) getValueAtXYdif(xd, yd int) float64 {
	return l.getValueAtRadius(math.Pow(math.Pow(float64(xd), 2)+math.Pow(float64(yd), 2), 1.0/2.0))
}
func (l *Light) getValueAtRadius(r float64) float64 {
	val := float64(l.maximumIntesity) / (math.Pow(l.extinctionRate*LIGHT_DIS_FACTOR*r+1, LIGHT_EXTINCTION_POWER))
	if val > LIGHT_EXTINCTION_LEVEL {
		return val
	}
	return 0
}

func (l *Light) SetRadiusByMaxI(r float64) {
	l.extinctionRate = ((math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCTION_LEVEL, 1.0/LIGHT_EXTINCTION_POWER) - 1) / (r * LIGHT_DIS_FACTOR))
	l.radiusNeedsUpdate = true
	l.matrixNeedsUpdate = true
	l.changed = true
}

func (l *Light) GetMaximumIntesity() int16 {
	return l.maximumIntesity
}
func (l *Light) GetExtinctionRate() float64 {
	return l.extinctionRate
}
func (l *Light) SetMaximumIntesity(r int16) {
	l.maximumIntesity = r
	l.radiusNeedsUpdate = true
	l.matrixNeedsUpdate = true
	l.changed = true
}
func (l *Light) SetExtinctionRate(r float64) {
	l.extinctionRate = r
	l.radiusNeedsUpdate = true
	l.matrixNeedsUpdate = true
	l.changed = true
}
func (l *Light) CalcRadius() float64 {
	l.radius = (math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCTION_LEVEL, 1.0/LIGHT_EXTINCTION_POWER) - 1) / (l.extinctionRate * LIGHT_DIS_FACTOR)
	return l.radius
}
func (l *Light) GetRadius() float64 {
	if l.radiusNeedsUpdate {
		l.CalcRadius()
	}
	return l.radius
}

func (l *Light) Loc() *Point {
	return l.location
}
func (l *Light) SetLoc(pnt *Point) {
	l.location = pnt
	l.changed = true
}
func (l *Light) Dir() *Vector {
	return l.direction
}
func (l *Light) SetDir(pnt *Vector) {
	l.direction = pnt
	l.matrixNeedsUpdate = true
	l.changed = true
}
func (l *Light) Angle() float64 {
	return l.angle
}
func (l *Light) SetAngle(a float64) {
	l.angle = a
	l.matrixNeedsUpdate = true
	l.changed = true
}
func (l *Light) SetDark() {
	l.SetMatrix(GetMatrix(1, 1, 0))
	l.changed = true
}
func (l *Light) Matrix() *Matrix {
	return l.lightMat
}
func (l *Light) SetMatrix(m *Matrix) {
	l.lightMat = m
}

func (l *Light) Changed() bool {
	return l.changed
}
func (l *Light) SetChanged(c bool) {
	l.changed = c
}
