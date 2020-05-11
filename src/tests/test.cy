var int x = 2
println(x)

var float y = 3.0
println(y)

var []int z = @[1,2,4,7,8]
println(z[3])

var list h = @{"one":100, "two":200, "three":300}
println(h$two)

var c = class {
    a:int
    b = method() {
        println("Hey")
    }
}
c.b()
