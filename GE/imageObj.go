package GE

import (
	"os"
	"image"
	_ "image/jpeg"
	_ "image/png"
	//"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	//"github.com/golang/freetype/truetype"
	//"golang.org/x/image/font"
	"bytes"
	"github.com/nfnt/resize"
	"image/color"
	"fmt"
	"math"
	"errors"
)
//Loads the Image in Original Resolution
func LoadImgObj(path string, width, height, x, y, angle float64) *ImageObj {
	_, img := LoadImg(path)
	//scaledImg := resize.Resize(uint(width), uint(height), *img, resize.NearestNeighbor)
	return &ImageObj{ImgToEbitenImg(img), img, width, height, x, y, angle}
}
//Stores the original, ebiten image and dimensions
type ImageObj struct {
	Img *ebiten.Image; OriginalImg *image.Image
	W, H, X, Y, Angle float64
}
//Sets the middle of the Image to x,y
func (obj *ImageObj) SetMiddle(x,y float64) {
	obj.X = x-obj.W/2; obj.Y = y-obj.H/2
}
//Sets the middle of the Image to x
func (obj *ImageObj) SetMiddleX(x float64) {
	obj.X = x-obj.W/2
}
//Sets the middle of the Image to y
func (obj *ImageObj) SetMiddleY(y float64) {
	obj.Y = y-obj.H/2
}
//Puts the Image in between x1,y1 and x2,y2
func (obj *ImageObj) PutInBounds(x1,y1, x2,y2 float64) {
	obj.X = x1; obj.Y = y1
	obj.W = x2-x1; obj.H = y2-y1
}
//Returns Information
func (obj *ImageObj) Print() string {
	return fmt.Sprintf("OW: %v, OH: %v, W: %f, H: %f, X: %f, Y: %f, Angle: %f", (*obj.OriginalImg).Bounds().Max.X, (*obj.OriginalImg).Bounds().Max.Y, obj.W, obj.H, obj.X, obj.Y, obj.Angle)
}
//Copys the ImageObj
func (obj *ImageObj) Copy() *ImageObj {
	return &ImageObj{obj.Img, obj.OriginalImg, obj.W, obj.H, obj.X, obj.Y, obj.Angle}
}
//Scales the Original Image in a dimension keeping the aspect ratio
func (obj *ImageObj) ScaleOriginal(width, height float64) {
	scaledImg := resize.Resize(uint(width), uint(height), *obj.OriginalImg, resize.NearestNeighbor)
	obj.OriginalImg = &scaledImg
	obj.Img = ImgToEbitenImg(&scaledImg)
}
func (obj *ImageObj) ScaleToOriginalSize() {
	w,h := obj.Img.Size()
	obj.W = float64(w)
	obj.H = float64(h)
}
//Scales the Image in a dimension keeping the aspect ratio
func (obj *ImageObj) ScaleDim(newval float64, dim int) {
	if dim == 0 {
		obj.ScaleToX(newval)
	}else if dim == 1 {
		obj.ScaleToY(newval)
	}
}
//Scales the Image in the X dimension keeping the aspect ratio
func (obj *ImageObj) ScaleToX(newWidth float64) {
	obj.H *= newWidth/obj.W
	obj.W = newWidth
}
//Scales the Image in the Y dimension keeping the aspect ratio
func (obj *ImageObj) ScaleToY(newHeight float64) {
	obj.W *= newHeight/obj.H
	obj.H = newHeight
}
func (obj *ImageObj) GetFrame(thickness float64, alpha uint8) (frame *ImageObj) {
	frame = &ImageObj{X:obj.X, Y:obj.Y, W:obj.W, H:obj.H}
	
	frameImg, _ := ebiten.NewImage(int(obj.W), int(obj.H), ebiten.FilterDefault)
	frameImg.Fill(&color.RGBA{0,0,0,0})
	
	left := GetLineOfPoints(0,0, 0,obj.H, thickness)
	top := GetLineOfPoints(0,0, obj.W,0, thickness)
	right := GetLineOfPoints(obj.W,0, obj.W,obj.H, thickness)
	bottom := GetLineOfPoints(0,obj.H, obj.W,obj.H, thickness)
	col := &color.RGBA{0,0,0,alpha}
	
	left.Fill(frameImg, col)
	top.Fill(frameImg, col)
	right.Fill(frameImg, col)
	bottom.Fill(frameImg, col)
	frame.Img = frameImg
	
	return
}
//Draws the ImageObj on a screen
func (obj *ImageObj) DrawImageObj(screen *ebiten.Image) {
	obj.PanicIfNil()
	w, h := obj.Img.Size()
	op := &ebiten.DrawImageOptions{}
	xScale := obj.W/(float64)(w)
	yScale := obj.H/(float64)(h)

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
	op.GeoM.Scale(xScale, yScale)

	op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)
	screen.DrawImage(obj.Img, op)
}
//Draws a blured Image on a screen (box=3; alphaScale=(box*2+1)^2/2=25)
func (obj *ImageObj) DrawImageBlured(screen *ebiten.Image, box int, alphaScale float64) {
	obj.PanicIfNil()
	// https://en.wikipedia.org/wiki/Box_blur
	for j := -box; j <= box; j++ {
		for i := -box; i <= box; i++ {
			w, h := obj.Img.Size()
			op := &ebiten.DrawImageOptions{}
			xScale := obj.W/(float64)(w)
			yScale := obj.H/(float64)(h)
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
			op.GeoM.Scale(xScale, yScale)
			op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)
			
			op.GeoM.Translate(float64(i), float64(j))
			op.ColorM.Scale(1, 1, 1, alphaScale)
			screen.DrawImage(obj.Img, op)
		}
	}
}
//Draws the ImageObj with a certain transparency on a screen
func (obj *ImageObj) DrawImageObjAlpha(screen *ebiten.Image, alpha float64) {
	obj.PanicIfNil()
	w, h := obj.Img.Size()
	op := &ebiten.DrawImageOptions{}
	xScale := obj.W/(float64)(w)
	yScale := obj.H/(float64)(h)

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
	op.GeoM.Scale(xScale, yScale)

	op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)
	
	op.ColorM.Scale(1, 1, 1, alpha)
	screen.DrawImage(obj.Img, op)
}

//Loads an image.Image
func LoadImg(path string) (error, *image.Image) {
	f, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	img, format, err2 := image.Decode(f)
	if err2 != nil {
		return errors.New(fmt.Sprintf("%v, Format: %s",err2,format)), nil
	}
	return nil, &img
}
//Converts an image.Image to an ebiten.Image
func ImgToEbitenImg(img *image.Image) (*ebiten.Image) {	
	gophersImage, err3 := ebiten.NewImageFromImage(*img, ebiten.FilterDefault)
	if err3 != nil {
		panic(fmt.Sprintf("Cannot Create Ebiten Image: %v",err3))
	}
	return gophersImage
}
//Loads an ebiten.Image
func LoadEbitenImg(path string) (*ebiten.Image) {	
	err, img := LoadImg(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot Load Image: %v",err))
	}
	return ImgToEbitenImg(img)
}
func LoadEbitenImgFromBytes(im []byte) (*ebiten.Image) {
	img, _, err := image.Decode(bytes.NewReader(im))
	if err != nil {
		panic(err)
	}
	runnerImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return runnerImage
}
//Loads all Icons from a path with a list of sizes and a fileformat ("./64.png")
func InitIcons(path string, sizes []int, fileformat string) (error, []image.Image) {
	imgs := make([]image.Image, len(sizes))
	for i,_ := range(imgs) {
		err, img := LoadImg(fmt.Sprintf("%s/%v.%s", path, sizes[i], fileformat))
		if err != nil {
			return err, nil
		}
		imgs[i] = *img
	}
	return nil, imgs
}
func (obj *ImageObj) PanicIfNil() {
	if obj.Img == nil {
		panic("ImageObj is nil")
	}
}