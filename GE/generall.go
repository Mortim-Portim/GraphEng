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
	"marvin/GraphEng/GC"
	"strings"
	//"fmt"
	"math"
)
const StandardFontSize = 64

type UpdateFunc func(screen *ebiten.Image, frame int)

type UpdateAble interface {
	Init(screen *ebiten.Image, data interface{})		//Called once at the beginning
	Start(screen *ebiten.Image, data interface{}) 		//Called whenever the View becomes visible
	Stop(screen *ebiten.Image, data interface{})		//Called whenever the View becomes invisible
	Update()									//Called every Frame if the View is visible
	Draw(screen *ebiten.Image)				//Called every Frame if the View is visible
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
	w, h := MeasureString(str, mplusNormalFont)
	
	popUpBack, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	popUpBack.Fill(backCol)
	xP, yP := h/6, h/4*3
	text.Draw(popUpBack, str, mplusNormalFont, int(xP), int(yP), textCol)
	return popUpBack
}

//Returns an ImageObj with text on it
func GetTextImage(textStr string, X, Y, H float64, ttf *truetype.Font, txtCol, backCol color.Color) (*ImageObj) {
	imgo := &ImageObj{H:H, X:X, Y:Y}
	if len(textStr) > 0 {
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
	}else{
		Back, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
		Back.Fill(backCol); imgo.Img = Back
	}
	
	return imgo
}

func GetTextLinesImages(textStr string, X, Y, lineHeight float64, ttf *truetype.Font, txtCol, backCol color.Color) (lineImgs []*ImageObj, maxWidth float64) {
	lines := strings.Split(textStr, "\n")
	lineImgs = make([]*ImageObj, len(lines))
	maxWidth = 0
	for i,str := range(lines) {
		for str[0] == " "[0] {
			str = str[1:]
		}
		lineImgs[i] = GetTextImage(str, X, Y+float64(i)*lineHeight, lineHeight, ttf, txtCol, backCol)
		if lineImgs[i].W > maxWidth {
			maxWidth = lineImgs[i].W
		}
	}
	return 
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

const allLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var faceHeight = make(map[font.Face]int)
func MeasureString(str string, faceTTF font.Face) (x, y int) {
	h, ok := faceHeight[faceTTF]
	if !ok {
		rectAll := text.BoundString(faceTTF, allLetters)
		h = rectAll.Max.Y-rectAll.Min.Y
		faceHeight[faceTTF] = h
	}
	rect := text.BoundString(faceTTF, str+"#")
	x, y = rect.Max.X-rect.Min.X, h*(strings.Count(str, "\n")+1)+h/6
	//fmt.Println(rect.String(), ":     ", x, ":", y)
	return 
}

func genVertices(X,Y,R float64, num int) *Points {
	centerX := X
	centerY := Y
	r       := R

	vs := make([]*GC.Vector,0)
	for i := 0; i <= num; i++ {
		rate := float64(i) / float64(num)
		vs = append(vs, &GC.Vector{
			X:float64(r*math.Cos(2*math.Pi*rate)) + centerX,
			Y:float64(r*math.Sin(2*math.Pi*rate)) + centerY,
			Z:0})
	}

	vs = append(vs, &GC.Vector{
		X:centerX,
		Y:centerY,
		Z:0})
	ps := Points(vs)
	return &ps
}