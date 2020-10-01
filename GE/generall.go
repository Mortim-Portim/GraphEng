package GE

import (
	//"os"
	//"image"
	//_ "image/jpeg"
	"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/font"
	//"github.com/nfnt/resize"
	"image/color"
	//"fmt"
	//"math"
)
const StandardFontSize = 64

type UpdateFunc func(screen *ebiten.Image, frame int)

type UpdateAble interface {
	Init(screen *ebiten.Image, data interface{})		//Called once at the beginning
	Start(screen *ebiten.Image, data interface{}) 		//Called whenever the View becomes visible
	Stop(screen *ebiten.Image, data interface{})		//Called whenever the View becomes invisible
	Update(frame int)									//Called every Frame if the View is visible
	Draw(screen *ebiten.Image, frame int)				//Called every Frame if the View is visible
}

var StandardFont *truetype.Font

//Initializes the Graphics Engine and the Standard Font (use "" to load standard font)
func Init(FontPath string) {
	if len(FontPath) > 0 {
		font, err := ioutil.ReadFile(FontPath)
	   	if err != nil {
	   		panic(err)
	   	}
		tt, err := truetype.Parse(font)
		CheckErr(err)
		StandardFont = tt
	}else{
		//fonts.ArcadeN_ttf
		tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
		CheckErr(err)
		StandardFont = tt
	}
	
	InitParams()
}

//Returns an Image with text
func MakePopUp(str string, size float64, ttf *truetype.Font, textCol, backCol color.Color) (*ebiten.Image) {
	mplusNormalFont := truetype.NewFace(ttf, &truetype.Options{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(str, mplusNormalFont)
	
	w, h := int(float64(pnt.X)), int(float64(pnt.Y))
	popUpBack, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	popUpBack.Fill(backCol)
	xP, yP := 0, pnt.Y/4*3
	text.Draw(popUpBack, str, mplusNormalFont, int(xP), int(yP), textCol)
	return popUpBack
}

//Returns an ImageObj with text on it
func GetTextImage(textStr string, X, Y, H float64, ttf *truetype.Font, txtCol, backCol color.Color) (*ImageObj) {
	imgo := &ImageObj{H:H, X:X, Y:Y}
	
	textImg := MakePopUp(textStr, StandardFontSize, ttf, txtCol, &color.RGBA{0,0,0,0})
	w,h := textImg.Size(); W := float64(w)*H/float64(h)
	imgo.W = W
	
	Back, _ := ebiten.NewImage(int(W), int(H), ebiten.FilterDefault)
	Back.Fill(backCol); imgo.Img = Back
	
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Scale(H/float64(h), H/float64(h))
	//op.GeoM.Translate(H*0.25,0)
	imgo.Img.DrawImage(textImg, op)
	
	return imgo
}

func GetTextLinesImage(textStr string, X, Y, H float64, ttf *truetype.Font, txtCol, backCol color.Color) (*ebiten.Image) {
	lines := HasLines(textStr)
	mplusNormalFont := truetype.NewFace(ttf, &truetype.Options{
		Size:    StandardFontSize,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	pnt := text.MeasureString(textStr, mplusNormalFont)
	
	w, h := pnt.X, pnt.Y
	Back, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	xP, yP := 0, float64(h)/float64(lines)
	text.Draw(Back, textStr, mplusNormalFont, int(xP), int(yP), txtCol)
	
	W := float64(w)*H/float64(h)
	smallImg, _ := ebiten.NewImage(int(W), int(H), ebiten.FilterDefault)
	smallImg.Fill(backCol)
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Scale(H/float64(h), H/float64(h))
	smallImg.DrawImage(Back, op)
	return smallImg
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

//Panics if an error occured
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SameCols(col, col2 color.Color) bool {
	r,g,b,a := col.RGBA()
	r2,g2,b2,a2 := col2.RGBA()
	if r == r2 && g == g2 && b == b2 && a == a2 {
		return true
	}
	return false
}

func ReduceColor(col color.Color, delta int) color.Color {
	r,g,b,a := col.RGBA()
	newR := int(r)-delta
	if newR < 0 {
		newR = 0
	}
	newG := int(g)-delta
	if newG < 0 {
		newG = 0
	}
	newB := int(b)-delta
	if newB < 0 {
		newB = 0
	}
	return &color.RGBA{uint8(newR),uint8(newG),uint8(newB),uint8(a)}
}

func ReduceColorImage(img *ebiten.Image, val int) (reduced *ebiten.Image) {
	W, H := img.Size()
	Back, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	reduced = Back
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1,1,1, float64(255-val)/255.0)
	reduced.DrawImage(img, op)
	return
}

//Returns true if e is in s
func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}