package btchain

import (
	"encoding/base64"
	"github.com/axengine/btchain/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
)

type SPOnce struct {
	flag uint32
	op   define.SpecialOP
}

func (app *BTApplication) SpecialOP(tx *define.Transaction) error {
	var op define.SpecialOP

	action := tx.Actions[0]
	if action == nil {
		return errors.New("CodeType_InvalidTx")
	}
	b, err := base64.StdEncoding.DecodeString(action.Data)
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return errors.New("error validator pubkey")
	}
	power, err := strconv.Atoi(action.Memo)
	if err != nil {
		return err
	}

	op.Type = "ed25519"
	op.PubKey = b
	op.Power = uint32(power)

	if atomic.LoadUint32(&app.sp.flag) == 1 {
		app.logger.Info("SpecialOP", zap.String("SP Flag", "== 1"))
		return errors.New("Locked")
	}

	app.sp.op = op
	atomic.StoreUint32(&app.sp.flag, 1)

	return nil
}
