/*
Package numgo provides implementations of n-dimensional array objects and operations.

Two types of numgo arrays are currently supported:
Array64 holds float64 values
Arrayb holds boolean values


Basic usage

Array64 objects can be created from a slice or as an empty array.  Arange and Identity are convenience creation functions.

 array := numgo.NewArray64(nil,1,2,3)  // This 1x2x3 array will be filled with zeros by default

 sl := []float64{2,2,3,3,4,4}
 wrap := numgo.NewArray64(sl)          // Wrap a copy of the slice as a 1-D Array64
 wrap2 := numgo.NewArray64(sl,3,2)     // Wrap a copy of the slice as a 3x2 Array64
 wrap3 := numgo.NewArray64(sl,2,2)     // Only wraps the first 4 values in sl as a 2x2 Array64

 arange := numgo.Arange(100)           // Simple 1-D array filled with incrementing numbers
 arange.Reshape(2,5,10)                // Changes the shape from 1-D to 3-D
 arange.Mean(2)                        // Mean across axis 2, returning a 2-D (2x5) array
 arange.Sum()                          // An empty call operates on all data (Grand Total)


Fold and Map operations

Map takes a function of type MapFunc and applies it across all data elements.  Fold and FoldCC take a function of type FoldFunc and applies it in contracting the data across one or more axes.

 // Increment MapFunc definition
 increment := func(input float64) float64 {
        return input + 1
 }
 array.Map(increment)  // Using map with defined function
 array.AddC(1)   // Using built-in add method


 // Create a FoldFunc
    sumFn := func(input []float64) float64 {
        var t float64 :=0
        for _,v := range input {
            t += v
        }
        return t
    }

 // Pass it into the Map method and give it any number of axes to map over

 // Returns an n-dimensional array object
 // with the count of NaN values on 2nd and 4th axes. (0-based axis count)
 array.Fold(sumFn, 2,4)

 // No axes operates over all data on all axes
 array.Fold(sumfn)

Function Chaining

numgo is designed with function-chaining at its core, to allow different actions on different axes and at different points in the calculation. Errors are maintained by the object and can be checked and handled using HasErr() and GetErr():

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

Debugging options

Debugging can be enabled by calling numgo.Debug(true). This will give detailed error strings by using GetDebug() instead of GetErr(). This makes debugging chained method calls much easier.

 numgo.Debug(true)
 nilp := new(Array64)     // Forgot to inintialize the array.

 nilp.Set(12, 1,4,0).AddC(2).DivC(6).At(1,4,0)
 if nilp.HasErr(){
     err, debug, trace := nilp.GetDebug()

     // Prints generic error: "Nil pointer received."
     fmt.Println(err)

     // Prints debug info:
     // "Nil pointer received by GetDebug().  Source array was not initialized."
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
	// Prints debug info:
        // "Reshape() cannot change data size.  Dimensions: [2,5] reshape: [3,3]"
	fmt.Println(debug)
	// Prints stack trace for the call to Reshape()
	fmt.Println(trace)
 }
*/
package numgo
