package coin

import (
	"config"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestLTCAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	ba, err := ag.GetBalanceByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestLTCAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	c, err := ag.IsTransactionConfirmedRPC("9d5675ffe98872f1fd02b548032de6e40a48259430ae8c00f185df883ff18d81")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestLTCAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	c, err := ag.IsAddressValidRPC("n1uaazxSzWochahoAnKbGPAxh34MhxNe1J")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestLTCAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestLTCCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	config.IsTestEnvironment = true
	addrStr, _ := LTCCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestLTCGetRedeemScriptByPubKeys(t *testing.T) {
	needCount := 2
	// keyindex 55,56,57
	pubKey1Str := "20c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802db5269fc144b5ed59f3a345b10f01c4c73effbed6171ec0742e6d53b87e89df85"
	pubKey2Str := "df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff54076b584c2c56f9666ac9dcfd54f7dc6a7f60b4302448ced00e29c6bad371ff"
	pubKey3Str := "7472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d768c7556d72bf9aef5bacc20c501733b08e4a4af20c4eda7a9f615cbd89e1c1b0e"
	pubKeyStrList := []string{pubKey1Str, pubKey2Str, pubKey3Str}
	redeemScript, _ := LTCGetRedeemScriptByPubKeys(needCount, pubKeyStrList)
	fmt.Println(redeemScript)
}

func TestLTCGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653ae"
	config.IsTestEnvironment = true
	addrStr, _ := LTCGetMultiSignAddressByRedeemScript(redeemScriptStr)
	fmt.Println(addrStr)
}

func TestLTCGetMultiSignScriptByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653ae"
	p2shScript, _ := LTCGetMultiSignScriptByRedeemScript(redeemScriptStr)
	fmt.Println(hex.EncodeToString(p2shScript))
}

func TestLTCAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestLTCAgent_BuildTrxFromUtxosRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	feeCost, inPuts, outPuts, err := ag.BuildTrxFromUtxosRPC(utxos, "mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestLTCAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestLTCUnPackRawTransaction(t *testing.T) {
	rawTrx := "0200000001dd31ca8a5c4417867072f7c795e482b95677e942d1b81db00e50fe60ac19a8de0100000000ffffffff0240420f00000000001976a91404b80434d66065d51639e05ac0af1b885058b3de88ac30941201000000001976a914ab9ffaedde51a4b4720ee51208b568d1834b68ef88ac00000000"
	trx, _ := LTCUnPackRawTransaction(rawTrx)
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

func TestLTCSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	keyIndex := uint16(3)
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

func TestLTCAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	keyIndex := uint16(3)

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

func TestLTCAgent_CreateRawTransaction2(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("QVaVKPfpfufUXandNxwvekTwRUipC5xQMv", "Qdj5TPifwxwSx4awSbGw6px1GhRNZUGp9b", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestLTCAgent_MultiSignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("QVaVKPfpfufUXandNxwvekTwRUipC5xQMv", "Qdj5TPifwxwSx4awSbGw6px1GhRNZUGp9b", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(55)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig1Hex:", trxSig1Hex)

	keyIndex2 := uint16(56)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig2Hex:", trxSig2Hex)
}

func TestLTCAgent_MultiVerifySignRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	trxSig1Hex := "0200000001242f703fa5bd6ef0183e7045bdb77c5ee778e3381fd9728f121037a3f071fd7401000000b500483045022100da315d993794d0fe406f84dc2b6a1ee33012d59cbfbefd30e4783d106625a5ed02204c89226dbb34e5903be0fea0ae42db64fb771fb474fde3d3b4864f8288cc165b014c6952210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653aeffffffff02184a89000000000017a9146270aea1b8db2faf44733fa1844c0a0c45fe5de78740420f000000000017a914bbd18356b2777fd198db0f8cdca3dfc2fd6097f18700000000"
	pubKey1 := "20c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802db5269fc144b5ed59f3a345b10f01c4c73effbed6171ec0742e6d53b87e89df85"
	err := ag.MultiVerifySignRawTransactionRPC(trxSig1Hex, pubKey1, nil)
	if err == nil {
		fmt.Println("Verify trxSig1Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig1Hex fail")
	}

	trxSig2Hex := "0200000001242f703fa5bd6ef0183e7045bdb77c5ee778e3381fd9728f121037a3f071fd7401000000b40047304402203d5d1853f03f57ba1aa4a0cb2cc31970b6ce85315ce13df6436b2b7833e149a902205921d84ebfddfb91fe5d33bea028da2eebf2358953001e743f0e1fa0cc98d977014c6952210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653aeffffffff02184a89000000000017a9146270aea1b8db2faf44733fa1844c0a0c45fe5de78740420f000000000017a914bbd18356b2777fd198db0f8cdca3dfc2fd6097f18700000000"
	pubKey2 := "df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff54076b584c2c56f9666ac9dcfd54f7dc6a7f60b4302448ced00e29c6bad371ff"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig2Hex, pubKey2, nil)
	if err == nil {
		fmt.Println("Verify trxSig2Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig2Hex fail")
	}
}

func TestLTCAgent_CombineRawTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	trxSig1Hex := "0200000001242f703fa5bd6ef0183e7045bdb77c5ee778e3381fd9728f121037a3f071fd7401000000b500483045022100da315d993794d0fe406f84dc2b6a1ee33012d59cbfbefd30e4783d106625a5ed02204c89226dbb34e5903be0fea0ae42db64fb771fb474fde3d3b4864f8288cc165b014c6952210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653aeffffffff02184a89000000000017a9146270aea1b8db2faf44733fa1844c0a0c45fe5de78740420f000000000017a914bbd18356b2777fd198db0f8cdca3dfc2fd6097f18700000000"
	trxSig2Hex := "0200000001242f703fa5bd6ef0183e7045bdb77c5ee778e3381fd9728f121037a3f071fd7401000000b40047304402203d5d1853f03f57ba1aa4a0cb2cc31970b6ce85315ce13df6436b2b7833e149a902205921d84ebfddfb91fe5d33bea028da2eebf2358953001e743f0e1fa0cc98d977014c6952210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653aeffffffff02184a89000000000017a9146270aea1b8db2faf44733fa1844c0a0c45fe5de78740420f000000000017a914bbd18356b2777fd198db0f8cdca3dfc2fd6097f18700000000"

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestLTCAgent_BroadcastTransactionRPC2(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")

	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("QVaVKPfpfufUXandNxwvekTwRUipC5xQMv", "Qdj5TPifwxwSx4awSbGw6px1GhRNZUGp9b", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	redeemScriptStr := "52210320c6217a25195a849bd838c9be3b62f50ef95ad882e47b43f1b7443b8831802d2103df58fcfa18f18d9e37f37d7506fcc89257b959574620ced178d117fc2ae3abff21027472bd877971f37eb46a5bbd801b11d232d633fc461e1f0dc18c2e8f0c137d7653ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(55)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex2 := uint16(56)
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
