package block

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axengine/btchain/browser/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	types "github.com/axengine/btchain/define"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type Block struct {
	*Header `json:"header"`
	*Data   `json:"data"`
}

type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

type Data struct {
	Txs []string `json:"txs"`
	//ExTxs []byte   `json:"extxs"`
	// Volatile
	//hash []byte
}

type Header struct {
	ChainID        string    `json:"chain_id"`
	Height         string    `json:"height"`
	Time           time.Time `json:"time"`
	NumTxs         string    `json:"num_txs"`          // XXX: Can we get rid of this?
	LastCommitHash string    `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       string    `json:"data_hash"`        // transactions
	ValidatorsHash string    `json:"validators_hash"`  // validators for the current block
	AppHash        string    `json:"app_hash"`         // state after txs from the previous block
	//ReceiptsHash   string    `json:"recepits_hash"`    // recepits_hash from previous block
	LastBlockID BlockID `json:"last_block_id"`
}

type BlockMeta struct {
	BlockID BlockID `json:"block_id"` // The block hash
	Header  *Header `json:"header"`   // The block's Header
}

type BlockID struct {
	Hash string `json:"hash"`
}

type Status struct {
	NodeInfo *NodeInfo `json:"node_info"`
	SyncInfo *SyncInfo `json:"sync_info"`
}

type NodeInfo struct {
	NetWork string `json:"network"`
}

type SyncInfo struct {
	LatestBlockHeight string `json:"latest_block_height"`
}

type Metas struct {
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockData struct {
	Hash           string
	ChainID        string
	Height         string
	Time           time.Time
	NumTxs         string
	LastCommitHash string
	DataHash       string
	ValidatorsHash string
	AppHash        string
	Txs            []*types.Transaction
}

func GetLastNumChain(num uint64) (blocks []*BlockData, err error) {
	var (
		stat Status
	)

	if stat, err = GetStatus(); err != nil {
		return
	}
	curHeight, _ := strconv.ParseUint(stat.SyncInfo.LatestBlockHeight, 10, 64)
	log.Logger.Debug("GetLastNumChain", zap.Uint64("curHeight", curHeight))
	return GetBlockChainList(curHeight-num+1, curHeight)
}

func GetBlockChainList(min, max uint64) (blocks []*BlockData, err error) {

	var (
		body        []byte
		metas       Metas
		resultBlock ResultBlock
	)

	url := fmt.Sprintf("%s/blockchain?minHeight=%d&maxHeight=%d", beego.AppConfig.String("chainaddr"), min, max)

	if body, err = GetHTTPResp(url); err != nil {
		log.Logger.Error("GetHTTPResp", zap.Error(err), zap.String("data", string(body)))
		return
	}

	if err = json.Unmarshal(body, &metas); err != nil {
		log.Logger.Error("GetBlockChainList", zap.Error(err), zap.String("data", string(body)))
		return
	}

	for _, o := range metas.BlockMetas {

		blk := &BlockData{
			Hash:           strings.ToLower(o.BlockID.Hash),
			ChainID:        o.Header.ChainID,
			Height:         o.Header.Height,
			Time:           o.Header.Time.Local(),
			NumTxs:         o.Header.NumTxs,
			LastCommitHash: strings.ToLower(o.Header.LastCommitHash),
			DataHash:       strings.ToLower(o.Header.DataHash),
			ValidatorsHash: strings.ToLower(o.Header.ValidatorsHash),
			AppHash:        strings.ToLower(o.Header.AppHash),
		}

		numTxs, _ := strconv.ParseUint(o.Header.NumTxs, 10, 32)
		if numTxs > 0 {
			if resultBlock, err = GetBlock(o.Header.Height); err != nil {
				log.Logger.Error("GetBlock", zap.Error(err))
				return
			}
			for _, v := range resultBlock.Block.Data.Txs {
				txBytes, err1 := base64.StdEncoding.DecodeString(v)
				if err1 != nil {
					log.Logger.Error("base64", zap.Error(err))
					continue
				}

				tx := new(types.Transaction)
				if err = rlp.DecodeBytes(txBytes, tx); err != nil {
					log.Logger.Error("rlp", zap.Error(err))
					beego.Error("Decode Transaction tx bytes: [%v], error : %v\n", common.FromHex(v), err)
					return
				}

				blk.Txs = append(blk.Txs, tx)
			}
		}

		blocks = append(blocks, blk)
	}

	return
}

func GetBlock(height string) (result ResultBlock, err error) {
	url := fmt.Sprintf("%s/block?height=%v", beego.AppConfig.String("chainaddr"), height)
	bytez, errB := GetHTTPResp(url)
	if errB != nil {
		err = errB
		return
	}
	err = json.Unmarshal(bytez, &result)
	if err != nil {
		log.Logger.Error("GetBlock", zap.Error(err), zap.String("data", string(bytez)))
	}
	return
}

func GetStatus() (status Status, err error) {
	var (
		body []byte
	)
	url := fmt.Sprintf("%s/status", beego.AppConfig.String("chainaddr"))

	if body, err = GetHTTPResp(url); err != nil {
		beego.Error("GetHTTPResp failed: %s", err.Error())
		return
	}
	if err = json.Unmarshal(body, &status); err != nil {
		log.Logger.Error("GetStatus", zap.Error(err), zap.String("data", string(body)))
		beego.Error("json.Unmarshal(Status) failed: %s", err.Error())
		return
	}
	return
}

type HTTPResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
}

func GetHTTPResp(url string) (bytez []byte, err error) {

	resp, errR := http.Get(url)
	if errR != nil {
		err = errR
		return
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var hr HTTPResponse
	err = json.Unmarshal(bytez, &hr)
	if err != nil {
		return
	}
	if hr.Result == nil {
		err = errors.New(fmt.Sprintf("json.Unmarshal (%s)HTTPResponse wrong ,maybe you need config 'chain_id'", url))
		return
	}
	bytez, err = hr.Result.MarshalJSON()
	if err != nil {
		return
	}
	return
}
