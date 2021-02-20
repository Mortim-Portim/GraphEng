package GE

import (
	"bytes"
	"fmt"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten"
	"github.com/icza/mjpeg"
)

//if ebiten.IsKeyPressed(ebiten.KeyC) && !g.rec.IsSaving() {
//		g.rec.Save("./res/out")
//	}
//	g.rec.NextFrame(screen)
func GetNewRecorder(frames, XRES, YRES, fps int) (r *Recorder) {
	r = &Recorder{frames: frames, current: 0, fps: fps}
	r.video = make([]*ebiten.Image, frames)
	r.drawer = &ImageObj{nil, nil, float64(XRES), float64(YRES), 0, 0, 0}
	r.back = GetColoredImg(XRES, YRES, color.RGBA{0, 0, 0, 255})
	r.saving = false
	return
}

type Recorder struct {
	frames, current, fps int
	back                 *ebiten.Image
	video                []*ebiten.Image
	drawer               *ImageObj
	saving               bool
}

func (r *Recorder) Delete() {
	if r == nil {
		return
	}
	r.saving = true
	r.back = nil
	for i := range r.video {
		r.video[i] = nil
	}
	r.video = nil
	r.drawer.Img = nil
	r = nil
}
func (r *Recorder) NextFrame(img *ebiten.Image) {
	if !r.saving {
		idx := r.current
		r.current++
		if r.current >= r.frames {
			r.current = 0
			debug.FreeOSMemory()
		}
		//copys the background
		newImg := DeepCopyEbitenImage(r.back)
		//draws the screen on the background
		r.drawer.Img = img
		r.drawer.DrawImageObj(newImg)
		r.video[idx] = newImg
	}
}

func (r *Recorder) IsSaving() bool {
	return r.saving
}

//MAY take a long time
func (r *Recorder) Save(path string, done chan bool) {
	r.saving = true
	go func() {
		aw, err := mjpeg.New(path+".avi", 200, 100, int32(r.fps))
		ShitImDying(err)
		counter := 0
		for i := r.current; i < r.current+r.frames; i++ {
			idx := i
			if idx >= r.frames {
				idx -= r.frames
			}
			counter++
			img := r.video[idx]
			r.video[idx] = nil
			if img != nil {
				buf := &bytes.Buffer{}
				err := jpeg.Encode(buf, img, nil)
				ShitImDying(err)
				err = aw.AddFrame(buf.Bytes())
				ShitImDying(err)
			}
		}
		ShitImDying(aw.Close())
		r.saving = false
		if done != nil {
			done <- true
		}
	}()
}

func getFFMPEGstring(i int) (out string) {
	out = ""
	if i < 10 {
		out = "000"
	} else if i < 100 {
		out = "00"
	} else if i < 1000 {
		out = "0"
	}
	out += fmt.Sprintf("%v", i)
	return
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
