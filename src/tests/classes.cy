define class myClass {
    properties {
        int a
        int b
    }
    methods {
        void myClass() {
            this.a = 6
            this.b = 4
        }
        int sum(x:int y:int) {
            return x+y+this.a+this.b
        }
    }
}

var myClass x = new(myClass)
println(x.sum(3,4))
// Should print '17'
