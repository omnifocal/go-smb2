package smb2

// import . "github.com/omnifocal/go-smb2/internal/smb2"

import (
	"bytes"
	"encoding/binary"

	// Thank you C-Sto for your suffering
	"github.com/C-Sto/goWMIExec/pkg/uuid"
)

const (
	UUID_32BIT_NDR_V2         = "8a885d04-1ceb-11c9-9fe8-08002b104860"
	UUID_64BIT_NDR_V1         = "71710533-beba-4937-8319-b5dbef9ccc36"
	UUID_BINDTIME_FEATURENEGO = "6cb71c2c-9812-4540-0300-000000000000"
)

type ctxItem struct {
	ContextId         uint16
	NumTransItems     uint16
	InterfaceUuid     uuid.UUID
	InterfaceVer      uint16
	InterfaceVerMinor uint16

	transItems []transItem
}

type transItem struct {
	TransSyntax uuid.UUID
	Ver         uint32
}

type bindRequest struct {
	Version      byte
	MinorVersion byte
	PacketType   byte
	PacketFlags  byte
	DataRepr     uint32
	FragLen      uint16
	AuthLen      uint16
	CallId       uint32
	MaxXmit      uint16
	MaxRecv      uint16
	AssocGroup   uint32
	NumCtxItems  uint32

	ctxItems []ctxItem
}

func (ci *ctxItem) marshalLE() []byte {
	var err error
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, ci.ContextId)
	err = binary.Write(buf, binary.LittleEndian, ci.NumTransItems)
	err = binary.Write(buf, binary.LittleEndian, ci.InterfaceUuid) // Byte array gets written in array order
	err = binary.Write(buf, binary.LittleEndian, ci.InterfaceVer)
	err = binary.Write(buf, binary.LittleEndian, ci.InterfaceVerMinor)
	if err != nil {
		panic(err)
	}

	out := buf.Bytes()

	for _, v := range ci.transItems {
		out = append(out, v.marshalLE()...)
	}

	return out
}

func (ti *transItem) marshalLE() []byte {
	var err error
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, ti.TransSyntax)
	err = binary.Write(buf, binary.LittleEndian, ti.Ver)
	if err != nil {
		panic(err)
	}

	out := buf.Bytes()

	return out
}

func (br *bindRequest) marshalLE() []byte {
	var err error
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.BigEndian, br.Version)
	err = binary.Write(buf, binary.BigEndian, br.MinorVersion)
	err = binary.Write(buf, binary.BigEndian, br.PacketType)
	err = binary.Write(buf, binary.BigEndian, br.PacketFlags)
	err = binary.Write(buf, binary.BigEndian, br.DataRepr)
	// Little endian after DataRepr
	err = binary.Write(buf, binary.LittleEndian, br.FragLen)
	err = binary.Write(buf, binary.LittleEndian, br.AuthLen)
	err = binary.Write(buf, binary.LittleEndian, br.CallId)
	err = binary.Write(buf, binary.LittleEndian, br.MaxXmit)
	err = binary.Write(buf, binary.LittleEndian, br.MaxRecv)
	err = binary.Write(buf, binary.LittleEndian, br.AssocGroup)
	err = binary.Write(buf, binary.LittleEndian, br.NumCtxItems)
	if err != nil {
		panic(err)
	}

	out := buf.Bytes()

	for _, v := range br.ctxItems {
		out = append(out, v.marshalLE()...)
	}

	return out
}

// func (fd *RemoteFile) rpcBind(intf uuid.UUID, intfVer uint16, intfVerMinor uint16) (ret ctxItem, err error) {
// 	ndr32, err := uuid.FromString(UUID_32BIT_NDR_V2)
// 	ndr64, err := uuid.FromString(UUID_64BIT_NDR_V1)
// 	bindNego, err := uuid.FromString(UUID_BINDTIME_FEATURENEGO)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req := &bindRequest{
// 		Version:      5,
// 		MinorVersion: 0,
// 		PacketType:   11,         // Bind
// 		PacketFlags:  0x03,       // First frag and last frag set
// 		DataRepr:     0x10000000, // Little endian ASCII, IEEE float
// 		FragLen:      160,
// 		AuthLen:      0,
// 		CallId:       2, // Don't know what this is about
// 		MaxXmit:      4280,
// 		MaxRecv:      4280,
// 		AssocGroup:   0,
// 		NumCtxItems:  3,
// 		ctxItems: []ctxItem{
// 			ctxItem{
// 				ContextId:         0,
// 				NumTransItems:     1,
// 				InterfaceUuid:     intf,
// 				InterfaceVer:      intfVer,
// 				InterfaceVerMinor: intfVerMinor,
// 				transItems: []transItem{
// 					transItem{
// 						TransSyntax: ndr32,
// 						Ver:         2,
// 					},
// 				},
// 			},
// 			ctxItem{
// 				ContextId:         1,
// 				NumTransItems:     1,
// 				InterfaceUuid:     intf,
// 				InterfaceVer:      intfVer,
// 				InterfaceVerMinor: intfVerMinor,
// 				transItems: []transItem{
// 					transItem{
// 						TransSyntax: ndr64,
// 						Ver:         1,
// 					},
// 				},
// 			},
// 			ctxItem{
// 				ContextId:         2,
// 				NumTransItems:     1,
// 				InterfaceUuid:     intf,
// 				InterfaceVer:      intfVer,
// 				InterfaceVerMinor: intfVerMinor,
// 				transItems: []transItem{
// 					transItem{
// 						TransSyntax: bindNego,
// 						Ver:         1,
// 					},
// 				},
// 			},
// 		},
// 	}
// }
