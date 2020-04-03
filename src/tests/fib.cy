var fib = function(n:int) {
  if n < 2 {
    return n
  } else {
    return fib(n - 1) + fib(n - 2)
  }
}
var i = 32
//for i=0 to 27 {

    print fib(i)
//}
