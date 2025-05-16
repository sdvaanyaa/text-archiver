package vlc

import (
	"reflect"
	"testing"
)

func Test_prepareText(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "base test",
			str:  "My name is Bob",
			want: "!my name is !bob",
		},
		{
			name: "emtpy string",
			str:  "",
			want: "",
		},
		{
			name: "no uppercase",
			str:  "hello, i'm fine",
			want: "hello, i'm fine",
		},
		{
			name: "all uppercase",
			str:  "GOLANG AWESOME",
			want: "!g!o!l!a!n!g !a!w!e!s!o!m!e",
		},
		{
			name: "mixed case",
			str:  "GoLaNg",
			want: "!go!la!ng",
		},
		{
			name: "multiple spaces",
			str:  "Golang           Awesome",
			want: "!golang           !awesome",
		},
		{
			name: "symbol and digits",
			str:  "How much is 2 + 2 ?",
			want: "!how much is 2 + 2 ?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareText(tt.str); got != tt.want {
				t.Errorf("prepareText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeBin(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "base test",
			str:  "!bob",
			want: "0010000000010100010000010",
		},
		{
			name: "empty string",
			str:  "",
			want: "",
		},
		{
			name: "space only",
			str:  " ",
			want: "11",
		},
		{
			name: "single letter",
			str:  "e",
			want: "101",
		},
		{
			name: "double exclamation",
			str:  "!!",
			want: "001000001000",
		},
		{
			name: "unknown symbol",
			str:  "A",
			want: "panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.want != "panic" {
						t.Errorf("unexpected panic: %v", r)
					}
				} else {
					if tt.want == "panic" {
						t.Errorf("expected panic but none occurred")
					}
				}
			}()

			got := encodeBin(tt.str)

			if tt.want != "panic" && got != tt.want {
				t.Errorf("encodeBin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitByChunks(t *testing.T) {
	type args struct {
		bStr      string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want BinaryChunks
	}{
		{
			name: "not multiple of chunkSize",
			args: args{
				bStr:      "0010000000010100010000010",
				chunkSize: 8,
			},
			want: BinaryChunks{"00100000", "00010100", "01000001", "00000000"},
		},
		{
			name: "exact chunkSize",
			args: args{
				bStr:      "11110000",
				chunkSize: 8,
			},
			want: BinaryChunks{"11110000"},
		},
		{
			name: "empty string",
			args: args{
				bStr:      "",
				chunkSize: 8,
			},
			want: BinaryChunks{},
		},
		{
			name: "shorter than chunkSize",
			args: args{
				bStr:      "101",
				chunkSize: 8,
			},
			want: BinaryChunks{"10100000"},
		},
		{
			name: "chunkSize 1",
			args: args{
				bStr:      "101",
				chunkSize: 1,
			},
			want: BinaryChunks{"1", "0", "1"},
		},
		{
			name: "chunkSize larger than input",
			args: args{
				bStr:      "10101",
				chunkSize: 10,
			},
			want: BinaryChunks{"1010100000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByChunks(tt.args.bStr, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryChunks_ToHex(t *testing.T) {
	tests := []struct {
		name string
		bcs  BinaryChunks
		want HexChunks
	}{
		{
			name: "base test",
			bcs:  BinaryChunks{"0101111", "10000000"},
			want: HexChunks{"2F", "80"},
		},
		{
			name: "empty slice",
			bcs:  BinaryChunks{},
			want: HexChunks{},
		},
		{
			name: "all zeros",
			bcs:  BinaryChunks{"00000000", "00000000"},
			want: HexChunks{"00", "00"},
		},
		{
			name: "all ones",
			bcs:  BinaryChunks{"11111111"},
			want: HexChunks{"FF"},
		},
		{
			name: "non-full 8-bit chunk padded",
			bcs:  BinaryChunks{"101"},
			want: HexChunks{"05"},
		},
		{
			name: "mixed bits",
			bcs:  BinaryChunks{"00001111", "10101010"},
			want: HexChunks{"0F", "AA"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcs.ToHex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "base test",
			str:  "My name is Ted",
			want: "20 30 3C 18 77 4A E4 4D 28",
		},
		{
			name: "empty string",
			str:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.str); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
