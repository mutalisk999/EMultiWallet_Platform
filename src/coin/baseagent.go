package coin

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/ybbus/jsonrpc"
	"math"
	"strconv"
	"strings"
)

type UTXODetail struct {
	TxId          string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	Account       string `json:"account"`
	ScriptPubKey  string `json:"scriptPubKey"`
	RedeemScript  string `json:"redeemScript"`
	Amount        int64  `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Spendable     bool   `json:"spendable"`
	Solvable      bool   `json:"solvable"`
}
type UTXOsDetail []UTXODetail

type InPut struct {
	TxId         string `json:"txid"`
	Vout         int    `json:"vout"`
	ScriptPubKey string `json:"scriptPubKey"`
}

type InPuts []InPut
type OutPuts map[string]string

func (o UTXOsDetail) Len() int {
	return len(o)
}

func (o UTXOsDetail) Less(i int, j int) bool {
	if o[i].Confirmations >= o[j].Confirmations {
		return true
	}
	return false
}

func (o UTXOsDetail) Swap(i int, j int) {
	o[i], o[j] = o[j], o[i]
}

type BaseAgent interface {
	Type() string
	CoinType() string
	Init(string)
	DoHttpJsonRpcCallType1(string, ...interface{}) (*jsonrpc.RPCResponse, error)
	GetBalanceByAddressRPC(string) (string, error)
	GetUtxosByAddressRPC(string) ([]UTXODetail, error)
	ImportAddressRPC(string) error
	BroadcastTransactionRPC(string) (string, error)
	GetRawTransactionRPC(string)(string,error)
	IsTransactionConfirmedRPC(string) (bool, error)
	IsAddressValidRPC(string) (bool, error)
	BuildTrxInPutsOutPutsRPC(string, string, string, string) (string, InPuts, OutPuts, error)
	BuildTrxFromUtxosRPC([]UTXODetail, string, string, string, string) (string, InPuts, OutPuts, error)
	CreateRawTransactionRPC(InPuts, OutPuts) (string, error)
	SignRawTransactionRPC(string, string, uint16, []UTXODetail) (string, error)
	MultiSignRawTransactionRPC(string, string, uint16, []UTXODetail) (string, error)
	MultiVerifySignRawTransactionRPC(string, string, []UTXODetail) error
	CombineRawTransactionRPC([]string) (string, error)
}

func AgentFactory(coinSymbol string) BaseAgent {
	switch coinSymbol {
	case "BTC":
		return new(BTCAgent)
	case "LTC":
		return new(LTCAgent)
	case "BCH":
		return new(BCHAgent)
	case "UB":
		return new(UBAgent)
	default:
		return nil
	}
	return nil
}

func ToPrecisionAmount(value string, nPrecision int) (int64, error) {
	precision := int64(math.Pow10(nPrecision))
	strArray := strings.Split(value, ".")
	if len(strArray) == 1 {
		quotient, err := strconv.Atoi(strArray[0])
		if err != nil {
			return 0, errors.New("invalid value: invalid quotient part")
		}
		return int64(quotient) * precision, nil
	} else if len(strArray) == 2 {
		quotient, err := strconv.Atoi(strArray[0])
		if err != nil {
			return 0, errors.New("invalid value: invalid quotient part")
		}

		remainderStr := strArray[1]
		for i := len(remainderStr); i < nPrecision; i++ {
			remainderStr = remainderStr + "0"
		}
		remainderStr = remainderStr[0:nPrecision]

		remainder, err := strconv.Atoi(remainderStr)
		if err != nil {
			return 0, errors.New("invalid value: invalid remainder part")
		}
		return int64(quotient)*precision + int64(remainder), nil
	} else {
		return 0, errors.New("invalid value: too many point")
	}
}

func FromPrecisionAmount(amount int64, nPrecision int) string {
	precision := int64(math.Pow10(nPrecision))
	quotient := amount / precision
	remainder := amount % precision
	return fmt.Sprintf("%d.%08d", quotient, remainder)
}
