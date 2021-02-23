package compression

//Merge (changing [][]byte, static ...[]byte) (merged []byte) - adds byte slices of static and changing length
func Merge(changing [][]byte, static ...[]byte) (merged []byte) {
	merged = make([]byte, 0)
	if len(changing) <= 0 && len(static) <= 0 {
		return
	}
	for _, bs := range static {
		merged = append(merged, bs...)
	}
	for _, bs := range changing {
		merged = append(merged, Int64ToBytes(int64(len(bs)))...)
		merged = append(merged, bs...)
	}
	return
}

//Demerge (comp []byte, length []int) (BS [][]byte) - Decodes merged bytes, returning static ones first
func Demerge(comp []byte, length []int) (BS [][]byte) {
	BS = make([][]byte, 0)
	if len(comp) <= 0 && len(length) <= 0 {
		return
	}
	currIdx := 0
	for _, l := range length {
		if currIdx+l > len(comp) {
			break
		}
		BS = append(BS, comp[currIdx:currIdx+l])
		currIdx += l
	}

	for true {
		length := int(BytesToInt64(comp[currIdx : currIdx+8]))
		currIdx += 8
		BS = append(BS, comp[currIdx:currIdx+length])
		currIdx += length
		if currIdx >= len(comp) {
			break
		}
	}
	return
}

//CompressAll (changing [][]byte, static ...[]byte) (comp []byte) - Compresses byte slices of static and changing length
func CompressAll(changing [][]byte, static ...[]byte) (comp []byte) {
	comp = Merge(changing, static...)
	comp, _ = CompressBytes(comp)
	return
}

//DecompressAll (comp []byte, length []int - Decompresses bytes, returning static ones first
func DecompressAll(comp []byte, length []int) (BS [][]byte) {
	comp, _ = DecompressBytes(comp)
	BS = Demerge(comp, length)
	return
}
