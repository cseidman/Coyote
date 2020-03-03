package main

import (
	."../common"
	."../compiler"
	."../vm"
	"flag"
	"fmt"
	//"os"
	"runtime/debug"
)

func main() {

	flag.Bool("debug", false, "debug mode")

	source := flag.String("f", "", "source file")

	flag.Parse()

	//DbgMode = *dbg
	SrcFile := *source

	debug.SetGCPercent(-1)

	// Load Globals slot
	LoadGlobals()

	if SrcFile == "" {
		repl()
	} else {
		RunFile(SrcFile)
	}

}

func repl() {

	var line string
	fmt.Println("Coyote Copyright (C) 2020  Claude Seidman")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; for details type 'show w'.")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Println("under certain conditions; type 'show c' for details.")
	fmt.Println()
	for {
		fmt.Printf("> ")
		if _, err := fmt.Scanln(&line); err != nil {
			Exec(&line)
		}
	}
}

func RunFile(path string) {

	debug.SetGCPercent(-1)

	source := ReadFile(path) + "\n"
	Exec(&source)
/*
	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}
*/

}


