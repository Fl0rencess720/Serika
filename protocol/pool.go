package protocol

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{New: func() interface{} {
		return &Header{}
	}}
	ResponsePool = sync.Pool{New: func() interface{} {
		return &Header{}
	}}
}
