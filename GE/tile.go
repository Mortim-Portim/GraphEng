package GE

import (
	"github.com/hajimehoshi/ebiten"
)

const INDEXFILENAME = "#index.txt"

/**
#index.txt:
tree
house


folderPath
	--> index.txt
	--> tree.png
	--> house.png
**/
//Reads all tiles from a directory
func ReadTiles(folderPath string) ([]*Tile, error) {
	if folderPath[len(folderPath)-1:] != "/" {
		folderPath += "/"
	}
	idxList := &List{}; idxList.LoadFromFile(folderPath+INDEXFILENAME)
	slc := idxList.GetSlice()
	ts := make([]*Tile, len(slc))
	names := make([]string, 0)
	for i,name := range(slc) {
	        if !containsS(names, name) {
				names = append(names, name)
				DNImg, err := LoadDayNightImg(folderPath+name+".png",0,0,0,0,0)
				if err != nil {
					return nil, err
				}
				DNImg.ScaleToOriginalSize()
				ts[i] = &Tile{DNImg, name}
            }
    }
    return ts, nil
}

/**
Tile represents a drawable struct that is one tile on the map
**/
type Tile struct {
	Img *DayNightImg
	Name string
}

func (t *Tile) Draw(screen *ebiten.Image, drawer *ImageObj, lightlevel int16) {
	drawer.CopyXYWHToDN(t.Img)
	alpha := float64(lightlevel)/float64(255)
	t.Img.Draw(screen, alpha)
}

/**
//Reads all tiles from a directory
func ReadTiles(folderPath string) ([]*Tile, error) {
	files, err1 := ioutil.ReadDir(folderPath)
    if err1 != nil {
    	return nil, err1
    }
	ts := make([]*Tile, 0)
	names := make([]string, 0)
    for _, f := range files {
            name := f.Name()[:len(f.Name())-4]
	        if !containsS(names, name) {
				names = append(names, name)
				DNImg := LoadDayNightImg(folderPath+name+".png",0,0,0,0,0)
				DNImg.ScaleToOriginalSize()
				ts = append(ts, &Tile{DNImg, name})
            }
    }
    return ts, nil
}
**/