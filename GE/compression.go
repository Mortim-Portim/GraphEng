package GE

import (
	"bytes"
	"io/ioutil"
	"compress/gzip"
	"encoding/binary"
)

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
//Converts  a []byte slice into a slice of ints
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
