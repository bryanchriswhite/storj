package kademlia

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/anacrolix/utp"
	"github.com/ccding/go-stun/stun"
)

var (
	errorValueNotFound = errors.New("Value not found")
)

type Network interface {
	sendMessage(*message, bool, int64) (*expectedResponse, error)
	getMessage() chan (*message)
	messagesFin()
	timersFin()
	getDisconnect() chan (int)
	init(self *NetworkNode)
	setup(host string, port string, useStun bool, stunAddr string) (publicAddr net.Addr, err error)
	listen() error
	disconnect() error
	cancelResponse(*expectedResponse)
	isInitialized() bool
	getNetworkAddr() string
}

type UtpNetwork struct {
	socket        *utp.Socket
	sendChan      chan (*message)
	recvChan      chan (*message)
	dcStartChan   chan (int)
	dcEndChan     chan (int)
	dcTimersChan  chan (int)
	dcMessageChan chan (int)
	address       *net.UDPAddr
	connection    *net.UDPConn
	mutex         *sync.Mutex
	connected     bool
	initialized   bool
	responseMap   map[int64]*expectedResponse
	aliveConns    *sync.WaitGroup
	self          *NetworkNode
	msgCounter    int64
	remoteAddress string
	messageCodec  messageCodec
}

type expectedResponse struct {
	ch    chan (*message)
	query *message
	node  *NetworkNode
	id    int64
}

func NewUtpNetwork() Network {
	codec := &gobCodec{}

	return &UtpNetwork{
		messageCodec: codec,
	}
}

func (rn *UtpNetwork) init(self *NetworkNode) {
	netMsgInit()

	rn.self = self
	rn.mutex = &sync.Mutex{}
	rn.sendChan = make(chan (*message))
	rn.recvChan = make(chan (*message))
	rn.dcStartChan = make(chan (int), 10)
	rn.dcEndChan = make(chan (int))
	rn.dcTimersChan = make(chan (int))
	rn.dcMessageChan = make(chan (int))
	rn.responseMap = make(map[int64]*expectedResponse)
	rn.aliveConns = &sync.WaitGroup{}
	rn.connected = false
	rn.initialized = true
}

func (rn *UtpNetwork) isInitialized() bool {
	return rn.initialized
}

func (rn *UtpNetwork) getMessage() chan (*message) {
	return rn.recvChan
}

func (rn *UtpNetwork) getNetworkAddr() string {
	return rn.remoteAddress
}

func (rn *UtpNetwork) messagesFin() {
	rn.dcMessageChan <- 1
}

func (rn *UtpNetwork) getDisconnect() chan (int) {
	return rn.dcStartChan
}

func (rn *UtpNetwork) timersFin() {
	rn.dcTimersChan <- 1
}

func (rn *UtpNetwork) setup(host string, port string, useStun bool, stunAddr string) (publicAddr net.Addr, err error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	if rn.connected {
		return nil, errors.New("already connected")
	}

	remoteAddress := "[" + host + "]" + ":" + port

	socket, err := utp.NewSocket("udp", remoteAddress)
	if err != nil {
		return nil, err
	}

	if useStun {
		c := stun.NewClientWithConnection(socket)

		if stunAddr != "" {
			c.SetServerAddr(stunAddr)
		}

		_, h, err := c.Discover()
		if err != nil {
			return nil, err
		}

		_, err = c.Keepalive()
		if err != nil {
			return nil, err
		}

		host = h.IP()
		port = strconv.Itoa(int(h.Port()))
		remoteAddress = "[" + host + "]" + ":" + port
	}

	rn.remoteAddress = remoteAddress

	rn.connected = true

	rn.socket = socket

	var intPort int
	if intPort, err = strconv.Atoi(port); err != nil {

	}

	addr := &net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: intPort,
	}
	return addr, nil
}

func (rn *UtpNetwork) sendMessage(msg *message, expectResponse bool, id int64) (*expectedResponse, error) {
	rn.mutex.Lock()
	if id == -1 {
		id = rn.msgCounter
		rn.msgCounter++
	}
	msg.ID = id
	rn.mutex.Unlock()

	conn, err := rn.socket.DialTimeout(rn.self.Addr.String(), time.Second)
	if err != nil {
		return nil, err
	}

	data, err := rn.messageCodec.serialize(msg)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	if expectResponse {
		rn.mutex.Lock()
		defer rn.mutex.Unlock()
		expectedResponse := &expectedResponse{
			ch:    make(chan (*message)),
			node:  msg.Receiver,
			query: msg,
			id:    id,
		}
		// TODO we need a way to automatically clean these up as there are
		// cases where they won't be removed manually
		rn.responseMap[id] = expectedResponse
		return expectedResponse, nil
	}

	return nil, nil
}

func (rn *UtpNetwork) cancelResponse(res *expectedResponse) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	close(rn.responseMap[res.query.ID].ch)
	delete(rn.responseMap, res.query.ID)
}

func (rn *UtpNetwork) disconnect() error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	if !rn.connected {
		return errors.New("not connected")
	}
	rn.dcStartChan <- 1
	rn.dcStartChan <- 1
	<-rn.dcTimersChan
	<-rn.dcMessageChan
	close(rn.sendChan)
	close(rn.recvChan)
	close(rn.dcTimersChan)
	close(rn.dcMessageChan)
	err := rn.socket.CloseNow()
	rn.connected = false
	rn.initialized = false
	close(rn.dcEndChan)
	return err
}

func (rn *UtpNetwork) listen() error {
	for {
		conn, err := rn.socket.Accept()

		if err != nil {
			rn.disconnect()
			<-rn.dcEndChan
			return err
		}

		go func(conn net.Conn) {
			for {
				// Wait for messages
				msg, err := rn.messageCodec.deserialize(conn)
				if err != nil {
					if err.Error() == "EOF" {
						// Node went bye bye
					}
					// TODO should we penalize this node somehow ? Ban it ?
					return
				}

				isPing := msg.Type == messageTypePing

				if !areNodesEqual(msg.Receiver, rn.self, isPing) {
					// TODO should we penalize this node somehow ? Ban it ?
					continue
				}

				if msg.ID < 0 {
					// TODO should we penalize this node somehow ? Ban it ?
					continue
				}

				rn.mutex.Lock()
				if rn.connected {
					if msg.IsResponse {
						if rn.responseMap[msg.ID] == nil {
							// We were not expecting this response
							rn.mutex.Unlock()
							continue
						}

						if !areNodesEqual(rn.responseMap[msg.ID].node, msg.Sender, isPing) {
							// TODO should we penalize this node somehow ? Ban it ?
							rn.mutex.Unlock()
							continue
						}

						if msg.Type != rn.responseMap[msg.ID].query.Type {
							close(rn.responseMap[msg.ID].ch)
							delete(rn.responseMap, msg.ID)
							rn.mutex.Unlock()
							continue
						}

						if !msg.IsResponse {
							close(rn.responseMap[msg.ID].ch)
							delete(rn.responseMap, msg.ID)
							rn.mutex.Unlock()
							continue
						}

						resChan := rn.responseMap[msg.ID].ch
						rn.mutex.Unlock()
						resChan <- msg
						rn.mutex.Lock()
						close(rn.responseMap[msg.ID].ch)
						delete(rn.responseMap, msg.ID)
						rn.mutex.Unlock()
					} else {
						assertion := false
						switch msg.Type {
						case messageTypeFindNode:
							_, assertion = msg.Data.(*queryDataFindNode)
						case messageTypeFindValue:
							_, assertion = msg.Data.(*queryDataFindValue)
						case messageTypeStore:
							_, assertion = msg.Data.(*queryDataStore)
						default:
							assertion = true
						}

						if !assertion {
							fmt.Printf("Received bad message %v from %+v", msg.Type, msg.Sender)
							close(rn.responseMap[msg.ID].ch)
							delete(rn.responseMap, msg.ID)
							rn.mutex.Unlock()
							continue
						}

						rn.recvChan <- msg
						rn.mutex.Unlock()
					}
				} else {
					rn.mutex.Unlock()
				}
			}
		}(conn)
	}
}
