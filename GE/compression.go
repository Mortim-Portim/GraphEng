package GE

import (
	"bytes"
	"io/ioutil"
	"compress/gzip"
	"encoding/binary"
	"math/big"
)
func Float64ToBytes(f float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, f)
	if err != nil {panic(err)}
	return buf.Bytes()
}
func BytesToFloat64(bs []byte) (f float64) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(bs)
	if err != nil {panic(err)}
	err2 := binary.Read(buf, binary.LittleEndian, &f)
	if err2 != nil {panic(err2)}
	return
}
//Converts an int64 into a [8]byte array
func Int64ToBytes(i int64) (b []byte) {
	bi := &big.Int{}; bi.SetInt64(i)
	b = make([]byte, 8)
	bi.FillBytes(b)
	return
}
//Converts a [8]byte array into an int64
func BytesToInt64(b []byte) (int64) {
	i := &big.Int{}
	i.SetBytes(b)
	return i.Int64()
}
//converts an int64 into a [8]byte array
func BigIntToBytes(i *big.Int) (b []byte) {
	b = make([]byte, 8)
	i.FillBytes(b)
	return
}
//converts an int64 into a [8]byte array
func BytesToBigInt(b []byte) (*big.Int) {
	i := &big.Int{}
	i.SetBytes(b)
	return i
}
//Converts an int16 into a [2]byte array
func Int16ToBytes(i int16) (b []byte) {
	b = make([]byte, 2)
	ui := uint16(int64(i)+32768)
	binary.LittleEndian.PutUint16(b, ui)
	return
}
//Converts a [2]byte array into an int16
func BytesToInt16(b []byte) (int16) {
	return int16(int64(binary.LittleEndian.Uint16(b))-32768)
}
//Converts a slice of ints into a []byte slice
func Int16sToBytes(is []int16) (bs []byte) {
	bs = make([]byte, 2*len(is))
	for idx,i := range(is) {
		b := Int16ToBytes(i)
		bs[idx*2] = b[0]
		bs[idx*2+1] = b[1]
	}
	return
}
//Converts a []byte slice into a slice of ints
func BytesToInt16s(bs []byte) (is []int16) {
	is = make([]int16, len(bs)/2)
	for idx,_ := range(bs) {
		if idx%2 == 0 {
			is[idx/2] = BytesToInt16(bs[idx:idx+2])
		}
	}
	return
}
//Appends an int16 to a []byte slice
func AppendInt16ToBytes(i int16, bs []byte) []byte {
	ib := Int16ToBytes(i)
	bs = append(bs, ib[0]);bs = append(bs, ib[1])
	return bs
}

//Compresses a []byte slice using gzip
func CompressBytes(bs []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err0 := gzip.NewWriterLevel(&b, gzip.DefaultCompression)
	if err0 != nil {
		return nil, err0
	}
	_, err1 := w.Write(bs)
	if err1 != nil {
		return nil, err1
	}
	err2 := w.Close()
	if err2 != nil {
		return nil, err2
	}
	return b.Bytes(), nil
}
//Decompresses a []byte slice using gzip
func DecompressBytes(bs []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err1 := b.Write(bs)
	if err1 != nil {
		return nil, err1
	}
	r,err2 := gzip.NewReader(&b)
	if err2 != nil {
		return nil, err2
	}
	content, err3 := ioutil.ReadAll(r)
	if err3 != nil {
		return nil, err3
	}
	err4 := r.Close()
	if err4 != nil {
		return nil, err4
	}
	return content, nil
}
