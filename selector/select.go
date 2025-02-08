package selector

import "github.com/Fl0rencess720/Serika/server"

type SelectMode uint8

const (
	Random SelectMode = iota
	RoundRobin
	P2C
	IPHash
)

type Selector interface {
	SelectService(string) (*server.Metadata, error)
}

func NewSelector(mode SelectMode, ms []*server.Metadata) Selector {
	switch mode {
	case Random:
		return newRandomSelector(ms)
	case RoundRobin:
		return newRoundRobinSelector(ms)
	case P2C:
		return newP2CSelector(ms)
	case IPHash:
		return newIPHashSelector(ms)
	default:
		return newRandomSelector(ms)
	}
}
