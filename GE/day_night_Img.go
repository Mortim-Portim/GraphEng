package GE

import (
	"image"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

/**
DayNightImg represents two images that are drawn on top of each other
the lightlevel at which they are drawn defines how transparent the bright image is made, when drawing onto the dark one

DayNightImg implements UpdateAble
**/
type DayNightImg struct {
	day, night *ImageObj
}

func (i *DayNightImg) Loc() (float64, float64) {
	return i.day.X, i.day.Y
}
func (i *DayNightImg) Bds() (float64, float64) {
	return i.day.W, i.day.H
}
func (i *DayNightImg) Size() (w int, h int) {
	w, h = i.day.Img.Size()
	return
}
func (i *DayNightImg) Clone() *DayNightImg {
	return &DayNightImg{i.day.Copy(), i.night.Copy()}
}
func LoadDayImg(day_path string, width, height, x, y, angle float64) (*DayNightImg, error) {
	img := &DayNightImg{}
	di, err := LoadImgObj(day_path, width, height, x, y, angle)
	if err != nil {
		return nil, err
	}
	img.day = di
	img.night = di
	return img, nil
}
func LoadDayNightImg(path string, width, height, x, y, angle float64) (*DayNightImg, error) {
	img := &DayNightImg{}
	dn, err := LoadEbitenImg(path)
	if err != nil {
		return nil, err
	}
	w, h := dn.Size()
	img.day = &ImageObj{Img: dn.SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image), W: width, H: height, X: x, Y: y, Angle: angle}
	img.night = &ImageObj{Img: dn.SubImage(image.Rect(0, h/2, w, h)).(*ebiten.Image), W: width, H: height, X: x, Y: y, Angle: angle}
	return img, nil
}
func CreateDayNightImg(dn *ebiten.Image, width, height, x, y, angle float64) (img *DayNightImg) {
	img = &DayNightImg{}
	w, h := dn.Size()
	img.day = &ImageObj{Img: dn.SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image), W: width, H: height, X: x, Y: y, Angle: angle}
	img.night = &ImageObj{Img: dn.SubImage(image.Rect(0, h/2, w, h)).(*ebiten.Image), W: width, H: height, X: x, Y: y, Angle: angle}
	return
}

func (obj *ImageObj) CopyXYWHToDN(obj2 *DayNightImg) {
	obj.CopyXYWHTo(obj2.day)
	obj.CopyXYWHTo(obj2.night)
}

func (i *DayNightImg) SetParams(x, y, w, h float64) {
	if i.day == nil || i.night == nil {
		i.day = &ImageObj{}
		i.night = &ImageObj{}
	}
	i.day.X = x
	i.night.X = x
	i.day.Y = y
	i.night.Y = y
	i.day.W = w
	i.night.W = w
	i.day.H = h
	i.night.H = h
}

//Lightlevel 1: day, 0: night
func (i *DayNightImg) Draw(screen *ebiten.Image, lightlevel float64) {
	i.night.DrawImageObj(screen)
	i.day.DrawImageObjAlpha(screen, lightlevel)
}

//Sets the middle of the Image to x,y
func (i *DayNightImg) SetMiddle(x, y float64) {
	i.SetMiddleX(x)
	i.SetMiddleY(y)
}

//Sets the middle of the Image to x
func (i *DayNightImg) SetMiddleX(x float64) {
	i.day.X = x - i.day.W/2
	i.night.X = i.day.X
}

func (i *DayNightImg) SetRotation(angle float64) {
	i.day.Angle = angle
	i.night.Angle = angle
}
func (i *DayNightImg) Rotation() float64 {
	return i.day.Angle
}

//Sets the middle of the Image to y
func (i *DayNightImg) SetMiddleY(y float64) {
	i.day.Y = y - i.day.H/2
	i.night.Y = i.day.Y
}

//Puts the Image in between x1,y1 and x2,y2
func (i *DayNightImg) PutInBounds(x1, y1, x2, y2 float64) {
	i.day.X = x1
	i.day.Y = y1
	i.day.W = x2 - x1
	i.day.H = y2 - y1
	i.night.X = i.day.X
	i.night.Y = i.day.Y
	i.night.W = i.day.W
	i.night.H = i.day.H
}

//Copys the ImageObj
func (i *DayNightImg) Copy() *DayNightImg {
	return &DayNightImg{i.day.Copy(), i.night.Copy()}
}

//Scales the Original Image in a dimension keeping the aspect ratio
func (i *DayNightImg) ScaleOriginal(width, height float64) {
	scaledImg := resize.Resize(uint(width), uint(height), *i.day.OriginalImg, resize.NearestNeighbor)
	i.day.OriginalImg = &scaledImg
	i.day.Img = ImgToEbitenImg(&scaledImg)
	scaledImg2 := resize.Resize(uint(width), uint(height), *i.night.OriginalImg, resize.NearestNeighbor)
	i.night.OriginalImg = &scaledImg2
	i.night.Img = ImgToEbitenImg(&scaledImg2)
}
func (i *DayNightImg) ScaleToOriginalSize() {
	w, h := i.day.Img.Size()
	i.day.W = float64(w)
	i.day.H = float64(h)
	i.night.W = i.day.W
	i.night.H = i.day.H
}

//Scales the Image in a dimension keeping the aspect ratio
func (i *DayNightImg) ScaleDim(newval float64, dim int) {
	if dim == 0 {
		i.ScaleToX(newval)
	} else if dim == 1 {
		i.ScaleToY(newval)
	}
}

//Scales the Image in the X dimension keeping the aspect ratio
func (i *DayNightImg) ScaleToX(newWidth float64) {
	i.day.H *= newWidth / i.day.W
	i.day.W = newWidth
	i.night.W = i.day.W
	i.night.H = i.day.H
}

//Scales the Image in the Y dimension keeping the aspect ratio
func (i *DayNightImg) ScaleToY(newHeight float64) {
	i.day.W *= newHeight / i.day.H
	i.day.H = newHeight
	i.night.W = i.day.W
	i.night.H = i.day.H
}
func (obj *DayNightImg) SetDay(i *ebiten.Image) {
	obj.day.Img = i
}
func (obj *DayNightImg) SetNight(i *ebiten.Image) {
	obj.night.Img = i
}
func (obj *DayNightImg) GetDay() (i *ebiten.Image) {
	return obj.day.Img
}
func (obj *DayNightImg) GetNight() (i *ebiten.Image) {
	return obj.night.Img
}

func (obj *DayNightImg) GetDayIO() (i *ImageObj) {
	return obj.day
}
func (obj *DayNightImg) GetNightIO() (i *ImageObj) {
	return obj.night
}
