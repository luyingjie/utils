package tun

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/songgao/water"
)

type Interface struct {
	TUN *water.Interface
}

func NewInterface() (*Interface, error) {
	iface := &Interface{}

	ifconfig := water.Config{
		DeviceType: water.TUN,
	}

	for i := 0; i < 10; i++ {
		ifconfig.Name = fmt.Sprintf("agent.%d", i)

		ifce, err := water.New(ifconfig)
		if err != nil {
			return nil, fmt.Errorf("new interface %s fail: %v", ifconfig.Name, err)
			time.Sleep(time.Second * 1)
			continue
		}

		iface.TUN = ifce
		return iface, nil
	}
	return nil, fmt.Errorf("new interface %s fail", ifconfig.Name)
}

func (iface *Interface) SetMTU(mtu int) error {
	out, err := ExecCmd("ifconfig", []string{iface.TUN.Name(), "mtu", fmt.Sprintf("%d", mtu)})
	if err != nil {
		return fmt.Errorf("set mtu fail: %s %v", out, err)
	}
	return nil
}

func (iface *Interface) Up() error {
	switch runtime.GOOS {
	case "linux":
		out, err := ExecCmd("ifconfig", []string{iface.TUN.Name(), "up"})
		if err != nil {
			return fmt.Errorf("ifconfig fail: %s %v", out, err)
		}

	default:
		return fmt.Errorf("unsupported: %s %s", runtime.GOOS, runtime.GOARCH)

	}

	return nil
}

func (iface *Interface) Read() ([]byte, error) {
	buf := make([]byte, 2048)
	n, err := iface.TUN.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func (iface *Interface) Write(buf []byte) (int, error) {
	return iface.TUN.Write(buf)
}

func (iface *Interface) Close() {
	iface.TUN.Close()
}

func ExecCmd(cmd string, args []string) (string, error) {
	b, err := exec.Command(cmd, args...).CombinedOutput()
	return string(b), err
}
