package nexproto

import (
	"errors"
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// SecureProtocolID is the protocol ID for the Secure Connection protocol
	JsonProtocolID = 0x75

	JsonRequest  = 0x1
	JsonRequest2 = 0x2 // ?
)

// JsonProtocol handles the Json requests
type JsonProtocol struct {
	server              *nex.Server
	ConnectionIDCounter *nex.Counter
	JSONRequestHandler  func(err error, client *nex.Client, callID uint32, rawJson string)
	JSONRequest2Handler func(err error, client *nex.Client, callID uint32, rawJson string)
}

// Setup initializes the protocol
func (jsonProtocol *JsonProtocol) Setup() {
	nexServer := jsonProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if JsonProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case JsonRequest:
				go jsonProtocol.handleRequest(packet)
			case JsonRequest2:
				go jsonProtocol.handleRequest2(packet)
			default:
				log.Printf("Unsupported Secure method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// Register sets the Register handler function
func (jsonProtocol *JsonProtocol) JSONRequest(handler func(err error, client *nex.Client, callID uint32, rawJson string)) {
	jsonProtocol.JSONRequestHandler = handler
}

func (jsonProtocol *JsonProtocol) JSONRequest2(handler func(err error, client *nex.Client, callID uint32, rawJson string)) {
	jsonProtocol.JSONRequest2Handler = handler
}

func (jsonProtocol *JsonProtocol) handleRequest(packet nex.PacketInterface) {
	if jsonProtocol.JSONRequestHandler == nil {
		log.Println("[Warning] JsonProtocol::JSONRequest not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, jsonProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[JsonProtocol::JSONRequest] Json missing length")
		go jsonProtocol.JSONRequestHandler(err, client, callID, "")
		return
	}

	rawJson, err := parametersStream.Read4ByteString()

	if err != nil {
		go jsonProtocol.JSONRequestHandler(err, client, callID, "")
		return
	}

	go jsonProtocol.JSONRequestHandler(nil, client, callID, rawJson)
}

func (jsonProtocol *JsonProtocol) handleRequest2(packet nex.PacketInterface) {
	if jsonProtocol.JSONRequestHandler == nil {
		log.Println("[Warning] JsonProtocol::JSONRequest2 not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, jsonProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[JsonProtocol::JSONRequest2] Json missing length")
		go jsonProtocol.JSONRequest2Handler(err, client, callID, "")
		return
	}

	rawJson, err := parametersStream.Read4ByteString()

	if err != nil {
		go jsonProtocol.JSONRequest2Handler(nil, client, callID, "[]")
		return
	}

	go jsonProtocol.JSONRequest2Handler(nil, client, callID, rawJson)
}

// NewSecureProtocol returns a new SecureProtocol
func NewJsonProtocol(server *nex.Server) *JsonProtocol {
	jsonProtocol := &JsonProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	jsonProtocol.Setup()

	return jsonProtocol
}
