package bean

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var (
	BASE_API_URL = "http://192.168.1.2:10000/v1/"
)

func Test_validatorUpdate(t *testing.T) {

	var (
		pubKeyB64  = "JYnWLYSQkDxfrblXY1gXwxDGBfMwHr4tGORLSjsaOtY="
		privKeyB64 = "g+8E8OkR1JId2XISJAPLkQH9/fNWYGL+IckOLkNXQZElidYthJCQPF+tuVdjWBfDEMYF8zAevi0Y5EtKOxo61g=="
	)

	pubkey, _ := base64.StdEncoding.DecodeString(pubKeyB64)
	privkey, _ := base64.StdEncoding.DecodeString(privKeyB64)
	signBytes := ed25519.Sign(privkey, pubkey)

	signHex := common.Bytes2Hex(signBytes)
	fmt.Println(len(signBytes))
	var op ValidatorReq
	op.Pubkey = pubKeyB64
	op.Power = "10"
	op.Sign = signHex
	fmt.Println(signHex)

	b, _ := json.Marshal(&op)
	resp, err := http.Post(BASE_API_URL+"validatorUpdate", "application/json", bytes.NewReader(b[:]))
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

func Test_sign(t *testing.T) {
	pubkeyB64 := "46aGd1erHRoDu/aKTqAerWKYEaxxaJiRI/xdFt3Anyc="
	privB64 := "6ax5zscgMX8HYrgRzS94FJQlVlsKRFzCYxYH+tzqbQSO89emyT4oYduQsNiEQ8aIh50Ot7xJrVuH9SqIFyzwxw=="

	pubkey, _ := base64.StdEncoding.DecodeString(pubkeyB64)
	privkey, _ := base64.StdEncoding.DecodeString(privB64)

	toSin := []byte("hello world")
	p := ed25519.PrivateKey(privkey)

	sig := ed25519.Sign(p, toSin)
	ss := hex.EncodeToString(sig[:])
	fmt.Println(ss)

	b := ed25519.Verify(pubkey, toSin, sig)
	fmt.Println(b)
}
