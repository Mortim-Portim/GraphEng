package GE

import (
	"strings"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	//"image"
)

/**
TextView displays a text of an arbitrary number of lines

The number of lines displayed can be specified, as well as the height of each line

It is recommended to format the text to a maximum number of runes using FormatTextToWidth()

TextView implements UpdateAble
**/

type TextView struct {
	X, Y, W, H                     float64
	text                           string
	lineHeight                     float64
	lines, displayLines, scrollIdx int

	lineImages []*ImageObj
}

//Resets the TextView
func (v *TextView) Reset() {
	v.scrollIdx = 0
}

//Initializes the TextView
func (v *TextView) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	v.Reset()
	return v.Update, v.Draw
}

//Starts the TextView
func (v *TextView) Start(screen *ebiten.Image, data interface{}) {
	v.Reset()
}

//Stops the TextView
func (v *TextView) Stop(screen *ebiten.Image, data interface{}) {}

//Updates the TextView
func (v *TextView) Update(frame int) {
	x, y := ebiten.CursorPosition()
	if int(v.X) <= x && x < int(v.X+v.W) && int(v.Y) <= y && y < int(v.Y+v.H) {
		_, dy := ebiten.Wheel()
		v.scrollIdx -= int(dy)
		if v.scrollIdx < 0 {
			v.scrollIdx = 0
		}
		if v.scrollIdx >= v.lines-v.displayLines {
			v.scrollIdx = v.lines - v.displayLines
		}
	}
}

//Draws the TextView
func (v *TextView) Draw(screen *ebiten.Image) {
	for i := 0; i < v.displayLines; i++ {
		idx := i + v.scrollIdx
		v.lineImages[idx].Y = v.Y + float64(i)*v.lineHeight
		v.lineImages[idx].DrawImageObj(screen)
	}
}

//Returns the number of lines a string has
func HasLines(text string) int {
	return strings.Count(text, "\n") + 1
}

//Formats a string to a specific number of maximum runes per line
func FormatTextToWidth(text string, maxRunes int, hardBreak bool) string {
	if hardBreak {
		return formatTextToWidthHardBreak(text, maxRunes)
	}
	return formatTextToWidthByWords(text, maxRunes)
}

func formatTextToWidthHardBreak(text string, maxRunes int) (formatet string) {
	formatet = ""

	currentLength := 0
	for _, r := range text {
		if string(r) != "\n" {
			currentLength++
			if currentLength <= maxRunes {
				formatet += string(r)
			} else {
				formatet += "\n" + string(r)
				currentLength = 1
			}
		} else {
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
	for _, word := range words {
		currentLength += len(word) + 1
		if currentLength < maxRunes && strings.Index(word, "\n") == -1 {
			newWords = append(newWords, word)
		} else {
			word = strings.ReplaceAll(word, "\n", "")
			if len(word)+1 > maxRunes {
				newWords = append(newWords, "\n"+string(word[:maxRunes]))
				newWords = append(newWords, "\n"+string(word[maxRunes:]))
				currentLength = len(word[maxRunes:]) + 1
			} else {
				newWords = append(newWords, "\n"+word)
				currentLength = len(word) + 1
			}
		}
	}
	formatet = strings.Join(newWords, " ")
	return
}
