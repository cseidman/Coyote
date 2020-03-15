package main

type Chunk struct {
	Count          int
	Capacity       int
	Lines          []int
	Code           []byte
	ConstantsCount int
	Constants      []Obj
}

func NewChunk() Chunk {
	return Chunk{
		Count:          0,
		Capacity:       65000,
		Lines:          make([]int, 65000),
		Code:           make([]byte, 65000),
		ConstantsCount: 0,
		Constants:      make([]Obj, 65000),
	}
}
