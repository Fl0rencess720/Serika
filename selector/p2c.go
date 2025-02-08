package selector

import (
	"hash/crc32"
	"math/rand/v2"
	"sync"

	"github.com/Fl0rencess720/Serika/server"
)

type P2CSelector struct {
	mutex   sync.Mutex
	hosts   []*host
	loadMap map[string]*host
}

type host struct {
	metadata *server.Metadata
	load     uint64
}

func newP2CSelector(ms []*server.Metadata) *P2CSelector {
	s := &P2CSelector{}
	s.loadMap = make(map[string]*host)
	for _, m := range ms {
		s.hosts = append(s.hosts, &host{
			metadata: m,
			load:     0,
		})
		s.loadMap[m.Address] = s.hosts[len(s.hosts)-1]
	}

	return s
}

func (s *P2CSelector) SelectService(key string) (*server.Metadata, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.hosts) == 0 {
		return nil, ErrNoMetadata
	}
	n1, n2 := s.hash(key)
	host := n2
	if s.loadMap[n1].load < s.loadMap[n2].load {
		host = n1
	}
	return &server.Metadata{Network: s.loadMap[host].metadata.Network, Address: host}, nil
}

const Salt = "&*##"

func (s *P2CSelector) hash(key string) (string, string) {
	var n1, n2 string
	if len(key) > 0 { // 请求IP为不为空的情况
		saltKey := key + Salt
		n1 = s.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(s.hosts))].metadata.Address
		n2 = s.hosts[crc32.ChecksumIEEE([]byte(saltKey))%uint32(len(s.hosts))].metadata.Address
		return n1, n2
	}
	// 请求IP为空的情况
	n1 = s.hosts[rand.IntN(len(s.hosts))].metadata.Address
	n2 = s.hosts[rand.IntN(len(s.hosts))].metadata.Address
	return n1, n2
}
