package GE

import (
	"image"
	"os"

	//_ "image/jpeg"
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/nfnt/resize"
)

func NewImageObj(img *image.Image, eimg *ebiten.Image, X, Y, W, H, angle float64) (iobj *ImageObj) {
	if img == nil {
		oimg := (image.Image)(eimg)
		img = &oimg
	}
	return &ImageObj{eimg, img, W, H, X, Y, angle}
}

//Loads the Image in Original Resolution
func LoadImgObj(path string, width, height, x, y, angle float64) (*ImageObj, error) {
	err, img := LoadImg(path)
	if err != nil {
		return nil, err
	}
	//scaledImg := resize.Resize(uint(width), uint(height), *img, resize.NearestNeighbor)
	eimg := ImgToEbitenImg(img)
	return &ImageObj{eimg, img, width, height, x, y, angle}, nil
}
func LoadImgObjFromBytes(bs []byte, width, height, x, y, angle float64) (*ImageObj, error) {
	eimg, err := LoadEbitenImgFromBytes(bs)
	if err != nil {
		return nil, err
	}
	img := (image.Image)(eimg)
	return &ImageObj{eimg, &img, width, height, x, y, angle}, nil
}
func EbitenImgToImgObj(img *ebiten.Image) *ImageObj {
	width, height := img.Size()
	oimg := (image.Image)(img)
	return &ImageObj{img, &oimg, float64(width), float64(height), 0, 0, 0}
}

//Stores the original, ebiten image and dimensions
type ImageObj struct {
	Img               *ebiten.Image
	OriginalImg       *image.Image
	W, H, X, Y, Angle float64
}

func (obj *ImageObj) Rectangle() *Rectangle {
	return GetRectangle(obj.X, obj.Y, obj.X+obj.W, obj.Y+obj.H)
}
func (obj *ImageObj) SetBottomRight(x, y float64) {
	obj.X = x - obj.W
	obj.Y = y - obj.H
}

//Sets the middle of the Image to x,y
func (obj *ImageObj) SetMiddle(x, y float64) {
	obj.X = x - obj.W/2
	obj.Y = y - obj.H/2
}
func (obj *ImageObj) GetMiddle() (float64, float64) {
	return obj.X + obj.W/2, obj.Y + obj.H/2
}

//Sets the middle of the Image to x
func (obj *ImageObj) SetMiddleX(x float64) {
	obj.X = x - obj.W/2
}

//Sets the middle of the Image to y
func (obj *ImageObj) SetMiddleY(y float64) {
	obj.Y = y - obj.H/2
}

//Puts the Image in between x1,y1 and x2,y2
func (obj *ImageObj) PutInBounds(x1, y1, x2, y2 float64) {
	obj.X = x1
	obj.Y = y1
	obj.W = x2 - x1
	obj.H = y2 - y1
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
	w, h := obj.Img.Size()
	obj.W = float64(w)
	obj.H = float64(h)
}

//Scales the Image in a dimension keeping the aspect ratio
func (obj *ImageObj) ScaleDim(newval float64, dim int) {
	if dim == 0 {
		obj.ScaleToX(newval)
	} else if dim == 1 {
		obj.ScaleToY(newval)
	}
}

//Scales the Image in the X dimension keeping the aspect ratio
func (obj *ImageObj) ScaleToX(newWidth float64) {
	obj.H *= newWidth / obj.W
	obj.W = newWidth
}

//Scales the Image in the Y dimension keeping the aspect ratio
func (obj *ImageObj) ScaleToY(newHeight float64) {
	obj.W *= newHeight / obj.H
	obj.H = newHeight
}

//Returns a frame of a given thickness and alpha value for the ImageObj
func (obj *ImageObj) GetFrame(thickness float64, alpha uint8, scale float64) (frame *ImageObj) {
	frame = &ImageObj{X: obj.X, Y: obj.Y, W: obj.W, H: obj.H}

	frameImg := ebiten.NewImage(int(obj.W), int(obj.H))
	frameImg.Fill(&color.RGBA{0, 0, 0, 0})

	left := GetLineOfPoints(0, 0, 0, obj.H, thickness)
	top := GetLineOfPoints(0, 0, obj.W, 0, thickness)
	right := GetLineOfPoints(obj.W, 0, obj.W, obj.H, thickness)
	bottom := GetLineOfPoints(0, obj.H, obj.W, obj.H, thickness)
	col := &color.RGBA{0, 0, 0, alpha}

	left.Fill(frameImg, col)
	top.Fill(frameImg, col)
	right.Fill(frameImg, col)
	bottom.Fill(frameImg, col)
	oImg := image.Image(frameImg)
	frame.OriginalImg = &oImg
	frame.ScaleOriginal(obj.W*scale, obj.H*scale)
	frame.ScaleToOriginalSize()
	return
}

//Draws the ImageObj on a screen
func (obj *ImageObj) DrawImageObj(screen *ebiten.Image) {
	obj.PanicIfNil()
	w, h := obj.Img.Size()
	op := &ebiten.DrawImageOptions{}
	xScale := obj.W / (float64)(w)
	yScale := obj.H / (float64)(h)

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
			xScale := obj.W / (float64)(w)
			yScale := obj.H / (float64)(h)
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
	xScale := obj.W / (float64)(w)
	yScale := obj.H / (float64)(h)

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(obj.Angle) * 2 * math.Pi / 360)
	op.GeoM.Scale(xScale, yScale)

	op.GeoM.Translate(obj.X+obj.W/2, obj.Y+obj.H/2)

	op.ColorM.Scale(1, 1, 1, alpha)
	screen.DrawImage(obj.Img, op)
}
func (obj *ImageObj) FillArea(screen *ebiten.Image, width, height float64) {
	oX := obj.X
	oY := obj.Y
	for obj.X+obj.W <= width {
		for obj.Y+obj.H <= height {
			obj.Draw(screen)
			obj.Y += obj.H
		}
		obj.Y = oY
		obj.X += obj.W
	}
	obj.X = oX
}
func (img *ImageObj) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return img.Update, img.Draw
}
func (img *ImageObj) Start(screen *ebiten.Image, data interface{}) {}

func (img *ImageObj) Stop(screen *ebiten.Image, data interface{}) {}

func (img *ImageObj) Update(frame int) {}

func (img *ImageObj) Draw(screen *ebiten.Image) {
	img.DrawImageObj(screen)
}
func DrawImageOnImage(dst, src *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

//Loads an image.Image
func LoadImg(path string) (error, *image.Image) {
	f, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	img, format, err2 := image.Decode(f)
	if err2 != nil {
		return errors.New(fmt.Sprintf("%v, Format: %s", err2, format)), nil
	}
	return nil, &img
}

//Converts an image.Image to an ebiten.Image
func ImgToEbitenImg(img *image.Image) *ebiten.Image {
	gophersImage := ebiten.NewImageFromImage(*img)
	return gophersImage
}

//Loads an ebiten.Image
func LoadEbitenImg(path string) (*ebiten.Image, error) {
	err, img := LoadImg(path)
	if err != nil {
		return nil, err
	}
	return ImgToEbitenImg(img), nil
}
func LoadEbitenImgFromBytes(im []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(im))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

//Returns an empty image
func GetEmptyImage(w, h int) (img *ebiten.Image) {
	img = ebiten.NewImage(w, h)
	img.Fill(color.RGBA{0, 0, 0, 0})
	return
}
func GetColoredImg(w, h int, col color.Color) (img *ebiten.Image) {
	img = ebiten.NewImage(w, h)
	img.Fill(col)
	return
}
func DeepCopyEbitenImage(img *ebiten.Image) (img2 *ebiten.Image) {
	w, h := img.Size()
	img2 = GetEmptyImage(w, h)
	op := &ebiten.DrawImageOptions{}
	img2.DrawImage(img, op)
	return
}

//Takes time
func SaveImage(path string, img image.Image) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	png.Encode(outputFile, img)

	// Don't forget to close files
	outputFile.Close()
	return nil
}

//Takes time
func SaveEbitenImage(path string, img *ebiten.Image) error {
	return SaveImage(path, (image.Image)(img))
}
func GetScaleOfEbitenImage(img *ebiten.Image) float64 {
	w, h := img.Size()
	return float64(w) / float64(h)
}

//Loads all Icons from a path with a list of sizes and a fileformat ("./64.png")
func InitIcons(path string, sizes []int, fileformat string) ([]image.Image, error) {
	imgs := make([]image.Image, len(sizes))
	for i := range imgs {
		err, img := LoadImg(fmt.Sprintf("%s/%v.%s", path, sizes[i], fileformat))
		if err != nil {
			return nil, err
		}
		imgs[i] = *img
	}
	return imgs, nil
}
func (obj *ImageObj) PanicIfNil() {
	if obj.Img == nil {
		panic("ImageObj is nil")
	}
}
func (obj *ImageObj) CopyXYWHTo(obj2 *ImageObj) {
	obj2.X = obj.X
	obj2.Y = obj.Y
	obj2.W = obj.W
	obj2.H = obj.H
}
