# BTChain API文档  v0.1

## BTChain 特性
- 帐户格式：采用ETH的帐户体系，使用地址和私钥，私钥Hex长64字节，地址Hex长42字节（包含0x前缀，其他任何地方使用时需要0x前缀）
`{"privkey":"c62949b9232fb7bf314c6156a4d33c56d956dbb3eadb61c94c5c23c4e1e707ea","address":"0xab96C28315D77b7786C8b973bb203896F75b07A3"}`
- 并发：支持帐户并发，一个tx包含多个action时，要么全成功要么全失败；
- 顺序：一个tx包含多个action时，action.ID需从0开始递增，序列化为JSON不需要对JSON Array排序；
- 幂等：API接口会将JSON tx解析为特定的对象并计算HASH，相同tx无论何时都会有相同的hash，相同hash在300秒内不可重复；

## 一、生成帐户地址和私钥

`GET http://192.168.1.2:10000/v1/genkey`

返回:

```
{
	"isSuccess": true,
	"result": {
		"privkey": "c62949b9232fb7bf314c6156a4d33c56d956dbb3eadb61c94c5c23c4e1e707ea",
		"address": "0xab96C28315D77b7786C8b973bb203896F75b07A3"
	}
}
```

说明：genkey只是根据算法生成地址和私钥对，不在链上创建地址，当发生交易时如果dst地址不存在自动创建

## 二、交易
```
http://192.168.1.2:10000/v1/transactionsCommit 交易全部提交后返回
http://192.168.1.2:10000/v1/transactionsSync   (推荐)初步校验后返回（交易帐户余额等，如果一笔交易中有多个action，且src有多个action时在最终执行可能出现余额不足）
http://192.168.1.2:10000/v1/transactionsAsync  总是返回成功（区块链接收成功就返回成功）
```

```
{
	"BaseFee": "0",
	"Actions": [{
		"ID": 0,
		"Src": "0x061a060880BB4E5AD559350203d60a4349d3Ecd6",
		"Priv": "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033",
		"Dst": "0xA15d837e862cD9CacEF81214617D7Bb3B6270701",
		"Amount": "100",
		"Data": "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑",
		"Memo": ""
	}]
}
```


| 参数名      | 类型   | 说明           | 备注              |
| ----------- | ------ | -------------- | ----------------- |
| ID      | uint8 | ACTION序号       | 从0开始，最大255 |
| From | String | 操作源帐户地址 |                   |
| Priv      | String | 操作源帐户私钥       |  |
| To | String | 目标帐户地址 |       
| Amount      | String | 金额       | 必须>0 |
| Data | String | BU行为数据 |     <=256字节  |
| Memo | String | 备注 |   <=256字节|

### 返回

```json
{
	"isSuccess": true,
	"result": "0x8bf191b30e45c60763f4b51965776dc45887c956b100bf7b0730fdaa25c0c33d"
}
```
| 参数名      | 类型   | 说明           | 备注              |
| ----------- | ------ | -------------- | ----------------- |
| tx      | String | 交易Hash      |  |

## 二（1）、签名交易
提供对应私钥交易的三个方法
```
/v1/signedTransactionsCommit
/v1/signedTransactionsSync
/v1/signedTransactionsAsync
```

额外参数
- `Time` 使用RFC3339格式，请转为+8时区
- `Sign` 签名信息

签名算法
- sign = secp256k1(hash,priv)

密钥
priv = crypto.ToECDSA(Hex2Bytes(PrivHex))

编解码
- 方法 RLP
- 参数 (SignHex字段为全0)
```
type Transaction struct {
	Type    uint8     //0-默认
	Actions []*Action //有序的action 按照ID ASC序
}
type Action struct {
	ID        uint8          // 最大支持255笔交易
	CreatedAt uint64         // 时间
	Src       ethcmn.Address //
	Dst       ethcmn.Address //
	Amount    *big.Int       //
	Data      string         // 用户数据
	Memo      string         // 备注
	SignHex   [65]byte       // 签名 65 bytes
}
tosignByte = rlp([]interface{}{tx.Type,tx.Actions})
```
hash算法
```
hw := sha3.NewKeccak256()
hw.Write(tosignByte)
hash = hw.Sum(nil)
```


请求示例
```
{
	"BaseFee": "0",
	"Actions": [{
		"ID": 0,
		"Src": "0x061a060880BB4E5AD559350203d60a4349d3Ecd6",
		"Priv": "5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033",
		"Dst": "0xa7b6fB0e8a56d96A37C96796dcdfcA694387dfcA",
		"Amount": "10",
		"Data": "admin init",
		"Memo": "BT Account",
		"Time": "2018-12-27T10:01:33+08:00",
		"Sign": "f81bc913aad18014ce98261342d29630b49b7eb6c5be68d4421939b2cce39e0e27e4980adf7a4aa0ae76de47bbfad67a3be2a9969fbcb113655e25e444d32c9001"
	}, {
		"ID": 1,
		"Src": "0x7eb2b9686F0393A924772588eb915472F11Ea274",
		"Priv": "17fa8fcbf4d07bbf182c50c73bb5096ba82cfa1358437129240472153e4fbf6f",
		"Dst": "0xa7b6fB0e8a56d96A37C96796dcdfcA694387dfcA",
		"Amount": "10",
		"Data": "admin init",
		"Memo": "BT Account",
		"Time": "2018-12-27T10:01:33+08:00",
		"Sign": "a5c09fbf85e12462682b363a2e02da8922f4a6d1dcb25cae42d57a508240aa6903a84467eca889dce79600ab707f3689ec0abf941d536c90431dc028528bc46200"
	}]
}
```

## 三、查询交易 by hash
```
GET
http://192.168.1.2:10000/v1/transactions/0x80fcfcc7bcb4107fa4c47152983fdff717676c6d665c4e89146d41a3c764babe
```
返回
```
{
	"isSuccess": true,
	"result": [{
		"tx_id": 2,
		"tx_hash": "0x80fcfcc7bcb4107fa4c47152983fdff717676c6d665c4e89146d41a3c764babe",
		"block_height": 114,
		"block_hash": "0xcd3c4ae9a2cc9c80c6c64f238c9ada23a432e94de4bb5236df5ad66f87887c87",
		"action_count": 2,
		"action_id": 1,
		"src": "0x061a060880bb4e5ad559350203d60a4349d3ecd6",
		"dst": "0xa15d837e862cd9cacef81214617d7bb3b6270701",
		"nonce": 1,
		"amount": 100,
		"result_code": 0,
		"result_msg": "",
		"created_at": 1545208976752997413,
		"jdata": "天苍苍野茫茫",
		"memo": ""
	}, {
		"tx_id": 1,
		"tx_hash": "0x80fcfcc7bcb4107fa4c47152983fdff717676c6d665c4e89146d41a3c764babe",
		"block_height": 114,
		"block_hash": "0xcd3c4ae9a2cc9c80c6c64f238c9ada23a432e94de4bb5236df5ad66f87887c87",
		"action_count": 2,
		"action_id": 0,
		"src": "0x061a060880bb4e5ad559350203d60a4349d3ecd6",
		"dst": "0xa15d837e862cd9cacef81214617d7bb3b6270701",
		"nonce": 0,
		"amount": 100,
		"result_code": 0,
		"result_msg": "",
		"created_at": 1545208976752820733,
		"jdata": "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑",
		"memo": ""
	}]
}
```

## 四、查询帐户
```
GET
http://192.168.1.2:10000/v1/accounts/0x061a060880BB4E5AD559350203d60a4349d3Ecd6
```
返回

```
{
	"isSuccess": true,
	"result": {
		"address": "0x061a060880BB4E5AD559350203d60a4349d3Ecd6",
		"balance": "999999999999999999999999999999999999999999999999999999999999999999999999999999999999955012"
	}
}
```

## 五、查询所有交易
http://192.168.8.144:10000/v1/transactions
```
{
	"isSuccess": true,
	"result": [{
		"tx_id": 2,
		"tx_hash": "0x38bb606de3d285845a2429bfb8e86fec2a58835d92dc45f29fb3fb4961da33e5",
		"block_height": 21,
		"block_hash": "0xc534d4105b49deed99f32abfc7cee6ec9f336253ab3af51cb0df68c1138973bf",
		"action_count": 1,
		"action_id": 0,
		"src": "0x061a060880bb4e5ad559350203d60a4349d3ecd6",
		"dst": "0xa15d837e862cd9cacef81214617d7bb3b6270701",
		"nonce": 1,
		"amount": 100,
		"result_code": 0,
		"result_msg": "",
		"created_at": 1545162956435109209,
		"jdata": "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑",
		"memo": ""
	}, {
		"tx_id": 1,
		"tx_hash": "0xab047f0582aba32010ba7c701715f3963cd2fa0d98356a6d8a96cebea00d5745",
		"block_height": 13,
		"block_hash": "0xd773c64e1aca59cfe30041023d065aa877fa5c253f7aa8756016bdc266faab0b",
		"action_count": 1,
		"action_id": 0,
		"src": "0x061a060880bb4e5ad559350203d60a4349d3ecd6",
		"dst": "0xa15d837e862cd9cacef81214617d7bb3b6270701",
		"nonce": 0,
		"amount": 100,
		"result_code": 0,
		"result_msg": "",
		"created_at": 1545162948260890110,
		"jdata": "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑",
		"memo": ""
	}]
}
```

## 六、查询帐户交易
http://192.168.8.144:10000/v1/accounts/0x061a060880BB4E5AD559350203d60a4349d3Ecd6/transactions

说明：只查询0x061a060880BB4E5AD559350203d60a4349d3Ecd6作为src的数据

## 六、查询帐户交易 by 方向
- http://192.168.8.144:10000/v1/accounts/0x061a060880BB4E5AD559350203d60a4349d3Ecd6/transactions/0 同上，address作为src
- http://192.168.8.144:10000/v1/accounts/0x061a060880BB4E5AD559350203d60a4349d3Ecd6/transactions/1 address作为dst

## 七、错误处理
当发生错误时有两种数据返回：

- 由API层返回的错误 - 一般可以认为交易完全失败
```
{
	"isSuccess": false,
	"message": "unexpected EOF"
}
错误类型：API链接区块链之间网络失败、API内部处理失败
当使用transactionsCommit方法进行交易时，可能出现API等待区块链返回交易结果超时或网络失败
```

- 由区块链返回的错误
```
{
	"isSuccess": false,
	"message": {
		"code": 10,
		"log": "CodeType_AccountNotFound:0x061a060880bB4E5AD559350203D60A4349D3eCD7"
	}
}
```
区块链返回的错误代码
```
const (
	CodeType_OK uint32 = 0
	// General response codes, 0 ~ 99
	CodeType_InternalError     uint32 = 1  //内部错误
	CodeType_EncodingError     uint32 = 2  //编解码错误
	CodeType_BadNonce          uint32 = 3  //nonce错误
	CodeType_Unauthorized      uint32 = 4  //未授权
	CodeType_InsufficientFunds uint32 = 5  //资金不足
	CodeType_UnknownRequest    uint32 = 6  //未知请求
	CodeType_InvalidTx         uint32 = 7  //交易不合法
	CodeType_UnknownAccount    uint32 = 8  //未知帐户
	CodeType_AccountExist      uint32 = 9  //帐户已存在
	CodeType_AccountNotFound   uint32 = 10 //帐户不存在
	CodeType_OutOfOrder        uint32 = 11 //action顺序错误
	CodeType_UnknownError      uint32 = 12 //未知错误
	CodeType_SignerFaild       uint32 = 13 //签名错误
)
```