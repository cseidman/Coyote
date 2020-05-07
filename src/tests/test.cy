var x = newarray(5,int)
x[0] = 100
x[1] = 101
x[2] = 1
println(x[0])
println(x[1])
println(x[x[2]])

var y = @{"One":100,"Two":200,"Three":300}
println(y$Three)

y$Four = 400
println(y$Four)
