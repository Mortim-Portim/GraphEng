package GE

import (
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
)

type Structure struct {
	*DayNightAnim
	squareSize int
	Collides, Background bool
	HitboxW,HitboxH float64
	Name string
}
/**
Params.txt:
spriteWidth: 	[1-NaN]
updatePeriod: 	[0-NaN]
hitBoxWidth:	[1-NaN]
hitBoxHeight:	[1-NaN]
Collides:		[true/false]
squareSize:		[1-NaN]
Background:		[true/false]
**/
func GetStructFromParams(img *ebiten.Image, p *Params) (s *Structure) {
	anim := GetDayNightAnim(1,1,1,1, int(p.Get("spriteWidth")), int(p.Get("updatePeriod")), img)
	collides := false
	if p.GetS("Collides") != "false" {
		collides = true
	}
	Background := false
	if p.GetS("Background") != "false" {
		Background = true
	}
	s = GetStructure(anim, p.Get("hitBoxWidth")-1,p.Get("hitBoxHeight")-1, int(p.Get("squareSize")), collides, Background)
	return
}
//Returns a StructureObj
func GetStructure(anim *DayNightAnim, HitboxW,HitboxH float64, squareSize int, Collides, Background bool) (o *Structure) {
	o = &Structure{DayNightAnim:anim, HitboxW:HitboxW ,HitboxH:HitboxH, squareSize:squareSize, Collides:Collides, Background:Background}
	o.Update(0)
	return
}
func (s *Structure) Clone() *Structure {
	return &Structure{DayNightAnim:s.DayNightAnim.Clone(), squareSize:s.squareSize, Collides:s.Collides, Background:s.Background, HitboxW:s.HitboxW, HitboxH:s.HitboxH, Name:s.Name}
}
/**
Reads a slice of Structures from a folder like this:
folder
----> obj1.png
----> obj1.txt
----> obj2.png
----> obj2.txt
**/
func ReadStructures(folderPath string) ([]*Structure, error) {
	files, err1 := ioutil.ReadDir(folderPath)
    if err1 != nil {
    	return nil, err1
    }
	ts := make([]*Structure, 0)
	names := make([]string, 0)
    for _, f := range files {
            name := f.Name()[:len(f.Name())-4]
	        if !containsS(names, name) {
				names = append(names, name)
				img, err := LoadEbitenImg(folderPath+name+".png")
				if err != nil {
					return nil, err
				}
				p := &Params{}
				err2 := p.LoadFromFile(folderPath+name+".txt")
				if err2 != nil {
					return nil, err2
				}
				obj := GetStructFromParams(img, p)
				obj.Name = name
				ts = append(ts, obj)
            }
    }
    return ts, nil
}