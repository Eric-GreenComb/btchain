package bean

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var (
	BASE_API_URL = "http://192.168.1.2:10000/v1/"
)

func Test_specialop(t *testing.T) {
	var op SpecilOP
	//op.Pubkey = "46aGd1erHRoDu/aKTqAerWKYEaxxaJiRI/xdFt3Anyc="
	op.Pubkey = "JYnWLYSQkDxfrblXY1gXwxDGBfMwHr4tGORLSjsaOtY="
	op.Power = "0"

	b, _ := json.Marshal(&op)
	resp, err := http.Post(BASE_API_URL+"specialop", "application/json", bytes.NewReader(b[:]))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(b))
}
