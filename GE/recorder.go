package GE

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/mjpeg"
)

/**
Recorder can log a specific number of *ebiten.Images and on demand save them as a ".avi" file or ".png"
the screen images should be added each frame and are rescaled automatically

be careful with high resolutions and long recording times, as the images will be stored in the RAM
when saving the CPU as well as the GPU might be heavily used for a long time
**/
func GetNewRecorder(frames, XRES, YRES, fps int) (r *Recorder) {
	r = &Recorder{frames: frames, current: 0, fps: fps, XRES: XRES, YRES: YRES}
	r.video = make([]*ebiten.Image, frames)
	r.saving = false
	return
}

type Recorder struct {
	frames, current, fps int
	XRES, YRES           int
	video                []*ebiten.Image
	saving               bool
}

func (r *Recorder) Delete() {
	if r == nil {
		return
	}
	r.saving = true
	for i := range r.video {
		r.video[i] = nil
	}
	r.video = nil
	r = nil
}
func (r *Recorder) NextFrame(img *ebiten.Image) {
	if !r.saving {
		idx := r.current
		r.current++
		if r.current >= r.frames {
			r.current = 0
			//runtime.GC()
		}
		if idx >= 0 && idx < len(r.video) {
			if r.video[idx] != nil {
				r.video[idx].Dispose()
			}
			r.video[idx] = DeepCopyScaleEbitenImage(img, r.XRES, r.YRES)
		}
	}
}

func (r *Recorder) IsSaving() bool {
	return r.saving
}

func (r *Recorder) SaveScreenShot(path string) error {
	idx := r.current - 1
	if idx < 0 {
		idx = r.frames - 1
	}
	if r.video[idx] != nil {
		return SaveEbitenImage(path+".png", r.video[idx])
	}
	return nil
}

//MAY take a long time
func (r *Recorder) Save(path string, done chan bool) {
	r.saving = true
	go func() {
		aw, err := mjpeg.New(path+".avi", int32(r.XRES), int32(r.YRES), int32(r.fps))
		ShitImDying(err)
		counter := 0
		for i := r.current; i < r.current+r.frames; i++ {
			idx := i
			if idx >= r.frames {
				idx -= r.frames
			}
			counter++

			if r.video[idx] != nil {
				buf := &bytes.Buffer{}
				ShitImDying(jpeg.Encode(buf, r.video[idx], nil))
				ShitImDying(aw.AddFrame(buf.Bytes()))
				buf.Reset()
				buf = nil
				r.video[idx].Dispose()
				r.video[idx] = nil
			}
		}
		ShitImDying(aw.Close())
		r.saving = false
		if done != nil {
			done <- true
		}
		aw = nil
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
