package vlc

import (
	"reflect"
	"testing"
)

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

func TestNewHexChunks(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want HexChunks
	}{
		{
			name: "base test",
			str:  "20 30 3C 18",
			want: HexChunks{"20", "30", "3C", "18"},
		},
		{
			name: "empty string",
			str:  "",
			want: HexChunks{""},
		},
		{
			name: "single hex",
			str:  "AB",
			want: HexChunks{"AB"},
		},
		{
			name: "extra spaces",
			str:  "AA  BB   CC",
			want: HexChunks{"AA", "", "BB", "", "", "CC"},
		},
		{
			name: "trailing space",
			str:  "12 34 ",
			want: HexChunks{"12", "34", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHexChunks(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHexChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexChunk_ToBinary(t *testing.T) {
	tests := []struct {
		name string
		hc   HexChunk
		want BinaryChunk
	}{
		{
			name: "base test",
			hc:   HexChunk("2F"),
			want: BinaryChunk("00101111"),
		},
		{
			name: "base test",
			hc:   HexChunk("80"),
			want: BinaryChunk("10000000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hc.ToBinary(); got != tt.want {
				t.Errorf("ToBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexChunks_ToBinary(t *testing.T) {
	tests := []struct {
		name string
		hcs  HexChunks
		want BinaryChunks
	}{
		{
			name: "base test",
			hcs:  HexChunks{"2F", "80"},
			want: BinaryChunks{"00101111", "10000000"},
		},
		{
			name: "base test",
			hcs:  HexChunks{"00", "20", "40", "00"},
			want: BinaryChunks{"00000000", "00100000", "01000000", "00000000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hcs.ToBinary(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryChunks_Join(t *testing.T) {
	tests := []struct {
		name string
		bcs  BinaryChunks
		want string
	}{
		{
			name: "base test",
			bcs:  BinaryChunks{"0101111", "10000000"},
			want: "010111110000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bcs.Join(); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
