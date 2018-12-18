package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"net/http"
	"sync/atomic"
	"testing"
)

// 1 用genkey生成10000个帐户
// 2 并发使用创世帐户对10000个帐户转帐
var (
	succ uint64
	fail uint64
)

func Benchmark_transactions(t *testing.B) {
	clients := make([]*http.Client, 100)
	for i := 0; i < 100; i++ {
		clients[i] = newHttpClient()
	}

	addrs := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		privkey, _ := crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(privkey.PublicKey).Hex()
	}

	for i := 0; i < 10000; i++ {
		cli := clients[i%100]
		if err := transaction(cli, "0x061a060880BB4E5AD559350203d60a4349d3Ecd6", addrs[i%10000], "1", "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"); err != nil {
			atomic.AddUint64(&fail, 1)
		} else {
			atomic.AddUint64(&succ, 1)
		}
	}
	log.Println("succ:", succ, " fail:", fail)
}
