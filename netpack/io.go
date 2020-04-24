package netpack

import (
	"github.com/OrlovEvgeny/sctun/mux"
	"io"
	"log"
	"net"
)

//muxRoute this is a router inside one slave node
func (ph *proxyNode) muxRoute() {
	defer ph.shutdownNode()

	for {
		//parse mux header frame.
		stream, err := mux.UpgradeMux(ph.dstSession)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		buf := make([]byte, stream.LengthRead())
		if _, err = io.ReadAtLeast(stream, buf, stream.LengthRead()); err != nil {
			//TODO for debug, maybe need close stream.
			log.Println(err)
			return
		}

		//find client with session id
		if val, ok := ph.selfKV.Load(stream.SID()); ok {
			conn := val.(net.Conn)
			if _, err := conn.Write(buf); err != nil {
				//TODO for debug, maybe need close stream.
				log.Println(err)
			}
			continue
		}
	}
}

//handle
func (ph *proxyNode) handle(conn net.Conn) {
	//generate new session id for new connect
	sid := mux.SIDUint32()

	//save proxy client net.Conn
	ph.selfKV.Store(sid, conn)
	defer func() {
		ph.selfKV.Delete(sid)
		conn.Close()
		ph.waitGroup.Done()
	}()

	//new stream with current session id
	stream := mux.OpenStream(sid, ph.dstSession)
	buf := make([]byte, bufSize)
	for {
		select {
		case <-ph.quitSrv:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}

		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("handle for %d sid, read - %s\n", sid, err.Error())
			}
			return
		}

		if _, err := stream.Write(buf[:n]); err != nil {
			log.Printf("handle for %d sid, write - %s\n", sid, err.Error())
			return
		}
	}
}
