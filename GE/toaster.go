package GE

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
)

/**
Toaster represents a Factory for Toasts(small messages, that are displayed only for a limited time)
A Toast can be timed using a specific number of frames or a channel that signals when the toast is supposed to stop displaying
**/

type Toast interface {
	ImageObj() *ImageObj
	Update()
	TimeIsOver() bool
}

func MakeToastTimed(img *ImageObj, duration int) *toastTimed {
	return &toastTimed{img, duration, 0}
}

type toastTimed struct {
	imgO              *ImageObj
	duration, counter int
}

func (t *toastTimed) ImageObj() *ImageObj { return t.imgO }
func (t *toastTimed) Update() {
	t.counter++
}
func (t *toastTimed) TimeIsOver() bool {
	return t.counter > t.duration
}

func MakeToastChan(img *ImageObj, c chan bool) *toastChan {
	t := &toastChan{img, false}
	go func() {
		<-c
		t.isOver = true
	}()
	return t
}

type toastChan struct {
	imgO   *ImageObj
	isOver bool
}

func (t *toastChan) ImageObj() *ImageObj { return t.imgO }
func (t *toastChan) Update()             {}
func (t *toastChan) TimeIsOver() bool {
	return t.isOver
}

func GetNewToaster(XRES, YRES, RelScreenPos, RelToastH float64, TTF *truetype.Font, BackCol, TextCol color.Color) (t *Toaster) {
	t = &Toaster{XRES: XRES, YRES: YRES, RelScreenPos: RelScreenPos, RelToastH: RelToastH, TTF: TTF, BackCol: BackCol, TextCol: TextCol}
	t.Start(nil, nil)
	return
}

type Toaster struct {
	XRES, YRES              float64
	RelScreenPos, RelToastH float64
	TTF                     *truetype.Font
	BackCol, TextCol        color.Color

	xposm, tH float64
	Toasts    []Toast
}

func (t *Toaster) New(msg string, frames int, c chan bool) {
	img := t.new(msg)
	if c == nil {
		t.Toasts = append(t.Toasts, MakeToastTimed(img, frames))
	} else {
		t.Toasts = append(t.Toasts, MakeToastChan(img, c))
	}
}

func (t *Toaster) new(msg string) *ImageObj {
	imgO := GetTextImage(msg, t.xposm, 0, t.tH, t.TTF, t.TextCol, t.BackCol)
	imgO.SetMiddleX(t.xposm)
	if imgO.X < 0 {
		imgO.X = 0
	} else if imgO.X+imgO.W > t.XRES {
		imgO.X = t.XRES - imgO.W
	}
	return imgO
}

func (t *Toaster) RemoveToast(i int) {
	t.Toasts[i] = t.Toasts[len(t.Toasts)-1]
	t.Toasts = t.Toasts[:len(t.Toasts)-1]
}
func (t *Toaster) Update(frame int) {
	rems := 0
	for idx := range t.Toasts {
		t.Toasts[idx-rems].Update()
		if t.Toasts[idx-rems].TimeIsOver() {
			t.RemoveToast(idx - rems)
			rems++
		}
	}
}
func (t *Toaster) Draw(screen *ebiten.Image) {
	for i, tst := range t.Toasts {
		tst.ImageObj().Y = float64(i) * t.tH
		tst.ImageObj().Draw(screen)
	}
}

func (t *Toaster) Start(screen *ebiten.Image, data interface{}) {
	t.xposm = t.XRES * t.RelScreenPos
	t.tH = t.YRES * t.RelToastH
}
func (t *Toaster) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return t.Update, t.Draw
}
func (t *Toaster) Stop(screen *ebiten.Image, data interface{}) {}
