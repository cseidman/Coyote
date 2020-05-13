var int x = 2
println(x)

var float y = 3.0
println(y)

var []int z = @[1,2,4,7,8]
println(z[3])

var list h = @{"one":100, "two":200, "three":300}
println(h$two)

var class c = class {
    a:int
    b = method() {
        println("Hey")
    }
    sum = method(x:int, y:int) int {
        return x+y+this.a
    }
}
c.b()
c.a = 100
println(c.sum(4,5))

var int s = 200
var func d = func(x:int, y:int) int {
    var int s = 100
    return x*y+s
}
println(d(2,6))
println(s)

var []float dist = dnorm(100,5.0,2.0)
for i = 0 to 20 {
    println(dist[i])
}