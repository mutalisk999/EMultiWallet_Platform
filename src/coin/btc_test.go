package coin

import (
	"config"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestBTCAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	ba, err := ag.GetBalanceByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestBTCAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	c, err := ag.IsTransactionConfirmedRPC("ed410a120a23c9c8e078f2ed8f43aa99f05bd95090e61befee0de00644e578ee")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBTCAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	c, err := ag.IsAddressValidRPC("mhuXoAkUNLPcboTFu9PDtGapc3wnZAGeyw")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBTCAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestBTCCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	config.IsTestEnvironment = true
	addrStr, _ := BTCCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestBTCGetRedeemScriptByPubKeys(t *testing.T) {
	needCount := 2
	// keyindex 50,51,52
	pubKey1Str := "1bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96030d8ccad50d875a84daefc0b03856a4ce10b571d8609631f3378fbe0daa251ae"
	pubKey2Str := "2e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a645b7b8ba6bc5fdcabeee72a763034dcb4389a20bd90cc94032e9dc8334f12a8e0"
	pubKey3Str := "fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1771552ddd361a4cee7f8fe976a1026eb3cfbaf71892a669b050d9416eed196ef5"
	pubKeyStrList := []string{pubKey1Str, pubKey2Str, pubKey3Str}
	redeemScript, _ := BTCGetRedeemScriptByPubKeys(needCount, pubKeyStrList)
	fmt.Println(redeemScript)
}

func TestBTCGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScriptStr := "5221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753ae"
	config.IsTestEnvironment = true
	addrStr, _ := BTCGetMultiSignAddressByRedeemScript(redeemScriptStr)
	fmt.Println(addrStr)
}

func TestBTCGetMultiSignScriptByRedeemScript(t *testing.T) {
	redeemScriptStr := "5221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753ae"
	p2shScript, _ := BTCGetMultiSignScriptByRedeemScript(redeemScriptStr)
	fmt.Println(hex.EncodeToString(p2shScript))
}

func TestBTCAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBTCAgent_BuildTrxFromUtxosRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	feeCost, inPuts, outPuts, err := ag.BuildTrxFromUtxosRPC(utxos, "mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBTCAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBTCUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000016059331a0741e2be764f25f7c85fa78855aa63cdc780ca10bfcac168ad7fd7ad0000000000ffffffff0240420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a2887184a8900000000001976a91430e83ac38de345ce80f069055cb9f8e15b28e54d88ac00000000"
	trx, _ := BTCUnPackRawTransaction(rawTrx)
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

func TestBTCSignRawTransaction(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	keyIndex := uint16(1)

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

func TestBTCAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	keyIndex := uint16(1)
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

func TestBTCAgent_CreateRawTransaction2(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2MzetpTXo7sH1UntpUErQG22n451bFTXBuS", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBTCAgent_MultiSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2MzetpTXo7sH1UntpUErQG22n451bFTXBuS", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "5221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(50)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig1Hex:", trxSig1Hex)

	keyIndex2 := uint16(51)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig2Hex:", trxSig2Hex)
}

func TestBTCAgent_MultiVerifySignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	trxSig1Hex := "02000000016b967d9e18a642ca9cd99e7a5fce3da7374d567617015a5dc5f3aaa7964a3a8900000000b40047304402207feecf93f29f978f4a96fd411efef0e9df7f1b90fabfae0a4a5b6b6bb71996d90220689d08cff587fb3848cffc727d799ed107fed7741262faf3ec396b4bb897f9e2014c695221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753aeffffffff0278184c000000000017a9145142dc9bd57a2d1dfb26571df244c903dde7b7ad8740420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a288700000000"
	pubKey1 := "1bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96030d8ccad50d875a84daefc0b03856a4ce10b571d8609631f3378fbe0daa251ae"
	err := ag.MultiVerifySignRawTransactionRPC(trxSig1Hex, pubKey1, nil)
	if err == nil {
		fmt.Println("Verify trxSig1Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig1Hex fail")
	}

	trxSig2Hex := "02000000016b967d9e18a642ca9cd99e7a5fce3da7374d567617015a5dc5f3aaa7964a3a8900000000b40047304402204b544441688c5dd307cced704ba6107ca1f0f122f152f4695c2918f8828a6ae602201673b61f58063066597243e3438024f7a2fc2f2428a7190be898b170398d18fb014c695221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753aeffffffff0278184c000000000017a9145142dc9bd57a2d1dfb26571df244c903dde7b7ad8740420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a288700000000"
	pubKey2 := "2e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a645b7b8ba6bc5fdcabeee72a763034dcb4389a20bd90cc94032e9dc8334f12a8e0"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig2Hex, pubKey2, nil)
	if err == nil {
		fmt.Println("Verify trxSig2Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig2Hex fail")
	}
}

func TestBTCAgent_CombineRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	trxSig1Hex := "0200000001c83985d23cf5b7e967a051e87a43562cc807328664c970ffff33f018d9a52a0b00000000b40047304402202d79a4868148d3398fd76b9ba44b7ffd872e03dd8393a23fc2f6ce7fe113804f0220284a9b7349f3ebb261e30402bf3329f24c9fa9fd85f9734beb05fd0a9035e3ef014c695221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753aeffffffff0248b16a000000000017a9145142dc9bd57a2d1dfb26571df244c903dde7b7ad8740420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a288700000000"
	trxSig2Hex := "0200000001c83985d23cf5b7e967a051e87a43562cc807328664c970ffff33f018d9a52a0b00000000b400473044022025adfc97448c1b88aa2c692a50074cf114deeddd295d90395e5d8844bdd31d55022064b0c935bf9f7e13f221084544f03b116362205a8c9978623a63b90c582c96b1014c695221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753aeffffffff0248b16a000000000017a9145142dc9bd57a2d1dfb26571df244c903dde7b7ad8740420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a288700000000"

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestBTCAgent_BroadcastTransactionRPC2(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")

	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2MzetpTXo7sH1UntpUErQG22n451bFTXBuS", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "5221021bea5be697868ee1edb75e22df72aca99ed0de01e141f41996e50fec89dad96021022e4a01654e525a658b994d20c5af87f933bdbd5b718fcdc5bc7f0ea2b7855a642103fce46b4ce30de24d582c174009292092520b1e57a624c56f34447ead4c94ff1753ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(50)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex2 := uint16(51)
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
