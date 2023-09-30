package nexproto

import (
	"errors"
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// SecureProtocolID is the protocol ID for the Secure Connection protocol
	SecureProtocolID = 0xB

	// SecureMethodRegister is the method ID for the method Register
	SecureMethodRegister = 0x1

	// SecureMethodRequestConnectionData is the method ID for the method RequestConnectionData
	SecureMethodRequestConnectionData = 0x2

	// SecureMethodRequestURLs is the method ID for the method RequestURLs
	SecureMethodRequestURLs = 0x3

	// SecureMethodRegisterEx is the method ID for the method RegisterEx
	SecureMethodRegisterEx = 0x4

	// SecureMethodTestConnectivity is the method ID for the method TestConnectivity
	SecureMethodTestConnectivity = 0x5

	// SecureMethodUpdateURLs is the method ID for the method UpdateURLs
	SecureMethodUpdateURLs = 0x6

	// SecureMethodReplaceURL is the method ID for the method ReplaceURL
	SecureMethodReplaceURL = 0x7

	// SecureMethodSendReport is the method ID for the method SendReport
	SecureMethodSendReport = 0x8
)

// SecureProtocol handles the Secure Connection nex protocol
type SecureProtocol struct {
	server                       *nex.Server
	ConnectionIDCounter          *nex.Counter
	RegisterHandler              func(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL)
	RequestConnectionDataHandler func(err error, client *nex.Client, callID uint32, stationCID uint32, stationPID uint32)
	RequestURLsHandler           func(err error, client *nex.Client, callID uint32, stationCID uint32, stationPID uint32)
	RegisterExHandler            func(err error, client *nex.Client, callID uint32, stationUrls []string, className string, ticketData []byte)
	TestConnectivityHandler      func(err error, client *nex.Client, callID uint32)
	UpdateURLsHandler            func(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL)
	ReplaceURLHandler            func(err error, client *nex.Client, callID uint32, oldStation *nex.StationURL, newStation *nex.StationURL)
	SendReportHandler            func(err error, client *nex.Client, callID uint32, reportID uint32, report []byte)
}

// Setup initializes the protocol
func (secureProtocol *SecureProtocol) Setup() {
	nexServer := secureProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if SecureProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case SecureMethodRegister:
				go secureProtocol.handleRegister(packet)
			case SecureMethodRequestConnectionData:
				go secureProtocol.handleRequestConnectionData(packet)
			case SecureMethodRequestURLs:
				go secureProtocol.handleRequestURLs(packet)
			case SecureMethodRegisterEx:
				go secureProtocol.handleRegisterEx(packet)
			case SecureMethodTestConnectivity:
				go secureProtocol.handleTestConnectivity(packet)
			case SecureMethodUpdateURLs:
				go secureProtocol.handleUpdateURLs(packet)
			case SecureMethodReplaceURL:
				go secureProtocol.handleReplaceURL(packet)
			case SecureMethodSendReport:
				go secureProtocol.handleSendReport(packet)
			default:
				log.Printf("Unsupported Secure method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// Register sets the Register handler function
func (secureProtocol *SecureProtocol) Register(handler func(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL)) {
	secureProtocol.RegisterHandler = handler
}

// RequestConnectionData sets the RequestConnectionData handler function
func (secureProtocol *SecureProtocol) RequestConnectionData(handler func(err error, client *nex.Client, callID uint32, stationCID uint32, stationPID uint32)) {
	secureProtocol.RequestConnectionDataHandler = handler
}

// RequestURLs sets the RequestURLs handler function
func (secureProtocol *SecureProtocol) RequestURLs(handler func(err error, client *nex.Client, callID uint32, stationCID uint32, stationPID uint32)) {
	secureProtocol.RequestURLsHandler = handler
}

// RegisterEx sets the RegisterEx handler function
func (secureProtocol *SecureProtocol) RegisterEx(handler func(err error, client *nex.Client, callID uint32, stationUrls []string, className string, ticketData []byte)) {
	secureProtocol.RegisterExHandler = handler
}

// TestConnectivity sets the TestConnectivity handler function
func (secureProtocol *SecureProtocol) TestConnectivity(handler func(err error, client *nex.Client, callID uint32)) {
	secureProtocol.TestConnectivityHandler = handler
}

// UpdateURLs sets the UpdateURLs handler function
func (secureProtocol *SecureProtocol) UpdateURLs(handler func(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL)) {
	secureProtocol.UpdateURLsHandler = handler
}

// ReplaceURL sets the ReplaceURL handler function
func (secureProtocol *SecureProtocol) ReplaceURL(handler func(err error, client *nex.Client, callID uint32, oldStation *nex.StationURL, newStation *nex.StationURL)) {
	secureProtocol.ReplaceURLHandler = handler
}

// SendReport sets the SendReport handler function
func (secureProtocol *SecureProtocol) SendReport(handler func(err error, client *nex.Client, callID uint32, reportID uint32, report []byte)) {
	secureProtocol.SendReportHandler = handler
}

func (secureProtocol *SecureProtocol) handleRegister(packet nex.PacketInterface) {
	if secureProtocol.RegisterHandler == nil {
		log.Println("[Warning] SecureProtocol::Register not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[SecureProtocol::Register] Data missing list length")
		go secureProtocol.RegisterHandler(err, client, callID, make([]*nex.StationURL, 0))
		return
	}

	stationUrls, err := parametersStream.ReadListStationURL()

	if err != nil {
		go secureProtocol.RegisterHandler(err, client, callID, nil)
		return
	}

	go secureProtocol.RegisterHandler(nil, client, callID, stationUrls)
}

func (secureProtocol *SecureProtocol) handleRequestConnectionData(packet nex.PacketInterface) {
	if secureProtocol.RequestConnectionDataHandler == nil {
		log.Println("[Warning] SecureProtocol::RequestConnectionData not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 8 {
		err := errors.New("[SecureProtocol::RequestConnectionData] Data length too small")
		go secureProtocol.RequestConnectionDataHandler(err, client, callID, 0, 0)
		return
	}

	stationCID := parametersStream.ReadUInt32LE()
	stationPID := parametersStream.ReadUInt32LE()

	go secureProtocol.RequestConnectionDataHandler(nil, client, callID, stationCID, stationPID)
}

func (secureProtocol *SecureProtocol) handleRequestURLs(packet nex.PacketInterface) {
	if secureProtocol.RequestURLsHandler == nil {
		log.Println("[Warning] SecureProtocol::RequestURLs not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 8 {
		err := errors.New("[SecureProtocol::RequestURLs] Data length too small")
		go secureProtocol.RequestURLsHandler(err, client, callID, 0, 0)
		return
	}

	stationCID := parametersStream.ReadUInt32LE()
	stationPID := parametersStream.ReadUInt32LE()

	go secureProtocol.RequestURLsHandler(nil, client, callID, stationCID, stationPID)
}

func (secureProtocol *SecureProtocol) handleRegisterEx(packet nex.PacketInterface) {
	if secureProtocol.RegisterExHandler == nil {
		log.Println("[Warning] SecureProtocol::RegisterEx not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[SecureProtocol::RegisterEx] Data missing list length")
		go secureProtocol.RegisterExHandler(err, client, callID, make([]string, 0), "", make([]byte, 0))
		return
	}

	stationURLCount := parametersStream.ReadUInt32LE()
	stationUrls := make([]string, 0)

	for i := 0; i < int(stationURLCount); i++ {
		stationString, err := parametersStream.Read4ByteString()

		if err != nil {
			go secureProtocol.RegisterExHandler(err, client, callID, stationUrls, "", make([]byte, 0))
			return
		}
		stationUrls = append(stationUrls, stationString)
	}

	dataHolderType, err := parametersStream.Read4ByteString()

	if err != nil {
		go secureProtocol.RegisterExHandler(err, client, callID, stationUrls, dataHolderType, make([]byte, 0))
		return
	}

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 8 {
		err := errors.New("[SecureProtocol::RegisterEx] Data holder missing lengths")
		go secureProtocol.RegisterExHandler(err, client, callID, stationUrls, dataHolderType, make([]byte, 0))
		return
	}

	_ = parametersStream.ReadUInt32LE() // Length including next buffer length field
	dataHolderInner, err := parametersStream.ReadBuffer()

	if err != nil {
		go secureProtocol.RegisterExHandler(err, client, callID, stationUrls, dataHolderType, dataHolderInner)
		return
	}

	go secureProtocol.RegisterExHandler(nil, client, callID, stationUrls, dataHolderType, dataHolderInner)
}

func (secureProtocol *SecureProtocol) handleTestConnectivity(packet nex.PacketInterface) {
	if secureProtocol.TestConnectivityHandler == nil {
		log.Println("[Warning] SecureProtocol::TestConnectivity not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()

	go secureProtocol.TestConnectivityHandler(nil, client, callID)
}

func (secureProtocol *SecureProtocol) handleUpdateURLs(packet nex.PacketInterface) {
	if secureProtocol.UpdateURLsHandler == nil {
		log.Println("[Warning] SecureProtocol::UpdateURLs not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[SecureProtocol::UpdateURLs] Data missing list length")
		go secureProtocol.UpdateURLsHandler(err, client, callID, make([]*nex.StationURL, 0))
		return
	}

	stationURLCount := parametersStream.ReadUInt32LE()
	stationUrls := make([]*nex.StationURL, 0)

	for i := 0; i < int(stationURLCount); i++ {
		stationString, err := parametersStream.ReadString()

		if err != nil {
			go secureProtocol.UpdateURLsHandler(err, client, callID, stationUrls)
			return
		}

		station := nex.NewStationURL(stationString)
		stationUrls = append(stationUrls, station)
	}

	go secureProtocol.UpdateURLsHandler(nil, client, callID, stationUrls)
}

func (secureProtocol *SecureProtocol) handleReplaceURL(packet nex.PacketInterface) {
	if secureProtocol.ReplaceURLHandler == nil {
		log.Println("[Warning] SecureProtocol::ReplaceURL not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	oldStationString, err := parametersStream.ReadString()

	if err != nil {
		go secureProtocol.ReplaceURLHandler(err, client, callID, nex.NewStationURL(""), nex.NewStationURL(""))
		return
	}

	newStationString, err := parametersStream.ReadString()

	if err != nil {
		go secureProtocol.ReplaceURLHandler(err, client, callID, nex.NewStationURL(""), nex.NewStationURL(""))
		return
	}

	oldStation := nex.NewStationURL(oldStationString)
	newStation := nex.NewStationURL(newStationString)

	go secureProtocol.ReplaceURLHandler(nil, client, callID, oldStation, newStation)
}

func (secureProtocol *SecureProtocol) handleSendReport(packet nex.PacketInterface) {
	if secureProtocol.SendReportHandler == nil {
		log.Println("[Warning] SecureProtocol::SendReport not implemented")
		go respondNotImplemented(packet, SecureProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, secureProtocol.server)

	if len(parametersStream.Bytes()[parametersStream.ByteOffset():]) < 4 {
		err := errors.New("[SecureProtocol::SendReport] Data missing report ID")
		go secureProtocol.SendReportHandler(err, client, callID, 0, []byte{})
		return
	}

	reportID := parametersStream.ReadUInt32LE()
	report, err := parametersStream.ReadQBuffer()

	if err != nil {
		go secureProtocol.SendReportHandler(err, client, callID, 0, []byte{})
		return
	}

	go secureProtocol.SendReportHandler(nil, client, callID, reportID, report)
}

// NewSecureProtocol returns a new SecureProtocol
func NewSecureProtocol(server *nex.Server) *SecureProtocol {
	secureProtocol := &SecureProtocol{
		server:              server,
		ConnectionIDCounter: nex.NewCounter(10),
	}

	secureProtocol.Setup()

	return secureProtocol
}
