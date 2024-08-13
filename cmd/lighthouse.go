package main

import (
	"fmt"

	"github.com/ImTheCurse/lighthouse/pkg/bluetooth"
)

func main() {

	mac := "48:E7:29:9F:76:8A"

	device, err := bluetooth.ScanForDevice(mac)

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	fmt.Println(device.LocalName())

}
