package client

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Fl0rencess720/suzuRPC/codec"
	"github.com/Fl0rencess720/suzuRPC/compressor"
	"github.com/Fl0rencess720/suzuRPC/protocol"
	"github.com/Fl0rencess720/suzuRPC/serializer"
)

var ErrShutdown = errors.New("connection is shut down")

type Call struct {
	ServicePath   string
	ServiceMethod string     // The name of the service and method to call.
	Args          any        // The argument to the function (*struct).
	Reply         any        // The reply from the function (*struct).
	Error         error      // After completion, the error status.
	Done          chan *Call // Receives *Call when Go is complete.
}

type Client struct {
	codec codec.ClientCodec

	Conn        net.Conn
	DialTimeout time.Duration

	mutex    sync.Mutex // protects following
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop

}

type Option func(o *options)

type options struct {
	comprressor compressor.CompressType
	serializer  serializer.Serializer
	dialTimeout time.Duration
}

func WithCompressor(c compressor.CompressType) Option {
	return func(o *options) {
		o.comprressor = c
	}
}

func WithSerializer(s serializer.Serializer) Option {
	return func(o *options) {
		o.serializer = s
	}
}
func WithDialTimeout(t time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = t
	}
}

func NewClient(network, address string, opts ...Option) (*Client, error) {
	c := Client{}
	if err := c.connect(network, address); err != nil {
		return nil, err
	}
	options := options{
		comprressor: compressor.Raw,
		serializer:  serializer.PROTOBUF,
		dialTimeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(&options)
	}
	c.codec.Compressor = compressor.Compressors[options.comprressor]
	c.codec.Serializer = serializer.Serializer(options.serializer)
	c.DialTimeout = options.dialTimeout
	return &c, nil
}

func (c *Client) Call(servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	call := <-c.Go(servicePath, serviceMethod, args, reply, make(chan *Call, 10)).Done
	return call.Error

}

func (c *Client) Go(servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.Done = done
	c.send(call)
	return call
}

func (c *Client) send(call *Call) {
	c.mutex.Lock()
	if c.shutdown || c.closing {
		call.Error = ErrShutdown
		c.mutex.Unlock()
		call.done()
		return
	}
	if c.pending == nil {
		c.pending = make(map[uint64]*Call)
	}
	seq := c.seq
	c.seq++
	c.pending[seq] = call
	c.mutex.Unlock()

	header := protocol.RequestPool.Get().(*protocol.Header)

	header.ServiceMethod = call.ServiceMethod
	header.ServicePath = call.ServicePath
	body := &protocol.Body{}
	data, err := c.codec.EncodeRequest(header, body)
	if err != nil {
		c.mutex.Lock()
		call = c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
		return
	}
	_, err = c.Conn.Write(data)

	header.Reset()
	protocol.RequestPool.Put(header)

	if err != nil {
		if e, ok := err.(*net.OpError); ok {
			if e.Err != nil {
				err = fmt.Errorf("net.OpError: %s", e.Err.Error())
			} else {
				err = errors.New("net.OpError")
			}

		}
		c.mutex.Lock()
		call = c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
		return
	}

}

func (c *Client) input() {
	var err error
	for err == nil {

	}
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
}
