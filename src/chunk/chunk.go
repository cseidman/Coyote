package chunk

type Chunk struct {
	Count int
	Capacity int
	Lines []int
	Code []byte
}

func NewChunk() Chunk {
	return Chunk {
		Count: 0,
		Capacity: 1024,
		Lines: make([]int,1024),
		Code: make([]byte,1024),
	}
}

func (c *Chunk) WriteChunk(b byte, line int) {

	if c.Capacity < c.Count + 1 {
		c.Capacity+=1024
		c.Lines = append(c.Lines,make([]int,1024)...)
		c.Code = append(c.Code,make([]byte,1024)...)
	}

	c.Code[c.Count] = b
	c.Lines[c.Count] = line
	c.Count++

}
