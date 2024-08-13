package bluetooth

import (
	"fmt"
	"time"

	"github.com/amit7itz/goset"
	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var markedDevices = make(map[bluetooth.Address]bool)

type bluetoothDevice = bluetooth.ScanResult
type Sender interface {
	New() error
	GetUUID() uuid.UUID
	Scan() ([]*BLEDevice, error)
	SendData(to bluetooth.Address, msg []byte) error
	GetAddress() bluetooth.Address
}

type Reciever interface {
	New() error
	GetUUID() (uuid.UUID, error)
	GetAddress() bluetooth.Address
	RecieveData(to bluetooth.Address, msg []byte) ([]byte, error)
}

type BLEDevice struct {
	id        uuid.UUID
	neighbors []*BLEDevice
	device    bluetoothDevice
}

func (ble BLEDevice) GetUUID() uuid.UUID {
	return ble.id
}

func (ble BLEDevice) GetAddress() bluetooth.Address {
	return ble.device.Address
}

func (ble BLEDevice) getNeighbors() []*BLEDevice {
	return ble.neighbors
}

func NewBLEDevice(device bluetoothDevice) *BLEDevice {
	return &BLEDevice{uuid.New(), nil, device}
}

func Scan() ([]bluetoothDevice, error) {
	err := adapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE stack:" + err.Error())
	}

	BLEdevicesSet := goset.NewSet[*BLEDevice]()
	deviceSet := goset.NewSet[bluetoothDevice]()

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

		// TODO: add node traversal from current device to next device recursively
		newDevice := &BLEDevice{uuid.New(), nil, device}
		BLEdevicesSet.Add(newDevice)
		deviceSet.Add(device)

	}

	err = adapter.Scan(handleScan)
	if err != nil {
		fmt.Println("Failed to start scan:", err.Error())
	}
	return deviceSet.Items(), nil
}

func ScanForDevice(targetAddress string) (*bluetoothDevice, error) {
	err := adapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE stack:" + err.Error())
	}
	fmt.Println("Scanning...")

	var dev *bluetoothDevice = nil

	handleScan := func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		go func() {
			time.Sleep(5 * time.Second)
			err := adapter.StopScan()
			if err != nil {
				fmt.Println("Failed to start scan:", err.Error())
			}
		}()
		if device.Address.MACAddress.String() == targetAddress {
			fmt.Println("Found Bluetooth device:", device.Address.MACAddress.String())
			dev = &device
			return
		}
	}

	err = adapter.Scan(handleScan)
	if err != nil {
		fmt.Println("Failed to start scan:", err.Error())
	}

	if dev == nil {
		fmt.Println()
		return nil, fmt.Errorf("Couldn't find bluetooth device: %v", targetAddress)
	}

	return dev, nil

}

func (ble BLEDevice) RecieveData(from bluetooth.Address) ([]byte, error) {
	device, err := adapter.Connect(from, bluetooth.ConnectionParams{})

	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return nil, err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{})
	service := services[1]

	if err != nil {

		fmt.Println("Failed to discover services.", err.Error())
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{})
	if err != nil {
		fmt.Println("Failed to discover charctristics", err.Error())
	}
	buf := make([]byte, 23)
	chars[0].Read(buf)
	return buf, nil

}

func SendData(to bluetooth.Address, msg []byte) error {
	device, err := adapter.Connect(to, bluetooth.ConnectionParams{})

	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{})
	service := services[1]

	if err != nil {

		fmt.Println("Failed to discover services.", err.Error())
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{})
	if err != nil {
		fmt.Println("Failed to discover charctristics", err.Error())
	}

	_, err = chars[0].WriteWithoutResponse(msg)
	if err != nil {
		fmt.Errorf("Write error: %v", err)
	}
	return nil
}
