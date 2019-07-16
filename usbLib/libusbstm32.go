package usbLib

import (
	"fmt"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

// LibStmUsbFct: Stm32DevHelper implementation
type LibStmUsbFct struct {
	handlers []HotPlugHandler // add some other operate
	devs     []Stm32UsbDev
	ctx      *gousb.Context
}

func (l *LibStmUsbFct) get(desc *gousb.DeviceDesc) communication {
	device := libUsbStmComm{vid: desc.Vendor, pid: desc.Product}
	return &device
}

func (l *LibStmUsbFct) register(handler HotPlugHandler) {
	l.handlers = append(l.handlers, handler)
}

func (l *LibStmUsbFct) unregister(handler HotPlugHandler) {
	for i := range l.handlers {
		if handler == l.handlers[i] {
			l.handlers = append(l.handlers[:i], l.handlers[i+1:]...)
			return
		}
	}
}

func (l *LibStmUsbFct) Close() error {
	var err error
	if l.ctx != nil {
		err = l.ctx.Close()
		l.ctx = nil
	}
	return err
}

func (l *LibStmUsbFct) list() ([]Stm32UsbDev, error) {
	l.Init()
	return l.devs, nil
}

func (l *LibStmUsbFct) Init() {
	if l.ctx == nil {
		l.ctx = gousb.NewContext()
		l.platformInit()
	}
}

// work via hotplug
func (l *LibStmUsbFct) platformInit() {
	l.ctx.RegisterHotplug(l.usbevt)
}

func (l *LibStmUsbFct) usbevt(e gousb.HotplugEvent) {
	log.Debugf("Got event %v\n", e)

	usbDesc, err := e.DeviceDesc()
	if err != nil {
		log.Println(err)
		return
	}

	if isSupported(usbDesc) {
		var stmDev *Stm32UsbDev
		switch e.Type() {
		case gousb.HotplugEventDeviceArrived:
			stmDev, err = l.mkStm32UsbDev(usbDesc)
			if err != nil {
				log.Println(err)
				return
			}
			l.devs = append(l.devs, *stmDev)
		case gousb.HotplugEventDeviceLeft:
			uuid := deviceUUID(usbDesc)
			for i := range l.devs {
				if uuid == l.devs[i].UUID {
					stmDev = &l.devs[i]
					l.devs = append(l.devs[:i], l.devs[i+1:]...)
					break
				}
			}
		}

		for i := range l.handlers {
			switch e.Type() {
			case gousb.HotplugEventDeviceArrived:
				l.handlers[i].Added(*stmDev)
			case gousb.HotplugEventDeviceLeft:
				l.handlers[i].Removed(*stmDev)
			}
		}
	}
}

func (l *LibStmUsbFct) mkStm32UsbDev(usbDesc *gousb.DeviceDesc) (*Stm32UsbDev, error) {
	stm32Dev := Stm32UsbDev{Name: "stm32", UUID: deviceUUID(usbDesc)}
	cmm := l.get(usbDesc)
	err := cmm.Open()
	if err != nil {
		log.Errorf("Error opening %s : %s\n", usbDesc, err)

		cmm.Close()
		return nil, err
	}

	// can initial some other ...

	return &stm32Dev, nil
}

func deviceUUID(d *gousb.DeviceDesc) string {
	return fmt.Sprintf("%d:%d", d.Bus, d.Address)
}

func isSupported(dd *gousb.DeviceDesc) bool {
	return dd.Vendor == gousb.ID(apogee) && dd.Product == gousb.ID(sp420)
}
