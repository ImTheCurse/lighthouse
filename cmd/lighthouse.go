package main

import (
	"fmt"
	"github.com/ImTheCurse/lighthouse/pkg/bluetooth"
)

func main() {

	devices, _ := bluetooth.Scan()
	for i := range devices {
		fmt.Println(devices[i].GetAddress())
	}

}
