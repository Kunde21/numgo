package numgo

import (
	"fmt"
	"testing"
)

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
	a := Create(2, 3, 4, 5)
	fmt.Println(a.Equals(Arange(5*4*3*2).Reshape(2, 3, 4, 5)))
	fmt.Println(a.Equals(Arange(5*4*3*2).Reshape(2, 3, 4, 5)).Any(0, 2))
	fmt.Println(a.Equals(Arange(5*4*3*2).Reshape(2, 3, 4, 5)).All(0, 2))
	fmt.Println(a.Equals(Create(2, 3, 4, 5)).Any(0, 3))
}
