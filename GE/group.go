package GE

import (
	"github.com/hajimehoshi/ebiten"
)
func GetGroup(members ...UpdateAble) (g *Group) {
	g = &Group{}
	g.Member = members
	return
}
func (g *Group) Add(members ...UpdateAble) {
	for _,mmb := range(members) {
		f_u, f_d := mmb.Init(g.InitScreen, g.InitData)
		if f_u != nil {
			g.UpdateFuncs = append(g.UpdateFuncs, f_u)
		}
		if f_d != nil {
			g.DrawFuncs = append(g.DrawFuncs, f_d)
		}
	}
}
type Group struct {
	Member []UpdateAble
	UpdateFuncs []UpdateFunc
	DrawFuncs []DrawFunc
	
	InitScreen *ebiten.Image
	InitData interface{}
}

func (g *Group) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	g.InitScreen = screen; g.InitData = data
	g.UpdateFuncs = make([]UpdateFunc, 0)
	g.DrawFuncs = make([]DrawFunc, 0)
	for _,mmb := range(g.Member) {
		f_u, f_d := mmb.Init(screen, data)
		if f_u != nil {
			g.UpdateFuncs = append(g.UpdateFuncs, f_u)
		}
		if f_d != nil {
			g.DrawFuncs = append(g.DrawFuncs, f_d)
		}
	}
	return g.Update, g.Draw
}
func (g *Group) Start(screen *ebiten.Image, data interface{}) {
	for _,mmb := range(g.Member) {
		mmb.Start(screen, data)
	}
}
func (g *Group) Stop(screen *ebiten.Image, data interface{}) {
	for _,mmb := range(g.Member) {
		mmb.Stop(screen, data)
	}
}
func (g *Group) Update(frame int) {
	for _,fnc := range(g.UpdateFuncs) {
		fnc(frame)
	}
}
func (g *Group) Draw(screen *ebiten.Image) {
	for _,fnc := range(g.DrawFuncs) {
		fnc(screen)
	}
}