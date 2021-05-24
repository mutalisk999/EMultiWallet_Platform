package coin

import (
	"bufio"
	"bytes"
	"config"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/mutalisk999/bitcoin-lib/src/blob"
	"github.com/mutalisk999/bitcoin-lib/src/keyid"
	"github.com/mutalisk999/bitcoin-lib/src/pubkey"
	"github.com/mutalisk999/bitcoin-lib/src/script"
	"github.com/mutalisk999/bitcoin-lib/src/serialize"
	"github.com/mutalisk999/bitcoin-lib/src/transaction"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"github.com/ybbus/jsonrpc"
	"io"
	"sort"
	"strconv"
)

type OMNIAgent struct {
	CoinSymbol string
	ServerUrl  string
}

func (agent *OMNIAgent) Type() string {
	return "OMNIAgent"
}

func (agent *OMNIAgent) SetCoinType(coinSymbol string) {
	agent.CoinSymbol = coinSymbol
}

func (agent *OMNIAgent) CoinType() string {
	return agent.CoinSymbol
}

func (agent *OMNIAgent) Init(urlstr string) {
	agent.ServerUrl = urlstr
}

func (agent *OMNIAgent) DoHttpJsonRpcCallType1(method string, args ...interface{}) (*jsonrpc.RPCResponse, error) {
	rpcClient := jsonrpc.NewClient(agent.ServerUrl)
	rpcResponse, err := rpcClient.Call(method, args)
	if err != nil {
		return nil, err
	}
	return rpcResponse, nil
}

func (agent *OMNIAgent) GetBalanceByAddressRPC(addr string) (string, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return "", errors.New("not support coin")
	}
	nPropertyId := supportCoin.OmniPropertyId

	res, err := agent.DoHttpJsonRpcCallType1("omni_getbalance", addr, nPropertyId)
	if err != nil {
		return "0", err
	}
	if res.Error != nil {
		return "0", errors.New(res.Error.Message)
	}

	resmap, ok := res.Result.(map[string]interface{})
	if ok == false {
		return "0", errors.New("parse response error")
	}
	balanceStr := resmap["balance"].(string)
	if err != nil {
		return "0", err
	}
	return balanceStr, nil
}

func (agent *OMNIAgent) GetFeeBalanceByAddressRPC(addr string) (string, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return "", errors.New("not support coin")
	}
	nPrec := supportCoin.Precision

	res, err := agent.DoHttpJsonRpcCallType1("listunspent", 0, 99999999, []string{addr})
	if err != nil {
		return "0", err
	}
	if res.Error != nil {
		return "0", errors.New(res.Error.Message)
	}

	sum := int64(0)
	for _, i := range res.Result.([]interface{}) {
		out := i.(map[string]interface{})

		amountva, ok := out["amount"]
		if ok == false {
			continue
		}
		amount, err := amountva.(json.Number).Float64()
		if err != nil {
			continue
		}
		amountStr := strconv.FormatFloat(amount, 'f', nPrec, 64)
		amountPrec, err := ToPrecisionAmount(amountStr, nPrec)
		if err != nil {
			continue
		}

		sum += amountPrec
	}
	return FromPrecisionAmount(sum, nPrec), nil
}

func (agent *OMNIAgent) GetUtxosByAddressRPC(addr string) ([]UTXODetail, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return nil, errors.New("not support coin")
	}
	nPrec := supportCoin.Precision

	res, err := agent.DoHttpJsonRpcCallType1("listunspent", 0, 99999999, []string{addr})
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, errors.New(res.Error.Message)
	}

	var utxos UTXOsDetail
	for _, i := range res.Result.([]interface{}) {
		var utxo UTXODetail
		utxo.Address = addr
		out := i.(map[string]interface{})

		amount, ok := out["amount"]
		if ok == false {
			continue
		}
		txid, ok := out["txid"]
		if ok == false {
			continue
		}
		vout, ok := out["vout"]
		if ok == false {
			continue
		}
		scriptPubKey, ok := out["scriptPubKey"]
		if ok == false {
			continue
		}
		confirmations, ok := out["confirmations"]
		if ok == false {
			continue
		}

		amountValue, err := amount.(json.Number).Float64()
		if err != nil {
			continue
		}
		amountStr := strconv.FormatFloat(amountValue, 'f', nPrec, 64)
		amountPrec, err := ToPrecisionAmount(amountStr, nPrec)
		if err != nil {
			continue
		}

		if amountPrec == 0 {
			continue
		}
		utxo.Amount = amountPrec

		txidValue := txid.(string)
		utxo.TxId = txidValue

		i64, err := vout.(json.Number).Int64()
		if err != nil {
			continue
		}
		utxo.Vout = int(i64)

		scriptPubKeyValue := scriptPubKey.(string)
		utxo.ScriptPubKey = scriptPubKeyValue

		i64, err = confirmations.(json.Number).Int64()
		if err != nil {
			continue
		}
		utxo.Confirmations = int(i64)

		utxos = append(utxos, utxo)
	}

	// sort by confirmations desc
	sort.Sort(utxos)

	return utxos, nil
}

func (agent *OMNIAgent) ImportAddressRPC(address string) error {

	res, err := agent.DoHttpJsonRpcCallType1("importaddress", address, "", false)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.Message)
	}
	return nil
}

func (agent *OMNIAgent) BroadcastTransactionRPC(rawtrx string) (string, error) {
	res, err := agent.DoHttpJsonRpcCallType1("sendrawtransaction", rawtrx)
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

func (agent *OMNIAgent) IsFeeTransactionConfirmedRPC(trxId string) (bool, error) {
	res, err := agent.DoHttpJsonRpcCallType1("gettransaction", trxId)
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
	cfm, err := resmap["confirmations"].(json.Number).Int64()
	if err != nil {
		return false, err
	}

	coin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return false, errors.New("not support coin")
	}
	if cfm >= int64(coin.ConfirmCount) {
		return true, nil
	}
	return false, nil

}

func (agent *OMNIAgent) IsTransactionConfirmedRPC(trxId string) (bool, error) {
	res, err := agent.DoHttpJsonRpcCallType1("omni_gettransaction", trxId)
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
	cfm, err := resmap["confirmations"].(json.Number).Int64()
	if err != nil {
		return false, err
	}

	isvalid, isfound := resmap["valid"]
	if !isfound {
		return false, nil
	}
	trxvalid := isvalid.(bool)
	if !trxvalid {
		return true, errors.New("invalid transaction")
	}

	coin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return false, errors.New("not support coin")
	}
	if cfm >= int64(coin.ConfirmCount) {
		return true, nil
	}
	return false, nil

}

func (agent *OMNIAgent) GetRawTransactionRPC(trxid string) (string, error) {
	res, err := agent.DoHttpJsonRpcCallType1("getrawtransaction", trxid)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", errors.New(res.Error.Message)
	}
	rawtx, err := res.GetString()
	if err != nil {
		return "", nil
	}
	return rawtx, err
}

func (agent *OMNIAgent) IsAddressValidRPC(address string) (bool, error) {
	res, err := agent.DoHttpJsonRpcCallType1("validateaddress", address)
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
	isValid := resmap["isvalid"].(bool)
	if err != nil {
		return false, err
	}
	return isValid, nil

}

func OMNIGetUnCompressPubKey(pubKeyBytes []byte) ([]byte, error) {
	if len(pubKeyBytes) != 64 {
		return nil, errors.New("invalid pubKeyBytes size")
	}

	pubkeyUnCompress := make([]byte, 65, 65)
	pubkeyUnCompress[0] = 0x4
	copy(pubkeyUnCompress[1:], pubKeyBytes[0:64])

	return pubkeyUnCompress, nil
}

func OMNIGetCompressPubKey(pubKeyBytes []byte) ([]byte, error) {
	if len(pubKeyBytes) != 64 {
		return nil, errors.New("invalid pubKeyBytes size")
	}

	pubkeyCompress := make([]byte, 33, 33)
	if pubKeyBytes[63]%2 == 0 {
		pubkeyCompress[0] = 0x2
	} else {
		pubkeyCompress[0] = 0x3
	}
	copy(pubkeyCompress[1:], pubKeyBytes[0:32])

	fmt.Println("pubkeyCompress:", hex.EncodeToString(pubkeyCompress))

	return pubkeyCompress, nil
}

func OMNICalcAddressByPubKey(pubKeyStr string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return "", err
	}

	pubkeyCompress, err := OMNIGetCompressPubKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	pubKey := new(pubkey.PubKey)
	pubKey.SetPubKeyData(pubkeyCompress)

	keyIdBytes, err := pubKey.CalcKeyIDBytes()
	if err != nil {
		return "", err
	}
	keyId := new(keyid.KeyID)
	keyId.SetKeyIDData(keyIdBytes)

	var version byte
	if config.IsTestEnvironment {
		version = 111
	} else {
		version = 0
	}
	addrStr, err := keyId.ToBase58Address(version)
	if err != nil {
		return "", err
	}

	return addrStr, nil
}

func OMNIGetRedeemScriptByPubKeys(needCount int, pubKeyStrList []string) (string, error) {
	if needCount <= 0 || needCount > 16 {
		return "", errors.New("OMNIGetRedeemScriptByPubKeys error: invalid needCount")
	}
	if len(pubKeyStrList) == 0 || len(pubKeyStrList) > 16 {
		return "", errors.New("OMNIGetRedeemScriptByPubKeys error: invalid pubKeyStrList size")
	}
	if needCount > len(pubKeyStrList) {
		return "", errors.New("OMNIGetRedeemScriptByPubKeys error: needCount greater than pubKeyStrList size")
	}

	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := serialize.PackUint8(bufWriter, uint8(needCount+0x50))
	if err != nil {
		return "", err
	}
	for _, pubKeyStr := range pubKeyStrList {
		pubKeyBytes, err := hex.DecodeString(pubKeyStr)
		if err != nil {
			return "", err
		}

		pubKeyCpsBytes, err := OMNIGetCompressPubKey(pubKeyBytes)
		if err != nil {
			return "", err
		}

		pubKey := new(pubkey.PubKey)
		pubKey.SetPubKeyData(pubKeyCpsBytes)

		err = pubKey.Pack(bufWriter)
		if err != nil {
			return "", err
		}
	}
	err = serialize.PackUint8(bufWriter, uint8(len(pubKeyStrList)+0x50))
	if err != nil {
		return "", err
	}
	err = serialize.PackUint8(bufWriter, uint8(script.OP_CHECKMULTISIG))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytesBuf.Bytes()), nil
}

func OMNIGetMultiSignAddressByRedeemScript(redeemScriptStr string) (string, error) {
	redeemScript, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		return "", err
	}

	scriptIdBytes := utility.Hash160(redeemScript)
	keyId := new(keyid.KeyID)
	keyId.SetKeyIDData(scriptIdBytes)

	var version byte
	if config.IsTestEnvironment {
		version = 196
	} else {
		version = 5
	}

	addrStr, err := keyId.ToBase58Address(version)
	if err != nil {
		return "", err
	}
	return addrStr, nil
}

func OMNIGetMultiSignScriptByRedeemScript(redeemScriptStr string) ([]byte, error) {
	redeemScript, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		return nil, err
	}
	scriptIdBytes := utility.Hash160(redeemScript)

	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err = serialize.PackUint8(bufWriter, script.OP_HASH160)
	if err != nil {
		return nil, err
	}
	multiSigScript := new(script.Script)
	multiSigScript.SetScriptBytes(scriptIdBytes)
	err = multiSigScript.Pack(bufWriter)
	if err != nil {
		return nil, err
	}
	err = serialize.PackUint8(bufWriter, script.OP_EQUAL)
	if err != nil {
		return nil, err
	}

	return bytesBuf.Bytes(), nil
}

func (agent *OMNIAgent) BuildTrxInPutsOutPutsRPC(addrFromStr string, addrToStr string, amountTransferStr string, feeRateStr string) (string, InPuts, OutPuts, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return "0", nil, nil, errors.New("not support coin")
	}
	nPrec := supportCoin.Precision

	feeCoin, ok := config.GlobalSupportCoinMgr["BTC"]
	if !ok {
		return "0", nil, nil, errors.New("not support fee coin")
	}
	nFeePrec := feeCoin.Precision

	balanceStr, err := agent.GetBalanceByAddressRPC(addrFromStr)
	if err != nil {
		return "0", nil, nil, err
	}
	balance, err := ToPrecisionAmount(balanceStr, nPrec)
	if err != nil {
		return "0", nil, nil, err
	}
	amountTransfer, err := ToPrecisionAmount(amountTransferStr, nPrec)
	if err != nil {
		return "0", nil, nil, err
	}
	feeRate, err := ToPrecisionAmount(feeRateStr, nFeePrec)
	if err != nil {
		return "0", nil, nil, err
	}

	if balance <= amountTransfer {
		return "0", nil, nil, errors.New("not enough balance")
	}

	utxos, err := agent.GetUtxosByAddressRPC(addrFromStr)
	if err != nil {
		return "0", nil, nil, err
	}

	spentBalance := int64(0)
	change := int64(0)
	feeCost := int64(0)
	dustTransfer := int64(546)
	trxBytes := 0
	balanceOk := false
	inputs := make([]InPut, 0)
	outputs := make(map[string]string)

	// trx size 100k
	for _, utxo := range utxos {
		// ignore dust
		if utxo.Amount < 546 {
			continue
		}
		inputs = append(inputs, InPut{TxId: utxo.TxId, Vout: utxo.Vout, ScriptPubKey: utxo.ScriptPubKey})
		spentBalance = spentBalance + utxo.Amount

		// vout: dest, change, opreturn
		trxBytes = len(inputs)*180 + 40 + 40 + 40
		if trxBytes > 100*1000 {
			return "0", nil, nil, errors.New("too large trx size")
		}
		feeCost = int64(float64(trxBytes) / 1000.0 * float64(feeRate))
		if spentBalance >= dustTransfer+feeCost {
			balanceOk = true
			break
		}
	}
	if balanceOk != true {
		return "0", nil, nil, errors.New("not enough fee balance")
	}

	if addrToStr == addrFromStr {
		return "0", nil, nil, errors.New("not support to transfer omni property to myself")
	} else {
		outputs[addrToStr] = fmt.Sprintf(FromPrecisionAmount(dustTransfer, nPrec))
		change = spentBalance - dustTransfer - feeCost
	}

	if change > 546 {
		outputs[addrFromStr] = fmt.Sprintf(FromPrecisionAmount(change, nPrec))
	}

	return FromPrecisionAmount(feeCost, nPrec), inputs, outputs, nil
}

func (agent *OMNIAgent) BuildTrxFromUtxosRPC(utxos []UTXODetail, addrFromStr string, addrToStr string, amountTransferStr string, feeRateStr string) (string, InPuts, OutPuts, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return "0", nil, nil, errors.New("not support coin")
	}
	nPrec := supportCoin.Precision

	feeCoin, ok := config.GlobalSupportCoinMgr["BTC"]
	if !ok {
		return "0", nil, nil, errors.New("not support fee coin")
	}
	nFeePrec := feeCoin.Precision

	// different with other utxo-type coins
	balanceStr, err := agent.GetBalanceByAddressRPC(addrFromStr)
	if err != nil {
		return "0", nil, nil, err
	}
	balance, err := ToPrecisionAmount(balanceStr, nPrec)
	if err != nil {
		return "0", nil, nil, err
	}
	amountTransfer, err := ToPrecisionAmount(amountTransferStr, nPrec)
	if err != nil {
		return "0", nil, nil, err
	}
	feeRate, err := ToPrecisionAmount(feeRateStr, nFeePrec)
	if err != nil {
		return "0", nil, nil, err
	}

	if balance <= amountTransfer {
		return "0", nil, nil, errors.New("not enough balance")
	}

	spentBalance := int64(0)
	change := int64(0)
	feeCost := int64(0)
	dustTransfer := int64(546)
	trxBytes := 0
	balanceOk := false
	inputs := make([]InPut, 0)
	outputs := make(map[string]string)

	// trx size 100k
	for _, utxo := range utxos {
		if utxo.Address != addrFromStr {
			continue
		}
		// ignore dust
		if utxo.Amount < 546 {
			continue
		}
		inputs = append(inputs, InPut{TxId: utxo.TxId, Vout: utxo.Vout, ScriptPubKey: utxo.ScriptPubKey})
		spentBalance = spentBalance + utxo.Amount

		// vout: dest, change, opreturn
		trxBytes = len(inputs)*180 + 40 + 40 + 40
		if trxBytes > 100*1000 {
			return "0", nil, nil, errors.New("too large trx size")
		}
		feeCost = int64(float64(trxBytes) / 1000.0 * float64(feeRate))
		if spentBalance >= dustTransfer+feeCost {
			balanceOk = true
			break
		}
	}
	if balanceOk != true {
		return "0", nil, nil, errors.New("not enough fee balance")
	}

	if addrToStr == addrFromStr {
		return "0", nil, nil, errors.New("not support to transfer omni property to myself")
	} else {
		outputs[addrToStr] = fmt.Sprintf(FromPrecisionAmount(dustTransfer, nPrec))
		change = spentBalance - dustTransfer - feeCost
	}

	if change > 546 {
		outputs[addrFromStr] = fmt.Sprintf(FromPrecisionAmount(change, nPrec))
	}

	return FromPrecisionAmount(feeCost, nPrec), inputs, outputs, nil
}

func (agent *OMNIAgent) CreateRawTransactionRPC(inputs InPuts, outputs OutPuts) (string, error) {
	res, err := agent.DoHttpJsonRpcCallType1("createrawtransaction", inputs, outputs)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", errors.New(res.Error.Message)
	}
	return res.Result.(string), nil
}

func (agent *OMNIAgent) CreateRawTransactionOpReturnRPC(rawTrx string, msgType uint16, msgVersion uint16,
	propertyId uint32, amountTransferStr string) (string, error) {
	supportCoin, ok := config.GlobalSupportCoinMgr[agent.CoinType()]
	if !ok {
		return "", errors.New("not support coin")
	}
	nPrec := supportCoin.Precision
	amountTransfer, err := ToPrecisionAmount(amountTransferStr, nPrec)
	if err != nil {
		return "", err
	}

	payLoadFormat := "%04x%04x%08x%016x"
	payLoad := fmt.Sprintf(payLoadFormat, msgType, msgVersion, propertyId, amountTransfer)
	res, err := agent.DoHttpJsonRpcCallType1("omni_createrawtx_opreturn", rawTrx, payLoad)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", errors.New(res.Error.Message)
	}
	return res.Result.(string), nil
}

func OMNIUnPackRawTransaction(rawTrx string) (*transaction.Transaction, error) {
	Blob := new(blob.Byteblob)
	err := Blob.SetHex(rawTrx)
	if err != nil {
		return nil, err
	}
	bytesBuf := bytes.NewBuffer(Blob.GetData())
	bufReader := io.Reader(bytesBuf)
	trx := new(transaction.Transaction)
	err = trx.UnPack(bufReader)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func OMNIPackRawTransaction(trxSig transaction.Transaction) (string, error) {
	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := trxSig.Pack(bufWriter)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytesBuf.Bytes()), nil
}

func OMNICombineSignatureAndPubKey(signature []byte, pubKey []byte) []byte {
	scriptSig := make([]byte, 0, 1+len(signature)+1+len(pubKey))
	scriptSig = append(scriptSig, byte(len(signature)))
	scriptSig = append(scriptSig, signature...)
	scriptSig = append(scriptSig, byte(len(pubKey)))
	scriptSig = append(scriptSig, pubKey...)
	fmt.Println("scriptSig ", hex.EncodeToString(scriptSig))
	return scriptSig
}

func OMNICombineSignatureAndRedeemScript(signature []byte, redeemScriptBytes []byte) ([]byte, error) {
	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err := serialize.PackUint8(bufWriter, script.OP_0)
	if err != nil {
		return nil, err
	}
	signatureScript := new(script.Script)
	signatureScript.SetScriptBytes(signature)
	err = signatureScript.Pack(bufWriter)
	if err != nil {
		return nil, err
	}

	if len(redeemScriptBytes) < int(script.OP_PUSHDATA1) {
	} else {
		opPushData := uint8(0)
		if len(redeemScriptBytes) <= 0xff {
			opPushData = script.OP_PUSHDATA1
		} else if len(redeemScriptBytes) <= 0xffff {
			opPushData = script.OP_PUSHDATA2
		} else {
			opPushData = script.OP_PUSHDATA4
		}
		err = serialize.PackUint8(bufWriter, opPushData)
		if err != nil {
			return nil, err
		}
	}

	redeemScrip := new(script.Script)
	redeemScrip.SetScriptBytes(redeemScriptBytes)
	err = redeemScrip.Pack(bufWriter)
	if err != nil {
		return nil, err
	}
	return bytesBuf.Bytes(), nil
}

func (agent *OMNIAgent) SignRawTransactionRPC(rawTrx string, pubKeyStr string, keyIndex uint16, utxos []UTXODetail) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return "", err
	}
	pubkeyCompress, err := OMNIGetCompressPubKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	trx, err := OMNIUnPackRawTransaction(rawTrx)
	if err != nil {
		return "", err
	}

	signedDataList := make([][]byte, len(trx.Vin))

	// add scriptPubKey
	for i := 0; i < len(trx.Vin); i++ {
		trxTemp, err := OMNIUnPackRawTransaction(rawTrx)
		if err != nil {
			return "", err
		}

		vinFound := false
		for j := 0; j < len(utxos); j++ {
			if trxTemp.Vin[i].PrevOut.Hash.GetHex() == utxos[j].TxId {
				vinFound = true
				script, err := hex.DecodeString(utxos[j].ScriptPubKey)
				if err != nil {
					return "", errors.New("invalid ScriptPubKey")
				}
				trxTemp.Vin[i].ScriptSig.SetScriptBytes(script)
				fmt.Println("script len", len(script))
				break
			}
		}
		if vinFound != true {
			return "", errors.New("can not found valid utxo for rawTrx")
		}

		rawTrxWithScript, err := OMNIPackRawTransaction(*trxTemp)
		if err != nil {
			return "", err
		}

		rawTrxBytes, err := hex.DecodeString(rawTrxWithScript)
		if err != nil {
			return "", err
		}
		// append SIGHASH_ALL
		rawTrxBytes = append(rawTrxBytes, []byte{0x1, 0x0, 0x0, 0x0}...)
		hashBytes := utility.Sha256(utility.Sha256(rawTrxBytes))

		fmt.Println("rawTrxBytes:", hex.EncodeToString(rawTrxBytes))
		fmt.Println("hashBytes:", hex.EncodeToString(hashBytes))

		// signature
		var rsBytes []byte

		for {
			rsBytes, err = CoinSignTrx('1', hashBytes, keyIndex)
			if err != nil {
				return "", err
			}
			if len(rsBytes) != 64 {
				return "", errors.New("invalid r/s lens")
			}
			if rsBytes[32] < 128 {
				break
			}
		}

		verifyOk, err := CoinVerifyTrx('1', keyIndex, hashBytes, rsBytes)
		if err != nil {
			return "", err
		}
		if !verifyOk {
			return "", errors.New("verify signature error")
		}
		fmt.Println("rsHex:", hex.EncodeToString(rsBytes))

		rBytes := rsBytes[0:32]
		sBytes := rsBytes[32:64]

		// serialize r,s to der encoding
		signedData, err := SerializeDerEncoding(rBytes, sBytes)
		if err != nil {
			return "", err
		}
		fmt.Println("signedData:", hex.EncodeToString(signedData))

		// append SIGHASH_ALL
		signedData = append(signedData, 0x1)

		scriptSig := OMNICombineSignatureAndPubKey(signedData, pubkeyCompress)

		signedDataList[i] = scriptSig
	}

	for i := 0; i < len(trx.Vin); i++ {
		trx.Vin[i].ScriptSig.SetScriptBytes(signedDataList[i])
	}

	trxSigStr, err := OMNIPackRawTransaction(*trx)
	if err != nil {
		return "", err
	}

	return trxSigStr, nil
}

func (agent *OMNIAgent) MultiSignRawTransactionRPC(rawTrx string, redeemScriptStr string, keyIndex uint16, utxos []UTXODetail) (string, error) {
	redeemScriptBytes, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		return "", err
	}

	trx, err := OMNIUnPackRawTransaction(rawTrx)
	if err != nil {
		return "", err
	}

	signedDataList := make([][]byte, len(trx.Vin))

	redeemScript, err := hex.DecodeString(redeemScriptStr)
	if err != nil {
		return "", errors.New("invalid redeemScript")
	}

	// add scriptPubKey
	for i := 0; i < len(trx.Vin); i++ {
		trxTemp, err := OMNIUnPackRawTransaction(rawTrx)
		if err != nil {
			return "", err
		}

		trxTemp.Vin[i].ScriptSig.SetScriptBytes(redeemScript)

		rawTrxWithScript, err := OMNIPackRawTransaction(*trxTemp)
		if err != nil {
			return "", err
		}

		rawTrxBytes, err := hex.DecodeString(rawTrxWithScript)
		if err != nil {
			return "", err
		}
		// append SIGHASH_ALL
		rawTrxBytes = append(rawTrxBytes, []byte{0x1, 0x0, 0x0, 0x0}...)
		hashBytes := utility.Sha256(utility.Sha256(rawTrxBytes))

		fmt.Println("rawTrxBytes:", hex.EncodeToString(rawTrxBytes))
		fmt.Println("hashBytes:", hex.EncodeToString(hashBytes))

		// signature
		var rsBytes []byte

		for {
			rsBytes, err = CoinSignTrx('1', hashBytes, keyIndex)
			if err != nil {
				return "", err
			}
			if len(rsBytes) != 64 {
				return "", errors.New("invalid r/s lens")
			}
			if rsBytes[32] < 128 {
				break
			}
		}

		fmt.Println("hashBytes", hex.EncodeToString(hashBytes))
		fmt.Println("rsBytes", hex.EncodeToString(rsBytes))

		verifyOk, err := CoinVerifyTrx('1', keyIndex, hashBytes, rsBytes)
		if err != nil {
			return "", err
		}
		if !verifyOk {
			return "", errors.New("verify signature error")
		}
		fmt.Println("rsHex:", hex.EncodeToString(rsBytes))

		rBytes := rsBytes[0:32]
		sBytes := rsBytes[32:64]

		// serialize r,s to der encoding
		signedData, err := SerializeDerEncoding(rBytes, sBytes)
		if err != nil {
			return "", err
		}
		fmt.Println("signedData:", hex.EncodeToString(signedData))

		// append SIGHASH_ALL
		signedData = append(signedData, 0x1)

		scriptSig, err := OMNICombineSignatureAndRedeemScript(signedData, redeemScriptBytes)
		if err != nil {
			return "", err
		}
		signedDataList[i] = scriptSig
	}

	for i := 0; i < len(trx.Vin); i++ {
		trx.Vin[i].ScriptSig.SetScriptBytes(signedDataList[i])
	}

	trxSigStr, err := OMNIPackRawTransaction(*trx)
	if err != nil {
		return "", err
	}

	return trxSigStr, nil
}

func (agent *OMNIAgent) MultiVerifySignRawTransactionRPC(signedRawTrx string, pubKeyStr string, utxos []UTXODetail) error {
	trx, err := OMNIUnPackRawTransaction(signedRawTrx)
	if err != nil {
		return err
	}

	signedDataList := make([]script.Script, 0, len(trx.Vin))
	for i := 0; i < len(trx.Vin); i++ {
		signedDataList = append(signedDataList, trx.Vin[i].ScriptSig)
	}

	for i := 0; i < len(trx.Vin); i++ {
		trxTemp, err := OMNIUnPackRawTransaction(signedRawTrx)
		if err != nil {
			return err
		}
		// set empty script
		for i := 0; i < len(trx.Vin); i++ {
			emptyScript := new(script.Script)
			trxTemp.Vin[i].ScriptSig.SetScriptBytes(emptyScript.GetScriptBytes())
		}

		bytesBuf := bytes.NewBuffer(signedDataList[i].GetScriptBytes())
		bufReader := io.Reader(bytesBuf)
		u8, err := serialize.UnPackUint8(bufReader)
		if err != nil {
			return err
		}
		if u8 != 0 {
			return errors.New("invalid multisig script, not started with 0x0")
		}

		signatureScript := new(script.Script)
		err = signatureScript.UnPack(bufReader)
		if err != nil {
			return err
		}

		signatureScriptBytes := signatureScript.GetScriptBytes()
		if signatureScriptBytes[len(signatureScriptBytes)-1] != 0x1 {
			return errors.New("invalid signature, not ended with 0x1[ALL]")
		}
		signatureScriptBytes = signatureScriptBytes[:len(signatureScriptBytes)-1]

		// skip oppushdata
		tmpBufReader := bufio.NewReader(bufReader)
		opPushDataBytes, err := tmpBufReader.Peek(1)
		if err != nil {
			return err
		}

		bufReader = io.Reader(tmpBufReader)
		if opPushDataBytes[0] == script.OP_PUSHDATA1 || opPushDataBytes[0] == script.OP_PUSHDATA2 || opPushDataBytes[0] == script.OP_PUSHDATA4 {
			_, err = serialize.UnPackUint8(bufReader)
			if err != nil {
				return err
			}
		}

		redeemScript := new(script.Script)
		err = redeemScript.UnPack(bufReader)
		if err != nil {
			return err
		}

		// get compress address
		pubKeyBytes, err := hex.DecodeString(pubKeyStr)
		if err != nil {
			return err
		}
		compressPubKey, err := OMNIGetCompressPubKey(pubKeyBytes)
		if err != nil {
			return err
		}
		solverOk, whichType, PubKeys := script.Solver(*redeemScript)
		if solverOk != true || whichType != script.TX_MULTISIG || len(PubKeys) == 0 {
			return errors.New("invalid redeemscript, solver error")
		}

		// is compressPubKey in PubKeys
		isInPubKeys := false
		for _, pubKey := range PubKeys {
			if 0 == bytes.Compare(compressPubKey, pubKey) {
				isInPubKeys = true
				break
			}
		}
		if !isInPubKeys {
			return errors.New("arguments pubKeyStr of function MultiVerifySignRawTransaction is no match to redeemscript")
		}

		// set Vin script pubkey
		trxTemp.Vin[i].ScriptSig.SetScriptBytes(redeemScript.GetScriptBytes())

		rawTrxWithScript, err := OMNIPackRawTransaction(*trxTemp)
		if err != nil {
			return err
		}

		rawTrxBytes, err := hex.DecodeString(rawTrxWithScript)
		if err != nil {
			return err
		}
		// append SIGHASH_ALL
		rawTrxBytes = append(rawTrxBytes, []byte{0x1, 0x0, 0x0, 0x0}...)
		hashBytes := utility.Sha256(utility.Sha256(rawTrxBytes))

		fmt.Println("rawTrxBytes:", hex.EncodeToString(rawTrxBytes))
		fmt.Println("hashBytes:", hex.EncodeToString(hashBytes))
		fmt.Println("signatureScriptBytes", hex.EncodeToString(signatureScriptBytes))

		// do not verify by crypto device
		//verifyOk, err := CoinVerifyTrxWithOutsidePubkey('2', pubKeyBytes, hashBytes, signatureScriptBytes)

		verifyOk, err := CoinVerifyTrx2(compressPubKey, hashBytes, signatureScriptBytes)
		if err != nil {
			return err
		}
		if !verifyOk {
			return errors.New("verify signature error")
		}
	}

	return nil
}

func (agent *OMNIAgent) CombineRawTransactionRPC(signedRawTrxs []string) (string, error) {
	trxs := make([]*transaction.Transaction, 0)
	for _, signedRawTrx := range signedRawTrxs {
		trx, err := OMNIUnPackRawTransaction(signedRawTrx)
		if err != nil {
			return "", err
		}
		trxs = append(trxs, trx)
	}

	vinCount := 0
	if len(trxs) == 0 {
		return "", errors.New("transaction count is 0")
	} else if len(trxs) == 1 {
		return signedRawTrxs[0], nil
	} else {
		vinCount = len(trxs[0].Vin)
		for i := 1; i < len(trxs); i++ {
			if vinCount != len(trxs[i].Vin) {
				return "", errors.New("signedRawTrxs with different vins")
			}
		}
	}

	tempTrx, err := OMNIUnPackRawTransaction(signedRawTrxs[0])
	if err != nil {
		return "", err
	}

	validRedeemScript := false
	redeemScript := new(script.Script)
	for i := 0; i < vinCount; i++ {
		signatureScripts := make([]*script.Script, 0)
		for j := 0; j < len(trxs); j++ {
			bytesBuf := bytes.NewBuffer(trxs[j].Vin[i].ScriptSig.GetScriptBytes())
			bufReader := io.Reader(bytesBuf)
			u8, err := serialize.UnPackUint8(bufReader)
			if err != nil {
				return "", err
			}
			if u8 != 0 {
				return "", errors.New("invalid multisig script, not started with 0x0")
			}

			signatureScript := new(script.Script)
			err = signatureScript.UnPack(bufReader)
			if err != nil {
				return "", err
			}

			signatureScriptBytes := signatureScript.GetScriptBytes()
			if signatureScriptBytes[len(signatureScriptBytes)-1] != 0x1 {
				return "", errors.New("invalid signature, not ended with 0x1[ALL]")
			}
			signatureScriptBytes = signatureScriptBytes[:len(signatureScriptBytes)-1]

			// skip oppushdata
			tmpBufReader := bufio.NewReader(bufReader)
			opPushDataBytes, err := tmpBufReader.Peek(1)
			if err != nil {
				return "", err
			}

			if opPushDataBytes[0] == script.OP_PUSHDATA1 || opPushDataBytes[0] == script.OP_PUSHDATA2 || opPushDataBytes[0] == script.OP_PUSHDATA4 {
				bufReader = io.Reader(tmpBufReader)
				_, err = serialize.UnPackUint8(bufReader)
				if err != nil {
					return "", err
				}
			}

			if !validRedeemScript {
				err = redeemScript.UnPack(bufReader)
				if err != nil {
					return "", err
				}
				validRedeemScript = true
			}

			signatureScripts = append(signatureScripts, signatureScript)
		}

		// combine signature
		bytesBuf := bytes.NewBuffer([]byte{})
		bufWriter := io.Writer(bytesBuf)
		err := serialize.PackUint8(bufWriter, 0)
		if err != nil {
			return "", err
		}

		for _, signatureScript := range signatureScripts {
			err = signatureScript.Pack(bufWriter)
			if err != nil {
				return "", err
			}
		}

		// oppushdata
		opPushData := uint8(0)
		if redeemScript.GetScriptLength() < int(script.OP_PUSHDATA1) {
			opPushData = uint8(0)
		} else if redeemScript.GetScriptLength() >= int(script.OP_PUSHDATA1) && redeemScript.GetScriptLength() <= 0xff {
			opPushData = script.OP_PUSHDATA1
		} else if redeemScript.GetScriptLength() > 0xff && redeemScript.GetScriptLength() <= 0xffff {
			opPushData = script.OP_PUSHDATA2
		} else if redeemScript.GetScriptLength() > 0xffff {
			opPushData = script.OP_PUSHDATA4
		}

		if opPushData != uint8(0) {
			err := serialize.PackUint8(bufWriter, opPushData)
			if err != nil {
				return "", err
			}
		}

		// redeemscript
		err = redeemScript.Pack(bufWriter)
		if err != nil {
			return "", err
		}

		tempTrx.Vin[i].ScriptSig.SetScriptBytes(bytesBuf.Bytes())
	}

	bytesBuf := bytes.NewBuffer([]byte{})
	bufWriter := io.Writer(bytesBuf)
	err = tempTrx.Pack(bufWriter)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytesBuf.Bytes()), nil
}
