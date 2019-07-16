package usbLib

import (
	"errors"
)

const (
	apogee = 0xba57 // 0xba57
	sp420  = 0x1001 // 0x1001
)

// HotPlugHandler hotplug apogee wrapper
type HotPlugHandler interface {
	Added(stmUsb Stm32UsbDev)
	Removed(stmUsb Stm32UsbDev)
}

func NewStm32Manager() *LibStmUsbFct {
	return &LibStmUsbFct{}
}

// FindUsbOne - Return first attached device
// Or error if no device found
// Note that device is opened and must be closed with Close() method after use
func FindUsbOne(factory *LibStmUsbFct) (*Stm32UsbDev, error) {

	devices, err := FindUsb(factory)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, errors.New("no device found")
	}
	return &devices[0], err
}

// FindUsb - Return list of attached Apogee devices
// Note that each device must be closed with Close() method after use
func FindUsb(factory *LibStmUsbFct) ([]Stm32UsbDev, error) {
	return factory.list()
}

// RegisterHandler register hotplug handler
func RegisterHandler(factory *LibStmUsbFct, handler HotPlugHandler) {
	factory.register(handler)
}

// UnRegisterHandler unregister hotplug handler
func UnRegisterHandler(factory *LibStmUsbFct, handler HotPlugHandler) {
	factory.unregister(handler)
}
