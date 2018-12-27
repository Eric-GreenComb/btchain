// 定义对外开放的HTTP API
package bean

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/tendermint/go-amino"
	"strconv"
)

type Transaction struct {
	BaseFee string    //基础手续费 暂不支持
	Actions []*Action //请求时可以不排序，DAPP处理时按照ACTION ID 排ASC序
}

type Action struct {
	ID     int    // 最大255,必须从0开始填写顺序递增
	Src    string // src 地址
	Priv   string // src 私钥
	Dst    string // 目标帐户地址
	Amount string // 金额 无小数
	Data   string // 行为价值 最大长度256,大于256只存256字节
	Memo   string // 最大长度256,大于256只存256字节
	Time   string // 签名交易时需要此参数 RFC2399
	Sign   string // 签名信息（签名交易时必选）
}

func (p Transaction) String() string {
	var str string
	for _, v := range p.Actions {
		str = str + ":" + strconv.Itoa(v.ID)
	}
	return str
}

func (p Transaction) Len() int {
	return len(p.Actions)
}

func (p Transaction) Less(i, j int) bool {
	return p.Actions[i].ID < p.Actions[j].ID
}

func (p Transaction) Swap(i, j int) {
	p.Actions[i], p.Actions[j] = p.Actions[j], p.Actions[i]
}

func (p Transaction) Hash() (h ethcmn.Hash) {
	hw := sha3.NewKeccak256()
	b, _ := amino.MarshalBinaryBare(&p)
	hw.Write(b)
	hw.Sum(h[:0])
	return
}

type Response struct {
	IsSuccess bool        `json:"isSuccess"`
	Msg       string      `json:"msg"`
	Result    interface{} `json:"result"`
}

type TransResp struct {
	Tx          string       `json:"tx"`
	ActionResps []ActionResp `json:"actions"`
}

// 同步请求时
// 异步请求
type ActionResp struct {
	ID   int    `json:"id"`
	Code int    `json:"code"` //0-success -1:失败 1：HOLD
	Msg  string `json:"msg"`
}
