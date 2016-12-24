package numgo

import "testing"

func init() {
	Debug(true)
}

func TestDebug(t *testing.T) {
	Debug(true)
	var nilp *Array64
	nilp.Set(12, 1, 4, 0)
	nilp.AddC(2)
	nilp.DivC(6)
	nilp.At(1, 4, 0)
	if !nilp.HasErr() {
		t.FailNow()
		err, debug, stack := nilp.GetDebug()
		t.Log(err)   // Prints generic error: "Nil pointer received."
		t.Log(debug) // Prints debug info: "Nil pointer received by SetE()."
		t.Log(stack)
		t.Fail()
	}
	nilp = MinSet(Arange(10).Reshape(2, 5).(*Array64), Arange(10))
	if err, debug, stack := nilp.GetDebug(); err != ShapeError {
		t.Log(err)
		t.Log(debug)
		t.Log(stack)
		t.Fail()
	}

}

func TestEncodeDecode(t *testing.T) {
	for _, v := range []error{
		nil,
		NilError,
		ShapeError,
		ReshapeError,
		NegativeAxis,
		IndexError,
		InvIndexError,
		InvIndexError,
		FoldMapError,
	} {
		if e := decodeErr(encodeErr(v)); v != e {
			t.Log("Failed:", v)
			t.Log("Received:", e)
			t.Fail()
		}
	}

	if e := encodeErr(&ngError{"Not pkg error"}); e != -1 {
		t.Log("Incorrect error returned: ", e)
		t.Fail()
	}

	if e := decodeErr(-1); e.Error() != "Unknown error Unmarshaled: -1" {
		t.Log("Incorrect error returned: ", e)
		t.Fail()
	}
}

func TestGetErr(t *testing.T) {
	var a *Array64
	var b *Arrayb

	c, d := new(Array64), new(Arrayb)

	if !a.HasErr() || !b.HasErr() {
		t.Log("Bad error state", a.HasErr(), b.HasErr())
		t.Fail()
	}
	if !c.HasErr() || !d.HasErr() {
		t.Log("Bad error state (new) ", c.HasErr(), d.HasErr())
		t.Fail()
	}
	switch {
	case a.GetErr() != NilError:
		t.Log("a failed", a.GetErr())
		t.Fail()
	case a.getErr() != NilError:
		t.Log("a failed", a.GetErr())
		t.Fail()
	case b.GetErr() != NilError:
		t.Log("b failed", b.GetErr())
		t.Fail()
	case c.GetErr() != NilError:
		t.Log("c failed", c.GetErr())
		t.Fail()
	case d.GetErr() != NilError:
		t.Log("d failed", d.GetErr())
		t.Fail()
	}

	d.err = InvIndexError
	if e := d.GetErr(); e != InvIndexError {
		t.Log("Error storage failed", e)
		t.Fail()
	}
	c.err = InvIndexError
	if e := c.GetErr(); e != InvIndexError {
		t.Log("Error storage failed", e)
		t.Fail()
	}
}
