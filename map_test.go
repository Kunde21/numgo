package numgo

import "testing"

func TestCleanAxis(t *testing.T) {

	tests := []struct {
		a, b []int
	}{
		{[]int{}, []int{}},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1}, []int{1}},
		{[]int{150}, []int{150}},
		{[]int{10, 9, 8, 7, 7, 8, 9, 10}, []int{10, 9, 8, 7}},
		{[]int{1, 2, 3, 4, 5, 5, 5, 5, 5, 1, 1, 1, 1, 1, 1}, []int{1, 2, 3, 4, 5}},
	}

	for i, v := range tests {
		clean := *(cleanAxis(&v.a))
		for j, k := range clean {
			if k != v.b[j] {
				t.Logf("Test %d failed.  Expected %v received %v\n", i, v.b, clean)
				t.Fail()
				break
			}
		}
	}
}

func TestValAxis(t *testing.T) {
	t.Parallel()
tests:
	for i, v := range []struct {
		a      *Array64
		ax, re []int
		err    error
	}{
		{new(Array64), []int{0}, []int{0}, NilError},
		{Arange(10), []int{0}, []int{}, nil},
		{Arange(10), []int{1}, []int{1}, IndexError},
		{Arange(10).Reshape(2, 5), []int{1, -1}, []int{1, -1}, IndexError},
		{Arange(10).Reshape(2, 5), []int{1, 0}, []int{}, nil},
		{&Array64{err: DivZeroError}, []int{}, []int{}, DivZeroError},
	} {
		if v.a.valAxis(&v.ax, "Test"); v.a.getErr() != v.err {
			t.Log("Error mismatch.  Expected", v.err, "Got", v.a.getErr())
			t.Fail()
		}
		if len(v.ax) != len(v.re) {
			t.Log("Length incorrect.  Expected", v.re, "Got", v.ax)
			t.Fail()
			continue tests
		}
		for idx, ax := range v.ax {
			if ax != v.re[idx] {
				t.Log("Result incorrect.  Expected", v.re, "Got", v.ax)
				t.Fail()
				continue tests
			}
		}
	}
}

func TestFoldCC(t *testing.T) {

	num := func(i float64) FoldFunc {
		return func(d []float64) float64 {
			return i
		}
	}

	sum := func(d []float64) (r float64) {
		for i := range d {
			r += d[i]
		}
		return r
	}

	pan := func(d []float64) (r float64) {
		for i := range d {
			r += d[i]
		}
		return d[len(d)+1]
	}

	a := Arange(100).Reshape(2, 1, 5, 2, 5)

	tests := []struct {
		a, b *Array64
		ax   []int
		f    FoldFunc
		err  error
	}{
		{a.C(), a.C().Reshape(2, 5, 2, 5), []int{1}, sum, nil},
		{a.C(), full(2, 2, 2, 5), []int{1, 2}, num(2), nil},
		{a.C(), a, []int{1}, pan, FoldMapError},
		{a.C(), a.C().Sum(0, 1, 2), []int{0, 1, 2}, sum, nil},
		{a.C().Flatten(), a, []int{1}, sum, IndexError},
		{a.C().Flatten(), a, []int{0, 1, 2}, sum, ShapeError},
	}

	for i, v := range tests {
		r := v.a.FoldCC(v.f, v.ax...)
		if !r.Equals(v.b).All().At(0) && !r.HasErr() {
			t.Logf("Test %d failed.  \nExpected:\n %v \nReceived:\n %v\n", i, v.b, r)
			t.Fail()
		}
		if e, d, s := r.GetDebug(); e != v.err {
			t.Logf("Test %d error failed.  Expected: %v Received: %v\n", i, v.err, e)
			t.Log(d, "\n", s, "\n", r)
			t.Fail()
		}
	}

}

func TestFold(t *testing.T) {

	num := func(i float64) FoldFunc {
		return func(d []float64) float64 {
			return i
		}
	}

	sum := func(d []float64) (r float64) {
		for i := range d {
			r += d[i]
		}
		return r
	}

	pan := func(d []float64) (r float64) {
		for i := range d {
			r += d[i]
		}
		return d[len(d)+1]
	}

	a := Arange(100).Reshape(2, 1, 5, 2, 5)

	tests := []struct {
		a, b *Array64
		ax   []int
		f    FoldFunc
		err  error
	}{
		{a.C(), a.C().Reshape(2, 5, 2, 5), []int{1}, sum, nil},
		{a.C(), full(2, 2, 2, 5), []int{1, 2}, num(2), nil},
		{a.C(), a, []int{1}, pan, FoldMapError},
		{a.C(), a.C().Sum(0, 1, 2), []int{0, 1, 2}, sum, nil},
		{a.C().Flatten(), a, []int{1}, sum, IndexError},
		{a.C().Flatten(), a, []int{0, 1, 2}, sum, ShapeError},
		{a.C(), a.C().Sum(), []int{}, sum, nil},
	}

	for i, v := range tests {
		r := v.a.Fold(v.f, v.ax...)
		if !r.Equals(v.b).All().At(0) && !r.HasErr() {
			t.Logf("Test %d failed.  \nExpected:\n %v \nReceived:\n %v\n", i, v.b, r)
			t.Fail()
		}
		if e, d, s := r.GetDebug(); e != v.err {
			t.Logf("Test %d error failed.  Expected: %v Received: %v\n", i, v.err, e)
			t.Log(d, "\n", s, "\n", r)
			t.Fail()
		}
	}
}

func TestMap(t *testing.T) {
	num := func(i float64) MapFunc {
		return func(d float64) float64 {
			return i
		}
	}

	inc := func(i float64) MapFunc {
		return func(d float64) float64 {
			return d + i
		}
	}

	pan := func(d float64) float64 {
		var f *float64
		return *f
	}

	a := Arange(100).Reshape(2, 1, 5, 2, 5)

	tests := []struct {
		a, b *Array64
		f    MapFunc
		err  error
	}{
		{a.C(), a.C(), inc(0), nil},
		{a.C(), full(2, a.shape...), num(2), nil},
		{a.C(), a, pan, FoldMapError},
		{a.C(), a.C().AddC(5), inc(5), nil},
		{a.C().Reshape(100, 100), a, num(1), ReshapeError},
	}

	for i, v := range tests {
		r := v.a.Map(v.f)
		if !r.Equals(v.b).All().At(0) && !r.HasErr() {
			t.Logf("Test %d failed.  \nExpected:\n %v \nReceived:\n %v\n", i, v.b, r)
			t.Fail()
		}
		if e, d, s := r.GetDebug(); e != v.err {
			t.Logf("Test %d error failed.  Expected: %v Received: %v\n", i, v.err, e)
			t.Log(d, "\n", s, "\n", r)
			t.Fail()
		}
	}
}
