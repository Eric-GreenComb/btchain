package controllers

import (
	"github.com/astaxie/beego"
	"github.com/axengine/btchain/browser/chain"
	"github.com/axengine/btchain/browser/datamanage"
	"github.com/axengine/btchain/browser/log"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DisplayNum = 25
	//BlockHashLen = 40
	TxHashLen  = 64
	AccountLen = 40
	TxLimit    = 25
)

type WebController struct {
	beego.Controller
	Logger *zap.Logger
}

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (wc *WebController) Index() {
	wc.Layout = "layout.html"
	wc.TplName = "web.tpl"
}

func (wc *WebController) Latest() {
	blocks, err := datamanage.GetBlock(DisplayNum)
	if err != nil {
		log.Logger.Error("datamanage.GetBlock", zap.Error(err))
		return
	}
	wc.Data["json"] = &Result{
		Success: true,
		Data:    blocks,
	}
	wc.ServeJSON()
}

func (wc *WebController) Action() {

	var (
		actions []datamanage.Action
		err     error
	)

	stxId := wc.GetString(":txid")

	log.Logger.Debug("Action /trans/txid/", zap.String("txid", stxId))

	txId, _ := strconv.ParseUint(stxId, 10, 64)

	wc.Layout = "layout.html"
	wc.TplName = "action.tpl"

	if actions, err = datamanage.Getactions(txId); err != nil {
		goto errDeal
	}

	wc.Data["TxID"] = stxId
	wc.Data["Actions"] = actions
	return
errDeal:
	log.Logger.Error("Action", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

// tx hash
func (wc *WebController) TxByTxHash() {
	var (
		err    error
		result *chain.Result
		datas  []chain.TxD_Ex
	)
	wc.Layout = "layout.html"
	wc.TplName = "search_txs.tpl"
	//timeLayout := "2006-01-02 15:04:05"

	hash := strings.TrimPrefix(wc.GetString(":hash"), "0x")

	result, err = chain.GetTxByHash(beego.AppConfig.String("chain_api"), hash)
	if err != nil {
		goto errDeal
	}

	for _, v := range result.TxDs {
		var x chain.TxD_Ex
		x.TxID = v.TxID
		x.TxHash = v.TxHash
		x.BlockHeight = v.BlockHeight
		x.BlockHash = v.BlockHash
		x.ActionCount = v.ActionCount
		x.ActionID = v.ActionID
		x.Src = v.Src
		x.Dst = v.Dst
		x.Nonce = v.Nonce
		x.Amount = v.Amount
		x.ResultCode = v.ResultCode
		x.ResultMsg = v.ResultMsg
		x.CreateAt = v.CreateAt
		x.JData = v.JData
		x.Memo = v.Memo
		x.TimeStr = time.Unix(int64(x.CreateAt)/1e9, 0).Format(time.RFC3339)
		datas = append(datas, x)
	}

	wc.Data["TransDetail"] = datas
	return
errDeal:
	log.Logger.Error("TxByTxHash", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

func (wc *WebController) TxInByAddress() {
	var (
		err    error
		result *chain.Result
		datas  []chain.TxD_Ex
	)
	wc.Layout = "layout.html"
	wc.TplName = "search_txs.tpl"
	//timeLayout := "2006-01-02 15:04:05"

	address := strings.TrimPrefix(wc.GetString(":address"), "0x")
	result, err = chain.GetTxByAccount(beego.AppConfig.String("chain_api"), address, "income")
	if err != nil {
		goto errDeal
	}
	for _, v := range result.TxDs {
		var x chain.TxD_Ex
		x.TxID = v.TxID
		x.TxHash = v.TxHash
		x.BlockHeight = v.BlockHeight
		x.BlockHash = v.BlockHash
		x.ActionCount = v.ActionCount
		x.ActionID = v.ActionID
		x.Src = v.Src
		x.Dst = v.Dst
		x.Nonce = v.Nonce
		x.Amount = v.Amount
		x.ResultCode = v.ResultCode
		x.ResultMsg = v.ResultMsg
		x.CreateAt = v.CreateAt
		x.JData = v.JData
		x.Memo = v.Memo
		x.TimeStr = time.Unix(int64(x.CreateAt)/1e9, 0).Format(time.RFC3339)
		datas = append(datas, x)
	}
	wc.Data["TransDetail"] = datas
	return
errDeal:
	log.Logger.Error("TxInByAddress", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

func (wc *WebController) TxOutByAddress() {
	var (
		err    error
		result *chain.Result
		datas  []chain.TxD_Ex
	)
	wc.Layout = "layout.html"
	wc.TplName = "search_txs.tpl"
	//timeLayout := "2006-01-02 15:04:05"

	address := strings.TrimPrefix(wc.GetString(":address"), "0x")
	result, err = chain.GetTxByAccount(beego.AppConfig.String("chain_api"), address, "payout")
	if err != nil {
		goto errDeal
	}
	for _, v := range result.TxDs {
		var x chain.TxD_Ex
		x.TxID = v.TxID
		x.TxHash = v.TxHash
		x.BlockHeight = v.BlockHeight
		x.BlockHash = v.BlockHash
		x.ActionCount = v.ActionCount
		x.ActionID = v.ActionID
		x.Src = v.Src
		x.Dst = v.Dst
		x.Nonce = v.Nonce
		x.Amount = v.Amount
		x.ResultCode = v.ResultCode
		x.ResultMsg = v.ResultMsg
		x.CreateAt = v.CreateAt
		x.JData = v.JData
		x.Memo = v.Memo
		x.TimeStr = time.Unix(int64(x.CreateAt)/1e9, 0).Format(time.RFC3339)
		datas = append(datas, x)
	}

	wc.Data["TransDetail"] = datas
	return
errDeal:
	log.Logger.Error("TxOutByAddress", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

func (wc *WebController) Block() {

	var (
		txs []datamanage.Trans
		err error
	)
	hash := strings.TrimPrefix(wc.GetString(":hash"), "0x")

	blk, err := datamanage.GetBlockByHash(hash)
	if err != nil {
		goto errDeal
	}

	wc.Layout = "layout.html"
	wc.TplName = "block.tpl"
	wc.Data["Block"] = blk

	if txs, err = datamanage.GetTxData(hash); err != nil {
		goto errDeal
	}
	wc.Data["Transactions"] = txs

	return
errDeal:
	log.Logger.Error("Block", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

func (wc *WebController) TxsPage() {
	var (
		err    error
		result *chain.Result
		datas  []chain.TxD_Ex
	)
	wc.Layout = "layout.html"
	wc.TplName = "search_txs.tpl"
	//timeLayout := "2006-01-02 15:04:05"

	if result, err = chain.GetTx(beego.AppConfig.String("chain_api"), "desc", 0, TxLimit); err != nil {
		goto errDeal
	}
	for _, v := range result.TxDs {
		var x chain.TxD_Ex
		x.TxID = v.TxID
		x.TxHash = v.TxHash
		x.BlockHeight = v.BlockHeight
		x.BlockHash = v.BlockHash
		x.ActionCount = v.ActionCount
		x.ActionID = v.ActionID
		x.Src = v.Src
		x.Dst = v.Dst
		x.Nonce = v.Nonce
		x.Amount = v.Amount
		x.ResultCode = v.ResultCode
		x.ResultMsg = v.ResultMsg
		x.CreateAt = v.CreateAt
		x.JData = v.JData
		x.Memo = v.Memo
		x.TimeStr = time.Unix(int64(x.CreateAt)/1e9, 0).Format(time.RFC3339)
		datas = append(datas, x)
	}

	wc.Data["TransDetail"] = datas
	return
errDeal:
	log.Logger.Error("TxsPage", zap.Error(err))
	wc.TplName = "error.tpl"
	return
}

func (wc *WebController) ContractPage() {
	wc.Layout = "layout.html"
	wc.TplName = "contract.tpl"
}

func (wc *WebController) Search() {
	hash := strings.TrimPrefix(wc.GetString(":hash"), "0x")

	reg := regexp.MustCompile(`^[1-9][0-9]*$`)

	if reg.MatchString(hash) {

		height, _ := strconv.ParseUint(hash, 10, 64)

		block, err := datamanage.GetBlockByHeight(height)
		if err == nil {
			hash = block.Hash
			wc.Redirect("/view/blocks/hash/"+hash, 302)
		}
	} else {
		switch len(hash) {
		case 40, 42: //account
			wc.Redirect("/view/accounts/"+hash, 302)
			break
		case 66, 64: // tx or block hash
			blk, _ := datamanage.GetBlockByHash(hash)
			if blk.Height != 0 {
				wc.Redirect("/view/blocks/hash/"+hash, 302)
			} else {
				wc.Redirect("/view/trans/detail/"+hash, 302)
			}
			break
		default:
			log.Logger.Error("Search", zap.String("hash", hash))
		}
	}
}
