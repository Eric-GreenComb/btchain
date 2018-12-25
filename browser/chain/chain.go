package chain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axengine/btchain/browser/log"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Result struct {
	IsSuccess bool  `json:"isSuccess"`
	TxDs      []TxD `json:"result"`
}

type TxD_Ex struct {
	TxD
	TimeStr string
}

type TxD struct {
	TxID        uint64 `json:"tx_id"`        //交易ID
	TxHash      string `json:"tx_hash"`      //交易HASH - 重复交易不会被处理
	BlockHeight uint64 `json:"block_height"` //区块高度
	BlockHash   string `json:"block_hash"`   //区块HASH
	ActionCount uint32 `json:"action_count"` //一笔交易多个action
	ActionID    uint32 `json:"action_id"`    //action id
	Src         string `json:"src"`          //用户ID (if dir==0,uid 表示转入方，否则表示转出方)
	Dst         string `json:"dst"`          //关联的用户ID
	Nonce       uint64 `json:"nonce"`        //对应操作源帐户(转出方)NONCE
	Amount      uint64 `json:"amount"`       //金额
	ResultCode  uint   `json:"result_code"`  //应答码 0-success
	ResultMsg   string `json:"result_msg"`   //应答消息
	CreateAt    uint64 `json:"created_at"`   //入库时间
	JData       string `json:"jdata"`        //数据部分 建议JSON序列化
	Memo        string `json:"memo"`         //交易备注
}

func get(url string) (body []byte, err error) {
	var (
		rsp *http.Response
	)
	if rsp, err = http.Get(url); err != nil {
		return
	}
	if rsp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http status not ok :%v,%v", url, rsp.StatusCode))
		return
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, rsp.Body)
	body = buf.Bytes()
	return
}

func GetTxByHash(host, hash string) (result *Result, err error) {

	url := host + "/v1/transactions/" + hash

	body, err := get(url)
	if err != nil {
		return
	}
	result = new(Result)

	if err = json.Unmarshal(body, result); err != nil {
		return
	}
	return

}

func GetTxByAccount(host, account, dir string) (result *Result, err error) {
	url := host + "/v1/accounts/" + account + "/transactions"
	if dir == "income" {
		url = url + "/1"
	}
	body, err := get(url)
	if err != nil {
		return
	}
	result = new(Result)

	if err = json.Unmarshal(body, result); err != nil {
		return
	}
	return

}

func GetTx(host, order string, cursor, limit uint64) (result *Result, err error) {

	url := fmt.Sprintf("%v/v1/transactions?cursor=%v&limit=%v&order=%v", host, cursor, limit, order)

	body, err := get(url)
	if err != nil {
		return
	}
	result = new(Result)

	if err = json.Unmarshal(body, result); err != nil {
		log.Logger.Error("GetTx", zap.Error(err))
		return
	}
	return
}
