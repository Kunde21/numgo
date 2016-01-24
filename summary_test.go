package numgo

import (
	"fmt"
	"testing"
)

func init() {
	debug = true
}

func TestSum(t *testing.T) {
	sz := []int{7, 3, 4, 5, 6}
	a := Arange(10).Reshape(2, 5).Sum(0)
	for i, v := range a.data {
		if int(v) != 5+2*i {
			t.Log("Incorrect result. Expected:", 5+2*i, "Received:", v)
		}
	}

	a = Arange(3 * 4 * 5 * 6 * 7).Reshape(sz...)

	str := uint64(1)
	for i, v := range a.shape {
		if int(v) != sz[i] {
			t.Log("Shape incorrect.  Expected:", sz, "Received:", a.shape)
			t.Fail()
		}
		str *= v
	}
	for i, v := range a.strides {
		if v != str {
			t.Log("Stride incorrect at ", i, "Expected:", str, "Received:", v)
			t.Fail()
		}
		if i < len(a.shape) {
			str /= a.shape[i]
		}
	}

	a = a.Sum(0, 2, 4)
	for i, v := range a.data {
		if int(v) != (i%5+i/5*20)*1008+189420 {
			t.Log("Incorrect result.  Expected: ", (i%5+i/5*20)*1008+189420, "Received", v)
			t.Fail()
		}
	}
}

func TestCount(t *testing.T) {
	for _, v := range Arange(10).Reshape(2, 5).Count(0).data {
		if v != 2 {
			t.Log("Incorrect count. Expected: 2 Received:", v)
			t.Fail()
		}
	}
	for _, v := range Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6).Count(3, 4).data {
		if v != 30 {
			t.Log("Incorrect count. Expected: 30 Received:", v)
			t.Fail()
		}
	}

}

func TestMean(t *testing.T) {
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	if tmp := a.Count().At(0); tmp != 3*4*5*6*7 {
		t.Log("Incorrect size.  Expected ", 3*4*5*6*7, "Received ", tmp)
		t.Fail()
	}

	for i, v := range a.C().Sum(3, 4).Div(a.Count(3, 4)).data {
		if v != 14.5+float64(i*30) {
			t.Log("Incorrect result.  Expected", 14.5+float64(i*30), "Received ", v)
			t.Fail()
		}
	}

	for _, v := range a.Count(0, 2, 4).data {
		if v != 168 {
			t.Log("Incorrect Count. Expected 168 Received", v)
			t.Fail()
		}
	}

	for i, v := range a.Mean(0, 2, 4).data {
		if v != 1127.5+float64(i%5*6+i/5*120) {
			t.Log("Incorrect result.  Expected", 1127.5+float64(i%5*6+i/5*120), "Received ", v)
			t.Fail()
		}
	}
}

func TestFoldMean(t *testing.T) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.Mean(1, 3, 5, 7).Equals(a.Fold(sm, 1, 3, 5, 7)).All().data[0] {
		t.Fail()
	}
}

func TestCollapse(t *testing.T) {
	a := Arange(2*3*4*5).Reshape(2, 3, 4, 5)
	sm := func(d []float64) (r float64) {
		for _, v := range d {
			r += v
		}
		return r
	}
	if !a.C().Sum(2, 1, 0).Equals(a.Fold(sm, 2, 1, 0)).All().data[0] {
		t.Log("Failed on 3 Axis.")
		t.FailNow()
	}

	if !a.C().Sum(0).Equals(a.Fold(sm, 0)).All().data[0] {
		fmt.Println("Failed on 1 Axis.")
		t.Log("Failed on 1 Axis.")
		t.FailNow()
	}

	if !a.C().Sum(2, 0).Equals(a.Fold(sm, 2, 0)).All().data[0] {
		t.Log("Failed on 2 Axis.")
		t.FailNow()
	}
}

func TestNanMean(t *testing.T) {
	a := Arange(2*3*4*5).Reshape(2, 3, 4, 5)
	if !a.Mean(1, 3).Equals(a.NaNMean(1, 3)).All().data[0] {
		t.Log("NaNMean producing different results than Mean")
		t.Log(a.Mean(1, 3))
		t.Log(a.NaNMean(1, 3))
		t.Fail()
	}
}

func BenchmarkCollapse(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.collapse([]int{5, 3, 7, 8})
	}
}

func BenchmarkFoldMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.C().Mean(1, 3, 5, 7).Equals(a.Fold(sm, 1, 3, 5, 7)).All().data[0] {
		b.FailNow()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Fold(sm, 5, 3, 7, 8)
	}
}

func BenchmarkFoldCCMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.C().Mean(1, 3, 5, 7).Equals(a.FoldCC(sm, 1, 3, 5, 7)).All().data[0] {
		b.FailNow()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.FoldCC(sm, 5, 3, 7, 8)
	}
}

func BenchmarkMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mean(5, 3, 7, 8)
	}
}
