package define

import (
	"crypto/ecdsa"
	"errors"
	"github.com/axengine/go-amino"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type TxSignature struct {
	Sig []byte
}

// todo 待测试
func (tx *Transaction) SigHash() (h ethcmn.Hash) {
	var signTx Transaction
	for _, v := range tx.Actions {
		var action Action
		action.ID = v.ID
		action.CreatedAt = v.CreatedAt
		action.Src = v.Src
		action.Dst = v.Dst
		action.Amount = v.Amount
		action.Data = v.Data
		signTx.Actions = append(signTx.Actions, &action)
	}
	hw := sha3.NewKeccak256()
	b, _ := amino.MarshalBinaryBare(signTx)
	hw.Write(b)
	hw.Sum(h[:0])
	return
}

func (tx *Transaction) sign(privkeys []*ecdsa.PrivateKey) ([]TxSignature, error) {
	hash := tx.SigHash()
	signatures := make([]TxSignature, len(privkeys))

	for i, v := range privkeys {
		sig, err := crypto.Sign(hash.Bytes(), v)
		if err != nil {
			return signatures, err
		}
		//fmt.Println("sign hash:", tx.SigHash().Hex(), " priv:", v, " sign:", sig)
		signatures[i] = TxSignature{Sig: sig}
	}

	return signatures, nil
}

func (tx *Transaction) Sign(privkeys []*ecdsa.PrivateKey) error {
	signatures, err := tx.sign(privkeys)
	for i, v := range tx.Actions {
		copy(v.SignHex[:], signatures[i].Sig)
	}

	return err
}

func Signer(tx *Transaction, sig []byte) (ethcmn.Address, error) {
	if len(sig) != 65 {
		return ethcmn.Address{}, errors.New("invalid signature length")
	}

	sigHash := tx.SigHash()
	//fmt.Println("signer hash:", tx.SigHash().Hex(), " sign:", sig)
	publicKey, err := crypto.Ecrecover(sigHash.Bytes(), sig)
	if err != nil {
		return ethcmn.Address{}, err
	}
	if len(publicKey) == 0 || publicKey[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}

	return ethcmn.BytesToAddress(publicKey), nil
}

func (tx *Transaction) CheckSig() error {
	for _, v := range tx.Actions {
		if _, err := Signer(tx, v.SignHex[:]); err != nil {
			return err
		}
	}
	return nil
}
