var myClass = class {
    int a
    int b
    sum(x:int y:int) int {
        return x+y+this.a+this.b
    }

}

var x = new myClass
x.a = 6
x.b = 4

println(x.sum(3,4))
// 17
