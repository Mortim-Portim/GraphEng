package GE

import (
	"fmt"
)


func GetMatrix(x,y int,v int8) (m *Matrix) {
	m = &Matrix{x:x,y:y}
	m.Init(v)
	m.ResetFocus()
	return
}
type Matrix struct {
	x,y int
	list []int8
	focus *Rectangle
}
func (m *Matrix) W() int {
	return int(m.focus.Bounds().X)
}
func (m *Matrix) H() int {
	return int(m.focus.Bounds().Y)
}
func (m *Matrix) Init(standard int8) {
	m.list = make([]int8, m.x*m.y)
	for i,_ := range(m.list) {
		m.list[i] = standard
	}
}
func (m *Matrix) InitIdx() {
	m.list = make([]int8, m.x*m.y)
	for i,_ := range(m.list) {
		m.list[i] = int8(i)
	}
}
func (m *Matrix) Clone() *Matrix {
	return &Matrix{m.x,m.y,m.list,m.focus}
}
func (m *Matrix) GetAbs(x,y int) int {
	idx := x+m.x*y
	if idx < 0 || idx >= int(len(m.list)) {
		return -1
	}
	return int(m.list[idx])
}
func (m *Matrix) Get(x, y int) int {
	idx := (x+int(m.focus.Min().X))+m.x*(y+int(m.focus.Min().Y))
	if idx < 0 || idx >= int(len(m.list)) {
		return -1
	}
	return int(m.list[idx])
}
func (m *Matrix) Set(x, y int, v int8) {
	m.list[(x+int(m.focus.Min().X))+m.x*(y+int(m.focus.Min().Y))] = v
}
func (m *Matrix) Add(x,y int, v int8) {
	m.Set(x,y, int8(m.Get(x,y))+v)
}
func (m *Matrix) Fill(x1,y1,x2,y2 int, v int8) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			m.Set(x,y,v)
		}
	}
}
func (m *Matrix) AddToAll(v int8) {
	for x := 0; x < m.W(); x++ {
		for y := 0; y < m.H(); y++ {
			m.Add(x,y,v)
		}
	}
}
func (m *Matrix) ResetFocus() {
	m.focus = GetRectangle(0,0,float64(m.x),float64(m.y))
}
func (m *Matrix) SetFocus(x1,y1,x2,y2 int) {
	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 > m.x {
		x2 = m.x
	}
	if y2 > m.y {
		y2 = m.y
	}
	m.focus = GetRectangle(float64(x1),float64(y1), float64(x2), float64(y2))
}
func (m *Matrix) Focus() *Rectangle {
	return m.focus
}
func (m *Matrix) SubMatrix(x1,y1,x2,y2 int) (newM *Matrix) {
	newM = m.Clone()
	newM.SetFocus(x1,y1,x2,y2)
	return
}

func (m *Matrix) Print() string {
	out := ""
	for x := 0; x < m.W(); x++ {
		for y := 0; y < m.H(); y++ {
			valStr := fmt.Sprintf("%v", m.Get(x,y))
			for i := 0; i < 3-len(valStr); i++ {
				out += " "
			}
			out += valStr
		}
		out += "\n"
	}
	return out
}