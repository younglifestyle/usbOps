package main

import (
	"fmt"
	"time"
	"usbOps/usbLib"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	testDev := usbLib.NewStm32Manager()
	stmUsbDev, err := usbLib.FindUsbOne(testDev)
	if err != nil {
		fmt.Println("find error :", err)
		return
	}
	defer stmUsbDev.Close()

	// simple test
	bytes := []byte{0xe0, 0x80, 0x02}
	// send CBW
	writeBt, err := stmUsbDev.SendMassStorageCommand(bytes, 20, uint32(time.Now().Unix()))
	fmt.Println("send command : ", writeBt, err)

	// read data
	buffer := make([]byte, 52)
	writeBt, err = stmUsbDev.Read(buffer)
	fmt.Println(writeBt, err)

	// read CSW
	writeBt, err = stmUsbDev.Read(buffer)
	fmt.Println(writeBt, err)

	select {}
}
