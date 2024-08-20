package main

import (
	//	"encoding/binary"
	"fmt"

	"github.com/ImTheCurse/lighthouse/pkg/bluetooth"
	"github.com/ImTheCurse/lighthouse/pkg/mesh"
	"github.com/cilium/cilium/pkg/mac"
)

func main() {

	macAdd := "48:E7:29:9F:76:8A"
	localMac := "04:7F:0E:3C:73:69"
	localAdd, _ := mac.ParseMAC(localMac)
	loc, _ := localAdd.Uint64()
	remoteAdd, _ := mac.ParseMAC(macAdd)
	rem, _ := remoteAdd.Uint64()
	device, err := bluetooth.ScanForDevice(macAdd)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
	ble, _ := mesh.NewMeshBLE(*device)
	ble.FormatData(mesh.SEND, mesh.UNLOCKED, loc, rem)

	buf, _ := ble.RecieveData(ble.GetAddress())

	fmt.Printf("Buffer: % x", buf)

	/*

		ble, _ := bluetooth.NewBLEDevice(*device)
		err = ble.WriteDataToLocalBuffer([]byte("we cooking?"))
		if err != nil {
			fmt.Println(err)
		}

		checkdata, _ := ble.RecieveData(device.Address)
		fmt.Println(string(checkdata[:]))

		bluetooth.SendData(device.Address, []byte("hello world!"))
		data, _ := ble.RecieveData(device.Address)
		fmt.Println(string(data[:]))
	*/
}
