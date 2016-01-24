package numgo

import "testing"

func init() {
	Debug(true)
}

func TestDebug(t *testing.T) {
	Debug(true)
	var nilp *Array64
	nilp.Set(12, 1, 4, 0).AddC(2).DivC(6).At(1, 4, 0)
	if !nilp.HasErr() {
		t.FailNow()
		err, debug, stack := nilp.GetDebug()
		t.Log(err)   // Prints generic error: "Nil pointer received."
		t.Log(debug) // Prints debug info: "Nil pointer received by SetE()."
		t.Log(stack)
		t.Fail()
	}
	nilp = MinSet(Arange(10).Reshape(2, 5), Arange(10))
	if err, debug, stack := nilp.GetDebug(); err != ShapeError {
		t.Log(err)
		t.Log(debug)
		t.Log(stack)
		t.Fail()
	}

}
