var recurse_fibonacci = function(n:int) int {
    if n <= 1 {
        return n
    } else {
        return recurse_fibonacci(n-1) + recurse_fibonacci(n-2)
    }
}
for i = 0 to 34 {
//var i = 0
//while i<= 34 {
    print recurse_fibonacci(i)
   // i = i+1
}

