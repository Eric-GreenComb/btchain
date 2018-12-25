package datamanage

import (
	"database/sql"
	"github.com/axengine/btchain/log"
	"go.uber.org/zap"
	"time"

	"github.com/axengine/btchain/browser/block"
)

const ONCENUMS = 10

var (
	Block_Attr = []string{"height", "hash", "chain_id", "time", "num_txs",
		"last_commit_hash", "data_hash", "validators_hash", "app_hash"}
	Transaction_Attr = []string{"hash", "blkhash", "basefee", "txtype", "actions", "createtime"}
	Action_Attr      = []string{"actid", "txid", "txhash", "src", "dst", "amount", "body", "memo", "createtime"}
)

type Block struct {
	Height        uint64 `db:"height"`
	Hash          string `db:"hash"`
	ChainID       string `db:"chain_id"`
	Time          string `db:"time"`
	NumTx         int    `db:"num_txs"`
	LastComitHash string `db:"last_commit_hash"`
	DataHash      string `db:"data_hash"`
	ValidHash     string `db:"validators_hash"`
	AppHash       string `db:"app_hash"`
}

type Trans struct {
	TxID       uint64 `db:"txid"`
	Hash       string `db:"hash"`
	BlockHash  string `db:"blkhash"`
	BaseFee    uint64 `db:"basefee"`
	TxType     uint8  `db:"txtype"`
	Actions    uint8  `db:"actions"`
	CreateTime uint64 `db:"createtime"`
}

type Action struct {
	ID         uint64 `db:"id"`
	ActID      string `db:"actid"`
	TxID       uint64 `db:"txid"`
	TxHash     string `db:"txhash"`
	Src        string `db:"src"`
	Dst        string `db:"dst"`
	Amount     string `db:"amount"`
	Body       string `db:"body"`
	Memo       string `db:"memo"`
	CreateTime string `db:"createtime"`
}

func GeneratField(names []string, args ...interface{}) (fields []Field) {
	for i, n := range names {
		f := Field{
			Name:  n,
			Value: args[i],
		}
		fields = append(fields, f)
	}
	return
}

func InitBlock() {
	if err := CreateDB(); err != nil {
		log.Logger.Error("CreateDB", zap.Error(err))
		return
	}
	for {
		if err := GetBlockData(); err != nil {
			log.Logger.Error("GetBlockData", zap.Error(err))
		}
		time.Sleep(time.Second * 5)
	}
	return
}

//获取指定num数量的block数据从数据库中
func GetBlock(num uint64) (blks []Block, err error) {
	order := &Order{
		Type:   "desc",
		Feilds: []string{"height"},
	}
	if err = SelectRowsOffset(BLOCK_TABLE, nil, order, 0, num, &blks); err != nil {
		return
	}
	return
}

//根据height获取block
func GetBlockByHeight(height uint64) (blk Block, err error) {
	where := []Where{
		Where{
			Name:  "height",
			Op:    "=",
			Value: height,
		},
	}
	var blks []Block

	if err = SelectRowsOffset(BLOCK_TABLE, where, nil, 0, 0, &blks); err != nil {
		return
	}
	if len(blks) == 0 {
		return
	}
	return blks[0], nil
}

//根据hash获取block
func GetBlockByHash(hash string) (blk Block, err error) {
	where := []Where{
		Where{
			Name:  "hash",
			Op:    "=",
			Value: hash,
		},
	}
	var blks []Block

	if err = SelectRowsOffset(BLOCK_TABLE, where, nil, 0, 0, &blks); err != nil {
		return
	}
	if len(blks) == 0 {
		return
	}
	return blks[0], nil
}

//从数据库中获取指定hash值得tx交易信息
func GetTxData(hash string) (trans []Trans, err error) {

	where := []Where{
		Where{
			Name:  "blkhash",
			Op:    "=",
			Value: hash,
		},
	}

	if err = SelectRowsOffset(TRANSACTION_TABLE, where, nil, 0, 0, &trans); err != nil {
		return
	}
	return
}

//根据txid从数据库获取operations数据
func Getactions(txid uint64) (opts []Action, err error) {
	where := []Where{
		Where{
			Name:  "txid",
			Op:    "=",
			Value: txid,
		},
	}
	if err = SelectRowsOffset(ACTION_TABLE, where, nil, 0, 0, &opts); err != nil {
		return
	}
	return

}

//获取数据库中height的最大值
func GetLastHeight() (height uint64, err error) {

	order := &Order{
		Type:   "desc",
		Feilds: []string{"height"},
	}

	var blocks []Block

	if err = SelectRowsOffset(BLOCK_TABLE, nil, order, 0, 1, &blocks); err != nil {
		return
	}

	if len(blocks) == 0 {
		height = 0
	} else {
		height = blocks[0].Height

	}
	return

}

//从chain获取从数据库height最大值起ONCENUMS的块数据，并分析tx，operation，数据插入数据库
func GetBlockData() (err error) {

	var (
		result sql.Result
		height uint64
		blks   []*block.BlockData
	)

	if height, err = GetLastHeight(); err != nil {
		log.Logger.Error("GetLastHeight", zap.Error(err))
		return
	}
	log.Logger.Info("GetBlockData", zap.Uint64("height", height))
	if height > 0 {
		if blks, err = block.GetBlockChainList(height+1, height+ONCENUMS); err != nil {
			return
		}
	} else {
		if blks, err = block.GetLastNumChain(ONCENUMS); err != nil {
			return
		}
	}

	for _, blk := range blks {
		if _, err = InsertData(BLOCK_TABLE, GeneratField(Block_Attr, blk.Height, blk.Hash, blk.ChainID, blk.Time,
			blk.NumTxs, blk.LastCommitHash, blk.DataHash, blk.ValidatorsHash, blk.AppHash)); err != nil {
			log.Logger.Error("InsertData BLOCK_TABLE", zap.Error(err))
			continue
		}
		for _, tx := range blk.Txs {
			if result, err = InsertData(TRANSACTION_TABLE, GeneratField(Transaction_Attr, tx.SigHash().Hex(), blk.Hash, 0,
				tx.Type, len(tx.Actions), blk.Time.Unix())); err != nil {
				log.Logger.Error("InsertData TRANSACTION_TABLE", zap.Error(err))
				continue
			}

			txid, _ := result.LastInsertId()

			for _, action := range tx.Actions {
				if _, err = InsertData(ACTION_TABLE, GeneratField(Action_Attr, action.ID,
					txid, tx.SigHash().Hex(), action.Src.Hex(), action.Dst.Hex(), action.Amount.String(),
					action.Data, action.Memo, action.CreatedAt)); err != nil {
					log.Logger.Error("InsertData ACTION_TABLE", zap.Error(err))
					continue
				}
			}
		}
	}
	return
}

func GetTypeString(itype int) string {
	switch itype {
	case 1:
		return "payment"
	default:
		return "payment"
	}
}
