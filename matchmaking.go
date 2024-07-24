package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	MatchmakingProtocolID = 0x15 // the first matchmaking service protocol

	RegisterGathering   = 0x1  // registers a gathering with the server
	TerminateGathering  = 0x2  // ends a gathering
	UpdateGathering     = 0x4  // updates a gathering
	Participate         = 0xB  // used to denote that a player is in a particular session
	CancelParticipation = 0xC  // used to denote that a player is no longer in a particular session
	FindBySingleID      = 0x15 // looks up a gathering by its ID
	LaunchSession       = 0x1A // "launches" the session and makes it officially active
	SetState            = 0x1E // sets the state of a gathering (in song, etc.)
)

// JsonProtocol handles the Json requests
type MatchmakingProtocol struct {
	server                     *nex.Server
	ConnectionIDCounter        *nex.Counter
	RegisterGatheringHandler   func(err error, client *nex.Client, callID uint32, gathering []byte)
	UpdateGatheringHandler     func(err error, client *nex.Client, callID uint32, gathering []byte, gatheringID uint32)
	ParticipateHandler         func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	CancelParticipationHandler func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	LaunchSessionHandler       func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	TerminateGatheringHandler  func(err error, client *nex.Client, callID uint32, gatheringID uint32)
	SetStateHandler            func(err error, client *nex.Client, callID uint32, gatheringID uint32, state uint32)
	FindBySingleIDHandler      func(err error, client *nex.Client, callID uint32, gatheringID uint32)
}

func (matchmakingProtocol *MatchmakingProtocol) Setup() {
	nexServer := matchmakingProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if MatchmakingProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case RegisterGathering:
				go matchmakingProtocol.handleRegisterGathering(packet)
			case UpdateGathering:
				go matchmakingProtocol.handleUpdateGathering(packet)
			case Participate:
				go matchmakingProtocol.handleParticipate(packet)
			case CancelParticipation:
				go matchmakingProtocol.handleCancelParticipation(packet)
			case LaunchSession:
				go matchmakingProtocol.handleLaunchSession(packet)
			case TerminateGathering:
				go matchmakingProtocol.handleTerminateGathering(packet)
			case SetState:
				go matchmakingProtocol.handleSetState(packet)
			case FindBySingleID:
				go matchmakingProtocol.handleFindBySingleID(packet)
			default:
				log.Printf("Unsupported Matchmaking method ID: %#v\n", request.MethodID())
			}
		}
	})
}

func (matchmakingProtocol *MatchmakingProtocol) RegisterGathering(handler func(err error, client *nex.Client, callID uint32, gathering []byte)) {
	matchmakingProtocol.RegisterGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) UpdateGathering(handler func(err error, client *nex.Client, callID uint32, gathering []byte, gatheringID uint32)) {
	matchmakingProtocol.UpdateGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) Participate(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.ParticipateHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) CancelParticipation(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.CancelParticipationHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) LaunchSession(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.LaunchSessionHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) TerminateGathering(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.TerminateGatheringHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) SetState(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32, state uint32)) {
	matchmakingProtocol.SetStateHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) FindBySingleID(handler func(err error, client *nex.Client, callID uint32, gatheringID uint32)) {
	matchmakingProtocol.FindBySingleIDHandler = handler
}

func (matchmakingProtocol *MatchmakingProtocol) handleRegisterGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::RegisterGathering not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	parametersStream.Read4ByteString()
	parametersStream.ReadUInt32LE()
	gathering, err := parametersStream.ReadBuffer()

	if err != nil {
		log.Println("Could not read gathering data")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	go matchmakingProtocol.RegisterGatheringHandler(nil, client, callID, gathering)
}

func (matchmakingProtocol *MatchmakingProtocol) handleUpdateGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::UpdateGathering not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	parametersStream.Read4ByteString()
	parametersStream.ReadUInt32LE()
	gathering, err := parametersStream.ReadBuffer()

	gatheringStream := NewStreamIn(gathering, matchmakingProtocol.server)

	gatheringID := gatheringStream.ReadUInt32LE()

	if err != nil {
		log.Println("Could not read gathering data")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	go matchmakingProtocol.UpdateGatheringHandler(nil, client, callID, gathering, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleParticipate(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::Participate not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.ParticipateHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleCancelParticipation(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::CancelParticipation not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.CancelParticipationHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleLaunchSession(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::LaunchSession not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.LaunchSessionHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleTerminateGathering(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::TerminateGathering not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.TerminateGatheringHandler(nil, client, callID, gatheringID)
}

func (matchmakingProtocol *MatchmakingProtocol) handleSetState(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::SetState not implemented")
		go respondNotImplemented(packet, MatchmakingProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()
	state := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.SetStateHandler(nil, client, callID, gatheringID, state)

}

func (matchmakingProtocol *MatchmakingProtocol) handleFindBySingleID(packet nex.PacketInterface) {
	if matchmakingProtocol.RegisterGatheringHandler == nil {
		log.Println("[Warning] MatchmakingProtocol::FindBySingleIDs not implemented")
		go respondNotImplemented(packet, MatchmakingProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, matchmakingProtocol.server)

	gatheringID := parametersStream.ReadUInt32LE()

	go matchmakingProtocol.FindBySingleIDHandler(nil, client, callID, gatheringID)
}

// NewMatchmakingProtocol returns a new MatchmakingProtocol
func NewMatchmakingProtocol(server *nex.Server) *MatchmakingProtocol {
	matchmakingProtocol := &MatchmakingProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	matchmakingProtocol.Setup()

	return matchmakingProtocol
}
