package btchain

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
	"testing"
)

func Test_commit(t *testing.T) {
	db := ethdb.NewMemDatabase()

	statedb, _ := state.New(ethcmn.Hash{}, state.NewDatabase(db))
	addr := ethcmn.HexToAddress("0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
	statedb.AddBalance(addr, big.NewInt(1000000000000))
	//root := statedb.IntermediateRoot(false)
	appHash, err := statedb.Commit(false)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("root hash:", root.Hex(), " appHash:", appHash.Hex())
	err = statedb.Database().TrieDB().Commit(appHash, true)
	if err != nil {
		t.Fatal(err)
	}
	statedb, err = state.New(appHash, state.NewDatabase(db))
	if err != nil {
		t.Fatal(err)
	}
}
