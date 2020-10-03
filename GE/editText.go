package GE

import (
	"fmt"
	"strings"
	"strconv"
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/golang/freetype/truetype"
)


type EditText struct {
	ImageObj
	text, placeHolderText string
	counter, MaxRunes int
	ttf *truetype.Font
	colors []color.Color
	
	currentColor int
	
	IsSelected, imageNeedsUpdate,Underscore bool
}
func (t *EditText) Print() string {
	return fmt.Sprintf("Text: %s, Placeholder: %s, Counter: %v, MaxRunes: %v", t.text, t.placeHolderText, t.counter, t.MaxRunes)
}

func (t *EditText) Init(screen *ebiten.Image, data interface{}) {}
func (t *EditText) Start(screen *ebiten.Image, data interface{}) {}
func (t *EditText) Stop(screen *ebiten.Image, data interface{}) {}
func (t *EditText) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		xi,yi := ebiten.CursorPosition(); x := float64(xi); y := float64(yi)
		if x > t.X && x < t.X+t.W && y > t.Y && y < t.Y+t.H {
			t.IsSelected = true
		}
	}
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyKPEnter) {
		t.IsSelected = false
		t.imageNeedsUpdate = true
	}
	if t.IsSelected {
		newText := string(ebiten.InputChars())
		if len(newText) > 0 && strings.ContainsAny(newText, allLetters+" ") {
			t.imageNeedsUpdate = true
			t.text += newText
		}
		if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(t.text) >= 1 {
				t.text = t.text[:len(t.text)-1]
				t.imageNeedsUpdate = true
			}
		}
		if t.counter%30 == 0 {
			t.Underscore = !t.Underscore
			t.imageNeedsUpdate = true
		}
		t.counter++
		if t.counter > 30 {
			t.counter -= 30
		}
	}
	if t.imageNeedsUpdate {
		t.UpdateImg()
		t.imageNeedsUpdate = false
	}
}
func (t *EditText) Draw(screen *ebiten.Image) {
	if t.ImageObj.Img != nil {
		t.ImageObj.DrawImageObj(screen)
	}
}

func (t *EditText) UpdateImg() {
	if len(t.text) > t.MaxRunes {
		t.text = t.text[:t.MaxRunes]
	}
	text := t.text
	col := t.colors[0]
	//fmt.Println(len(t.text), ":", t.IsSelected)
	if len(t.text) <= 0 && !t.IsSelected {
		text = t.placeHolderText
		col = t.colors[1]
	}
	if t.IsSelected && t.Underscore {
		text += "_"
	}
	t.ImageObj = *GetTextImage(text, t.X, t.Y, t.H, t.ttf, col, EditText_Back_Col)
}

func (t *EditText) SetText(text string) {
	t.text = text
}
func (t *EditText) GetText() string {
	return t.text
}
func (t *EditText) SetInt(i int) {
	t.SetText(strconv.Itoa(i))
}
func (t *EditText) GetInt() int {
	i, _ := strconv.Atoi(t.GetText())
	return i
}
func (t *EditText) SetUint8(i uint8) {
	t.SetText(strconv.Itoa(int(i)))
}
func (t *EditText) GetUint8() uint8 {
	i, _ := strconv.Atoi(t.GetText())
	return uint8(i)
}