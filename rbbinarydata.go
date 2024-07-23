package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	RBBinaryDataProtocolID = 0x76

	// uploads a blob of binary data
	SaveBinaryData = 1

	// retrieves a blob of binary data
	GetBinaryData = 2
)

type RBBinaryDataProtocol struct {
	server                *nex.Server
	SaveBinaryDataHandler func(err error, client *nex.Client, callID uint32, metadata string, blob []byte)
	GetBinaryDataHandler  func(err error, client *nex.Client, callID uint32, metadata string)
}

// Setup initializes the protocol
func (rbBinaryDataProtocol *RBBinaryDataProtocol) Setup() {
	nexServer := rbBinaryDataProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if RBBinaryDataProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case SaveBinaryData:
				go rbBinaryDataProtocol.handleSaveBinaryData(packet)
			case GetBinaryData:
				go rbBinaryDataProtocol.handleGetBinaryData(packet)
			default:
				log.Printf("Unsupported RBBinaryData method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// SaveBinaryData sets the SaveBinaryData handler function
func (rbBinaryDataProtocol *RBBinaryDataProtocol) SaveBinaryData(handler func(err error, client *nex.Client, callID uint32, metadata string, blob []byte)) {
	rbBinaryDataProtocol.SaveBinaryDataHandler = handler
}

// GetBinaryData sets the GetBinaryData handler function
func (rbBinaryDataProtocol *RBBinaryDataProtocol) GetBinaryData(handler func(err error, client *nex.Client, callID uint32, metadata string)) {
	rbBinaryDataProtocol.GetBinaryDataHandler = handler
}

func (rbBinaryDataProtocol *RBBinaryDataProtocol) handleSaveBinaryData(packet nex.PacketInterface) {
	if rbBinaryDataProtocol.SaveBinaryDataHandler == nil {
		log.Println("[Warning] RBBinaryDataProtocol::SaveBinaryData not implemented")
		go respondNotImplemented(packet, RBBinaryDataProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, rbBinaryDataProtocol.server)

	metadata, err := parametersStream.Read4ByteString()
	if err != nil {
		go rbBinaryDataProtocol.SaveBinaryDataHandler(err, client, callID, "", nil)
		return
	}

	blob, err := parametersStream.ReadBuffer()
	if err != nil {
		go rbBinaryDataProtocol.SaveBinaryDataHandler(err, client, callID, "", nil)
		return
	}

	go rbBinaryDataProtocol.SaveBinaryDataHandler(nil, client, callID, metadata, blob)
}

func (rbBinaryDataProtocol *RBBinaryDataProtocol) handleGetBinaryData(packet nex.PacketInterface) {
	if rbBinaryDataProtocol.GetBinaryDataHandler == nil {
		log.Println("[Warning] RBBinaryDataProtocol::GetBinaryData not implemented")
		go respondNotImplemented(packet, RBBinaryDataProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, rbBinaryDataProtocol.server)

	metadata, err := parametersStream.Read4ByteString()
	if err != nil {
		go rbBinaryDataProtocol.GetBinaryDataHandler(err, client, callID, "")
		return
	}

	go rbBinaryDataProtocol.GetBinaryDataHandler(nil, client, callID, metadata)
}

// NewRBBinaryDataProtocol returns a new RBBinaryDataProtocol
func NewRBBinaryDataProtocol(server *nex.Server) *RBBinaryDataProtocol {
	rbBinaryDataProtocol := &RBBinaryDataProtocol{server: server}

	rbBinaryDataProtocol.Setup()

	return rbBinaryDataProtocol
}
