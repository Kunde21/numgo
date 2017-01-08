package numgo

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
)

func init() {
	debug = true
}

func TestNewArray64(t *testing.T) {
	t.Parallel()
	shp := []int{2, 3, 4}
	a := NewArray64(nil, shp...)
	if len(a.data) != 24 {
		t.Logf("Length %d, expected %d", len(a.data), 24)
		t.FailNow()
	}

	for _, v := range a.data {
		if v != 0 {
			t.Errorf("Value %f, expected %d", v, 0)
		}
	}
	a = NewArray64(nil)
	if e := a.GetErr(); e != nil {
		t.Error("Unexpected error:", e)
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4})
	if e := a.Equals(Arange(5)); !e.All().At(0) {
		t.Error("Slice Assignment Failed", a.GetErr(), e)
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4}, 3)
	if e := a.Equals(Arange(3)); !e.All().At(0) {
		t.Error("Slice Assignment Failed", a.GetErr(), e)
	}

	a = NewArray64([]float64{0, 1, 2, 3, 4, 5}, 2, -1, 3)
	if e := a.GetErr(); e != NegativeAxis {
		t.Error("Expected NegativeAxis, got:", e)
	}

	a = NewArray64(nil, 1, 2, 5, 9)
	if e := a.Equals(newArray64(1, 2, 5, 9)); !e.All().At(0) {
		t.Error("Creation has different results:", e)
	}
}
func TestFull(t *testing.T) {
	t.Parallel()
	shp := []int{2, 3, 4}
	a := FullArray64(1, shp...)
	if len(a.data) != 24 {
		t.Errorf("Length %d, expected %d\n", len(a.data), 24)
	}

	for _, v := range a.data {
		if v != 1 {
			t.Errorf("Value %f, expected %d\n", v, 1)
			break
		}
	}

	if e := a.Equals(full(1, 2, 3, 4)); !e.All().At(0) {
		t.Error("Full creation has different results:", e)
	}
	if e := FullArray64(0, shp...).Equals(full(0, 2, 3, 4)); !e.All().At(0) {
		t.Error("Full creation has different results:", e)
	}
}

func TestShapes(t *testing.T) {
	shp := []int{3, 3, 4, 7}
	a := NewArray64(nil, shp...)
	for i, v := range a.shape {
		if shp[i] != v {
			t.Log(a.shape, "!=", shp)
			t.FailNow()
		}
	}
}

func TestRandArray64(t *testing.T) {
	t.Parallel()
	a := RandArray64(0, 2, []int{2, 3, -7, 12}...)
	if e := a.GetErr(); e != NegativeAxis {
		t.Error("Expected NegativeAxis, got:", e)
	}
}

func TestArange(t *testing.T) {
	t.Parallel()
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
		t.Error("Arange generating incorrect ranges", e)
	}

	a = Arange(-24)
	if len(a.data) != 24 {
		t.Logf("Length %d.  Expected size %d\n", len(a.data), 24)
	}
	if len(a.shape) != 1 {
		t.Logf("Axis %d.  Expected %d\n", len(a.shape), 1)
	}
	for i, v := range a.data {
		if -float64(len(a.data)-i) != v {
			t.Logf("Value %f.  Expected %d\n", v, i)
		}
	}
	if e := a.Equals(Arange(1, 25).SubtrC(1)); e.All().At(0) {
		t.Error("Arange generating incorrect ranges", e)
	}

	a = Arange(24, 0)
	for i := 1; i < len(a.data); i++ {
		if a.data[i]-a.data[i-1] != -1 {
			t.Error("Stepping incorrect for negative range.", a)
		}
	}

	if e := a.Equals(Arange(-25).MultC(-1)); !e.All().At(0) {
		t.Error("Negative Arange failed", e)
	}

	a = Arange(24, 0, 2)
	if e := a.GetErr(); e != ShapeError {
		t.Error("Expected ShapeError, got", e)
	}

	a = Arange(0)
	if a.shape[0] != 1 {
		t.Error("Arange(0) shape error:", a.shape[0])
	}

	a = Arange()
	if a.shape[0] != 0 {
		t.Error("Arange() shape error:", a.shape[0])
	}
}

func TestIdent(t *testing.T) {
	t.Parallel()
	var tmp *Array64
	for k := 0; k < 5; k++ {
		tmp = Identity(k)
		if len(tmp.shape) != 2 {
			t.Error("Incorrect identity shape.", tmp.shape)
		}
		if tmp.shape[0] != k || tmp.shape[1] != k {
			t.Error("Incorrect shape values. I()", k, tmp.shape)
		}
		if len(tmp.data) != k*k {
			t.Error("Data array incorrect.", tmp.data)
		}
		for i := 0; i < k; i++ {
			for j := 0; j < k; j++ {
				if i != j && tmp.At(i, j) != 0 {
					t.Error("Data Value incorrect at", i, j, len(tmp.data))
				}
				if i == j && tmp.At(i, j) != 1 {
					t.Error("Data Value incorrect at", i, j, len(tmp.data))
				}
			}
		}
	}

	tmp = Identity(-10)
	if e := tmp.GetErr(); e != NegativeAxis {
		t.Error("Error failed.  Expected NegativeAxis, got", e)
	}
}

func TestSubArray(t *testing.T) {
	t.Parallel()
	a := Arange(100).Reshape(2, 5, 10)
	b := Arange(50).Reshape(5, 10)
	c := a.SubArr(0)
	if !c.Equals(b).All().At(0) {
		t.Error("Subarray incorrect. Expected\n", b, "\nReceived\n", c)
	}

	b = Arange(50).AddC(50).Reshape(5, 10)
	c = a.SubArr(1)
	if !c.Equals(b).All().At(0) {
		t.Error("Subarray incorrect. Expected\n", b, "\nReceived\n", c)
	}
}

func TestString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a   *Array64
		str string
	}{
		{nil, "<nil>"},
		{newArray64(0), "[]"},
		{&Array64{err: InvIndexError}, "Error: " + InvIndexError.s},
		{Arange(10), fmt.Sprint(Arange(10).data)},
		{Arange(10).Reshape(2, 5), "[[0 1 2 3 4] \n [5 6 7 8 9]]"},
		{Arange(20).Reshape(2, 2, 5), "[[[0 1 2 3 4]  \n  [5 6 7 8 9]] \n\n [[10 11 12 13 14]  \n  [15 16 17 18 19]]]"},
		{&Array64{}, "<nil>"},
	}

	for i, tst := range tests {
		if !strings.EqualFold(tst.a.String(), tst.str) {
			t.Log("String() gave unexpected results in test", i)
			t.Log(tst.a)
			t.Error(tst.str)
		}
	}
}

func TestReshape(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a   *Array64
		sh  []int
		err error
	}{
		{Arange(10), []int{2, 5}, nil},
		{Arange(11), []int{2, 5}, ReshapeError},
		{Arange(10), []int{2, -5}, NegativeAxis},
		{&Array64{err: InvIndexError}, []int{0}, InvIndexError},
		{nil, []int{1}, NilError},
	}

	for i, tst := range tests {
		tst.a.Reshape(tst.sh...)
		if e := tst.a.GetErr(); e != tst.err {
			t.Error("Error incorrect in test", i, ", expected", tst.err, "\ngot", e)
		}
		if tst.err != nil {
			continue
		}
		for j, v := range tst.a.shape {
			if v != tst.sh[j] {
				t.Error("Reshape incorrect in test", i, ", expected", tst.sh, "got", tst.a.shape)
				break
			}
		}
	}
}

func TestJSON(t *testing.T) {
	//t.Parallel()

	tests := []*Array64{
		NewArray64(nil, 0),
		Arange(10),
		RandArray64(0, 2, ([]int{10, 10})...).Div(Arange(10)),
		Arange(10).Reshape(2, 2),
		FullArray64(math.NaN(), 10),
		FullArray64(math.Inf(1), 10),
		FullArray64(math.Inf(-1), 10),
	}
	for i, v := range tests {
		b, err := json.Marshal(v)
		if err != nil {
			t.Error("Marshal Error in test", i, ":", err)
			continue
		}
		tmp := new(Array64)
		err = json.Unmarshal(b, tmp)
		if err != nil {
			t.Error("Unmarshal Errorin test", i, ":", err)
			continue
		}

		e1, e2 := v.GetErr(), tmp.GetErr()
		if e1 != e2 {
			t.Log("Error mismatch in test", i)
			t.Log("From:", e1)
			t.Error("To:", e2)
		}

		if e := tmp.Equals(v); !e.All().At(0) {
			t.Log("Value changedin test", i)
			t.Log(string(b))
			t.Log(v)
			t.Error(tmp)
		}
	}

	var v *Array64
	b, err := json.Marshal(v)
	if err != nil {
		t.Error("Marshal Error in nil test:", err)
	}
	tmp := new(Array64)
	err = json.Unmarshal(b, tmp)
	if err != nil {
		t.Error("Unmarshal Error in nil test:", err)
	}

	e1, e2 := v.GetErr(), tmp.GetErr()
	if e1 != e2 {
		t.Log("Error mismatch in nil test")
		t.Log("From:", e1)
		t.Error("To:", e2)
	}

	b, err = json.Marshal(Arange(10))
	v = nil
	e1 = json.Unmarshal(b, v)
	if e1 == nil {
		t.Log("Empty unmarshal didn't return error:")
		t.Error("Res:", v)
	}

	v = new(Array64)
	e1 = json.Unmarshal([]byte(`{"junk": "This will not pass."}`), v)
	if e1 != nil || v.err != NilError {
		t.Log("Error unmarshal didn't error correctly:")
		t.Error(v)
	}
}
