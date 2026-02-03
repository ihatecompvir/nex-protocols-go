package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	MessageDeliveryProtocolID = 0x1B
	DeliverMessage            = 0x1
)

type MessageDeliveryProtocol struct {
	server                *nex.Server
	ConnectionIDCounter   *nex.Counter
	DeliverMessageHandler func(err error, client *nex.Client, callID uint32, data []byte)
}

func (unknownProtocol *MessageDeliveryProtocol) Setup() {
	nexServer := unknownProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if MessageDeliveryProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case DeliverMessage:
				go unknownProtocol.handleDeliverMessage(packet)
			default:
				log.Printf("Unsupported Message Delivery method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (messageDeliveryProtocol *MessageDeliveryProtocol) DeliverMessage(handler func(err error, client *nex.Client, callID uint32, data []byte)) {
	messageDeliveryProtocol.DeliverMessageHandler = handler
}

func (messageDeliveryProtocol *MessageDeliveryProtocol) handleDeliverMessage(packet nex.PacketInterface) {
	if messageDeliveryProtocol.DeliverMessageHandler == nil {
		log.Println("[Warning] MessageDeliveryProtocol::DeliverMessageHandler not implemented")
		go respondNotImplemented(packet, MessageDeliveryProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	go messageDeliveryProtocol.DeliverMessageHandler(nil, client, callID, parameters)
}

// NewMessageDeliveryProtocol returns a new MessageDeliveryProtocol
func NewMessageDeliveryProtocol(server *nex.Server) *MessageDeliveryProtocol {
	messageDeliveryProtocol := &MessageDeliveryProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	messageDeliveryProtocol.Setup()

	return messageDeliveryProtocol
}
