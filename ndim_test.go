package numgo

import (
	"reflect"
	"testing"
)

func Test_newNDim(t *testing.T) {
	tests := []struct {
		name  string
		shape []int
		want  nDimMetadata
	}{
		{"Empty", []int{0},
			nDimMetadata{shape: []int{0}, strides: []int{0, 0}, err: nil}},
		{"Vector", []int{5},
			nDimMetadata{shape: []int{5}, strides: []int{5, 1}, err: nil}},
		{"Multi-dim", []int{1, 2, 3},
			nDimMetadata{shape: []int{1, 2, 3}, strides: []int{6, 6, 3, 1}, err: nil}},
		{"Negative Axis", []int{1, -2, 3},
			nDimMetadata{shape: nil, strides: []int{0}, err: NegativeAxis}},
	}
	for _, tt := range tests {
		got := newNDim(tt.shape)
		switch {
		case !reflect.DeepEqual(got.shape, tt.want.shape):
			t.Errorf("%q. shape newNDim() = %v, want %v", tt.name,
				got.shape, tt.want.shape)
		case !reflect.DeepEqual(got.strides, tt.want.strides):
			t.Errorf("%q. strides newNDim() = %v, want %v", tt.name,
				got.strides, tt.want.strides)
		case got.err != tt.want.err:
			t.Errorf("%q. err newNDim() = %v, want %v", tt.name,
				got.err, tt.want.err)
		}
	}
}

func Test_nDimMetadata_reshape(t *testing.T) {
	tests := []struct {
		name  string
		shape []int
		ndim  nDimMetadata
		want  nDimMetadata
	}{
		{name: "Flatten",
			shape: []int{105},
			ndim:  nDimMetadata{shape: []int{3, 5, 7}, strides: []int{105, 35, 7, 1}, err: nil},
			want:  nDimMetadata{shape: []int{105}, strides: []int{105, 1}, err: nil},
		},
		{name: "Open",
			shape: []int{3, 5, 7},
			ndim:  nDimMetadata{shape: []int{105}, strides: []int{105, 1}, err: nil},
			want: nDimMetadata{shape: []int{3, 5, 7},
				strides: []int{105, 35, 7, 1}, err: nil},
		},
		{name: "Negative Axis",
			shape: []int{3, -5, 7},
			ndim:  nDimMetadata{shape: []int{3, 5, 7}, strides: []int{105, 35, 7, 1}, err: nil},
			want: nDimMetadata{shape: []int{3, 5, 7},
				strides: []int{105, 35, 7, 1}, err: NegativeAxis},
		},
		{name: "Resize",
			shape: []int{3, 15, 7},
			ndim:  nDimMetadata{shape: []int{3, 5, 7}, strides: []int{105, 35, 7, 1}, err: nil},
			want: nDimMetadata{shape: []int{3, 5, 7},
				strides: []int{105, 35, 7, 1}, err: ReshapeError},
		},
	}
	for _, tt := range tests {
		tt.ndim.reshape(tt.shape)
		switch {
		case !reflect.DeepEqual(tt.ndim.shape, tt.want.shape):
			t.Errorf("%q. shape ndim.reshape() = %v, want %v", tt.name,
				tt.ndim.shape, tt.want.shape)
		case !reflect.DeepEqual(tt.ndim.strides, tt.want.strides):
			t.Errorf("%q. strides ndim.reshape() = %v, want %v", tt.name,
				tt.ndim.strides, tt.want.strides)
		case tt.ndim.err != tt.want.err:
			t.Errorf("%q. err ndim.reshape() = %v, want %v", tt.name,
				tt.ndim.err, tt.want.err)
		}
	}
}
