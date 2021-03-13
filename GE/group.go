package GE

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/**
Group represents a group of structs implementing UpdateAble
all members of the group can be Started, Stopped, Updated and Drawn using one function call

Example:
A EditText and Animation form a group, that is always updated and drawn at the same time

Group implements UpdateAble
**/

func GetGroup(members ...UpdateAble) (g *Group) {
	g = &Group{}
	g.UpdateFuncs = make([]UpdateFunc, 0)
	g.DrawFuncs = make([]DrawFunc, 0)
	g.Members = make([]UpdateAble, 0)
	g.Add(members...)
	return
}
func (g *Group) Add(members ...UpdateAble) {
	for _, mmb := range members {
		f_u, f_d := mmb.Init(g.InitScreen, g.InitData)
		if f_u != nil {
			g.UpdateFuncs = append(g.UpdateFuncs, f_u)
		}
		if f_d != nil {
			g.DrawFuncs = append(g.DrawFuncs, f_d)
		}
	}
	g.Members = append(g.Members, members...)
}
func (g *Group) Set(member UpdateAble, idx int) {
	if idx < 0 || idx >= len(g.Members) {
		return
	}
	f_u, f_d := member.Init(g.InitScreen, g.InitData)
	if f_u != nil {
		g.UpdateFuncs[idx] = f_u
	}
	if f_d != nil {
		g.DrawFuncs[idx] = f_d
	}
	g.Members[idx] = member
}

type Group struct {
	Members     []UpdateAble
	UpdateFuncs []UpdateFunc
	DrawFuncs   []DrawFunc

	InitScreen *ebiten.Image
	InitData   interface{}
}

func (g *Group) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	g.InitScreen = screen
	g.InitData = data
	return g.Update, g.Draw
}
func (g *Group) Start(screen *ebiten.Image, data interface{}) {
	for _, mmb := range g.Members {
		mmb.Start(screen, data)
	}
}
func (g *Group) Stop(screen *ebiten.Image, data interface{}) {
	for _, mmb := range g.Members {
		mmb.Stop(screen, data)
	}
}
func (g *Group) Update(frame int) {
	for _, fnc := range g.UpdateFuncs {
		fnc(frame)
	}
}
func (g *Group) Draw(screen *ebiten.Image) {
	for _, fnc := range g.DrawFuncs {
		fnc(screen)
	}
}
