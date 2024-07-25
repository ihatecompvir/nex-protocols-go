package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	NintendoManagementProtocolID = 0x53

	GetConsoleUsernames = 2
)

type NintendoManagementProtocol struct {
	server                     *nex.Server
	GetConsoleUsernamesHandler func(err error, client *nex.Client, callID uint32)
}

// Setup initializes the protocol
func (nintendoManagementProtocol *NintendoManagementProtocol) Setup() {
	nexServer := nintendoManagementProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if NintendoManagementProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case GetConsoleUsernames:
				go nintendoManagementProtocol.handleGetConsoleUsernames(packet)
			default:
				log.Printf("Unsupported RBBinaryData method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// GetConsoleUsernames sets the GetConsoleUsernames handler function
func (nintendoManagementProtocol *NintendoManagementProtocol) GetConsoleUsernames(handler func(err error, client *nex.Client, callID uint32)) {
	nintendoManagementProtocol.GetConsoleUsernamesHandler = handler
}

func (nintendoManagementProtocol *NintendoManagementProtocol) handleGetConsoleUsernames(packet nex.PacketInterface) {
	if nintendoManagementProtocol.GetConsoleUsernamesHandler == nil {
		log.Println("[Warning] NintendoManagementProtocol::GetConsoleUsernames not implemented")
		go respondNotImplemented(packet, NintendoManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()

	go nintendoManagementProtocol.GetConsoleUsernamesHandler(nil, client, callID)
}

// NewRBBinaryDataProtocol returns a new RBBinaryDataProtocol
func NewNintendoManagementProtocol(server *nex.Server) *NintendoManagementProtocol {
	nintendoManagementProtocol := &NintendoManagementProtocol{server: server}

	nintendoManagementProtocol.Setup()

	return nintendoManagementProtocol
}
