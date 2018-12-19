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
	"testing"
	"time"
)

var (
	BASE_API_URL = "http://192.168.1.2:10000/v1/"
)

func Test_transaction(t *testing.T) {
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.Src = "0x061a060880BB4E5AD559350203d60a4349d3Ecd7"
	action.Dst = "0xA15d837e862cD9CacEF81214617D7Bb3B6270701"
	action.Amount = "100"
	action.Priv = "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"
	action.Data = "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑"

	tdata.Actions = append(tdata.Actions, &action)
	action2 := action
	action2.ID = 1
	action2.Data = "天苍苍野茫茫"
	tdata.Actions = append(tdata.Actions, &action2)

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

	ethcmn.HexToAddress("0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
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
	Payment(address.(string), "1000000000000")
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
