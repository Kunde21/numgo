package numgo

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Sum(0).data)
	fmt.Println(Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6).Sum(3, 4))
}

func TestCount(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Count(0).data)
	fmt.Println(Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6).Count(3, 4))
}

func TestMean(t *testing.T) {
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	fmt.Println(a.Count())

	fmt.Println(a.C().Sum(3, 4).Div(a.Count(3, 4)))
	fmt.Println(a.Mean(3, 4))
}
