var tt = class {
    private multi:int
    Sumit = method(x:int,y:int) int {
        return (x+y)*this.multi
    }
}
tt.multi = 2
print tt.Sumit(7,8)
