package compression

import (
	"bytes"
	"encoding/binary"
)

//GetBytesOfVar (i interface{}) []byte - Converts any variable into bytes
func GetBytesOfVar(i interface{}) []byte {
	switch v := i.(type) {
	case int:
		i = int64(v)
	case *int:
		i = int64(*v)
	case uint:
		i = uint64(v)
	case *uint:
		i = uint64(*v)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//GetVarFromBytes (bs []byte, i interface{}) (rerr error) - Sets i from the given bytes, i is pointer to the variable (e.g.: *int)
func GetVarFromBytes(bs []byte, i interface{}) (rerr error) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(bs)
	if err != nil {
		panic(err)
	}

	err2 := binaryReadTo(buf, i)
	if err2 != nil {
		switch i := i.(type) {
		case *int:
			var v int64
			rerr = binaryReadTo(buf, &v)
			*i = int(v)
			break
		case *uint:
			var v uint64
			rerr = binaryReadTo(buf, &v)
			*i = uint(v)
			break
		}
	}
	return
}
func binaryReadTo(buf *bytes.Buffer, i interface{}) error {
	return binary.Read(buf, binary.LittleEndian, i)
}
