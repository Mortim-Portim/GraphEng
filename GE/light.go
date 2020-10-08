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
	extinctionRate float64
	radius float64
	radiusNeedsUpdate bool
}

func GetLightSource(loc *Rectangle, maxI int16, extRate float64) (l *Light) {
	l = &Light{Location:loc, maximumIntesity:maxI, extinctionRate:extRate}
	l.CalcRadius()
	return
}

func (l *Light) LightMatrix(mat *Matrix) {
	rad := l.GetRadius()
	md := l.Location.GetMiddle()
	for x := int(md.X-rad-1); x < int(md.X+rad+1); x++ {
		for y := int(md.Y-rad-1); y < int(md.Y+rad+1); y++ {
			pnt := &Point{float64(x), float64(y)}
			if !pnt.InBounds(l.Location) {
				mat.Add(x,y, int16(l.GetValueAtXY(x,y)))
			}else{
				mat.Add(x,y, l.maximumIntesity)
			}
		}
	}
}
func (l *Light) GetValueAtXY(x,y int) float64 {
	return l.getValueAtRadius(l.Location.DistanceTo(&Point{float64(x),float64(y)}))
}
func (l *Light) getValueAtXYdif(xd,yd int) float64 {
	return l.getValueAtRadius(math.Pow(math.Pow(float64(xd), 2)+math.Pow(float64(yd), 2),1.0/2.0))
}
func (l *Light) getValueAtRadius(r float64) float64 {
	val := float64(l.maximumIntesity)/(math.Pow(l.extinctionRate*LIGHT_DIS_FACTOR*r+1, 3))
	if val > LIGHT_EXTINCION_LEVEL {
		return val
	}
	return 0
}

func (l *Light) SetRadius(r float64) {
	l.extinctionRate = ((math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCION_LEVEL, 1.0/3.0)-1)/(r*LIGHT_DIS_FACTOR))
	l.radiusNeedsUpdate = true
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
}
func (l *Light) SetExtinctionRate(r float64) {
	l.extinctionRate = r
	l.radiusNeedsUpdate = true
}
func (l *Light) CalcRadius() float64 {
	l.radius = (math.Pow(float64(l.maximumIntesity)/LIGHT_EXTINCION_LEVEL, 1.0/3.0)-1)/(l.extinctionRate*LIGHT_DIS_FACTOR)
	return l.radius
}
func (l *Light) GetRadius() float64 {
	if l.radiusNeedsUpdate {
		l.CalcRadius()
	}
	return l.radius
}