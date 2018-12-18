package btchain

import (
	"errors"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
)

const (
	QUERY_TX      = "/tx"
	QUERY_ACCOUNT = "/account"
)

var (
	ZERO_ADDRESS = ethcmn.Address{}
	ZERO_HASH    = ethcmn.Hash{}
)

//
//type Result struct {
//	Code int32  `json:"Code"`
//	Data []byte `json:"Data"`
//	Log  string `json:"Log"` // Can be non-deterministic
//}
//
//func NewResultOK(data []byte, log string) Result {
//	return Result{
//		Code: 0,
//		Data: data,
//		Log:  log,
//	}
//}
//
//func NewError(code int32, log string) Result {
//	return Result{
//		Code: code,
//		Log:  log,
//	}
//}

func (app *BTApplication) QueryTx(tx []byte) ([]byte, error) {
	var query define.TxQuery
	if err := rlp.DecodeBytes(tx, &query); err != nil {
		app.logger.Debug("rlp.DecodeBytes", zap.Error(err))
		return nil, err
	}

	if query.Account != ZERO_ADDRESS {
		result, err := app.dataM.QueryAccountTxs(&query.Account, query.Cursor, query.Limit, query.Order)
		b, err := rlp.EncodeToBytes(result)
		return b, err
	}
	if query.TxHash != ZERO_HASH {
		result, err := app.dataM.QuerySingleTx(&query.TxHash)
		b, err := rlp.EncodeToBytes(result)
		return b, err
	}

	result, err := app.dataM.QueryAllTxs(query.Cursor, query.Limit, query.Order)
	b, err := rlp.EncodeToBytes(result)
	return b, err
}

func (app *BTApplication) QueryAccount(from []byte) ([]byte, error) {
	address := ethcmn.BytesToAddress(from)

	if !app.stateDup.state.Exist(address) {
		return nil, errors.New("account not exist")
	}

	balance := app.stateDup.state.GetBalance(address)
	b, err := rlp.EncodeToBytes(balance)
	return b, err
}
