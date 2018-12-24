package datamanage

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"errors"

	"github.com/astaxie/beego"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sqlx.DB
	lockdb sync.RWMutex
)

const (
	createBlockCollectSQL = `CREATE TABLE IF NOT EXISTS block_t
    (
		height				INTEGER NOT NULL,
		hash				VARCHAR(64) PRIMARY KEY,
		chain_id			VARCHAR(32) NOT NULL,
		time				DATETIME NOT NULL,
		num_txs				INTEGER NOT NULL,
		last_commit_hash    VARCHAR(64) NOT NULL,
		data_hash			VARCHAR(64) NOT NULL,
		validators_hash		VARCHAR(64) NOT NULL,
		app_hash			VARCHAR(64) NOT NULL
	);`

	createTxCollectSQL = `CREATE TABLE IF NOT EXISTS transaction_t
    (
		txid 			INTEGER PRIMARY KEY	AUTOINCREMENT,
		hash			VARCHAR(64),
		blkhash 		VARCHAR(64),
		basefee 		INTEGER,
		txtype 			INTEGER,
		actions			INTEGER,
		createtime 		INTEGER
    );`

	createOptionsSQL = `CREATE TABLE IF NOT EXISTS action_t
    (
		id				INTEGER PRIMARY KEY	AUTOINCREMENT,
		actid 			INTEGER NOT NULL,
		txid 			INTEGER NOT NULL,
		txhash			VARCHAR(64),
		src				VARCHAR(42),
		dst				VARCHAR(42),
		amount			VARCHAR(100),
		body 			VARCHAR(256),
		memo 			TEXT DEFAULT '',
		createtime		INTEGER
    );`
)

const (
	BLOCK_TABLE       = "block_t"
	TRANSACTION_TABLE = "transaction_t"
	ACTION_TABLE      = "action_t"
)

var (
	createIndexs = []string{
		"CREATE INDEX IF NOT EXISTS transaction_blk_hash ON transaction_t (blkhash);",
		"CREATE INDEX IF NOT EXISTS action_txid ON action_t (txid);",
	}
)

type Field struct {
	Name  string
	Value interface{}
}
type Where struct {
	Name  string
	Value interface{}
	Op    string // can be =、>、<、<> and any operator supported by sql-database
}

// GetOp get operator of current where clause, default =
func (w *Where) GetOp() string {
	if w.Op == "" {
		return "="
	}
	return w.Op
}

// Order  used to identify query order
type Order struct {
	Type   string   // "asc" or "desc"
	Feilds []string // order by x
}

func CreateDB() (err error) {
	lockdb.Lock()
	defer lockdb.Unlock()

	os.Remove(fmt.Sprintf("./%s.db", beego.AppConfig.String("appname")))

	if db, err = sqlx.Open("sqlite3", fmt.Sprintf("./%s.db", beego.AppConfig.String("appname"))); err != nil {
		goto errDeal
	}

	if _, err = db.Exec(createBlockCollectSQL); err != nil {
		goto errDeal
	}
	if _, err = db.Exec(createTxCollectSQL); err != nil {
		goto errDeal
	}
	if _, err = db.Exec(createOptionsSQL); err != nil {
		goto errDeal
	}

	for _, sqlIndex := range createIndexs {
		if _, err = db.Exec(sqlIndex); err != nil {
			goto errDeal
		}
	}
errDeal:
	return
}

func InsertData(tableName string, fields []Field) (sql.Result, error) {

	lockdb.Lock()
	defer lockdb.Unlock()

	var sqlBuff bytes.Buffer

	sqlBuff.WriteString(fmt.Sprintf("insert into %v (", tableName))

	for i := 0; i < len(fields)-1; i++ {
		sqlBuff.WriteString(fmt.Sprintf("%s,", fields[i].Name))
	}
	sqlBuff.WriteString(fmt.Sprintf("%s) values (", fields[len(fields)-1].Name))

	// fill field value
	for i := 0; i < len(fields)-1; i++ {
		sqlBuff.WriteString("?,")
	}
	sqlBuff.WriteString("?);")

	// execute
	values := make([]interface{}, len(fields))
	for i, v := range fields {
		values[i] = v.Value
	}
	return db.Exec(sqlBuff.String(), values...)
}

func SelectRowsOffset(table string, where []Where, order *Order, offset, limit uint64, result interface{}) error {
	lockdb.Lock()
	defer lockdb.Unlock()
	if table == "" {
		return errors.New("table name is required")
	}
	//	if order != nil && (len(order.Feilds) == 0 || order.Type == "") {
	//		return errors.New("order type and fields is required")
	//	}

	values := make([]interface{}, len(where))
	for i, v := range where {
		values[i] = v.Value
	}

	var sqlBuff bytes.Buffer
	sqlBuff.WriteString(fmt.Sprintf("select * from %s where 1 = 1", table))
	for i := 0; i < len(where); i++ {
		sqlBuff.WriteString(fmt.Sprintf(" and %s %s ? ", where[i].Name, where[i].GetOp()))
	}
	if order != nil {
		// append order by clause for ordering
		sqlBuff.WriteString(fmt.Sprintf(" order by %s ", order.Feilds[0]))
		for i := 1; i < len(order.Feilds); i++ {
			sqlBuff.WriteString(fmt.Sprintf(" , %s ", order.Feilds[i]))
		}
		sqlBuff.WriteString(order.Type)

		// append limit clause for paging
		sqlBuff.WriteString(fmt.Sprintf(" limit %d offset %d ", limit, offset))
	}

	return db.Select(result, sqlBuff.String(), values...)
}
