package life

import (
	"bytes"
	"unsafe"
)

/*
#cgo CFLAGS: -IC:/home/h1b/dev/sandbox/go-ws-life/life
#cgo LDFLAGS: /home/h1b/dev/sandbox/conways-game-of-life/target/release/liblife.so
#include "liblife.h"
*/
import "C"

// func TestTheThing() {
// 	fmt.Println("Start")

// 	state_ptr := C.init_state_random(4, 4, 4)
// 	txt_ptr := C.next_state_img(state_ptr)

// 	data := []byte(C.GoString((*C.char)(txt_ptr)))

// 	fmt.Printf("GO:state_ptr:%s\n", state_ptr)
// 	fmt.Printf("GO:txt_ptr:%s\n", txt_ptr)

// 	fmt.Printf("GO:data:%s\n", data)

// 	//

// 	txt_ptr = C.next_state_img(state_ptr)

// 	data = []byte(C.GoString((*C.char)(txt_ptr)))

// 	fmt.Printf("GO:state_ptr:%s\n", state_ptr)
// 	fmt.Printf("GO:txt_ptr:%s\n", txt_ptr)

// 	fmt.Printf("GO:data:%s\n", data)

// 	C.free_char_p(txt_ptr)
// 	C.free_void_p(state_ptr)

// 	fmt.Println("End")
// }

type LifeNativeRust struct {
	Width int `json:"width"`
	Size  int `json:"size"`
	ptr   unsafe.Pointer
}

func (life *LifeNativeRust) Free() {
	C.free_void_p(life.ptr)
}

func (life *LifeNativeRust) Life() {
	// Dummy

	txt_ptr := C.next_state_img(life.ptr)
	C.free_char_p(txt_ptr)
}

func (life *LifeNativeRust) DrawImageDataUrl() ([]byte, error) {
	txt_ptr := C.next_state_img(life.ptr)
	defer C.free_char_p(txt_ptr)

	// result := []byte(C.GoString((*C.char)(txt_ptr)))

	result := bytes.Clone([]byte(C.GoString((*C.char)(txt_ptr))))

	return result, nil
}

func NewLifeNativeRust(width int) *LifeNativeRust {
	state_ptr := C.init_state_random(C.int(width), C.int(width), 4)

	return &LifeNativeRust{
		Width: width,
		Size:  width * width,
		ptr:   state_ptr,
	}
}
