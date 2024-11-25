package haffman

import "sort"

type Node struct {
	Value byte
	Freq  uint32
	Left  *Node
	Right *Node
}

func NewNode(value byte, freq uint32, left *Node, right *Node) *Node {
	return &Node{value, freq, left, right}
}

type ByFreq []*Node

func (a ByFreq) Len() int { return len(a) }

func (a ByFreq) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByFreq) Less(i, j int) bool {
	return (a[i].Freq > a[j].Freq) || ((a[i].Freq == a[j].Freq) && (a[i].Value < a[j].Value))
}

func Insert(a *ByFreq, x *Node) {
	*a = append(*a, x)
}

func CountBytes(data []byte) *map[byte]uint32 {
	frequency := make(map[byte]uint32)
	for _, c := range data {
		frequency[c]++
	}

	return &frequency
}

func BuildHuffmanTree(data []byte, frequency *map[byte]uint32) *Node {
	pq := &ByFreq{}

	for char, freq := range *frequency {
		node := NewNode(char, freq, nil, nil)
		Insert(pq, node)
	}

	var x, y *Node

	for pq.Len() > 1 {
		sort.Stable(pq)
		x = (*pq)[pq.Len()-1]
		y = (*pq)[pq.Len()-2]
		(*pq) = (*pq)[:pq.Len()-2]
		node := NewNode(0, x.Freq+y.Freq, x, y)
		Insert(pq, node)
	}
	return (*pq)[0]
}

func BuildHuffmanCodesInternal(root *Node, s string, huffmanCodes *map[byte]string, ch chan struct{}) {
	if root.Left == nil && root.Right == nil {
		(*huffmanCodes)[root.Value] = s
		close(ch)
		return
	}

	leftCh := make(chan struct{})
	rightCh := make(chan struct{})

	BuildHuffmanCodesInternal(root.Left, s+"0", huffmanCodes, leftCh)
	BuildHuffmanCodesInternal(root.Right, s+"1", huffmanCodes, rightCh)
	<-leftCh
	<-rightCh
	close(ch)
}

func BuildHuffmanCodes(root *Node) *map[byte]string {
	huffmanCodes := make(map[byte]string)
	ch := make(chan struct{})

	BuildHuffmanCodesInternal(root, "", &huffmanCodes, ch)

	<-ch
	return &huffmanCodes
}

func Encode(data []byte) []byte {
	countedBytes := CountBytes(data)
	haffmanTree := BuildHuffmanTree(data, countedBytes)
	codes := BuildHuffmanCodes(haffmanTree)

	encoded := make([]byte, 0, len(data))

	for countedByte, frequency := range *countedBytes {
		encoded = append(encoded, countedByte)
		encoded = append(encoded, byte(frequency>>24))
		encoded = append(encoded, byte(frequency>>16))
		encoded = append(encoded, byte(frequency>>8))
		encoded = append(encoded, byte(frequency))
	}

	for i := 0; i < 5; i++ { //5 нулевых байт разделяют заголовок и закодированные данные
		encoded = append(encoded, 0)
	}

	var bufferedByte byte = 0
	var counter byte = 0
	for _, dataByte := range data {
		currentCode := (*codes)[dataByte]
		for _, bit := range currentCode {
			if bit == '1' {
				bufferedByte |= (1 << (7 - counter))
			}
			counter++
			if counter == 8 {
				encoded = append(encoded, bufferedByte)
				bufferedByte = 0
				counter = 0
			}
		}
	}
	encoded = append(encoded, bufferedByte)
	encoded = append(encoded, counter)
	return encoded
}

func allZeroes(arr []byte) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

func bytesToUint32(arr []byte) uint32 {
	var result uint32 = 0
	for i, b := range arr {
		result |= uint32(b) << ((3 - i) * 8)
	}
	return result
}

func BuildHuffmanDecodesInternal(root *Node, s string, huffmanCodes *map[string]byte, ch chan struct{}) {
	if root.Left == nil && root.Right == nil {
		(*huffmanCodes)[s] = root.Value
		close(ch)
		return
	}

	leftCh := make(chan struct{})
	rightCh := make(chan struct{})

	BuildHuffmanDecodesInternal(root.Left, s+"0", huffmanCodes, leftCh)
	BuildHuffmanDecodesInternal(root.Right, s+"1", huffmanCodes, rightCh)
	<-leftCh
	<-rightCh
	close(ch)
}

func BuildHuffmanDecodes(root *Node) *map[string]byte {
	huffmanCodes := make(map[string]byte)
	ch := make(chan struct{})

	BuildHuffmanDecodesInternal(root, "", &huffmanCodes, ch)

	<-ch
	return &huffmanCodes
}

func Decode(data []byte) []byte {
	var terminateBuffer = [5]byte{}
	var dataStartBios = 0

	countedBytes := make(map[byte]uint32)

	for ; true; dataStartBios += 5 {
		terminateBuffer[0] = data[dataStartBios]
		terminateBuffer[1] = data[dataStartBios+1]
		terminateBuffer[2] = data[dataStartBios+2]
		terminateBuffer[3] = data[dataStartBios+3]
		terminateBuffer[4] = data[dataStartBios+4]

		// fmt.Println(terminateBuffer)

		if allZeroes(terminateBuffer[:]) {
			dataStartBios += 5
			break
		}

		countedBytes[terminateBuffer[0]] = bytesToUint32(terminateBuffer[1:])
	}

	haffmanTree := BuildHuffmanTree(data, &countedBytes)
	codes := BuildHuffmanDecodes(haffmanTree)

	maxCodeLength := 0
	for code := range *codes {
		if len(code) > maxCodeLength {
			maxCodeLength = len(code)
		}
	}

	decoded := make([]byte, 0, len(data))

	bufferedCode := ""
	for _, dataByte := range data[dataStartBios : len(data)-2] {
		for counter := 0; counter < 8; counter++ {
			if ((int(dataByte) >> (7 - counter)) & 1) == 1 {
				bufferedCode += "1"
			} else {
				bufferedCode += "0"
			}
			if val, ok := (*codes)[bufferedCode]; ok {
				decoded = append(decoded, val)
				bufferedCode = ""
			} else if len(bufferedCode) > maxCodeLength {
				bufferedCode = ""
			}
		}
	}

	lastByteValidBitsAmount := data[len(data)-1]
	if lastByteValidBitsAmount < 8 {
		for counter := 0; counter < int(lastByteValidBitsAmount); counter++ {
			if int(data[len(data)-2])>>(7-counter)&1 == 1 {
				bufferedCode += "1"
			} else {
				bufferedCode += "0"
			}
			if val, ok := (*codes)[bufferedCode]; ok {
				decoded = append(decoded, val)
				bufferedCode = ""
			}
		}
	}

	return decoded
}
