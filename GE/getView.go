package GE

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

/**
This file contains getter for the most important UI objects
**/
//BUTTONS -------------------------------------------------------------------------------------------------------------------------------
func LoadButton(pathU, pathD string, X, Y, W, H float64, ChangeDrawDarkOnLeft bool) (*Button, error) {
	up, err := LoadEbitenImg(pathU)
	if err != nil {
		return nil, err
	}
	down, err := LoadEbitenImg(pathD)
	if err != nil {
		return nil, err
	}
	btn := GetUpDownImageButton(up, down, X, Y, W, H)
	btn.ChangeDrawDarkOnLeft = ChangeDrawDarkOnLeft
	return btn, nil
}

//Returns a Button showing a ImageObj
func GetButton(img *ImageObj, dark *ebiten.Image) *Button {
	b := &Button{}
	b.Img = img
	b.dark = dark
	b.light = img.Img
	b.Active = true
	return b
}

//Returns a Button with Text of a specific color on it
func GetTextButton(str, downStr string, ttf *truetype.Font, X, Y, H float64, textCol, backCol color.Color) *Button {
	img := GetTextImage(str, X, Y, H, ttf, textCol, backCol)
	dark := GetTextImage(downStr, X, Y, H, ttf, textCol, ReduceColor(backCol, ReduceColOnButtonDown))
	return GetButton(img, dark.Img)
}
func GetUpDownImageButton(up, down *ebiten.Image, X, Y, W, H float64) *Button {
	img := &ImageObj{up, nil, W, H, X, Y, 0}
	return GetButton(img, down)
}

//Returns a standard Button of an Image
func GetImageButton(eimg *ebiten.Image, X, Y, W, H float64) *Button {
	img := &ImageObj{eimg, nil, W, H, X, Y, 0}
	dark := ReduceColorImage(img.Img, ReduceColOnButtonDown)
	return GetButton(img, dark)
}

func GetSizedTextButton(str string, ttf *truetype.Font, X, Y, W, H float64, textCol, backCol color.Color) *Button {
	light := ebiten.NewImage(int(W), int(H))
	light.Fill(backCol)
	tImg := GetTextImage(str, 0, 0, H, ttf, textCol, color.RGBA{0, 0, 0, 0})
	tImg.X = (W - tImg.W) / 2
	tImg.Draw(light)
	btn := GetImageButton(light, X, Y, W, H)
	return btn
}

//EDITTEXT ------------------------------------------------------------------------------------------------------------------------------
//Returns a EditText, with a placeHolderText and a maximum number of runes
func GetEditText(placeHolderText string, X, Y, H float64, maxRunes int, ttf *truetype.Font, cols ...color.Color) (et *EditText) {
	imgo := ImageObj{H: H, X: X, Y: Y}
	et = &EditText{imgo, "", placeHolderText, 0, maxRunes, ttf, cols, 0, false, true, true, nil}
	return
}

//TEXTVIEW ------------------------------------------------------------------------------------------------------------------------------

//Returns a TextView displaying a specifc number of lines of a specific height
func GetTextView(text string, X, Y, lineHeight float64, displayLines int, ttf *truetype.Font, txtCol, backCol color.Color) (v *TextView) {
	v = &TextView{X: X, Y: Y, text: text, lineHeight: lineHeight, displayLines: displayLines}
	v.lines = HasLines(text)
	v.H = float64(v.lines) * lineHeight
	v.lineImages, v.W = GetTextLinesImages(text, X, Y, lineHeight, ttf, txtCol, backCol)
	return
}

//TABVIEW -------------------------------------------------------------------------------------------------------------------------------

//TabViewParams are used when creating a TabView, storing all necassary information
type TabViewParams struct {
	Back, Text color.Color
	TTF        *truetype.Font

	Dis              float64
	X, Y, W, H, TabH float64
	Curr             int
	Scrs             []UpdateAble

	Nms  []string
	Imgs []*ebiten.Image
	Dark []*ebiten.Image
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
		p.TabH = p.H * TabsHeight
	}
	if p.TTF == nil {
		p.TTF = StandardFont
	}
}

//Returns a TabView Created using TabViewParams
func GetTabView(p *TabViewParams) *TabView {
	p.fillDefault()
	if len(p.Imgs) > 0 {
		if len(p.Dark) > 0 {
			return getTabViewWithTwoImages(p.Imgs, p.Dark, p.Scrs, p.X, p.Y, p.W, p.H, p.TabH, p.Dis, p.Curr)
		} else {
			return getTabViewWithImages(p.Imgs, p.Scrs, p.X, p.Y, p.W, p.H, p.TabH, p.Dis, p.Curr)
		}
	}
	return getTabView(p.Nms, p.Scrs, p.X, p.Y, p.W, p.H, p.TabH, p.TTF, p.Text, p.Back, p.Dis, p.Curr)
}

//SCROLLBAR ------------------------------------------------------------------------------------------------------------------------------
//Returns a horizontal ScrollBar
func GetImageScrollbarFromPath(X, Y, W, H float64, bar, pointer string, min, max, current int, ttf *truetype.Font) (*ScrollBar, error) {
	b := &ScrollBar{min: min, max: max, current: current, ttf: ttf}
	bar_i, err := LoadEbitenImg(bar)
	if err != nil {
		return nil, err
	}
	pointer_i, err := LoadEbitenImg(pointer)
	if err != nil {
		return nil, err
	}
	b.Img = bar_i
	b.X = X
	b.Y = Y
	b.W = W
	b.H = H
	b.pointer = &ImageObj{Y: Y - H/2, W: H * 2, H: H * 2}
	b.value = &ImageObj{Y: Y - H/2, W: H * 2, H: H * 2}
	b.pointer.Img = pointer_i
	b.length = max - min
	b.stepsize = W / float64(b.length)
	b.UpdatePos()
	return b, nil
}

//Returns a horizontal ScrollBar
func GetImageScrollbar(X, Y, W, H float64, bar, pointer *ebiten.Image, min, max, current int, ttf *truetype.Font) (b *ScrollBar) {
	b = &ScrollBar{min: min, max: max, current: current, ttf: ttf}
	b.Img = bar
	b.X = X
	b.Y = Y
	b.W = W
	b.H = H
	b.pointer = &ImageObj{Y: Y - H/2, W: H * 2, H: H * 2}
	b.value = &ImageObj{Y: Y - H/2, W: H * 2, H: H * 2}
	b.pointer.Img = pointer
	b.length = max - min
	b.stepsize = W / float64(b.length)
	b.UpdatePos()
	return
}

//Returns a standard horizontal ScrollBar
func GetStandardScrollbar(X, Y, W, H float64, min, max, current int, ttf *truetype.Font) (b *ScrollBar) {

	bar := ebiten.NewImage(int(W), int(H))
	bar.Fill(&color.RGBA{0, 0, 0, 0})
	line := GetLineOfPoints(0, H/2, W, H/2, H/6)
	line.Fill(bar, &color.RGBA{0, 200, 50, 255})

	pointer := ebiten.NewImage(int(H*2), int(H*2))
	pointer.Fill(&color.RGBA{0, 0, 0, 0})
	pnts := genVertices(H, H, H, 100)
	pnts.Fill(pointer, &color.RGBA{200, 0, 50, 255})

	return GetImageScrollbar(X, Y, W, H, bar, pointer, min, max, current, ttf)
}

//ANIMATION ------------------------------------------------------------------------------------------------------------------------------

//Returns an Animation
func GetAnimation(X, Y, W, H float64, spriteWidth int, fps float64, sprites *ebiten.Image) (anim *Animation) {
	uP := 0
	if fps > 0 {
		uP = int(FPS / fps)
	}
	w, h := sprites.Size()
	anim = &Animation{ImageObj{X: X, Y: Y, W: W, H: H}, int(float64(w) / float64(spriteWidth)), 0, spriteWidth, h, uP, -1, sprites}
	anim.Update(0)
	return
}

//Returns an Animation
func GetAnimationFromParams(X, Y, W, H float64, p *Params, img *ebiten.Image) (anim *Animation) {
	w, h := img.Size()
	spriteWidth := p.Get("spriteWidth")
	fps := p.Get("fps")
	uP := 0
	if fps > 0 {
		uP = int(FPS / fps)
	}
	anim = &Animation{ImageObj{X: X, Y: Y, W: W, H: H}, int(float64(w) / float64(spriteWidth)), 0, int(spriteWidth), h, uP, -1, img}
	anim.Update(0)
	return
}

//Returns an DayNightAnimation
func GetDayNightAnim(X, Y, W, H float64, spriteWidth int, fps float64, sprites *ebiten.Image) (anim *DayNightAnim) {
	uP := 0
	if fps > 0 {
		uP = int(FPS / fps)
	}
	w, h := sprites.Size()
	dnimg := &DayNightImg{&ImageObj{}, &ImageObj{}}
	dnimg.SetParams(X, Y, W, H)
	anim = &DayNightAnim{dnimg, int(float64(w) / float64(spriteWidth)), 0, spriteWidth, h, uP, -1, 255, sprites, nil}
	anim.Update(0)
	return
}
func GetDayNightAnimFromParams(X, Y, W, H float64, pPath, imgPath string) (*DayNightAnim, error) {
	p := &Params{}
	err := p.LoadFromFile(pPath)
	if err != nil {
		return nil, err
	}
	img, err := LoadEbitenImg(imgPath)
	if err != nil {
		return nil, err
	}
	return GetDayNightAnim(X, Y, W, H, int(p.Get("spriteWidth")), p.Get("fps"), img), nil
}
