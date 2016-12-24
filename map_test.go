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
		{Arange(10).Reshape(2, 5).(*Array64), []int{1, -1}, []int{1, -1}, IndexError},
		{Arange(10).Reshape(2, 5).(*Array64), []int{1, 0}, []int{}, nil},
		{&Array64{nDimFields{err: InvIndexError}}, []int{}, []int{}, InvIndexError},
	} {
		if v.a.valAxis(&v.ax, "Test"); v.a.getErr() != v.err {
			t.Log("Error mismatch.", i, "Expected", v.err, "Got", v.a.getErr())
			t.Fail()
		}
		if len(v.ax) != len(v.re) {
			t.Log("Length incorrect.", i, "Expected", v.re, "Got", v.ax)
			t.Fail()
			continue tests
		}
		for idx, ax := range v.ax {
			if ax != v.re[idx] {
				t.Log("Result incorrect.", i, "Expected", v.re, "Got", v.ax)
				t.Fail()
				continue tests
			}
		}
	}
}

func TestFoldCC(t *testing.T) {

	num := func(i nDimElement) FoldFunc {
		return func(d []nDimElement) nDimElement {
			return i
		}
	}

	sum := func(d []nDimElement) nDimElement {
		r := 0.0
		for i := range d {
			r += d[i].(float64)
		}
		return nDimElement(r)
	}

	pan := func(d []nDimElement) nDimElement {
		r := 0.0
		//TODO check this loop, isn't it useless to assign to r?
		for i := range d {
			r += d[i].(float64)
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
		{a.C().(*Array64), a.C().Reshape(2, 5, 2, 5).(*Array64), []int{1}, sum, nil},
		{a.C().(*Array64), full(2, 2, 2, 5), []int{1, 2}, num(2), nil},
		{a.C().(*Array64), a.(*Array64), []int{1}, pan, FoldMapError},
		{a.C().(*Array64), a.C().(*Array64).Sum(0, 1, 2), []int{0, 1, 2}, sum, nil},
		{a.C().(*Array64).Flatten().(*Array64), a.(*Array64), []int{1}, sum, IndexError},
		{a.C().(*Array64).Flatten().(*Array64), a.(*Array64), []int{0, 1, 2}, sum, ShapeError},
	}

	for i, v := range tests {
		r := v.a.FoldCC(v.f, v.ax...)
		if !r.Equals(v.b).All().At(0).(bool) && !r.HasErr() {
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

	num := func(i nDimElement) FoldFunc {
		return func(d []nDimElement) nDimElement {
			return i
		}
	}

	sum := func(d []nDimElement) nDimElement {
		r := 0.0
		for i := range d {
			r += d[i].(float64)
		}
		return r
	}

	pan := func(d []nDimElement) nDimElement {
		r := 0.0
		for i := range d {
			r += d[i].(float64)
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
		{a.C().(*Array64), a.C().Reshape(2, 5, 2, 5).(*Array64), []int{1}, sum, nil},
		{a.C().(*Array64), full(2, 2, 2, 5), []int{1, 2}, num(2), nil},
		{a.C().(*Array64), a.(*Array64), []int{1}, pan, FoldMapError},
		{a.C().(*Array64), a.C().(*Array64).Sum(0, 1, 2), []int{0, 1, 2}, sum, nil},
		{a.C().Flatten().(*Array64), a.(*Array64), []int{1}, sum, IndexError},
		{a.C().Flatten().(*Array64), a.(*Array64), []int{0, 1, 2}, sum, ShapeError},
		{a.C().(*Array64), a.C().(*Array64).Sum(), []int{}, sum, nil},
	}

	for i, v := range tests {
		r := v.a.Fold(v.f, v.ax...)
		if !r.Equals(v.b).All().At(0).(bool) && !r.HasErr() {
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
	num := func(i nDimElement) MapFunc {
		return func(d nDimElement) nDimElement {
			return i
		}
	}

	inc := func(i nDimElement) MapFunc {
		return func(d nDimElement) nDimElement {
			return d.(float64) + i.(float64)
		}
	}

	pan := func(d nDimElement) nDimElement {
		var f *float64
		return *f
	}

	a := Arange(100).Reshape(2, 1, 5, 2, 5)

	tests := []struct {
		a, b *Array64
		f    MapFunc
		err  error
	}{
		{a.C().(*Array64), a.C().(*Array64), inc(0), nil},
		{a.C().(*Array64), full(2, a.fields().shape...), num(2), nil},
		{a.C().(*Array64), a.(*Array64), pan, FoldMapError},
		{a.C().(*Array64), a.C().(*Array64).AddC(5), inc(5), nil},
		{a.C().Reshape(100, 100).(*Array64), a.(*Array64), num(1), ReshapeError},
	}

	for i, v := range tests {
		r := v.a.Map(v.f)
		if !r.Equals(v.b).All().At(0).(bool) && !r.HasErr() {
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
