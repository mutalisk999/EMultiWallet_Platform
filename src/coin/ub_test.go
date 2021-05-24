package coin

import (
	"config"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestUBAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	ba, err := ag.GetBalanceByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestUBAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	c, err := ag.IsTransactionConfirmedRPC("02faec58bfbcb47a34b4670b25fe2117d69e8aa234eeee2d62b262a75fd63d65")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestUBAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	c, err := ag.IsAddressValidRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestUBAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestUBCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	config.IsTestEnvironment = true
	addrStr, _ := UBCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestUBGetRedeemScriptByPubKeys(t *testing.T) {
	needCount := 2
	// keyindex 65,66,67
	pubKey1Str := "24fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b719917e4a60cfa23e02f4a0963c4b6353b6e1b6b2be28894115a57690cbcb7a410"
	pubKey2Str := "60b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb60944a966aab0bd60187862950a1991377106bf35c8ed6d726adb5db5bb08dc6a8a76"
	pubKey3Str := "9cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d0096fa52c237baa36e689d79949c20d0259737234fa8f049cbef55c7d0e57f958ce"
	pubKeyStrList := []string{pubKey1Str, pubKey2Str, pubKey3Str}
	redeemScript, _ := UBGetRedeemScriptByPubKeys(needCount, pubKeyStrList)
	fmt.Println(redeemScript)
}

func TestUBGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953ae"
	config.IsTestEnvironment = true
	addrStr, _ := UBGetMultiSignAddressByRedeemScript(redeemScriptStr)
	fmt.Println(addrStr)
}

func TestUBGetMultiSignScriptByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953ae"
	p2shScript, _ := UBGetMultiSignScriptByRedeemScript(redeemScriptStr)
	fmt.Println(hex.EncodeToString(p2shScript))
}

func TestUBAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestUBAgent_BuildTrxFromUtxosRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	feeCost, inPuts, outPuts, err := ag.BuildTrxFromUtxosRPC(utxos, "mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestUBAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestUBUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000017d704bb25c19e317c2db09b95d364a57ecd4a81b3584f647dd7e4dbf0599e9060100000000ffffffff0240420f00000000001976a914130c91e9def87a4445440ab430f7a98e76424d9188ac98472677000000001976a9141b46aa5c903f3dc8eb9592876af0a061db4b3bed88ac00000000"
	trx, _ := UBUnPackRawTransaction(rawTrx)
	fmt.Println("trx version:", trx.Version)
	fmt.Println("trx locktime", trx.LockTime)
	fmt.Println("trx vin size:", len(trx.Vin))
	for i := 0; i < len(trx.Vin); i++ {
		fmt.Println("vin prevout:", trx.Vin[i].PrevOut.Hash.GetHex(), trx.Vin[i].PrevOut.N)
		fmt.Println("vin scriptsig:", trx.Vin[i].ScriptSig)
		fmt.Println("vin sequence:", trx.Vin[i].Sequence)
		fmt.Println("vin scriptwitness:", trx.Vin[i].ScriptWitness)
	}
}

func TestUBSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	keyIndex := uint16(6)

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	trxSigHex, err := ag.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestUBAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	keyIndex := uint16(6)
	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	trxSigHex, err := ag.SignRawTransactionRPC(rawTrx, pubKeyStr, keyIndex, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSigHex:", trxSigHex)
	trxId, err := ag.BroadcastTransactionRPC(trxSigHex)
	if err == nil {
		fmt.Println("trxId:", trxId)
	} else {
		fmt.Println("err:", err)
	}
}

func TestUBAgent_CreateRawTransaction2(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N4XK3G52pnVK3BV5QoApN4jgN8jPCv6p1N", "mg95NR86A5qWbshDB97pLCkcPyiV8XUMpX", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestUBAgent_MultiSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N4XK3G52pnVK3BV5QoApN4jgN8jPCv6p1N", "mg95NR86A5qWbshDB97pLCkcPyiV8XUMpX", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(65)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig1Hex:", trxSig1Hex)

	keyIndex2 := uint16(66)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig2Hex:", trxSig2Hex)
}

func TestUBAgent_MultiVerifySignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	trxSig1Hex := "020000000128d0bd6fa2162103fed9ba3a28f6290b480e1b956573734b2ad29636e8426e0801000000b40047304402204594e5252c29575baec301d27dbf1ea7f8dab92e27ee285b75a7281fb8b300f702206088e862862cb24c5a622f2ac5fdb7448dea71bb0ee29cafbd05ed5eeeb6b1dd094c6952210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953aeffffffff02184a89000000000017a9147bb44171b973baf49ae70a401af4be888b852e028740420f00000000001976a91406d4e176d8fa2028a0f50c3ce3863b55def4ea3788ac00000000"
	pubKey1 := "24fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b719917e4a60cfa23e02f4a0963c4b6353b6e1b6b2be28894115a57690cbcb7a410"
	err := ag.MultiVerifySignRawTransactionRPC(trxSig1Hex, pubKey1, nil)
	if err == nil {
		fmt.Println("Verify trxSig1Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig1Hex fail")
	}

	trxSig2Hex := "020000000128d0bd6fa2162103fed9ba3a28f6290b480e1b956573734b2ad29636e8426e0801000000b500483045022100ef29820292c19a5cb2e827f1e44f50c4f3943c6dbacfeddeca82b54cb038adbb02204c8b08874ac20449c32cdf688001876ee43a43130ca85bc2c4017c415d976332094c6952210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953aeffffffff02184a89000000000017a9147bb44171b973baf49ae70a401af4be888b852e028740420f00000000001976a91406d4e176d8fa2028a0f50c3ce3863b55def4ea3788ac00000000"
	pubKey2 := "60b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb60944a966aab0bd60187862950a1991377106bf35c8ed6d726adb5db5bb08dc6a8a76"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig2Hex, pubKey2, nil)
	if err == nil {
		fmt.Println("Verify trxSig2Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig2Hex fail")
	}
}

func TestUBAgent_CombineRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	trxSig1Hex := "020000000128d0bd6fa2162103fed9ba3a28f6290b480e1b956573734b2ad29636e8426e0801000000b40047304402204594e5252c29575baec301d27dbf1ea7f8dab92e27ee285b75a7281fb8b300f702206088e862862cb24c5a622f2ac5fdb7448dea71bb0ee29cafbd05ed5eeeb6b1dd094c6952210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953aeffffffff02184a89000000000017a9147bb44171b973baf49ae70a401af4be888b852e028740420f00000000001976a91406d4e176d8fa2028a0f50c3ce3863b55def4ea3788ac00000000"
	trxSig2Hex := "020000000128d0bd6fa2162103fed9ba3a28f6290b480e1b956573734b2ad29636e8426e0801000000b500483045022100ef29820292c19a5cb2e827f1e44f50c4f3943c6dbacfeddeca82b54cb038adbb02204c8b08874ac20449c32cdf688001876ee43a43130ca85bc2c4017c415d976332094c6952210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953aeffffffff02184a89000000000017a9147bb44171b973baf49ae70a401af4be888b852e028740420f00000000001976a91406d4e176d8fa2028a0f50c3ce3863b55def4ea3788ac00000000"

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestUBAgent_BroadcastTransactionRPC2(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")

	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N4XK3G52pnVK3BV5QoApN4jgN8jPCv6p1N", "mg95NR86A5qWbshDB97pLCkcPyiV8XUMpX", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210224fc56b357ebd4c87e2e5825c376d1c30fe030f71407eacd45dac55f00241b71210260b48a1d709dcf63125bc5d5c70045b75d30572ff0e451069b67256efbb6094421029cff24f5619badc6c382a9cbebade5762ff6c21a3944aa2d4443e15d8a57d00953ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(65)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex2 := uint16(66)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})

	fmt.Println("trxSigHex:", trxSigHex)
	trxId, err := ag.BroadcastTransactionRPC(trxSigHex)
	if err == nil {
		fmt.Println("trxId:", trxId)
	} else {
		fmt.Println("err:", err)
	}
}
