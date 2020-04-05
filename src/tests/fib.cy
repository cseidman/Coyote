var fib = function(n:int) {
  if n < 2 {
    return n
  } else {
    return fib(n - 1) + fib(n - 2)
  }
}
var i = 6
print fib(i)
