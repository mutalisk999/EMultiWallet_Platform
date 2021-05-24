package coin

import (
	"config"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestBCHAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	ba, err := ag.GetBalanceByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestBCHAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	c, err := ag.IsTransactionConfirmedRPC("4a6f14efecc0f10a8afff7dd1a9d1c2b69e2c1588ab512fd85b023f0c19fcc1f")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBCHAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	c, err := ag.IsAddressValidRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBCHAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestBCHCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	config.IsTestEnvironment = true
	addrStr, _ := BCHCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestBCHGetRedeemScriptByPubKeys(t *testing.T) {
	needCount := 2
	// keyindex 70,71,72
	pubKey1Str := "21717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13fb20b528475074cf4598d0da500a921d3b79ce807d3cc3fc03ec54cc521d003e0"
	pubKey2Str := "d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000ecfe0c6415df571093ac9cd7caa02128e0df3fb2a2818a5b6815c6b354bc22d3"
	pubKey3Str := "36ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b28d25689b2172ce873a38d59b4554f77d660daebea4e280a58a214f15fb81b07"
	pubKeyStrList := []string{pubKey1Str, pubKey2Str, pubKey3Str}
	redeemScript, _ := BCHGetRedeemScriptByPubKeys(needCount, pubKeyStrList)
	fmt.Println(redeemScript)
}

func TestBCHGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53ae"
	config.IsTestEnvironment = true
	addrStr, _ := BCHGetMultiSignAddressByRedeemScript(redeemScriptStr)
	fmt.Println(addrStr)
}

func TestBCHGetMultiSignScriptByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53ae"
	p2shScript, _ := BCHGetMultiSignScriptByRedeemScript(redeemScriptStr)
	fmt.Println(hex.EncodeToString(p2shScript))
}

func TestBCHAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBCHAgent_BuildTrxFromUtxosRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	feeCost, inPuts, outPuts, err := ag.BuildTrxFromUtxosRPC(utxos, "bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBCHAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBCHUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000011fcc9fc1f023b085fd12b58a58c1e2692b1c9d1addf7ff8a0af1c0ecef146f4a0100000000ffffffff0240420f00000000001976a9146180cf77dca5a0a5f761e9cef0b1a5c917384e8088ac18a4eb02000000001976a9149e6827242840668d9bb217498c3a29207f9ed44888ac00000000"
	trx, _ := BCHUnPackRawTransaction(rawTrx)
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

func TestBCHSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	keyIndex := uint16(4)
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

func TestBCHAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	keyIndex := uint16(4)

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

func TestBCHAgent_CreateRawTransaction2(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu", "bchtest:qpqdwvcvrg5yntyn03rn9wlv0527hl9ct5qm37qzff", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBCHAgent_MultiSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu", "bchtest:qpqdwvcvrg5yntyn03rn9wlv0527hl9ct5qm37qzff", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	utxos, err := ag.GetUtxosByAddressRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex1 := uint16(70)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig1Hex:", trxSig1Hex)

	keyIndex2 := uint16(71)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig2Hex:", trxSig2Hex)
}

func TestBCHAgent_MultiVerifySignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	utxos, err := ag.GetUtxosByAddressRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	trxSig1Hex := "02000000015e45017e55aa6f5eee7827ddd0be037cd2d81006a262fba623e5750a2ee5065300000000b40047304402202a77b9393d33fae540ca1c5f2eaa2fe7b07d710af033926bb43d8c81a6fc47db0220303354af1bc045f4a0939436c32a50c78618608b628038b52554c9c814da7094414c6952210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53aeffffffff0248b16a000000000017a9148b9197e344178c1fe0dc826fc0ff6baca7b682bb8740420f00000000001976a91440d7330c1a2849ac937c4732bbec7d15ebfcb85d88ac00000000"
	pubKey1 := "21717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13fb20b528475074cf4598d0da500a921d3b79ce807d3cc3fc03ec54cc521d003e0"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig1Hex, pubKey1, utxos)
	if err == nil {
		fmt.Println("Verify trxSig1Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig1Hex fail")
	}

	trxSig2Hex := "02000000015e45017e55aa6f5eee7827ddd0be037cd2d81006a262fba623e5750a2ee5065300000000b400473044022015346be776e3d4d1727c8b4f0fa3a0c6480baf8bf144c69ec0c48fb4dea5b860022003e6933afc6f894f6605404c5f891ad09b348e87620664a8c2388960907ac186414c6952210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53aeffffffff0248b16a000000000017a9148b9197e344178c1fe0dc826fc0ff6baca7b682bb8740420f00000000001976a91440d7330c1a2849ac937c4732bbec7d15ebfcb85d88ac00000000"
	pubKey2 := "d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000ecfe0c6415df571093ac9cd7caa02128e0df3fb2a2818a5b6815c6b354bc22d3"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig2Hex, pubKey2, utxos)
	if err == nil {
		fmt.Println("Verify trxSig2Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig2Hex fail")
	}
}

func TestBCHAgent_CombineRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	trxSig1Hex := "02000000015e45017e55aa6f5eee7827ddd0be037cd2d81006a262fba623e5750a2ee5065300000000b40047304402202a77b9393d33fae540ca1c5f2eaa2fe7b07d710af033926bb43d8c81a6fc47db0220303354af1bc045f4a0939436c32a50c78618608b628038b52554c9c814da7094414c6952210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53aeffffffff0248b16a000000000017a9148b9197e344178c1fe0dc826fc0ff6baca7b682bb8740420f00000000001976a91440d7330c1a2849ac937c4732bbec7d15ebfcb85d88ac00000000"
	trxSig2Hex := "02000000015e45017e55aa6f5eee7827ddd0be037cd2d81006a262fba623e5750a2ee5065300000000b400473044022015346be776e3d4d1727c8b4f0fa3a0c6480baf8bf144c69ec0c48fb4dea5b860022003e6933afc6f894f6605404c5f891ad09b348e87620664a8c2388960907ac186414c6952210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53aeffffffff0248b16a000000000017a9148b9197e344178c1fe0dc826fc0ff6baca7b682bb8740420f00000000001976a91440d7330c1a2849ac937c4732bbec7d15ebfcb85d88ac00000000"

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestBCHAgent_BroadcastTransactionRPC2(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")

	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu", "bchtest:qpqdwvcvrg5yntyn03rn9wlv0527hl9ct5qm37qzff", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210221717d130011fd20b5c6dcec17c02bb567b89f1c337e8a46e63533af5656a13f2103d262f0c9af4bb714cd3f096c0e2dac54ad640fb459f785c27d6b0566172c0000210336ee3fe30f11d3b96274a6592f83e03a3469eb940817c18205728a74b1a9121b53ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	utxos, err := ag.GetUtxosByAddressRPC("bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex1 := uint16(70)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex2 := uint16(71)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, utxos)
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
