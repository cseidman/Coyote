package common

import (
	. "../parser"
	"fmt"
	"io/ioutil"
)

func ReadFile(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	ErrCheck(err)

	return string(dat)
}

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func PrintParser(parser Parser) {
	fmt.Printf("Previous %d:%s ", parser.Previous.Type, parser.Previous.ToString())
	fmt.Printf("Current %d:%s\n", parser.Current.Type, parser.Current.ToString())
}
