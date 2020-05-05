package main

import (
	"gonum.org/v1/gonum/mat"
)

type ObjMatrix struct {
	Data mat.Matrix
}

func MakeDense() {
	x := &ObjMatrix{
		Data: mat.NewDense(3, 5, nil),
	}

}
