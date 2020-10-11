package GE

import (
	"io/ioutil"
)

//Converts the World into a []byte slice
func (p *WorldStructure) ToBytes() ([]byte, error) {
	idxBs, err1 := p.TileMat.Compress()
	if err1 != nil {return nil, err1}
	
	ligBs, err2 := p.LIdxMat.Compress()
	if err2 != nil {return nil, err2}
	mats := append(idxBs, ligBs...)
	
	bs, err3 := CompressBytes(append(mats, AppendInt16ToBytes( int16(len(ligBs)), Int16ToBytes(int16(len(idxBs))) )...))
	if err3 != nil {return nil, err3}
	return bs, nil
}
//Converts a []byte slice into a WorldStructure
func (p *WorldStructure) FromBytes(data []byte) error {
	bs, err2 := DecompressBytes(data)
   	if err2 != nil {return err2}
   	
   	lenIdx := BytesToInt16(bs[len(bs)-4:len(bs)-2])
   	lenLig := BytesToInt16(bs[len(bs)-2:len(bs)])
   	
   	err3 := p.TileMat.Decompress(bs[:lenIdx])
   	if err3 != nil {return err3}
   	
   	err4 := p.LIdxMat.Decompress(bs[lenIdx:lenIdx+lenLig])
   	if err4 != nil {return err4}
   	return nil
}

//Saves the world in a highly compressed way to the file system
func (p *WorldStructure) Save(path string) error {
	bs, err := p.ToBytes()
	if err != nil {return err}
	return ioutil.WriteFile(path, bs, 0644)
}
//Loads the world from the file system
func (p *WorldStructure) Load(path string) error {
	data, err1 := ioutil.ReadFile(path)
   	if err1 != nil {return err1}
   	err2 := p.FromBytes(data)
   	if err2 != nil {return err2}
   	return nil
}