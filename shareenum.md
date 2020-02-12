# SMB2 share enumeration request design notes

- Use go-smb2 for session and encapsulation
- Maybe github.com/c-sto/wmiexec for DCERPC?

## Process flow
1. Instantiate smb2.Dialer and dial 445 on target
2. TreeConnect to `\\<TARGET>\IPC$`
3. CreateRequest for `srvsvc`
    - No oplock
    - ImpersonationLevel `2`
    - CreateFlags `0x0`
    - DesiredAccess `0x0012019f`
    - FileAttributes: `0x0`
    - ShareAccess `0x00000007`
        - shared for read/write/delete
    - Disposition: Open `1`
    - CreateOptions: `0x00000000`
    - Filename (UTF-16): `srvsvc`
    - No CreateContexts
4. Handle CreateResponse
5. GetInfo request for the file handle
    - Class: FILE_INFO (`0x01`)
    - InfoLevel: SMB2_FILE_STANDARD_INFO (`0x05`)
6. Handle GetInfo Response
7. DCE/RPC Bind call to SRVSVC
    - SMB WriteRequest (`0x09`)
        - Channel: None (`0x0`)
        - Write Flags: `0x0`
    - DCE/RPC Bind
        - Packet type: Bind (`0x0b`)
        - Data Repr: `0x10000000`
            - Little endian
            - ASCII
            - IEEE
        - Call ID: 2
        - Assoc Group: `0x0`
        - Num Ctx Items: 3
        - Ctx Item[1]: Context ID:0, SRVSVC, 32bit NDR
            - Context ID: 0
            - Num Trans Items: 1
            - Abstract Syntax: SRVSVC V3.0
                - Interface: SRVSVC UUID: 4b324fc8-1670-01d3-1278-5a47bf6ee188
                - Interface Ver: 3
                - Interface Ver Minor: 0
            - Transfer Syntax[1]: 32bit NDR V2
                - Transfer Syntax: 32bit NDR UUID:8a885d04-1ceb-11c9-9fe8-08002b104860
                - ver: 2
        - Ctx Item[2]: Context ID:1, SRVSVC, 64bit NDR
            - Context ID: 1
            - Num Trans Items: 1
            - Abstract Syntax: SRVSVC V3.0
                - Interface: SRVSVC UUID: 4b324fc8-1670-01d3-1278-5a47bf6ee188
                - Interface Ver: 3
                - Interface Ver Minor: 0
            - Transfer Syntax[1]: 64bit NDR V1
                - Transfer Syntax: 64bit NDR UUID:71710533-beba-4937-8319-b5dbef9ccc36
                - ver: 1
        - Ctx Item[3]: Context ID:2, SRVSVC, Bind Time Feature Negotiation
            - Context ID: 2
            - Num Trans Items: 1
            - Abstract Syntax: SRVSVC V3.0
                - Interface: SRVSVC UUID: 4b324fc8-1670-01d3-1278-5a47bf6ee188
                - Interface Ver: 3
                - Interface Ver Minor: 0
            - Transfer Syntax[1]: Bind Time Feature Negotiation V1
                Transfer Syntax: Bind Time Feature Negotiation - UUID:6cb71c2c-9812-4540-0300-000000000000
                Bind Time Features: 0x0003, Security Context Multiplexing Supported, Keep - Connection On Orphan Supported
                - ver: 1
8. Handle SMB WriteResponse
9. Send SMB ReadRequest for srvsvc
    - Reserved: `0x5000`
    - Read Length: 1024
    - Offset: 0
    - Channel: None
10. Handle ReadResponse (containing DCERPC Bind_ack)
11. Send NetShareEnumAll RPC request
    - Context ID: 1
    - Opnum: 15
    - Server Service, NetShareEnumAll
        - Operation: NetShareEnumAll (15)
        - [Response in frame: 55]
        - Pointer to Server Unc (uint16)
        - Pointer to Level (uint32)
            - Level: 0
        - Pointer to Ctr (srvsvc_NetShareCtr)
            - srvsvc_NetShareCtr
                - Ctr
                - NULL Pointer: Pointer to Ctr0 (srvsvc_NetShareCtr0)
        - Max Buffer: 131072
        - NULL Pointer: Pointer to Resume Handle (uint32)
        - Long frame
12. Handle NetShareEnumAll response
13. Close srvsvc
14. Tree disconnect
15. Session logoff
16. Kill connection
