package mux

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	version = 4 //mux protocol version
)

const ( // cmds
	cmdSYN byte = iota // stream open
	cmdFIN             // stream close, a.k.a EOF mark
	cmdPSH             // data push
	cmdNOP             // no operation
)

const (
	sizeOfVer    = 1
	sizeOfCmd    = 1
	sizeOfLength = 4
	sizeOfSid    = 4
	headerSize   = sizeOfVer + sizeOfCmd + sizeOfSid + sizeOfLength
)

type header struct {
	version int    // version
	length  uint32 // payload size (little endian)
	cmd     byte   // command
	sid     uint32 // session id (little endian)
}

//newHeaderFrame you may set - 0, 0 for arg by default
func newHeaderFrame(sid, length uint32) *header {
	hdr := &header{
		version: version,
		length:  length,
		cmd:     cmdPSH,
		sid:     sid,
	}

	return hdr
}

//write
func (h *header) write() []byte {
	buf := make([]byte, headerSize)
	buf[0] = version
	buf[1] = cmdPSH
	binary.LittleEndian.PutUint32(buf[2:], h.length)
	binary.LittleEndian.PutUint32(buf[6:headerSize], h.sid)
	return buf
}

//reade
func (h *header) read(conn net.Conn) error {
	buf := make([]byte, headerSize)
	_, err := io.ReadAtLeast(conn, buf, headerSize)
	if err != nil {
		return err
	}
	h.version = int(buf[0])
	h.cmd = buf[1]
	h.length = binary.LittleEndian.Uint32(buf[2:])
	h.sid = binary.LittleEndian.Uint32(buf[6:headerSize])

	if h.version != version {
		return errors.New("version not support")
	}
	return nil
}
