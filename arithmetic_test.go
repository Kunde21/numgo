package numgo

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	a, b := Arange(20), Arange(20)
	a.Add(b)
	if c := a.Add(b.Reshape(2, 10)); c.err != ShapeError {
		t.Log("Shape tests failed.  Expected nil, returned:", c)
		t.Fail()
	}

	if c := Arange(20).Reshape(2, 10).Add(b.Reshape(2, 10)); c.Equals(Arange(0, 40, 2).Reshape(2, 10)).All().data[0] == false {
		t.Log("Addition on multiple axis failed.  Expected 0-40 with increment by 2.")
		t.Log("Recieved:", c)
		t.Fail()
	}

	a = Arange(20).Reshape(2, 2, 5).Add(Arange(5))
	for i, v := range a.data {
		if int(v) != i%5*2+i/5*5 {
			t.Log("Incorrect result at:", i, "Expected", i%5*2+i/5*5, "Got", v)
			t.Fail()
		}
	}
}

func TestSubtr(t *testing.T) {
	a := Arange(20).Reshape(2, 10).Subtr(Arange(10))
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i)
		for _, v := range tmp {
			if v != float64(i*10) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("At", i, "Expected: ", float64(i*10), "Got:", v)
				t.Fail()
			}
		}
	}

	a = Arange(20).Reshape(2, 2, 5).Subtr(Arange(5))
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i)
		for _, v := range tmp {
			if v != float64(i/5*5) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("Expected: ", i*5, "Got:", v)
				t.Fail()
			}
		}
	}

	a = Arange(20).Reshape(2, 1, 1, 2, 5).Subtr(Arange(5))
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i, 0, 0)
		for _, v := range tmp {
			if v != float64(i/5*5) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("Expected: ", i*5, "Got:", v)
				t.Fail()
			}
		}
	}
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
