package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type BinaryChunks []BinaryChunk

type BinaryChunk string

type HexChunks []HexChunk

type HexChunk string

type encodingTable map[rune]string

const chunksSize = 8

func Encode(str string) string {
	// prepare text: A -> !a, B -> !b
	str = prepareText(str)

	// encode to binary: some text -> 10010101
	bStr := encodeBin(str)

	// split binary by chunks (8): bits to bytes -> '10010101 10010101 10010101'
	chunks := splitByChunks(bStr, chunksSize)

	// bytes to hex -> '20 30 3C'
	hexChunksStr := chunks.ToHex()

	return hexChunksStr.ToString()
}

func (hcs HexChunks) ToString() string {
	const sep = " "

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var b strings.Builder
	b.WriteString(string(hcs[0]))

	for _, chunk := range hcs[1:] {
		b.WriteString(sep)
		b.WriteString(string(chunk))
	}

	return b.String()
}

func (bcs BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		hexChunk := chunk.ToHex()
		res = append(res, hexChunk)
	}

	return res
}

func (bc BinaryChunk) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunksSize)
	if err != nil {
		panic("can't parse binary chunk: " + err.Error())
	}

	res := fmt.Sprintf("%X", num)

	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
}

// splitByChunks splits binary string by chunks with given size,
// i.g.: '100101011001010110010101' -> '10010101 10010101 10010101'
func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)
	chunksCount := strLen / chunkSize

	if strLen%chunkSize != 0 {
		chunksCount++
	}

	res := make(BinaryChunks, 0, chunksCount)

	var b strings.Builder

	for i, ch := range bStr {
		b.WriteRune(ch)

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(b.String()))
			b.Reset()
		}
	}

	if b.Len() != 0 {
		lastChunk := b.String()
		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))

		res = append(res, BinaryChunk(lastChunk))
	}

	return res
}

// encodeBin converts the input string into a continuous binary string (no spaces)
func encodeBin(str string) string {
	var b strings.Builder

	for _, ch := range str {
		b.WriteString(bin(ch))
	}

	return b.String()
}

func bin(ch rune) string {
	table := getEncodingTable()

	res, ok := table[ch]
	if !ok {
		panic("unknown character: " + string(ch))
	}

	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		'!': "001000",
		'a': "011",
		'b': "0000010",
		'c': "000101",
		'd': "00101",
		'e': "101",
		'f': "000100",
		'g': "0000100",
		'h': "0011",
		'i': "01001",
		'j': "000000001",
		'k': "0000000001",
		'l': "001001",
		'm': "000011",
		'n': "10000",
		'o': "10001",
		'p': "0000101",
		'q': "000000000001",
		'r': "01000",
		's': "0101",
		't': "1001",
		'u': "00011",
		'v': "00000001",
		'w': "0000011",
		'x': "00000000001",
		'y': "0000001",
		'z': "000000000000",
	}
}

// prepareText formats the input for encoding:
// replaces each uppercase letter with '!' followed by its lowercase equivalent
// i.g., "My name is Ben" -> "!my name is !ben"
func prepareText(str string) string {
	var b strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			b.WriteRune('!')
			b.WriteRune(unicode.ToLower(ch))
		} else {
			b.WriteRune(ch)
		}
	}

	return b.String()
}
