package mux

import "net"

//UpgradeMux just read connection header and parse mux protocol, return updated net.Conn as stream
func UpgradeMux(conn net.Conn) (*Stream, error) {
	stream := &Stream{conn: conn, hdr: newHeaderFrame(0, 0)}
	if err := stream.readHeader(); err != nil {
		return stream, err
	}
	return stream, nil
}

//OpenStream create new header with you SID, return updated net.Conn as stream
func OpenStream(sid uint32, conn net.Conn) *Stream {
	return &Stream{conn: conn, hdr: newHeaderFrame(sid, 0)}
}
