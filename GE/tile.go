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
	ts := make([]*Tile, len(files))
    for i, f := range files {
            path := folderPath+f.Name()
            img, err := LoadEbitenImg(path)
            if err != nil {
	            return ts, err
            }
            ts[i] = &Tile{img, f.Name()[:len(f.Name())-4]}
    }
    return ts, nil
}


type Tile struct {
	Img *ebiten.Image
	Name string
}

func (t *Tile) Draw(screen *ebiten.Image, x, y, w, h float64, myLayer, drawLayer int, drawer, frame *ImageObj) {
	drawer.Img = t.Img
	drawer.X = x
	drawer.Y = y
	if myLayer == drawLayer {
		drawer.DrawImageObj(screen)
	} else if myLayer < drawLayer {
		dif := 1.0 / float64(drawLayer-myLayer+1)
		drawer.DrawImageObjAlpha(screen, dif)
	} else {
		//box := float64(1 + myLayer - drawLayer)
		//sq := box*2 + 1
		//drawer.DrawImageBlured(screen, int(box), 1.0/((sq*sq)*0.35))
		drawer.DrawImageObj(screen)
	}
	if frame != nil {
		frame.X = drawer.X
		frame.Y = drawer.Y
		frame.DrawImageObj(screen)
	}
}