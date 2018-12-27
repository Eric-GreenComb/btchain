package gunit

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/alex023/httpcli"
	"github.com/axengine/btchain/api/bean"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
	"testing"
	"time"
)

var (
	BASE_API_URL = "http://192.168.8.145:10000/v1/"
	//BASE_API_URL = "https://btapi.ibdt.tech/v1/"
)

func Test_transaction(t *testing.T) {
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.Src = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
	action.Dst = "0xEBC5D91e9b3c8ea8194ec6b0A63ce1548F9eA448"
	action.Amount = "1"
	action.Priv = "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"
	action.Data = "admin init"
	action.Memo = "BT Account"

	tdata.Actions = append(tdata.Actions, &action)

	b, _ := json.Marshal(&tdata)

	fmt.Println(string(b))

	resp, err := http.Post(BASE_API_URL+"transactionsCommit", "application/json", bytes.NewReader(b[:]))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

}

func Test_signedTransaction(t *testing.T) {
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.Src = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
	action.Dst = "0xa7b6fB0e8a56d96A37C96796dcdfcA694387dfcA"
	action.Amount = "10"
	action.Priv = "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"
	action.Data = "admin init"
	action.Memo = "BT Account"
	action.Time = time.Now().Format(time.RFC3339)

	action2 := action
	action2.ID = 1
	action2.Src = "0x7eb2b9686F0393A924772588eb915472F11Ea274"
	action2.Priv = "17fa8fcbf4d07bbf182c50c73bb5096ba82cfa1358437129240472153e4fbf6f"
	tdata.Actions = append(tdata.Actions, &action)
	tdata.Actions = append(tdata.Actions, &action2)

	tx, err := signtx(&tdata)
	if err != nil {
		t.Fatal("sign failed")
	}

	for i, v := range tx.Actions {
		signature := ethcmn.Bytes2Hex(v.SignHex[:])
		tdata.Actions[i].Sign = signature
	}

	b, _ := json.Marshal(&tdata)

	fmt.Println(string(b))

	resp, err := http.Post(BASE_API_URL+"signedTransactionsCommit", "application/json", bytes.NewReader(b[:]))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

}

func signtx(tdata *bean.Transaction) (*define.Transaction, error) {
	sort.Sort(tdata)

	ops := len(tdata.Actions)

	var tx define.Transaction

	tx.Actions = make([]*define.Action, ops)
	var privkeys []*ecdsa.PrivateKey
	for i, v := range tdata.Actions {
		var action define.Action
		action.ID = uint8(v.ID)
		t, err := time.Parse(time.RFC3339, v.Time)
		if err != nil {
			return &tx, err
		}
		action.CreatedAt = uint64(t.UnixNano())
		action.Src = ethcmn.HexToAddress(v.Src)
		action.Dst = ethcmn.HexToAddress(v.Dst)
		action.Amount, _ = new(big.Int).SetString(v.Amount, 10)
		action.Data = v.Data
		action.Memo = v.Memo
		copy(action.SignHex[:], ethcmn.Hex2Bytes(v.Sign))

		tx.Actions[i] = &action

		privkey, err := crypto.ToECDSA(ethcmn.Hex2Bytes(v.Priv))
		if err != nil {
			return &tx, err
		}
		privkeys = append(privkeys, privkey)
	}
	tx.Sign(privkeys)
	return &tx, nil
}

func Test_address(t *testing.T) {
	var tx define.Transaction
	var action define.Action
	var privkeys []*ecdsa.PrivateKey

	action.CreatedAt = uint64(time.Now().UnixNano())
	action.ID = 0
	action.Src = ethcmn.HexToAddress("0x061a060880BB4E5AD559350203d60a4349d3Ecd6")
	action.Dst = ethcmn.HexToAddress("0x061a060880BB4E5AD559350203d60a4349d3Ecd7")
	action.Amount, _ = new(big.Int).SetString("1", 10)
	action.Data = "hello world"
	action.Memo = "for test"
	tx.Actions = append(tx.Actions, &action)

	privkey, err := crypto.ToECDSA(ethcmn.Hex2Bytes("5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"))
	if err != nil {
		t.Fatal(err)
	}
	privkeys = append(privkeys, privkey)

	tx.Sign(privkeys)

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		t.Fatal(err)
	}
	begin := time.Now()
	client := client.NewHTTP("192.168.8.144:26657", "/websocket")
	//result, err := client.BroadcastTxAsync(b)  //mempool.checkTx
	result, err := client.BroadcastTxSync(b) //checkTx
	//result, err := client.BroadcastTxCommit(b) //commit
	if err != nil {
		t.Fatal(err)
	}
	log.Println("use", time.Since(begin))
	log.Println("code", result.Code)
	log.Println("data", ethcmn.BytesToHash(result.Data).Hex())
	log.Println("log", result.Log)
	log.Println("hash", ethcmn.BytesToHash(result.Hash).Hex())
}

func Test_genaccount(t *testing.T) {
	c := httpcli.Get(BASE_API_URL + "genkey")
	resp, err := c.Response()
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	resp.Json(&data)
	log.Println(data)
	resultx, _ := data["result"]
	result := resultx.(map[string]interface{})
	address, _ := result["address"]
	privkey, _ := result["privkey"]
	Payment(address.(string), "10000")
	log.Println("address:", address.(string), " priv:", privkey.(string))
}

func Test_Payment(t *testing.T) {
	if err := Payment("0x061a060880BB4E5AD559350203d60a4349d3Ecd6", "1000000"); err != nil {
		t.Fatal(err)
	}
}

func Payment(to, amount string) error {
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.Src = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
	action.Priv = "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"
	action.Dst = to
	action.Amount = amount
	action.Data = "admin payment"
	tdata.Actions = append(tdata.Actions, &action)
	b, _ := json.Marshal(&tdata)
	resp, err := http.Post(BASE_API_URL+"transactionsCommit", "application/json", bytes.NewReader(b[:]))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
