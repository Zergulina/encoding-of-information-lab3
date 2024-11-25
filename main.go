package main

import (
	"LAB3/haffman"
	"LAB3/lzw"
	"os"
)

func main() {
	data, _ := os.ReadFile("D:/bmp/rgb.bmp")

	encodedData := haffman.Encode(lzw.Encode(&data, 600))

	os.WriteFile("D:/bmp/rgb.zap", *encodedData, os.ModeAppend)
	os.WriteFile("D:/bmp/rgb.zap", *encodedData, os.ModeAppend)

	zippedData, _ := os.ReadFile("D:/bmp/rgb.zap")

	decoded := lzw.Decode(haffman.Decode(&zippedData))

	os.WriteFile("D:/bmp/rgb1.bmp", *decoded, os.ModeAppend)
}
