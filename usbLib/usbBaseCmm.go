package usbLib

import (
	"errors"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

type communication interface {
	Open() error
	IsOpen() bool
	Read(buffer []byte) (int, error)
	Write(bytes string) (int, error)
	Close() error
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
	if err != nil {
		log.Errorf("read %v", err)
	} else {
		log.Debugf("read %x\n", []byte(buffer))
	}
	return
}

func (st *libUsbStmComm) Write(bytes string) (writeBt int, err error) {
	writeBt, err = st.output.Write([]byte(bytes))
	if err != nil {
		log.Errorf("write %v", err)
	} else {
		log.Debugf("write %x\n", []byte(bytes))
	}
	return
}

func (st *libUsbStmComm) Close() error {
	var err error

	if st.intf != nil {
		st.intf.Close()
		st.intf = nil
	}
	if st.conf != nil {
		st.conf.Close()
		st.conf = nil
	}
	if st.dev != nil {
		st.dev.Close()
		st.dev = nil
	}

	if st.ctx != nil {
		err = st.ctx.Close()
		st.ctx = nil
	}

	return err
}
