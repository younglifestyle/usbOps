package usbLib

import (
	"fmt"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

type Stm32DevHelper interface {
	list() ([]Stm32UsbDev, error)
	register(handler HotPlugHandler)
	unregister(handler HotPlugHandler)
	Close() error
	Init()
}

// as stm32-usb Instance, including the operation, eg. readSome, writeSome
type Stm32UsbDev struct {
	Name string
	UUID string
	cmm  communication

	// can set define device specify parameter, eg. sn, hard-version ...
}

// Close device
func (s *Stm32UsbDev) Close() error {
	return s.cmm.Close()
}

func (s *Stm32UsbDev) Control(rType uint8, request uint8, val uint16, idx uint16, data []byte) (int, error) {
	if s.cmm == nil {
		fmt.Println("cmm nil")
		return 0, nil
	}

	return s.cmm.Control(rType, request, val, idx, data)
}

func (s *Stm32UsbDev) Read(buffer []byte) (readBt int, err error) {

	return s.cmm.Read(buffer)
}

func (s *Stm32UsbDev) Write(input []byte) (writeBt int, err error) {

	log.Debug("is write ")

	return s.cmm.Write(input)
}

func (s *Stm32UsbDev) SendMassStorageCommand(input []byte) (writeBt int, err error) {

	var cbw CommandBlockWrapper
	cbw.dCBWSignature[0] = 'U'
	cbw.dCBWSignature[1] = 'S'
	cbw.dCBWSignature[2] = 'B'
	cbw.dCBWSignature[3] = 'C'
	cbw.dCBWTag = 0x89
	cbw.dCBWDataTransferLength = 32
	cbw.bmCBWFlags = gousb.ControlOut

	// gousb.ControlOut : 0x0 ,   gousb.ControlIn : 0x80

	cbw.bCBWLUN = 0
	cbw.bCBWCBLength = uint8(len(input))
	for index, val := range input {
		cbw.CBWCB[index] = val
	}
	bytes := MyStructToBytes(&cbw)

	return s.cmm.Write(bytes[:31])
}
