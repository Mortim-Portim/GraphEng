package GE

import (
	"math"
)

var strght = [][2]int{{0, -1}, {+1, 0}, {0, +1}, {-1, 0}}
var diag = [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

func GetNeighboursInBounds(x, y, rx1, ry1, rx2, ry2 int, collMat *Matrix) (ns [][2]int) {
	ns = make([][2]int, 0)
	for _, n := range strght {
		p := [2]int{n[0] + x, n[1] + y}
		if p[0] >= rx1 && p[1] >= ry1 && p[0] < rx2 && p[1] < ry2 {
			ns = append(ns, n)
		}
	}

	for _, n := range diag {
		cld1, err1 := collMat.Get(x+n[0], y)
		cld2, err2 := collMat.Get(x, y+n[1])
		if err1 == nil && cld1 == 0 && err2 == nil && cld2 == 0 {
			ns = append(ns, n)
		}
	}
	return
}

//Start and end must be greater or equal to 0 and in the bounds of the matrix focus
func FindPathMat(wrld *WorldStructure, start, end [2]int) (nodes [][2]int) {
	rx1 := BiggestInt(SmallestInt(start[0], end[0])-5, 0)
	ry1 := BiggestInt(SmallestInt(start[1], end[1])-5, 0)
	rx2 := BiggestInt(start[0], end[0]) + 5
	ry2 := BiggestInt(start[1], end[1]) + 5
	w := rx2 - rx1 + 1
	h := ry2 - ry1 + 1

	collMat := GetMatrix(w, h, 0)
	objctIDs := wrld.GetObjectsInField(rx1, ry1, w, h)

	for _, objID := range objctIDs {
		hb := wrld.Objects[objID].Hitbox
		collMat.Fill(int(hb.min.X)-rx1, int(hb.min.Y)-ry1, int(math.Ceil(hb.max.X-1))-rx1, int(math.Ceil(hb.max.Y-1))-ry1, 1)
	}

	result, _ := FindPath(start, end, func(pos [2]int) (nbs [][2]int) {
		nbsr := GetNeighboursInBounds(pos[0], pos[1], rx1, ry1, rx2, ry2, collMat)
		nbs = make([][2]int, 0)
		for _, nb := range nbsr {
			X := nb[0] + pos[0]
			Y := nb[1] + pos[1]
			val, err := collMat.Get(X-rx1, Y-ry1)
			if val <= 0 && err == nil {
				nbs = append(nbs, [2]int{X, Y})
			}
		}
		return
	})
	return result
}

//Returns the smallest int
func SmallestInt(iarr ...int) int {
	s := iarr[0]

	for i := 1; i < len(iarr); i++ {
		if iarr[i] < s {
			s = iarr[i]
		}
	}

	return s
}

func BiggestInt(iarr ...int) int {
	s := iarr[0]

	for i := 1; i < len(iarr); i++ {
		if iarr[i] > s {
			s = iarr[i]
		}
	}

	return s
}
