package usbLib

import (
	"reflect"
	"unsafe"
)

// CBW  DATA  CSW
// Section 5.1: Command Block Wrapper (CBW)
type CommandBlockWrapper struct {
	dCBWSignature          [4]uint8
	dCBWTag                uint32
	dCBWDataTransferLength uint32
	bmCBWFlags             uint8
	bCBWLUN                uint8
	bCBWCBLength           uint8
	CBWCB                  [16]uint8
}

// Section 5.2: Command Status Wrapper (CSW)
//struct command_status_wrapper {
//uint8_t dCSWSignature[4];
//uint32_t dCSWTag;
//uint32_t dCSWDataResidue;
//uint8_t bCSWStatus;
//};

var sizeOfMyStruct = int(unsafe.Sizeof(CommandBlockWrapper{}))

func MyStructToBytes(s *CommandBlockWrapper) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfMyStruct
	x.Cap = sizeOfMyStruct
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}
