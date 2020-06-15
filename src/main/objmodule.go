package main
/* ---------------------------------------------------------------------------------------------
The module is the fundamental code block which contains functions, variables, and other objects
all inside of it's own namespace. The initial, unnamed module is 'main' which is the starting
point for the application. As modules are loaded, their objects become visible, either implicitly,
or via the namespace identifier
---------------------------------------------------------------------------------------------------*/

/*
The module object is very simple, but it serves as a reference point for module level
objects such as functions and
*/
type ObjModule struct {
	Parent *ObjModule // Modules can be nested. This refers to the parent module if there is one
	Name string // Initial name of the module
}

func (o ObjModule) ShowValue() string {
	panic("implement me")
}

func (o ObjModule) Type() ValueType {
	panic("implement me")
}

func (o ObjModule) ToBytes() []byte {
	panic("implement me")
}

func (o ObjModule) ToValue() interface{} {
	panic("implement me")
}

func (o ObjModule) Print() string {
	panic("implement me")
}

