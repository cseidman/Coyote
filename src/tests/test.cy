var int[,] x = int[3,3]
var int y = 0
for i = 0 to 2 {
    for j = 0 to 2 {
        x[i,j] = y
        y = y + 1
    }
}
println(x[2,2])
println(x[1,2])
println(x[0,0])

var int[] t = int[3]
t[0] = 0
t[1] = 1
t[2] = 2

println(t[0])
println(t[1])
println(t[2])

var int[,,] m = int[2,3,4]
y = 0
for i = 0 to 1 {
    for j = 0 to 2 {
        for v = 0 to 3 {
            m[i,j,v] = y
            y = y + 1
        }
    }
}

println(m[1,1,1])