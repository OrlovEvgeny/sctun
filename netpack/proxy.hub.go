package netpack

import (
	"log"
	"net"
	"strconv"
	"sync"
)

const (
	bufSize   = 2 << 14 // default socks5 read buffer size
	fromRange = 10001   // server port range from
	toRange   = 10501   // server port range from
)

//portStack
type portStack struct {
	sync.Mutex
	from  int
	to    int
	stack []int
}

//newPortStack this is tcp port pool
func newPortStack() *portStack {
	return &portStack{
		from:  fromRange,
		to:    toRange,
		stack: make([]int, 0, toRange-fromRange),
	}
}

//next get new port for server
func (p *portStack) next() int {
	p.Lock()
	defer p.Unlock()
	if len(p.stack) == 0 {
		for i := p.from; i < p.to; i++ {
			p.stack = append(p.stack, i)
		}
	}
	next := p.stack[0]
	p.stack = p.stack[1:]
	return next
}

//put returns port to pool
func (p *portStack) put(port int) {
	p.Lock()
	defer p.Unlock()
	p.stack = append(p.stack, port)
}

//Hub for node pool
type Hub struct {
	portStack *portStack
}

//NewHub
func NewHub() *Hub {
	return &Hub{
		portStack: newPortStack(),
	}
}

//proxyNode for pool proxy slave servers
type proxyNode struct {
	hub        *Hub
	port       int
	proxyAddr  string
	selfKV     *sync.Map //map for match client socks5 server and remote slave proxy node, key: sid; value: net.Conn
	dstSession net.Conn  //destination proxy slave

	quitSrv chan struct{} // shutdown server chan
}

//newNode
func (hub *Hub) newNode(conn net.Conn) *proxyNode {
	nextport := hub.portStack.next()
	return &proxyNode{
		selfKV:     &sync.Map{},
		dstSession: conn,
		port:       nextport,
		hub:        hub,
		proxyAddr:  "0.0.0.0:" + strconv.Itoa(nextport),
		quitSrv:    make(chan struct{}, 1),
	}
}

//shutDown shutdown proxy node
func (ph *proxyNode) shutDown() {
	ph.quitSrv <- struct{}{}
}

//runSrv proxy listener
func (ph *proxyNode) runSrv() {
	addr, _ := net.ResolveTCPAddr("tcp", ph.proxyAddr)
	log.Printf("starting new socks5 proxy-server on %s\n", addr.String())
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		log.Printf("proxy %s has been shutdown\n", addr.String())
		ph.hub.portStack.put(ph.port)
	}()
	wg := &sync.WaitGroup{}
	for {
		select {
		case <-ph.quitSrv:
			listener.Close()
			wg.Wait()
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		wg.Add(1)
		go ph.handle(wg, conn)
	}
}

//JoinToSrv spawn new proxy node with proxy listener on new port
func (hub *Hub) JoinToSrv(conn net.Conn) {
	ph := hub.newNode(conn)
	go ph.runSrv()
	go ph.muxRoute()
}