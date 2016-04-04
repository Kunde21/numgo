package numgo

import "testing"

func TestDotProd(t *testing.T) {
	a, b := Arange(11), Arange(11).AddC(1)
	fmaSet := fmaSupt
	for i, v := range []struct {
		a *Array64
		b *Array64
		c float64
		e error
	}{
		{a, b, 440, nil},
		{a.C().Resize(10), b.C().Resize(10), 330, nil},
	} {
		if c := v.a.DotProd(v.b); c.At(0) != v.c {
			t.Log("Test", i, "Expected", v.c, "Got", c.At(0))
			t.Fail()
		}
		fmaSupt = false

		if c := v.a.DotProd(v.b); c.At(0) != v.c {
			t.Log("Test", i, "No FMA expected", v.c, "Got", c.At(0))
			t.Fail()
		}
		fmaSupt = fmaSet

		if c := a.GetErr(); c != v.e {
			t.Log("Error test", i, "Expected", v.e, "Got", c)
			t.Fail()
		}
	}
}

func BenchmarkDotProd(t *testing.B) {
	a, b := Arange(1000000), Arange(1000000)

	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = a.DotProd(b)
	}
}
