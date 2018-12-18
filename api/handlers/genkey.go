package handlers

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func (hd *Handler) GenKey(ctx *gin.Context) {
	privkey, err := crypto.GenerateKey()
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var result struct {
		Privkey string `json:"privkey"`
		Address string `json:"address"`
	}

	buff := make([]byte, 32)
	copy(buff[32-len(privkey.D.Bytes()):], privkey.D.Bytes())
	result.Privkey = ethcmn.Bytes2Hex(buff)
	result.Address = ethcmn.Bytes2Hex(crypto.FromECDSAPub(&privkey.PublicKey))

	hd.responseWrite(ctx, true, result)
}
