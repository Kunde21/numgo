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
	}{
		{a, a.C(), true, true},
		{a, a.C().AddC(1), false, false},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true},
		{a, Arange(0, 20, 2), true, false},
		{a, Arange(27, 7, -2), true, false},
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
	}
}

func TestLess(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
	}{
		{a, a.C(), false, false},
		{a, a.C().AddC(1), true, true},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false},
		{a, Arange(0, 20, 2), true, false},
		{a, Arange(27, 7, -2), true, false},
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
	}
}

func TestLessEq(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
	}{
		{a, a.C(), true, true},
		{a, a.C().AddC(1), true, true},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true},
		{a, Arange(0, 20, 2), true, true},
		{a, Arange(27, 7, -2), true, true},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.LessEq(v.b)
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
	}
}

func TestGreater(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
	}{
		{a, a.C(), false, false},
		{a, a.C().AddC(-1), true, true},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), false, false},
		{a, Arange(0, 20, 2), false, false},
		{a, Arange(27, 7, -2), false, false},
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
	}
}

func TestGreaterEq(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
	}{
		{a, a.C(), true, true},
		{a, a.C().AddC(-1), true, true},
		{a.C().Reshape(2, 5), a.C().Reshape(2, 5), true, true},
		{a, Arange(0, 20, 2), true, false},
		{a, Arange(27, 7, -2), true, false},
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
		{a.C(), &Array64{err: DivZeroError}, true, DivZeroError},
		{&Array64{err: DivZeroError}, a.C(), true, DivZeroError},
		{a.C().Reshape(5, 2), a.C().Reshape(2, 5), true, ShapeError},
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

	for i := 0; i < 50; i++ {
		idx := rand.Intn(len(a.data))
		a.data[idx] = true
		b := a.C().Any(1)
		if b.At(0) == b.At(1) {
			t.Logf("Any #%d failed.  Index %d gave %v, %v\n", i, idx, b.At(0), b.At(1))
			t.Log(a.data, a.shape)
			t.Log(b.data)
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

func TestA(t *testing.T) {
	sz := []int{2, 3, 4, 5}
	a := NewArray64(nil, sz...)
	b := a.Equals(Arange(5 * 4 * 3 * 2).Reshape(sz...))

	for i, v := range b.shape {
		if int(v) != sz[i] {
			t.Log("Shape incorrect")
			t.Log("Expected:", sz)
			t.Log("Received:", b.shape)
		}
	}

	b = a.Equals(Arange(5*4*3*2).Reshape(2, 3, 4, 5)).Any(0, 2)
	for i, v := range b.data {
		if i == 0 && !v {
			t.Log("First value. Expected true, got", v)
			t.Fail()
		}
		if i > 0 && v {
			t.Log("First value. Expected false, got", v)
			t.Fail()
		}
	}

	b = a.Equals(Arange(5*4*3*2).Reshape(2, 3, 4, 5)).All(0, 2)
	for _, v := range b.data {
		if v {
			t.Log("First value. Expected false, got", v)
			t.Fail()
		}
	}

	b = a.Equals(NewArray64(nil, 2, 3, 4, 5)).Any(0, 3)
	for _, v := range b.data {
		if !v {
			t.Log("First value. Expected true, got", v)
			t.Fail()
		}
	}
}

func TestDebug(t *testing.T) {
	Debug(true)
	var nilp *Array64
	nilp.Set(12, 1, 4, 0).AddC(2).DivC(6).At(1, 4, 0)
	if !nilp.HasErr() {
		t.FailNow()
		err, debug, stack := nilp.GetDebug()
		t.Log(err)   // Prints generic error: "Nil pointer received."
		t.Log(debug) // Prints debug info: "Nil pointer received by SetE()."
		t.Log(stack)
		t.Fail()
	}
	nilp = MinSet(Arange(10).Reshape(2, 5), Arange(10))
	if err, debug, stack := nilp.GetDebug(); err != ShapeError {
		t.Log(err)
		t.Log(debug)
		t.Log(stack)
		t.Fail()
	}

}
