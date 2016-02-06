package numgo

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"testing"
)

func init() {
	debug = true
}

func TestAddC(t *testing.T) {
	a := Arange(21)

	fmt.Println(a.C().AddC(5))
	a.Resize(20)
	fmt.Println(a.C().AddC(5))
	a = Arange(22)
	fmt.Println(a.C().AddC(5))
	a.Resize(23)
	fmt.Println(a.C().AddC(5))

	runtime.GC()
}

func TestSubtrC(t *testing.T) {
	a := Arange(21)

	fmt.Println(a.C().SubtrC(5))
	a.Resize(20)
	fmt.Println(a.C().SubtrC(5))
	a = Arange(22)
	fmt.Println(a.C().SubtrC(5))
	a.Resize(23)
	fmt.Println(a.C().SubtrC(5))
	runtime.GC()
}

func TestAdd(t *testing.T) {
	a, b := Arange(20), Arange(20)
	a.Add(b)
	runtime.GC()

	if c := a.Add(b.Reshape(2, 10)); c.err != ShapeError {
		t.Log("Shape tests failed.  Expected nil, returned:", c)
		t.Fail()
	}
	runtime.GC()
	if c := Arange(20).Reshape(2, 10).Add(b.Reshape(2, 10)); c.Equals(Arange(0, 40, 2).Reshape(2, 10)).All().At(0) == false {
		t.Log("Addition on multiple axis failed.  Expected 0-40 with increment by 2.")
		t.Log("Recieved:", c)
		t.Fail()
	}
	runtime.GC()

	log.Println("Mismatch")
	a = Arange(20)
	a.Reshape(2, 2, 5)
	runtime.GC()
	a.Add(Arange(5))
	runtime.GC()
	for i, v := range a.data {
		if int(v) != i%5*2+i/5*5 {
			t.Log("Incorrect result at:", i, "Expected", i%5*2+i/5*5, "Got", v)
			t.Fail()
		}
	}
	runtime.GC()
	a.AddC(1)
	tst := Arange(20).Reshape(2, 2, 5).Add(Arange(1, 6))
	runtime.GC()
	if e := a.Equals(tst); e.All().At(0) {
		fmt.Println("Add passed")
	} else {
		fmt.Println(e.data)
		fmt.Println(Arange(1, 6))
		fmt.Println(Arange(20))
		fmt.Println(a.Flatten())
		fmt.Println(tst.Flatten())
	}
	runtime.GC()
}

func TestSubtr(t *testing.T) {
	a := Arange(20).Reshape(2, 10).Subtr(Arange(10))
	runtime.GC()
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i)
		for _, v := range tmp {
			if v != float64(i*10) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("At", i, "Expected: ", float64(i*10), "Got:", v)
				t.Fail()
			}
		}
	}

	a = Arange(20).Reshape(2, 2, 5).Subtr(Arange(5))
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i)
		for _, v := range tmp {
			if v != float64(i/5*5) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("Expected: ", i*5, "Got:", v)
				t.Fail()
			}
		}
	}

	a = Arange(20).Reshape(2, 1, 1, 2, 5).Subtr(Arange(5))
	for i := 0; i < 2; i++ {
		tmp := a.SliceElement(i, 0, 0)
		for _, v := range tmp {
			if v != float64(i/5*5) {
				t.Log("Subtract resulted in unexpected value.")
				t.Log("Expected: ", i*5, "Got:", v)
				t.Fail()
			}
		}
	}
}

func TestMult(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5 {
		t.Log("Expected:", (100-1)/.5, "Got:", len(a.data))
		t.FailNow()
	}

	a = Full(math.NaN(), 4, 4)
	for _, v := range a.data {
		if !math.IsNaN(v) {
			t.Log("Expected NaN, got ", v)
			t.FailNow()
		}
	}
}

func BenchmarkAddC(b *testing.B) {
	a := Arange(1003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.AddC(5)
	}
	runtime.GC()
}

func BenchmarkSubtrC(b *testing.B) {
	a := Arange(1003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.SubtrC(5)
	}
	runtime.GC()
}

func BenchmarkMultC(b *testing.B) {
	a := Arange(1003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.MultC(2)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkDivC(b *testing.B) {
	a := Arange(1003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.SubtrC(5)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkAdd(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 12), Arange(12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Add(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkSubtr(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 12), Arange(12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Subtr(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkMult(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 12), Arange(1, 13)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Mult(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkDiv(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 12), Arange(1, 13)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Div(n)
	}
	b.StopTimer()
	runtime.GC()
}
