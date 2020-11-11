package Compression

import (
	"bytes"
	"io/ioutil"
	"compress/gzip"
	"encoding/binary"
	"math/big"
)

func BoolToByte(b bool) byte {
	if b {
		return byte(0)
	}else{
		return byte(1)
	}
}
func ByteToBool(b byte) (bool) {
	return b == 0
}
//Converts an float64 into a [8]byte array
func Float64ToBytes(f float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, f)
	if err != nil {panic(err)}
	return buf.Bytes()
}
//Converts a [8]byte array into a float64
func BytesToFloat64(bs []byte) (f float64) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(bs)
	if err != nil {panic(err)}
	err2 := binary.Read(buf, binary.LittleEndian, &f)
	if err2 != nil {panic(err2)}
	return
}
func Float64sToBytes(fs ...float64) (bs []byte) {
	lenFs := len(fs)
	done := make(chan bool)
	bs = make([]byte, lenFs*8)
	for i,f := range(fs) {
		start := i*8
		val := f
		go func(){
			copy(bs[start:start+8], Float64ToBytes(val))
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<- done
	}
	return
}
func BytesToFloat64s(bs []byte) (fs []float64) {
	lenFs := len(bs)/8
	done := make(chan bool)
	fs = make([]float64, lenFs)
	for i,_ := range(fs) {
		start := i*8
		idx := i
		go func(){
			fs[idx] = BytesToFloat64(bs[start:start+8])
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<- done
	}
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
func Int64sToBytes(is ...int64) (bs []byte) {
	lenFs := len(is)
	done := make(chan bool)
	bs = make([]byte, lenFs*8)
	for i,f := range(is) {
		start := i*8
		val := f
		go func(){
			copy(bs[start:start+8], Int64ToBytes(val))
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<- done
	}
	return
}
func BytesToInt64s(bs []byte) (is []int64) {
	lenFs := len(bs)/8
	done := make(chan bool)
	is = make([]int64, lenFs)
	for i,_ := range(is) {
		start := i*8
		idx := i
		go func(){
			is[idx] = BytesToInt64(bs[start:start+8])
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<- done
	}
	return
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
//Converts a 2d slice of ints into a []byte slice
func Int16s2DToBytes(is [][]int16) (bs []byte) {
	bs = make([]byte, 0)
	bs = append(bs, Int64ToBytes(int64(len(is[0])))...)
	
	for _,il := range(is) {
		bs = append(bs, Int16sToBytes(il)...)
	}
	
	return
}
//Converts a []byte slice into a 2d slice of ints
func BytesToInt16s2D(bs []byte) (is [][]int16) {
	is = make([][]int16, 0)
	xl := int(BytesToInt64(bs[:8]))*2
	bs = bs[8:]
	
	for true {
		bts := bs[:xl]
		is = append(is, BytesToInt16s(bts))
		if xl >= len(bs) {
			break
		}
		bs = bs[xl:]
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
	if len(bs) <= 0 {
		return bs, nil
	}
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
	if len(bs) <= 0 {
		return bs, nil
	}
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
