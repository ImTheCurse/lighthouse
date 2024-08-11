package bluetooth

import (
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/google/uuid"
	"time"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

type Sender interface {
	New() error
	GetUUID() uuid.UUID
	Scan() ([]*BLEDevice, error)
	Send(to uuid.UUID)
	GetAddress() bluetooth.Address
}

type Reciever interface {
	New() error
	GetUUID() (uuid.UUID, error)
	GetAddress() bluetooth.Address
}

type BLEDevice struct {
	id         uuid.UUID
	macAddress bluetooth.Address
	payload    []byte
}

func (ble BLEDevice) GetUUID() uuid.UUID {
	return ble.id
}

func (ble BLEDevice) GetAddress() bluetooth.Address {
	return ble.macAddress
}

func Scan() ([]*BLEDevice, error) {
	err := adapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE stack:" + err.Error())
	}

	devicesSet := goset.NewSet[*BLEDevice]()

	fmt.Println("Scanning...")
	handleScan := func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		go func() {
			time.Sleep(5 * time.Second)
			err := adapter.StopScan()
			if err != nil {
				fmt.Println("Failed to start scan:", err.Error())
			}
		}()
		fmt.Println("Found device:", device.Address.String(), device.RSSI, device.LocalName())
		payload := device.AdvertisementPayload.Bytes()
		newDevice := &BLEDevice{uuid.New(), device.Address, payload}
		devicesSet.Add(newDevice)

	}

	err = adapter.Scan(handleScan)
	if err != nil {
		fmt.Println("Failed to start scan:", err.Error())
	}
	return devicesSet.Items(), nil

}
