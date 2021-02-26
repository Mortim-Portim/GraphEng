package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"math"
	"math/big"
)

var bitvals = []int{1, 2, 4, 8, 16, 32, 64, 128}

//BoolsToBytes (bls ...bool) (bs []byte) - puts up to 8 bools in one byte
func BoolsToBytes(bls ...bool) (bs []byte) {
	bs = []byte{}
	for i := 0; i < len(bls); i += 8 {
		b := 0
		for i2 := 0; i2 < 8; i2++ {
			if i+i2 < len(bls) && bls[i+i2] {
				b += bitvals[i2]
			}
		}
		bs = append(bs, byte(b))
	}
	return
}

//BytesToBools (bs []byte) (bls []bool) - returns a slice with a lenght of multiples of 8, standard value being false
func BytesToBools(bs []byte) (bls []bool) {
	bls = make([]bool, len(bs)*8)
	for i, b := range bs {
		for i2 := 0; i2 < 8; i2++ {
			bls[i*8+i2] = (int(b) & bitvals[i2]) == bitvals[i2]
		}
	}
	return
}

//Float2DToBytesRound (f1, f2 float64) (bs []byte) - converts two float64 to byte[6]
func Float2DToBytesRound(f1, f2 float64) (bs []byte) {
	bs = append(FloatToBytesRound(f1), FloatToBytesRound(f2)...)
	return
}

//BytesToFloat2DRound (bs []byte) (f1, f2 float64) - converts byte[6] to float64
func BytesToFloat2DRound(bs []byte) (f1, f2 float64) {
	f1 = BytesToFloatRound(bs[0:3])
	f2 = BytesToFloatRound(bs[3:6])
	return
}

//FloatToBytesRound (f float64) (bs []byte) - Converts a float64 in [0;2^16] with a precision of 1/127 to byte[3]
func FloatToBytesRound(f float64) (bs []byte) {
	ff := math.Round(f)
	fr := (f - ff) * 255
	bs = UInt16ToBytes(uint16(ff))
	bs = append(bs, Int8ToBytes(int8(fr))...)
	return
}

//BytesToFloatRound (bs []byte) (f float64) - Converts byte[3] into a float64 in [0;2^16] with a precision of 1/127
func BytesToFloatRound(bs []byte) (f float64) {
	f = float64(BytesToUInt16(bs[0:2]))
	f += float64(BytesToInt8(bs[2:3])) / 255
	return
}

//FloatToBytesRoundFP (f float64) (bs []byte) - Converts a float64 in [0;255] with a precision of 1/65535 to byte[3]
func FloatToBytesRoundFP(f float64) (bs []byte) {
	ff := math.Floor(f)
	fr := int16((f - ff) * 65535)
	bs = append(Int16ToBytes(fr), byte(ff))
	return
}

//BytesToFloatRoundFP (bs []byte) (f float64) - Converts byte[3] into a float64 in [0;255] with a precision of 1/65535
func BytesToFloatRoundFP(bs []byte) (f float64) {
	f = float64(BytesToInt16(bs[0:2])) / 65535
	f += float64(bs[2])
	return
}

//UInt16ToBytes (i uint16) (b []byte) - converts a uint16 to [2]byte
func UInt16ToBytes(i uint16) (b []byte) {
	b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return
}

//BytesToUInt16 (b []byte) uint16 - converts a [2]byte to uint16
func BytesToUInt16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

//UInt64ToBytes (i uint64) (b []byte) - converts a uint16 to [2]byte
func UInt64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 2)
	binary.LittleEndian.PutUint64(b, i)
	return
}

//BytesToUInt64 (b []byte) uint64 - converts a [2]byte to uint16
func BytesToUInt64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

//Int8ToBytes (i int8) []byte - converts a int8 to [1]byte
func Int8ToBytes(i int8) []byte {
	return []byte{byte(int(i) + 128)}
}

//BytesToInt8 (b []byte) int8 - converts a [1]byte to int8
func BytesToInt8(b []byte) int8 {
	return int8(int(b[0]) - 128)
}

//BoolToByte (b bool) byte - converts one bool into one byte
func BoolToByte(b bool) byte {
	if b {
		return byte(0)
	}
	return byte(1)
}

//ByteToBool (b byte) bool - converts one byte into one bool
func ByteToBool(b byte) bool {
	return b == 0
}

//Float32ToBytes (f float32) []byte - Converts an float64 into a [4]byte slice
func Float32ToBytes(f float32) []byte {
	return GetBytesOfVar(f)
}

//BytesToFloat32 (bs []byte) (f float32) - Converts a [4]byte slice into a float32
func BytesToFloat32(bs []byte) (f float32) {
	GetVarFromBytes(bs, &f)
	return
}

//Float32sToBytes (fs ...float32) (bs []byte) - converts multiple float32 into bytes
func Float32sToBytes(fs ...float32) (bs []byte) {
	bs = make([]byte, len(fs)*4)
	for i, f := range fs {
		copy(bs[i*4:i*4+4], Float32ToBytes(f))
	}
	return
}

//BytesToFloat32s (bs []byte) (fs []float32) - converts bytes into multiple float32
func BytesToFloat32s(bs []byte) (fs []float32) {
	fs = make([]float32, len(bs)/4)
	for i := range fs {
		fs[i] = BytesToFloat32(bs[i*4 : i*4+4])
	}
	return
}

//Float64ToBytes (f float64) []byte - Converts an float64 into a [8]byte array
func Float64ToBytes(f float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, f)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//BytesToFloat64 (bs []byte) (f float64) - Converts a [8]byte array into a float64
func BytesToFloat64(bs []byte) (f float64) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(bs)
	if err != nil {
		panic(err)
	}
	err2 := binary.Read(buf, binary.LittleEndian, &f)
	if err2 != nil {
		panic(err2)
	}
	return
}

//Float64sToBytes (fs ...float64) (bs []byte) - converts multiple float64 into bytes
func Float64sToBytes(fs ...float64) (bs []byte) {
	lenFs := len(fs)
	done := make(chan bool)
	bs = make([]byte, lenFs*8)
	for i, f := range fs {
		start := i * 8
		val := f
		go func() {
			copy(bs[start:start+8], Float64ToBytes(val))
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<-done
	}
	return
}

//BytesToFloat64s (bs []byte) (fs []float64) - converts bytes into multiple float64
func BytesToFloat64s(bs []byte) (fs []float64) {
	lenFs := len(bs) / 8
	done := make(chan bool)
	fs = make([]float64, lenFs)
	for i := range fs {
		start := i * 8
		idx := i
		go func() {
			fs[idx] = BytesToFloat64(bs[start : start+8])
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<-done
	}
	return
}

//Int64ToBytes (i int64) (b []byte) - Converts an int64 into a [8]byte array
func Int64ToBytes(i int64) (b []byte) {
	bi := &big.Int{}
	bi.SetInt64(i)
	b = make([]byte, 8)
	bi.FillBytes(b)
	return
}

//BytesToInt64 (b []byte) int64 - Converts a [8]byte array into an int64
func BytesToInt64(b []byte) int64 {
	i := &big.Int{}
	i.SetBytes(b)
	return i.Int64()
}

//Int64sToBytes (is ...int64) (bs []byte) - converts multiple int64 into bytes
func Int64sToBytes(is ...int64) (bs []byte) {
	lenFs := len(is)
	done := make(chan bool)
	bs = make([]byte, lenFs*8)
	for i, f := range is {
		start := i * 8
		val := f
		go func() {
			copy(bs[start:start+8], Int64ToBytes(val))
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<-done
	}
	return
}

//BytesToInt64s (bs []byte) (is []int64) - converts bytes into multiple int64
func BytesToInt64s(bs []byte) (is []int64) {
	lenFs := len(bs) / 8
	done := make(chan bool)
	is = make([]int64, lenFs)
	for i := range is {
		start := i * 8
		idx := i
		go func() {
			is[idx] = BytesToInt64(bs[start : start+8])
			done <- true
		}()
	}
	for i := 0; i < lenFs; i++ {
		<-done
	}
	return
}

//BigIntToBytes (i *big.Int) (b []byte) - converts a big int into a [8]byte array
func BigIntToBytes(i *big.Int) (b []byte) {
	b = make([]byte, 8)
	i.FillBytes(b)
	return
}

//BytesToBigInt (b []byte) *big.Int - converts a [8]byte array into a big int
func BytesToBigInt(b []byte) *big.Int {
	i := &big.Int{}
	i.SetBytes(b)
	return i
}

//Int16ToBytes (i int16) (b []byte) - Converts an int16 into a [2]byte array
func Int16ToBytes(i int16) (b []byte) {
	b = make([]byte, 2)
	ui := uint16(int64(i) + 32768)
	binary.LittleEndian.PutUint16(b, ui)
	return
}

//BytesToInt16 (b []byte) int16 - Converts a [2]byte array into an int16
func BytesToInt16(b []byte) int16 {
	return int16(int64(binary.LittleEndian.Uint16(b)) - 32768)
}

//Uint32ToBytes (i uint32) (b []byte) - Converts an uint32 into a [4]byte array
func Uint32ToBytes(i uint32) (b []byte) {
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return
}

//BytesToUint32 (b []byte) uint32 - Converts a [4]byte array into an uint32
func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

//Int16sToBytes (is ...int16) (bs []byte) - Converts a slice of ints into a []byte slice
func Int16sToBytes(is ...int16) (bs []byte) {
	bs = make([]byte, 2*len(is))
	for idx, i := range is {
		b := Int16ToBytes(i)
		bs[idx*2] = b[0]
		bs[idx*2+1] = b[1]
	}
	return
}

//BytesToInt16s (bs []byte) (is []int16) - Converts a []byte slice into a slice of ints
func BytesToInt16s(bs []byte) (is []int16) {
	is = make([]int16, len(bs)/2)
	for idx := range bs {
		if idx%2 == 0 {
			is[idx/2] = BytesToInt16(bs[idx : idx+2])
		}
	}
	return
}

//Int16s2DToBytes (is ...[]int16) (bs []byte) - Converts a 2d slice of ints into a []byte slice
func Int16s2DToBytes(is ...[]int16) (bs []byte) {
	bs = make([]byte, 0)
	bs = append(bs, Int64ToBytes(int64(len(is[0])))...)

	for _, il := range is {
		bs = append(bs, Int16sToBytes(il...)...)
	}

	return
}

//BytesToInt16s2D (bs []byte) (is [][]int16) - Converts a []byte slice into a 2d slice of ints
func BytesToInt16s2D(bs []byte) (is [][]int16) {
	is = make([][]int16, 0)
	xl := int(BytesToInt64(bs[:8])) * 2
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

//AppendInt16ToBytes (i int16, bs []byte) []byte - Appends an int16 to a []byte slice
func AppendInt16ToBytes(i int16, bs []byte) []byte {
	ib := Int16ToBytes(i)
	bs = append(bs, ib[0])
	bs = append(bs, ib[1])
	return bs
}

//CompressBytes (bs []byte) ([]byte, error) - Compresses a []byte slice using gzip
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

//DecompressBytes (bs []byte) ([]byte, error) - Decompresses a []byte slice using gzip
func DecompressBytes(bs []byte) ([]byte, error) {
	if len(bs) <= 0 {
		return bs, nil
	}
	var b bytes.Buffer
	_, err1 := b.Write(bs)
	if err1 != nil {
		return nil, err1
	}
	r, err2 := gzip.NewReader(&b)
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
