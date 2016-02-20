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

	if b := a.AddC(2).Equals(Arange(2, 22)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
	a.Reshape(10, 10).AddC(1)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestSubtrC(t *testing.T) {
	a := Arange(21)
	if b := a.SubtrC(2).Equals(Arange(-2, 18)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
	a.Reshape(10, 10).SubtrC(2)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestMultC(t *testing.T) {
	a := Arange(21)

	if b := a.MultC(2).Equals(Arange(0, 40, 2)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
	a.Reshape(10, 10).MultC(1)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestDivC(t *testing.T) {
	a := Arange(0, 40, 2)
	if b := a.DivC(2).Equals(Arange(21)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
	a.DivC(0)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestPowC(t *testing.T) {
	a := Arange(0, 40, 2)
	if b := a.PowC(0).Equals(Full(1, 21)); !b.All().At(0) {
		t.Log(b)
		t.Log(a)
		t.Fail()
	}
	a.Reshape(10, 10).PowC(0)
	if !a.HasErr() {
		t.Log(a.GetErr())
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
	if c := Arange(20).Reshape(2, 10).Add(b.Reshape(2, 10)); !c.Equals(Arange(0, 38, 2).Reshape(2, 10)).All().At(0) {
		t.Log("Addition on multiple axis failed.  Expected 0-40 with increment by 2.")
		t.Log("Recieved:", c)
		t.Log(Arange(0, 38, 2).Reshape(2, 10))
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
	tst := Arange(20).Reshape(2, 2, 5).Add(Arange(1, 5))
	runtime.GC()
	if e := a.Equals(tst); !e.All().At(0) {
		t.Log(e.data)
		t.Log(Arange(1, 5))
		t.Log(Arange(20))
		t.Log(a.Flatten())
		t.Log(tst.Flatten())
		t.Fail()
	}

	a.Reshape(10, 10).Add(Arange(5))
	if !a.HasErr() {
		t.Log(a.GetErr())
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

	a.Reshape(10, 10).Subtr(Arange(5))
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}

}

func TestMult(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5+1 {
		t.Log("Expected:", (100-1)/.5+1, "Got:", len(a.data))
		t.FailNow()
	}

	a = Full(math.NaN(), 4, 4)
	for _, v := range a.data {
		if !math.IsNaN(v) {
			t.Log("Expected NaN, got ", v)
			t.FailNow()
		}
	}
	a.Reshape(10, 10).Mult(Arange(5))
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestDiv(t *testing.T) {
	a := Arange(1, 100, .5)
	if len(a.data) != (100-1)/.5+1 {
		t.Log("Expected:", (100-1)/.5+1, "Got:", len(a.data))
		t.FailNow()
	}

	a = a.C().Div(a)
	if e := a.Equals(Full(1, 200)); e.All().At(0) {
		t.Log(e)
		t.Fail()
	}

	a.Reshape(10, 10).Div(Arange(5))
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestPow(t *testing.T) {
	a := Arange(1, 100.5, .5)
	if len(a.data) != (100.5-1)/.5+1 {
		t.Log("Expected:", (100.5-1)/.5+1, "Got:", len(a.data))
		t.FailNow()
	}

	a = a.Reshape(2, 25, 4).Pow(Full(0, 4))
	if e := a.Equals(Full(1, 2, 25, 4)); !e.All().At(0) {
		t.Log(e)
		t.Fail()
	}
	if a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}

	a.Reshape(10, 10).Pow(Arange(5))
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestFMA12(t *testing.T) {
	a := Arange(1, 100.5, .5)
	if len(a.data) != (100.5-1)/.5+1 {
		t.Log("Expected:", (100.5-1)/.5+1, "Got:", len(a.data))
		t.FailNow()
	}

	b := a.C().FMA12(2, a)
	c := a.C().MultC(2).Add(a)
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}

	b = a.C().Reshape(2, 4, 25).FMA12(2, Arange(25))
	c = a.C().Reshape(2, 4, 25).MultC(2).Add(Arange(25))
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}

	a.Reshape(10, 10).FMA12(2, a)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestFMA21(t *testing.T) {
	a := Arange(1, 100.5, .5)
	if len(a.data) != (100.5-1)/.5+1 {
		t.Log("Expected:", (100.5-1)/.5+1, "Got:", len(a.data))
		t.FailNow()
	}

	b := a.C().FMA21(2, a)
	c := a.C().Mult(a).AddC(2)
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}

	b = a.C().Reshape(2, 4, 25).FMA21(2, Arange(25))
	c = a.C().Reshape(2, 4, 25).Mult(Arange(25)).AddC(2)
	if e := b.Equals(c); !e.All().At(0) {
		t.Log(a, "\n", b, "\n", c)
		t.Fail()
	}

	a.Reshape(10, 10).FMA21(2, a)
	if !a.HasErr() {
		t.Log(a.GetErr())
		t.Fail()
	}
}

func TestValRith(t *testing.T) {
	var a, b *Array64
	if !a.valRith(b, "") {
		t.Fail()
	}

	a = NewArray64(nil, 0)
	if !a.valRith(b, "") {
		t.Fail()
	}

	_ = a.GetErr()
	b = &Array64{err: DivZeroError}
	if !a.valRith(b, "") {
		t.Fail()
	}

	_ = a.GetErr()
	b = NewArray64(nil, 2, 2, 2)
	if !a.valRith(b, "") {
		t.Fail()
	}

	_ = a.GetErr()
	a = NewArray64(nil, 2, 2, 4)
	if !a.valRith(b, "") {
		t.Fail()
	}

	_ = a.GetErr()
	a.Resize(2, 2, 2)
	if a.valRith(b, "") {
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
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.AddC(5)
	}
	runtime.GC()
}

func BenchmarkAddC_noAVX(b *testing.B) {
	if !avxSupt {
		b.Skip()
	}
	avxSupt = false
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.AddC(5)
	}
	b.StopTimer()
	avxSupt = true
	runtime.GC()
}

func BenchmarkMultC(b *testing.B) {
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.MultC(2)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkDivC(b *testing.B) {
	a := Arange(5003)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.SubtrC(5)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkAdd(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 10, 12), Arange(12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Add(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkSubtr(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 10, 12), Arange(12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Subtr(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkMult(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 10, 12), Arange(1, 12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Mult(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkDiv(b *testing.B) {
	a, n := NewArray64(nil, 10, 10, 10, 12), Arange(1, 12)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Div(n)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkFMA12_FMA(b *testing.B) {
	a := Arange(0, 1000000, .5)
	if len(a.data) != (1000000)/.5+1 {
		b.Log("Expected:", (1000000)/.5+1, "Got:", len(a.data))
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

func BenchmarkFMA12_FMAB(b *testing.B) {
	a := Arange(0, 1000000, .5).Resize(2, 2, 500, 1000)
	if len(a.data) != (1000000)/.5 {
		b.Log("Expected:", (1000000)/.5, "Got:", len(a.data))
		b.FailNow()
	}
	c := Arange(1000)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.FMA12(2, c)
	}
	b.StopTimer()
	runtime.GC()
}

func BenchmarkCopy(b *testing.B) {
	a := Arange(0, 1000000, .5)
	if len(a.data) != (1000000)/.5+1 {
		b.Log("Expected:", (1000000)/.5+1, "Got:", len(a.data))
		b.FailNow()
	}
	c := a.C()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		copy(c.data, a.data)
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
	if len(a.data) != (1000000-1)/.5+1 {
		b.Log("Expected:", (1000000-1)/.5+1, "Got:", len(a.data))
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
