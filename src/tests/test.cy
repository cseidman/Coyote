module main

import TestModule

var x = 100
var y = TestModule::DoubleUp(x)
println(y)
// 200
