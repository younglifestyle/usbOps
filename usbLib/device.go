package usbLib

import (
	"fmt"
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
	if s.cmm == nil{
		fmt.Println("cmm nil")
		return 0, nil
	}

	return s.cmm.Control(rType, request, val, idx, data)
}

func (s *Stm32UsbDev) Read() (readBt int, err error) {

	var buffer []byte

	return s.cmm.Read(buffer)
}

func (s *Stm32UsbDev) Write(input []byte) (writeBt int, err error) {

	log.Debug("is write ")

	return s.cmm.Write(string(input))
}

func (s *Stm32UsbDev) WriteContorl() (writeBt int, err error) {

	s.cmm.Close()

	return
}