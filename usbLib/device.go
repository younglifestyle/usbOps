package usbLib

import (
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

func (s *Stm32UsbDev) Read() (readBt int, err error) {

	log.Debug("is read ")

	return
}

func (s *Stm32UsbDev) Write() (writeBt int, err error) {

	log.Debug("is write ")

	return
}
