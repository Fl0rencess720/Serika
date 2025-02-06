package server

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/Fl0rencess720/suzuRPC/codec"
	"github.com/Fl0rencess720/suzuRPC/compressor"
	"github.com/Fl0rencess720/suzuRPC/protocol"
	"github.com/Fl0rencess720/suzuRPC/serializer"
)

type Server struct {
	codec    codec.ServerCodec
	ln       net.Listener
	services map[string]reflect.Value
	mutex    sync.Mutex
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]reflect.Value),
	}
}

func (s *Server) Register(serviceName string, service interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.services[serviceName]; exists {
		return errors.New("service already registered")
	}
	s.services[serviceName] = reflect.ValueOf(service)
	return nil
}

func (s *Server) Serve(network, address string) error {
	ln, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer ln.Close()

	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	s.codec = codec.ServerCodec{
		Compressor: compressor.Compressors[compressor.Raw],
		Serializer: serializer.Serializers[serializer.JSON],
	}

	for {
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		if err != nil {
			return
		}

		header := protocol.RequestPool.Get().(*protocol.Header)
		body := &protocol.Body{}
		if err := s.codec.DecodeRequest(data[:n], header, body); err != nil {
			return
		}
		response := s.handleRequest(header, body)
		_, err = conn.Write(response)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
}

func (s *Server) handleRequest(header *protocol.Header, body *protocol.Body) []byte {
	service, method := header.ServicePath, header.ServiceMethod
	// 获取服务实例
	serviceValue, exists := s.services[service]
	if !exists {
		return encodeErrorResponse(s.codec, header.ID, "service not found")
	}

	// 获取方法
	methodValue := serviceValue.MethodByName(method)
	if !methodValue.IsValid() {
		return encodeErrorResponse(s.codec, header.ID, "method not found")
	}

	// 检查方法参数数量
	numArgs := methodValue.Type().NumIn()
	if numArgs != 2 {
		return encodeErrorResponse(s.codec, header.ID, fmt.Sprintf("method requires %d arguments, but got 2", numArgs))
	}

	// 创建 args 和 reply 实例
	argsType := methodValue.Type().In(0)
	replyType := methodValue.Type().In(1)

	// 如果 argsType 是指针类型，需要解引用
	var args interface{}
	if argsType.Kind() == reflect.Ptr {
		args = reflect.New(argsType.Elem()).Interface()
	} else {
		args = reflect.New(argsType).Elem().Interface()
	}

	reply := reflect.New(replyType.Elem()).Interface()
	if err := serializer.Serializers[serializer.JSON].Decode(body.Payload, args); err != nil {
		return encodeErrorResponse(s.codec, header.ID, fmt.Sprintf("decode error: %v", err))
	}
	results := methodValue.Call([]reflect.Value{
		reflect.ValueOf(args),
		reflect.ValueOf(reply),
	})
	var err error
	if len(results) > 1 && !results[1].IsNil() {
		err = results[1].Interface().(error)
	}
	if err != nil {
		return encodeErrorResponse(s.codec, header.ID, err.Error())
	}
	return encodeSuccessResponse(s.codec, header.ID, reply)
}

func encodeErrorResponse(codec codec.ServerCodec, id uint64, errMsg string) []byte {
	return encodeResponse(codec, id, nil, errors.New(errMsg))
}

func encodeSuccessResponse(codec codec.ServerCodec, id uint64, reply interface{}) []byte {
	return encodeResponse(codec, id, reply, nil)
}

func encodeResponse(codec codec.ServerCodec, id uint64, reply interface{}, err error) []byte {
	header := protocol.ResponsePool.Get().(*protocol.Header)
	defer protocol.ResponsePool.Put(header)

	header.ID = id
	if err != nil {
		header.Status = 1
	} else {
		header.Status = 0
	}

	body := &protocol.Body{}
	if reply != nil {
		data, _ := serializer.Serializers[serializer.JSON].Encode(reply)
		body.Payload = data
	}
	data, _ := codec.EncodeResponse(reply, header, body)

	return data
}
