package btchain

import (
	"github.com/axengine/btchain/define"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"sync/atomic"
)

type SPOnce struct {
	flag uint32
	op   define.SpecialOP
}

func (app *BTApplication) SpecialOP(tx []byte) define.Result {
	var op define.SpecialOP
	if err := rlp.DecodeBytes(tx, &op); err != nil {
		app.logger.Error("rlp.DecodeBytes", zap.Error(err))
		return define.NewError(define.CodeType_EncodingError, err.Error())
	}

	if atomic.LoadUint32(&app.sp.flag) == 1 {
		app.logger.Info("SpecialOP", zap.String("SP Flag", "== 1"))
		return define.NewError(define.CodeType_InternalError, "SP Flag==1,wait a moment")
	}

	app.sp.op = op
	atomic.StoreUint32(&app.sp.flag, 1)

	return MakeResultData("waiting 2 height...")
}
