package mesh

import (
	"encoding/binary"
	"fmt"

	ble "github.com/ImTheCurse/lighthouse/pkg/bluetooth"
	"github.com/cilium/cilium/pkg/mac"
	"tinygo.org/x/bluetooth"
)

type bluetoothDevice = bluetooth.ScanResult
type opcode uint16
type lock uint16
type meshBLE struct {
	ble.BLEDevice
}

const HEADER_SIZE = 20
const (
	NOTIFY opcode = iota + 1
	TRAVERSE
	SEND
	RECV
	ACK
	ERROR
)
const (
	UNLOCKED lock = iota + 1
	LOCKED
)

func NewMeshBLE(device bluetoothDevice) (*meshBLE, error) {
	dev, err := ble.NewBLEDevice(device)
	if err != nil {
		return nil, err
	}
	return &meshBLE{*dev}, nil
}

// Format buffer data in device to the format:
//
//	|  OPCODE  |  LOCK  |  TARGET ADDRESS  | BYTE PADDING |  LOCAL ADDRESS | BYTE PADDING |   DATA   |
//	   2 byte    2 byte       6 byte	    2 byte          6 bytes        2 byte        44 byte
//
// NOTE: TARGET ADDRESS and LOCAL ADDRESS are in little endian byte order, as defined in reference: https://pkg.go.dev/github.com/cilium/cilium@v1.16.1/pkg/mac#MAC.Uint64
// NOTE: data is sent from device's buffer, and data is in little endian byte order as defined in reference: https://www.bluetooth.org/DocMan/handlers/DownloadDoc.ashx?doc_id=587177
// NOTE: as defined in vol.3 chapter G, part 2.4 bluetooth core specification.
func (device *meshBLE) FormatData(code opcode, lockStatus lock, localAddress mac.Uint64MAC, targetAddress mac.Uint64MAC) error {
	buf, err := device.GetDeviceBuffer()
	if err != nil {
		return fmt.Errorf("Unable to read buffer, Error: %v", buf)
	}
	if err != nil {
		return fmt.Errorf("Unable to decode buffer. Error: %v", err)
	}

	// FIX: Hardcoded local device address, need to parse target address and pass it to sendData.
	address := device.GetAddress()
	data, err := device.RecieveData(address)
	if err != nil {
		return fmt.Errorf("Unable to format data, Error: %v", err)
	}
	header := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint16(header[:2], uint16(code))
	binary.BigEndian.PutUint16(header[2:4], uint16(lockStatus))
	binary.LittleEndian.PutUint64(header[4:12], uint64(targetAddress))
	binary.LittleEndian.PutUint64(header[12:20], uint64(localAddress))

	// NOTE: We assume that there wasn't unformatted data before hand, as we initialize the BLE devices with formatted data.
	n := copy(buf[0:20], header)
	if n != 20 {
		return fmt.Errorf("Error: Header wasen't copied to data.")
	}
	data = data[20:]
	if err != nil {
		return fmt.Errorf("Unable to format data: %v", err)
	}
	n = copy(buf[20:], data)
	fmt.Println()
	fmt.Printf("Buffer length: %v", len(buf))
	fmt.Println()
	fmt.Printf("Data Length: %v", len(data))
	fmt.Printf("\nData: %x", data)
	fmt.Printf("\nBuf: %x", buf)
	ble.SendData(address, buf)
	return nil
}

// TODO: initialize ble device with a handler for each opcode.
func (device *meshBLE) SendFormattedData(localAddress mac.Uint64MAC, targetAddress mac.Uint64MAC) error {
	err := device.FormatData(SEND, UNLOCKED, localAddress, targetAddress)
	buf, _ := device.GetDeviceBuffer()
	fmt.Printf("From SendFormat: %x", buf)
	if err != nil {
		return err
	}
	neighbors := device.BLEDevice.GetNeighbors()
	if neighbors == nil {
		neighbors, err = ble.Scan()
		if err != nil {
			return err
		}
	}
	for _, nei := range neighbors {
		buf, err := device.GetDeviceBuffer()
		if err != nil {
			return fmt.Errorf("Unable to send formatted data. Error: %v", err)
		}
		ble.SendData(nei.Address, buf)
	}
	return nil
}
