package main

import (
	"fmt"
	"usbMonitor/usbLib"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	testDev := usbLib.NewStm32Manager()

	stmUsbDev, err := usbLib.FindUsbOne(testDev)
	if err != nil {
		fmt.Println("find error :", err)
	}
	defer stmUsbDev.Close()

	select {}
}
