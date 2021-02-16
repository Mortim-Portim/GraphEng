package GE

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"log"
	"runtime/pprof"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/mortim-portim/GraphEng/res"
	"golang.org/x/image/font"
)

const allLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz/.:"
const StandardFontSize = 64

var StandardFont *truetype.Font

/**
UpdateAble is an interface that can be initialized, started and stoped
An UpdateAble may register an Update and Draw function when initialized

Init should be called only once
Start should be called when UpdateAble becomes visible
Start should be called when UpdateAble becomes invisible
**/
type UpdateFunc func(frame int)
type DrawFunc func(screen *ebiten.Image)
type UpdateAble interface {
	Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc)
	Start(screen *ebiten.Image, data interface{})
	Stop(screen *ebiten.Image, data interface{})
}

//Initializes the Standard Font (use "" to load standard font), the Random generator, the audio context and the parameter
func Init(FontPath string) {
	if len(FontPath) > 0 {
		StandardFont = ParseFont(FontPath)
	} else {
		StandardFont = ParseFontFromBytes(res.VT323)
	}
	InitAudioContext()
	rand.Seed(time.Now().UnixNano())
	InitParams(nil)
}

//Parses a font from bytes
func ParseFontFromBytes(fnt []byte) *truetype.Font {
	tt, err := truetype.Parse(fnt)
	CheckErr(err)
	return tt
}

//Parses a font from the filesystem
func ParseFont(path string) *truetype.Font {
	font, err1 := ioutil.ReadFile(path)
	CheckErr(err1)
	tt, err2 := truetype.Parse(font)
	CheckErr(err2)
	return tt
}

var TTF_FACES = make(map[*truetype.Font]font.Face)
func GetFace(ttf *truetype.Font) font.Face {
	fc, ok := TTF_FACES[ttf]
	if ok {return fc}
	mplusNormalFont := truetype.NewFace(ttf, &truetype.Options{
		Size:    StandardFontSize,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	TTF_FACES[ttf] = mplusNormalFont
	return mplusNormalFont
}
//Draws text of the given font on an Image
func MakePopUp(str string, ttf *truetype.Font, textCol, backCol color.Color) *ebiten.Image {
	fnt := GetFace(ttf)
	w, h := MeasureString(str, fnt)

	popUpBack, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	popUpBack.Fill(backCol)
	xP, yP := h/6, h/4*3
	text.Draw(popUpBack, str, fnt, int(xP), int(yP), textCol)
	return popUpBack
}

func MakePopUpOld(str string, size float64, ttf *truetype.Font, textCol, backCol color.Color) *ebiten.Image {
	fnt := truetype.NewFace(ttf, &truetype.Options{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	w, h := MeasureString(str, fnt)

	popUpBack, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	popUpBack.Fill(backCol)
	xP, yP := h/6, h/4*3
	text.Draw(popUpBack, str, fnt, int(xP), int(yP), textCol)
	return popUpBack
}

//Returns an ImageObj with a single line text on it
func GetTextImage(textStr string, X, Y, H float64, ttf *truetype.Font, txtCol, backCol color.Color) *ImageObj {
	imgo := &ImageObj{H: H, X: X, Y: Y}
	if len(textStr) > 0 {
		textImg := MakePopUp(textStr, ttf, txtCol, &color.RGBA{0, 0, 0, 0})
		w, h := textImg.Size()
		W := float64(w) * H / float64(h)
		imgo.W = W

		Back, _ := ebiten.NewImage(int(W), int(H), ebiten.FilterDefault)
		Back.Fill(backCol)
		imgo.Img = Back

		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(H/float64(h), H/float64(h))
		//op.GeoM.Translate(H*0.25,0)
		imgo.Img.DrawImage(textImg, op)
	} else {
		Back, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
		Back.Fill(backCol)
		imgo.Img = Back
	}

	return imgo
}

//Returns slice of ImageObjs that all represent a line of textStr
func GetTextLinesImages(textStr string, X, Y, lineHeight float64, ttf *truetype.Font, txtCol, backCol color.Color) (lineImgs []*ImageObj, maxWidth float64) {
	lines := strings.Split(textStr, "\n")
	lineImgs = make([]*ImageObj, len(lines))
	maxWidth = 0
	for i, str := range lines {
		for str[0] == " "[0] {
			str = str[1:]
		}
		lineImgs[i] = GetTextImage(str, X, Y+float64(i)*lineHeight, lineHeight, ttf, txtCol, backCol)
		if lineImgs[i].W > maxWidth {
			maxWidth = lineImgs[i].W
		}
	}
	return
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 20
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func areKeysPressed(keys ...ebiten.Key) bool {
	keysPressed := 0
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			keysPressed++
		}
	}
	if keysPressed < len(keys) {
		return false
	}
	for _, key := range keys {
		if inpututil.KeyPressDuration(key) == 1 {
			return true
		}
	}
	return false
}
func IsKeyJustDown(key ebiten.Key) bool {
	if inpututil.KeyPressDuration(key) == 1 {
		return true
	}
	return false
}

//Panics if an error occured
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Checks if two colors have the same red, green, blue and alpha value
func SameCols(col, col2 color.Color) bool {
	r, g, b, a := col.RGBA()
	r2, g2, b2, a2 := col2.RGBA()
	if r == r2 && g == g2 && b == b2 && a == a2 {
		return true
	}
	return false
}

//Reduces the color values of an image to a minimum of 0
func ReduceColor(col color.Color, delta int) color.Color {
	r, g, b, a := col.RGBA()
	newR := int(r) - delta
	if newR < 0 {
		newR = 0
	}
	newG := int(g) - delta
	if newG < 0 {
		newG = 0
	}
	newB := int(b) - delta
	if newB < 0 {
		newB = 0
	}
	return &color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)}
}

//Reduces the alpha value of an Image, making it more transparent
func ReduceColorImage(img *ebiten.Image, val int) (reduced *ebiten.Image) {
	W, H := img.Size()
	Back, _ := ebiten.NewImage(W, H, ebiten.FilterDefault)
	reduced = Back
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, float64(255-val)/255.0)
	reduced.DrawImage(img, op)
	return
}

//Returns true if e is in s
func containsI(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Returns true if e is in s
func containsS(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//measures a string, storing the maximum height of that Face in a map to be reused
var faceHeight = make(map[font.Face]int)

func MeasureString(str string, faceTTF font.Face) (x, y int) {
	h, ok := faceHeight[faceTTF]
	if !ok {
		rectAll := text.BoundString(faceTTF, allLetters)
		h = rectAll.Max.Y - rectAll.Min.Y
		faceHeight[faceTTF] = h
	}
	rect := text.BoundString(faceTTF, str+"#")
	x, y = rect.Max.X-rect.Min.X, h*(strings.Count(str, "\n")+1)+h/6
	//fmt.Println(rect.String(), ":     ", x, ":", y)
	return
}

//Generates a slice of Points forming a circle
func genVertices(X, Y, R float64, num int) *Points {
	centerX := X
	centerY := Y
	r := R

	vs := make([]*Vector, 0)
	for i := 0; i <= num; i++ {
		rate := float64(i) / float64(num)
		vs = append(vs, &Vector{
			X: float64(r*math.Cos(2*math.Pi*rate)) + centerX,
			Y: float64(r*math.Sin(2*math.Pi*rate)) + centerY,
			Z: 0})
	}

	vs = append(vs, &Vector{
		X: centerX,
		Y: centerY,
		Z: 0})
	ps := Points(vs)
	return &ps
}

func GetAllFiles(root string) ([]string, error) {
	if root[len(root)-1:] != "/" {root += "/"}
	var files []string
	f, err := os.Open(root)
	if err != nil {return files, err}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {return files, err}
	err = AppendAllFilesInAllFolders(&files, root, fileInfo)
	return files, nil
}
func AppendAllFilesInAllFolders(files *[]string, root string, folderFs []os.FileInfo) error {
	if root[len(root)-1:] != "/" {root += "/"}
	for _, file := range folderFs {
		if !file.IsDir() {
			*files = append(*files, root+file.Name())
		}else{
			path := root+file.Name()+"/"
			f, err := os.Open(path)
			if err != nil {return err}
			fileInfo, err := f.Readdir(-1)
			f.Close()
			if err != nil {return err}
			err = AppendAllFilesInAllFolders(files, path, fileInfo)
			if err != nil {return err}
		}
	}
	return nil
}
func OSReadDir(root string) ([]string, error) {
	if root[len(root)-1:] != "/" {
		root += "/"
	}
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

//Calls handleFile for each different filename
func ReadAllFiles(dir string, handleFile func(name string)) error {
	files, err := OSReadDir(dir)
	if err != nil {
		return err
	}
	var currentError error
	names := make([]string, 0)
	for _, f := range files {
		n := strings.Split(f, ".")[0]
		if !IsStringInList(n, names) {
			currentError = err
			names = append(names, n)
			handleFile(n)
		}
	}
	return currentError
}

func ShutDown() {
	if runtime.GOOS == "linux" {
		if err := exec.Command("shutdown", "now").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
		}
	} else if runtime.GOOS == "windows" {
		if err := exec.Command("cmd", "/C", "shutdown", "/t", "0", "/s").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
		}
	} else {
		panic("Wat is dat denn fürn Betriebssystem????")
	}
}

var LOGFILE *os.File

func SetLogFile(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	LOGFILE = f
}
func LogToFile(text string) {
	if LOGFILE == nil {
		return
	}
	if _, err := LOGFILE.WriteString(text); err != nil {
		panic(err)
	}
}
func CloseLogFile() {
	if LOGFILE == nil {
		return
	}
	LOGFILE.Close()
}
func ShitImDying(err error) {
	if err != nil {
		defer func() {
			help := fmt.Sprintf("ShitImDying: %v\nStacktrace: %s", err, string(debug.Stack()))
			if LOGFILE != nil {
				LogToFile(help)
			}
		}()
		panic(err)
	}
}
/**
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
**/
var CPU_PROF_F *os.File
func StartProfiling(cpuprofile *string) {
	if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal("could not create CPU profile: ", err)
        }
        CPU_PROF_F = f
        if err := pprof.StartCPUProfile(CPU_PROF_F); err != nil {
            log.Fatal("could not start CPU profile: ", err)
        }
    }
}
func StopProfiling(cpuprofile, memprofile *string) {
	if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
    }
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
		CPU_PROF_F.Close()
	}
}