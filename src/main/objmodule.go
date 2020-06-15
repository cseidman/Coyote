package main

import(
	"fmt"
)


/* ---------------------------------------------------------------------------------------------
The module is the fundamental code block which contains functions, variables, and other objects
all inside of it's own namespace. The initial, unnamed module is 'main' which is the starting
point for the application. As modules are loaded, their objects become visible, either implicitly,
or via the namespace identifier
---------------------------------------------------------------------------------------------------*/

/*
On the compiler side, this is where we keep track of all nitty-grotty details about module allocation
*/
type ObjModule struct {
	ParentModule *ObjModule
	Name string
	IsUsed bool // Is it used anywhere?
	LoadedModules []*ObjModule // The modules that this module loads
	ModuleCount int // How many loaded modules we have
	ModuleId []byte // This is an id value that contains information starting from the parent
	MainFunction *ObjFunction // Starting function in the module
}

var moduleSequence = int16(0)

func (c *Compiler) ImportStatement() {
	// import <modulepath>
	c.Consume(TOKEN_IDENTIFIER,"Expect module name after 'import'")
	moduleName := c.Parser.Previous.ToString()
	for {
		if !c.Match(TOKEN_SLASH) {
			break
		}
	}

	var moduleObj *ObjModule

	idx := c.ResolveModule(moduleName)
	if idx != -1 {
		moduleObj = c.CompileModule(moduleName, false)
		
	}

	c.CurrentModule.LoadedModules[c.CurrentModule.ModuleCount] = &c.Modules[idx]
	c.EmitInstr(OP_IMPORT, int16(idx))
	c.WriteComment(fmt.Sprintf("Import module %s",moduleName))

}

func (c *Compiler) ResolveModule(moduleName string) int {

	curMod := c.CurrentModule
	modName := moduleName

	// Checks to see if this module had been registered
	for i:=0;i<c.ModuleCount-1;i++ {
		// We already have this module by comparing the name and the parental lineage
		for c.Modules[i].Name == modName {//&& bytes.Compare(c.Modules[i].ParentModule.ModuleId,curMod.ParentModule.ModuleId) {
			// Check to see of this is a compound definition
			// TODO: When we start adding nested modules, get back to this
			//if c.Match(TOKEN_DOT) {
				// If so, check the next definition level
			//	curMod = c.Modules[i]
			//	modName = curMod.Name
			//	continue
			//}
			return i
		}
	}
	return -1
}

func (c *Compiler) RegisterModule(moduleName string) int {

	idx := c.ResolveModule(moduleName)
	// If we found an existing module
	if idx > -1 {
		return idx
	}

	// Get the value of the parent module. This is a string of bytes that provide a lineage trail
	// that uniquely identfies a module of the same name but with different parental pathway. Each
	// New module gets a 2-byte number which then is tacked on to the parent's unique identifier.
	// That makes it possible to differentiate "file.csv.open" from "stream.csv.open" because
	// The two references could have Id values (ex:) 010203 vs 040506
	// Note: This is relevant only when we get to nested modules
	modId := append(c.CurrentModule.ModuleId,Int16ToBytes(moduleSequence)...)
	moduleSequence++

	mod := ObjModule{
		ParentModule: c.CurrentModule,
		Name:         moduleName,
		IsUsed:       false,
		ModuleId: 	  modId,
	}

	c.Modules[c.ModuleCount] = mod
	c.CurrentModule = &mod
	c.ModuleCount++
	return c.ModuleCount-1
}

func (c *Compiler) DeclareModule() {
	c.Consume(TOKEN_IDENTIFIER,"Expect module name after 'module'")
	moduleName := c.Parser.Previous.ToString()

	c.RegisterModule(moduleName)
}


/*
The module object is very simple, but it serves as a reference point for module level
objects such as functions and
*/

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

