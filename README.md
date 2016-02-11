# numgo [![GoDoc](https://godoc.org/github.com/Kunde21/numgo?status.svg)](https://godoc.org/github.com/Kunde21/numgo) [![Build Status](https://travis-ci.org/Kunde21/numgo.svg?branch=master)](https://travis-ci.org/Kunde21/numgo) [![Go Report Card](https://goreportcard.com/badge/github.com/Kunde21/numgo)](https://goreportcard.com/report/github.com/Kunde21/numgo) [![codecov.io](https://codecov.io/github/Kunde21/numgo/coverage.svg?branch=master)](https://codecov.io/github/Kunde21/numgo?branch=master)

An n-dimensional array package implemented in Go.  

Note:  Under heavy development.  API will not stabilize until v0.1 tag.

## Installation 

```
go get github.com/Kunde21/numgo
```

## Using numgo

Most of the functionality resembles numpy's API, with some changes to work with Go's type system.  

```go
var array := numgo.NewArray64(nil,/*array shape*/1,2,3)	// This will be filled with zeros by default
var arange := numgo.Arange(100)                         // Simple 1-D array filled with incrementing numbers
arange.Reshape(2,5,10)                                  // Changes the shape from 1-D to 3-D
arange.Mean(2)                                          // Mean across axis 2, returning a 2-D (2x5) array
arange.Sum()                                            // An empty call operates over all data on all axes
```

Any formula can be created and mapped onto one or more axes within the array:

```go
	// Create a FoldFunc
	countNaNFn := func(input []float64) float64 {
	   var i float64 :=0
	       for _,v := range input {
	       	       if math.IsNaN(v) {
				i++
			}
		}
		return i
	}

	// Pass it into the Fold method and give it any number of axes to fold over

	// Returns an n-dimensional array object 
	// with the count of NaN values on 2nd and 4th axes. (0-based axis count)
	array.Fold(countNaNFn, 2,4) 
	// No axes operates over all data on all axes
	array.Fold(countNanfn)
```

## Function chaining

numgo is designed to allow chaining of functions, to allow different actions on different axes and at different points in the calculation.  Errors are maintained by the object and can be checked and handled using `HasErr()` or `GetErr()`:

```go
	// Errors are not handled on each call, 
	// but, instead, can be checked and handled after a block of calculations
	ng := numgo.Arange(100).Reshape(2,5,10).Mean(2).Min(1).Max()
	
	// Non-allocating style
	if ng.HasErr() {
	   log.Println(ng.GetErr())  // GetErr() clears the error flag
	 }
	   
	 // Allocation style
	if err = ng.GetErr(); err != nil {  
		log.Println(err)
	}
	// ng.GetErr() will always return nil here, 
	// so avoid stacking this type of error handling 
```

## Debugging option

Debugging can be enabled by calling `numgo.Debug(true)`.  This will give detailed error strings and stack traces by using `GetDebug()` instead of `GetErr()`.  This makes debugging chained method calls much easier.

```go
	numgo.Debug(true)
	nilp := new(Array64)		// Forgot to inintialize the array.
	
	nilp.Set(12, 1,4,0).AddC(2).DivC(6).At(1,4,0)
	if nilp.HasErr(){
		err, debug, trace := nilp.GetDebug()
		// Prints generic error: "Nil pointer received."
		fmt.Println(err)
		// Prints debug info: "Nil pointer received by GetDebug().  Source array was not initialized."
		fmt.Println(debug)
		// Prints stack trace for the call to GetDebug()
		fmt.Println(trace)
	}

	resz := NewArray64(nil, 2, 5)   // 2-D (2x5) array of zeros
	
	// Reshape would change the capacity of the array, which should use Resize
	resz.AddC(10).DivC(2).Reshape(3,3).Mean(1)  

	if resz.HasErr() {
	   	err, debug, trace := resz.GetDebug()
		// Prints generic error: "New shape cannot change the size of the array."
		fmt.Println(err)
		// Prints debug info: "Reshape() cannot change data size.  Dimensions: [2,5] reshape: [3,3]"
		fmt.Println(debug)
		// Prints stack trace for the call to Reshape()
		fmt.Println(trace)
	}

```

## Contributions

If you have any suggestions, corrections, bug reports, or design ideas please create an issue so that we can discuss and imporove the code.  
