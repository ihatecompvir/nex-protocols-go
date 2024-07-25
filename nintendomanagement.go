package nexproto

import (
	"fmt"
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	NintendoManagementProtocolID = 0x53

	GetConsoleUsernames = 2 // returns a list of users on a particular console by its friend code
)

type NintendoManagementProtocol struct {
	server                     *nex.Server
	GetConsoleUsernamesHandler func(err error, client *nex.Client, callID uint32, friendCode string)
}

func reverseBytes(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	output := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		output[i] = input[len(input)-1-i]
	}
	return output
}

func bytesToUint64(input []byte) uint64 {
	var result uint64
	for i := 0; i < len(input); i++ {
		result |= uint64(input[i]) << (8 * i)
	}
	return result
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
func (nintendoManagementProtocol *NintendoManagementProtocol) GetConsoleUsernames(handler func(err error, client *nex.Client, callID uint32, friendCode string)) {
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
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, nintendoManagementProtocol.server)

	friendCode := make([]byte, 7)

	for i := 0; i < 7; i++ {
		friendCode[i] = parametersStream.ReadUInt8()
	}

	finalFriendCode := fmt.Sprintf("%d", bytesToUint64(reverseBytes(friendCode)))

	go nintendoManagementProtocol.GetConsoleUsernamesHandler(nil, client, callID, finalFriendCode)
}

// NewRBBinaryDataProtocol returns a new RBBinaryDataProtocol
func NewNintendoManagementProtocol(server *nex.Server) *NintendoManagementProtocol {
	nintendoManagementProtocol := &NintendoManagementProtocol{server: server}

	nintendoManagementProtocol.Setup()

	return nintendoManagementProtocol
}
