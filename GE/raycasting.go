package GE

import (
	"math"
)

func (l *Light) ApplyRaycasting(collMat *Matrix, factor float64) {
	if !l.static || l.LightMat == nil {
		rad := l.GetRadius()
		CollMat := collMat.SubMatrix(int(l.Location.X-rad), int(l.Location.Y-rad), int(l.Location.X+rad), int(l.Location.Y+rad))
		length := rad*2+1
		l.LightMat = GetMatrix(int(length), int(length), 0)
		dira := l.direction.GetRotationZ()
		pnts := getPointsForBox(int(length))
		mina := dira-l.angle/2
		maxa := dira+l.angle/2
		
		for _,op := range(pnts) {
			pnt := &Vector{rad, rad, 0}
			a := op.Sub(pnt).GetRotationZ(); aI1 := a-360; aI2 := a+360
			
			if isInBounds(a, mina, maxa) || isInBounds(aI1, mina, maxa) || isInBounds(aI2, mina, maxa) {
				dx := float64(int(op.X-pnt.X))
				dy := float64(int(op.Y-pnt.Y))
				
				l.iterateOverLine(dx,dy,rad,factor,pnt,CollMat)
			}
		}
	}
}

func (l *Light) iterateOverLine(dx, dy, rad, factor float64, pnt *Vector, CollMat *Matrix) {
	xStep := 1.0
	if dx < 0 {
		xStep = -1.0
	}
	yStep := 1.0
	if dy < 0 {
		yStep = -1.0
	}
	stepYpX := dy/dx
	stepXpY := dx/dy
	for pnt.X >= 0 && pnt.Y >= 0 && pnt.X <= 2*rad && pnt.Y <= 2*rad {
		x,y := int(pnt.X), int(pnt.Y)
		if CollMat.Get(x,y) != NON_COLLIDING_IDX {
			break
		}
		if l.LightMat.Get(x,y) == 0 {
			val := int16(l.getValueAtXYdif(x-int(rad),y-int(rad))*factor)
			l.LightMat.Set(x,y, val)
		}
		if math.Abs(dx) > math.Abs(dy) {
			pnt.X += xStep
			pnt.Y += xStep*stepYpX
		}else{
			pnt.Y += yStep
			pnt.X += yStep*stepXpY
		}
	}
}

func isInBounds(val, bl, bm float64) bool {
	if val >= bl && val <= bm {
		return true
	}
	return false
}

func getPointsForBox(length int) (ps []*Vector) {
	ps = make([]*Vector, 0)
	for x := 0; x < length; x++ {
		ps = append(ps, &Vector{float64(x), 0,0})
	}
	for y := 0; y < length; y++ {
		ps = append(ps, &Vector{0, float64(y),0})
	}
	for x := 0; x < length; x++ {
		ps = append(ps, &Vector{float64(x), float64(length),0})
	}
	for y := 0; y < length; y++ {
		ps = append(ps, &Vector{float64(length), float64(y),0})
	}
	return
}