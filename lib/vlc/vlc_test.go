package vlc

import (
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

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		encodedText string
		want        string
	}{
		{
			name:        "base test",
			encodedText: "20 30 3C 18 77 4A E4 4D 28",
			want:        "My name is Ted",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.encodedText); got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
