package lzw

type DictUnit {
	key string,
	val uint32 
}

func Encode(data []byte, dictionarySize uint32 ) []byte {
	bufferedCode := ""
	dictionary := make([]DictUnit, dictionarySize)
	encoded := make([]byte, 0, len(data))
}