var z = class {
    a:int
    b = method(x:int, y:int) int {
        return (x+y) * this.a
    }
}
z.a = 2
println(z.b(5,4))
