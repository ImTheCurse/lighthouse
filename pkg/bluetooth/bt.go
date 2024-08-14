package bluetooth

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
	"syscall"
)

type BTDevice struct {
	id  uuid.UUID
	buf []byte
	sd  int
}

// TODO: add windows support.
func NewBTDevice() (*BTDevice, error) {
	sd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		fmt.Errorf("Unable to create Socket. Error: %v", err)
		return nil, err
	}
	device := &BTDevice{uuid.New(), make([]byte, 1024), sd}

	return device, nil
}

func (device *BTDevice) Bind(localMacAddress string) {
	unix.Bind(device.sd, &unix.SockaddrRFCOMM{
		Channel: 1,
		Addr:    str2ba(localMacAddress),
	})
}

func (device *BTDevice) Listen() error {
	err := unix.Listen(device.sd, 1)
	if err != nil {
		fmt.Errorf("Failed to listen to socket. Error: %v", err)
		return err
	}
	return nil
}

func (device *BTDevice) Accept() {
	nfd, sa, _ := unix.Accept(device.sd)
	fmt.Printf("conn addr=%v fd=%d", sa.(*unix.SockaddrRFCOMM).Addr, nfd)
	unix.Read(nfd, device.buf)
}

func (device *BTDevice) Connect(to string) error {
	err := unix.Connect(device.sd, &unix.SockaddrRFCOMM{
		Channel: 1,
		Addr:    str2ba(to),
	})
	if err != nil {
		fmt.Errorf("Failed to connect to Address: %v Error: %v", to, err)
		return err
	}
	return nil
}

func (device *BTDevice) Send(data []byte) error {
	_, err := unix.Write(device.sd, data)
	if err != nil {
		fmt.Errorf("Failed to send data, Error: %v", err)
		return err
	}
	return nil
}

// str2ba converts MAC address string representation to little-endian byte array
func str2ba(addr string) [6]byte {
	a := strings.Split(addr, ":")
	var b [6]byte
	for i, tmp := range a {
		u, _ := strconv.ParseUint(tmp, 16, 8)
		b[len(b)-1-i] = byte(u)
	}
	return b
}
