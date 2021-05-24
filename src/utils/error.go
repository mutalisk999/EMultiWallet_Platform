package utils

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
)

var GlobalError map[int]string

type Error struct {
	ErrCode int    `json:"code"`
	ErrMsg  string `json:"message"`
}

func SetInternalError(ctx iris.Context, errorStr string) {
	ctx.Values().Set("error", errorStr)
	ctx.StatusCode(iris.StatusInternalServerError)
}

func FormatSysError(err error) *Error {
	sysError := new(Error)
	if err == nil {
		sysError.ErrCode = 0
		sysError.ErrMsg = ""
	} else {
		sysError.ErrCode = 999999
		sysError.ErrMsg = err.Error()
	}
	return sysError
}

func GetErrorString(err *Error) string {
	if err == nil {
		return "no error"
	}
	return fmt.Sprintf("errcode: %d, error msg: %s", err.ErrCode, err.ErrMsg)
}

func InitGlobalError() {
	GlobalError = make(map[int]string)

	GlobalError[100000] = "integer check [%s:%d], setvalue is out of range [%d,%d]"
	GlobalError[100001] = "string length check [%s:%s], length is out of range [%d,%d]"
	GlobalError[100002] = "invalid datetime string format [%s:%s]"
	GlobalError[100010] = "invalid mobile phone number %s"
	GlobalError[100011] = "invalid identity card number %s"
	GlobalError[200000] = "not support json rpc function [%s] for path [%s]"
	GlobalError[200001] = "invalid json rpc params"
	GlobalError[200002] = "should not with session id field while login"
	GlobalError[200003] = "missing session id field"
	GlobalError[200004] = "invalid session id"
	GlobalError[200005] = "forbidden query another user's log"
	GlobalError[200006] = "no permission, not admin"
	GlobalError[200007] = "invalid [%s:%d] state"
	GlobalError[200008] = "forbidden query another user's wallet"
	GlobalError[200009] = "no permission, not accountant"
	GlobalError[200010] = "dst addr %s not in wallet:[%d] transfer whitelist"
	GlobalError[200011] = "forbidden revoke another user's transaction"
	GlobalError[200012] = "accountant required"
	GlobalError[200013] = "can not confirm trx, reach the account confirmations needed"
	GlobalError[200014] = "can not confirm trx, account confirmed before"
	GlobalError[200015] = "forbidden update another user's notification"
	GlobalError[200016] = "fail to create address from pubkey"
	GlobalError[200017] = "can not confirm trx, server confirmed before"
	GlobalError[200018] = "can not confirm trx, reach the server confirmations needed"
	GlobalError[200019] = "no user has relation with wallet:[%d]"
	GlobalError[200100] = "choose unused server pubkey fail"
	GlobalError[200101] = "fail to create multisig address from pubkeys"
	GlobalError[200102] = "invalid serverkeys, walletid [%d]"
	GlobalError[200103] = "get utxo from wallet fail"
	GlobalError[200104] = "invalid serverid, serverid [%d]"
	GlobalError[300001] = "db error, dbname [%s], optype [%s], detail [%s]"
	GlobalError[300010] = "change notification state log fail"
	GlobalError[300011] = "create log of trx-creating fail"
	GlobalError[300012] = "create log of trx-revoking fail"
	GlobalError[300013] = "create log of trx-confirming fail"
	GlobalError[300014] = "Pack transaction fail"
	GlobalError[300015] = "GetMultiSignRedeemScript fail"
	GlobalError[300016] = "MultiSignRawTransaction fail"
	GlobalError[400001] = "invalid auth code"
	GlobalError[400002] = "invalid signature"
	GlobalError[400003] = "invalid pubkey,user didn't exist"
	GlobalError[400004] = "generate session id failed"
	GlobalError[400005] = "loginId or MgmtId is invalid"
	GlobalError[400006] = "unkown db error"
	GlobalError[400007] = "invalid account id,user didn't exist"
	GlobalError[400008] = "forbidden query another user's info"
	GlobalError[400009] = "administrator required"
	GlobalError[400010] = "verify sequence failed"
	GlobalError[400011] = "register param is duplication"
	GlobalError[400012] = "only allow active account login"
	GlobalError[400013] = "invalid pubkey format"
	GlobalError[500000] = "Wallet not found [%d]"
	GlobalError[500001] = "invalid coin id [%d]"
	GlobalError[500002] = "invalid symbol [%s]"
	GlobalError[500003] = "invalid coin address [%s]"
	GlobalError[500004] = "NeedSigCount must bigger than 0"
	GlobalError[500005] = "Wallet [%d] already deleted"
	GlobalError[600001] = "unsupport coin symbol[%s]"
	GlobalError[700000] = "key pool used up"
	GlobalError[800000] = "Coin service response error [%s]"
	GlobalError[900001] = "server is being maintained.please wait."
	GlobalError[900002] = "create anti-tamper signature error"
	GlobalError[900003] = "verify anti-tamper signature error"
	GlobalError[900004] = "anti-tamper signature format error"
	GlobalError[900010] = "json Marshal error"
	GlobalError[900011] = "json Unmarshal error"
	GlobalError[900020] = "web socket request error [%s]"
	GlobalError[900021] = "web socket response error [%s]"
	GlobalError[900040] = "different coin symbols with the same coin id"
}

func InvalidCoinSymbol(coinSymbol string) error {
	return errors.New(fmt.Sprintf("Invaild Coin Symbol [%s]", coinSymbol))
}

func MakeError(errCode int, errArgs ...interface{}) *Error {
	err := new(Error)
	err.ErrCode = errCode
	format, ok := GlobalError[errCode]
	if ok {
		err.ErrMsg = fmt.Sprintf(format, errArgs...)
	}
	return err
}
