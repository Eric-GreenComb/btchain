package handlers

import (
	"crypto/ecdsa"
	"github.com/axengine/btchain/bean"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/big"
	"time"
)

func (hd *Handler) SendTransactions(ctx *gin.Context) {
	var tdata bean.Transaction
	if err := ctx.BindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	ops := len(tdata.Actions)

	var tx define.Transaction
	tx.Actions = make([]*define.Action, ops)

	var privkeys []*ecdsa.PrivateKey
	for i, v := range tdata.Actions {
		var action define.Action
		action.ID = uint8(v.ID)
		action.Time = time.Now()
		action.From = ethcmn.HexToAddress(v.From)
		action.To = ethcmn.HexToAddress(v.To)
		action.Amount, _ = new(big.Int).SetString(v.Amount, 10)
		action.Behavior.GenAt = v.Behavior.GenAt
		copy(action.Behavior.OrderID[:], []byte(v.Behavior.OrderID))
		copy(action.Behavior.NodeID[:], []byte(v.Behavior.NodeID))
		copy(action.Behavior.PartnerID[:], []byte(v.Behavior.PartnerID))
		copy(action.Behavior.BehaviorID[:], []byte(v.Behavior.BehaviorID))
		action.Behavior.Direction = v.Behavior.Direction
		copy(action.Behavior.Memo[:], []byte(v.Behavior.Memo))

		tx.Actions[i] = &action
		privkey, _ := crypto.ToECDSA(ethcmn.Hex2Bytes(v.Priv))
		privkeys = append(privkeys, privkey)
	}

	//签名
	tx.Sign(privkeys)

	b, err := rlp.EncodeToBytes(tx)
	//b, err := json.Marshal(tx)
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

	hd.responseWrite(ctx, true, string(result.DeliverTx.Data))
}
