package numgo

import (
	"math/rand"
	"testing"
)

func init() {
	debug = true
}

func rnd() (sz []int) {
	sz = make([]int, rand.Intn(8)+1)
	for i := range sz {
		sz[i] = rand.Intn(10) + 1
	}
	return
}

func TestFlatten(t *testing.T) {
	for i := 0; i < 20; i++ {
		a := RandArray64(rand.Float64()*100, rand.Float64()*100, rnd()...)
		if v := a.C().Count().Subtr(a.C().Flatten().Count()); v.At(0) != 0 {
			t.Log("Size Changed", v)
			t.Fail()
		}
	}

	if e := Arange(10).Reshape(0).Flatten().GetErr(); e != ReshapeError {
		t.Log("Error failed.  Expected ReshapeError received ", e)
		t.Fail()
	}
}

func TestC(t *testing.T) {
	for i := 0; i < 20; i++ {
		r := rnd()
		a := RandArray64(rand.Float64()*100, rand.Float64()*100, r...)
		b := a.C()
		if v := a.Count().Subtr(a.C().Count()); v.At(0) != 0 {
			t.Log("Size Changed", v)
			t.Fail()
		}
		if v := a.Equals(b); !v.All().At(0) {
			t.Log("Data Changed", v)
			t.Fail()
		}
	}
	if e := Arange(10).Reshape(0).C().GetErr(); e != ReshapeError {
		t.Log("Error failed.  Expected ReshapeError received ", e)
		t.Fail()
	}
}

func TestShape(t *testing.T) {
	var a *Array64
	for i := 0; i < 20; i++ {
		sz := rnd()
		a = NewArray64(nil, sz...)
		for i, v := range a.Shape() {
			if a.shape[i] != v {
				t.Log("Change at", i, "was", a.shape[i], "is", v)
				t.Fail()
			}
		}
		ch := rand.Intn(len(sz))
		sh := a.Shape()
		sh[ch]--
		for i, v := range a.shape {
			if sh[i] != v && i != ch {
				t.Log("Change at", i, "was", a.shape[i], "is", v)
				t.Fail()
			}
			if sh[i] == v && i == ch {
				t.Log("Change propagated at", i, "was", a.shape[i], "is", v)
				t.Fail()
			}
		}
	}
	sh := a.Reshape(-1).Shape()
	if !a.HasErr() || sh != nil {
		t.Log("Shape() error handling incorrect")
		t.Log("Shape:", sh, "Err:", a.getErr())
		t.Fail()
	}
}

func TestAt(t *testing.T) {
	t.Parallel()
	a := Arange(125).Reshape(5, 5, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		v := a.At(x, y, z)
		if v != float64(x*25+y*5+z) && !a.HasErr() {
			t.Logf("Value %d failed.  Expected: %v Received: %v", i, float64(x*25+y*5+z), v)
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

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				if a.at([]int{i, j, k}) != a.At(i, j, k) {
					t.Log("at failed for index", i, j, k)
					t.Log(a.at([]int{i, j, k}), "!=", a.At(i, j, k))
					t.Fail()
				}
			}
		}
	}
}

func TestSliceElement(t *testing.T) {
	a := Arange(125).Reshape(5, 5, 5)
	for i := 0; i < 20; i++ {
		x, y := rand.Intn(6), rand.Intn(6)
		val := a.SliceElement(x, y)
		for i, v := range val {
			if v != float64(x*25+y*5+i) && !a.HasErr() {
				t.Logf("Value %d failed.  Expected: %v Received: %v", i, float64(x*25+y*5+i), v)
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
}

func TestSubArr(t *testing.T) {
	a := Arange(125).Reshape(5, 5, 5)
	g, b := false, false

	for i := 0; i < 20 || !g || !b; i++ {
		x, y := rand.Intn(6), rand.Intn(6)
		val := a.SubArr(x, y)
		if v := val.Equals(Arange(float64(x*25+y*5), float64(x*25+y*5+5))); !v.All().At(0) && !v.HasErr() {
			t.Logf("Value %d failed.  Expected: %v Received: %v", i, float64(x*25+y*5+i), v)
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

func TestSet(t *testing.T) {
	a := NewArray64(nil, 5, 5, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		val := rand.Float64() * 100
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

	_ = a.Reshape(0).Set(0, 1, 1, 1)
	if e, d, s := a.GetDebug(); e != ReshapeError {
		t.Log("ReshapeError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
	_ = a.Set(0, 0, 0, 0, 0)
	if e, d, s := a.GetDebug(); e != InvIndexError {
		t.Log("InvIndexError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
}

func TestSetSliceElement(t *testing.T) {
	a := NewArray64(nil, 5, 5, 3, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(4)
		val := RandArray64(5, 100, []int{5}...).SliceElement()
		v := a.SetSliceElement(val, x, y, z)
		if !a.HasErr() {
			for j, k := range v.SliceElement(x, y, z) {
				if k != val[j] {
					t.Logf("Value %d failed.  Expected: %v Received: %v", i, v.At(x, y, z), val)
					t.Log(x, y, z)
					t.Fail()
				}
			}
		}
		if e := a.GetErr(); (x > 4 || y > 4 || z > 2) && e != IndexError {
			t.Log("Error failed.  Expected IndexErr Received", e)
			t.Log(x, y, z)
			t.Fail()
		}
	}

	_ = a.Reshape(0).SetSliceElement(nil, 1, 1, 1)
	if e, d, s := a.GetDebug(); e != ReshapeError {
		t.Log("ReshapeError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
	_ = a.SetSliceElement(nil, 0, 0, 0, 0)
	if e, d, s := a.GetDebug(); e != InvIndexError {
		t.Log("InvIndexError failed.  Received", e)
		t.Log(d, "\n", s)
		t.Fail()
	}
}

func TestSetSubArr(t *testing.T) {
	a := NewArray64(nil, 5, 5, 3, 5)

	b := Arange(15).Reshape(3, 5)
	a.SetSubArr(b, 0, 1)
	if a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
	for i := range a.data {
		if i >= 15 && i < 30 && a.data[i] != float64(i-15) {
			t.Log("Failed at:", i, a.data[i])
			t.Fail()
		}
		if !(i >= 15 && i < 30) && a.data[i] != 0 {
			t.Log("Not 0 Failed at:", i, a.data[i])
			t.Fail()
		}
	}
	a.SetSubArr(b, 0)
	for i := range a.data {
		if i >= 0 && i < 5*15 && a.data[i] != float64(i%15) {
			t.Log("Failed 2 at:", i, a.data[i])
			t.Fail()
		}
		if !(i >= 0 && i < 5*15) && a.data[i] != 0 {
			t.Log("Not 0 Failed 2 at:", i, a.data[i])
			t.Fail()
		}
	}
	a.SetSubArr(b, 1, 1, 1)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Did not error correctly.  Expected InvIndexError, got ", e)
		t.Fail()
	}

	a.SetSubArr(b.Reshape(5, 3), 0, 1)
	if e := a.GetErr(); e != ShapeError {
		t.Log("Did not error correctly.  Expected ShapeError, got ", e)
		t.Fail()
	}
	b.err = InvIndexError
	a.SetSubArr(b.Reshape(3, 5), 0, 1)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Did not error correctly.  Expected InvIndexError, got ", e)
		t.Fail()
	}
	b.err, a = nil, nil
	a.SetSubArr(b, 0, 1)
	if e := a.GetErr(); e != NilError {
		t.Log("Did not error correctly.  Expected NilError, got ", e)
		t.Fail()
	}
}

func TestResize(t *testing.T) {
	a := NewArray64(nil, 5, 5, 3, 5)

	a.Resize(-1)
	if e := a.GetErr(); e != NegativeAxis {
		t.Log("Negative axis failed to error", e)
		t.Fail()
	}
	a.Resize(5, 3, 2, -10)
	if e := a.GetErr(); e != NegativeAxis {
		t.Log("Negative axis failed to error", e)
		t.Fail()
	}

	a.Set(1, 0, 0, 0, 2).Resize(5, 5)
	if a.HasErr() {
		t.Log("Error in set/resize", a.GetErr())
		t.Fail()
	}
	_ = a.At(0, 0, 0, 2)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Bad Error after resize", e)
		t.Fail()
	}
	if c := a.At(0, 2); c != 1 {
		t.Log("Data didn't move correctly in reduction.  Expected 1, got", c)
		t.Fail()
	}
	if c := a.Resize(5, 5, 3).At(0, 0, 2); c != 1 {
		t.Log("Data didn't move correctly in small expansion.  Expected 1, got", c)
		t.Fail()
	}
	if c := a.Resize(5, 5, 3, 5, 10).At(0, 0, 0, 0, 2); c != 1 {
		t.Log("Data didn't move correctly in large expansion.  Expected 1, got", c)
		t.Fail()
	}
	a.Resize().At(0)
	if e := a.GetErr(); e != IndexError {
		t.Log("Did not error correctly.  Expected IndexError, got ", e)
		t.Fail()
	}

	a.err = InvIndexError
	if e := a.Resize(10).GetErr(); e != InvIndexError {
		t.Log("Error didn't pass through correctly.  Expected InvIndexError, got", e)
		t.Fail()
	}
}

func TestAppend(t *testing.T) {
	a := NewArray64(nil, 1, 2, 3, 4, 5)
	b := Arange(120)

	a.Append(nil, 1)
	if e := a.GetErr(); e != NilError {
		t.Log("Expected NilError, received", e)
		t.Fail()
	}

	a.Append(b, 1)
	if e := a.GetErr(); e != ShapeError {
		t.Log("Expected ShapeError, received", e)
		t.Fail()
	}

	a.Append(b, 5)
	if e := a.GetErr(); e != IndexError {
		t.Log("Expected IndexError, received", e)
		t.Fail()
	}

	a.Append(nil, -1)
	if e := a.GetErr(); e != IndexError {
		t.Log("Expected IndexError, received", e)
		t.Fail()
	}

	a.Append(b.Reshape(5, 4, 3, 2, 1), 1)
	if e := a.GetErr(); e != ShapeError {
		t.Log("Expected ShapeError, received", e)
		t.Fail()
	}

	a.Append(b.Reshape(1, 2, 1, 3, 4, 5), 2)
	if e := a.GetErr(); e != ShapeError {
		t.Log("Expected ShapeError, received", e)
		t.Fail()
	}

	a.Append(b.Reshape(1, 2, 3, 4, 5), 2)
	if e := a.GetErr(); e != nil {
		t.Log("Unexpected Error: ", e)
		t.Fail()
	}
	if a.shape[2] != 6 {
		t.Log("Shape updated incorrectly.  Expected 6, got", a.shape[2])
		t.Fail()
	}

	a.Resize(1, 2, 3, 4, 5).Append(b, 0)
	if e := a.GetErr(); e != nil {
		t.Log("Unexpected Error: ", e)
		t.Fail()
	}
	if a.shape[0] != 2 {
		t.Log("Shape updated incorrectly.  Expected 2, got", a.shape[0])
		t.Fail()
	}

	a.err = InvIndexError
	a.Append(b, 0)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Expected InvIndexError, received", e)
		t.Fail()
	}

	a = NewArray64([]float64{1, 2, 3, 3, 2, 1, 2, 1, 3}, 3, 3)
	b = NewArray64([]float64{5, 6, 5, 4, 6, 5}, 3, 2)
	c := NewArray64([]float64{
		1, 2, 3, 5, 6,
		3, 2, 1, 5, 4,
		2, 1, 3, 6, 5},
		3, 5)

	if !a.Append(b, 1).Equals(c).All().At(0) {
		t.Log("Append gave unexpected results")
		t.Log(a.Append(b, 1).Equals(c))
		t.Fail()
	}
}
