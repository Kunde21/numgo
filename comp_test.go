package numgo

import "testing"

func TestMax(t *testing.T) {
	a := Arange(20).Reshape(2, 5, 2)

	tests := []struct {
		a, b *Array64
		ax   []int
	}{
		{a.C(), NewArray64([]float64{19}), []int{}},
		{a.C(), Arange(10, 20).Reshape(5, 2), []int{0}},
		{a.C(), NewArray64([]float64{8, 9, 18, 19}, 2, 2), []int{1}},
		{a.C(), Arange(0, 20, 2).AddC(1).Reshape(2, 5), []int{2}},
		{a.C(), NewArray64([]float64{9, 19}), []int{1, 2}},
	}
	for i, v := range tests {
		if c := v.a.Max(v.ax...).Equals(v.b); !c.All().At(0) {
			t.Logf("Test %d Failed:\n %v == %v : %v\n", i, v.a.Max(v.ax...), v.b, c.All().At(0))
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
		{a.C(), Arange(0, 10).Reshape(5, 2), []int{0}, nil},
		{a.C(), NewArray64([]float64{0, 1, 10, 11}, 2, 2), []int{1}, nil},
		{a.C(), Arange(0, 20, 2).Reshape(2, 5), []int{2}, nil},
		{a.C(), NewArray64([]float64{0, 10}), []int{1, 2}, nil},
		{nil, nil, []int{}, NilError},
	}
	for i, v := range tests {
		if c := v.a.Min(v.ax...).Equals(v.b); !c.All().At(0) && !c.HasErr() {
			t.Logf("Test %d Failed:\n %v == %v : %v\n", i, v.a.Min(v.ax...), v.b, c.All().At(0))
			t.Fail()
		}
		if e := v.a.GetErr(); e != v.err {
			t.Logf("Test %d Error Failed: Expected %#v got %#v\n", i, v.err, e)
			t.Fail()
		}
	}
}
