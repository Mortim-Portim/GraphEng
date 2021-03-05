package GE

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"

	cmp "github.com/mortim-portim/GraphEng/compression"
)

var ERROR_COORDS_NOT_ON_MATRIX = errors.New("Coordinates not on matrix")

/**
Matrix is a struct that can store a 2 dimensional array of a maximum width and height of 65535
the stored values are 16-bit integers

A matrix can have a focus, that lets all functions (except WAbs, HAbs, GetAbs) access only values and indexes withhin that focus
If a value is outside the matrix Get will return -1

Example:
whole Matrix:
 10  8  8  6  6  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  8  8  8  6  6  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  8  8  8  6  6  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  6  6  6  6  6  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  6  6  6  6  6  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  4  4  4  4  4  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  4  4  4  4  4  4  4  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  2  2  2  2  2  2  2  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  2  2  2  2  2  2  2  2  2  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  0  0  0  0  0  0  0  0  0  0  0 -2 -2 -4 -4 -6 -6  0  0  0
  0  0  0  0  0  0  0  0  0  0  0 -2 -2 -4 -4 -6 -6  0  0  0
 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -4 -4 -6 -6  0  0  0
 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2 -4 -4 -6 -6  0  0  0
 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -6 -6  0  0  0
 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -4 -6 -6  0  0  0
 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6  0  0  0
 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6 -6  0  0  0
  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0

subMatrix with focus (3,2,15,12)
  6  6  4  4  2  2  0  0 -2 -2
  6  6  4  4  2  2  0  0 -2 -2
  6  6  4  4  2  2  0  0 -2 -2
  4  4  4  4  2  2  0  0 -2 -2
  4  4  4  4  2  2  0  0 -2 -2
  2  2  2  2  2  2  0  0 -2 -2
  2  2  2  2  2  2  0  0 -2 -2
  0  0  0  0  0  0  0  0 -2 -2
  0  0  0  0  0  0  0  0 -2 -2
 -2 -2 -2 -2 -2 -2 -2 -2 -2 -2
**/

//Returns a Matrix of width=x, height=y and initial value=v
func GetMatrix(x, y int, v int64) (m *Matrix) {
	m = &Matrix{}
	m.x = uint64(x)
	m.y = uint64(y)
	m.Init(v)
	m.ResetFocus()
	return
}
func GetMatrixFromImage(img *image.Gray) (m *Matrix) {
	m = &Matrix{}
	m.x = uint64(img.Bounds().Dx())
	m.y = uint64(img.Bounds().Dy())
	m.ResetFocus()
	m.FillFromImage(img)
	return
}

type Matrix struct {
	x, y  uint64
	list  []int64
	focus *Rectangle
}

func (m *Matrix) ScaleTo(x, y int, v int64) {
	nL := make([]int64, (m.W()-x)*(m.H()-y))
	for i := range nL {
		nL[i] = v
	}
	m.x = uint64(x)
	m.y = uint64(y)
	m.list = append(m.list, nL...)
}

//Returns the width of the focused Matrix
func (m *Matrix) W() int {
	return int(m.focus.Bounds().X)
}

//Returns the height of the focused Matrix
func (m *Matrix) H() int {
	return int(m.focus.Bounds().Y)
}

//Returns the absolute width of the Matrix
func (m *Matrix) WAbs() int {
	return int(m.x)
}

//Returns the absolute height of the Matrix
func (m *Matrix) HAbs() int {
	return int(m.y)
}

//Initializes m with a certain value
func (m *Matrix) Init(standard int64) {
	m.list = make([]int64, m.WAbs()*m.HAbs())
	for i := range m.list {
		m.list[i] = standard
	}
}

//Initializes m with the indexes (used for debugging)
func (m *Matrix) InitIdx() {
	m.list = make([]int64, m.WAbs()*m.HAbs())
	for i := range m.list {
		m.list[i] = int64(i)
	}
}

//Returns a matrix with the same fields
func (m *Matrix) Clone() *Matrix {
	return &Matrix{m.x, m.y, m.list, m.focus}
}

//Returns the value of the matrix at the absolute x and y coordinates
func (m *Matrix) GetAbs(x, y int) (int64, error) {
	if x < 0 || y < 0 || x >= m.WAbs() || y >= m.HAbs() {
		return 0, ERROR_COORDS_NOT_ON_MATRIX
	}
	return m.list[x+m.WAbs()*y], nil
}

//Sets the value of the matrix at the x and y coordinates
func (m *Matrix) SetAbs(x, y int, v int64) {
	if x < 0 || y < 0 || x >= m.WAbs() || y >= m.HAbs() {
		return
	}
	m.list[x+m.WAbs()*y] = v
}
func (m *Matrix) XYtoIDX(x, y int) (int, error) {
	if x >= m.WAbs() || y >= m.HAbs() {
		return 0, ERROR_COORDS_NOT_ON_MATRIX
	}
	return x + m.WAbs()*y, nil
}
func (m *Matrix) IDXtoXY(idx int) (int, int, error) {
	csm1 := m.WAbs() - 1
	x := idx % m.WAbs()
	y := (idx - (idx % csm1)) / csm1
	if x >= m.WAbs() || y >= m.HAbs() {
		return 0, 0, ERROR_COORDS_NOT_ON_MATRIX
	}
	return x, y, nil
}

//Returns the value of the matrix at the index
func (m *Matrix) GetAbsByIdx(idx int) (int64, error) {
	if idx < 0 || idx >= len(m.list) {
		return 0, fmt.Errorf("Index: %v not on matrix with w:%v, h:%v, wAbs:%v, hAbs:%v, lx:%v, ly:%v", idx, int(m.focus.Bounds().X), int(m.focus.Bounds().Y), m.WAbs(), m.HAbs(), int(m.focus.Min().X), int(m.focus.Min().Y))
	}
	return m.list[idx], nil
}

//Returns the value of the focused matrix at the x and y coordinates
func (m *Matrix) Get(x, y int) (int64, error) {
	xl, yl := int(m.focus.Min().X)+x, int(m.focus.Min().Y)+y
	if xl < 0 || xl >= m.WAbs() || yl < 0 || yl >= m.HAbs() {
		return 0, ERROR_COORDS_NOT_ON_MATRIX //fmt.Errorf("Coordinates (%v:%v) not on matrix with w:%v, h:%v, xl:%v, yl:%v, wAbs:%v, hAbs:%v, lx:%v, ly:%v", x, y, int(m.focus.Bounds().X), int(m.focus.Bounds().Y), xl, yl, m.WAbs(), m.HAbs(), int(m.focus.Min().X), int(m.focus.Min().Y))
	}
	idx := xl + m.WAbs()*yl
	if idx >= len(m.list) {
		panic(fmt.Sprintf("Len list: %v Index: %v with Coordinates (%v:%v) not on matrix with w:%v, h:%v, xl:%v, yl:%v, wAbs:%v, hAbs:%v, lx:%v, ly:%v", len(m.list), idx, x, y, int(m.focus.Bounds().X), int(m.focus.Bounds().Y), xl, yl, m.WAbs(), m.HAbs(), int(m.focus.Min().X), int(m.focus.Min().Y)))
	}
	return m.list[idx], nil
}

//Returns the nearest value of a border of the matrix to a point
func (m *Matrix) GetNearest(x, y int) (int64, error) {
	v, err := m.Get(x, y)
	if err == nil {
		return v, nil
	}
	if x >= 0 && x < m.W() && y >= 0 && y < m.H() {
		return 0, nil
	}
	X, Y, W, H := float64(x), float64(y), float64(m.W()), float64(m.H())

	//TopLeft
	if X < 0 && Y < 0 {
		return m.Get(0, 0)
	}
	//TopRight
	if X >= W && Y < 0 {
		return m.Get(m.W()-1, 0)
	}
	//BottomLeft
	if X < 0 && Y >= H {
		return m.Get(0, m.H()-1)
	}
	//BottomRight
	if X >= W && Y >= H {
		return m.Get(m.W()-1, m.H()-1)
	}

	//Left
	if X < 0 {
		return m.Get(0, y)
	}
	//Right
	if X >= W {
		return m.Get(m.W()-1, y)
	}
	//Top
	if Y < 0 {
		return m.Get(x, 0)
	}
	//Bottom
	if Y >= H {
		return m.Get(x, m.H()-1)
	}
	return 0, fmt.Errorf("Should not be reached: x:%v, y:%v, w:%v, h:%v", X, Y, W, H)
}

//Sets the value of the focused matrix at the x and y coordinates
func (m *Matrix) Set(x, y int, v int64) {
	xl, yl := int(m.focus.Min().X)+x, int(m.focus.Min().Y)+y
	if xl < 0 || xl >= m.WAbs() || yl < 0 || yl >= m.HAbs() {
		return
	}
	m.list[xl+m.WAbs()*yl] = v
}

//Adds a value to the value of the focused matrix at the x and y coordinate
func (m *Matrix) Add(x, y int, v int64) {
	ov, _ := m.Get(x, y)
	m.Set(x, y, int64(ov)+v)
}

//Adds a value to the value of the focused matrix at the x and y coordinate
func (m *Matrix) AddAbs(x, y int, v int64) {
	ov, _ := m.GetAbs(x, y)
	m.SetAbs(x, y, int64(ov)+v)
}

//Swaps the value of two positions
func (m *Matrix) Swap(x1, y1, x2, y2 int) error {
	p1v, err := m.Get(x1, y1)
	if err != nil {
		return err
	}
	p2v, err := m.Get(x2, y2)
	if err != nil {
		return err
	}
	m.Set(x1, y1, p2v)
	m.Set(x2, y2, p1v)
	return nil
}

//Swaps the value of two positions
func (m *Matrix) SwapAbs(x1, y1, x2, y2 int) error {
	p1v, err := m.GetAbs(x1, y1)
	if err != nil {
		return err
	}
	p2v, err := m.GetAbs(x2, y2)
	if err != nil {
		return err
	}
	m.SetAbs(x1, y1, p2v)
	m.SetAbs(x2, y2, p1v)
	return nil
}

//Clears the matrix
func (m *Matrix) Clear(v int64) {
	m.FillAbs(0, 0, m.WAbs()-1, m.HAbs()-1, v)
}

//Fills a Rectangle with a value
func (m *Matrix) Fill(x1, y1, x2, y2 int, v int64) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			m.Set(x, y, v)
		}
	}
}
func (m *Matrix) FillAll(v int64) {
	for i := range m.list {
		m.list[i] = v
	}
}
func (m *Matrix) FillFromImage(img *image.Gray) {
	for x := 0; x < int(m.x); x++ {
		for y := 0; y < int(m.y); y++ {
			if (image.Point{x, y}.In(img.Rect)) {
				m.SetAbs(x, y, int64(img.GrayAt(x, y).Y))
			}
		}
	}
}

//Fills a Rectangle with a value
func (m *Matrix) FillAbs(x1, y1, x2, y2 int, v int64) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			m.SetAbs(x, y, v)
		}
	}
}

//Adds a value to a Rectangle
func (m *Matrix) AddToRect(x1, y1, x2, y2 int, v int64) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			m.Add(x, y, v)
		}
	}
}

//Adds a value to all focused values
func (m *Matrix) AddToAll(v int64) {
	for x := 0; x < m.W(); x++ {
		for y := 0; y < m.H(); y++ {
			m.Add(x, y, v)
		}
	}
}

//Adds a value to all values
func (m *Matrix) AddToAllAbs(v int64) {
	for x := 0; x < m.WAbs(); x++ {
		for y := 0; y < m.HAbs(); y++ {
			m.Add(x, y, v)
		}
	}
}

//Resets the focus (0,0, m.x, m.y)
func (m *Matrix) ResetFocus() {
	m.focus = GetRectangle(0, 0, float64(m.WAbs()), float64(m.HAbs()))
}

//Sets the focus of the matrix
func (m *Matrix) SetFocus(x1, y1, x2, y2 int) {
	//	if x1 < 0 {
	//		x1 = 0
	//	}
	//	if y1 < 0 {
	//		y1 = 0
	//	}
	//	if x2 > int(m.x.Int64()) {
	//		x2 = int(m.x.Int64())
	//	}
	//	if y2 > int(m.y.Int64()) {
	//		y2 = int(m.y.Int64())
	//	}
	m.focus = GetRectangle(float64(x1), float64(y1), float64(x2), float64(y2))
}

//Sets the focus of the matrix
func (m *Matrix) SetFocusRect(r *Rectangle) {
	m.focus = r
}

//Copys the focus of the matrix
func (m *Matrix) CopyFocus(m2 *Matrix) {
	//m2.SetFocus(int(m.focus.Min().X), int(m2.focus.Min().Y), int(m2.focus.Max().X), int(m2.focus.Max().Y))
	m2.focus = m.focus.Copy()
}

//Returns the focus of the matrix
func (m *Matrix) Focus() *Rectangle {
	return m.focus
}

//Returns a copy of the matrix focused on a specific rectangle
func (m *Matrix) SubMatrix(x1, y1, x2, y2 int) (newM *Matrix) {
	newM = m.Clone()
	newM.SetFocus(x1, y1, x2, y2)
	return
}
func (m *Matrix) SubMatFromMiddle(x, y, r int) (newM *Matrix) {
	newM = m.Clone()
	newM.SetFocus(x-r, y-r, x+r, y+r)
	return
}

//Prints a matrix with maximum values of 999
func (m *Matrix) Print() string {
	out := ""
	for y := 0; y < m.H(); y++ {
		for x := 0; x < m.W(); x++ {
			v, _ := m.Get(x, y)
			valStr := fmt.Sprintf("%v", v)
			for i := 0; i < 4-len(valStr); i++ {
				out += " "
			}
			out += valStr
		}
		out += "\n\n"
	}
	out += fmt.Sprintf("W:%v, H:%v, Focus:%s", m.W(), m.H(), m.focus.Print())
	return out
}

//Converts a Matrix to a []byte slice
func (m *Matrix) ToBytes() []byte {
	b := cmp.Int64sToBytes(m.list...)
	b = append(b, cmp.Int16sToBytes(int16(m.focus.Min().X), int16(m.focus.Min().Y), int16(m.focus.Max().X), int16(m.focus.Max().Y))...)
	b = append(b, cmp.UInt64ToBytes(m.x)...)
	b = append(b, cmp.UInt64ToBytes(m.y)...)
	return b
}

//Loads a Matrix from a []byte slice
func (m *Matrix) FromBytes(bs []byte) {
	m.list = cmp.BytesToInt64s(bs[:len(bs)-24])
	is := cmp.BytesToInt16s(bs[len(bs)-24 : len(bs)-16])
	m.x = cmp.BytesToUInt64(bs[len(bs)-16 : len(bs)-8])
	m.y = cmp.BytesToUInt64(bs[len(bs)-8:])
	m.focus = GetRectangle(float64(is[0]), float64(is[1]), float64(is[2]), float64(is[3]))

	//Ensure backcompatibility
	if len(m.list) != m.WAbs()*m.HAbs() {
		int16L := cmp.BytesToInt16s(bs[:len(bs)-24])
		m.list = make([]int64, len(int16L))
		for i, i16 := range int16L {
			m.list[i] = int64(i16)
		}
	}
}

//Compresses a Matrix to a []byte slice
func (m *Matrix) Compress() ([]byte, error) {
	return cmp.CompressBytes(m.ToBytes())
}

//Decompresses a []byte slice, that was compressed by m.Compress()
func (m *Matrix) Decompress(bs []byte) error {
	content, err := cmp.DecompressBytes(bs)
	if err != nil {
		return err
	}
	m.FromBytes(content)
	return nil
}

//Loads a compressed matrix from the file system
func (m *Matrix) Load(path string) error {
	dat, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err2
	}
	m.Decompress(dat)
	return nil
}

//Saves a matrix in compressed form to the file system
func (m *Matrix) Save(path string) error {
	bs, err := m.Compress()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, 0644)
}

//Loads an uncompressed matrix from the file system
func (m *Matrix) LoadUnCompressed(path string) error {
	dat, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err2
	}
	m.FromBytes(dat)
	return nil
}

//Saves a matrix in uncompressed form to the file system
func (m *Matrix) SaveUnCompressed(path string) error {
	return ioutil.WriteFile(path, m.ToBytes(), 0644)
}
