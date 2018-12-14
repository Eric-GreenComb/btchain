package define

import (
	"github.com/axengine/go-amino"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"math/big"
)

// TransactionData 交易表
// 一次HTTP请求算一笔交易，拥有唯一的交易HASH
// 单笔交易可以携带多个ACTION，每次ACTION有一个ID，ACTIONID由调用方指定，注意JSON Array没有顺序
// 每个ACTION保护出让方和受让方，“创建帐户”不是ACTION
// 单笔交易多个ACTION时可以部分失败，部分失败返回{成功的action列表，失败的action列表}
type TransactionData struct {
	TxID        uint64         `json:"tx_id"`        //交易ID
	TxHash      ethcmn.Hash    `json:"tx_hash"`      //交易HASH - 重复交易不会被处理
	BlockHeight uint64         `json:"block_height"` //区块高度
	BlockHash   ethcmn.Hash    `json:"block_hash"`   //区块HASH
	ActionCount uint32         `json:"action_count"` //一笔交易多个action
	ActionID    uint32         `json:"action_id"`    //action id
	UID         ethcmn.Address `json:"uid"`          //用户ID (if dir==0,uid 表示转入方，否则表示转出方)
	RelatedUID  ethcmn.Address `json:"related_uid"`  //关联的用户ID
	Direction   uint8          `json:"direction"`    //方向
	Nonce       uint64         `json:"nonce"`        //对应操作源帐户(转出方)NONCE
	Amount      *big.Int       `json:"amount"`       //金额
	ResultCode  uint           `json:"result_code"`  //应答码 0-success
	ResultMsg   string         `json:"result_msg"`   //应答消息
	CreateAt    uint64         `json:"created_at"`   //入库时间
	JData       string         `json:"jdata"`        //行为描述 json of struct Behavior
	Memo        string         `json:"memo"`         //交易备注
}

type Behavior struct {
	GenAt      uint64 `json:"created_at"`  //行为发生时间
	OrderID    uint32 `json:"order_id"`    //订单ID
	NodeID     uint32 `json:"node_id"`     //节点ID
	PartnerID  uint32 `json:"partner_id"`  //商户ID
	BehaviorID uint32 `json:"behavior_id"` //行为ID
	Direction  uint8  `json:"direction"`   //行为方向 0:FROM->TO 1:TO->FROM
}

// 通过Transaction 计算txhash
type Transaction struct {
	Actions []Action //有序的action 按照ID ASC序
}

// todo 待测试
func (tx *Transaction) Hash() (h ethcmn.Hash) {
	hw := sha3.NewKeccak256()
	b, _ := amino.MarshalBinaryBare(tx)
	hw.Write(b)
	hw.Sum(h[:0])
	return
}

type Action struct {
	ID       uint8 //最大支持255笔交易
	Type     uint8 //操作类型 0-开户 1-交易
	From     ethcmn.Address
	To       ethcmn.Address
	Amount   *big.Int
	Behavior Behavior
	SignHex  [65]byte //65 bytes
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
