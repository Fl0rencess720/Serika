package selector

import (
	"errors"
	"math/rand/v2"
	"sync"

	"github.com/Fl0rencess720/Serika/server"
)

var ErrNoMetadata = errors.New("there is no metadata in slice")

type RandomSelector struct {
	mutex sync.Mutex
	ms    []*server.Metadata
}

func newRandomSelector(ms []*server.Metadata) *RandomSelector {
	return &RandomSelector{ms: ms}
}

func (s *RandomSelector) SelectService(string) (*server.Metadata, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.ms) == 0 {
		return nil, ErrNoMetadata
	}
	randInt := rand.IntN(len(s.ms))
	return s.ms[randInt], nil
}
