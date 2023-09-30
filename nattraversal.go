package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	NATTraversalProtocolID = 0x3 // the first matchmaking service protocol

	RequestProbeInitiation = 0x1
	InitiateProbe          = 0x2
)

// JsonProtocol handles the Json requests
type NATTraversalProtocol struct {
	server                        *nex.Server
	ConnectionIDCounter           *nex.Counter
	RequestProbeInitiationHandler func(err error, client *nex.Client, callID uint32, stationURLs []string)
}

func (natTraversalProtocol *NATTraversalProtocol) Setup() {
	nexServer := natTraversalProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if NATTraversalProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go natTraversalProtocol.handleRequestProbeInitiation(packet)
			default:
				log.Printf("Unsupported NAT traversal method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (natTraversalProtocol *NATTraversalProtocol) RequestProbeInitiation(handler func(err error, client *nex.Client, callID uint32, stationURLs []string)) {
	natTraversalProtocol.RequestProbeInitiationHandler = handler
}

func (natTraversalProtocol *NATTraversalProtocol) handleRequestProbeInitiation(packet nex.PacketInterface) {
	if natTraversalProtocol.RequestProbeInitiationHandler == nil {
		log.Println("[Warning] NATTraversal::RequestProbeInitiation not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, natTraversalProtocol.server)

	numStationURLs := parametersStream.ReadUInt32LE()
	urlSlice := make([]string, numStationURLs)

	for i := 0; i < int(numStationURLs); i++ {
		url, err := parametersStream.Read4ByteString()

		if err != nil {
			go natTraversalProtocol.RequestProbeInitiationHandler(nil, client, callID, make([]string, 0))
			return
		}

		urlSlice[i] = url
	}

	go natTraversalProtocol.RequestProbeInitiationHandler(nil, client, callID, urlSlice)
}

// NewSecureProtocol returns a new SecureProtocol
func NewNATTraversalProtocol(server *nex.Server) *NATTraversalProtocol {
	natTraversalProtocol := &NATTraversalProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	natTraversalProtocol.Setup()

	return natTraversalProtocol
}
