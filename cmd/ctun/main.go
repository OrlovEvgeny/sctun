package main

import (
	"flag"
	"fmt"
	"github.com/OrlovEvgeny/sctun/mux"
	"github.com/OrlovEvgeny/sctun/socks5"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const defaultAddr = "127.0.0.1:8080"

var (
	addr string

	//delay reconnect counter (in seconds), step 5 sec.
	delay uint32

	//map for pipe io.PipeWriter, key: mux protocol SID (session id); value: ptr on io.PipeWriter.
	pipeKV = &sync.Map{}
)

func init() {
	flag.StringVar(&addr, "master", defaultAddr, "--master <ip_master_server>")
	flag.Parse()
}

func main() {
	cerr := make(chan error, 1000)
	go logger(cerr)

	srvDialer(cerr)
}

//logger
func logger(cerr <-chan error) {
	for {
		err := <-cerr

		if err == io.EOF {
			log.Println("server close connect")
		}

		log.Println(err)
	}
}

//srvDialer client loop dial to master server
func srvDialer(cerr chan<- error) {
	for {
		//check reconnect counter
		if err := delayDial(); err != nil {
			log.Fatal(err)
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Fatalf("tcp addr incorect, see --help. Error: %s", err.Error())
		}
		conn, err := net.Dial("tcp", tcpAddr.String())
		if err != nil {
			cerr <- err
			atomic.AddUint32(&delay, 5)
			continue
		}

		log.Printf("dial to master server %s ok.\n", tcpAddr.String())

		handle(conn, cerr)
		if err := conn.Close(); err != nil {
			cerr <- err
		}
		atomic.StoreUint32(&delay, 0)
	}
}

//delayDial delay counter before start try reconnect to server (step 5 sec).
func delayDial() error {
	if count := atomic.LoadUint32(&delay); count != 0 {
		if count >= 300 {
			return fmt.Errorf("reconnect timed out")
		}
		fmt.Printf("attempt to connect %d sec. ...\n", count)
		time.Sleep(time.Duration(count) * time.Second)
	}
	return nil
}

//handle
func handle(conn net.Conn, cerr chan<- error) {
	//TODO maybe need have make custom config, or no :)
	srvSocks5, err := socks5.New(&socks5.Config{})
	if err != nil {
		cerr <- err
		return
	}
	for {
		//parse mux header frame.
		stream, err := mux.UpgradeMux(conn)
		if err != nil {
			cerr <- err
			return
		}

		//non overhead allocate buffer, with fix size
		buf := make([]byte, stream.LengthRead())
		if _, err := stream.Read(buf); err != nil {
			cerr <- err
			return
		}

		/*
			search pipeWriter in map,
			if found it means that the socks5 connection already exists and wait data,
			try write buffer
		*/
		if value, ok := pipeKV.Load(stream.SID()); ok {
			w := value.(*io.PipeWriter)
			if _, err := w.Write(buf); err != nil {
				log.Println("pipe write error")
				cerr <- err
			}
			continue
		}

		//create new rw pipe
		rpipe, wpipe := io.Pipe()
		//save new pipeWriter as key SID
		pipeKV.Store(stream.SID(), wpipe)

		//spawn new mux for one uniq SID
		go func(conn *mux.Stream, r *io.PipeReader, w *io.PipeWriter) {
			defer func() {
				log.Println("pipe shutdown")
				pipeKV.Delete(conn.SID())
				r.Close()
				w.Close()
			}()

			if err := srvSocks5.MuxConn(conn, r); err != nil {
				log.Println(err)
			}

		}(stream, rpipe, wpipe)

		//write welcome message to new pipeWriter
		if _, err := wpipe.Write(buf); err != nil {
			cerr <- err
		}
	}
}
