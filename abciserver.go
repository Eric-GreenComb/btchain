package btchain

import (
	"encoding/json"
	"fmt"
	"github.com/axengine/btchain/code"
	"github.com/axengine/btchain/define"
	"github.com/axengine/btchain/version"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	abciversion "github.com/tendermint/tendermint/version"
	"go.uber.org/zap"
	"log"
	"sort"
	"strconv"
)

func (app *BTApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	log.Println("=====>>InitChain")
	return abcitypes.ResponseInitChain{}
}

// Info
// TM Core启动时会查询chain信息 需要返回lastBlock相关信息，否则会从第一个块replay
func (app *BTApplication) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	if app.currentHeader.Height != 0 {
		resInfo.LastBlockHeight = int64(app.currentHeader.Height)
		resInfo.LastBlockAppHash = app.currentHeader.PrevHash.Bytes()
	}
	resInfo.Data = fmt.Sprintf("{\"size\":%v}", app.currentHeader.Height)
	resInfo.Version = abciversion.ABCIVersion
	resInfo.AppVersion = version.APPVersion
	app.logger.Info("ABCI Info", zap.Uint64("height", app.currentHeader.Height), zap.String("PrevHash", app.currentHeader.PrevHash.Hex()))
	return resInfo
}

// CheckTx
// 初步检查，如果check失败，将不会被打包
func (app *BTApplication) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	//app.logger.Debug("ABCI CheckTx", zap.ByteString("tx", tx[:]))
	var t define.Transaction
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		app.logger.Warn("rlp unmarshal", zap.Error(err), zap.ByteString("tx", tx))
		return abcitypes.ResponseCheckTx{Code: code.CodeTypeEncodingError, Log: "CodeTypeEncodingError"}
	}
	sort.Sort(t)
	app.logger.Debug("ABCI CheckTx", zap.String("tx", t.String()))

	//检查每个操作是否合法
	for i, action := range t.Actions {
		if i != int(action.ID) {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeOutOfOrder:"+strconv.Itoa(int(action.ID))))
			return abcitypes.ResponseCheckTx{Code: code.CodeOutOfOrder, Log: "CodeOutOfOrder:" + strconv.Itoa(int(action.ID))}
		}
		if !app.stateDup.state.Exist(action.From) {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeAccountNotFound:"+action.From.Hex()))
			return abcitypes.ResponseCheckTx{Code: code.CodeAccountNotFound, Log: "CodeAccountNotFound:" + action.From.Hex()}
		}

		balance := app.stateDup.state.GetBalance(action.From)
		if balance.Cmp(action.Amount) < 0 {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeNotEnoughMoney:"+balance.String()))
			return abcitypes.ResponseCheckTx{Code: code.CodeNotEnoughMoney, Log: "CodeNotEnoughMoney:" + balance.String()}
		}
	}

	//检查签名是否合法
	if err := t.CheckSig(); err != nil {
		app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeSignerFaild:"+err.Error()))
		return abcitypes.ResponseCheckTx{Code: code.CodeSignerFaild, Log: "CodeSignerFaild:" + err.Error()}
	}

	return abcitypes.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

// BeginBlock
// 区块开始，记录区块高度和hash
func (app *BTApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.logger.Info("ABCI BeginBlock", zap.Int64("height", req.Header.Height), zap.String("hash", ethcmn.Bytes2Hex(req.Hash)))
	app.tempHeader.Height = uint64(req.Header.Height)
	app.tempHeader.BlockHash = ethcmn.BytesToHash(req.Hash)
	return abcitypes.ResponseBeginBlock{}
}

func (app *BTApplication) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	var (
		t define.Transaction
	)
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		app.logger.Warn("rlp unmarshal", zap.Error(err), zap.ByteString("tx", tx))
		return abcitypes.ResponseDeliverTx{Code: code.CodeTypeEncodingError, Log: "CodeTypeEncodingError"}
	}
	sort.Sort(t)
	app.logger.Debug("ABCI DeliverTx", zap.String("tx", t.String()))

	//创建快照
	stateSnapshot := app.stateDup.state.Snapshot()

	app.tempHeader.TxCount = app.tempHeader.TxCount + 1
	app.tempHeader.OpCount = app.tempHeader.OpCount + uint64(t.Len())

	txHash := t.SigHash()
	actionCount := t.Len()
	for _, action := range t.Actions {
		//自动创建to
		if !app.stateDup.state.Exist(action.To) {
			app.stateDup.state.CreateAccount(action.To)
		}

		//取nonce
		nonce := app.stateDup.state.GetNonce(action.From)

		//必须再次校验余额
		balance := app.stateDup.state.GetBalance(action.From)
		if balance.Cmp(action.Amount) < 0 {
			app.stateDup.state.RevertToSnapshot(stateSnapshot)
			app.logger.Warn("ABCI DeliverTx", zap.String("err", "CodeNotEnoughMoney"), zap.String("from", action.From.Hex()), zap.String("amount", balance.String()))
			return abcitypes.ResponseDeliverTx{Code: code.CodeNotEnoughMoney, Log: "not enough money"}
		}

		//资金操作
		app.stateDup.state.SubBalance(action.From, action.Amount)
		app.stateDup.state.AddBalance(action.To, action.Amount)
		app.stateDup.state.SetNonce(action.From, nonce+1)

		var txData define.TransactionData
		txData.TxHash = txHash
		txData.BlockHeight = app.tempHeader.Height
		txData.BlockHash = app.tempHeader.BlockHash
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

	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK, Data: []byte(txHash.Hex())}
}

func (app *BTApplication) commitState() (ethcmn.Hash, error) {
	var (
		stateRoot ethcmn.Hash
		err       error
	)

	app.stateDup.lock.Lock()
	defer app.stateDup.lock.Unlock()

	// 更新stateRoot
	stateRoot = app.stateDup.state.IntermediateRoot(false)
	if _, err = app.stateDup.state.Commit(false); err != nil {
		return stateRoot, err
	}
	if err = app.stateDup.state.Database().TrieDB().Commit(stateRoot, true); err != nil {
		return stateRoot, err
	}

	return stateRoot, nil
}

func (app *BTApplication) Commit() abcitypes.ResponseCommit {
	var (
		stateRoot ethcmn.Hash
		err       error
	)

	// 更新stateRoot
	stateRoot, err = app.commitState()
	if err != nil {
		app.logger.Error("ABCI Commit state commit", zap.Error(err))
		app.SaveLastBlock(app.currentHeader.Hash(), app.currentHeader)
		return abcitypes.ResponseCommit{Data: app.currentHeader.Hash()}
	}

	//计算 appHash
	appHash := app.tempHeader.Hash()

	//保存lastblock
	app.SaveLastBlock(appHash, app.tempHeader)

	//更新currentHeader 保护现场
	app.currentHeader = app.tempHeader

	//SQLITE3保存记录
	if err := app.SaveDBData(); err != nil {
		app.logger.Error("ABCI SaveDBData", zap.Error(err), zap.Uint64("height", app.tempHeader.Height))
	}

	//  -------------准备新区块---------

	// 更新stateDup
	app.stateDup.lock.Lock()
	app.stateDup.state, err = state.New(stateRoot, state.NewDatabase(app.chainDb))
	app.stateDup.lock.Unlock()
	if err != nil {
		app.logger.Error("ABCI state.New", zap.Error(err), zap.String("stateRoot", stateRoot.Hex()))
		panic(err)
	}

	//清理区块执行现场
	app.blockExeInfo = &blockExeInfo{}
	app.tempHeader = &define.AppHeader{}

	//下个区块需要上个的hash
	app.tempHeader.PrevHash = ethcmn.BytesToHash(appHash)
	app.tempHeader.StateRoot = stateRoot

	app.logger.Info("ABCI Commit", zap.String("hash", ethcmn.BytesToHash(appHash).Hex()), zap.Uint64("height", app.tempHeader.Height))
	return abcitypes.ResponseCommit{Data: appHash}
}

func (app *BTApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	app.logger.Debug("ABCI Query:", zap.String("path", reqQuery.Path), zap.String("data", string(reqQuery.Data)),
		zap.Int64("Height", reqQuery.Height),
		zap.Bool("Prove", reqQuery.Prove))

	switch reqQuery.Path {
	case QUERY_TX:
		result := app.QueryTx(reqQuery.Data)
		b, err := rlp.EncodeToBytes(&result)
		if err != nil {
			return abcitypes.ResponseQuery{Code: code.CodeTypeEncodingError, Log: err.Error()}
		}
		return abcitypes.ResponseQuery{Value: b}
	case QUERY_ACCOUNT:
		result := app.QueryAccount(reqQuery.Data)
		b, err := rlp.EncodeToBytes(&result)
		if err != nil {
			return abcitypes.ResponseQuery{Code: code.CodeTypeEncodingError, Log: err.Error()}
		}
		return abcitypes.ResponseQuery{Value: b}
	default:
		app.logger.Warn("ABCI Query", zap.String("code", "CodeUnknownPath"))
		return abcitypes.ResponseQuery{Code: code.CodeUnknownPath, Log: "CodeUnknownPath"}
	}
	return abcitypes.ResponseQuery{Value: []byte(fmt.Sprintf("%v", app.currentHeader.PrevHash))}
}
