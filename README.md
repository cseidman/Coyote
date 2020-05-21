# Coyote
Welcome to Coyote - a fast, lightweight language designed for data engineers in mind. It lets you use the best features of both Functional and Object-Oriented languages. The philosophy of the Coyote language is to incorporate the power of a full-fledged language with built-in SQL databases and OLAP so that the tight integration between both provides a seamless experience that adds power to Data Science and Data Analytics.

Unlike scripting languages, Coyote makes you declare your variables and provides the security that comes from compile-time type checking. At the same time, it offers extensive data manipulation structures that can be combined as needed. Multi-dimensional arrays can contain classes and functions, functions can be stored in Data Frames, and all objects can be passed as parameters to functions or methods.

To get started, you may download the latest version from https://github.com/cseidman/Coyote/releases 

## Table of Contents
* [Quick Introdution](#quick-introduction)
* [Variable Types](#variable-types)
* [Composite Types](#composite-types)
* [Declaring Variables](#declaring-variables)

# Quick Introduction

### Variable Types
| Type | Description  |
|--|--|
| int | 64 bit integer |
| float | 64 bit float  |
| string | string |
| byte | 8 bit byte |
| bool | boolean true/false |

### Composite Types
| Type | Description  |
|--|--|
| class | lightweight OOP-style class |
| enum | enum of type byte  |
| array | collection of variables |
| list | associative array/hash table  |
| matrix | mathematical matrix |

## Declaring Variables

Variables are created by using the `var` statement followed by the variable type and name. Ex: 
`var int MyVariable` or you can initialize it at the same time like this: `var MyVariable = 45` and the variable will be initialized by the value type

In Coyote, you must assign a type to a variable, which stays the same until it's re-declared. Ex:
```{Coyote}
var lastName string
lastName = "Jones"

// This is fine
var firstName = "Fred"

var x int
x = 100
println(x)

var y = 200

// Error
firstName = 100
[line 9] Error at '100': Variable firstName is a Scalar of type string: cannot assign a Scalar of type integer
```

#### Arrays
Declaring arrays 
```
var x int[] 
var y float[] 
var z bool[] 
...
```
You can also declare it as a sized array
```
var x = new int[3]
x[0] = 100
x[1] = 101
x[2] = 1

println(x[1])
// 101

```
If the array is initialized at the same time as it's declared, it takes the size of the initializer
```
var y = @[10,20,30]
println(y[1])
// 20
```
A declared but uninitialized array can be initialized (and sized) later 
```
var z int[]
z = @[200,201,2]
println(z[1])
// 201
```

#### Multi-Dimensional Arrays
Dimensions in an array are delimited by commas. One comma indicates two dimension, two commas mean three dimensions, and so on. There is no practical limit to how many dimensions you can declare in an array. 
```
var x = new int[3,3]
var y = 0
for i = 0 to 2 {
    for j = 0 to 2 {
        x[i,j] = y
        y = y + 1
    }
}
println(x[1,1])
// 4
```
As with regular arrays, the variable can be sized in advance
```
var m = new int[2,3,4]
var y = 0
for i = 0 to 1 {
    for j = 0 to 2 {
        for v = 0 to 3 {
            m[i,j,v] = y
            y = y + 1
        }
    }
}
println(m[1,1,1])
// 17

```
To declare a multi-dimensional array and initialize it at the same time, you can add ```[int,int]``` at the beginning of the declaration of the array elements:
```
var x = @[[3,3]0,1,2,3,4,5,6,7,8]
x[2,2] = 4
x[1,1] = 1
x[0,0] = 100
println(x[2,2])
```  
#### Lists 
Lists contain elements of different types like − numbers, strings, arrays and even another list inside it. A list can also contain a matrix or a function as its elements. List is created as follows:
```
var l = @{"One":1, "Two":2, "Three":3}

var veggies = list[string,float]
veggies$Tomatoes = 2.00
veggies$Celery = 3.50
veggies$Spinach = 2.75

println(veggies$Celery)
// 3.50000
```
**List of arrays:**
```
var x = @{
        "Q1":@["Jan","Feb","Mar"],
        "Q2":@["Apr","May","June"]
        }
println(x$Q2[1])
// May
```  
**Array of Lists:**
```
var food = @[
    @{"Carrots":1.75,"Celery":3.50, "Onions":0.75},
    @{"Beef":4.55,"Pork":5.75,"Chicken":2.80}
]
println(food[0]$Celery)
println(food[1]$Pork)
// 3.5000
// 5.7500
```
#### Enums
Enums elements represent ```int``` values
```
var size = enum {
    XTRALARGE,
    LARGE,
    MEDIUM,
    SMALL
}

var x = size.LARGE

if x == size.LARGE {
    println("It's large")
}
// It's large
```

### Functions
Functions don't have names, they return a function *type* which is stored in a variable. If you pass parameters, you must use a name:type expression followed by a return type if there is one. If there is a declared return type, it must be explicitely returned with the ```return``` keyword
```
var f = func(x:int, y:int) int {
    return x * y
}

println(f(4,5))
// 20

```
Functions can be passed as parameters to other functions:
```
var f = func(fn:func) {
    fn()
}

var SayHello = func() {
    println("Hi")
}

var SayBye = func() {
    println("Bye")
}

f(SayHello)
f(SayBye)

// Hi
// Bye
```
A closure in Coyote is a function that is able to bind objects the closure used in the environment is was created in. These functions maintain access to the scope in which they were defined, allowing for powerful design patterns similar to concepts of functional programing

Suppose you want a function that adds 2 to its argument. You would likely write something like this:
```
var add_2 = func(y:int) int {
    return 2 + y
}
add_2(5)
// 7
```
Now suppose you need another function that instead adds 5 to its argument. The natural thing to do would be to write another function, just like add_2, where the 2 is replaced with a 7. But this would be grossly inefficient: if in the future you discover that you made a mistake and you in fact need to multiply the values instead of add them, you would be forced to change the code in two places. In this trivial example, that may not be much trouble, but for more complicated projects, duplicating code is a recipe for disaster.

A better idea would be to write a function that takes one argument, x, that returns another function which adds its argument, y, to x. In other words, something like this:

```
var add_x = func(x:int) func {
   return func(y:int) int {
        return x+y
   }
}

var f = add_x(7)
println(f(5))
// 12

var g = add_x(10)
println(g(5))

// 15

```

