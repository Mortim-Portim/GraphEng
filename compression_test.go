package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/mortim-portim/GraphEng/GE"
	cmp "github.com/mortim-portim/GraphEng/compression"
)

func TestBoolsToBytes(t *testing.T) {
	fmt.Println("\nTestBoolsToBytes")
	bools := []bool{true, false, true, true, false, false, false}
	bs := cmp.BoolsToBytes(bools...)
	fmt.Printf("%v -> %v\n", bools, bs)
	newBools := cmp.BytesToBools(bs)
	fmt.Printf("%v -> %v\n", bs, newBools)
}

func TestUIn16ToBytes(t *testing.T) {
	fmt.Println("\nTestUIn16ToBytes")
	i := uint16(60000)
	bs := cmp.UInt16ToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToUInt16(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestUIn64ToBytes(t *testing.T) {
	fmt.Println("\nTestUIn64ToBytes")
	i := uint64(60000)
	bs := cmp.UInt64ToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToUInt64(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestBoolToByte(t *testing.T) {
	fmt.Println("\nTestBoolToByte")
	b := false
	bs := cmp.BoolToByte(b)
	fmt.Printf("%v -> %v\n", b, bs)
	newB := cmp.ByteToBool(bs)
	fmt.Printf("%v -> %v\n", bs, newB)
}

func TestFloat32ToByte(t *testing.T) {
	fmt.Println("\nTestFloat32ToByte")
	f := float32(3545.235342)
	bs := cmp.Float32ToBytes(f)
	fmt.Printf("%v -> %v\n", f, bs)
	newF := cmp.BytesToFloat32(bs)
	fmt.Printf("%v -> %v\n", bs, newF)
}

func TestFloat64ToByte(t *testing.T) {
	fmt.Println("\nTestFloat64ToByte")
	f := float64(3545.235342)
	bs := cmp.Float64ToBytes(f)
	fmt.Printf("%v -> %v\n", f, bs)
	newF := cmp.BytesToFloat64(bs)
	fmt.Printf("%v -> %v\n", bs, newF)
}

func TestInt64ToByte(t *testing.T) {
	fmt.Println("\nTestIn64ToBytes")
	i := int64(600007365)
	bs := cmp.Int64ToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToInt64(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestInt64sToByte(t *testing.T) {
	fmt.Println("\nTestIn64sToBytes")
	is := []int64{8912374, 928745, 91874365, 9823476}
	bs := cmp.Int64sToBytes(is...)
	fmt.Printf("%v -> %v\n", is, bs)
	newIs := cmp.BytesToInt64s(bs)
	fmt.Printf("%v -> %v\n", bs, newIs)
}

func TestBigIntToByte(t *testing.T) {
	fmt.Println("\nTestBigIntToByte")
	i := big.NewInt(982735)
	bs := cmp.BigIntToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToBigInt(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestInt16ToByte(t *testing.T) {
	fmt.Println("\nTestInt16ToByte")
	i := int16(23876)
	bs := cmp.Int16ToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToInt16(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestUint32ToByte(t *testing.T) {
	fmt.Println("\nTestUint32ToByte")
	i := uint32(23876)
	bs := cmp.Uint32ToBytes(i)
	fmt.Printf("%v -> %v\n", i, bs)
	newI := cmp.BytesToUint32(bs)
	fmt.Printf("%v -> %v\n", bs, newI)
}

func TestInt16sToByte(t *testing.T) {
	fmt.Println("\nTestIn16sToBytes")
	is := []int16{8912, 9287, 9187, 9823}
	bs := cmp.Int16sToBytes(is...)
	fmt.Printf("%v -> %v\n", is, bs)
	newIs := cmp.BytesToInt16s(bs)
	fmt.Printf("%v -> %v\n", bs, newIs)
}

func TestGeneral(t *testing.T) {
	fmt.Println("\nTestGeneral")
	BYTE := byte(53)
	UINT8 := uint8(26)
	INT8 := int8(29)
	INT16 := int16(876)
	UINT16 := uint16(9369)
	UINT32 := uint32(987263)
	BOOL := true
	BOOLS := []bool{true, false, true, true}

	BYTE_BS := cmp.GetBytesOfVar(BYTE)
	UINT8_BS := cmp.GetBytesOfVar(UINT8)
	INT8_BS := cmp.GetBytesOfVar(INT8)
	INT16_BS := cmp.GetBytesOfVar(INT16)
	UINT16_BS := cmp.GetBytesOfVar(UINT16)
	UINT32_BS := cmp.GetBytesOfVar(UINT32)
	BOOL_BS := cmp.GetBytesOfVar(BOOL)
	BOOLS_BS := cmp.GetBytesOfVar(BOOLS)

	BYTE_N := byte(0)
	UINT8_N := uint8(0)
	INT8_N := int8(0)
	INT16_N := int16(0)
	UINT16_N := uint16(0)
	UINT32_N := uint32(0)
	BOOL_N := false
	BOOLS_N := []bool{}
	GE.ShitImDying(cmp.GetVarFromBytes(BYTE_BS, &BYTE_N))
	GE.ShitImDying(cmp.GetVarFromBytes(UINT8_BS, &UINT8_N))
	GE.ShitImDying(cmp.GetVarFromBytes(INT8_BS, &INT8_N))
	GE.ShitImDying(cmp.GetVarFromBytes(INT16_BS, &INT16_N))
	GE.ShitImDying(cmp.GetVarFromBytes(UINT16_BS, &UINT16_N))
	GE.ShitImDying(cmp.GetVarFromBytes(UINT32_BS, &UINT32_N))
	GE.ShitImDying(cmp.GetVarFromBytes(BOOL_BS, &BOOL_N))
	GE.ShitImDying(cmp.GetVarFromBytes(BOOLS_BS, &BOOLS_N))

	fmt.Printf("byte: %v -> %v -> %v\n", BYTE, BYTE_BS, BYTE_N)
	fmt.Printf("uint8: %v -> %v -> %v\n", UINT8, UINT8_BS, UINT8_N)
	fmt.Printf("int8: %v -> %v -> %v\n", INT8, INT8_BS, INT8_N)
	fmt.Printf("int16: %v -> %v -> %v\n", INT16, INT16_BS, INT16_N)
	fmt.Printf("uint16: %v -> %v -> %v\n", UINT16, UINT16_BS, UINT16_N)
	fmt.Printf("uint32: %v -> %v -> %v\n", UINT32, UINT32_BS, UINT32_N)
	fmt.Printf("bool: %v -> %v -> %v\n", BOOL, BOOL_BS, BOOL_N)
	fmt.Printf("bools: %v -> %v -> %v\n", BOOLS, BOOLS_BS, BOOLS_N)
}
