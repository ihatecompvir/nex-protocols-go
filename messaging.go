package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	MessagingProtocolID = 0x17

	GetMessageHeaders = 0x3
)

type MessagingProtocol struct {
	server                   *nex.Server
	ConnectionIDCounter      *nex.Counter
	GetMessageHeadersHandler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, rangeOffset uint32, rangeSize uint32)
}

func (unknownProtocol *MessagingProtocol) Setup() {
	nexServer := unknownProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if MessagingProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case GetMessageHeaders:
				go unknownProtocol.handleGetMessageHeaders(packet)
			default:
				log.Printf("Unsupported Messaging method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (messagingProtocol *MessagingProtocol) GetMessageHeaders(handler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, rangeOffset uint32, rangeSize uint32)) {
	messagingProtocol.GetMessageHeadersHandler = handler
}

func (messagingProtocol *MessagingProtocol) handleGetMessageHeaders(packet nex.PacketInterface) {
	if messagingProtocol.GetMessageHeadersHandler == nil {
		log.Println("[Warning] MessagingProtocol::GetMessageHeadersHandler not implemented")
		go respondNotImplemented(packet, MessagingProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, messagingProtocol.server)

	pid := parametersStream.ReadUInt32LE()
	recipientType := parametersStream.ReadUInt32LE() // 1 = PID,  2 = gathering ID
	rangeOffset := parametersStream.ReadUInt32LE()
	rangeSize := parametersStream.ReadUInt32LE()

	go messagingProtocol.GetMessageHeadersHandler(nil, client, callID, pid, recipientType, rangeOffset, rangeSize)
}

// NewSecureProtocol returns a new SecureProtocol
func NewMessagingProtocol(server *nex.Server) *MessagingProtocol {
	messagingProtocol := &MessagingProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	messagingProtocol.Setup()

	return messagingProtocol
}
