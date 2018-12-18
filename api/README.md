# BTChain API文档  v0.1

## BTChain 特性
- BTChain 基于TM开发，BFT共识
- 采用ETH的帐户体系，使用地址和私钥，私钥Hex长64字节，地址Hex长42字节（包含0x前缀，其他任何地方使用时需要0x前缀）
`{"privkey":"c62949b9232fb7bf314c6156a4d33c56d956dbb3eadb61c94c5c23c4e1e707ea","address":"0xab96C28315D77b7786C8b973bb203896F75b07A3"}`
-

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
| Behavior | String | BU行为数据 |     <=256字节
| Memo | String | 备注 |   <=256字节

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

## 三、查询交易 by hash
```
GET
http://192.168.1.2:10000/v1/transactions/0x8bf191b30e45c60763f4b51965776dc45887c956b100bf7b0730fdaa25c0c33d
```
返回
```
{
	"isSuccess": true,
	"result": {
		"action_count": 1,
		"action_id": 0,
		"amount": 100,
		"block_hash": "0x6f51074d78307f7a91626f9067cf17865f1cd8fe3d73ae17f4011663c8596366",
		"block_height": 680,
		"created_at": 1545133165555261000,
		"dst": "0xa15d837e862cd9cacef81214617d7bb3b6270701",
		"jdata": "沧海一声笑,滔滔两岸潮,浮沉随浪只记今朝;苍天笑,纷纷世上潮,谁负谁胜出天知晓;江山笑,烟雨遥,涛浪淘尽红尘俗世几多娇;清风笑,竟惹寂寥,豪情还剩了一襟晚照;苍生笑,不再寂寥,豪情仍在痴痴笑笑",
		"memo": "",
		"nonce": 44888,
		"result_code": 0,
		"result_msg": "",
		"src": "0x061a060880bb4e5ad559350203d60a4349d3ecd6",
		"tx_hash": "0x8bf191b30e45c60763f4b51965776dc45887c956b100bf7b0730fdaa25c0c33d",
		"tx_id": 44889
	}
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
