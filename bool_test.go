package numgo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func init() {
	debug = true
}

func rndBool() (sz []bool) {
	sz = make([]bool, rand.Intn(100)+10)
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
	t.Parallel()
	shp := []int{3, 3, 4, 7}
	a := NewArrayB(nil, shp...)
	for i, v := range a.shape {
		if shp[i] != v {
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
		{&Arrayb{err: InvIndexError}, "Error: " + InvIndexError.s},
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
		{&Arrayb{err: InvIndexError}, []int{0}, InvIndexError},
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
			if v != tst.sh[j] {
				t.Log("Reshape incorrect in test", i, ", expected", tst.sh, "got", tst.a.shape)
				t.Fail()
				break
			}
		}
	}
}

func TestCb(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		a := NewArrayB(rndBool())
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestSetSliceElementb(t *testing.T) {
	t.Parallel()
	a := NewArrayB(nil, 5, 5, 3, 5)

	for i := 0; i < 20; i++ {
		x, y, z := rand.Intn(6), rand.Intn(6), rand.Intn(4)
		val := rndBool()[:5]
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

func TestSetSubArrb(t *testing.T) {
	a := NewArrayB(nil, 5, 5, 3, 5)

	b := Fullb(true, 3, 5)
	a.SetSubArr(b, 0, 1)
	if a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
	for i := range a.data {
		if i >= 15 && i < 30 && !a.data[i] {
			t.Log("Failed at:", i, a.data[i])
			t.Fail()
		}
		if !(i >= 15 && i < 30) && a.data[i] {
			t.Log("Not 0 Failed at:", i, a.data[i])
			t.Fail()
		}
	}
	a.SetSubArr(b, 0)
	for i := range a.data {
		if i >= 0 && i < 5*15 && !a.data[i] {
			t.Log("Failed 2 at:", i, a.data[i])
			t.Fail()
		}
		if !(i >= 0 && i < 5*15) && a.data[i] {
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

func TestResizeb(t *testing.T) {
	a := NewArrayB(nil, 5, 5, 3, 5)

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

	a.Set(true, 0, 0, 0, 2).Resize(5, 5)
	if a.HasErr() {
		t.Log("Error in set/resize", a.GetErr())
		t.Fail()
	}
	_ = a.At(0, 0, 0, 2)
	if e := a.GetErr(); e != InvIndexError {
		t.Log("Bad Error after resize", e)
		t.Fail()
	}
	if c := a.At(0, 2); !c {
		t.Log("Data didn't move correctly in reduction.  Expected 1, got", c)
		t.Fail()
	}
	if c := a.Resize(5, 5, 3).At(0, 0, 2); !c {
		t.Log("Data didn't move correctly in small expansion.  Expected 1, got", c)
		t.Fail()
	}
	if c := a.Resize(5, 5, 3, 5, 10).At(0, 0, 0, 0, 2); !c {
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

func TestAppendb(t *testing.T) {
	a := NewArrayB(nil, 1, 2, 3, 4, 5)
	b := Fullb(true, 120)

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
}

func TestJSONb(t *testing.T) {
	t.Parallel()

	tests := []*Arrayb{
		NewArrayB(nil, 0),
		fullb(true, 10),
		newArrayB(10).Reshape(2, 2),
		Fullb(false, 10),
		Fullb(true, 10),
	}
	for i, v := range tests {
		b, err := json.Marshal(v)
		if err != nil {
			t.Log("Marshal Error in test", i, ":", err)
			t.Fail()
			continue
		}
		tmp := new(Arrayb)
		err = json.Unmarshal(b, tmp)
		if err != nil {
			t.Log("Unmarshal Errorin test", i, ":", err)
			t.Fail()
			continue
		}

		e1, e2 := v.GetErr(), tmp.GetErr()
		if e1 != e2 {
			t.Log("Error mismatch in test", i)
			t.Log("From:", e1)
			t.Log("To:", e2)
			t.Fail()
		}

		if e := tmp.Equals(v); e1 == nil && !e.All().At(0) {
			t.Log("Value changedin test", i)
			t.Log(string(b))
			t.Log(v)
			t.Log(tmp)
			t.Fail()
		}
	}

	var v *Arrayb
	b, err := json.Marshal(v)
	if err != nil {
		t.Log("Marshal Error in nil test:", err)
		t.Fail()
	}
	tmp := new(Arrayb)
	err = json.Unmarshal(b, tmp)
	if err != nil {
		t.Log("Unmarshal Error in nil test:", err)
		t.Fail()
	}

	e1, e2 := v.GetErr(), tmp.GetErr()
	if e1 != e2 {
		t.Log("Error mismatch in nil test")
		t.Log("From:", e1)
		t.Log("To:", e2)
		t.Fail()
	}

	b, err = json.Marshal(newArrayB(10))
	v = nil
	e1 = json.Unmarshal(b, v)
	if e1 == nil {
		t.Log("Empty unmarshal didn't return error:")
		t.Log("Res:", v)
		t.Fail()
	}

	v = new(Arrayb)
	e1 = json.Unmarshal([]byte(`{"junk": "This will not pass."}`), v)
	if e1 != nil || v.err != NilError {
		t.Log("Error unmarshal didn't error correctly:")
		t.Log(v)
		t.Fail()
	}
}
