package bluetooth

/*
#cgo CFLAGS: -w -I .
#cgo LDFLAGS: -L../../lib -lbtferret -Wl,-rpath=./lib
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "btlib.h"
#include "util.h"


*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/google/uuid"
)

type BTDevice struct {
	id  uuid.UUID
	buf []byte
}

func (device *BTDevice) GetBuffer() []byte {
	return device.buf
}

func NewBTDevice() *BTDevice {
	return &BTDevice{uuid.New(), make([]byte, 0)}
}

func (device *BTDevice) InitServer() error {
	file := C.CString("devices.txt")
	if C.init_blue(file) == 0 {
		return fmt.Errorf("Error: Unable to initalize classic bluetooth server.")
	}
	return nil
}

func (device *BTDevice) StartServer() {

	//for more information on availeable flags, reference: https://github.com/petzval/btferret?tab=readme-ov-file#4-2-2-classic_server
	C.classic_server(C.ANY_DEVICE, (*[0]byte)(C.handleConnection), 10, C.serverKeyflag)
	dataPtr := unsafe.Pointer(C.serverDataBuffer)
	defer C.free(dataPtr)
	data := C.GoString((*C.char)(dataPtr))
	device.buf = []byte(data)
	fmt.Printf("Data from golang!: %v", data)
	C.close_all()

}
