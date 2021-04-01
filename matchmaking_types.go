package nexproto

import nex "github.com/ihatecompvir/nex-go"

/*
	NEX and Rendez-Vous have multiple protocols for match making
	These protocols all share the same types
	In an effort to keep this library organized, each type used in all match making protocols is defined here
*/

// Gathering holds information about a matchmake gathering
type Gathering struct {
	ID                  uint32
	OwnerPID            uint32
	HostPID             uint32
	MinimumParticipants uint16
	MaximumParticipants uint16
	ParticipationPolicy uint32
	PolicyArgument      uint32
	Flags               uint32
	State               uint32
	Description         string

	*nex.Structure
}

// ExtractFromStream extracts a Gathering structure from a stream
func (gathering *Gathering) ExtractFromStream(stream *nex.StreamIn) error {
	var err error

	gathering.ID = stream.ReadUInt32LE()
	gathering.OwnerPID = stream.ReadUInt32LE()
	gathering.HostPID = stream.ReadUInt32LE()
	gathering.MinimumParticipants = stream.ReadUInt16LE()
	gathering.MaximumParticipants = stream.ReadUInt16LE()
	gathering.ParticipationPolicy = stream.ReadUInt32LE()
	gathering.PolicyArgument = stream.ReadUInt32LE()
	gathering.Flags = stream.ReadUInt32LE()
	gathering.State = stream.ReadUInt32LE()
	gathering.Description, err = stream.ReadString()

	if err != nil {
		return err
	}

	return nil
}

// NewGathering returns a new Gathering
func NewGathering() *Gathering {
	return &Gathering{}
}

// MatchmakeSession holds information about a matchmake session
type MatchmakeSession struct {
	GameMode              uint32
	Attributes            []uint32
	OpenParticipation     bool
	MatchmakeSystemType   uint32
	ApplicationData       []byte
	ParticipationCount    uint32
	ProgressScore         uint8           // NEX v3.5.0+
	SessionKey            []byte          // NEX v3.0.0+
	Option                uint32          // NEX v3.5.0+
	MatchmakeParam        *MatchmakeParam // NEX v4.0.0+
	StartedTime           *nex.DateTime   // NEX v4.0.0+
	UserPassword          string          // NEX v4.0.0+
	ReferGID              uint32          // NEX v4.0.0+
	UserPasswordEnabled   bool            // NEX v4.0.0+
	SystemPasswordEnabled bool            // NEX v4.0.0+
	CodeWord              string          // NEX v4.0.0+

	hierarchy []nex.StructureInterface
	*Gathering
}

// GetHierarchy returns the Structure hierarchy
func (matchmakeSession *MatchmakeSession) GetHierarchy() []nex.StructureInterface {
	return matchmakeSession.hierarchy
}

// ExtractFromStream extracts a MatchmakeSession structure from a stream
func (matchmakeSession *MatchmakeSession) ExtractFromStream(stream *nex.StreamIn) error {
	var err error
	server := stream.Server

	matchmakeSession.GameMode = stream.ReadUInt32LE()
	matchmakeSession.Attributes = stream.ReadListUInt32LE()
	matchmakeSession.OpenParticipation = stream.ReadUInt8() == 1
	matchmakeSession.MatchmakeSystemType = stream.ReadUInt32LE()
	matchmakeSession.ApplicationData, err = stream.ReadBuffer()

	if err != nil {
		return err
	}

	matchmakeSession.ParticipationCount = stream.ReadUInt32LE()

	if server.NexVersion() >= 30500 {
		matchmakeSession.ProgressScore = stream.ReadUInt8()
	}

	if server.NexVersion() >= 30000 {
		matchmakeSession.SessionKey, err = stream.ReadBuffer()

		if err != nil {
			return err
		}
	}

	if server.NexVersion() >= 30500 {
		matchmakeSession.Option = stream.ReadUInt32LE()
	}

	if server.NexVersion() >= 40000 {
		matchmakeParam, err := stream.ReadStructure(NewMatchmakeParam())

		if err != nil {
			return err
		}

		matchmakeSession.MatchmakeParam = matchmakeParam.(*MatchmakeParam)
		matchmakeSession.StartedTime = nex.NewDateTime(stream.ReadUInt64LE())
		matchmakeSession.UserPassword, err = stream.ReadString()

		if err != nil {
			return err
		}

		matchmakeSession.ReferGID = stream.ReadUInt32LE()
		matchmakeSession.UserPasswordEnabled = stream.ReadUInt8() == 1
		matchmakeSession.SystemPasswordEnabled = stream.ReadUInt8() == 1
		matchmakeSession.CodeWord, err = stream.ReadString()

		if err != nil {
			return err
		}
	}

	return nil
}

// NewMatchmakeSession returns a new MatchmakeSession
func NewMatchmakeSession() *MatchmakeSession {
	matchmakeSession := &MatchmakeSession{}

	gathering := NewGathering()

	matchmakeSession.Gathering = gathering

	matchmakeSession.hierarchy = []nex.StructureInterface{
		gathering,
	}

	return matchmakeSession
}

// MatchmakeParam holds parameters for a matchmake session
type MatchmakeParam struct {
	parameters map[interface{}]interface{}

	*nex.Structure
}

// ExtractFromStream extracts a MatchmakeParam structure from a stream
func (matchmakeParam *MatchmakeParam) ExtractFromStream(stream *nex.StreamIn) error {
	var err error

	matchmakeParam.parameters, err = stream.ReadMap(stream.ReadString, stream.ReadVariant)

	if err != nil {
		return err
	}

	return nil
}

// NewMatchmakeParam returns a new MatchmakeParam
func NewMatchmakeParam() *MatchmakeParam {
	return &MatchmakeParam{}
}
