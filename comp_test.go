package numgo

import "testing"

func init() {
	debug = true
}

func TestMax(t *testing.T) {
	a := Arange(20).Reshape(2, 5, 2)

	tests := []struct {
		a, b *Array64
		ax   []int
		err  error
	}{
		{a.C(), NewArray64([]float64{19}), []int{}, nil},
		{a.C(), Arange(10, 20).Reshape(5, 2), []int{0}, nil},
		{a.C(), NewArray64([]float64{8, 9, 18, 19}, 2, 2), []int{1}, nil},
		{a.C(), Arange(0, 20, 2).AddC(1).Reshape(2, 5), []int{2}, nil},
		{a.C(), NewArray64([]float64{9, 19}), []int{1, 2}, nil},
		{nil, nil, []int{}, NilError},
		{a.C(), nil, []int{1, 2, 3, 4}, ShapeError},
		{a.C(), nil, []int{4}, IndexError},
	}
	for i, v := range tests {
		if c := v.a.Max(v.ax...).Equals(v.b); !c.All().At(0) && !c.HasErr() {
			t.Logf("Test %d Failed:\n %v == %v : %v\n", i, v.a.Max(v.ax...), v.b, c.All().At(0))
			t.Fail()
		}
		if e := v.a.GetErr(); e != v.err {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.err, e)
			t.Fail()
		}

	}
}

func TestMin(t *testing.T) {
	a := Arange(20).Reshape(2, 5, 2)

	tests := []struct {
		a, b *Array64
		ax   []int
		err  error
	}{
		{a.C(), NewArray64([]float64{0}), []int{}, nil},
		{a.C(), Arange(10).Reshape(5, 2), []int{0}, nil},
		{a.C(), NewArray64([]float64{0, 1, 10, 11}, 2, 2), []int{1}, nil},
		{a.C(), Arange(0, 20, 2).Reshape(2, 5), []int{2}, nil},
		{a.C(), NewArray64([]float64{0, 10}), []int{1, 2}, nil},
		{Arange(20, -1), NewArray64([]float64{-1}), []int{}, nil},
		{nil, nil, []int{}, NilError},
		{a.C(), nil, []int{1, 2, 3, 4}, ShapeError},
		{a.C(), nil, []int{4}, IndexError},
	}
	for i, v := range tests {
		if c := v.a.Min(v.ax...).Equals(v.b); !c.All().At(0) && !c.HasErr() {
			t.Logf("Test %d Failed:\n %v == %v : %v\n", i, v.a.Min(v.ax...), v.b, c.All().At(0))
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

func TestMinSet(t *testing.T) {
	a := Arange(20)

	tests := []struct {
		a []*Array64
		m *Array64
		e error
	}{
		{[]*Array64{a}, a, nil},
		{[]*Array64{a, a.C().AddC(1)}, a, nil},
		{[]*Array64{a.C().AddC(1), a}, a, nil},
		{[]*Array64{}, nil, NilError},
		{[]*Array64{a, nil}, nil, NilError},
		{[]*Array64{a, {nDimMetadata{err: InvIndexError}, nil}}, nil, InvIndexError},
		{[]*Array64{a, a.C().Reshape(2, 10)}, nil, ShapeError},
	}

	for i, v := range tests {
		m := MinSet(v.a...)
		if c := m.Equals(v.m); !c.All().At(0) && !c.HasErr() {
			t.Logf("Test %d Failed:\n %v\n", i, c)
			t.Fail()
		}
		if e := m.GetErr(); e != v.e {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.e, e)
			t.Fail()
		}
	}

}

func TestMaxSet(t *testing.T) {
	a := Arange(20)

	tests := []struct {
		a []*Array64
		m *Array64
		e error
	}{
		{[]*Array64{a}, a, nil},
		{[]*Array64{a, a.C().AddC(-1)}, a, nil},
		{[]*Array64{a.C().AddC(-1), a}, a, nil},
		{[]*Array64{}, nil, NilError},
		{[]*Array64{a, nil}, nil, NilError},
		{[]*Array64{a, {nDimMetadata{err: InvIndexError}, nil}}, nil, InvIndexError},
		{[]*Array64{a, a.C().Reshape(2, 10)}, nil, ShapeError},
	}

	for i, v := range tests {
		m := MaxSet(v.a...)
		if c := m.Equals(v.m); !c.All().At(0) && !c.HasErr() {
			t.Logf("Test %d Failed:\n %v\n", i, c)
			t.Fail()
		}
		if e := m.GetErr(); e != v.e {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.e, e)
			t.Fail()
		}
	}

}
