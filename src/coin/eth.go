package coin

import (
	"bytes"
	"config"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ybbus/jsonrpc"
	"math/big"
	"strconv"
	"strings"
)

type ETHAgent struct {
	ServerUrl string
}

func (agent *ETHAgent) Init(urlstr string) {
	agent.ServerUrl = urlstr
}
func (agent *ETHAgent) DoHttpJsonRpcCallType1(method string, args ...interface{}) (*jsonrpc.RPCResponse, error) {
	rpcClient := jsonrpc.NewClient(agent.ServerUrl)
	rpcResponse, err := rpcClient.Call(method, args)
	if err != nil {
		return nil, err
	}
	return rpcResponse, nil
}

func (agent *ETHAgent) CreateTransaction(coinSymbol, from, to, gas, gasprice, value, data string, keyindex uint16) (string, error) {
	nonce, err := agent.GetTransactionCount(from, "pending")
	if err != nil {
		return "", err
	}
	to_address := common.HexToAddress(to)
	AccountNonce := uint64(nonce)
	Price := big.NewInt(0)
	Price.SetString(gasprice, 10)
	GasLimit := big.NewInt(0)
	GasLimit.SetString(gas, 10)
	Amount := big.NewInt(0)
	amountStr := ConvertStringToBigNumber(coinSymbol, value)
	Amount.SetString(amountStr, 10)
	Payload, err := hex.DecodeString(data)
	//types.NewEIP155Signer() if change to formal chain  must use EIP155 and set chainid
	var signer = types.HomesteadSigner{}
	my_trx := types.NewTransaction(AccountNonce, to_address, Amount, GasLimit.Uint64(), Price, Payload)
	hash_data := signer.Hash(my_trx)
	success := 0
	for {
		res, err := CoinSignTrx('1', hash_data.Bytes(), keyindex)

		for i := 0; i < 1; i++ {
			sigdata := make([]byte, 0, 65)
			sigdata = append(sigdata, res...)
			sigdata = append(sigdata, byte(i))
			my_trx, err = my_trx.WithSignature(signer, sigdata)
			if err != nil {
				fmt.Println("withsignature", err.Error())
			}
			sig_from, err := signer.Sender(my_trx)
			if err != nil {
				fmt.Println(err.Error())
			}
			if sig_from.Hex() == from {
				success = 1
				break
			}
		}
		if success > 0 {
			break
		}

	}

	raw_writer := new(bytes.Buffer)
	my_trx.EncodeRLP(raw_writer)

	return "0x" + hex.EncodeToString(raw_writer.Bytes()), nil
}

//cal call contract cost
func (agent *ETHAgent) EstimateGas(to string, value string, data string) (*big.Int, error) {
	if to == "" {
		to = "0x0"
	}
	res, err := agent.DoHttpJsonRpcCallType1("eth_estimateGas", map[string]string{"to": to, "value": value, "data": data})
	if err != nil {
		return big.NewInt(0), err
	}
	esti_fee := big.NewInt(0)

	esti_fee_str, err := res.GetString()
	if err != nil {
		return big.NewInt(0), err
	}
	esti_fee.SetString(esti_fee_str[2:], 16)

	return esti_fee, nil
}

func (agent *ETHAgent) GetBalanceByAddress(addr string) (*big.Int, error) {

	res, err := agent.DoHttpJsonRpcCallType1("eth_getBalance", addr, "latest")
	if err != nil {
		return big.NewInt(0), err
	}
	balance := big.NewInt(0)

	balance_str, err := res.GetString()
	if err != nil {
		return big.NewInt(0), err
	}
	balance.SetString(balance_str[2:], 16)

	return balance, nil
}

func ConvertBigNumberToString(coinSymbol string, input *big.Int) string {
	precision := big.NewInt(10)
	for i := 0; i < config.GlobalSupportCoinMgr[coinSymbol].Precision-1; i++ {
		precision.Mul(precision, big.NewInt(10))
	}
	modRes := big.NewInt(0)
	divValue, modRes := input.DivMod(input, precision, modRes)

	res := fmt.Sprintf("%d.%018d", divValue, modRes)
	res = strings.TrimRight(res, "0")
	res = strings.TrimRight(res, ".")
	return res
}

func AddZero(origin string, count int) string {
	for i := 0; i < count; i++ {
		origin += "0"
	}
	return origin
}

func ConvertStringToBigNumber(coinSymbol string, input string) string {
	precision := big.NewInt(10)
	for i := 0; i < config.GlobalSupportCoinMgr[coinSymbol].Precision-1; i++ {
		precision.Mul(precision, big.NewInt(10))
	}
	position := strings.Index(input, ".")
	if position == -1 {
		return AddZero(input, config.GlobalSupportCoinMgr[coinSymbol].Precision)
	}
	res := input[:position]
	count := config.GlobalSupportCoinMgr[coinSymbol].Precision - (len(input) - position - 1)
	res += input[position+1:]
	return AddZero(res, count)
}

func (agent *ETHAgent) GetTransactionRealCost(coinSymbol, trxId string) (string, error) {

	res, err := agent.DoHttpJsonRpcCallType1("eth_getTransactionReceipt", trxId)
	if err != nil {
		return "0", err
	}
	gasUsed := big.NewInt(0)
	gasUsed_str, exist := res.Result.(map[string]interface{})["gasUsed"].(string)
	if exist != true {
		return "0", err
	}
	gasUsed.SetString(gasUsed_str[2:], 16)

	res, err = agent.DoHttpJsonRpcCallType1("eth_getTransactionByHash", trxId)
	if err != nil {
		return "0", err
	}
	gasPrice := big.NewInt(0)
	gasPrice_str, exist := res.Result.(map[string]interface{})["gasPrice"].(string)
	if exist != true {
		return "0", err
	}
	gasPrice.SetString(gasPrice_str[2:], 16)
	totalCost := big.NewInt(0)
	totalCost.Mul(gasPrice, gasUsed)
	return ConvertBigNumberToString(coinSymbol, totalCost), nil
}

//tag latest pending
func (agent *ETHAgent) GetTransactionCount(addr string, tag string) (int64, error) {
	res, err := agent.DoHttpJsonRpcCallType1("eth_getTransactionCount", addr, tag)
	if err != nil {
		return 0, err
	}
	trxC, err := res.GetString()
	if err != nil {
		return 0, nil
	}
	txCount, err := strconv.ParseInt(trxC[2:], 16, 64)
	if err != nil {
		return 0, nil
	}

	return txCount, nil
}

func (agent *ETHAgent) BroadcastTransaction(rawtrx string) (string, error) {
	res, err := agent.DoHttpJsonRpcCallType1("eth_sendRawTransaction", rawtrx)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", errors.New(res.Error.Message)
	}
	txid, err := res.GetString()
	if err != nil {
		return "", nil
	}
	return txid, err
}

func (agent *ETHAgent) IsTransactionConfirmed(rawtrx string) (bool, error) {
	res, err := agent.DoHttpJsonRpcCallType1("eth_getTransactionByHash", rawtrx)
	if err != nil {
		return false, err
	}
	if res.Error != nil {
		return false, errors.New(res.Error.Message)
	}
	resmap, ok := res.Result.(map[string]interface{})
	if ok == false {
		return false, errors.New("parse response error")
	}
	cbn, ok := resmap["blockNumber"].(string)
	if err != nil {
		return false, err
	}
	confirm_height, err := strconv.ParseInt(cbn[2:], 16, 32)
	if confirm_height > 0 {
		res, err := agent.DoHttpJsonRpcCallType1("eth_blockNumber")
		if err != nil {
			return false, err
		}
		if res.Error != nil {
			return false, errors.New(res.Error.Message)
		}
		curBlcokNumber, ok := res.Result.(string)
		if ok == false {
			return false, errors.New("parse response error")
		}
		height_number, err := strconv.ParseInt(curBlcokNumber[2:], 16, 32)
		if err != nil {
			return false, err
		}
		if height_number-confirm_height > int64(config.GlobalSupportCoinMgr["ETH"].ConfirmCount) {
			return true, nil
		}
	}
	return false, nil

}

func ETHCalcAddressByPubkey(pubKeyStr string) (string, error) {
	pubKeyBytes, err := hex.DecodeString("04" + pubKeyStr)
	if err != nil {
		return "", err
	}

	pub_key, err := crypto.UnmarshalPubkey(pubKeyBytes)
	//fmt.Println(pub_key)
	//fmt.Println(hex.EncodeToString(crypto.FromECDSAPub(pub_key)))
	addrStr := crypto.PubkeyToAddress(*pub_key)

	return addrStr.Hex(), nil
}

func ETHValidateAddress(address string) bool {
	if address == "0x0" {
		return true
	}
	if strings.HasPrefix(address, "0x") {
		if len(address) != 42 {
			return false
		}
		_, err := hex.DecodeString(address[2:])
		if err != nil {
			return false
		}
	} else {
		if len(address) != 40 {
			return false
		}
		_, err := hex.DecodeString(address[:])
		if err != nil {
			return false
		}
	}
	return true
}