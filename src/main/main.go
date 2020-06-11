/*
Copyright (C) 2020  Claude Seidman

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main


import (
	"flag"
	"fmt"
	"runtime/debug"
	"time"
)

func main() {

	dbg := flag.Bool("debug", false, "debug mode")
	source := flag.String("f", "", "source file")

	flag.Parse()

	dbgMode := *dbg
	SrcFile := *source

	debug.SetGCPercent(-1)

	if SrcFile == "" {
		repl()
	} else {
		start := time.Now()
		RunFile(SrcFile, dbgMode)
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("\nElapsed: %v\n", elapsed.String())

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
			Exec(&line, false)
		}
	}
}

func RunFile(path string, dbgMode bool) {

	debug.SetGCPercent(-1)

	source := ReadFile(path) + "\n"
	Exec(&source, dbgMode)
	/*
		if result == INTERPRET_COMPILE_ERROR {
			os.Exit(65)
		}
		if result == INTERPRET_RUNTIME_ERROR {
			os.Exit(70)
		}
	*/

}
