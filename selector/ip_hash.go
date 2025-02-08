package selector

import (
	"hash/crc32"
	"sync"

	"github.com/Fl0rencess720/Serika/server"
)

type IPHashSelector struct {
	mutex sync.Mutex
	ms    []*server.Metadata
}

func newIPHashSelector(ms []*server.Metadata) *IPHashSelector {
	return &IPHashSelector{ms: ms}
}

func (s *IPHashSelector) SelectService(key string) (*server.Metadata, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.ms) == 0 {
		return nil, ErrNoMetadata
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(s.ms))
	return s.ms[value], nil
}
