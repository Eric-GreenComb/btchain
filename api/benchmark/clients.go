package main

import (
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type ClientMgr struct {
	num     uint64
	idx     uint64
	clients []*http.Client
}

func NewClientMgr(num int) *ClientMgr {
	var mgr ClientMgr
	mgr.num = uint64(num)
	//mgr.clients = make([]*http.Client, num)
	for i := 0; i < num; i++ {
		mgr.clients = append(mgr.clients, newHttpClient())
	}
	return &mgr
}

func newHttpClient() *http.Client {
	c := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(60 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	return &c
	//return http.DefaultClient
}

func (p *ClientMgr) Get() *http.Client {
	cli := p.clients[p.idx%p.num]
	atomic.AddUint64(&p.idx, 1)
	return cli
}
