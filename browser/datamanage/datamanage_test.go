package datamanage

import (
	"github.com/astaxie/beego"
	"testing"
	"time"
)

func init() {
	beego.AppConfig.Set("chainaddr", "http://192.168.8.144:26701")
	beego.AppConfig.Set("appname", "datamanage")
	err := CreateDB()
	if err != nil {
		panic(err)
	}
}

func TestGeneratField(t *testing.T) {

	fs := GeneratField([]string{"name", "age"}, "fhy", 28)

	t.Log(fs)

	return
}

func TestGetBlockData(t *testing.T) {

	for {
		if err := GetBlockData(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}

	return
}

func Test_SyncBlock(t *testing.T) {
	if err := GetBlockData(); err != nil {
		t.Fatal(err)
	}
}

func TestGetLastHeight(t *testing.T) {

	height, err := GetLastHeight()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(height)

	return
}

func TestGetTxData(t *testing.T) {

	blks, err := GetTxData("7426e17fa7b9f6d55771487b350c81bf869e3167")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blks)
	return
}

func TestGetBlock(t *testing.T) {
	blks, err := GetBlock(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blks)
	return

}

func TestGetOperations(t *testing.T) {
	blks, err := Getactions(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blks)
	return
}
