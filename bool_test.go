package numgo

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func init() {
	debug = true
}

func rnd_bool() (sz []bool) {
	sz = make([]bool, rand.Intn(100)+1)
	for i := range sz {
		sz[i] = rand.Intn(1) == 1
	}
	return
}

func TestCreateb(t *testing.T) {
	t.Parallel()
	shp := []int{2, 3, 4}
	a := NewArrayB(nil, shp...)
	if len(a.data) != 24 {
		t.Logf("Length %d, expected %d", len(a.data), 24)
		t.FailNow()
	}

	for _, v := range a.data {
		if v {
			t.Logf("Value %v, expected %v", v, false)
			t.Fail()
		}
	}
	a = NewArrayB(nil)
	if e := a.GetErr(); e != nil {
		t.Log("Unexpected error:", e)
		t.Fail()
	}

	a = NewArrayB([]bool{false, false, true, false, false})
	if e := a.Equals(NewArrayB(nil, 5).Set(true, 2)); !e.All().At(0) {
		t.Log("Slice Assignment Failed", a.GetErr(), e)
		t.Fail()
	}

	a = NewArrayB([]bool{false, false, false, false, true}, 3)
	if e := a.Equals(NewArrayB(nil, 3)); !e.All().At(0) {
		t.Log("Slice Assignment Failed", a.GetErr(), e)
		t.Fail()
	}

	a = NewArrayB([]bool{false, false, false, false, true}, 2, -1, 3)
	if e := a.GetErr(); e != NegativeAxis {
		t.Log("Expected NegativeAxis, got:", e)
		t.Fail()
	}

	a = NewArrayB(nil, 1, 2, 5, 9)
	if e := a.Equals(newArrayB(1, 2, 5, 9)); !e.All().At(0) {
		t.Log("Creation has different results:", e)
		t.Fail()
	}
}

func TestFullb(t *testing.T) {
	t.Parallel()
	shp := []int{2, 3, 4}
	a := Fullb(true, shp...)
	if len(a.data) != 24 {
		t.Logf("Length %d, expected %d\n", len(a.data), 24)
		t.Fail()
	}

	for _, v := range a.data {
		if !v {
			t.Logf("Value %v, expected %v\n", v, true)
			t.Fail()
			break
		}
	}

	if e := a.Equals(fullb(true, 2, 3, 4)); !e.All().At(0) {
		t.Log("Full creation has different results:", e)
		t.Fail()
	}
	if e := Fullb(false, shp...).Equals(fullb(false, 2, 3, 4)); !e.All().At(0) {
		t.Log("Full creation has different results:", e)
		t.Fail()
	}
}

func TestShapesB(t *testing.T) {
	shp := []int{3, 3, 4, 7}
	a := NewArrayB(nil, shp...)
	for i, v := range a.shape {
		if uint64(shp[i]) != v {
			t.Log(a.shape, "!=", shp)
			t.FailNow()
		}
	}
}

func TestStringB(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a   *Arrayb
		str string
	}{
		{nil, "<nil>"},
		{newArrayB(0), "[]"},
		{&Arrayb{err: DivZeroError}, "Error: " + DivZeroError.s},
		{Fullb(true, 10), fmt.Sprint(Fullb(true, 10).data)},
		{Fullb(false, 10).Reshape(2, 5), "[[false false false false false] \n [false false false false false]]"},
		{Fullb(true, 20).Reshape(2, 2, 5), "[[[true true true true true]  \n  [true true true true true]] \n\n [[true true true true true]  \n  [true true true true true]]]"},
		{&Arrayb{}, "<nil>"},
	}

	for i, tst := range tests {
		if !strings.EqualFold(tst.a.String(), tst.str) {
			t.Log("String() gave unexpected results in test", i)
			t.Log(tst.a)
			t.Log(tst.str)
			t.Fail()
		}
	}
}

func TestReshapeB(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a   *Arrayb
		sh  []int
		err error
	}{
		{Fullb(false, 10), []int{2, 5}, nil},
		{Fullb(false, 11), []int{2, 5}, ReshapeError},
		{Fullb(false, 10), []int{2, -5}, NegativeAxis},
		{&Arrayb{err: DivZeroError}, []int{0}, DivZeroError},
		{nil, []int{1}, NilError},
	}

	for i, tst := range tests {
		tst.a.Reshape(tst.sh...)
		if e := tst.a.GetErr(); e != tst.err {
			t.Log("Error incorrect in test", i, ", expected", tst.err, "\ngot", e)
			t.Fail()
		}
		if tst.err != nil {
			continue
		}
		for j, v := range tst.a.shape {
			if v != uint64(tst.sh[j]) {
				t.Log("Reshape incorrect in test", i, ", expected", tst.sh, "got", tst.a.shape)
				t.Fail()
				break
			}
		}
	}
}

func TestCb(t *testing.T) {
	for i := 0; i < 20; i++ {
		a := NewArrayB(rnd_bool())
		b := a.C()
		if v := a.strides[0] - a.C().strides[0]; v != 0 {
			t.Log("Size Changed", v)
			t.Fail()
		}
		if v := a.Equals(b); !v.All().At(0) {
			t.Log("Data Changed", v)
			t.Fail()
		}
	}
	if e := newArrayB(10).Reshape(0).C().GetErr(); e != ReshapeError {
		t.Log("Error failed.  Expected ReshapeError received ", e)
		t.Fail()
	}
}

func TestAtb(t *testing.T) {
	a := fullb(true, 5, 5, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		v := a.At(x, y, z)
		if !v && !a.HasErr() {
			t.Logf("Value %d failed.  Expected: %v Received: %v", i, true, v)
			t.Log(x, y, z)
			t.Fail()
		}
		if e := a.GetErr(); (x > 4 || y > 4 || z > 4) && e != IndexError {
			t.Log("Error failed.  Expected IndexErr Received", e)
			t.Log(x, y, z)
			t.Fail()
		}
	}

	_ = a.At(3, 2, 1, 0)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Error failed.  Expected InvIndexErr Received", e)
		t.Fail()
	}
}

func TestSliceElementb(t *testing.T) {
	a := Fullb(true, 5, 5, 5)
	for i := 0; i < 20; i++ {
		x, y := rand.Intn(6), rand.Intn(6)
		val := a.SliceElement(x, y)
		for i, v := range val {
			if !v && !a.HasErr() {
				t.Logf("Value %d failed.  Expected: %v Received: %v", i, true, v)
				t.Log(x, y)
				t.Fail()
			}
		}
		if e := a.GetErr(); (x > 4 || y > 4) && e != IndexError {
			t.Log("Error failed.  Expected IndexErr Received", e)
			t.Log(x, y)
			t.Log(val)
			t.Fail()
		}
	}
	val := a.SliceElement(0, 0, 0)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Error failed.  Expected InvIndexErr Received", e)
		t.Log(val)
		t.Fail()
	}

}

func TestSubArrb(t *testing.T) {
	a := newArrayB(5, 5, 5)
	g, b := false, false

	for i := 0; i < 20 || !g || !b; i++ {
		x, y := rand.Intn(6), rand.Intn(6)
		val := a.SubArr(x, y)
		if val.Any().At(0) && !val.HasErr() {
			t.Logf("Value %d failed.  Expected: %v Received: %v", i, false, val)
			t.Log(x, y)
			t.Fail()
		}
		if a.HasErr() {
			g = true
		} else {
			b = true
		}
		if e := a.GetErr(); (x > 4 || y > 4) && e != IndexError {
			t.Log("Error failed.  Expected IndexErr Received", e)
			t.Log(x, y)
			t.Log(val)
			t.Fail()
		}
	}
	_ = a.Reshape(0).SubArr(0)
	if e, d, s := a.GetDebug(); e != ReshapeError {
		t.Log("ReshapeError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
	_ = a.SubArr(0, 3, 2, 1)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("InvIndexError failed.  Received", e)
		t.Fail()
	}
}

func TestSetb(t *testing.T) {
	a := NewArrayB(nil, 5, 5, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		val := rand.Intn(1) == 1
		v := a.Set(val, x, y, z)
		if v.At(x, y, z) != val && !a.HasErr() {
			t.Logf("Value %d failed.  Expected: %v Received: %v", i, v.At(x, y, z), val)
			t.Log(x, y, z)
			t.Fail()
		}
		if e := a.GetErr(); (x > 4 || y > 4 || z > 4) && e != IndexError {
			t.Log("Error failed.  Expected IndexErr Received", e)
			t.Log(x, y, z)
			t.Fail()
		}
	}

	_ = a.Reshape(0).Set(true, 1, 1, 1)
	if e, d, s := a.GetDebug(); e != ReshapeError {
		t.Log("ReshapeError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
	_ = a.Set(true, 0, 0, 0, 0)
	if e, d, s := a.GetDebug(); e != InvIndexError {
		t.Log("InvIndexError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
}
