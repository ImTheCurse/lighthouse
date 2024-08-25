package bluetooth

import (
	"fmt"
	"time"

	"github.com/amit7itz/goset"
	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

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
	neighbors []bluetoothDevice
	device    bluetoothDevice
}

// Get device mesh network id.
func (ble BLEDevice) GetUUID() uuid.UUID {
	return ble.id
}

// Get local device MAC address.
func (ble BLEDevice) GetAddress() bluetooth.Address {
	return ble.device.Address
}

// Get all previously scanned neighbors.
func (ble BLEDevice) GetNeighbors() []bluetoothDevice {
	return ble.neighbors
}

// Construct new BLEDevice.
func NewBLEDevice(device bluetoothDevice) (*BLEDevice, error) {
	nei, err := Scan()
	if err != nil {
		return nil, fmt.Errorf("Unable to create BLE device. Error: %v", err)
	}
	return &BLEDevice{uuid.New(), nei, device}, nil
}

// Scan for all of the devices in the physical area.
func Scan() ([]bluetoothDevice, error) {
	err := adapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE stack:" + err.Error())
	}

	deviceSet := goset.NewSet[bluetoothDevice]()

	fmt.Println("Scanning...")
	handleScan := func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		go func() {
			// FIX: hardcoded to 5 seconds for debugging purposes, to be removed.
			time.Sleep(5 * time.Second)
			err := adapter.StopScan()
			if err != nil {
				fmt.Println("Failed to start scan:", err.Error())
			}
		}()
		fmt.Println("Found device:", device.Address.String(), device.RSSI, device.LocalName())
		deviceSet.Add(device)

	}

	err = adapter.Scan(handleScan)
	if err != nil {
		fmt.Println("Failed to start scan:", err.Error())
	}
	return deviceSet.Items(), nil
}

// Scan for a specific device in the physical area.
func ScanForDevice(targetAddress string) (*bluetoothDevice, error) {
	err := adapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE stack:" + err.Error())
	}
	fmt.Println("Scanning...")

	var dev *bluetoothDevice = nil

	handleScan := func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		go func() {
			// FIX: hardcoded to 5 seconds for debugging purposes, to be removed.
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

// Get device default charctristic buffer
func (ble BLEDevice) GetDeviceBuffer() ([]byte, error) {
	return ble.RecieveData(ble.device.Address)
}

// Read data from a given target ble device using its mac address.
func (ble BLEDevice) RecieveData(from bluetooth.Address) ([]byte, error) {
	device, err := adapter.Connect(from, bluetooth.ConnectionParams{})
	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return nil, err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{})
	if len(services) == 0 {
		return make([]byte, 0), fmt.Errorf("Error: found no services.")
	}
	service := services[1]

	if err != nil {

		fmt.Println("Failed to discover services.", err.Error())
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{})
	if err != nil {
		fmt.Println("Failed to discover charctristics", err.Error())
	}
	buf := make([]byte, 64)
	chars[0].Read(buf)
	return buf, nil

}

// Write data to target ble device. Write directly to device, with no data format.
func SendData(to bluetooth.Address, msg []byte) error {
	device, err := adapter.Connect(to, bluetooth.ConnectionParams{})
	fmt.Println()
	fmt.Printf("Sent message: %x\n", msg)
	//defer device.Disconnect()
	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{})
	if len(services) == 0 {
		return fmt.Errorf("Error: No services available.")
	}
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
