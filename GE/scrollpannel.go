package GE

import (
    "github.com/hajimehoshi/ebiten"
)

func GetScrollPanel(x, y, w, h float64, buttons ...*Button) (sp *ScrollPanel) {
    ebitimg := GetEmptyImage(int(w), int(h))
    panel := &ImageObj{ebitimg, nil, w, h, x, y, 0}

    sp = &ScrollPanel{buttons, panel}
    sp.Redraw()
    return
}

type ScrollPanel struct {
    content []*Button
    panel   *ImageObj
}

func (sp *ScrollPanel) Add(button *Button) {
    sp.content = append(sp.content, button)
    sp.Redraw()
}

func (sp *ScrollPanel) Scroll(dy float64) {
    for _, cont := range sp.content {
        cont.Img.Y += dy
    }
    sp.Redraw()
}

func (sp *ScrollPanel) Redraw() {
    screen := sp.panel.Img
    screen.Clear()

    for  _, cont := range sp.content {
        if cont.Img.Y < sp.panel.Y+sp.panel.W && cont.Img.Y+cont.Img.H > sp.panel.Y {
            cont.Active = true
            copy := cont.Img.Copy()
            copy.X -= sp.panel.X
            copy.Y -= sp.panel.Y
            copy.Draw(screen)
        } else {
            cont.Active = false
        }
    }
}

func (sp *ScrollPanel) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
    return sp.Update, sp.Draw
}

func (sp *ScrollPanel) Start(screen *ebiten.Image, data interface{}) {}

func (sp *ScrollPanel) Stop(screen *ebiten.Image, data interface{}) {}

func (sp *ScrollPanel) Update(frame int) {
    for _, cont := range sp.content {
        cont.Update(frame)
    }

    x, y := ebiten.CursorPosition()
    _, scroll := ebiten.Wheel()
    hasFocus := int(sp.panel.X) <= x && x < int(sp.panel.X+sp.panel.W) && int(sp.panel.Y) <= y && y < int(sp.panel.Y+sp.panel.H)

    if hasFocus && scroll > 0 {
        sp.Scroll(10)
    }

    if hasFocus && scroll < 0 {
        sp.Scroll(-10)
    }
}

func (sp *ScrollPanel) Draw(screen *ebiten.Image) {
    sp.panel.Draw(screen)
}
