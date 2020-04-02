var x = class {
    a:int
    b = method() {
        print this.a
    }
    c = method() {
        this.a = this.a + 11
        this.b()
    }


}

var y = class {
    a:int
    b = method() {
        print this.a
    }
    c = method() {
        this.a = this.a + 11
        this.b()
    }


}


x.a = 100
y.a = 200
x.b() // 100
y.b() // 200
x.c() // 111
y.c() // 211

