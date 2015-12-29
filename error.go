package numgo

type ngError struct {
	s string
}

func (n *ngError) Error() string {
	return n.s
}

// Debug sets the error reporting level for the library.
// To get debugging data from the library, set this to true
// and use GetDebug() in place of GetErr().
//
// Debugging information includes the function call that generated
// the error and the values involved in that function call.
// This will add overhead to error reporting and handling, so
// use it for development and debugging purposes.
//
// Performance and memory will be better when using GetErr() and Debug = false.
var Debug = false

var (
	NilError     = &ngError{"Nil pointer recieved."}
	ShapeError   = &ngError{"Array shapes don't match and can't be broadcast."}
	ReshapeError = &ngError{"New shape cannot change the size of the array."}
	NegativeAxis = &ngError{"Negative axis length received."}
	IndexError   = &ngError{"Index or Axis out of range."}
	DivZeroError = &ngError{"Division by zero encountered."}
)

// HasErr tests for the existence of an error on the Arrayf object.
//
// Errors will be maintained through a chain of function calls,
// so only the first error will be returned when GetErr() is called.
// Use HasErr() as a gate for the GetErr() or GetDebug() choice in
// error handling code.
func (a *Arrayf) HasErr() bool {
	if a == nil {
		a = new(Arrayf)
		a.err = NilError
	}
	return a.err != nil
}

// GetErr returns the error object and clears the error from the array.
//
// This will only return an error value once per error instance.  Do not use
// it in the if statment to test for the existence of an error.  HasErr() is
// provided for that purpose.
func (a *Arrayf) GetErr() (err error) {
	if a == nil {
		a = new(Arrayf)
		a.err = NilError
	}
	err = a.err
	a.err, a.debug = nil, ""
	return
}

// GetDebug returns and clears the error object from the array object.  The returned debug string
// will include the function that generated the error and the arguments that caused it.
//
// This debug information will only be generated and returned if numgo.Debug is set to true
// before the function call that causes the error.
func (a *Arrayf) GetDebug() (err error, debug string) {
	if a == nil {
		a = new(Arrayf)
		a.err = NilError
	}
	err, debug = a.err, a.debug
	a.err, a.debug = nil, ""
	return
}

// encodeErr is a supporting function for MarshalJSON
func encodeErr(err *ngError) int8 {
	switch err {
	case nil:
		return 0
	case NilError:
		return 1
	case ShapeError:
		return 2
	case ReshapeError:
		return 3
	case NegativeAxis:
		return 4
	case IndexError:
		return 5
	case DivZeroError:
		return 6
	}
	return -1
}

// decodeErr is a supporting method for UnmarshalJSON
func (a *ngError) decodeErr(err int8) {
	switch err {
	case 0:
		a = nil
	case 1:
		a = NilError
	case 2:
		a = ShapeError
	case 3:
		a = ReshapeError
	case 4:
		a = NegativeAxis
	case 5:
		a = IndexError
	case 6:
		a = DivZeroError
	default:
		a = &ngError{"Unknown error Unmarshaled."}
	}
}

// HasErr tests for the existence of an error on the Arrayb object.
//
// Errors will be maintained through a chain of function calls,
// so only the first error will be returned when GetErr() is called.
// Use HasErr() as a gate for the GetErr() or GetDebug() choice in
// error handling code.
func (a *Arrayb) HasErr() bool {
	if a == nil {
		a = new(Arrayb)
		a.err = NilError
	}
	return a.err != nil
}

// GetErr returns the error object and clears the error from the array.
//
// This will only return an error value once per error instance.  Do not use
// it in the if statment to test for the existence of an error.  HasErr() is
// provided for that purpose.
func (a *Arrayb) GetErr() (err error) {
	if a == nil {
		a = new(Arrayb)
		a.err = NilError
	}
	err = a.err
	a.err, a.debug = nil, ""
	return
}

// GetDebug returns and clears the error object from the array object.  The returned debug string
// will include the function that generated the error and the arguments that caused it.
//
// This debug information will only be generated and returned if numgo.Debug is set to true
// before the function call that causes the error.
func (a *Arrayb) GetDebug() (err error, debug string) {
	if a == nil {
		a = new(Arrayb)
		a.err = NilError
	}
	err, debug = a.err, a.debug
	a.err, a.debug = nil, ""
	return
}
