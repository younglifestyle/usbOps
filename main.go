package main

import (
	"fmt"
	"usbMonitor/usbLib"

	log "github.com/sirupsen/logrus"
)

type tyUSBCmd_GetBIBInfo struct {
	bPrivCMD    uint8  //SCSI_PrivateCMD
	bOP         uint8  //控制命令码:OP_BIBInfo
	dwInfoIndex uint32 //指明需要返回的哪些信息
}

func main() {

	log.SetLevel(log.DebugLevel)

	testDev := usbLib.NewStm32Manager()

	stmUsbDev, err := usbLib.FindUsbOne(testDev)
	if err != nil {
		fmt.Println("find error :", err)
		return
	}
	defer stmUsbDev.Close()

	bytes := []byte{0xf0, 0x80, 0x02, 0, 0, 0}

	writeBt, err := stmUsbDev.Write(bytes)
	fmt.Println(writeBt, err)

	writeBt, err = stmUsbDev.Read()
	fmt.Println(writeBt, err)

	//i, err := stmUsbDev.Control(gousb.ControlOut, 0, 0, 0,
	//	bytes)
	//fmt.Println(i, err)

	select {}
}
