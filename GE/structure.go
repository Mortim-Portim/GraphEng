package GE

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Structure struct {
	squareSize int
	Collides bool
	HitboxW,HitboxH float64
	Name string
	understandable, IsUnderstood bool
	NUA, UA *DayNightAnim
	layer int8
}
/**
Params.txt:
spriteWidth: 	[1-NaN]
updatePeriod: 	[0-NaN]
hitBoxWidth:	[1-NaN]
hitBoxHeight:	[1-NaN]
Collides:		[true/false]
squareSize:		[1-NaN]
layer:			[-128-128]
understandable: [true/false]
**/
func GetStructFromParams(img *ebiten.Image, p *Params) (s *Structure) {
	collides := p.GetBool("Collides", false)
	layer :=	int8(p.Get("layer"))
	understandable := p.GetBool("understandable", false)
	spW := int(p.Get("spriteWidth"))
	uP := int(p.Get("updatePeriod"))
	var anim *DayNightAnim
	var u_Img *DayNightAnim
	if understandable {
		w, h := img.Size()
		ui := DeepCopyEbitenImage(img.SubImage(image.Rect(0,	0, w, h/2)).(*ebiten.Image))
		bi := DeepCopyEbitenImage(img.SubImage(image.Rect(0,  h/2, w,   h)).(*ebiten.Image))
		anim = GetDayNightAnim(1,1,1,1, spW, uP, ui)
		u_Img = GetDayNightAnim(1,1,1,1, spW, uP, bi)
	}else{
		anim = GetDayNightAnim(1,1,1,1, spW, uP, img)
		u_Img = nil
	}
	
	s = GetStructure(anim, u_Img, p.Get("hitBoxWidth")-1,p.Get("hitBoxHeight")-1, int(p.Get("squareSize")), collides, understandable, layer)
	return
}
//Returns a StructureObj
func GetStructure(NUA, UA *DayNightAnim, HitboxW,HitboxH float64, squareSize int, Collides, understandable bool, layer int8) (o *Structure) {
	o = &Structure{NUA:NUA, UA:UA, HitboxW:HitboxW ,HitboxH:HitboxH, squareSize:squareSize, Collides:Collides, understandable:understandable, layer:layer}
	if NUA != nil {
		o.NUA.Update(0)
	}
	if UA != nil {
		o.UA.Update(0)
	}
	return
}
func (s *Structure) Clone() *Structure {
	return &Structure{NUA:s.NUA.Copy(), UA:s.UA.Copy(), squareSize:s.squareSize, Collides:s.Collides, layer:s.layer, HitboxW:s.HitboxW, HitboxH:s.HitboxH, Name:s.Name}
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