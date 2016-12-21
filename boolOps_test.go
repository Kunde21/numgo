package numgo

import (
	"math/rand"
	"testing"
)

func init() {
	debug = true
}

func TestEquals(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), true, true, nil},
		{a, a.C().AddC(1), false, false, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true, nil},
		{a, Arange(0, 18, 2), true, false, nil},
		{a, Arange(27, 9, -2), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Equals(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestNotEq(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), false, false, nil},
		{a, a.C().AddC(1), true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false, nil},
		{a, Arange(0, 18, 2), true, false, nil},
		{a, Arange(27, 9, -2), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.NotEq(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestLess(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), false, false, nil},
		{a, a.C().AddC(1), true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false, nil},
		{a, Arange(0, 18, 2), true, false, nil},
		{a, Arange(27, 9, -2), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Less(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestLessEq(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), true, true, nil},
		{a, a.C().AddC(1), true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true, nil},
		{a, Arange(0, 18, 2), true, true, nil},
		{a, Arange(27, 9, -2), true, true, nil},
		{a, Arange(-10, 18, 3), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.LessEq(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, "\n", v.b.data, "\n", c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestGreater(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), false, false, nil},
		{a, a.C().AddC(-1), true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false, nil},
		{a, Arange(0, 18, 2), false, false, nil},
		{a, Arange(27, 9, -2), false, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Greater(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestGreaterEq(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), true, true, nil},
		{a, a.C().AddC(-1), true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true, nil},
		{a, Arange(0, 18, 2), true, false, nil},
		{a, Arange(27, 9, -2), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.GreaterEq(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestCompValid(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b *Array64
		e    bool
		err  error
	}{
		{a, a.C(), false, nil},
		{a.C().Reshape(2, 5), a.C(), true, ShapeError},
		{a.C(), a.C().Reshape(2, 5), true, ShapeError},
		{a.C(), nil, true, NilError},
		{nil, a.C(), true, NilError},
		{a.C(), &Array64{nDimMetadata{err: InvIndexError}, nil}, true, InvIndexError},
		{&Array64{nDimMetadata{err: InvIndexError}, nil}, a.C(), true, InvIndexError},
		{a.C().Reshape(5, 2), a.C().Reshape(2, 5), true, ShapeError},
		{a.C().Reshape(5, 5), a, true, ReshapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Equals(v.b)
		e := c.HasErr()
		if e != v.e {
			t.Logf("HasErr failed in test %d.  Expected %v got %v\n", i, v.e, c.HasErr())
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d)
			t.Log(s)
			t.Log(v.a)
			t.Fail()
		}
	}
}

func TestAny(t *testing.T) {
	a := newArrayB(10).Reshape(2, 5)

	tests := []struct {
		a, b *Arrayb
		ax   []int
		err  error
	}{
		{a.C().Set(true, 0, 1), NewArrayB([]bool{true, false}), []int{1}, nil},
		{a.C().Set(true, 1, 4), NewArrayB([]bool{false, true}), []int{1}, nil},
		{a.C(), NewArrayB(nil, 5), []int{0}, nil},
		{nil, nil, []int{}, NilError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Any(v.ax...)
		if d := c.Equals(v.b); !d.All().At(0) && !c.HasErr() {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.b, d)
			t.Log(v.a, "\n", v.ax)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}

	for i := 0; i < 50; i++ {
		idx := rand.Intn(len(a.data))
		a.data[idx] = true
		b := a.C().Any(1)
		if b.At(0) == b.At(1) {
			t.Logf("Any #%d failed.  Index %d gave %v, %v\n", i, idx, b.At(0), b.At(1))
			t.Log(a)
			t.Log(b)
			t.Fail()
		}
		a.data[idx] = false
	}
}

func TestAll(t *testing.T) {
	a := fullb(true, 10).Reshape(2, 1, 5)

	for i := 0; i < 50; i++ {
		idx := rand.Intn(len(a.data))
		a.data[idx] = false
		b := a.C().All(2)
		if b.At(0, 0) == b.At(1, 0) {
			t.Logf("All #%d failed.  Index %d gave %v, %v\n", i, idx, b.At(0), b.At(1))
			t.Log(a.data, a.shape)
			t.Log(b.data)
			t.Fail()
		}
		a.data[idx] = true
	}

}

func TestBoolEquals(t *testing.T) {

	a, b := newArrayB(10), fullb(true, 10)

	tests := []struct {
		a, b     *Arrayb
		any, all bool
		err      error
	}{
		{a, a.C(), true, true, nil},
		{a, b, false, false, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true, nil},
		{a, a.C().Set(true, 5), true, false, nil},
		{b, b.C().Set(false, 7), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Equals(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestBoolNotEq(t *testing.T) {

	a, b := newArrayB(10), fullb(true, 10)

	tests := []struct {
		a, b     *Arrayb
		any, all bool
		err      error
	}{
		{a, a.C(), false, false, nil},
		{a, b, true, true, nil},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false, nil},
		{a.C().Reshape(2, 5), b.C().Reshape(2, 5), true, true, nil},
		{a, a.C().Set(true, 5), true, false, nil},
		{b, b.C().Set(false, 7), true, false, nil},
		{nil, a, false, false, NilError},
		{a, nil, false, false, NilError},
		{a.C().Reshape(2, 5), a, false, false, ShapeError},
		{a.C().Reshape(2, 5), a.C().Reshape(5, 2), false, false, ShapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.NotEq(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(v.a.data, v.b.data, c.data)
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error failed.  Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d, "\n", s, "\n", v.a, "\n", v.b)
			t.Fail()
		}
	}
}

func TestBoolCompValid(t *testing.T) {

	a := newArrayB(10)

	tests := []struct {
		a, b *Arrayb
		e    bool
		err  error
	}{
		{a, a.C(), false, nil},
		{a.C().Reshape(2, 5), a.C(), true, ShapeError},
		{a.C(), a.C().Reshape(2, 5), true, ShapeError},
		{a.C(), nil, true, NilError},
		{nil, a.C(), true, NilError},
		{a.C(), &Arrayb{nDimMetadata{err: InvIndexError}, nil}, true, InvIndexError},
		{&Arrayb{nDimMetadata{err: InvIndexError}, nil}, a.C(), true, InvIndexError},
		{a.C().Reshape(5, 2), a.C().Reshape(2, 5), true, ShapeError},
		{a.C().Reshape(5, 5), a, true, ReshapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Equals(v.b)
		e := c.HasErr()
		if e != v.e {
			t.Logf("HasErr failed in test %d.  Expected %v got %v\n", i, v.e, c.HasErr())
			t.Fail()
		}
		if e, d, s := c.GetDebug(); e != v.err {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d)
			t.Log(s)
			t.Log(v.a)
			t.Fail()
		}
	}
}

func TestBoolValAxis(t *testing.T) {

	a := newArrayB(10).Reshape(1, 1, 5, 1, 2)

	tests := []struct {
		a   *Arrayb
		ax  []int
		e   bool
		err error
	}{
		{a, []int{}, false, nil},
		{a.C().Reshape(2, 5), []int{1, 2, 3}, true, ShapeError},
		{a.C(), []int{1, 2, 4, 4, 5, 6, 7}, true, ShapeError},
		{a.C().Reshape(5, 2), []int{1, 1, 0, 0, 1}, false, nil},
		{nil, []int{1}, true, NilError},
		{a.C().Reshape(10), []int{1}, true, IndexError},
		{a.C(), []int{0, 5, 1}, true, IndexError},
		{&Arrayb{nDimMetadata{err: InvIndexError}, nil}, []int{0}, true, InvIndexError},
		{a.C().Reshape(5, 5), []int{1}, true, ReshapeError},
	}

	var c *Arrayb
	for i, v := range tests {
		v.a.valAxis(&v.ax, "")
		e := v.a.HasErr()
		if e != v.e {
			t.Logf("HasErr failed in test %d.  Expected %v got %v\n", i, v.e, c.HasErr())
			t.Fail()
		}
		if e, d, s := v.a.GetDebug(); e != v.err {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.err, e)
			t.Log("Debug:", d)
			t.Log(s)
			t.Log(v.a)
			t.Fail()
		}
	}
}
