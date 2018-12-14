package btchain

import (
	"encoding/json"
	"fmt"
	"github.com/axengine/btchain/code"
	"github.com/axengine/btchain/define"
	"github.com/axengine/btchain/version"
	"github.com/axengine/go-amino"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	abciversion "github.com/tendermint/tendermint/version"
	"log"
	"sort"
)

func (app *BTApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	log.Println("=====>>InitChain")
	return abcitypes.ResponseInitChain{}
}

func (app *BTApplication) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	log.Println("=====>Info height:", app.currentHeader.Height)
	return abcitypes.ResponseInfo{
		Data:       fmt.Sprintf("{\"size\":%v}", app.currentHeader.Height),
		Version:    abciversion.ABCIVersion,
		AppVersion: version.APPVersion,
	}
}

func (app *BTApplication) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	log.Println("=====>CheckTx tx:", tx)
	var t define.Transaction
	if err := amino.UnmarshalBinaryBare(tx, &t); err != nil {
		return abcitypes.ResponseCheckTx{Code: code.CodeTypeEncodingError, Log: "CodeTypeEncodingError"}
	}
	sort.Sort(t)

	//检查每个操作是否合法

	//检查帐户余额是否充足

	//检查签名是否合法

	return abcitypes.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *BTApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	log.Println("=====>BeginBlock h:", req.Header.Height, "hash:", ethcmn.Bytes2Hex(req.Hash))
	app.tempHeader.Height = uint64(req.Header.Height)
	return abcitypes.ResponseBeginBlock{}
}

func (app *BTApplication) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	log.Println("=====>DeliverTx tx:", tx)
	var (
		t define.Transaction
	)
	if err := amino.UnmarshalBinaryBare(tx, &t); err != nil {
		return abcitypes.ResponseDeliverTx{Code: code.CodeTypeEncodingError, Log: "CodeTypeEncodingError"}
	}
	sort.Sort(t)

	stateSnapshot := app.stateDup.state.Snapshot()

	app.tempHeader.TxCount = app.tempHeader.TxCount + 1
	app.tempHeader.OpCount = app.tempHeader.OpCount + uint64(t.Len())

	txHash := t.Hash()
	actionCount := t.Len()
	for _, action := range t.Actions {
		//do something
		if action.Type == 0 { //开户
			if app.stateDup.state.Exist(action.To) {
				app.stateDup.state.RevertToSnapshot(stateSnapshot)
				return abcitypes.ResponseDeliverTx{Code: code.CodeALREADY_EXIST, Log: "ACCOUNT ALREADY_EXIST"}
			}
			app.stateDup.state.CreateAccount(action.To)
			continue
		} else {
			//交易
			if action.Behavior.Direction == 1 {
				action.From, action.To = action.To, action.From
			}
			nonce := app.stateDup.state.GetNonce(action.From)

			balance := app.stateDup.state.GetBalance(action.From)
			if balance.Cmp(action.Amount) < 0 {
				app.stateDup.state.RevertToSnapshot(stateSnapshot)
				return abcitypes.ResponseDeliverTx{Code: code.CodeNotEnoughMoney, Log: "not enough money"}
			}

			app.stateDup.state.SubBalance(action.From, action.Amount)
			app.stateDup.state.AddBalance(action.To, action.Amount)

			var txData define.TransactionData
			txData.TxHash = txHash
			txData.BlockHeight = app.tempHeader.Height
			//txData.BlockHash = []byte("") //提交的时候才知道
			txData.ActionCount = uint32(actionCount)
			txData.ActionID = uint32(action.ID)
			txData.UID = action.From
			txData.RelatedUID = action.To
			txData.Nonce = nonce
			txData.Direction = action.Behavior.Direction
			txData.Amount = action.Amount
			txData.CreateAt = action.Behavior.GenAt
			b, _ := json.Marshal(&action.Behavior)
			txData.JData = string(b)

			app.blockExeInfo.txDatas = append(app.blockExeInfo.txDatas, &txData)
		}
	}

	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *BTApplication) Commit() abcitypes.ResponseCommit {

	var (
		stateRoot ethcmn.Hash
		err       error
	)

	// 更新stateRoot
	app.stateDup.lock.Lock()
	stateRoot, err = app.stateDup.state.Commit(false)
	app.stateDup.lock.Unlock()
	if err != nil {
		app.SaveLastBlock(app.currentHeader.Hash(), app.currentHeader)
		return abcitypes.ResponseCommit{Data: app.currentHeader.Hash()}
	}

	//提交时间
	//app.tempHeader.ClosedAt = time.Now()

	//计算 appHash
	appHash := app.tempHeader.Hash()

	log.Println("=====>Commit Hash:", ethcmn.BytesToHash(appHash).Hex(), " h:", app.tempHeader)

	//保存lastblock
	app.SaveLastBlock(appHash, app.tempHeader)

	//更新currentHeader 保护现场
	app.currentHeader = app.tempHeader

	//SQLITE3保存记录
	app.SaveDBData()

	//  -------------准备新区块---------

	// 更新stateDup
	app.stateDup.lock.Lock()
	app.stateDup.state, err = state.New(stateRoot, state.NewDatabase(app.chainDb))
	app.stateDup.lock.Unlock()

	//清理区块执行现场
	app.blockExeInfo = &blockExeInfo{}
	app.tempHeader = &define.AppHeader{}

	//下个区块需要上个的hash
	app.tempHeader.PrevHash = ethcmn.BytesToHash(appHash)
	app.tempHeader.StateRoot = stateRoot

	return abcitypes.ResponseCommit{Data: appHash}
}

func (app *BTApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	log.Printf("on Query data:%v height:%v path:%v prove:%v\n", string(reqQuery.Data), reqQuery.Height, reqQuery.Path, reqQuery.Prove)
	//if reqQuery.Prove {
	//	value := app.state.db.Get(prefixKey(reqQuery.Data))
	//	resQuery.Index = -1 // TODO make Proof return index
	//	resQuery.Key = reqQuery.Data
	//	resQuery.Value = value
	//	if value != nil {
	//		resQuery.Log = "exists"
	//	} else {
	//		resQuery.Log = "does not exist"
	//	}
	//	return
	//} else {
	//	resQuery.Key = reqQuery.Data
	//	value := app.state.db.Get(prefixKey(reqQuery.Data))
	//	resQuery.Value = value
	//	if value != nil {
	//		resQuery.Log = "exists"
	//	} else {
	//		resQuery.Log = "does not exist"
	//	}
	//	return
	//}
	return
}
