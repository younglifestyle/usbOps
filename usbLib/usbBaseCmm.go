package usbLib

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

type communication interface {
	Open() error
	IsOpen() bool
	Read(buffer []byte) (int, error)
	Write(bytes []byte) (int, error)
	Close() error
	Control(rType uint8, request uint8, val uint16, idx uint16, data []byte) (int, error)
}

// libUsbStmComm: libusb communication implementation
type libUsbStmComm struct {
	vid    gousb.ID
	pid    gousb.ID
	ctx    *gousb.Context
	dev    *gousb.Device
	intf   *gousb.Interface
	conf   *gousb.Config
	output *gousb.OutEndpoint
	input  *gousb.InEndpoint
}

func (st *libUsbStmComm) Open() error {
	var err error
	st.ctx = gousb.NewContext()
	st.dev, err = st.ctx.OpenDeviceWithVIDPID(st.vid, st.pid)
	if err != nil {
		log.Error("open device is failed,", err)
		st.Close()
		return err
	}
	if st.dev == nil {
		return errors.New("no device, no error")
	}

	log.Debug("Set auto-detach")
	st.dev.SetAutoDetach(true)

	st.conf, err = st.dev.Config(1)
	if err != nil {
		log.Println(err)
		st.Close()
		return err
	}

	st.intf, err = st.conf.Interface(0, 0)
	if err != nil {
		log.Error("interface error,", err)
		st.Close()
		return err
	}
	log.Debugf("interface : %+v", st.intf.Setting)

	for epNum, endpoints := range st.intf.Setting.Endpoints {

		fmt.Println("TransferType : ", endpoints.TransferType, "direction : ", endpoints.Direction)

		if endpoints.Direction == gousb.EndpointDirectionOut {
			st.output, err = st.intf.OutEndpoint(int(epNum))
			if err != nil {
				log.Error("set outpoint error,", err)
				st.Close()
				return err
			}
		} else {
			st.input, err = st.intf.InEndpoint(int(epNum))
			if err != nil {
				log.Error("set inpoint error,", err)
				st.Close()
				return err
			}
		}
	}

	return nil
}

func (st *libUsbStmComm) IsOpen() bool {
	return st.output != nil && st.input != nil
}

func (st *libUsbStmComm) Read(buffer []byte) (readBt int, err error) {
	readBt, err = st.input.Read(buffer)

	fmt.Println("read base : ", readBt)

	return
}

func (st *libUsbStmComm) Write(bytes []byte) (writeBt int, err error) {
	writeBt, err = st.output.Write([]byte(bytes))

	return
}

func (st *libUsbStmComm) Close() error {
	var err error

	if st.intf != nil {

		fmt.Println("inter close")

		st.intf.Close()
		st.intf = nil
	}
	if st.conf != nil {
		fmt.Println("conf close")

		st.conf.Close()
		st.conf = nil
	}
	if st.dev != nil {
		fmt.Println("dev close")

		st.dev.Close()
		st.dev = nil
	}

	if st.ctx != nil {
		fmt.Println("ctx close")

		err = st.ctx.Close()
		st.ctx = nil
	}

	return err
}

func (st *libUsbStmComm) Control(rType uint8, request uint8, val uint16, idx uint16, data []byte) (int, error) {

	if st.dev == nil {
		fmt.Println("dev nil")
	}
	if data == nil {
		fmt.Println("data nil")
	}

	return st.dev.Control(rType, request, val, idx, data)
}
