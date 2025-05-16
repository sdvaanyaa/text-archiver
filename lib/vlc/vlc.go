package vlc

import (
	"strings"
	"unicode"
)

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

func Decode(encodedText string) string {
	hChunks := NewHexChunks(encodedText)

	// hex chunks -> binary chunks
	bChunks := hChunks.ToBinary()

	// binary chunks -> binary string
	bString := bChunks.Join()

	// build decoding tree
	dTree := getEncodingTable().DecodingTree()

	// bString(dTree) -> text
	decoded := dTree.Decode(bString)

	return restoreText(decoded)
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

// restoreText is opposite to prepareText, it prepares decoded text to export:
// it changes: ! + <lower case letter> -> to upper case letter.
//
//	i.g.: !my name is !ted -> My name is Ted.
func restoreText(str string) string {
	var b strings.Builder

	var isCapital bool

	for _, ch := range str {
		if isCapital {
			b.WriteRune(unicode.ToUpper(ch))
			isCapital = false

			continue
		}

		if ch == '!' {
			isCapital = true

			continue
		} else {
			b.WriteRune(ch)
		}
	}

	return b.String()
}
