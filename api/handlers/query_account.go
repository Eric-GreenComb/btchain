package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/axengine/btchain"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	var data define.Result
	err = rlp.DecodeBytes(result.Response.Value, &data)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	resData := make(map[string]interface{}, 0)
	if err := json.Unmarshal(data.Data, &resData); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		hd.logger.Error("resData", zap.String("data", string(data.Data)))
		return
	}
	hd.responseWrite(ctx, true, &resData)
}
