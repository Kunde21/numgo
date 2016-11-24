package numgo

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	debug = true
	rand.Seed(time.Now().UnixNano())
}

func TestSum(t *testing.T) {
	t.Parallel()
	sz := []int{7, 3, 4, 5, 6}
	a := Arange(10).Reshape(2, 5).Sum(0)
	for i, v := range a.data {
		if int(v) != 5+2*i {
			t.Log("Incorrect result. Expected:", 5+2*i, "Received:", v)
		}
	}

	a = Arange(3*4*5*6*7).Reshape(sz...).Sum(0, 2, 4)
	for i, v := range a.data {
		if int(v) != (i%5+i/5*20)*1008+189420 {
			t.Log("Incorrect result.  Expected: ", (i%5+i/5*20)*1008+189420, "Received", v)
			t.Fail()
		}
	}

	if e := a.Reshape(0).Sum(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}
}

func BenchmarkSum(b *testing.B) {
	sz := []int{6, 3, 4, 8, 35}
	a := Arange(3 * 4 * 5 * 6 * 7 * 8).Reshape(sz...)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.C().Sum(0, 2, 4)
	}
	b.StopTimer()
	runtime.GC()
}

func TestNaNSum(t *testing.T) {
	t.Parallel()
	mask := func(i int, v float64) MapFunc {
		return func(d float64) float64 {
			if int(d)%i == 0 && d != 0 {
				return v
			}
			return d
		}
	}
	maskz := func(d float64) float64 {
		if math.IsNaN(d) {
			return 0
		}
		return d
	}

	a, nan := Arange(100).Reshape(2, 1, 2, 5, 5), math.NaN()
	r := rand.Intn
	for i := 1; i <= 20; i++ {
		v, x, y, z := r(100), r(4), r(4), r(4)
		if r := a.C().Map(mask(i, nan)).NaNSum(x, y, z).Map(maskz).Equals(a.C().Map(mask(i, 0)).Sum(x, y, z)); !r.All().At(0) {
			t.Logf("Test %d failed with v: %v x: %v y: %v z: %v\n", i, v, x, y, z)
			t.Log(r)
			t.Fail()
		}
	}

	if e := a.Reshape(0).NaNSum(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}

}

func TestCount(t *testing.T) {
	t.Parallel()
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

	if e := Arange(100).Reshape(0).Count(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}

	for i, v := range []struct {
		a  *Array64
		ax []int
	}{
		{NewArray64(nil, 1, 2, 3, 4, 5), []int{2, 3}},
		{NewArray64(nil, 2, 4, 6, 8), []int{0, 1, 2}},
		{NewArray64(nil, 2, 4, 8, 16), []int{}},
	} {
		if a, b := v.a.count(v.ax...), v.a.Count(v.ax...).data[0]; a != b {
			t.Log("Counts differed in test", i)
			t.Log(a, "!=", b)
			t.Fail()
		}
	}
}

func TestNaNCount(t *testing.T) {
	t.Parallel()
	mask := func(i int, v float64) MapFunc {
		return func(d float64) float64 {
			if int(d)%i == 0 && d != 0 {
				return v
			}
			return d
		}
	}

	a, nan := Arange(100).Reshape(2, 1, 2, 5, 5), math.NaN()
	r := rand.Intn
	for i := 1; i <= 20; i++ {
		v, x, y, z := r(99)+1, r(4), r(4), r(4)
		r := a.C().Map(mask(v, nan)).NaNCount(x, y, z).NaNSum()
		s := a.C().Map(mask(v, 0)).Count(x, y, z).Sum().SubtrC(float64(100 / v))
		if 100%v == 0 {
			r.SubtrC(1)
		}
		if !r.Equals(s).At(0) {
			t.Logf("Test %d failed with v: %v x: %v y: %v z: %v\n", i, v, x, y, z)
			t.Log(float64(100 / v))
			t.Log(a.C().Map(mask(v, nan)).NaNCount(x, y, z), "\n", r, s)
			t.Fail()
		}
	}

	if e := a.Reshape(0).NaNCount(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}
}

func TestMean(t *testing.T) {
	t.Parallel()
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

	for i, v := range a.Mean(0, 2, 4).data {
		if v != 1127.5+float64(i%5*6+i/5*120) {
			t.Log("Incorrect result.  Expected", 1127.5+float64(i%5*6+i/5*120), "Received ", v)
			t.Fail()
		}
	}
	if e := a.Reshape(0).Mean(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}

}

func TestNanMean(t *testing.T) {
	t.Parallel()
	a := Arange(2*3*4*5).Reshape(2, 3, 4, 5)
	if !a.Mean(1, 3).Equals(a.NaNMean(1, 3)).All().data[0] {
		t.Log("NaNMean producing different results than Mean")
		t.Log(a.Mean(1, 3))
		t.Log(a.NaNMean(1, 3))
		t.Fail()
	}
	if e := a.Reshape(0).NaNMean(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}
}

func TestNonzero(t *testing.T) {
	t.Parallel()
	maskz := func(d float64) float64 {
		if d == 0 {
			return math.NaN()
		}
		return d
	}

	sz, r := []int{4, 7, 3, 8, 2}, rand.Intn
	for i := 0; i < 10; i++ {
		x, y, z := r(len(sz)-1), r(len(sz)-1), r(len(sz)-1)
		a := RandArray64(0, 100, sz...).Append(NewArray64(nil, sz...), y)
		if j, k := a.C().Nonzero(x, y, z), a.C().Map(maskz).NaNCount(x, y, z); !j.Equals(k).All().At(0) {
			t.Logf("Test %d failed with x: %v y: %v z: %v\n", i, x, y, z)
			t.Log("Expected:\n", k)
			t.Log("Received:\n", j)
			t.Log(j.GetDebug())
			t.Log(a.count(x, y, z))
			t.FailNow()
		}
		if a.C().Nonzero().Equals(a.C().Count()).At(0) {
			i--
		}
	}
	if e := Arange(10).Reshape(0).Nonzero(1).GetErr(); e != ReshapeError {
		t.Log("Error failed", e)
		t.Fail()
	}

}

func TestFoldMean(t *testing.T) {
	t.Parallel()
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.Mean(1, 3, 5, 7).Equals(a.Fold(sm, 1, 3, 5, 7)).All().At(0) {
		t.Fail()
	}
}

func TestCollapse(t *testing.T) {
	t.Parallel()
	a := Arange(2*3*4*5).Reshape(2, 3, 4, 5)
	sm := func(d []float64) (r float64) {
		for _, v := range d {
			r += v
		}
		return r
	}
	if !a.C().Sum(2, 1, 0).Equals(a.Fold(sm, 2, 1, 0)).All().At(0) {
		t.Log("Failed on 3 Axis.")
		t.FailNow()
	}

	if !a.C().Sum(0).Equals(a.Fold(sm, 0)).All().At(0) {
		fmt.Println("Failed on 1 Axis.")
		t.Log("Failed on 1 Axis.")
		t.FailNow()
	}

	if !a.C().Sum(2, 0).Equals(a.Fold(sm, 2, 0)).All().At(0) {
		t.Log("Failed on 2 Axis.")
		t.FailNow()
	}
}

func BenchmarkCollapse(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.collapse([]int{5, 3, 7, 8})
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkFoldMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sm := func(d []float64) (r float64) {
		l := len(d)
		for i := 0; i < l; i++ {
			r += d[i]
		}
		return r / float64(l)
	}

	if !a.C().Mean(1, 3, 5, 7).Equals(a.Fold(sm, 1, 3, 5, 7)).All().data[0] {
		b.FailNow()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Fold(sm, 5, 3, 7, 8)
	}
	b.StopTimer()
	runtime.GC()
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
	b.StopTimer()
	runtime.GC()
}

func BenchmarkMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mean(5, 3, 7, 8)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkC(t *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10*11).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		_ = a.C()
	}
	t.StopTimer()
	runtime.GC()
}
