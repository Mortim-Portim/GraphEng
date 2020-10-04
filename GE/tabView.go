package GE

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type TabView struct {
	X, Y, W, H, TabH float64
	
	TabBtns *Group
	
	Screens *Group
	CurrentTab int
}
func (t *TabView) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	t.TabBtns.Init(screen, data)
	t.Screens.Init(screen, data)
	return t.Update, t.Draw
}
func (t *TabView) Start(screen *ebiten.Image, data interface{}) {
	t.TabBtns.Init(screen, data)
	t.Screens.Init(screen, data)
}
func (t *TabView) Stop(screen *ebiten.Image, data interface{}) {
	t.TabBtns.Init(screen, data)
	t.Screens.Init(screen, data)
}
func (t *TabView) Update(frame int) {
	t.Screens.UpdateFuncs[t.CurrentTab](frame)
	for i,mmb := range(t.TabBtns.Member) {
		btn := mmb.(*Button)
		btn.DrawDark = false
		if i == t.CurrentTab {
			btn.DrawDark = true
		}
		t.TabBtns.UpdateFuncs[i](frame)
	}
}
func (t *TabView) Draw(screen *ebiten.Image) {
	t.Screens.DrawFuncs[t.CurrentTab](screen)
	t.TabBtns.Draw(screen)
}

func (t *TabView) OnClick(b *Button) {
	if b.LPressed || b.RPressed {
		t.CurrentTab = b.Data.(int)
	}
}
func getTabView(Names []string, screens []UpdateAble, X, Y, W, H, TabH float64, ttf *truetype.Font, txtCol, backCol color.Color, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH, nil, GetGroup(screens...), curr}
	TabBtns := make([]UpdateAble, len(Names))
	for i,name := range(Names) {
		TabBtns[i] = GetTextButton(name, name, ttf, X, Y, TabH, txtCol, backCol)
		TabBtns[i].(*Button).RegisterOnEvent(v.OnClick)
		TabBtns[i].(*Button).Data = i
	}
	v.TabBtns = GetGroup(TabBtns...)
	for i,mmb := range(v.TabBtns.Member[1:]) {
		tab := mmb.(*Button)
		tabm1 := v.TabBtns.Member[i].(*Button)
		tab.Img.X = (tabm1.Img.X+tabm1.Img.W+W*dis)
	}
	return v
}
func getTabViewWithImages(imgs []*ebiten.Image, screens []UpdateAble, X, Y, W, H, TabH float64, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH,  nil, GetGroup(screens...), curr}
	TabBtns := make([]UpdateAble, len(imgs))
	for i,img := range(imgs) {
		TabBtns[i] = GetImageButton(img, X, Y, 0, 0)
		TabBtns[i].(*Button).RegisterOnEvent(v.OnClick)
		TabBtns[i].(*Button).Img.ScaleToOriginalSize()
		TabBtns[i].(*Button).Img.ScaleToY(TabH)
		TabBtns[i].(*Button).Data = i
	}
	v.TabBtns = GetGroup(TabBtns...)
	for i,mmb := range(v.TabBtns.Member[1:]) {
		tab := mmb.(*Button)
		tabm1 := v.TabBtns.Member[i].(*Button)
		tab.Img.X = (tabm1.Img.X+tabm1.Img.W+W*dis)
	}
	return v
}