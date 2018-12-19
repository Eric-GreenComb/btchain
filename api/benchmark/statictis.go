package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	DefaultAnalyzer = NewAnalyze()
)

type State struct {
	idx uint64
	t   time.Duration
}

type Analyze struct {
	mu    sync.Mutex
	datas map[uint64]*State
}

func NewAnalyze() *Analyze {
	return &Analyze{
		datas: make(map[uint64]*State),
	}
}

func (p *Analyze) Add(s *State) {
	p.mu.Lock()
	p.datas[s.idx] = s
	p.mu.Unlock()
}

func (p *Analyze) String() string {
	count := len(p.datas)
	if count == 0 {
		return fmt.Sprintf("count=%s", count)
	}
	var (
		du       time.Duration
		max, min time.Duration
	)
	for _, v := range p.datas {
		if min == 0 {
			min = v.t
		}
		if max == 0 {
			max = v.t
		}

		du = du + v.t
		if v.t > max {
			max = v.t
		}
		if v.t < min {
			min = v.t
		}
	}

	return fmt.Sprintf("\n\t\tcount=%v max=%v min=%v avg=%v ms", count, max, min, du.Nanoseconds()/int64(count)/1000000)
}
