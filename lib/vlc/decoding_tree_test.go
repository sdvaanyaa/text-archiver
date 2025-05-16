package vlc

import (
	"reflect"
	"testing"
)

func Test_encodingTable_DecodingTree(t *testing.T) {
	tests := []struct {
		name string
		et   encodingTable
		want DecodingTree
	}{
		{
			name: "base tree test",
			et: encodingTable{
				'a': "11",
				'b': "1001",
				'z': "0101",
			},
			want: DecodingTree{
				Zero: &DecodingTree{
					One: &DecodingTree{
						Zero: &DecodingTree{
							One: &DecodingTree{
								Value: "z",
							},
						},
					},
				},
				One: &DecodingTree{
					Zero: &DecodingTree{
						Zero: &DecodingTree{
							One: &DecodingTree{
								Value: "b",
							},
						},
					},
					One: &DecodingTree{
						Value: "a",
					},
				},
			},
		},
		{
			name: "empty tree test",
			et:   encodingTable{},
			want: DecodingTree{},
		},
		{
			name: "single element test",
			et: encodingTable{
				'x': "0",
			},
			want: DecodingTree{
				Zero: &DecodingTree{
					Value: "x",
				},
			},
		},
		{
			name: "deep tree test",
			et: encodingTable{
				'x': "00000",
			},
			want: DecodingTree{
				Zero: &DecodingTree{
					Zero: &DecodingTree{
						Zero: &DecodingTree{
							Zero: &DecodingTree{
								Zero: &DecodingTree{
									Value: "x",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.DecodingTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodingTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
