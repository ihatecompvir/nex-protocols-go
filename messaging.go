package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	MessagingProtocolID = 0x17

	GetNumberOfMessages = 0x2
	GetMessageHeaders   = 0x3
	RetrieveMessages    = 0x5
	DeleteMessages      = 0x6
)

type MessagingProtocol struct {
	server                     *nex.Server
	ConnectionIDCounter        *nex.Counter
	GetNumberOfMessagesHandler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32)
	GetMessageHeadersHandler   func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, rangeOffset uint32, rangeSize uint32)
	RetrieveMessagesHandler    func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, lastMessageIDs []uint32, leaveOnServer bool)
	DeleteMessagesHandler      func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, messageIDs []uint32)
}

func (unknownProtocol *MessagingProtocol) Setup() {
	nexServer := unknownProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if MessagingProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case GetNumberOfMessages:
				go unknownProtocol.handleGetNumberOfMessages(packet)
			case GetMessageHeaders:
				go unknownProtocol.handleGetMessageHeaders(packet)
			case RetrieveMessages:
				go unknownProtocol.handleRetrieveMessages(packet)
			case DeleteMessages:
				go unknownProtocol.handleDeleteMessages(packet)
			default:
				log.Printf("Unsupported Messaging method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (messagingProtocol *MessagingProtocol) GetNumberOfMessages(handler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32)) {
	messagingProtocol.GetNumberOfMessagesHandler = handler
}

func (messagingProtocol *MessagingProtocol) GetMessageHeaders(handler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, rangeOffset uint32, rangeSize uint32)) {
	messagingProtocol.GetMessageHeadersHandler = handler
}

func (messagingProtocol *MessagingProtocol) RetrieveMessages(handler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, lastMessageIDs []uint32, leaveOnServer bool)) {
	messagingProtocol.RetrieveMessagesHandler = handler
}

func (messagingProtocol *MessagingProtocol) DeleteMessages(handler func(err error, client *nex.Client, callID uint32, pid uint32, recipientType uint32, messageIDs []uint32)) {
	messagingProtocol.DeleteMessagesHandler = handler
}

func (messagingProtocol *MessagingProtocol) handleGetNumberOfMessages(packet nex.PacketInterface) {
	if messagingProtocol.GetNumberOfMessagesHandler == nil {
		log.Println("[Warning] MessagingProtocol::GetNumberOfMessagesHandler not implemented")
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
	go messagingProtocol.GetNumberOfMessagesHandler(nil, client, callID, pid, recipientType)
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

func (messagingProtocol *MessagingProtocol) handleRetrieveMessages(packet nex.PacketInterface) {
	if messagingProtocol.RetrieveMessagesHandler == nil {
		log.Println("[Warning] MessagingProtocol::RetrieveMessagesHandler not implemented")
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
	numLastMessageIds := parametersStream.ReadUInt32LE()
	lastMessageIDs := make([]uint32, numLastMessageIds)
	for i := 0; i < int(numLastMessageIds); i++ {
		lastMessageIDs[i] = parametersStream.ReadUInt32LE()
	}
	leaveOnServer := parametersStream.ReadUInt8() != 0

	go messagingProtocol.RetrieveMessagesHandler(nil, client, callID, pid, recipientType, lastMessageIDs, leaveOnServer)
}

func (messagingProtocol *MessagingProtocol) handleDeleteMessages(packet nex.PacketInterface) {
	if messagingProtocol.DeleteMessagesHandler == nil {
		log.Println("[Warning] MessagingProtocol::DeleteMessagesHandler not implemented")
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
	numMessageIDs := parametersStream.ReadUInt32LE()
	messageIDs := make([]uint32, numMessageIDs)
	for i := 0; i < int(numMessageIDs); i++ {
		messageIDs[i] = parametersStream.ReadUInt32LE()
	}
	go messagingProtocol.DeleteMessagesHandler(nil, client, callID, pid, recipientType, messageIDs)
}

// NewMessagingProtocol returns a new MessagingProtocol
func NewMessagingProtocol(server *nex.Server) *MessagingProtocol {
	messagingProtocol := &MessagingProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	messagingProtocol.Setup()

	return messagingProtocol
}
