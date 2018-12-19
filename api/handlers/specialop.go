package handlers

import (
	"encoding/base64"
	"github.com/axengine/btchain/api/bean"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func (hd *Handler) SpecialOP(ctx *gin.Context) {
	var op bean.SpecilOP
	if err := ctx.BindJSON(&op); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	b, err := base64.StdEncoding.DecodeString(op.Pubkey)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if len(b) != 32 {
		hd.responseWrite(ctx, false, "Err address format")
		return
	}
	_, err = strconv.Atoi(op.Power)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var tx define.Transaction
	var action define.Action
	action.Data = op.Pubkey
	action.Memo = op.Power

	tx.Type = 1
	tx.Actions = append(tx.Actions, &action)

	b, err = rlp.EncodeToBytes(&tx)
	if err != nil {
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
