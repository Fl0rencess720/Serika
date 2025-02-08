package selector

import (
	"sync"

	"github.com/Fl0rencess720/Serika/server"
)

type RoundRobinSelector struct {
	mutex sync.Mutex
	index uint64
	ms    []*server.Metadata
}

func newRoundRobinSelector(ms []*server.Metadata) *RoundRobinSelector {
	return &RoundRobinSelector{ms: ms}
}

func (s *RoundRobinSelector) SelectService(_ string) (*server.Metadata, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.ms) == 0 {
		return nil, ErrNoMetadata
	}

	s.index = s.index % uint64(len(s.ms))
	m := s.ms[s.index]
	s.index++
	return m, nil
}
