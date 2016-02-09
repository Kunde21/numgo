package numgo

import (
	"math"
	"runtime"
	"testing"
)

func init() {
	debug = true
}

func TestAddC(t *testing.T) {
	a := Arange(21)

	if b := a.AddC(2).Equals(Arange(2, 23)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
}

func TestSubtrC(t *testing.T) {
	a := Arange(21)
	if b := a.SubtrC(2).Equals(Arange(-2, 19)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}

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
	if e := a.Equals(tst); !e.All().At(0) {
		t.Log(e.data)
		t.Log(Arange(1, 6))
		t.Log(Arange(20))
		t.Log(a.Flatten())
		t.Log(tst.Flatten())
		t.Fail()
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

func TestFMA12(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5 {
		t.Log("Expected:", (100-1)/.5, "Got:", len(a.data))
		t.FailNow()
	}

	b := a.C().FMA12(2, a)
	c := a.C().MultC(2).Add(a)
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}
}

func TestFMA21(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5 {
		t.Log("Expected:", (100-1)/.5, "Got:", len(a.data))
		t.FailNow()
	}

	b := a.C().FMA21(2, a)
	c := a.C().Mult(a).AddC(2)
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}
}

func BenchmarkSubtrC(b *testing.B) {
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.SubtrC(5)
	}
	runtime.GC()
}

func BenchmarkAddC_AVX(b *testing.B) {
	avxSupt = true
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.AddC(5)
	}
	runtime.GC()
}

func BenchmarkAddC_noAVX(b *testing.B) {
	avxSupt = false
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.AddC(5)
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

func BenchmarkFMA12_FMA(b *testing.B) {
	a := Arange(1, 1000000, .5)
	if len(a.data) != (1000000-1)/.5 {
		b.Log("Expected:", (1000000-1)/.5, "Got:", len(a.data))
		b.FailNow()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.FMA12(2, a)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkFMA12_noFMA(b *testing.B) {
	if !fmaSupt {
		b.Skip()
	}
	tmp := fmaSupt
	fmaSupt = false
	a := Arange(1, 1000000, .5)
	if len(a.data) != (1000000-1)/.5 {
		b.Log("Expected:", (1000000-1)/.5, "Got:", len(a.data))
		b.FailNow()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.FMA12(2, a)
	}
	b.StopTimer()
	fmaSupt = tmp
	runtime.GC()
}
