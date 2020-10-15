package GE

import (
	"github.com/hajimehoshi/ebiten"
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

#index.txt:
obj1
obj2

folderPath
	----> index.txt
	----> obj1.png
	----> obj1.txt
	----> obj2.png
	----> obj2.txt
**/
func ReadStructures(folderPath string) ([]*Structure, error) {
	if folderPath[len(folderPath)-1:] != "/" {
		folderPath += "/"
	}
	idxList := &List{}; idxList.LoadFromFile(folderPath+INDEXFILENAME)
	slc := idxList.GetSlice()
	
	ts := make([]*Structure, len(slc))
	names := make([]string, 0)
    for i,name := range(slc) {
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
				ts[i] = obj
            }
    }
    return ts, nil
}