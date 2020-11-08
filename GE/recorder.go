package GE

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"os/exec"
	//"time"
	"fmt"
	"os"
	"io"
)
//cd ./tmp && ffmpeg -i %04d.png -vf fps=30 -y out.mp4

func GetNewRecorder(frames, XRES, YRES, fps int) (r *Recorder) {
	r = &Recorder{frames:frames, current:0, fps:fps}
	r.video = make([]*ebiten.Image, frames)
	r.drawer = &ImageObj{nil,nil,float64(XRES),float64(YRES),0,0,0}
	r.back = GetColoredImg(XRES,YRES, color.RGBA{0,0,0,255})
	r.saving = false
	return
}
type Recorder struct {
	frames, current,fps int
	back *ebiten.Image
	video []*ebiten.Image
	drawer *ImageObj
	saving bool
}
func (r *Recorder) NextFrame(img *ebiten.Image) {
	if !r.saving {
		idx := r.current
		r.current ++
		if r.current >= r.frames {
			r.current = 0
		}
		//copys the background
		newImg := DeepCopyEbitenImage(r.back)
		//draws the screen on the background
		r.drawer.Img = img
		r.drawer.DrawImageObj(newImg)
		r.video[idx] = newImg
	}
}

//MAY take a long time
func (r *Recorder) Save(copyPath string) {
	r.saving = true
	go func(){
		err := os.Mkdir("./tmp", 0777)
	    if err != nil {
	    	fmt.Println("Error creating dir: ", err)
	    }
	    
	    counter := 0
	    for i := r.current; i < r.current+r.frames; i++ {
			idx := i
			if idx >= r.frames {
				idx -= r.frames
			}
			counter ++
			idxStr := getFFMPEGstring(counter)
			img := r.video[idx]
			SaveEbitenImage(fmt.Sprintf("./tmp/%s.png",idxStr), img)
			r.video[idx] = nil
	    }
	    
		wdir,err := os.Getwd()
	    cmd := exec.Command("ffmpeg", "-i", "%04d.png", "-vf", fmt.Sprintf("fps=%v",r.fps), "-y", "out.mp4")
	    cmd.Dir = fmt.Sprintf("%s/tmp/", wdir)
	    err = cmd.Run()
	    if err != nil {
		    fmt.Println("Error running ffmpeg: ", err)
	    }
	    
	    CopyFile(fmt.Sprintf("%s/tmp/out.mp4", wdir), copyPath)
		
		err = os.RemoveAll(fmt.Sprintf("%s/tmp", wdir))
		if err != nil {
		    fmt.Println(err)
		}
		//fmt.Println("Removing tmp")
		r.saving = false
	}()
}

func getFFMPEGstring(i int) (out string) {
	out = ""
	if i < 10 {
		out = "000"
	}else if i < 100 {
		out = "00"
	}else if i < 1000 {
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