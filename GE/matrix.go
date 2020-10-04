package GE

import (
	"fmt"
)

type Matrix struct {
	X,Y,Z int
	list []int
}
func (m *Matrix) Init(standard int) {
	m.list = make([]int, m.X*m.Y*m.Z)
	for i,_ := range(m.list) {
		m.list[i] = standard
	}
}
func (m *Matrix) InitIdx() {
	m.list = make([]int, m.X*m.Y*m.Z)
	for i,_ := range(m.list) {
		m.list[i] = i
	}
}
func (m *Matrix) Get(x, y, z int) int {
	return m.list[x+m.X*y+m.X*m.Y*z]
}
func (m *Matrix) Set(x, y, z, v int) {
	m.list[x+m.X*y+m.X*m.Y*z] = v
}

func (m *Matrix) SubMatrix(x1,y1,z1,x2,y2,z2 int) *Matrix {
	newMat := &Matrix{x2-x1+1, y2-y1+1, z2-z1+1, nil}
	newMat.Init(0)
	
	for x := x1; x <= x2; x ++ {
		for y := y1; y <= y2; y ++ {
			for z := z1; z <= z2; z ++ {
				newMat.Set(x-x1, y-y1, z-z1, m.Get(x,y,z))
			}
		}
	}
	
	return newMat
}

func (m *Matrix) Print() string {
	out := ""
	for _,i := range(m.list) {
		out += fmt.Sprintf("%v, ", i)
	}
	return out
}