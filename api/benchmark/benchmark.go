package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"github.com/axengine/btchain/api/bean"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	clients *ClientMgr
	adds    *AddressMgr
	seq     uint64
	success uint64
	fail    uint64
)

func main() {
	tps := flag.Int("t", 100, "transactions per second")
	actions := flag.Int("actions", 10, "actions per tx")
	du := flag.Int("d", 300, "durations seconds")
	clientNum := flag.Int("c", 10, "http client num")
	flag.Parse()
	log.Println("tps=", *tps, " action=", *actions, " clientNum=", *clientNum, " du=", *du)

	adds = NewAddressMgr(10000)
	clients = NewClientMgr(*clientNum)

	var (
		wg sync.WaitGroup
	)

	tk := time.NewTicker(time.Second)
	exit := time.After(time.Second*time.Duration(*du) + time.Millisecond*500)
OUTLOOP:
	for {
		select {
		case <-tk.C:
			{
				for idx := 0; idx < *tps; idx++ {
					go func(w *sync.WaitGroup) {
						w.Add(1)
						defer w.Done()
						cli := clients.Get()
						to := adds.Get()
						begin := time.Now()
						if err := transaction(cli, "0x061a060880BB4E5AD559350203d60a4349d3Ecd6", to, "1",
							"5b416c67c05f67cdba1de4f1e993040aa7b4f6a6ef022186f3a5640f72e26033", *actions); err != nil {
							log.Println("======>", err)
							atomic.AddUint64(&fail, 1)
						} else {
							atomic.AddUint64(&success, 1)
						}
						atomic.AddUint64(&seq, 1)
						DefaultAnalyzer.Add(&State{idx: atomic.LoadUint64(&seq), t: time.Since(begin)})
					}(&wg)
				}
			}
		case <-exit:
			log.Println("send tx ok...")
			break OUTLOOP
		}
	}
	wg.Wait()
	log.Println("done... err/succ:", fail, "/", success, "\n\t\tanalyze:", DefaultAnalyzer.String())
}

func transaction(cli *http.Client, from, to string, amount string, priv string, actions int) error {
	//fmt.Printf("from:%v to:%v amount:%v priv:%v\n", from, to, amount, priv)
	var tdata bean.Transaction
	tdata.BaseFee = "0"

	var action bean.Action
	action.ID = 0
	action.Src = from
	action.Dst = to
	action.Amount = amount
	action.Priv = priv
	action.Data = "xxxxxxxxxxxxx"
	action.Memo = "xxxxxx"

	for i := 0; i < actions; i++ {
		actionnew := action
		actionnew.ID = i
		tdata.Actions = append(tdata.Actions, &actionnew)
	}

	b, _ := json.Marshal(&tdata)
	resp, err := cli.Post("http://192.168.8.144:8000/v1/"+"transactionsSync", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if is, ok := data["isSuccess"]; !ok || !is.(bool) {
		return errors.New(string(b))
	}
	log.Println("--->", string(b))
	return nil
}
