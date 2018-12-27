package handlers

import (
	"crypto/ecdsa"
	"errors"
	"github.com/axengine/btchain/api/bean"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/big"
	"sort"
	"time"
)

func (hd *Handler) SendTransactionsCommit(ctx *gin.Context) {
	tx, err := hd.makeTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxCommit(b)
	if err != nil {
		hd.logger.Error("BroadcastTxCommit", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.CheckTx.Code != define.CodeType_OK {
		hd.logger.Info("CheckTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWrite(ctx, false, result.CheckTx)
		return
	}
	if result.DeliverTx.Code != define.CodeType_OK {
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWrite(ctx, false, result.DeliverTx)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.DeliverTx.Data).Hex())
}

func (hd *Handler) SendTransactionsAsync(ctx *gin.Context) {
	tx, err := hd.makeTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxAsync(b)
	if err != nil {
		hd.logger.Error("BroadcastTxAsync", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
		hd.logger.Info("BroadcastTxAsync", zap.Uint32("code", result.Code))
		hd.responseWrite(ctx, false, result)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.Hash).Hex())
}

func (hd *Handler) SendTransactionsSync(ctx *gin.Context) {
	tx, err := hd.makeTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxSync(b)
	if err != nil {
		hd.logger.Error("BroadcastTxSync", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
		hd.logger.Info("BroadcastTxSync", zap.Uint32("code", result.Code))
		hd.responseWrite(ctx, false, result)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.Data).Hex())
}

func (hd *Handler) SendSignedTransactionsCommit(ctx *gin.Context) {
	tx, err := hd.makeSignedTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxCommit(b)
	if err != nil {
		hd.logger.Error("SendSignedTransactionsCommit", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.CheckTx.Code != define.CodeType_OK {
		hd.logger.Info("CheckTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWrite(ctx, false, result.CheckTx)
		return
	}
	if result.DeliverTx.Code != define.CodeType_OK {
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWrite(ctx, false, result.DeliverTx)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.DeliverTx.Data).Hex())
}

func (hd *Handler) SendSignedTransactionsAsync(ctx *gin.Context) {
	tx, err := hd.makeSignedTx(ctx)
	if err != nil {
		hd.logger.Error("makeSignedTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxAsync(b)
	if err != nil {
		hd.logger.Error("SendSignedTransactionsAsync", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
		hd.logger.Info("SendSignedTransactionsAsync", zap.Uint32("code", result.Code))
		hd.responseWrite(ctx, false, result)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.Hash).Hex())
}

func (hd *Handler) SendSignedTransactionsSync(ctx *gin.Context) {
	tx, err := hd.makeSignedTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		hd.logger.Error("MarshalBinaryBare", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.BroadcastTxSync(b)
	if err != nil {
		hd.logger.Error("SendSignedTransactionsSync", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
		hd.logger.Info("SendSignedTransactionsSync", zap.Uint32("code", result.Code))
		hd.responseWrite(ctx, false, result)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.Data).Hex())
}

// chckeTx return true if ok
func (hd *Handler) chechTx(hash string) bool {
	b, err := hd.cache.SetExNx(hash, []byte(hash))
	if err != nil {
		hd.logger.Error("chechTx", zap.Error(err))
	}
	return b
}

func (hd *Handler) makeTx(ctx *gin.Context) (*define.Transaction, error) {
	var tdata bean.Transaction
	if err := ctx.BindJSON(&tdata); err != nil {
		return nil, err
	}
	sort.Sort(&tdata)

	if !hd.chechTx(tdata.Hash().Hex()) {
		return nil, errors.New("repeat tx:" + tdata.Hash().Hex())
	}

	ops := len(tdata.Actions)

	var tx define.Transaction
	tx.Actions = make([]*define.Action, ops)

	var privkeys []*ecdsa.PrivateKey
	for i, v := range tdata.Actions {
		var action define.Action
		action.ID = uint8(v.ID)
		action.CreatedAt = uint64(time.Now().UnixNano())
		action.Src = ethcmn.HexToAddress(v.Src)
		action.Dst = ethcmn.HexToAddress(v.Dst)
		action.Amount, _ = new(big.Int).SetString(v.Amount, 10)
		action.Data = v.Data
		action.Memo = v.Memo

		if len(v.Src) != 42 || len(v.Dst) != 42 {
			return nil, errors.New("err address")
		}

		if len(v.Priv) != 64 {
			return nil, errors.New("err privkey")
		}

		if action.Amount.Cmp(big.NewInt(0)) <= 0 {
			return nil, errors.New("err amount")
		}

		tx.Actions[i] = &action
		privkey, err := crypto.ToECDSA(ethcmn.Hex2Bytes(v.Priv))
		if err != nil {
			return nil, err
		}
		privkeys = append(privkeys, privkey)
	}
	//签名
	tx.Sign(privkeys)

	return &tx, nil
}

func (hd *Handler) makeSignedTx(ctx *gin.Context) (*define.Transaction, error) {
	var tdata bean.Transaction
	if err := ctx.BindJSON(&tdata); err != nil {
		return nil, err
	}
	sort.Sort(&tdata)

	if !hd.chechTx(tdata.Hash().Hex()) {
		return nil, errors.New("repeat tx:" + tdata.Hash().Hex())
	}

	ops := len(tdata.Actions)

	var tx define.Transaction
	tx.Actions = make([]*define.Action, ops)

	for i, v := range tdata.Actions {
		var action define.Action
		action.ID = uint8(v.ID)
		t, err := time.Parse(time.RFC3339, v.Time)
		if err != nil {
			return nil, err
		}
		// 5min是结合重复交易缓存来设定的 交易在redis的过期时间是5分钟
		if time.Since(t) > time.Minute*5 {
			return nil, errors.New(" expired tx:" + tdata.Hash().Hex())
		}
		action.CreatedAt = uint64(t.UnixNano())
		action.Src = ethcmn.HexToAddress(v.Src)
		action.Dst = ethcmn.HexToAddress(v.Dst)
		action.Amount, _ = new(big.Int).SetString(v.Amount, 10)
		action.Data = v.Data
		action.Memo = v.Memo
		copy(action.SignHex[:], ethcmn.Hex2Bytes(v.Sign))

		if len(v.Src) != 42 || len(v.Dst) != 42 {
			return nil, errors.New("err address")
		}

		if len(v.Priv) != 64 {
			return nil, errors.New("err privkey")
		}

		if action.Amount.Cmp(big.NewInt(0)) <= 0 {
			return nil, errors.New("err amount")
		}

		if len(v.Sign) != 130 {
			return nil, errors.New("err signature")
		}

		tx.Actions[i] = &action
	}

	return &tx, nil
}
