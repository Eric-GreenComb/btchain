package define

import (
	"crypto/ecdsa"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"testing"
	"time"
)

func TestTransaction_Sign(t *testing.T) {
	tx := new(Transaction)
	var action Action
	action.Time = time.Now()
	action.ID = 0
	action.Type = 0 //开户

	privkey, _ := crypto.ToECDSA(ethcmn.Hex2Bytes("94ee3a554a8ba2e2b6d08dc8a1597cb21a5101519732e0a6c39cdb606ea07e80"))

	action.From = ethcmn.HexToAddress("021838d5bfcc571d47a9e008306620a4871105fbc13aa6cfcbf917ffccc2a28761")
	action.To = ethcmn.HexToAddress("03525411c4d3943c02d87af4836fa97846796c69552c98fa25715588c2d0371902")

	tx.Actions = append(tx.Actions, &action)

	var privkeys []*ecdsa.PrivateKey
	privkeys = append(privkeys, privkey)
	if err := tx.Sign(privkeys); err != nil {
		t.Fatal("sign err", err)
	}

	log.Println("after sign:", tx.Actions[0])

	if err := tx.CheckSig(); err != nil {
		t.Fatal(err)
	}

}
