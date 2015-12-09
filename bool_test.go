package numgo

import "testing"

func TestEquals(t *testing.T) {
	a, b := Arange(10), Arange(10)

	c := a.Equals(b)
	for _, v := range c.data {
		if !v {
			t.Log("Equals expected equivalence, got", c.data)
			t.FailNow()
		}
	}

	if !c.Any().data[0] {
		t.Log("Any expected true, got false.", c)
		t.FailNow()
	}
	if !c.All().data[0] {
		t.Log("All expected true, got false.", c)
		t.FailNow()
	}

	c = a.Equals(b.AddC(1))
	if c.Any().data[0] {
		t.Log("Any expected false, got true", c)
		t.FailNow()
	}
	if c.All().data[0] {
		t.Log("Any expected false, got true", c)
		t.FailNow()
	}

	c = a.Equals(Arange(0, 20, 2))
	if !c.Any().data[0] {
		t.Log("Any expected true, got false", c)
		t.FailNow()
	}
	if c.All().data[0] {
		t.Log("Any expected false, got true", c)
		t.FailNow()
	}

	c = a.Equals(Arange(27, 7, -2))
	if !c.Any().data[0] {
		t.Log("Any expected true, got false", c)
		t.FailNow()
	}
	if c.All().data[0] {
		t.Log("Any expected false, got true", c)
		t.FailNow()
	}
}

func TestA(t *testing.T) {
	sz := []int{2, 3, 4, 5}
	a := Create(sz...)
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

	b = a.Equals(Create(2, 3, 4, 5)).Any(0, 3)
	for _, v := range b.data {
		if !v {
			t.Log("First value. Expected true, got", v)
			t.Fail()
		}
	}
}
