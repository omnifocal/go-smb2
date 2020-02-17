// +build integration

package smb2

import (
	"fmt"
	"net"
	"os"

	"testing"

	"github.com/C-Sto/goWMIExec/pkg/uuid"
)

var (
	host = os.Getenv("SMB_TEST_HOST")
	user = os.Getenv("SMB_TEST_USER")
	pass = os.Getenv("SMB_TEST_PASS")
)

func TestRpcBind(t *testing.T) {
	share := "IPC$"

	d := &Dialer{
		Initiator: &NTLMInitiator{
			User:     user,
			Password: pass,
			Domain:   ".",
		},
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:445", host))
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	c, err := d.Dial(conn)
	if err != nil {
		t.Error(err)
	}
	defer c.Logoff()

	fs, err := c.Mount(fmt.Sprintf(`\\%s\%s`, host, share))
	if err != nil {
		t.Error(err)
	}
	defer fs.Umount()

	intf, err := uuid.FromString("4b324fc8-1670-01d3-1278-5a47bf6ee188")
	intfVer := uint16(3)
	intfVerMinor := uint16(0)

	file, err := fs.OpenPipe("srvsvc")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	file.rpcBind(intf, intfVer, intfVerMinor)

	return
}
