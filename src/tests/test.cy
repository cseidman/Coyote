var xx = function(a:int) function {
    var tt = function(x:int) int {
        return x * a
    }
    return tt
}

var y = xx(2)
var a = y(5)
var b = y(10)

if a == 10 and b == 30 {
    print "OK"
} else {
    print "Nope"
}
