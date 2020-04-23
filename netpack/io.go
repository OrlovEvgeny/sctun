package netpack

import (
	"fmt"
	"github.com/OrlovEvgeny/sctun/mux"
	"io"
	"log"
	"net"
	"sync"
)

//muxRoute this is a router inside one slave node
func (ph *proxyNode) muxRoute() {
	for {
		//parse mux header frame.
		stream, err := mux.UpgradeMux(ph.dstSession)
		if err != nil {
			ph.shutDown()
			return
		}

		//non overhead allocate buffer, with fix size
		buf := make([]byte, stream.LengthRead())
		if _, err = stream.Read(buf); err != nil {
			//TODO for debug, maybe need close stream.
			log.Println(err)
			continue
		}

		//find client with session id
		if val, ok := ph.selfKV.Load(stream.SID()); ok {
			conn := val.(net.Conn)
			if _, err := conn.Write(buf); err != nil {
				//TODO for debug, maybe need close stream.
				log.Println(err)
			}
		}
	}
}

//handle
func (ph *proxyNode) handle(wg *sync.WaitGroup, conn net.Conn) {
	//generate new session id for new connect
	sid := mux.SIDUint32()

	//save proxy client net.Conn
	ph.selfKV.Store(sid, conn)
	defer func() {
		wg.Done()
		ph.selfKV.Delete(sid)
		conn.Close()
	}()

	for {
		buf := make([]byte, bufSize)
		n, err := conn.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("handle read - ", err)
			return
		}
		//new stream with current session id
		stream := mux.OpenStream(sid, ph.dstSession)
		if _, err := stream.Write(buf[:n]); err != nil {
			log.Println("handle stream.Write - ", err)
			return
		}
	}
}
