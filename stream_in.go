package nexproto

import (
	nex "github.com/ihatecompvir/nex-go"
)

// StreamIn is an abstraction of StreamIn from github.com/ihatecompvir/nex-go
// Adds protocol-specific Structure list support
type StreamIn struct {
	*nex.StreamIn
}

// ReadListStationURL reads a list of StationURL structures
func (stream *StreamIn) ReadListStationURL() ([]*nex.StationURL, error) {
	length := stream.ReadUInt32LE()
	stationUrls := make([]*nex.StationURL, 0)

	for i := 0; i < int(length); i++ {
		stationString, err := stream.ReadString()

		if err != nil {
			return nil, err
		}

		station := nex.NewStationURL(stationString)
		stationUrls = append(stationUrls, station)
	}

	return stationUrls, nil
}

// NewStreamIn returns a new nexproto output stream
func NewStreamIn(data []byte, server *nex.Server) *StreamIn {
	return &StreamIn{
		StreamIn: nex.NewStreamIn(data, server),
	}
}
