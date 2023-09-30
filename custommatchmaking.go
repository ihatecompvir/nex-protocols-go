package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	CustomMatchmakingProtocolID = 0x6E

	CustomFind = 0x1
)

// JsonProtocol handles the Json requests
type CustomMatchmakingProtocol struct {
	server              *nex.Server
	ConnectionIDCounter *nex.Counter
	CustomFindHandler   func(err error, client *nex.Client, callID uint32, data []byte)
}

func (customMatchmakingProtocol *CustomMatchmakingProtocol) Setup() {
	nexServer := customMatchmakingProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if CustomMatchmakingProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go customMatchmakingProtocol.handleCustomFind(packet)
			default:
				log.Printf("Unsupported CustomMatchmaking method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (customMatchmakingProtocol *CustomMatchmakingProtocol) CustomFind(handler func(err error, client *nex.Client, callID uint32, data []byte)) {
	customMatchmakingProtocol.CustomFindHandler = handler
}

func (customMatchmakingProtocol *CustomMatchmakingProtocol) handleCustomFind(packet nex.PacketInterface) {
	if customMatchmakingProtocol.CustomFindHandler == nil {
		log.Println("[Warning] CustomMatchmakingProtocol::CustomFind not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	go customMatchmakingProtocol.CustomFindHandler(nil, client, callID, parameters)
}

// NewCustomMatchmakingProtocol returns a new CustomMatchmakingProtocol
func NewCustomMatchmakingProtocol(server *nex.Server) *CustomMatchmakingProtocol {
	customMatchmakingProtocol := &CustomMatchmakingProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	customMatchmakingProtocol.Setup()

	return customMatchmakingProtocol
}
