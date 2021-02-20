package GE

import (
	"fmt"
	"io/ioutil"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

const ERROR_WRONG_WORLD_VERSION = "Wrong World Version: %v"

const TO_BYTES_FUNC_VERSION = 1

//Converts the World into a []byte slice (VERSION 1)
func (p *WorldStructure) ToBytes() []byte {
	tilBs := p.TileMatToBytes()
	regBs := p.RegionMatToBytes()
	objBs := p.ObjectsToBytes()
	lghBs := p.LightsToBytes()
	changing := [][]byte{tilBs, objBs, lghBs, regBs}
	mdxBs, mdyBs, maxBs, timBs := p.StatsToBytes()
	return append([]byte{TO_BYTES_FUNC_VERSION}, cmp.CompressAll(changing, mdxBs, mdyBs, maxBs, timBs)...)
}

func SetWorldStructFromBytes(p *WorldStructure, data []byte) error {
	version := data[0]
	data = data[1:]
	fnc, ok := WorldStructureLoader[version]
	if !ok {
		return fmt.Errorf(ERROR_WRONG_WORLD_VERSION, version)
	}
	fnc(data, p)
	return nil
}

var WorldStructureLoader = map[byte]func([]byte, *WorldStructure) error{
	1: func(data []byte, p *WorldStructure) error {
		bs := cmp.DecompressAll(data, []int{8, 8, 2, 15})
		err := p.TileMatFromBytes(bs[4])
		if err != nil {
			return err
		}
		p.ResetMatrixesFromTileMat()
		p.xTilesAbs = p.TileMat.WAbs()
		p.yTilesAbs = p.TileMat.HAbs()
		p.SetDisplayWH(int(p.TileMat.Focus().Bounds().X)-2, int(p.TileMat.Focus().Bounds().Y)-2)
		p.BytesToObjects(bs[5])
		p.BytesToLights(bs[6])
		err = p.RegionMatFromBytes(bs[7])
		if err != nil {
			return err
		}
		p.StatsFromBytes(bs)
		return nil
	},
}

//fmt.Printf("Creating new World with w:%v, h:%v, fw:%v, fh:%v\n", tilMat.WAbs(), tilMat.HAbs(), int(tilMat.Focus().Bounds().X), int(tilMat.Focus().Bounds().Y))

//Saves the world in a highly compressed way to the file system
func (p *WorldStructure) Save(path string) error {
	return ioutil.WriteFile(path, p.ToBytes(), 0644)
}
func LoadWorldStructureFromBytes(X, Y, W, H float64, data []byte, tile_path, struct_path string) (*WorldStructure, error) {
	p := GetWorldStructure(X, Y, W, H, 1, 1, 1, 1)
	err := p.LoadTiles(tile_path)
	if err != nil {
		return p, err
	}
	err = p.LoadStructureObjs(struct_path)
	if err != nil {
		return p, err
	}
	err = SetWorldStructFromBytes(p, data)
	return p, err
}

//Loads the world from the file system ("./res/wrld.map", "./res/tiles", "./res/structs")
func LoadWorldStructure(X, Y, W, H float64, wrld_path, tile_path, struct_path string) (*WorldStructure, error) {
	p := GetWorldStructure(X, Y, W, H, 1, 1, 1, 1)
	data, err := ioutil.ReadFile(wrld_path)
	if err != nil {
		return p, err
	}
	return LoadWorldStructureFromBytes(X, Y, W, H, data, tile_path, struct_path)
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
func (p *WorldStructure) TileMatToBytes() []byte {
	bs, err := p.TileMat.Compress()
	CheckErr(err)
	return bs
}
func (p *WorldStructure) TileMatFromBytes(bs []byte) error {
	p.TileMat = GetMatrix(0, 0, 0)
	return p.TileMat.Decompress(bs)
}
func (p *WorldStructure) RegionMatToBytes() []byte {
	bs, err := p.RegionMat.Compress()
	CheckErr(err)
	return bs
}
func (p *WorldStructure) RegionMatFromBytes(bs []byte) error {
	p.RegionMat = GetMatrix(0, 0, 0)
	return p.RegionMat.Decompress(bs)
}

//(8),(8),(2),(15)
func (p *WorldStructure) StatsToBytes() (mdxBs, mdyBs, maxBs, timBs []byte) {
	mdxBs = cmp.Int64ToBytes(int64(p.middleX))
	mdyBs = cmp.Int64ToBytes(int64(p.middleY))
	maxBs = cmp.Int16ToBytes(int16(p.maxLightLevel))
	var err error
	timBs, err = p.CurrentTime.MarshalBinary()
	CheckErr(err)
	return
}
func (p *WorldStructure) StatsFromBytes(bs [][]byte) {
	p.SetMiddle(int(cmp.BytesToInt64(bs[0])), int(cmp.BytesToInt64(bs[1])), true)
	p.maxLightLevel = int16(cmp.BytesToInt16(bs[2]))
	p.CurrentTime.UnmarshalBinary(bs[3])
}
