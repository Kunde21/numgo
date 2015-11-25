package numgo

import (
	"fmt"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	a, b := Arange(20), Arange(20)
	a.Add(b)
	fmt.Println(a.Add(b.Reshape(2, 10)))
	fmt.Println(Arange(20).Reshape(2, 10).Add(b.Reshape(2, 10)))
	fmt.Println(Arange(20).Reshape(2, 2, 5).Add(Arange(5)))
}

func TestSubtr(t *testing.T) {
	fmt.Println(Arange(20).Reshape(2, 10).Subtr(Arange(10)))
	fmt.Println(Arange(20).Reshape(2, 10).Subtr(Arange(10)))
	fmt.Println(Arange(20).Reshape(2, 2, 5).Subtr(Arange(5)))
	fmt.Println(Arange(20).Reshape(2, 1, 1, 2, 5).Subtr(Arange(5)))
}

func TestMult(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5 {
		t.Log("Expected:", (100-1)/.5, "Got:", len(a.data))
		t.FailNow()
	}

	a = Full(math.NaN(), 4, 4)
	for _, v := range a.data {
		if !math.IsNaN(v) {
			t.Log("Expected NaN, got ", v)
			t.FailNow()
		}
	}
}
