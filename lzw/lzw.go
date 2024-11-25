package lzw

import (
	"fmt"
	"os"
)

func Encode(data *[]byte, dictionarySize uint32) *[]byte {
	bufferedCode := make([]byte, 0, dictionarySize/2)
	dictionary := make(map[string]uint32)
	encode := make([]byte, 0, len(*data))
	var counter uint64 = 0

	for i := 0; i < 4; i++ {
		encode = append(encode, byte((dictionarySize)>>(8*(3-i))))
	}

	for i := 0; i < len((*data))-1; i++ {
		bufferedCode = append(bufferedCode, (*data)[i])
		if _, ok := dictionary[string(append(bufferedCode, (*data)[i+1]))]; ok {
			continue
		}
		if len(bufferedCode) == 1 {
			for j := 0; j < 3; j++ {
				encode = append(encode, 0)
			}
			encode = append(encode, bufferedCode[0])
		} else {
			for j := 0; j < 4; j++ {
				encode = append(encode, byte((dictionary[string(bufferedCode)]+256)>>(8*(3-j))))
			}
		}

		for key, val := range dictionary {
			if val == uint32(counter) {
				delete(dictionary, key)
				break
			}
		}
		dictionary[string(append(bufferedCode, (*data)[i+1]))] = uint32(counter)

		bufferedCode = make([]byte, 0, dictionarySize/2)
		counter++
		if counter == uint64(dictionarySize) {
			counter = 0
		}
		// fmt.Println(counter)
		// if encode[len(encode)-4] == 0 && encode[len(encode)-3] == 0 && encode[len(encode)-2] == 0 && encode[len(encode)-1] == 0 {
		// 	fmt.Println(i, byte((dictionary[string(bufferedCode)]+256)>>24), byte((dictionary[string(bufferedCode)]+256)>>16), byte((dictionary[string(bufferedCode)]+256)>>8), byte(dictionary[string(bufferedCode)]+256), encode[len(encode)-4:len(encode)])
		// }

	}

	bufferedCode = append(bufferedCode, (*data)[len((*data))-1])
	if len(bufferedCode) == 1 {
		for i := 0; i < 3; i++ {
			encode = append(encode, 0)
		}
		encode = append(encode, bufferedCode[0])
	} else {
		for i := 0; i < 4; i++ {
			encode = append(encode, byte((dictionary[string(bufferedCode)]+256)>>(8*(3-i))))
		}
	}

	// fmt.Println(encode)

	file, err := os.Create("input.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
	}
	defer file.Close()

	// Итерация по мапе и запись в файл
	for key, value := range dictionary {
		// Преобразование ключа и значения в строку
		keyStr := fmt.Sprintf("%v", []byte(key))
		valueStr := fmt.Sprintf("%v", value)

		// Запись ключа и значения в файл
		line := keyStr + ": " + valueStr + "\n"
		_, err = file.WriteString(line)
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
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
	decode := make([]byte, 0, 4*len(*data))
	var dictionarySize uint32 = 0

	dictionary := make(map[uint64][]byte)

	dataBuffer := make([]byte, 0, len(*data))

	for i := 0; i < 4; i++ {
		dictionarySize |= (uint32((*data)[i]) << ((3 - i) * 8))
	}

	var counter uint64 = 0

	for i := 4; i < len(*data); i += 4 {
		_, ok := dictionary[counter]
		if ok {
			delete(dictionary, counter)
		}

		code := bytesToUint32((*data)[i : i+4])
		if code < 256 {
			decode = append(decode, byte(code))
			if len(dataBuffer) > 0 {
				dst := make([]byte, len(dataBuffer))
				copy(dst, dataBuffer)
				dictionary[counter] = append(dst, byte(code))
				counter++

			}
			// fmt.Println(decode[:25])
			dataBuffer = []byte{byte(code)}
		} else {
			if val, ok := dictionary[uint64(code)-256]; ok {
				decode = append(decode, val...)
				dst := make([]byte, len(dataBuffer))
				copy(dst, dataBuffer)
				dictionary[counter] = append(dst, val[0])
				dataBuffer = dictionary[uint64(code)-256]

				counter++
			} else {
				dst := make([]byte, len(dataBuffer))
				copy(dst, dataBuffer)
				dataBuffer = append(dst, dst[len(dataBuffer)-1])
				dst = make([]byte, len(dataBuffer))
				copy(dst, dataBuffer)
				decode = append(decode, dst...)
				dst = make([]byte, len(dataBuffer))
				copy(dst, dataBuffer)
				dictionary[counter] = dst
				counter++

			}
		}
		if counter == uint64(dictionarySize) {
			counter = 0
		}
	}

	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
	}
	defer file.Close()

	// Итерация по мапе и запись в файл
	for key, value := range dictionary {
		// Преобразование ключа и значения в строку
		keyStr := fmt.Sprintf("%d", key)
		valueStr := fmt.Sprintf("%v", value)

		// Запись ключа и значения в файл
		line := keyStr + ": " + valueStr + "\n"
		_, err = file.WriteString(line)
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
		}
	}

	return &decode
}

// func Decode(data *[]byte) *[]byte {
// 	decode := make([]byte, 0, 4*len(*data))
// 	var dictionarySize uint32 = 0

// 	dictionary := make(map[uint64][]byte)

// 	dataBuffer := make([]byte, 0, len(*data))

// 	dictionarySize = bytesToUint32((*data)[4:8])

// 	dataBuffer = append(dataBuffer, byte(bytesToUint32((*data)[4:8])))

// 	decode = append(decode, (*data)[7])

// 	var counter uint64 = 0

// 	for i := 8; i < len(*data); i += 4 {
// 		if _, ok := dictionary[counter]; ok {
// 			delete(dictionary, counter)
// 		}

// 		code := bytesToUint32((*data)[i : i+4])
// 		if code < 256 {
// 			decode = append(decode, byte(code))
// 			dst := make([]byte, len(dataBuffer))
// 			copy(dst, dataBuffer)
// 			dictionary[counter] = append(dst, byte(code))
// 			counter++
// 			dataBuffer = []byte{byte(code)}
// 		} else {
// 			if val, ok := dictionary[uint64(code)-256]; ok {
// 				dst1 := make([]byte, len(val))
// 				copy(dst1, val)
// 				decode = append(decode, dst1...)
// 				dst2 := make([]byte, len(dataBuffer))
// 				copy(dst2, dataBuffer)
// 				dictionary[counter] = append(dst2, dst1[0])
// 			}
// 		}
// 	}
// }
