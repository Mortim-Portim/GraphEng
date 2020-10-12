package GE

import (
	"math/big"
)


func CompressAll(changing [][]byte, static ...[]byte) (comp []byte) {
	comp = make([]byte, 0)
	for _,bs := range(static) {
		comp = append(comp, bs...)
	}
	for _,bs := range(changing) {
		i := &big.Int{}; i.SetInt64(int64(len(bs)))
		comp = append(comp, BigIntToBytes(i)...)
		comp = append(comp, bs...)
	}
	comp,_ = CompressBytes(comp)
	return
}

//static 0-n, changing 0-n
func DecompressAll(comp []byte, length []int) (BS [][]byte) {
	comp,_ = DecompressBytes(comp)
	BS = make([][]byte, 0)
	currIdx := 0
	for _,l := range(length) {
		BS = append(BS, comp[currIdx:currIdx+l])
		currIdx += l
	}
	
	for true {
		length := BytesToBigInt(comp[currIdx:currIdx+8])
		currIdx += 8
		BS = append(BS, comp[currIdx:currIdx+int(length.Int64())])
		currIdx += int(length.Int64())
		if currIdx >= len(comp) {
			break
		}
	}
	return
}