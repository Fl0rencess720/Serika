package client

import (
	"errors"
	"fmt"
	"io"
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
	Metadata      map[string]string
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
	serializer  serializer.SerializerType
	dialTimeout time.Duration
}

func WithCompressor(c compressor.CompressType) Option {
	return func(o *options) {
		o.comprressor = c
	}
}

func WithSerializer(s serializer.SerializerType) Option {
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
	c.codec = *codec.NewClientCodec(compressor.Compressors[options.comprressor], serializer.Serializers[options.serializer])
	c.DialTimeout = options.dialTimeout
	go c.input()
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
	header.ID = seq
	body := &protocol.Body{}
	data, err := c.codec.EncodeRequest(call.Args, header, body)
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
		header := protocol.ResponsePool.Get().(*protocol.Header)
		body := &protocol.Body{}
		data := make([]byte, 4096)
		n, err := c.Conn.Read(data)
		if err != nil {
			continue
		}
		fmt.Printf("len(data): %v\n", len(data[:n]))
		err = c.codec.DecodeResponse(data[:n], header, body)
		if err != nil {
			continue
		}
		c.mutex.Lock()
		call, ok := c.pending[header.ID]
		if ok {
			delete(c.pending, header.ID)
		}
		c.mutex.Unlock()

		if !ok {
			continue
		}

		if header.Status != 0 {
			call.Error = errors.New("rpc error: " + string(header.Status))
		} else {
			if len(body.Payload) > 0 {
				err = c.codec.Serializer.Decode(body.Payload, call.Reply)
				if err != nil {
					call.Error = fmt.Errorf("decode data failed: %v", err)
				}
			}
		}
		call.Metadata = body.Metadata
		header.Reset()
		protocol.ResponsePool.Put(header)
		call.done()
	}

	c.mutex.Lock()
	c.shutdown = true
	closing := c.closing
	c.Conn.Close()
	err = handleNetError(err, closing)

	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.pending = nil
	c.mutex.Unlock()
}

// 处理网络错误
func handleNetError(err error, closing bool) error {
	if e, ok := err.(*net.OpError); ok {
		if e.Err != nil {
			return fmt.Errorf("net.OpError: %s", e.Err.Error())
		}
		return errors.New("net.OpError")
	}
	if err == io.EOF {
		if closing {
			return ErrShutdown
		}
		return io.ErrUnexpectedEOF
	}
	return err
}

func (call *Call) done() {
	select {
	case call.Done <- call:
	default:
	}
}
