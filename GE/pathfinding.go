package GE

import (
	//"math"
	//"time"
	//"fmt"
)
var PNS = [][2]int{[2]int{0, -1},[2]int{+1, 0},[2]int{0, +1},[2]int{-1, 0}}
var PNDS = [][2]int{[2]int{-1, -1},[2]int{+1, -1},[2]int{+1, +1},[2]int{-1, +1}}
var PNDSS = append(PNS, PNDS...)
func GetNeighboursInBounds(x, y, rx1, ry1, rx2, ry2 int, diagonal bool) (ns [][2]int) {
	pns := append([][2]int{}, PNS...)
	if diagonal {
		pns = append(pns, PNDS...)
	}
	ns = make([][2]int, 0)
	for _,n := range(pns) {
		p := [2]int{n[0]+x, n[1]+y}
		if p[0] >= rx1 && p[1] >= ry1 && p[0] < rx2 && p[1] < ry2 {
			ns = append(ns, n)
		}
	}
	return
}

//Start and end must be greater or equal to 0 and in the bounds of the matrix focus
func FindPathMat(mat *Matrix, start, end [2]int, diagonal bool) (nodes [][2]int) {
	rx1 := int(mat.Focus().Min().X)
	ry1 := int(mat.Focus().Min().Y)
	rx2 := int(mat.Focus().Max().X)
	ry2 := int(mat.Focus().Max().Y)
	result, _ := FindPath(start, end, func(pos [2]int)(nbs [][2]int){
		nbsr := GetNeighboursInBounds(pos[0], pos[1], rx1, ry1, rx2, ry2, diagonal)
		nbs = make([][2]int, 0)
		for _,nb := range(nbsr) {
			X := nb[0]+pos[0]; Y := nb[1]+pos[1]
			val, err := mat.Get(X,Y)
			if val <= 0 && err == nil {
				nbs = append(nbs, [2]int{X,Y})
			}
		}
		return
	})
	return result
}