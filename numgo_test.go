package numgo

import (
	"fmt"
	"strings"
	"testing"
)

func init() {
	debug = true
}

func TestNewArray64(t *testing.T) {
	shp := []int{2, 3, 4}
	a := NewArray64(nil, shp...)
	if len(a.data) != 24 {
		t.Logf("Length %d, expected %d", len(a.data), 24)
		t.FailNow()
	}

	for _, v := range a.data {
		if v != 0 {
			t.Logf("Value %f, expected %d", v, 0)
			t.Fail()
		}
	}
	a = NewArray64(nil)
	if e := a.GetErr(); e != nil {
		t.Log("Unexpected error:", e)
		t.Fail()
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4})
	if e := a.Equals(Arange(5)); !e.All().At(0) {
		t.Log("Slice Assignment Failed", a.GetErr(), e)
		t.Fail()
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4}, 3)
	if e := a.Equals(Arange(3)); !e.All().At(0) {
		t.Log("Slice Assignment Failed", a.GetErr(), e)
		t.Fail()
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4, 5}, 2, -1, 3)
	if e := a.GetErr(); e != NegativeAxis {
		t.Log("Expected NegativeAxis, got:", e)
		t.Fail()
	}

	a = NewArray64(nil, 1, 2, 5, 9)
	if e := a.Equals(newArray64(1, 2, 5, 9)); !e.All().At(0) {
		t.Log("Creation has different results:", e)
		t.Fail()
	}
}
func TestFull(t *testing.T) {
	shp := []int{2, 3, 4}
	a := Full(1, shp...)
	if len(a.data) != 24 {
		t.Logf("Length %d, expected %d\n", len(a.data), 24)
		t.Fail()
	}

	for _, v := range a.data {
		if v != 1 {
			t.Logf("Value %f, expected %d\n", v, 1)
			t.Fail()
			break
		}
	}

	if e := a.Equals(full(1, 2, 3, 4)); !e.All().At(0) {
		t.Log("Full creation has different results:", e)
		t.Fail()
	}
	if e := Full(0, shp...).Equals(full(0, 2, 3, 4)); !e.All().At(0) {
		t.Log("Full creation has different results:", e)
		t.Fail()
	}
}

func TestShapes(t *testing.T) {
	shp := []int{3, 3, 4, 7}
	a := NewArray64(nil, shp...)
	for i, v := range a.shape {
		if uint64(shp[i]) != v {
			t.Log(a.shape, "!=", shp)
			t.FailNow()
		}
	}
}

func TestRandArray64(t *testing.T) {
	a := RandArray64(0, 2, []int{2, 3, -7, 12})
	if e := a.GetErr(); e != NegativeAxis {
		t.Log("Expected NegativeAxis, got:", e)
		t.Fail()
	}
}

func TestArange(t *testing.T) {
	a := Arange(24)
	if len(a.data) != 24 {
		t.Logf("Length %d.  Expected size %d\n", len(a.data), 24)
	}
	if len(a.shape) != 1 {
		t.Logf("Axis %d.  Expected %d\n", len(a.shape), 1)
	}
	for i, v := range a.data {
		if float64(i) != v {
			t.Logf("Value %f.  Expected %d\n", v, i)
		}
	}

	if e := a.Equals(Arange(1, 25).SubtrC(1)); e.All().At(0) {
		t.Log("Arange generating incorrect ranges", e)
		t.Fail()
	}

	a = Arange(24, 0)
	for i := 1; i < len(a.data); i++ {
		if a.data[i]-a.data[i-1] != -1 {
			t.Log("Stepping incorrect for negative range.", a)
			t.Fail()
		}
	}

	if e := a.Equals(Arange(-24).MultC(-1)); !e.All().At(0) {
		t.Log("Negative Arange failed", e)
		t.Fail()
	}

	a = Arange(24, 0, 2)
	if e := a.GetErr(); e != ShapeError {
		t.Log("Expected ShapeError, got", e)
		t.Fail()
	}

	a = Arange(0)
	if a.shape[0] != 1 {
		t.Log("Arange(0) shape error:", a.shape[0])
		t.Fail()
	}

	a = Arange()
	if a.shape[0] != 0 {
		t.Log("Arange() shape error:", a.shape[0])
		t.Fail()
	}
}

func TestIdent(t *testing.T) {
	tmp := Identity(0)
	if len(tmp.shape) != 2 {
		t.Log("Incorrect identity shape.", tmp.shape)
		t.Fail()
	}
	if tmp.shape[0] != 0 || tmp.shape[1] != 0 {
		t.Log("Incorrect shape values. I(0)", tmp.shape)
		t.Fail()
	}
	if len(tmp.data) > 0 {
		t.Log("Data array incorrect.", tmp.data)
		t.Fail()
	}

	tmp = Identity(1)
	if tmp.shape[0] != 1 || tmp.shape[1] != 1 {
		t.Log("Incorrect shape values. I(1)", tmp.shape)
		t.Fail()
	}
	if len(tmp.data) != 1 {
		t.Log("Data Length incorrect I(1)", len(tmp.data))
		t.Fail()
	}

	tmp = Identity(4)
	if tmp.shape[0] != 4 || tmp.shape[1] != 4 {
		t.Log("Incorrect shape values. I(4)", tmp.shape)
		t.Fail()
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i != j && tmp.At(i, j) != 0 {
				t.Log("Data Value incorrect at", i, j, len(tmp.data))
				t.Fail()
			}
			if i == j && tmp.At(i, j) != 1 {
				t.Log("Data Value incorrect at", i, j, len(tmp.data))
				t.Fail()
			}
		}
	}

	tmp = Identity(-10)
	if e := tmp.GetErr(); e != NegativeAxis {
		t.Log("Error failed.  Expected NegativeAxis, got", e)
		t.Fail()
	}
}

func TestSubArray(t *testing.T) {
	a := Arange(100).Reshape(2, 5, 10)
	b := Arange(50).Reshape(5, 10)
	c := a.SubArr(0)
	if !c.Equals(b).All().At(0) {
		t.Log("Subarray incorrect. Expected\n", b, "\nReceived\n", c)
		t.Fail()
	}

	b = Arange(50).AddC(50).Reshape(5, 10)
	c = a.SubArr(1)
	if !c.Equals(b).All().At(0) {
		t.Log("Subarray incorrect. Expected\n", b, "\nReceived\n", c)
		t.Fail()
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		a   *Array64
		str string
	}{
		{nil, "<nil>"},
		{newArray64(0), "[]"},
		{&Array64{err: DivZeroError}, "Error: " + DivZeroError.s},
		{Arange(10), fmt.Sprint(Arange(10).data)},
		{Arange(10).Reshape(2, 5), "[[0 1 2 3 4] \n [5 6 7 8 9]]"},
		{Arange(20).Reshape(2, 2, 5), "[[[0 1 2 3 4]  \n  [5 6 7 8 9]] \n\n [[10 11 12 13 14]  \n  [15 16 17 18 19]]]"},
	}

	for i, tst := range tests {
		if !tst.a.HasErr() && !strings.EqualFold(tst.a.String(), tst.str) {
			t.Log("String() gave unexpected results in test", i)
			t.Log(tst.a)
			t.Log(tst.str)
			t.Fail()
		}
	}
}
