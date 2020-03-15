package main

import (
	"encoding/binary"
	"math"
)

/* ---------------------------------------------------------------
These functions are all designed to convert specific datatypes to
and from their []byte values
------------------------------------------------------------------ */

func BoolToBytes(val bool) []byte {
	b := make([]byte, 1)
	if val {
		b[0] = 1
	} else {
		b[0] = 0
	}
	return b
}

func BytesToBool(b []byte) bool {

	if b[0] == 1 {
		return true
	} else {
		return false
	}
}

func Int32ToBytes(val int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(val))
	return b
}

func BytesToInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func Int16ToBytes(val int16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(val))
	return b
}

func BytesToInt16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

func Float32ToBytes(val float32) []byte {

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(val))
	return b
}

func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

func Int64ToBytes(val int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(val))
	return b
}

func BytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func Float64ToBytes(val float64) []byte {
	b := make([]byte, 8)
	u := math.Float64bits(val)
	binary.BigEndian.PutUint64(b[:], u)
	return b
}

func BytesToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

func OperandToBytes(opr int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(opr))
	return b
}

// Array operations

// Int32 Arrays
func Int32Array2Bytes(val []int32) []byte {
	size := len(val) * 4
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b = append(b, Int32ToBytes(val[i])...)
	}
	return b
}

func BytesToInt32Array(b []byte) []int32 {
	size := len(b) / 4
	val := make([]int32, size)
	for i := 0; i < size; i++ {
		val[i] = BytesToInt32(b[i*4 : (i*4)+4])
	}
	return val
}

// Float32 arrays

func Float32Array2Bytes(val []float32) []byte {
	size := len(val) * 4
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b = append(b, Float32ToBytes(val[i])...)
	}
	return b
}

func BytesToFloat32Array(b []byte) []float32 {
	size := len(b) / 4
	val := make([]float32, size)
	for i := 0; i < size; i++ {
		val[i] = BytesToFloat32(b[i*4 : (i*4)+4])
	}
	return val
}
