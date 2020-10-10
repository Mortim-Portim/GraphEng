package GE

import (
	"math"
	"marvin/GraphEng/GC"
)

func (l *Light) ApplyRaycasting(factor float64) {
	rad := l.GetRadius()
	l.LightMat = GetMatrix(int(rad*2+1), int(rad*2+1), 0)
	//dira := l.direction.GetRotationZ()
	rays := int(l.angle/l.accuracy)
	
	for i := 0; i < rays; i++ {
		a := float64(i)*l.accuracy-l.angle/2
		//vec := GC.GetVectorFromRot(a)
		startPnt := l.Location.Clone()
		endPnt := 
		for startPnt.X >= 0 && startPnt.Y >= 0 && startPnt.X <= 2*rad && startPnt.Y <= 2*rad {
			
		}
	}
	
	/**
	for x := int(md.X-rad-1); x < int(md.X+rad+1); x++ {
		for y := int(md.Y-rad-1); y < int(md.Y+rad+1); y++ {
			pnt := &Point{float64(x), float64(y)}
			if !pnt.InBounds(l.Location) {
				realPnt := &Point{float64(x)+0.5, float64(y)+0.5}
				ray := &line{md.X, md.Y, realPnt.X, realPnt.Y}
				shadow := false
				for x2 := int(md.X-rad-1); x2 < int(md.X+rad+1); x2++ {
					for y2 := int(md.Y-rad-1); y2 < int(md.Y+rad+1); y2++ {
						if collM.GetAbs(x2,y2) != NON_COLLIDING_IDX {
							ln1 := &line{float64(x2),float64(y2), float64(x2+1), float64(y2+1)}
							ln2 := &line{float64(x2),float64(y2+1), float64(x2+1), float64(y2)}
							_,_,i1 := intersection(ln1, ray)
							_,_,i2 := intersection(ln2, ray)
							if i1 || i2 {
								shadow = true
							}
						}
					}
				}
				if !shadow {
					lightM.AddAbs(x,y, int16(l.GetValueAtXY(x,y)*factor))
				}
			}else{
				lightM.AddAbs(x,y, int16(float64(l.maximumIntesity)*factor))
			}
		}
	}
	**/
}

type line struct {
	X1, Y1, X2, Y2 float64
}
func NewRay(x, y, length, angle float64) line {
	return line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

// intersection calculates the intersection of given two lines.
func intersection(l1, l2 *line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)
	return x, y, true
}