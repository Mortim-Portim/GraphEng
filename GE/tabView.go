package GE

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type TabView struct {
	X, Y, W, H, TabH float64
	
	TabBtns []*Button
	TabUpdateFuncs []UpdateFunc
	TabDrawFuncs []DrawFunc
	
	Screens []UpdateAble
	ScreenUpdateFuncs []UpdateFunc
	ScreenDrawFuncs []DrawFunc
	CurrentTab int
}
func (t *TabView) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	t.TabUpdateFuncs = make([]UpdateFunc, 0)
	t.ScreenUpdateFuncs = make([]UpdateFunc, 0)
	t.TabDrawFuncs = make([]DrawFunc, 0)
	t.ScreenDrawFuncs = make([]DrawFunc, 0)
	for _,btn := range(t.TabBtns) {
		f_u, f_d := btn.Init(screen, data)
		if f_u != nil {
			t.TabUpdateFuncs = append(t.TabUpdateFuncs, f_u)
		}
		if f_d != nil {
			t.TabDrawFuncs = append(t.TabDrawFuncs, f_d)
		}
	}
	for _,scr := range(t.Screens) {
		f_u, f_d := scr.Init(screen, data)
		if f_u != nil {
			t.ScreenUpdateFuncs = append(t.ScreenUpdateFuncs, f_u)
		}
		if f_d != nil {
			t.ScreenDrawFuncs = append(t.ScreenDrawFuncs, f_d)
		}
	}
	
	return t.Update, t.Draw
}
func (t *TabView) Start(screen *ebiten.Image, data interface{}) {}
func (t *TabView) Stop(screen *ebiten.Image, data interface{}) {}
func (t *TabView) Update(frame int) {
	t.ScreenUpdateFuncs[t.CurrentTab](frame)
	for i,btn := range(t.TabBtns) {
		btn.DrawDark = false
		if i == t.CurrentTab {
			btn.DrawDark = true
		}
		t.TabUpdateFuncs[i](frame)
	}
}
func (t *TabView) Draw(screen *ebiten.Image) {
	t.ScreenDrawFuncs[t.CurrentTab](screen)
	for _,fnc := range(t.TabDrawFuncs) {
		fnc(screen)
	}
}

func (t *TabView) OnClick(b *Button) {
	if b.LPressed || b.RPressed {
		t.CurrentTab = b.Data.(int)
	}
}
func getTabView(Names []string, screens []UpdateAble, X, Y, W, H, TabH float64, ttf *truetype.Font, txtCol, backCol color.Color, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH, nil, nil, nil, screens, nil, nil, curr}
	v.TabBtns = make([]*Button, len(Names))
	for i,name := range(Names) {
		v.TabBtns[i] = GetTextButton(name, name, ttf, X, Y, TabH, txtCol, backCol)
		v.TabBtns[i].RegisterOnEvent(v.OnClick)
		v.TabBtns[i].Data = i
	}
	for i,tab := range(v.TabBtns[1:]) {
		tab.Img.X = (v.TabBtns[i].Img.X+v.TabBtns[i].Img.W+W*dis)
	}
	return v
}
func getTabViewWithImages(imgs []*ebiten.Image, screens []UpdateAble, X, Y, W, H, TabH float64, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH, nil, nil, nil, screens, nil, nil, curr}
	v.TabBtns = make([]*Button, len(imgs))
	for i,img := range(imgs) {
		v.TabBtns[i] = GetImageButton(img, X, Y, 0, 0)
		v.TabBtns[i].RegisterOnEvent(v.OnClick)
		v.TabBtns[i].Img.ScaleToOriginalSize()
		v.TabBtns[i].Img.ScaleToY(TabH)
		v.TabBtns[i].Data = i
	}
	for i,tab := range(v.TabBtns[1:]) {
		tab.Img.X = (v.TabBtns[i].Img.X+v.TabBtns[i].Img.W+W*dis)
	}
	return v
}