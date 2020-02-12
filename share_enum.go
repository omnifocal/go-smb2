package smb2

import "errors"
import "fmt"
import . "github.com/omnifocal/go-smb2/internal/smb2"

type ShareInfo interface {
	Name() string
}

type shareInfo0 struct {
	shi0Netname string
}

func (shi *shareInfo0) Name() string {
	return shi.shi0Netname
}

func (fs *RemoteFileSystem) NetShareEnumAll(srvUnc string, level uint32) ShareInfo {
	// Only level 0 implemented currently
	if level != 0 {
		panic(errors.New("Invalid NetShareEnun level"))
	}

	createReq := &CreateRequest{
		SecurityFlags:        0,
		RequestedOplockLevel: SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel:   Impersonation,
		SmbCreateFlags:       0,
		DesiredAccess:        0x0012019f,
		FileAttributes:       0,
		ShareAccess:          0x00000007,
		CreateDisposition:    1,
		CreateOptions:        0,
	}
	srvsvc, err := fs.createFile("srvsvc", createReq, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(srvsvc.fd)

}
