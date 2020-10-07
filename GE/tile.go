package GE

import (
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
)

/**
Tile represents a drawable struct that is one tile on the map
**/

//Reads all tiles from a directory
func ReadTiles(folderPath string) ([]*Tile, error) {
	files, err1 := ioutil.ReadDir(folderPath)
    if err1 != nil {
    	return nil, err1
    }
	ts := make([]*Tile, 0)
	names := make([]string, 0)
    for _, f := range files {
            name := f.Name()[:len(f.Name())-6]
	        if !containsS(names, name) {
				names = append(names, name)
				DNImg := LoadDayNightImg(folderPath+name+"_D.png", folderPath+name+"_N.png",0,0,0,0,0)
				DNImg.ScaleToOriginalSize()
				ts = append(ts, &Tile{DNImg, name})
            }
    }
    return ts, nil
}


type Tile struct {
	Img *DayNightImg
	Name string
}

func (t *Tile) Draw(screen *ebiten.Image, drawer, frame *ImageObj, lightlevel uint8) {
	drawer.CopyXYWHToDN(t.Img)
	alpha := float64(255-lightlevel)/float64(255)
	t.Img.Draw(screen, alpha)
	if frame != nil {
		frame.X = drawer.X
		frame.Y = drawer.Y
		frame.DrawImageObj(screen)
	}
}