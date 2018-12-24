package block

import (
	"testing"

	"github.com/astaxie/beego"
)

func init() {
	beego.AppConfig.Set("chainaddr", "http://192.168.40.128:46657")
}

func TestGetStatus(t *testing.T) {
	s, _ := GetStatus()
	t.Log(s.NodeInfo.NetWork, s.LatestBlockHeight)
	return
}

func TestGetBlock(t *testing.T) {
	s, err := GetBlock(5991)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.Block.AppHash, s.BlockMeta.Hash)
	return
}

func TestGetBlockList(t *testing.T) {

	bls, err := GetBlockChainList(10, 20)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bls {
		t.Log(v.AppHash, v.Hash, v.NumTxs, len(v.Txs))
	}

	return
}

func TestGetLastNumChain(t *testing.T) {
	bls, err := GetLastNumChain(10)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bls {
		t.Log(v.Height, v.Hash, v.NumTxs, len(v.Txs))
	}

	return

}
