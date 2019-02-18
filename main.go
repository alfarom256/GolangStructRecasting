package main

import (
	"fmt"
	"math"
	"unsafe"
)

// make a big struct with some data we want to put into memory
// at the end with length N and total size K (struct + payload)
type larger_struct struct {
	a1 uint64
	a2 uint64
	a3 uint64
	a4 uint64
	a5 uint64
	a6 bool
	a7 uintptr
}

// make a smaller struct of size K - payload size and reinterpret the cast
// then when it is free-ed
type smaller_struct struct {
	x1 uint64
	x2 uint64
	x3 uint64
	x4 uint64
}

func main() {
	fmt.Println("Allocating objects...")
	large := larger_struct{
		a1: 0xFAFAFAFAFAFAFAFA,
		a2: 0xFFFFFFFFFFFFFFFF,
		a3: 0xFFFFFFFFFFFFFFFF,
		a4: 0xFFFFFFFFFFFFFFFF,
		a5: 0xDEADBEEFDEADBEEF,
		a6: true,
		a7: 0x1337133713371337,
	}
	large2 := larger_struct{
		a1: 0xFEFEFEFEFEFEFEFE,
		a2: 0xFEFEFEFEFEFEFEFE,
		a3: 0xFEF1F0FEF1F0FEF1,
		a4: 0xFFFFFFFFFFFFFFFF,
		a5: 0xDEADBEEFDEADBEEF,
		a6: true,
		a7: 0x1337133713371337,
	}
	large = dummyHeapTest(large)
	fmt.Printf("Address of struct: %p\n", &large)
	fmt.Printf("Address of struct2: %p\n", &large2)
	fmt.Printf("Contents of Memory (dumping 2x):\n")
	large_addr := uintptr(unsafe.Pointer(&large))

	// store a copy of the memory region
	large_orig := dumpMemSize(large_addr, 4*uint32(unsafe.Sizeof(large)))

	//pretty print it
	prettyMemPrint(large_orig, large_addr)

	fmt.Println("\nRecasting Struct 1 to smaller struct")

	// try to force cast reinterpretation
	small := smaller_struct{
		x1: 0x0,
		x2: 0x0,
		x3: 0x0,
		x4: 0x0,
	}
	large = dummyHeapTest(large)

	// storing address of large into a variable
	large_addr = uintptr(unsafe.Pointer(&large))

	// here we're attempting to reinterpret the cast of large
	// by taking the address and casting it to a smaller_struct
	// kinda meaningless since small is a new object
	small = *(*smaller_struct)(unsafe.Pointer(&large))

	// Invalid indirect
	//*small = unsafe.Pointer(&large)
	small_addr := uintptr(unsafe.Pointer(&small))

	fmt.Printf("Address of small: %x\n", small_addr)
	small.x1 = 18085043209519168250 // 0x FAFAFAFAFAFAFA...
	small.x2 = 18085043209519168250 //    '   '    '    '    '
	fmt.Println("Assigning new values to struct")
	fmt.Println("Contents of New Memory (dumping 2x):\n")

	large_new := dumpMemSize(large_addr, 4*uint32(unsafe.Sizeof(large)))
	prettyMemPrint(large_new, large_addr)

}

// "trying to make stuff on the heap", I said at 4am
func dummyHeapTest(large larger_struct) larger_struct {
	return larger_struct{
		large.a1,
		large.a2,
		large.a3,
		large.a4,
		large.a5,
		large.a6,
		large.a7,
	}
}

func dumpMem(begin uintptr, end uintptr) []byte {
	size_of_mem := uint(begin - end)
	out := make([]byte, size_of_mem)
	for i := range out {
		out[i] = *((*byte)(unsafe.Pointer(uintptr(begin) + uintptr(i))))
	}
	return out
}

func dumpMemSize(begin uintptr, size uint32) []byte {
	out := make([]byte, size)
	for i := range out {
		out[i] = *((*byte)(unsafe.Pointer(uintptr(begin) + uintptr(i))))
	}
	return out
}

func prettyMemPrint(mem []byte, begin uintptr) {
	count := len(mem)
	row_count := math.Ceil(float64(count) / 16)
	row_remainder := count % 16
	row_remainder--
	rows := make([][]byte, int(row_count))

	for i := range rows {
		rows[i] = make([]byte, 16)
	}

	for i := 0; i < int(row_count); i++ {
		_begin := i * 16
		_end := _begin + 15
		if i+1 >= int(row_count) {
			if row_remainder == -1 { // if the memory we're analyzing fits perfectly into 16 bytes, don't bother trimming
				continue
			}
			_end = _begin + row_remainder
		}
		rows[i] = mem[_begin:_end]
	}
	for i := range rows {
		fmt.Printf("\n")
		for j := range rows[i] {
			fmt.Printf("%x ", rows[i][j])
		}
	}
}
