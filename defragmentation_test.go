package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	usedMemory := make(map[unsafe.Pointer]struct{}, len(pointers))
	for i := 0; i < len(pointers); i++ {
		usedMemory[pointers[i]] = struct{}{}
	}

	write := 0
	for read := 0; read < len(memory); read++ {
		if _, ok := usedMemory[unsafe.Pointer(&memory[read])]; ok {
			if read != write {
				memory[write] = memory[read]
				for i := range pointers {
					if pointers[i] == unsafe.Pointer(&memory[read]) {
						pointers[i] = unsafe.Pointer(&memory[write])
					}
				}
			}
			write++
		}
	}

	for i := write; i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
