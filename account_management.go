package nexproto

import (
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// AccountManagementProtocolID is the protocol ID for the Account Management protocol
	AccountManagementProtocolID = 0x19

	SetStatus             = 0x11
	NintendoCreateAccount = 0x1B
)

// AccountManagementProtocol handles the Account Management nex protocol
type AccountManagementProtocol struct {
	server                       *nex.Server
	NintendoCreateAccountHandler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)
	SetStatusHandler             func(err error, client *nex.Client, callID uint32, status string)
}

// Setup initializes the protocol
func (accountManagementProtocol *AccountManagementProtocol) Setup() {
	nexServer := accountManagementProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if AccountManagementProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case NintendoCreateAccount:
				go accountManagementProtocol.handleNintendoCreateAccount(packet)
			case SetStatus:
				go accountManagementProtocol.handleSetStatus(packet)
			default:
				log.Printf("Unsupported AccountManagement method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// NintendoCreateAccount sets the NintendoCreateAccount handler function
func (accountManagementProtocol *AccountManagementProtocol) NintendoCreateAccount(handler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)) {
	accountManagementProtocol.NintendoCreateAccountHandler = handler
}

// SetStatus sets the SetStatus handler function
func (accountManagementProtocol *AccountManagementProtocol) SetStatus(handler func(err error, client *nex.Client, callID uint32, status string)) {
	accountManagementProtocol.SetStatusHandler = handler
}

func (accountManagementProtocol *AccountManagementProtocol) handleNintendoCreateAccount(packet nex.PacketInterface) {
	if accountManagementProtocol.NintendoCreateAccountHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::NintendoCreateAccount not implemented")
		go respondNotImplemented(packet, AccountManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, accountManagementProtocol.server)

	username, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	key, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	groups := parametersStream.ReadUInt32LE()
	email, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.NintendoCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	go accountManagementProtocol.NintendoCreateAccountHandler(nil, client, callID, username, key, groups, email)
}

func (accountManagementProtocol *AccountManagementProtocol) handleSetStatus(packet nex.PacketInterface) {
	if accountManagementProtocol.NintendoCreateAccountHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::SetStatus not implemented")
		go respondNotImplemented(packet, AccountManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, accountManagementProtocol.server)

	status, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.SetStatusHandler(err, client, callID, "")
		return
	}

	go accountManagementProtocol.SetStatusHandler(nil, client, callID, status)
}

// NewAccountManagementProtocol returns a new AccountManagementProtocol
func NewAccountManagementProtocol(server *nex.Server) *AccountManagementProtocol {
	accountManagementProtocol := &AccountManagementProtocol{server: server}

	accountManagementProtocol.Setup()

	return accountManagementProtocol
}
