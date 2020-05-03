var z = class {
    a:int
    b = method(x:int, y:int) int {
        return (x+y) * this.a
    }
}
z.a = 2
println(z.b(5,4))

var x = @{"a":10,"b":20,"c":30}
print(x["a"])

x["d"] = 100

print(x["d"])