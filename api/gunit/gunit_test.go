package gunit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/axengine/btchain/bean"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	BASE_API_URL = "http://192.168.8.144:10000/v1/"
)

func Test_transaction(t *testing.T) {
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.From = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
	action.To = "0xA15d837e862cD9CacEF81214617D7Bb3B6270701"
	action.Amount = "100"
	action.Priv = "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033"
	action.Behavior.GenAt = uint64(time.Now().UnixNano())
	action.Behavior.Direction = 1
	action.Behavior.NodeID = "NODE0000001"
	action.Behavior.PartnerID = "PARTNER000001"
	action.Behavior.BehaviorID = "12345678901212321"
	action.Behavior.OrderID = "IAMAORDERID000001233423000000000000001"
	action.Behavior.Memo = "FOR TEST"

	tdata.Actions = append(tdata.Actions, action)

	b, _ := json.Marshal(&tdata)

	resp, err := http.Post(BASE_API_URL+"transactions", "application/json", bytes.NewReader(b))
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
	//privateKey, _ := crypto.GenerateKey()
	//address := crypto.PubkeyToAddress(privateKey.PublicKey)
	//compress := crypto.CompressPubkey(&privateKey.PublicKey)
	//fmt.Println(address)

	//a8971729fbc199fb3459529cebcd8704791fc699d88ac89284f23ff8e7fca7d6
	//pubkey, err := crypto.DecompressPubkey(ethcmn.Hex2Bytes("02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//address := crypto.PubkeyToAddress(*pubkey)
	//fmt.Println(address.Hex())

	privkey, _ := crypto.GenerateKey()
	buff := make([]byte, 32)
	copy(buff[32-len(privkey.D.Bytes()):], privkey.D.Bytes())
	fmt.Println(ethcmn.Bytes2Hex(buff))
	fmt.Println(ethcmn.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey)))
	fmt.Println(crypto.PubkeyToAddress(privkey.PublicKey).Hex())

}
