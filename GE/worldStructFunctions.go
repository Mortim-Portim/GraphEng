package GE

import (
	"fmt"
	"io/ioutil"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

//Converts the World into a []byte slice
func (p *WorldStructure) ToBytes() ([]byte, error) {
	tilBs, err1 := p.TileMat.Compress()
	if err1 != nil {
		return nil, err1
	}
	regBs, err2 := p.RegionMat.Compress()
	if err2 != nil {
		return nil, err2
	}
	objBs := p.ObjectsToBytes()
	lghBs := p.LightsToBytes()
	changing := append([][]byte{tilBs}, objBs, lghBs, regBs)

	mdxBs := cmp.Int64ToBytes(int64(p.middleX))
	mdyBs := cmp.Int64ToBytes(int64(p.middleY))
	minBs := cmp.Int16ToBytes(int16(p.minLight))
	maxBs := cmp.Int16ToBytes(int16(p.maxLight))
	delBs := cmp.Float64ToBytes(p.deltaB)
	return cmp.CompressAll(changing, mdxBs, mdyBs, minBs, maxBs, delBs), nil
}

//Converts a []byte slice into a WorldStructure
func GetWorldStructureFromBytes(X, Y, W, H float64, data []byte, tile_path, struct_path string) (*WorldStructure, error) {
	bs := cmp.DecompressAll(data, []int{8, 8, 2, 2, 8})
	tilMat := GetMatrix(0, 0, 0)
	err := tilMat.Decompress(bs[5])
	if err != nil {
		return nil, err
	}

	fmt.Printf("Creating new World with w:%v, h:%v, fw:%v, fh:%v\n", tilMat.WAbs(), tilMat.HAbs(), int(tilMat.Focus().Bounds().X), int(tilMat.Focus().Bounds().Y))
	p := GetWorldStructure(X, Y, W, H, tilMat.WAbs(), tilMat.HAbs(), int(tilMat.Focus().Bounds().X), int(tilMat.Focus().Bounds().Y))
	p.TileMat = tilMat
	err = p.LoadTiles(tile_path)
	if err != nil {
		return nil, err
	}
	err = p.LoadStructureObjs(struct_path)
	if err != nil {
		return nil, err
	}
	p.BytesToObjects(bs[6])
	p.BytesToLights(bs[7])
	regMat := GetMatrix(0, 0, 0)
	if len(bs) >= 9 {
		err = regMat.Decompress(bs[8])
		if err != nil {return nil, err}
	}
	p.RegionMat = regMat
	p.SetMiddle(int(cmp.BytesToInt64(bs[0])), int(cmp.BytesToInt64(bs[1])), true)
	p.SetLightStats(int16(cmp.BytesToInt16(bs[2])), int16(cmp.BytesToInt16(bs[3])), cmp.BytesToFloat64(bs[4]))
	return p, nil
}

//Saves the world in a highly compressed way to the file system
func (p *WorldStructure) Save(path string) error {
	bs, err := p.ToBytes()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, 0644)
}

//Loads the world from the file system ("./res/wrld.map", "./res/tiles", "./res/structs")
func LoadWorldStructure(X, Y, W, H float64, wrld_path, tile_path, struct_path string) (*WorldStructure, error) {
	data, err1 := ioutil.ReadFile(wrld_path)
	if err1 != nil {
		return nil, err1
	}
	p, err2 := GetWorldStructureFromBytes(X, Y, W, H, data, tile_path, struct_path)
	if err2 != nil {
		return nil, err2
	}
	return p, nil
}
func (p *WorldStructure) LoadTiles(path string) error {
	//loads all tiles
	tiles, errT := ReadTiles(path)
	if errT != nil {
		return errT
	}
	p.Tiles = tiles
	return nil
}
func (p *WorldStructure) LoadStructureObjs(path string) error {
	//loads all objs
	objs, errO := ReadStructures(path)
	if errO != nil {
		return errO
	}
	p.AddStruct(objs...)
	return nil
}
func (p *WorldStructure) AddNamedStructureObj(name string, x, y float64) {
	p.AddStructObj(GetStructureObj(p.GetNamedStructure(name), x, y))
}
