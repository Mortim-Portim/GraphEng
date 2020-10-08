package GE

import (
	"math"
)
const LIGHT_DIS_FACTOR = 20
const LIGHT_EXTINCION_LEVEL = 10

type Light struct {
	Location *Rectangle
	//0-255
	maximumIntesity int16
	//0-0.05
	extinctionRate int16
	radius float64
	radiusNeedsUpdate bool
}

func (l *Light) LightMatrix(mat *Matrix) {
	rad := l.GetRadius()
	minX := l.Location.Min().X
	
}
func (l *Light) GetValueAtXY(x,y int) float64 {
	return l.getValueAtRadius(l.Location.DistanceTo(&Point{float64(x),float64(y)}))
}
func (l *Light) getValueAtXYdif(xd,yd int) float64 {
	return l.getValueAtRadius(math.Pow(math.Pow(float64(xd), 2)+math.Pow(float64(yd), 2),1.0/2.0))
}
func (l *Light) getValueAtRadius(r float64) float64 {
	val := float64(l.maximumIntesity)/(math.Pow(float64(l.extinctionRate)*LIGHT_DIS_FACTOR*r+1, 3))
	if val > LIGHT_EXTINCION_LEVEL {
		return val
	}
	return 0
}

func (l *Light) SetRadius(r float64) {
	l.extinctionRate = int16((math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCION_LEVEL, 1.0/3.0)-1)/(r*LIGHT_DIS_FACTOR))
	l.radiusNeedsUpdate = true
}

func (l *Light) GetMaximumIntesity() int16 {
	return l.maximumIntesity
}
func (l *Light) GetExtinctionRate() int16 {
	return l.extinctionRate
}
func (l *Light) SetMaximumIntesity(r int16) {
	l.maximumIntesity = r
	l.radiusNeedsUpdate = true
}
func (l *Light) SetExtinctionRate(r int16) {
	l.extinctionRate = r
	l.radiusNeedsUpdate = true
}
func (l *Light) CalcRadius() float64 {
	l.radius = (math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCION_LEVEL, 1.0/3.0)-1)/(float64(l.extinctionRate)*LIGHT_DIS_FACTOR)
	return l.radius
}
func (l *Light) GetRadius() float64 {
	if l.radiusNeedsUpdate {
		l.CalcRadius()
	}
	return l.radius
}