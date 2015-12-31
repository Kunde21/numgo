# numgo

An n-dimensional array package implemented in Go.  

## Installation 

```
go get github.com/Kunde21/numgo
```

## Using numgo

Most of the functionality resembles numpy's API, with some changes to work with Go's type system.  

```go
var array := numgo.Create(/*shape of the array*/1,2,3)	// This will be filled with zeros by default
var arange := numgo.Arange(100)                         // Simple 1-D array filled with incrementing numbers
arange.Reshape(2,5,10)                                  // Changes the shape from 1-D to 3-D
arange.Mean(2)                                          // Mean across axis 2, returning a 2-D (2x5) array
arange.Sum()                                            // An empty call operates over all data on all axes
```

Any formula can be created and mapped onto one or more axes within the array:

```go
	// Create a MapFunc
	countNaNFn := func(input []float64) float64 {
		var i float64 :=0
		for _,v := range input {
			if math.IsNaN(v) {
				i++
			}
		}
		return i
	}

	// Pass it into the Map method and give it any number of axes to map over

	 // Returns an n-dimensional array object 
	 // with the count of NaN values on 2nd and 4th axes. (0-based axis count)
	array.Map(countNaNFn, 2,4) 
	
	// No axes operates over all data on all axes
	array.Map(countNanfn)
```

## Function chaining

numgo is designed to allow chaining of functions, to allow different actions on different axes and at different points in the calculation.  Errors are maintained by the object and can be checked and handled using `Err()` or `GetErr()`:

```go
	// Errors are not handled on each call, 
	// but, instead, can be checked and handled after a block of calculations
	ng := numgo.Arange(100).Reshape(2,5,10).Mean(2).Min(1).Max()
	
	// Non-allocating style
	if ng.Err() {
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

Debugging can be enabled by calling `numgo.Debug(true)`.  This will give detailed error strings by using `GetDebug()` instead of `GetErr()`.  This makes debugging chained method calls much easier.

```go
	numgo.Debug(true)
	var nilp *Arrayf		// Forgot to create the array.
	
	nilp.SetE(12, 1,4,0).AddC(2).DivC(6).E(1,4,0)
	if nilp.HasErr(){
		err, debug := nilp.GetDebug()
		fmt.Println(err)		// Prints generic error: "Nil pointer received."
		fmt.Println(debug)		// Prints debug info: "Nil pointer received by SetE()."
	}
```

## Contributions

If you have any suggestions, corrections, bug reports, or design ideas please create an issue so that we can discuss and imporove the code.  
