package nexproto

import (
	"errors"
	"log"

	nex "github.com/ihatecompvir/nex-go"
)

const (
	// AccountManagementProtocolID is the protocol ID for the Account Management protocol
	AccountManagementProtocolID = 0x19
	DeleteAccount               = 0x02
	SetStatus                   = 0x11
	FindByNameLike              = 0x19
	LookupOrCreateAccount       = 0x1B // also used by Xbox 360 when multiple profiles are signed in
)

// AccountManagementProtocol handles the Account Management nex protocol
type AccountManagementProtocol struct {
	server                       *nex.Server
	DeleteAccountHandler         func(err error, client *nex.Client, callID uint32, pid uint32)
	LookupOrCreateAccountHandler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)
	SetStatusHandler             func(err error, client *nex.Client, callID uint32, status string)
	FindByNameLikeHandler        func(err error, client *nex.Client, callID uint32, uiGroups uint32, name string)
}

// Setup initializes the protocol
func (accountManagementProtocol *AccountManagementProtocol) Setup() {
	nexServer := accountManagementProtocol.server

	nexServer.On("Data", func(packet nex.PacketInterface) {
		request := packet.RMCRequest()

		if AccountManagementProtocolID == request.ProtocolID() {
			switch request.MethodID() {
			case DeleteAccount:
				go accountManagementProtocol.handleDeleteAccount(packet)
			case LookupOrCreateAccount:
				go accountManagementProtocol.handleLookupOrCreateAccount(packet)
			case SetStatus:
				go accountManagementProtocol.handleSetStatus(packet)
			case FindByNameLike:
				go accountManagementProtocol.handleFindByNameLike(packet)
			default:
				log.Printf("Unsupported AccountManagement method ID: %#v\n", request.MethodID())
			}
		}
	})
}

// DeleteAccount sets the DeleteAccount handler function
func (accountManagementProtocol *AccountManagementProtocol) DeleteAccount(handler func(err error, client *nex.Client, callID uint32, pid uint32)) {
	accountManagementProtocol.DeleteAccountHandler = handler
}

// CreateAccount sets the CreateAccount handler function
func (accountManagementProtocol *AccountManagementProtocol) LookupOrCreateAccount(handler func(err error, client *nex.Client, callID uint32, username string, key string, groups uint32, email string)) {
	accountManagementProtocol.LookupOrCreateAccountHandler = handler
}

// SetStatus sets the SetStatus handler function
func (accountManagementProtocol *AccountManagementProtocol) SetStatus(handler func(err error, client *nex.Client, callID uint32, status string)) {
	accountManagementProtocol.SetStatusHandler = handler
}

// FindByNameLike sets the FindByNameLike handler function
func (accountManagementProtocol *AccountManagementProtocol) FindByNameLike(handler func(err error, client *nex.Client, callID uint32, uiGroups uint32, name string)) {
	accountManagementProtocol.FindByNameLikeHandler = handler
}

func (accountManagementProtocol *AccountManagementProtocol) handleDeleteAccount(packet nex.PacketInterface) {
	if accountManagementProtocol.DeleteAccountHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::DeleteAccount not implemented")
		go respondNotImplemented(packet, AccountManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, accountManagementProtocol.server)

	pid := parametersStream.ReadUInt32LE()

	go accountManagementProtocol.DeleteAccountHandler(nil, client, callID, pid)
}

func (accountManagementProtocol *AccountManagementProtocol) handleLookupOrCreateAccount(packet nex.PacketInterface) {
	if accountManagementProtocol.LookupOrCreateAccountHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::LookupOrCreateAccount not implemented")
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
		go accountManagementProtocol.LookupOrCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	key, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.LookupOrCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	groups := parametersStream.ReadUInt32LE()
	email, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.LookupOrCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	dataHolderName, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.LookupOrCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	// I don't think PS3 can ever call this method, but just in case
	if dataHolderName != "NintendoToken" && dataHolderName != "XboxUserInfo" && dataHolderName != "SonyNPTicket" {
		err := errors.New("[AccountManagementProtocol::LookupOrCreateAccount] Data holder name does not match")
		go accountManagementProtocol.LookupOrCreateAccountHandler(err, client, callID, "", "", 0, "")
		return
	}

	go accountManagementProtocol.LookupOrCreateAccountHandler(nil, client, callID, username, key, groups, email)
}

func (accountManagementProtocol *AccountManagementProtocol) handleSetStatus(packet nex.PacketInterface) {
	if accountManagementProtocol.SetStatusHandler == nil {
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

func (accountManagementProtocol *AccountManagementProtocol) handleFindByNameLike(packet nex.PacketInterface) {
	if accountManagementProtocol.FindByNameLikeHandler == nil {
		log.Println("[Warning] AccountManagementProtocol::FindByNameLike not implemented")
		go respondNotImplemented(packet, AccountManagementProtocolID)
		return
	}

	client := packet.Sender()
	request := packet.RMCRequest()

	callID := request.CallID()
	parameters := request.Parameters()

	parametersStream := nex.NewStreamIn(parameters, accountManagementProtocol.server)

	uiGroups := parametersStream.ReadUInt32LE()
	name, err := parametersStream.Read4ByteString()
	if err != nil {
		go accountManagementProtocol.FindByNameLikeHandler(err, client, callID, 0, "")
		return
	}

	go accountManagementProtocol.FindByNameLikeHandler(nil, client, callID, uiGroups, name)
}

// NewAccountManagementProtocol returns a new AccountManagementProtocol
func NewAccountManagementProtocol(server *nex.Server) *AccountManagementProtocol {
	accountManagementProtocol := &AccountManagementProtocol{server: server}

	accountManagementProtocol.Setup()

	return accountManagementProtocol
}
