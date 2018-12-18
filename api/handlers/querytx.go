package handlers

import (
	"encoding/json"
	"github.com/axengine/btchain"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func (hd *Handler) QuerySingleTx(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	txhash := ctx.Param("txhash")

	if txhash == "" {
		hd.responseWrite(ctx, false, "param txhash is required")
		return
	}

	hd.queryTxs(ctx, "", "", txhash, cursor, limit, order)

}

func (hd *Handler) queryTxs(ctx *gin.Context, account, direction, txhash, cursor, limit, order string) {
	var err error
	var query define.TxQuery
	if len(cursor) != 0 {
		query.Cursor, err = strconv.ParseUint(cursor, 10, 0)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
	}
	if len(limit) != 0 {
		var tmplmt uint64
		tmplmt, err = strconv.ParseUint(limit, 10, 0)
		query.Limit = tmplmt
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
	}
	query.Order = order

	if account != "" {
		query.Account = ethcmn.HexToAddress(account)
		drcn, _ := strconv.Atoi(direction)
		query.Direction = uint8(drcn)
	}
	if txhash != "" {
		query.TxHash = ethcmn.HexToHash(txhash)
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.ABCIQuery(btchain.QUERY_TX, bys)
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

	log.Println("data.Data", string(data.Data))
	//resData := make(map[string]interface{}, 0)
	resData := make([]define.TransactionData, 0)
	if err := json.Unmarshal(data.Data, &resData); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.responseWrite(ctx, true, &resData)
}

func (hd *Handler) QueryTxs(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	hd.queryTxs(ctx, "", "", "", cursor, limit, order)
}

func (hd *Handler) QueryAccTxsByDirection(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	account := ctx.Param("address")
	direction := ctx.Param("direction")
	if account == "" {
		hd.responseWrite(ctx, false, "param account is required")
		return
	}

	hd.queryTxs(ctx, account, direction, "", cursor, limit, order)
}

func (hd *Handler) QueryAccTxs(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	account := ctx.Param("address")
	if account == "" {
		hd.responseWrite(ctx, false, "param account is required")
		return
	}

	hd.queryTxs(ctx, account, "", "", cursor, limit, order)
}
