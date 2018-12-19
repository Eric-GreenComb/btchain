package handlers

import (
	"encoding/base64"
	"github.com/axengine/btchain"
	"github.com/axengine/btchain/api/bean"
	"github.com/axengine/btchain/define"
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
	power, err := strconv.Atoi(op.Power)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if len(b) != 32 {
		hd.responseWrite(ctx, false, "Err address format")
		return
	}
	var query define.SpecialOP
	query.Type = "ed25519"
	query.Power = uint32(power)
	query.PubKey = b

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.ABCIQuery(btchain.SPECIAL_OP, bys)
	if err != nil {
		hd.logger.Error("SPECIAL_OP", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var data define.Result
	err = rlp.DecodeBytes(result.Response.Value, &data)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	hd.responseWrite(ctx, true, &data)
}
