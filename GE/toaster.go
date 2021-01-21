package GE

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type toast struct {
	imgO *ImageObj
	duration, counter int
}
func (t *toast) Update() {
	t.counter ++
}
func (t *toast) TimeIsOver() bool {
	return t.counter > t.duration
}

func GetNewToaster(XRES, YRES, RelScreenPos, RelToastH float64, TTF *truetype.Font, BackCol, TextCol color.Color) (t *Toaster) {
	t = &Toaster{XRES:XRES, YRES:YRES, RelScreenPos:RelScreenPos, RelToastH:RelToastH, TTF:TTF, BackCol:BackCol, TextCol:TextCol}
	t.Start(nil, nil)
	return 
}
type Toaster struct {
	XRES, YRES float64
	RelScreenPos, RelToastH float64
	TTF *truetype.Font
	BackCol, TextCol color.Color
	
	xposm, tH float64
	Toasts []*toast
}

func (t *Toaster) New(msg string, frames int) {
	imgO := GetTextImage(msg, t.xposm, 0, t.tH, t.TTF, t.TextCol, t.BackCol)
	imgO.SetMiddleX(t.xposm)
	if imgO.X < 0 {
		imgO.X = 0
	}else if imgO.X+imgO.W > t.XRES {
		imgO.X = t.XRES-imgO.W
	}
	t.Toasts = append(t.Toasts, &toast{imgO, frames, 0})
}

func (t *Toaster) RemoveToast(i int) {
	t.Toasts[i] = t.Toasts[len(t.Toasts)-1]
	t.Toasts = t.Toasts[:len(t.Toasts)-1]
}
func (t *Toaster) Update(frame int) {
	rems := 0
	for idx, _ := range t.Toasts {
		t.Toasts[idx-rems].Update()
		if t.Toasts[idx-rems].TimeIsOver() {
			t.RemoveToast(idx-rems)
			rems ++
		}
	}
}
func (t *Toaster) Draw(screen *ebiten.Image) {
	for i,tst := range(t.Toasts) {
		tst.imgO.Y = float64(i)*t.tH
		tst.imgO.Draw(screen)
	}
}

func (t *Toaster) Start(screen *ebiten.Image, data interface{}) {
	t.xposm = t.XRES*t.RelScreenPos
	t.tH = t.YRES*t.RelToastH
}
func (t *Toaster) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {return t.Update, t.Draw}
func (t *Toaster) Stop(screen *ebiten.Image, data interface{}) {}