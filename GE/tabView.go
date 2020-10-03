package GE

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type TabView struct {
	X, Y, W, H, TabH float64
	
	TabBtns []*Button
	Screens []UpdateAble
	CurrentTab int
}
func (t *TabView) Init(screen *ebiten.Image, data interface{}) {}
func (t *TabView) Start(screen *ebiten.Image, data interface{}) {}
func (t *TabView) Stop(screen *ebiten.Image, data interface{}) {}
func (t *TabView) Update() {
	if t.CurrentTab < len(t.Screens) {
		t.Screens[t.CurrentTab].Update()
	}
	for i,btn := range(t.TabBtns) {
		btn.DrawDark = false
		if i == t.CurrentTab {
			btn.DrawDark = true
		}
		btn.Update()
	}
}
func (t *TabView) Draw(screen *ebiten.Image) {
	if t.CurrentTab < len(t.Screens) {
		t.Screens[t.CurrentTab].Draw(screen)
	}
	for _,btn := range(t.TabBtns) {
		btn.Draw(screen)
	}
}

func (t *TabView) OnClick(b *Button) {
	if b.LPressed || b.RPressed {
		t.CurrentTab = b.Data.(int)
	}
}
func getTabView(Names []string, screens []UpdateAble, X, Y, W, H, TabH float64, ttf *truetype.Font, txtCol, backCol color.Color, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH, nil, screens, curr}
	v.TabBtns = make([]*Button, len(Names))
	for i,name := range(Names) {
		v.TabBtns[i] = GetTextButton(name, name, ttf, X, Y, TabH, txtCol, backCol, v.OnClick, v.OnClick)
		v.TabBtns[i].Data = i
	}
	for i,tab := range(v.TabBtns[1:]) {
		tab.Img.X = (v.TabBtns[i].Img.X+v.TabBtns[i].Img.W+W*dis)
	}
	return v
}
func getTabViewWithImages(imgs []*ebiten.Image, screens []UpdateAble, X, Y, W, H, TabH float64, dis float64, curr int) (v *TabView) {
	v = &TabView{X, Y, W, H, TabH, nil, screens, curr}
	v.TabBtns = make([]*Button, len(imgs))
	for i,img := range(imgs) {
		v.TabBtns[i] = GetImageButton(img, X, Y, 0, 0, v.OnClick, v.OnClick)
		v.TabBtns[i].Img.ScaleToOriginalSize()
		v.TabBtns[i].Img.ScaleToY(TabH)
		v.TabBtns[i].Data = i
	}
	for i,tab := range(v.TabBtns[1:]) {
		tab.Img.X = (v.TabBtns[i].Img.X+v.TabBtns[i].Img.W+W*dis)
	}
	return v
}