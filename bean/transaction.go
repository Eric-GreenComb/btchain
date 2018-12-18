// 定义对外开放的HTTP API
package bean

type Transaction struct {
	BaseFee string      //基础手续费 暂不支持
	Actions ActionSlice //请求时可以不排序，DAPP处理时按照ACTION ID 排ASC序
}

type Action struct {
	ID       int      //最大255,必须从1开始填写
	From     string   //from 公钥
	Priv     string   //from 私钥
	To       string   //目标帐户
	Amount   string   //金额 无小数
	Behavior Behavior //行为价值
}

type Behavior struct {
	GenAt      uint64 `json:"created_at"`  //行为发生时间
	OrderID    string `json:"order_id"`    //订单ID
	NodeID     string `json:"node_id"`     //节点ID
	PartnerID  string `json:"partner_id"`  //商户ID
	BehaviorID string `json:"behavior_id"` //行为ID
	Direction  uint8  `json:"behavior_id"` //
	Memo       string `json:"memo"`        //备注
}

type ActionSlice []Action

func (p ActionSlice) Len() int {
	return len(p)
}
func (p ActionSlice) Less(i, j int) bool {
	return p[i].ID < p[j].ID
}
func (p ActionSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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
