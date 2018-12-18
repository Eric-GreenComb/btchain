package handlers

import (
	"fmt"
	"github.com/axengine/btchain"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"strings"
)

func (hd *Handler) QueryAccount(ctx *gin.Context) {
	addressHex := ctx.Param("address")
	if !ethcmn.IsHexAddress(addressHex) {
		hd.responseWrite(ctx, false, fmt.Sprintf("Invalid address %s", addressHex))
		return
	}
	if strings.Index(addressHex, "0x") == 0 {
		addressHex = addressHex[2:]
	}

	result, err := hd.client.ABCIQuery(btchain.QUERY_ACCOUNT, []byte(addressHex))
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.responseWrite(ctx, true, result)
}
