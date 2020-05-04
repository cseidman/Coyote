var tt = 8+7
println(tt)

var t = class {
    a:int
    b = method(x:int,y:int) int {
        return (x+y)*this.a
    }
}
t.a = 2
println(t.b(5,4))

var x = enum{
ONE,
TWO,
THREE,
FOUR
}

print(x.TWO)


