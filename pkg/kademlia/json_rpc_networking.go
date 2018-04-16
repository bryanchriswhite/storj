package kademlia
//
// import (
// 	"errors"
// 	"sync"
// 	"net/rpc/jsonrpc"
// 	"net/rpc"
// 	"net"
// 	"fmt"
// 	"math/rand"
// 	"github.com/ccding/go-stun/stun"
// )
//
// type RpcMethods struct{}
//
// type JsonRpcResponse struct {
// 	result Result
// 	id     Id
// }
//
// type JsonRpcRequest struct {
// 	method Method
// 	params Params
// 	id     Id
// }
//
// type Method string
// type Id string
// type Nonce int64
// type Signature string
// type Address string
// type Port int
// type NodeId string
//
// type Params struct {
// 	contact   Contact
// 	nonce     Nonce
// 	signature Signature
// }
//
// type Result struct {
// 	contact   Contact
// 	nonce     Nonce
// 	signature Signature
// }
//
// type Contact struct {
// 	address  Address
// 	port     Port
// 	nodeId   NodeId
// 	protocol string
// }
//
// func NodeToContact(node *NetworkNode) (contact Contact) {
// 	contact.address = Address(node.Addr.IP.String())
// 	contact.port = Port(node.Addr.Port)
// 	contact.nodeId = NodeId(node.ID)
// 	contact.protocol = "1.2.0"
//
// 	return contact
// }
//
// func (rm *RpcMethods) Ping(args Params, reply *JsonRpcResponse) {
//
// }
//
// func NewNonce() Nonce {
// 	return Nonce(rand.Int63n(1e16))
// }
//
// func NewJsonRpcNetwork() Network {
// 	return &JsonRpcNetwork{}
// }
//
// type JsonRpcNetwork struct {
// 	self        *NetworkNode
// 	mutex       sync.Mutex
// 	connected   bool
// 	initialized bool
// 	publicAddr  net.TCPAddr
// 	client      *rpc.Client
// 	msgCounter  int64
// }
//
// func (jn *JsonRpcNetwork) sendMessage(msg *message, expectResponse bool, id int64) (res *expectedResponse, err error) {
// 	jn.mutex.Lock()
// 	if id == -1 {
// 		id = jn.msgCounter
// 		jn.msgCounter++
// 	}
// 	msg.ID = id
// 	jn.mutex.Unlock()
//
// 	var remoteMethod string
// 	switch msg.Type {
// 	case 0:
// 		remoteMethod = "RpcMethods.Ping"
// 	default:
// 		err = errors.New(fmt.Sprintf("Unknown request method iota: %d", msg.Type))
// 		return nil, err
// 	}
//
// 	reply := new(JsonRpcResponse)
// 	params := &Params{
// 		contact:   NodeToContact(msg.Sender),
// 		nonce:     NewNonce(),
// 		signature: "",
// 	}
//
// 	client, err := jsonrpc.Dial("tcp", msg.Receiver.Addr.String())
// 	client.Call(remoteMethod, params, reply)
// 	// conn, err := jn.socket.DialTimeout("["+msg.Receiver.IP.String()+"]:"+strconv.Itoa(msg.Receiver.Port), time.Second)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// data, err := jn.messageCodec.serialize(msg)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// _, err = conn.Write(data)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
//
// 	if expectResponse {
// 		jn.mutex.Lock()
// 		defer jn.mutex.Unlock()
// 		expectedResponse := &expectedResponse{
// 			ch:    make(chan (*message)),
// 			node:  msg.Receiver,
// 			query: msg,
// 			id:    id,
// 		}
// 		// TODO we need a way to automatically clean these up as there are
// 		// cases where they won't be removed manually
// 		jn.responseMap[id] = expectedResponse
// 		return expectedResponse, nil
// 	}
//
// 	return nil, nil
// }
//
// func (jn *JsonRpcNetwork) getMessage() chan int {
//
// }
//
// func (jn *JsonRpcNetwork) messagesFin() {
//
// }
//
// func (jn *JsonRpcNetwork) timersFin() {
//
// }
//
// func (jn *JsonRpcNetwork) getDisconnect() chan int {
// 	return jn.dcChan
// }
//
// func (jn *JsonRpcNetwork) init(self *NetworkNode) {
// 	jn.self = self
// 	jn.connected = false
// 	jn.initialized = true
// 	jn.dcChan = make(chan int)
// }
//
// func (jn *JsonRpcNetwork) setup(host string, port string, useStun bool, stunAddr string) (publicAddr net.Addr, err error) {
// 	jn.mutex.Lock()
// 	defer jn.mutex.Unlock()
// 	if jn.connected {
// 		return nil, errors.New("already connected")
// 	}
//
// 	addr := &net.TCPAddr{}
//
// 	if useStun {
// 		c := stun.NewClient()
//
// 		if stunAddr != "" {
// 			c.SetServerAddr(stunAddr)
// 		}
//
// 		_, h, err := c.Discover()
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		_, err = c.Keepalive()
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		addr.IP = net.ParseIP(h.IP())
// 		addr.Port = int(h.Port())
// 	}
//
// 	jn.publicAddr = *addr
//
// 	return addr, nil
// }
//
// func (jn *JsonRpcNetwork) listen() (err error) {
// 	server := rpc.NewServer()
//
// 	if err = rpc.Register(&RpcMethods{}); err != nil {
// 		return err
// 	}
//
// 	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
//
// 	var listener net.Listener
// 	if listener, err = net.Listen("tcp", jn.publicAddr.String()); err != nil {
// 		return err
// 	}
//
// 	for {
// 		var conn *net.Conn
// 		if *conn, err = listener.Accept(); err != nil {
// 			return err
// 		}
//
// 		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
// 	}
// }
//
// func (jn *JsonRpcNetwork) disconnect() error {
//
// }
//
// func (jn *JsonRpcNetwork) cancelResponse(*expectedResponse) {
//
// }
//
// func (jn *JsonRpcNetwork) isInitialized() bool {
// 	return jn.initialized
// }
//
// func (jn *JsonRpcNetwork) getNetworkingAddr() string {
//
// }
