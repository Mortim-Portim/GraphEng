package GE

import (
	"github.com/hajimehoshi/ebiten"
	"strings"
	"image"
)
type TextView struct {
	ImageObj
	text string
	lineHeight, realHeight float64
	scrollIdx float64
	lines, displayLines int
	
	textImg *ebiten.Image
}
func (v *TextView) Reset() {
	v.scrollIdx = 0
}
func (v *TextView) Init(screen *ebiten.Image, data interface{}) {
	v.Reset()
}
func (v *TextView) Start(screen *ebiten.Image, data interface{}) {
	v.Reset()
}
func (v *TextView) Stop(screen *ebiten.Image, data interface{}) {
	v.Reset()
}
func (v *TextView) Update(frame int) {
	x, y := ebiten.CursorPosition()
	if int(v.X) <= x && x < int(v.X+v.W) && int(v.Y) <= y && y < int(v.Y+v.H) {
		_, dy := ebiten.Wheel()
		v.scrollIdx -= dy
		if v.scrollIdx < 0 {
			v.scrollIdx = 0
		}
		if v.scrollIdx >= float64(v.lines-v.displayLines) {
			v.scrollIdx = float64(v.lines-v.displayLines)
		}
	}
}
func (v *TextView) Draw(screen *ebiten.Image, frame int) {
	w,_ := v.textImg.Size()
	yPnt := v.scrollIdx*v.lineHeight+v.lineHeight/4
	img := v.textImg.SubImage(image.Rectangle{image.Point{0, int(yPnt)}, image.Point{w, int(yPnt+v.H)}})
	v.Img = ImgToEbitenImg(&img)
	v.ImageObj.DrawImageObj(screen)
}


func HasLines(text string) int {
	return strings.Count(text, "\n")+1
}

func FormatTextToWidth(text string, maxRunes int, hardBreak bool) (string) {
	if hardBreak {
		return formatTextToWidthHardBreak(text, maxRunes)
	}
	return formatTextToWidthByWords(text, maxRunes)
}

func formatTextToWidthHardBreak(text string, maxRunes int) (formatet string) {
	formatet = ""
	
	currentLength := 0
	for _,r := range(text) {
		if string(r) != "\n" {
			currentLength ++
			if currentLength <= maxRunes {
				formatet += string(r)
			}else{
				formatet += "\n"+string(r)
				currentLength = 1
			}
		}else{
			formatet += "\n"
			currentLength = 0
		}
	}
	
	return
}

func formatTextToWidthByWords(text string, maxRunes int) (formatet string) {
	text = strings.ReplaceAll(text, "\n", " \n")
	words := strings.Split(text, " ")
	newWords := make([]string, 0)
	
	currentLength := 0
	for _,word := range(words) {
		currentLength += len(word)+1
		if currentLength < maxRunes && strings.Index(word, "\n") == -1 {
			newWords = append(newWords, word)
		}else{
			word = strings.ReplaceAll(word, "\n", "")
			if len(word)+1 > maxRunes {
				newWords = append(newWords, "\n"+string(word[:maxRunes]))
				newWords = append(newWords, "\n"+string(word[maxRunes:]))
				currentLength = len(word[maxRunes:])+1
			}else{
				newWords = append(newWords, "\n"+word)
				currentLength = len(word)+1
			}
		}
	}
	formatet = strings.Join(newWords, " ")
	return
}