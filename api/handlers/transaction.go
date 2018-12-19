package handlers

import (
	"crypto/ecdsa"
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

func (hd *Handler) makeTx(ctx *gin.Context) (*define.Transaction, error) {
	var tdata bean.Transaction
	if err := ctx.BindJSON(&tdata); err != nil {
		return nil, err
	}
	sort.Sort(&tdata)

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

func (hd *Handler) SendTransactionsCommit(ctx *gin.Context) {
	tx, err := hd.makeTx(ctx)
	if err != nil {
		hd.logger.Error("makeTx", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if hd.chechTx(tx.SigHash().Hex()) {
		hd.responseWrite(ctx, false, "repeat tx")
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
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.CheckTx.Code != define.CodeType_OK {
		hd.responseWrite(ctx, false, result.CheckTx)
		return
	}
	if result.DeliverTx.Code != define.CodeType_OK {
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

	if hd.chechTx(tx.SigHash().Hex()) {
		hd.responseWrite(ctx, false, "repeat tx")
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
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
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

	if hd.chechTx(tx.SigHash().Hex()) {
		hd.responseWrite(ctx, false, "repeat tx")
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
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if result.Code != define.CodeType_OK {
		hd.responseWrite(ctx, false, result)
		return
	}
	hd.responseWrite(ctx, true, ethcmn.BytesToHash(result.Data).Hex())
}

// chckeTx return true if the hash exist
func (hd *Handler) chechTx(hash string) bool {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	_, ok := hd.cache.Get(hash)
	if !ok {
		hd.cache.Set(hash, hash)
	}
	return ok
}
