var tt = class {
    private multi:int
    Sumit = method(x:int,y:int) int {
        return (x+y) * this.multi
    }
}
tt.multi = 2
print tt.Sumit(6,3)

var fn = function(x:class) int {
    return x.Sumit(32,5)
}

print fn(tt)
