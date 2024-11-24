package main

import (
	"LAB3/haffman"
	"os"
)

func main() {
	data, _ := os.ReadFile("D:/bmp/rgb.bmp")

	encodedData := haffman.HaffmanEncode(data)

	os.WriteFile("D:/bmp/rgb.zap", encodedData, os.ModeAppend)

	zippedData, _ := os.ReadFile("D:/bmp/rgb.zap")

	decoded := haffman.HaffmanDecode(zippedData)

	os.WriteFile("D:/bmp/rgb1.bmp", decoded, os.ModeAppend)

	// fmt.Println(len(encodedData), len(zippedData))

	// for i := 0; i < len(encodedData); i++ {
	// 	if encodedData[i] != zippedData[i] {
	// 		fmt.Println(i)
	// 	}
	// }
}
