package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"sync/atomic"
)

type AddressMgr struct {
	num      uint64
	idx      uint64
	addresss []string
}

func NewAddressMgr(num int) *AddressMgr {
	var mgr AddressMgr
	mgr.num = uint64(num)
	mgr.addresss = make([]string, num)
	for i := 0; i < num; i++ {
		privkey, _ := crypto.GenerateKey()
		mgr.addresss[i] = crypto.PubkeyToAddress(privkey.PublicKey).Hex()
	}
	return &mgr
}

func (p *AddressMgr) Get() string {
	cli := p.addresss[p.idx%p.num]
	atomic.AddUint64(&p.idx, 1)
	return cli
}
