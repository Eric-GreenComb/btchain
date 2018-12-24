package chain

import (
	"testing"
)

func TestGetTxByHash(t *testing.T) {
	result, err := GetTxByHash("http://10.253.6.90:8000", "0x550d7a3033bd642b81cd1c2dc51e2f7d3c54cb6b5a2d57c20b2c3eea002db341")

	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result.TxDs[0])
	return
}

func TestGetTxByAccount(t *testing.T) {
	result, err := GetTxByAccount("http://10.253.6.90:8000", "0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result.TxDs)
	return
}

func TestGetTx(t *testing.T) {
	result, err := GetTx("http://10.253.6.90:8000", "desc", 0, 5)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result.TxDs)
	return
}
