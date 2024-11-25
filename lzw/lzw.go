package lzw

import "fmt"

func Encode(data *[]byte, dictionarySize uint32) *[]byte {
	dictionary := make(map[string]uint32)
	for i := 0; i < 256; i++ {
		dictionary[string(byte(i))] = uint32(i)
	}

	currentSize := uint32(256)
	encode := make([]byte, 0, len(*data))
	w := ""

	for i := 0; i < 4; i++ {
		encode = append(encode, byte((dictionarySize)>>(8*(3-i))))
	}

	for _, k := range *data {
		wk := w + string(k)
		if _, ok := dictionary[wk]; ok {
			w = wk
		} else {
			for j := 0; j < 4; j++ {
				encode = append(encode, byte(dictionary[w]>>(8*(3-j))))
			}

			dictionary[wk] = uint32(currentSize)
			currentSize++
			w = string(k)
			if currentSize == dictionarySize+256 {
				currentSize = 256
				dictionary = make(map[string]uint32)
				for i := 0; i < 256; i++ {
					dictionary[string(byte(i))] = uint32(i)
				}
			}
		}
	}
	if len(w) > 0 {
		for j := 0; j < 4; j++ {
			encode = append(encode, byte(dictionary[w]>>(8*(3-j))))
		}
	}
	return &encode
}

func bytesToUint32(arr []byte) uint32 {
	var result uint32 = 0
	for i, b := range arr {
		result |= uint32(b) << ((3 - i) * 8)
	}
	return result
}

func Decode(data *[]byte) *[]byte {
	dictionary := make(map[uint32][]byte)
	for i := 0; i < 256; i++ {
		dictionary[uint32(i)] = []byte{byte(i)}
	}

	currentSize := uint32(256)

	decode := make([]byte, 0, len(*data)*4)

	dictionarySize := bytesToUint32((*data)[:4])

	fmt.Println(dictionarySize)

	decode = append(decode, (*data)[7])

	w := []byte{(*data)[7]}

	for i := 8; i < len(*data); i += 4 {
		var entry []byte
		k := bytesToUint32((*data)[i : i+4])
		if _, ok := dictionary[k]; ok {
			dst := make([]byte, len(dictionary[k]))
			copy(dst, dictionary[k])
			entry = dst
		} else {
			if k == uint32(currentSize) {
				dst := make([]byte, len(w))
				copy(dst, w)
				entry = append(dst, dst[0])
			} else {
				continue
			}
		}

		decode = append(decode, entry...)
		dst := make([]byte, len(w))
		copy(dst, w)
		dictionary[uint32(currentSize)] = append(dst, entry[0])
		currentSize++
		w = entry

		if currentSize == dictionarySize+256 {
			currentSize = 256
			dictionary = make(map[uint32][]byte)
			for i := 0; i < 256; i++ {
				dictionary[uint32(i)] = []byte{byte(i)}
			}
		}
	}

	return &decode
}
