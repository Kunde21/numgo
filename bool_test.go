package numgo

import (
	"fmt"
	"testing"
)

func TestEquals(t *testing.T) {

	a := Arange(10)

	tests := []struct {
		a, b     *Array64
		any, all bool
		err      error
	}{
		{a, a.C(), true, true, nil},
		{a, a.C().AddC(1), false, false, nil},
		{a.C().Reshape(2, 5), b.C().Reshape(2, 5), true, true, nil},
		{a, Arange(0, 20, 2), true, false, nil},
		{a, Arange(27, 7, -2), true, false, nil},
	}

	var c *Arrayb
	for i, v := range tests {
		c = v.a.Equals(v.b)
		if d := c.Any().At(0); d != v.any {
			t.Logf("Test %d failed.  Any expected %v got %v\n", i, v.any, d)
			t.Log(a.data, b.data, c.data)
			t.Fail()
		}
		if d := c.All().At(0); d != v.all {
			t.Logf("Test %d failed.  All expected %v got %v\n", i, v.all, d)
			t.Log(a.data, b.data, c.data)
			t.Fail()
		}
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
		err, debug := nilp.GetDebug()
		fmt.Println(err)   // Prints generic error: "Nil pointer received."
		fmt.Println(debug) // Prints debug info: "Nil pointer received by SetE()."
	}
}
