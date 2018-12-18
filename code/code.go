// ABCI RPC 内部错误代码 业务无关
package code

const (
	CodeTypeOK            uint32 = 0
	CodeTypeEncodingError uint32 = 10001
	CodeTypeBadNonce      uint32 = 10002
	CodeTypeUnauthorized  uint32 = 10003
	CodeTypeUnknownError  uint32 = 10004
	CodeTypeExec          uint32 = 10005
	CodeALREADY_EXIST     uint32 = 10006
	CodeNotEnoughMoney    uint32 = 10007
	CodeOutOfOrder        uint32 = 10008
	CodeAccountNotFound   uint32 = 10009
	CodeSignerFaild       uint32 = 10010
	CodeUnknownPath       uint32 = 10011
)
