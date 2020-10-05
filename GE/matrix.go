package GE

import (
	"fmt"
	"io/ioutil"
)

func GetMatrix(x,y,v int16) (m *Matrix) {
	m = &Matrix{x:x,y:y}
	m.Init(v)
	m.ResetFocus()
	return
}
type Matrix struct {
	x,y int16
	list []int16
	focus *Rectangle
}
func (m *Matrix) W() int16 {
	return int16(m.focus.Bounds().X)
}
func (m *Matrix) H() int16 {
	return int16(m.focus.Bounds().Y)
}
func (m *Matrix) Init(standard int16) {
	m.list = make([]int16, m.x*m.y)
	for i,_ := range(m.list) {
		m.list[i] = standard
	}
}
func (m *Matrix) InitIdx() {
	m.list = make([]int16, m.x*m.y)
	for i,_ := range(m.list) {
		m.list[i] = int16(i)
	}
}
func (m *Matrix) Clone() *Matrix {
	return &Matrix{m.x,m.y,m.list,m.focus}
}
func (m *Matrix) GetAbs(x,y int16) int {
	idx := x+m.x*y
	if idx < 0 || idx >= int16(len(m.list)) {
		return -1
	}
	return int(m.list[idx])
}
func (m *Matrix) Get(x, y int16) int {
	idx := (x+int16(m.focus.Min().X))+m.x*(y+int16(m.focus.Min().Y))
	if idx < 0 || idx >= int16(len(m.list)) {
		return -1
	}
	return int(m.list[idx])
}
func (m *Matrix) Set(x, y, v int16) {
	m.list[(x+int16(m.focus.Min().X))+m.x*(y+int16(m.focus.Min().Y))] = v
}
func (m *Matrix) Add(x,y, v int16) {
	m.Set(x,y, int16(m.Get(x,y))+v)
}
func (m *Matrix) Fill(x1,y1,x2,y2, v int16) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			m.Set(x,y,v)
		}
	}
}
func (m *Matrix) AddToAll(v int16) {
	for x := int16(0); x < m.W(); x++ {
		for y := int16(0); y < m.H(); y++ {
			m.Add(x,y,v)
		}
	}
}
func (m *Matrix) ResetFocus() {
	m.focus = GetRectangle(0,0,float64(m.x),float64(m.y))
}
func (m *Matrix) SetFocus(x1,y1,x2,y2 int16) {
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
func (m *Matrix) SubMatrix(x1,y1,x2,y2 int16) (newM *Matrix) {
	newM = m.Clone()
	newM.SetFocus(x1,y1,x2,y2)
	return
}

func (m *Matrix) Print() string {
	out := ""
	for x := int16(0); x < m.W(); x++ {
		for y := int16(0); y < m.H(); y++ {
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

func (m *Matrix) ToBytes() []byte {
	b := Int16sToBytes(m.list)
	b = AppendInt16ToBytes(m.x, b)
	b = AppendInt16ToBytes(m.y, b)
	
	b = AppendInt16ToBytes(int16(m.focus.Min().X), b)
	b = AppendInt16ToBytes(int16(m.focus.Min().Y), b)
	b = AppendInt16ToBytes(int16(m.focus.Max().X), b)
	b = AppendInt16ToBytes(int16(m.focus.Max().Y), b)
	return b
}
func (m *Matrix) FromBytes(bs []byte) {
	is := BytesToInt16s(bs)
	m.list = is[:len(is)-6]
	m.x = is[len(is)-6]
	m.y = is[len(is)-5]
	m.focus = GetRectangle(float64(is[len(is)-4]), float64(is[len(is)-3]), float64(is[len(is)-2]), float64(is[len(is)-1]))
}

func (m *Matrix) Compress() ([]byte, error) {
	return CompressBytes(m.ToBytes())
}

func (m *Matrix) Decompress(bs []byte) error {
	content, err := DecompressBytes(bs)
	if err != nil {
		return err
	}
	m.FromBytes(content)
	return nil
}
func (m *Matrix) Load(path string) error {
	dat, err2 := ioutil.ReadFile(path)
   	if err2 != nil {
   		return err2
   	}
   	m.Decompress(dat)
	return nil
}
func (m *Matrix) Save(path string) error {
	bs, err := m.Compress()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, 0644)
}
func (m *Matrix) LoadUnCompressed(path string) error {
	dat, err2 := ioutil.ReadFile(path)
   	if err2 != nil {
   		return err2
   	}
   	m.FromBytes(dat)
	return nil
}
func (m *Matrix) SaveUnCompressed(path string) error {
	return ioutil.WriteFile(path, m.ToBytes(), 0644)
}