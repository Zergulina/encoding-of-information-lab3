package main

import (
	"LAB3/lzw"
	"os"
)

func main() {
	data, _ := os.ReadFile("D:/bmp/rgb.bmp")

	encodedData := lzw.Encode(&data, 10000000)

	os.WriteFile("D:/bmp/rgb.zap", *encodedData, os.ModeAppend)

	zippedData, _ := os.ReadFile("D:/bmp/rgb.zap")

	decoded := lzw.Decode(&zippedData)

	// fmt.Println(data[213582:213612])
	// fmt.Println((*decoded)[213582:213612])

	os.WriteFile("D:/bmp/rgb1.bmp", *decoded, os.ModeAppend)

	// fmt.Println(*decoded)

	// fmt.Println(len(encodedData), len(zippedData))

	// for i := 0; i < len(encodedData); i++ {
	// 	if encodedData[i] != zippedData[i] {
	// 		fmt.Println(i)
	// 	}
	// }
}
