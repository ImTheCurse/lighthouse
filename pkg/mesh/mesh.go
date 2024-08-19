package mesh

import (
	"encoding/binary"
	"fmt"
	ble "github.com/ImTheCurse/lighthouse/pkg/bluetooth"
	"github.com/cilium/cilium/pkg/mac"
)

type opcode uint16
type lock uint16

const HEADER_SIZE = 28
const (
	NOTIFY opcode = iota
	TRAVERSE
	SEND
	RECV
	ACK
	ERROR
)
const (
	UNLOCKED lock = iota
	LOCKED
)

// Format buffer data in device to the format:
//
//	|  SEQ. NUM  |  OPCODE  |  LOCK  |  TARGET ADDRESS  | LOCAL ADDRESS |  DATA   |
//	   8 byte       2 byte    2 byte       8 byte            8 bytes     494 bytes
//
// NOTE: TARGET ADDRESS and LOCAL ADDRESS are in little endian byte order, as defined in reference: https://pkg.go.dev/github.com/cilium/cilium@v1.16.1/pkg/mac#MAC.Uint64
// NOTE: data is sent from device's buffer.
func formatData(device *ble.BLEDevice, seqNum uint64, code opcode, lockStatus lock, localAddress mac.Uint64MAC, targetAddress mac.Uint64MAC) error {
	buf, err := device.GetDeviceBuffer()
	if err != nil {
		return fmt.Errorf("Unable to format data, Error: %v", err)
	}
	header := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint64(header[:8], seqNum)
	binary.BigEndian.PutUint16(header[8:10], uint16(code))
	binary.BigEndian.PutUint16(header[10:12], uint16(lockStatus))
	binary.LittleEndian.PutUint64(header[12:20], uint64(targetAddress))
	binary.LittleEndian.PutUint64(header[20:28], uint64(localAddress))

	// NOTE: We assume that there wasn't unformatted data before hand, as we initialize the BLE devices with formatted data.
	n := copy(buf[:28], header)
	if n != 28 {
		return fmt.Errorf("Header wasen't copied to data. Error")
	}
	return nil
}
