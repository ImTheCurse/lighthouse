package main

import (
	"encoding/binary"
	"fmt"
	"github.com/ImTheCurse/lighthouse/pkg/bluetooth"
)

func main() {

	mac := "48:E7:29:9F:76:8A"

	device, err := bluetooth.ScanForDevice(mac)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	ble := bluetooth.NewBLEDevice(*device)

	data, _ := ble.RecieveData(device.Address)
	fmt.Println(binary.LittleEndian.Uint32(data[:]))

	bluetooth.SendData(device.Address, []byte("hello world!"))
	data, _ = ble.RecieveData(device.Address)
	fmt.Println(string(data[:]))

}
