package GE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/nfnt/resize"
)

type DayNightImg struct {
	day, night *ImageObj
}

func (i *DayNightImg) Size() (w int, h int) {
	w,h = i.day.Img.Size()
	return 
}

func LoadDayImg(day_path string, width, height, x, y, angle float64) (img *DayNightImg) {
	img = &DayNightImg{}
	img.day = LoadImgObj(day_path, width, height, x, y, angle)
	img.night = img.day
	return
}
func LoadDayNightImg(day_path, night_path string, width, height, x, y, angle float64) (img *DayNightImg) {
	img = &DayNightImg{}
	img.day = LoadImgObj(day_path, width, height, x, y, angle)
	img.night = LoadImgObj(night_path, width, height, x, y, angle)
	return
}

func (obj *ImageObj) CopyXYWHToDN(obj2 *DayNightImg) {
	obj.CopyXYWHTo(obj2.day)
	obj.CopyXYWHTo(obj2.night)
}

func (i *DayNightImg) SetParams(x,y,w,h float64) {
	if i.day == nil || i.night == nil {
		i.day = &ImageObj{}
		i.night = &ImageObj{}
	}
	i.day.X = x; i.night.X = x
	i.day.Y = y; i.night.Y = y
	i.day.W = w; i.night.W = w
	i.day.H = h; i.night.H = h
}

//Lightlevel 0: day, 1: night
func (i *DayNightImg) Draw(screen *ebiten.Image, lightlevel float64) {
	i.day.DrawImageObj(screen)
	i.night.DrawImageObjAlpha(screen, lightlevel)
}
//Sets the middle of the Image to x,y
func (i *DayNightImg) SetMiddle(x,y float64) {
	i.SetMiddleX(x)
	i.SetMiddleY(y)
}
//Sets the middle of the Image to x
func (i *DayNightImg) SetMiddleX(x float64) {
	i.day.X = x-i.day.W/2
	i.night.X = i.day.X
}
//Sets the middle of the Image to y
func (i *DayNightImg) SetMiddleY(y float64) {
	i.day.Y = y-i.day.H/2
	i.night.Y = i.day.Y
}
//Puts the Image in between x1,y1 and x2,y2
func (i *DayNightImg) PutInBounds(x1,y1, x2,y2 float64) {
	i.day.X = x1; i.day.Y = y1
	i.day.W = x2-x1; i.day.H = y2-y1
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
	i.day.Img,_ = ImgToEbitenImg(&scaledImg)
	scaledImg2 := resize.Resize(uint(width), uint(height), *i.night.OriginalImg, resize.NearestNeighbor)
	i.night.OriginalImg = &scaledImg2
	i.night.Img,_ = ImgToEbitenImg(&scaledImg2)
}
func (i *DayNightImg) ScaleToOriginalSize() {
	w,h := i.day.Img.Size()
	i.day.W = float64(w)
	i.day.H = float64(h)
	i.night.W = i.day.W
	i.night.H = i.day.H
}
//Scales the Image in a dimension keeping the aspect ratio
func (i *DayNightImg) ScaleDim(newval float64, dim int) {
	if dim == 0 {
		i.day.ScaleToX(newval)
	}else if dim == 1 {
		i.day.ScaleToY(newval)
	}
}
//Scales the Image in the X dimension keeping the aspect ratio
func (i *DayNightImg) ScaleToX(newWidth float64) {
	i.day.H *= newWidth/i.day.W
	i.day.W = newWidth
	i.night.W = i.day.W
	i.night.H = i.day.H
}
//Scales the Image in the Y dimension keeping the aspect ratio
func (i *DayNightImg) ScaleToY(newHeight float64) {
	i.day.W *= newHeight/i.day.H
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