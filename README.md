# btchain
A block chain for BT of JVV

BTChain Baseed on tendermint 0.27.0-9c236ffd

## 编译链程序
- 先按照tendermint标准编译流程进行编译，确保单独编译tendermint可以成功
- 修改tendermint/proxy/client.go,func:DefaultClientCreator,加入以下分支
```
	case "bt":
		return NewLocalClientCreator(btchain.NewBTApplication())
```
导入包
```
"github.com/axengine/btchain"
```
- 确保gcc可用，因为要开启cgo
- 进入tendermint目录执行`make build_c`，生成的二进制文件在build目录
- 配置文件config.toml与可执行程序在同一目录
```
[genesis]
account = "0x061a060880BB4E5AD559350203d60a4349d3Ecd6"
amount = "10000000000"


[db]
type = "sqlite3"
path = "./data/"

[log]
env = "debug" # production
path = "./log/" #需先创建
```
- 日志目录，与可执行程序同级的log目录，需先创建

## 编译API
- 进入github.com/axengine/btchain/api 执行`go build`即可
- 配置文件，在可执行程序所在目录的config目录下
```
bind = ":10000"
rpc = "127.0.0.1:26657"
writable = true #false时只有查询API
isAdmin = true #true时有validator更新API

[log]
path = "./log/"
```