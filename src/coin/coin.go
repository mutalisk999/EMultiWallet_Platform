package coin

import (
	"fmt"
	"github.com/Messer4/bchaddr"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kataras/iris/core/errors"
	"config"
	"github.com/mutalisk999/bitcoin-lib/src/script"
	"model"
	"utils"
	"github.com/btcsuite/btcd/btcec"
	"github.com/mutalisk999/go-lib/src/net/buffer_tcp"
	"github.com/mutalisk999/go-cryptocard/src/cryptocard"
	"github.com/go-xorm/xorm"
)


func GetAddressByPubKey(pubKeyStr string, coinSymbol string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "", utils.InvalidCoinSymbol(coinSymbol)
	}

	if coinSymbol == "BTC" {
		addrStr, err := BTCCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	}else if coinSymbol == "ETH"{
		addrStr,err := ETHCalcAddressByPubkey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinSymbol == "LTC" {
		addrStr, err := LTCCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinSymbol == "BCH" {
		addrStr, err := BCHCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.IsErc20{
		addrStr, err := ETHCalcAddressByPubkey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.IsOmni {
		addrStr, err := OMNICalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.CoinSymbol == "UB" {
		addrStr, err := UBCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	}
	return "", nil
}

func GetRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, trxid string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// broadcast transaction
		trxSigHex, err := agent.GetRawTransactionRPC(trxid)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// getrawtransaction
		trxSigHex, err := agent.GetRawTransactionRPC(trxid)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	}
	return "", nil
}

func GetKeyID(coinSymbol string ,addr string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "UB" {
		keystr := base58.Decode(addr)
		return string(keystr[1:21]), nil
	} else if coinSymbol == "BCH"{
		keysBase58,err := bchaddr.ToLegacyAddress(addr)
		if err != nil {
			return "", err
		}
		keystr := base58.Decode(keysBase58)
		return string(keystr[1:21]), nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		keystr := base58.Decode(addr)
		return string(keystr[1:21]), nil
	}
	return "", nil
}

func ExtractDestination(scriptPubKey script.Script) (bool,[]byte) {
	_, whichType, vSolutions := script.Solver(scriptPubKey)

	// p2pkh
	if whichType == script.TX_PUBKEYHASH {
		return true, vSolutions[0]
	}

	// p2sh
	if whichType == script.TX_SCRIPTHASH {
		return true, vSolutions[0]
	}

	return false, vSolutions[0]
}

func GetMultiSignRedeemScript(needCount int, pubKeyStrSlice []string, coinSymbol string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "", utils.InvalidCoinSymbol(coinSymbol)
	}

	if coinSymbol == "BTC" {
		redeemScriptStr, err := BTCGetRedeemScriptByPubKeys(needCount, pubKeyStrSlice)
		if err != nil {
			return "", err
		}
		return redeemScriptStr, nil
	}else if coinSymbol == "ETH"{
		return "", errors.New("GetMultiSignRedeemScript not support ETH")
	} else if coinSymbol == "LTC" {
		redeemScriptStr, err := LTCGetRedeemScriptByPubKeys(needCount, pubKeyStrSlice)
		if err != nil {
			return "", err
		}
		return redeemScriptStr, nil
	} else if coinSymbol == "BCH" {
		redeemScriptStr, err := BCHGetRedeemScriptByPubKeys(needCount, pubKeyStrSlice)
		if err != nil {
			return "", err
		}
		return redeemScriptStr, nil
	} else if coinObj.IsErc20{
		return "", errors.New("GetMultiSignRedeemScript not support ERC20")
	} else if coinObj.IsOmni {
		redeemScriptStr, err := OMNIGetRedeemScriptByPubKeys(needCount, pubKeyStrSlice)
		if err != nil {
			return "", err
		}
		return redeemScriptStr, nil
	} else if coinObj.CoinSymbol == "UB" {
		redeemScriptStr, err := UBGetRedeemScriptByPubKeys(needCount, pubKeyStrSlice)
		if err != nil {
			return "", err
		}
		return redeemScriptStr, nil
	}
	return "", nil
}

func GetMultiSignAddress(redeemScriptStr string, coinSymbol string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "", utils.InvalidCoinSymbol(coinSymbol)
	}

	if coinSymbol == "BTC" {
		multiSignAddress, err := BTCGetMultiSignAddressByRedeemScript(redeemScriptStr)
		if err != nil {
			return "", err
		}
		return multiSignAddress, nil
	}else if coinSymbol == "ETH"{
		return "", errors.New("GetMultiSignAddress not support ETH")
	} else if coinSymbol == "LTC" {
		multiSignAddress, err := LTCGetMultiSignAddressByRedeemScript(redeemScriptStr)
		if err != nil {
			return "", err
		}
		return multiSignAddress, nil
	} else if coinSymbol == "BCH" {
		multiSignAddress, err := BCHGetMultiSignAddressByRedeemScript(redeemScriptStr)
		if err != nil {
			return "", err
		}
		return multiSignAddress, nil
	} else if coinObj.IsErc20{
		return "", errors.New("GetMultiSignAddress not support ERC20")
	} else if coinObj.IsOmni {
		multiSignAddress, err := OMNIGetMultiSignAddressByRedeemScript(redeemScriptStr)
		if err != nil {
			return "", err
		}
		return multiSignAddress, nil
	} else if coinObj.CoinSymbol == "UB" {
		multiSignAddress, err := UBGetMultiSignAddressByRedeemScript(redeemScriptStr)
		if err != nil {
			return "", err
		}
		return multiSignAddress, nil
	}
	return "", nil
}

func GetMultiSignAddressByPubKeys(needCount int, pubKeyStrSlice []string, coinSymbol string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "", utils.InvalidCoinSymbol(coinSymbol)
	}

	if coinSymbol == "ETH" || coinObj.IsErc20  {
		return "", nil
	}

	redeemScriptStr, err := GetMultiSignRedeemScript(needCount, pubKeyStrSlice, coinSymbol)
	if err != nil {
		return "", err
	}
	multiSignAddress, err := GetMultiSignAddress(redeemScriptStr, coinSymbol)
	if err != nil {
		return "", err
	}
	return multiSignAddress, nil
}

func ImportAddress(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		err := agent.ImportAddressRPC(address)
		if err != nil {
			return err
		}
	} else if coinSymbol == "ETH"||coinObj.IsErc20{
		return nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		err := agent.ImportAddressRPC(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsTrxConfirmed(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, trxId string) (bool, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  false,utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmedRPC(trxId)
		return isConfirmed, err
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmed(trxId)
		return isConfirmed, err
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmed(trxId)
		return isConfirmed, err
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmedRPC(trxId)
		return isConfirmed, err
	}
	return false, nil
}

func IsAddressValid(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (bool, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  false,utils.InvalidCoinSymbol(coinSymbol)
	}
	if address ==""{
		return true,nil
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isValid, err := agent.IsAddressValidRPC(address)
		return isValid, err
	}else if coinSymbol == "ETH"{
		isValied := ETHValidateAddress(address)
		return isValied, nil
	}else if coinObj.IsErc20{
		isValied := ETHValidateAddress(address)
		return isValied, nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isValid, err := agent.IsAddressValidRPC(address)
		return isValid, err
	}
	return false, nil
}

func ConfirmTransaction(dbSession *xorm.Session, coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string,
	trxId int, rawTrxId string, trxState int) (error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" || coinObj.IsOmni {
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err := trxMgr.UpdateTransactionStateFeeCost(dbSession, trxId, trxState, nil)
		if err != nil {
			return err
		}
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		real_cost,err := agent.GetTransactionRealCost(coinSymbol,rawTrxId)
		if err != nil {
			return err
		}
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err = trxMgr.UpdateTransactionStateFeeCost(dbSession, trxId, trxState, &real_cost)
		if err != nil {
			return err
		}
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		real_cost,err := agent.GetTransactionRealCost(coinSymbol,rawTrxId)
		if err != nil {
			return err
		}
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err = trxMgr.UpdateTransactionStateFeeCost(dbSession, trxId, trxState, &real_cost)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBalance(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (string,string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		balance, err := agent.GetBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}

		return balance,"",nil
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		balance, err := agent.GetBalanceByAddress(address)
		if err != nil {
			return "","", err
		}

		return ConvertBigNumberToString(coinSymbol,balance),"", nil
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		erc20_balance,balance, err := agent.GetERC20BalanceByAddress(coinSymbol,address)
		if err != nil {
			return "","", err
		}

		return ConvertBigNumberToString(coinSymbol,erc20_balance),ConvertBigNumberToString("ETH",balance), nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		balance, err := agent.GetBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}
		feeBalance, err := agent.GetFeeBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}

		return balance,feeBalance,nil
	}
	return "","", nil
}

func Transfer(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, keyIndex uint16, pubKeyStr string,
	from string, to string, amountStr string, feeStr string, gasPrice string, gasLimit string) (string, string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get utxos
		utxos, err := agent.GetUtxosByAddressRPC(from)
		if err != nil {
			return "", "", err
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxInPutsOutPutsRPC(from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}

		fmt.Println("trxSigHex:", trxSigHex)

		// broadcast transaction
		trxId, err := agent.BroadcastTransactionRPC(trxSigHex)
		if err != nil {
			return "", "", err
		}

		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxId, nil
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateTransaction(coinSymbol,from,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil

	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateERC20Transaction(coinSymbol,from,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get utxos
		utxos, err := agent.GetUtxosByAddressRPC(from)
		if err != nil {
			return "", "", err
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxInPutsOutPutsRPC(from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}
		rawTrx, err = agent.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, uint32(coinObj.OmniPropertyId), amountStr)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}

		fmt.Println("trxSigHex:", trxSigHex)

		// broadcast transaction
		trxId, err := agent.BroadcastTransactionRPC(trxSigHex)
		if err != nil {
			return "", "", err
		}
		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxId, nil
	}
	return "", "", nil
}

func CreateMultiSignedTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, keyIndex uint16,
	needCount int, pubKeyStrSlice []string, utxos []UTXODetail,
	fromAddr string, to string, amountStr string, feeStr string, gasPrice string, gasLimit string) (string, string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get redeemscript
		redeemScript, err := GetMultiSignRedeemScript(needCount, pubKeyStrSlice, coinSymbol)
		if err != nil {
			return "", "", err
		}

		// get multisig address, check multisig address
		multiSigAddr , err := GetMultiSignAddress(redeemScript, coinSymbol)
		if err != nil {
			return "", "", err
		}
		if multiSigAddr != fromAddr {
			return "", "", errors.New("pubKeyStrSlice is not match to fromAddr")
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxFromUtxosRPC(utxos, fromAddr, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.MultiSignRawTransactionRPC(rawTrx, redeemScript, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}
		fmt.Println("trxSigHex:", trxSigHex)

		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxSigHex, nil
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateTransaction(coinSymbol,fromAddr,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil

	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateERC20Transaction(coinSymbol,fromAddr,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get redeemscript
		redeemScript, err := GetMultiSignRedeemScript(needCount, pubKeyStrSlice, coinSymbol)
		if err != nil {
			return "", "", err
		}

		// get multisig address, check multisig address
		multiSigAddr , err := GetMultiSignAddress(redeemScript, coinSymbol)
		if err != nil {
			return "", "", err
		}
		if multiSigAddr != fromAddr {
			return "", "", errors.New("pubKeyStrSlice is not match to fromAddr")
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxFromUtxosRPC(utxos, fromAddr, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}
		rawTrx, err = agent.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, uint32(coinObj.OmniPropertyId), amountStr)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.MultiSignRawTransactionRPC(rawTrx, redeemScript, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}
		fmt.Println("trxSigHex:", trxSigHex)

		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxSigHex, nil
	}
	return "", "", nil
}

func ListUnSpent(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string)  ([]UTXODetail, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return nil, utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		utxos, err := agent.GetUtxosByAddressRPC(address)
		return utxos, err
	}else if coinSymbol == "ETH"{
		return nil, errors.New("ListUnSpent not support ETH")
	}else if coinObj.IsErc20{
		return nil, errors.New("ListUnSpent not support ERC20")
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		utxos, err := agent.GetUtxosByAddressRPC(address)
		return utxos, err
	}
	return nil, nil
}

// transaction build (support single and multi sign)
func CreateRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, utxos []UTXODetail,
	from string, to string, amountStr string, feeStr string, gasPrice string, gasLimit string) (string, string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxFromUtxosRPC(utxos, from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}
		return feeCost, rawTrx, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxFromUtxosRPC(utxos, from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}
		rawTrx, err = agent.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, uint32(coinObj.OmniPropertyId), amountStr)
		if err != nil {
			return "", "", err
		}
		return feeCost, rawTrx, nil
	}
	return "", "", nil
}

// single sign transaction
func SignRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, utxos []UTXODetail,
	keyIndex uint16, pubKeyStr string, rawTrx string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// sign transaction
		trxSigHex, err := agent.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// sign transaction
		trxSigHex, err := agent.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	}
	return "", nil
}

// multi sign transaction
func MultiSignRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, utxos []UTXODetail,
	keyIndex uint16, redeemScript string, rawTrx string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// sign transaction
		trxSigHex, err := agent.MultiSignRawTransactionRPC(rawTrx, redeemScript, keyIndex, utxos)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// sign transaction
		trxSigHex, err := agent.MultiSignRawTransactionRPC(rawTrx, redeemScript, keyIndex, utxos)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	}
	return "", nil
}

// verify multi sign transaction with one signature
func MultiVerifySignRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, utxos []UTXODetail,
	pubKeyStr string, signedRawTrx string) error {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// verify signature
		err := agent.MultiVerifySignRawTransactionRPC(signedRawTrx, pubKeyStr, nil)
		if err != nil {
			return err
		}
		return nil
	} else if coinSymbol == "BCH" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// verify signature
		err := agent.MultiVerifySignRawTransactionRPC(signedRawTrx, pubKeyStr, utxos)
		if err != nil {
			return err
		}
		return nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// verify signature
		err := agent.MultiVerifySignRawTransactionRPC(signedRawTrx, pubKeyStr, nil)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func CombineRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, signedRawTrxs []string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// combine transaction
		trxSigHex, err := agent.CombineRawTransactionRPC(signedRawTrxs)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// combine transaction
		trxSigHex, err := agent.CombineRawTransactionRPC(signedRawTrxs)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	}
	return "", nil
}

func SendRawTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, signedRawTrx string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" {
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// broadcast transaction
		trxSigHex, err := agent.BroadcastTransactionRPC(signedRawTrx)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	} else if coinSymbol == "ETH"{
		//TODO
	} else if coinObj.IsErc20{
		//TODO
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// broadcast transaction
		trxSigHex, err := agent.BroadcastTransactionRPC(signedRawTrx)
		if err != nil {
			return "", err
		}
		return trxSigHex, nil
	}
	return "", nil
}

// sigType: '1' r+s,  '2' der
func CoinSignTrx(sigType byte, signData []byte, keyIndex uint16) ([]byte, error) {
	conn := new(buffer_tcp.BufferTcpConn)
	err := conn.TCPConnect(config.GlobalConfig.CryptoDeviceConfig.DeviceIp,
		config.GlobalConfig.CryptoDeviceConfig.DevicePort, float64(config.GlobalConfig.CryptoDeviceConfig.TimeOut))
	if err != nil {
		return nil, err
	}
	defer conn.TCPDisConnect()

	var l7req cryptocard.L7Request
	l7req.Set(sigType, keyIndex, nil, signData)
	err = l7req.Pack(conn)
	if err != nil {
		return nil, err
	}

	var l8resp cryptocard.L8Response
	err = l8resp.UnPack(conn)
	if err != nil {
		return nil, err
	}

	return l8resp.DataSigned, nil
}

// sigType: '1' r+s,  '2' der
func CoinVerifyTrx(sigType byte, keyIndex uint16, signData []byte, signedData []byte) (bool, error){
	conn := new(buffer_tcp.BufferTcpConn)
	err := conn.TCPConnect(config.GlobalConfig.CryptoDeviceConfig.DeviceIp,
		config.GlobalConfig.CryptoDeviceConfig.DevicePort, float64(config.GlobalConfig.CryptoDeviceConfig.TimeOut))
	if err != nil {
		return false, err
	}
	defer conn.TCPDisConnect()

	var l4req cryptocard.L4Request
	l4req.Set(sigType, keyIndex, nil, nil, signData, signedData)
	err = l4req.Pack(conn)
	if err != nil {
		return false, err
	}

	var l5resp cryptocard.L5Response
	err = l5resp.UnPack(conn)
	if err != nil {
		return false, err
	}

	if l5resp.ErrCode[0] == '0' && l5resp.ErrCode[1] == '0' {
		return true, nil
	} else {
		return false, nil
	}

	return true, nil
}

// sigType: '1' r+s,  '2' der
//func CoinVerifyTrxWithOutsidePubkey(sigType byte, outsidePubKey []byte, signData []byte, signedData []byte) (bool, error){
//	conn := new(buffer_tcp.BufferTcpConn)
//	err := conn.TCPConnect(config.GlobalConfig.CryptoDeviceConfig.DeviceIp,
//		config.GlobalConfig.CryptoDeviceConfig.DevicePort, float64(config.GlobalConfig.CryptoDeviceConfig.TimeOut))
//	if err != nil {
//		return false, err
//	}
//	defer conn.TCPDisConnect()
//
//	var l4req cryptocard.L4Request
//	l4req.Set(sigType, 9999, []byte("T1"), outsidePubKey, signData, signedData)
//	err = l4req.Pack(conn)
//	if err != nil {
//		return false, err
//	}
//
//	var l5resp cryptocard.L5Response
//	err = l5resp.UnPack(conn)
//	if err != nil {
//		return false, err
//	}
//
//	if l5resp.ErrCode[0] == '0' && l5resp.ErrCode[1] == '0' {
//		return true, nil
//	} else {
//		return false, nil
//	}
//
//	return true, nil
//}

// use btcec module
func CoinVerifyTrx2(pubKeyBytes []byte, signData []byte, signedData []byte) (bool, error){
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}
	signature, err := btcec.ParseSignature(signedData, btcec.S256())
	if err != nil {
		return false, err
	}
	verified := signature.Verify(signData, pubKey)
	if !verified {
		return false, nil
	}
	return true, nil
}

func SerializeDerEncoding(rBytes []byte, sBytes []byte) ([]byte, error) {
	if len(rBytes) != 32 {
		return nil, errors.New("invalid rBytes len")
	}
	if len(sBytes) != 32 {
		return nil, errors.New("invalid sBytes len")
	}

	var r []byte
	r = append(r, 0)
	r = append(r, rBytes...)
	var s []byte
	s = append(s, 0)
	s = append(s, sBytes...)
	for {
		if len(r) > 1 && r[0] == 0 && r[1] < 0x80 {
			r = r[1:]
		} else {
			break
		}
	}
	for {
		if len(s) > 1 && s[0] == 0 && s[1] < 0x80 {
			s = s[1:]
		} else {
			break
		}
	}

	size := 6 + len(r) + len(s)
	signedData := make([]byte, size, size)
	signedData[0] = 0x30
	signedData[1] = 4 + byte(len(r)) + byte(len(s))
	signedData[2] = 0x2
	signedData[3] = byte(len(r))
	copy(signedData[4:4+len(r)], r)
	signedData[4+len(r)] = 0x2
	signedData[5+len(r)] = byte(len(s))
	copy(signedData[6+len(r):6+len(r)+len(s)], s)

	return signedData, nil
}

func CalcBalanceFromUtxos(utxos []UTXODetail, address string, nPrec int) string {
	sum := int64(0)
	for _, utxo := range utxos {
		if utxo.Address != address {
			continue
		}
		sum += utxo.Amount
	}
	return FromPrecisionAmount(sum, nPrec)
}
