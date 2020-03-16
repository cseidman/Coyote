var sum = function(x:int, y:int) int {
    return x+y
}

var exec = function(f:function,x:int,y:int) int {
    return f(x,y)
}

print exec(sum,4,5)
