package datamanager

import (
	"database/sql"
	"github.com/axengine/btchain/database"
	"github.com/axengine/btchain/define"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

func (m *DataManager) PrepareTransaction() (*sql.Stmt, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash"},
		database.Feild{Name: "blockheight"},
		database.Feild{Name: "blockhash"},
		database.Feild{Name: "actioncount"},
		database.Feild{Name: "actionid"},
		database.Feild{Name: "uid"},
		database.Feild{Name: "relateduid"},
		database.Feild{Name: "direction"},
		database.Feild{Name: "nonce"},
		database.Feild{Name: "amount"},
		database.Feild{Name: "resultcode"},
		database.Feild{Name: "resultmsg"},
		database.Feild{Name: "createdat"},
		database.Feild{Name: "jdata"},
		database.Feild{Name: "memo"},
	}

	return m.qdb.Prepare(database.TableTransactions, fields)
}

func (m *DataManager) AddTransactionStmt(stmt *sql.Stmt, data *define.TransactionData) (err error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash", Value: data.TxHash.Hex()},
		database.Feild{Name: "blockheight", Value: data.BlockHeight},
		database.Feild{Name: "blockhash", Value: data.BlockHash.Hex()},
		database.Feild{Name: "actioncount", Value: data.ActionCount},
		database.Feild{Name: "actionid", Value: data.ActionID},
		database.Feild{Name: "uid", Value: data.UID},
		database.Feild{Name: "relateduid", Value: data.RelatedUID},
		database.Feild{Name: "direction", Value: data.Direction},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "amount", Value: data.Amount},
		database.Feild{Name: "resultcode", Value: data.ResultCode},
		database.Feild{Name: "resultmsg", Value: data.ResultMsg},
		database.Feild{Name: "createdat", Value: data.CreateAt},
		database.Feild{Name: "jdata", Value: data.JData},
		database.Feild{Name: "memo", Value: data.Memo},
	}
	_, err = m.qdb.Excute(stmt, fields)
	return err
}

// AddTransaction insert a tx record
func (m *DataManager) AddTransaction(data *define.TransactionData) (uint64, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash", Value: data.TxHash.Hex()},
		database.Feild{Name: "blockheight", Value: data.BlockHeight},
		database.Feild{Name: "blockhash", Value: data.BlockHash.Hex()},
		database.Feild{Name: "actioncount", Value: data.ActionCount},
		database.Feild{Name: "actionid", Value: data.ActionID},
		database.Feild{Name: "uid", Value: data.UID},
		database.Feild{Name: "relateduid", Value: data.RelatedUID},
		database.Feild{Name: "direction", Value: data.Direction},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "amount", Value: data.Amount.String()},
		database.Feild{Name: "resultcode", Value: data.ResultCode},
		database.Feild{Name: "resultmsg", Value: data.ResultMsg},
		database.Feild{Name: "createdat", Value: data.CreateAt},
		database.Feild{Name: "jdata", Value: data.JData},
		database.Feild{Name: "memo", Value: data.Memo},
	}

	sqlRes, err := m.qdb.Insert(database.TableTransactions, fields)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// QuerySingleTx query single tx record
func (m *DataManager) QuerySingleTx(txhash *ethcmn.Hash) (*define.TransactionData, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "txhash", Value: txhash.Hex()},
	}

	var result []database.TxData
	err := m.qdb.SelectRows(database.TableTransactions, where, nil, nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}
	if len(result) > 1 {
		// panic ?
	}

	/*
		database.Feild{Name: "txhash", Value: data.TxHash.Hex()},
		database.Feild{Name: "blockheight", Value: data.BlockHeight},
		database.Feild{Name: "blockhash", Value: data.BlockHash.Hex()},
		database.Feild{Name: "actioncount", Value: data.ActionCount},
		database.Feild{Name: "actionid", Value: data.ActionID},
		database.Feild{Name: "uid", Value: data.UID},
		database.Feild{Name: "relateduid", Value: data.RelatedUID},
		database.Feild{Name: "direction", Value: data.Direction},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "amount", Value: data.Amount},
		database.Feild{Name: "resultcode", Value: data.ResultCode},
		database.Feild{Name: "resultmsg", Value: data.ResultMsg},
		database.Feild{Name: "createdat", Value: data.CreateAt},
		database.Feild{Name: "jdata", Value: data.JData},
		database.Feild{Name: "memo", Value: data.Memo},
	*/

	r := result[0]
	td := define.TransactionData{
		TxID:        r.TxID,
		TxHash:      ethcmn.HexToHash(r.TxHash),
		BlockHeight: r.BlockHeight,
		BlockHash:   ethcmn.HexToHash(r.BlockHash),
		ActionCount: r.ActionCount,
		ActionID:    r.ActionID,
		UID:         ethcmn.HexToAddress(r.UID),
		RelatedUID:  ethcmn.HexToAddress(r.RelatedUID),
		Direction:   r.Direction,
		Nonce:       r.Nonce,
		Amount:      Str2Big(r.Amount),
		ResultCode:  r.ResultCode,
		ResultMsg:   r.ResultMsg,
		CreateAt:    r.CreateAt,
		JData:       r.JData,
		Memo:        r.Memo,
	}
	return &td, nil
}

func Str2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

// QueryAccountTxs query account's tx records
func (m *DataManager) QueryAccountTxs(accid *ethcmn.Address, cursor, limit uint64, order string) ([]define.TransactionData, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if accid != nil {
		where = append(where, database.Where{Name: "account", Value: accid.Hex()})
	}
	orderT, err := database.MakeOrder(order, "txid")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("txid", cursor, limit)

	var result []database.TxData
	err = m.qdb.SelectRows(database.TableTransactions, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	var res []define.TransactionData
	for _, r := range result {
		td := define.TransactionData{
			TxID:        r.TxID,
			TxHash:      ethcmn.HexToHash(r.TxHash),
			BlockHeight: r.BlockHeight,
			BlockHash:   ethcmn.HexToHash(r.BlockHash),
			ActionCount: r.ActionCount,
			ActionID:    r.ActionID,
			UID:         ethcmn.HexToAddress(r.UID),
			RelatedUID:  ethcmn.HexToAddress(r.RelatedUID),
			Direction:   r.Direction,
			Nonce:       r.Nonce,
			Amount:      Str2Big(r.Amount),
			ResultCode:  r.ResultCode,
			ResultMsg:   r.ResultMsg,
			CreateAt:    r.CreateAt,
			JData:       r.JData,
			Memo:        r.Memo,
		}

		res = append(res, td)
	}

	return res, nil
}

// QueryAllTxs query all tx records
func (m *DataManager) QueryAllTxs(cursor, limit uint64, order string) ([]define.TransactionData, error) {
	return m.QueryAccountTxs(nil, cursor, limit, order)
}
