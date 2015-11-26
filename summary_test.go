package numgo

import (
	"fmt"
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Sum(0).data)
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	fmt.Println(a.shape)
	fmt.Println(a.strides)
	fmt.Println(a.Sum(0, 2, 4))
}

func TestCount(t *testing.T) {
	fmt.Println(Arange(10).Reshape(2, 5).Count(0).data)
	fmt.Println(Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6).Count(3, 4))
}

func TestMean(t *testing.T) {
	a := Arange(3*4*5*6*7).Reshape(7, 3, 4, 5, 6)
	fmt.Println(a.Count())

	fmt.Println(a.C().Sum(3, 4).Div(a.Count(3, 4)))

	fmt.Println(a.shape)
	fmt.Println(a.strides)
	fmt.Println(a.Count(0, 2, 4))
	fmt.Println(a.Mean(0, 2, 4))
}

func TestClean(t *testing.T) {
	fmt.Println(cleanAxis(0, 1, 2, 3, 4))
	fmt.Println(cleanAxis(0, 1, 2, 2, 4))
	fmt.Println(cleanAxis(0, 1, 2, 2, 4, 4))
	fmt.Println(12.0 + math.NaN())
}

func TestCollapse(t *testing.T) {
	a := Arange(2*3*4*5).Reshape(2, 3, 4, 5)
	sm := func(d []float64) (r float64) {
		for _, v := range d {
			r += v
		}
		return r
	}
	if !a.C().Sum(2, 1, 0).Equals(a.Map(sm, 2, 1, 0)).All().data[0] {
		t.Log("Failed on 3 Axis.")
		t.FailNow()
	} else {
		fmt.Println("Success on 3 Axis.")
	}
	if !a.C().Sum(0).Equals(a.Map(sm, 0)).All().data[0] {
		fmt.Println("Failed on 1 Axis.")
		t.Log("Failed on 1 Axis.")
		t.FailNow()
	} else {
		fmt.Println("Success on 1 Axis.")
	}

	if !a.C().Sum(2, 0).Equals(a.Map(sm, 2, 0)).All().data[0] {
		t.Log("Failed on 2 Axis.")
		t.FailNow()
	} else {
		fmt.Println("Success on 3 Axis.")
	}
	fmt.Println(2 * 3 * 4 * 5 * 6 * 7 * 8 * 9 * 10)
}

func BenchmarkCollapse(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.collapse(1, 3, 5, 7)
	}
}

func BenchmarkMapMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.C().Sum(1, 3, 5, 7).Equals(a.Map(sm, 1, 3, 5, 7)).All().data[0] {
		b.FailNow()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Map(sm, 1, 3, 5, 7)
	}
}

func BenchmarkMapCCSum(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10)
	sm := func(d []float64) (r float64) {
		i := 0
		for _, v := range d {
			r += v
			i++
		}
		return r / float64(i)
	}

	if !a.C().Mean(1, 3, 5, 7).Equals(a.Map(sm, 1, 3, 5, 7)).All().data[0] {
		b.FailNow()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.MapCC(sm, 1, 3, 5, 7)
	}
}

func BenchmarkMean(b *testing.B) {
	a := Arange(2*3*4*5*6*7*8*9*10).Reshape(2, 3, 4, 5, 6, 7, 8, 9, 10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mean(1, 3, 5, 7)
	}
}
