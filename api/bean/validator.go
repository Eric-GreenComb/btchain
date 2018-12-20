package bean

type ValidatorReq struct {
	Pubkey string // pubkey hex
	Power  string // 0-100
	Sign   string // 采用私钥对 公钥 进行签名 hex
}
