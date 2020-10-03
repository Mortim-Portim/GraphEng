package GE

import (
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/golang/freetype/truetype"
	//"golang.org/x/image/font"
	"image/color"
)

//BUTTONS -------------------------------------------------------------------------------------------------------------------------------

func GetButton(img *ImageObj, dark *ebiten.Image, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	b := &Button{}
	b.Img = img; b.dark = dark; b.light = img.Img; b.onPressLeft = onPressLeft; b.onPressRight = onPressRight; b.Active = true
	return b
}
func GetTextButton(str, downStr string, ttf *truetype.Font, X, Y, H float64, textCol, backCol color.Color, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	img := GetTextImage(str, X, Y, H, ttf, textCol, backCol)
	dark := GetTextImage(downStr, X, Y, H, ttf, textCol, ReduceColor(backCol, ReduceColOnButtonDown))
	return GetButton(img, dark.Img, onPressLeft, onPressRight)
}
func GetImageButton(path [2]string, X,Y, W,H float64, onPressLeft func(b *Button), onPressRight func(b *Button)) *Button {
	img := LoadImgObj(path[0], W, H, X, Y, 0)
	dark := ReduceColorImage(img.Img, ReduceColOnButtonDown)
	return GetButton(img, dark, onPressLeft, onPressRight)
}
//EDITTEXT ------------------------------------------------------------------------------------------------------------------------------

func GetEditText(placeHolderText string, X, Y, H float64, maxRunes int, ttf *truetype.Font, cols ...color.Color) (et *EditText) {
	imgo := ImageObj{H:H, X:X, Y:Y}
	et = &EditText{imgo, "", placeHolderText, 0, maxRunes, ttf, cols, 0, false, true, true}
	return
}
//TEXTVIEW ------------------------------------------------------------------------------------------------------------------------------

func GetTextView(text string, X, Y, H, lineHeight float64, ttf *truetype.Font, txtCol, backCol color.Color) *TextView {
	v := &TextView{text:text}; v.X = X; v.Y = Y; v.lineHeight = lineHeight
	v.lines = HasLines(text); v.realHeight = lineHeight*float64(v.lines)
	v.textImg = GetTextLinesImage(text, X, Y, v.realHeight, ttf, txtCol, backCol)
	w,h := v.textImg.Size()
	v.W = float64(w)
	v.H = H
	if float64(h) < v.H {
		v.H = float64(h)
	}
	v.displayLines = int(v.H/v.lineHeight)
	//fmt.Println("Width:",w,", Height:",h,", realHeight:",v.realHeight,", lines",v.lines,", displayLines:",v.displayLines)
	return v
}
//TABVIEW -------------------------------------------------------------------------------------------------------------------------------
type TabViewParams struct {
	Back, Text color.Color
	TTF *truetype.Font
	
	Dis float64
	X, Y, W, H, TabH float64
	Curr int
	Scrs []UpdateAble
	
	Nms, Pths []string
}
func (p *TabViewParams) fillDefault() {
	if p.Back == nil {
		p.Back = TabBack_Col
	}
	if p.Text == nil {
		p.Text = TabText_Col
	}
	if p.Dis == 0 {
		p.Dis = TabsDistance
	}
	if p.TabH == 0 && p.H != 0 {
		p.TabH = p.H*TabsHeight
	}
	if p.TTF == nil {
		p.TTF = StandardFont
	}
}
func GetTabView(p *TabViewParams) (*TabView) {
	p.fillDefault()
	if p.Pths != nil {
		return getTabViewWithImages(p.Pths, p.Scrs, p.X, p.Y, p.W, p.H, p.TabH, p.Dis, p.Curr)
	}
	return getTabView(p.Nms, p.Scrs, p.X, p.Y, p.W, p.H, p.TabH, p.TTF, p.Text, p.Back, p.Dis, p.Curr)
}
//SCROLLBAR ------------------------------------------------------------------------------------------------------------------------------
func GetImageScrollbar(X, Y, W, H float64, bar, pointer *ebiten.Image, min, max, current int, ttf *truetype.Font) (b *ScrollBar) {
	b = &ScrollBar{min: min, max: max, current: current, ttf:ttf}
	b.Img = bar
	b.X = X
	b.Y = Y
	b.W = W
	b.H = H
	b.pointer = &ImageObj{Y:Y-H/2, W: H*2, H: H*2}
	b.value = &ImageObj{Y:Y-H/2, W: H*2, H: H*2}
	b.pointer.Img = pointer
	b.length = max-min
	b.stepsize = W / float64(b.length)
	b.UpdatePos()
	return
}
func GetStandardScrollbar(X, Y, W, H float64, min, max, current int, ttf *truetype.Font) (b *ScrollBar) {
	
	bar, _ := ebiten.NewImage(int(W), int(H), ebiten.FilterDefault)
	bar.Fill(&color.RGBA{0,0,255,255})
	line := GetLineOfPoints(0,H/2, W,H/2, H/6)
	line.Fill(bar, &color.RGBA{0,200,50,255})
	
	pointer, _ := ebiten.NewImage(int(H*2), int(H*2), ebiten.FilterDefault)
	pointer.Fill(&color.RGBA{0,255,0,100})
	pnts := genVertices(H,H,H, 11)
	pnts.Fill(pointer, &color.RGBA{200,0,50,255})
	
	return GetImageScrollbar(X,Y,W,H,bar,pointer,min,max,current,ttf)
}