One of the main tools for programming is to use variables. Variables store values
that may be manipulated further down the line. What seems like a simple
mechanism is actually a powerful tool. This two part series covers the main
ways to initialize variables in Go. In Go there is a lot of flexibility when it
comes to declaring variables. Sometimes there is confusion which syntax should
be used. This post aims to address that.

To declare and initialize a variable with a simple type like an integer or a
string, you can use the following syntax:

```go
// Option 1
var x int     // declare an integer, initialized with value 0
var b bool    // declare a bool, initialized with value false
var s string  // declare a string, initialized with value ""

// Option 2
var x int = 0       // declare and initialize an integer with value 0
var b bool = false  // declare and initialize a bool with value false
var s string = ""   // declare and initialize a string with value ""

// Option 3
var x = 0      // declare and initialize an integer with value 0
var b = false  // declare and initialize a bool with value false
var s = ""     // declare and initialize a string with value ""

// Option 4
x := 0    // declare and initialize an integer with value 0
b := 0.0  // declare and initialize a float with value false
s := ""   // declare and initialize a string with value ""
```

Keen readers might notice that the last two options seem very similar, but one
uses the equals operator `=` instead of the walrus`:=` operator. The short
variable declaration is considered to be idiomatic and it is the preferred
method to initialize a variable with a specific value. In this case, the type
can be safely inferred since the value needs to necessarily be set.

At this point, why have option 3 at all? Well, the `var` keyword enables
concise multiple variable declaration:

```go
// initialize multiple variables of different types
var x, b, s = 0, false, ""
```

Sometimes there is a need to keep track of multiple related values. Go offers
arrays and array slices, as well as hash based maps to tackle this need.

To declare an array the following syntax is used, assume `MAX` is an integer
constant previously initialized with value 3:
```go
// declare an array of type int with fixed size MAX
// where all elements have value 0
var coords [MAX]int

 // declare an array of type int with fixed size MAX
 // where all elements have value 0 using an array literal
var coords [MAX]int = [MAX]int{0,0,0}

// declare an array of type int with fixed size MAX
// where all elements have value 0
// using an array literal and short variable declaration
coords := [MAX]int{0,0,0}
```

But how do you declare an array of arrays (a matrix)? Well you use `[SIZE]`
notation followed by the type, in this case we want an array of type array. The
inner array can be of type integer, for example:

```go
// declare and initialize a matrix of type integer
// with fixed size MAX
var graph [MAX][MAX]int

// declare and initialize a matrix of type int
// with fixed size MAX
var graph [MAX][MAX]int = [MAX][MAX]int{
  [MAX]int{0, 0, 0},
  [MAX]int{0, 0, 0},
  [MAX]int{0, 0, 0},
}

// declare a matrix of type int with
// fixed size MAX using array literals
// and short variable declaration
graph := [MAX][MAX]int{
  [MAX]int{0, 0, 0},
  [MAX]int{0, 0, 0},
  [MAX]int{0, 0, 0},
}
```

The next post in this series shows how to declare and initialize a slice of
strings and hash based maps.