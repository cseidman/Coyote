package common

import "io/ioutil"

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
