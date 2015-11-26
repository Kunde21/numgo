package numgo

import (
	"fmt"
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Sum(0).data)
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	fmt.Println(a.shape)
	fmt.Println(a.strides)
	fmt.Println(a.Sum(0, 2, 4))
}

func TestCount(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Count(0).data)
	fmt.Println(Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6).Count(3, 4))
}

func TestMean(t *testing.T) {
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	fmt.Println(a.Count())

	fmt.Println(a.C().Sum(3, 4).Div(a.Count(3, 4)))

	fmt.Println(a.shape)
	fmt.Println(a.strides)
	fmt.Println(a.Count(0, 2, 4))
	fmt.Println(a.Mean(0, 2, 4))
}

func TestClean(t *testing.T) {
	a := []int{0, 1, 2, 3, 4, 5}
	b, c := a[:3], a[3:]
	fmt.Println(cleanAxis(0, 1, 2, 3, 4))
	fmt.Println(cleanAxis(0, 1, 2, 2, 4))
	fmt.Println(cleanAxis(0, 1, 2, 2, 4, 4))
	fmt.Println(12.0 + math.NaN())
	fmt.Println(b, c)
	b, c = c, b
	fmt.Println(a)
}
