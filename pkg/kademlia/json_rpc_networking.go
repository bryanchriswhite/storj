package kademlia

import (
	"errors"
	"sync"
	"net/rpc/jsonrpc"
	"net/rpc"
	"net"
	"fmt"
	"math/rand"
	"github.com/ccding/go-stun/stun"
)

type RpcMethods struct{}

type JsonRpcResponse struct {
	result Result
	id     Id
}

type JsonRpcRequest struct {
	method Method
	params Params
	id     Id
}

type Method string
type Id string
type Nonce int64
type Signature string
type Address string
type Port int
type NodeId string

type Params struct {
	contact   Contact
	nonce     Nonce
	signature Signature
}

type Result struct {
	contact   Contact
	nonce     Nonce
	signature Signature
}

type Contact struct {
	address  Address
	port     Port
	nodeId   NodeId
	protocol string
}

func NodeToContact(node *NetworkNode) (contact Contact) {
	contact.address = Address(node.Addr.IP.String())
	contact.port = Port(node.Addr.Port)
	contact.nodeId = NodeId(node.ID)
	contact.protocol = "1.2.0"

	return contact
}

func (rm *RpcMethods) Ping(args Params, reply *JsonRpcResponse) {

}

func NewNonce() Nonce {
	return Nonce(rand.Int63n(1e16))
}

type JsonRpcNetworking struct {
	mutex       sync.Mutex
	connected   bool
	initialized bool
	publicAddr  net.TCPAddr
	client      *rpc.Client
	msgCounter  int64
}

func (jn *JsonRpcNetworking) sendMessage(msg *message, expectResponse bool, id int64) (res *expectedResponse, err error) {
	jn.mutex.Lock()
	if id == -1 {
		id = jn.msgCounter
		jn.msgCounter++
	}
	msg.ID = id
	jn.mutex.Unlock()

	var remoteMethod string
	switch msg.Type {
	case 0:
		remoteMethod = "RpcMethods.Ping"
	default:
		err = errors.New(fmt.Sprintf("Unknown request method iota: %d", msg.Type))
		return nil, err
	}

	// reply := new(JsonRpcResponse)
	// contact := NodeToContact(msg.Receiver)
	params := &Params{
		contact:   NodeToContact(msg.Sender),
		nonce:     NewNonce(),
		signature: "",
	}

	client, err = jsonrpc.Dial("tcp", msg.Receiver.Addr.String())
	client.Call(remoteMethod, params, reply)
	// conn, err := jn.socket.DialTimeout("["+msg.Receiver.IP.String()+"]:"+strconv.Itoa(msg.Receiver.Port), time.Second)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// data, err := jn.messageCodec.serialize(msg)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// _, err = conn.Write(data)
	// if err != nil {
	// 	return nil, err
	// }

	if expectResponse {
		jn.mutex.Lock()
		defer jn.mutex.Unlock()
		expectedResponse := &expectedResponse{
			ch:    make(chan (*message)),
			node:  msg.Receiver,
			query: msg,
			id:    id,
		}
		// TODO we need a way to automatically clean these up as there are
		// cases where they won't be removed manually
		jn.responseMap[id] = expectedResponse
		return expectedResponse, nil
	}

	return nil, nil
}

func (jn *JsonRpcNetworking) getMessage() chan int {

}

func (jn *JsonRpcNetworking) messagesFin() {

}

func (jn *JsonRpcNetworking) timersFin() {

}

func (jn *JsonRpcNetworking) getDisconnect() chan int {

}

func (jn *JsonRpcNetworking) init(self *NetworkNode) {
	jn.connected = false
	jn.initialized = true
}

func (jn *JsonRpcNetworking) createSocket(host string, port string, useStun bool, stunAddr string) (publicAddr *net.TCPAddr, err error) {
	jn.mutex.Lock()
	defer jn.mutex.Unlock()
	if jn.connected {
		return nil, errors.New("already connected")
	}

	addr := &net.TCPAddr{}

	if useStun {
		c := stun.NewClient()

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

		addr.IP = net.ParseIP(h.IP())
		addr.Port = int(h.Port())
	}

	jn.publicAddr = *addr

	return addr, nil
}

func (jn *JsonRpcNetworking) listen() error {
	// listener, err := net.Listen("tcp", jn.remoteAddress)
	// server := rpc.NewServer()
	err := rpc.Register(&RpcMethods{})
	rpc.HandleHTTP()
	// listener, err := net.Listen("tcp", "")
	if err != nil {
		return err
	}

	for {
		go func() {
			// rpc.ServeCodec(jsonrpc.NewServerCodec(conn))

			// for {
			// 	// Wait for messages
			// 	msg, err := deserializeMessage(conn)
			// 	if err != nil {
			// 		if err.Error() == "EOF" {
			// 			// Node went bye bye
			// 		}
			// 		// TODO should we penalize this node somehow ? Ban it ?
			// 		return
			// 	}
			//
			// 	isPing := msg.Type == messageTypePing
			//
			// 	if !areNodesEqual(msg.Receiver, jn.self, isPing) {
			// 		// TODO should we penalize this node somehow ? Ban it ?
			// 		continue
			// 	}
			//
			// 	if msg.ID < 0 {
			// 		// TODO should we penalize this node somehow ? Ban it ?
			// 		continue
			// 	}
			//
			// 	jn.mutex.Lock()
			// 	if jn.connected {
			// 		if msg.IsResponse {
			// 			if jn.responseMap[msg.ID] == nil {
			// 				// We were not expecting this response
			// 				jn.mutex.Unlock()
			// 				continue
			// 			}
			//
			// 			if !areNodesEqual(jn.responseMap[msg.ID].node, msg.Sender, isPing) {
			// 				// TODO should we penalize this node somehow ? Ban it ?
			// 				jn.mutex.Unlock()
			// 				continue
			// 			}
			//
			// 			if msg.Type != jn.responseMap[msg.ID].query.Type {
			// 				close(jn.responseMap[msg.ID].ch)
			// 				delete(jn.responseMap, msg.ID)
			// 				jn.mutex.Unlock()
			// 				continue
			// 			}
			//
			// 			if !msg.IsResponse {
			// 				close(jn.responseMap[msg.ID].ch)
			// 				delete(jn.responseMap, msg.ID)
			// 				jn.mutex.Unlock()
			// 				continue
			// 			}
			//
			// 			resChan := jn.responseMap[msg.ID].ch
			// 			jn.mutex.Unlock()
			// 			resChan <- msg
			// 			jn.mutex.Lock()
			// 			close(jn.responseMap[msg.ID].ch)
			// 			delete(jn.responseMap, msg.ID)
			// 			jn.mutex.Unlock()
			// 		} else {
			// 			assertion := false
			// 			switch msg.Type {
			// 			case messageTypeFindNode:
			// 				_, assertion = msg.Data.(*queryDataFindNode)
			// 			case messageTypeFindValue:
			// 				_, assertion = msg.Data.(*queryDataFindValue)
			// 			case messageTypeStore:
			// 				_, assertion = msg.Data.(*queryDataStore)
			// 			default:
			// 				assertion = true
			// 			}
			//
			// 			if !assertion {
			// 				fmt.Printf("Received bad message %v from %+v", msg.Type, msg.Sender)
			// 				close(jn.responseMap[msg.ID].ch)
			// 				delete(jn.responseMap, msg.ID)
			// 				jn.mutex.Unlock()
			// 				continue
			// 			}
			//
			// 			jn.recvChan <- msg
			// 			jn.mutex.Unlock()
			// 		}
			// 	} else {
			// 		jn.mutex.Unlock()
			// 	}
			// }
		}()
	}
}

func (jn *JsonRpcNetworking) disconnect() error {

}

func (jn *JsonRpcNetworking) cancelResponse(*expectedResponse) {

}

func (jn *JsonRpcNetworking) isInitialized() bool {
	return jn.initialized
}

func (jn *JsonRpcNetworking) getNetworkingAddr() string {

}
