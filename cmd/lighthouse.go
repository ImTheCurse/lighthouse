package main

import (
	//"encoding/binary"

	"github.com/ImTheCurse/lighthouse/pkg/bluetooth"
)

func main() {
	bluetooth.InitServer()
	bluetooth.StartServer()
	/*
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
	*/

}
