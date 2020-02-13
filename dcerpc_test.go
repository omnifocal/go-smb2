package smb2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/C-Sto/goWMIExec/pkg/uuid"
)

func TestBRMarshal(t *testing.T) {
	intf, err := uuid.FromString("4b324fc8-1670-01d3-1278-5a47bf6ee188")
	intfVer := uint16(3)
	intfVerMinor := uint16(0)
	ndr32, err := uuid.FromString(UUID_32BIT_NDR_V2)
	ndr64, err := uuid.FromString(UUID_64BIT_NDR_V1)
	bindNego, err := uuid.FromString(UUID_BINDTIME_FEATURENEGO)
	if err != nil {
		t.Error(err)
	}

	req := &bindRequest{
		Version:      5,
		MinorVersion: 0,
		PacketType:   11,         // Bind
		PacketFlags:  0x03,       // First frag and last frag set
		DataRepr:     0x10000000, // Little endian ASCII, IEEE float
		FragLen:      160,
		AuthLen:      0,
		CallId:       2, // Don't know what this is about
		MaxXmit:      4280,
		MaxRecv:      4280,
		AssocGroup:   0,
		NumCtxItems:  3,
		ctxItems: []ctxItem{
			ctxItem{
				ContextId:         0,
				NumTransItems:     1,
				InterfaceUuid:     intf,
				InterfaceVer:      intfVer,
				InterfaceVerMinor: intfVerMinor,
				transItems: []transItem{
					transItem{
						TransSyntax: ndr32,
						Ver:         2,
					},
				},
			},
			ctxItem{
				ContextId:         1,
				NumTransItems:     1,
				InterfaceUuid:     intf,
				InterfaceVer:      intfVer,
				InterfaceVerMinor: intfVerMinor,
				transItems: []transItem{
					transItem{
						TransSyntax: ndr64,
						Ver:         1,
					},
				},
			},
			ctxItem{
				ContextId:         2,
				NumTransItems:     1,
				InterfaceUuid:     intf,
				InterfaceVer:      intfVer,
				InterfaceVerMinor: intfVerMinor,
				transItems: []transItem{
					transItem{
						TransSyntax: bindNego,
						Ver:         1,
					},
				},
			},
		},
	}

	want, err := hex.DecodeString("05000b0310000000a000000002000000b810b810000000000300000000000100c84f324b7016d30112785a47bf6ee18803000000045d888aeb1cc9119fe808002b1048600200000001000100c84f324b7016d30112785a47bf6ee1880300000033057171babe37498319b5dbef9ccc360100000002000100c84f324b7016d30112785a47bf6ee188030000002c1cb76c12984045030000000000000001000000")
	if err != nil {
		t.Error(err)
	}
	got := req.marshalLE()
	t.Log(hex.Dump(got))
	if bytes.Compare(want, got) != 0 {
		t.Fail()
	}
}
