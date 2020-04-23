package mux

import (
	"net"
	"time"
)

// Stream implements net.Conn
type Stream struct {
	hdr  *header
	conn net.Conn
}

//Close
func (s *Stream) Close() error {
	return s.conn.Close()
}

//Read
func (s *Stream) Read(b []byte) (n int, err error) {
	return s.conn.Read(b)
}

//Write
func (s *Stream) Write(b []byte) (n int, err error) {
	return s.WriteSid(s.hdr.sid, b)
}

//WriteSid write data with mux protocol header
func (s *Stream) WriteSid(sid uint32, b []byte) (n int, err error) {
	bsize := len(b)
	hdr := newHeaderFrame(sid, uint32(bsize))
	headFrame := hdr.write()

	buf := make([]byte, 0, len(headFrame)+bsize)
	buf = append(buf, headFrame...)
	buf = append(buf, b...)
	n, err = s.conn.Write(buf)
	if n > 0 {
		n -= headerSize
	}
	return n, err
}

//LocalAddr
func (s *Stream) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

//RemoteAddr
func (s *Stream) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

//SetDeadline
func (s *Stream) SetDeadline(t time.Time) error {
	return s.conn.SetDeadline(t)
}

//SetReadDeadline
func (s *Stream) SetReadDeadline(t time.Time) error {
	return s.conn.SetReadDeadline(t)
}

//SetWriteDeadline
func (s *Stream) SetWriteDeadline(t time.Time) error {
	return s.conn.SetWriteDeadline(t)
}

//SID return header session id
func (s *Stream) SID() uint32 {
	return s.hdr.sid
}

//LengthRead return current payload size
func (s *Stream) LengthRead() uint32 {
	return s.hdr.length
}

//readHeader
func (s *Stream) readHeader() error {
	return s.hdr.read(s.conn)
}
