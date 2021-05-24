package coin

import (
	"config"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestOMNIAgent_GetFeeBalanceByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	ba, err := ag.GetFeeBalanceByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestOMNIAgent_GetBalanceByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	ba, err := ag.GetBalanceByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestOMNIAgent_IsFeeTransactionConfirmed(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsFeeTransactionConfirmedRPC("08d75f23070a12bb1059e3498dab47a000ce303238f99584da412b6a4c6b9558")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_IsTransactionConfirmed(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsTransactionConfirmedRPC("362ec25c3f877f016918babca120cf80bfbc05051091e7079bb16b8ddf5807cf")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_IsAddressValidRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsAddressValidRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_GetUtxosByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestOMNICalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	config.IsTestEnvironment = true
	addrStr, _ := OMNICalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestOMNIGetRedeemScriptByPubKeys(t *testing.T) {
	needCount := 2
	// keyindex 60,61,62
	pubKey1Str := "04947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c14b0d63474bd312a29c94af32ebcea6a5a11870e8f502aa7719ec7a63bf8fa8b4d"
	pubKey2Str := "1d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c3430b5a987aa8eeca6f0fe1f973a65e0d4e2fcce2b8ed1386541e11d9ca3f20e29"
	pubKey3Str := "155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c268d84ec82a0a8a892e0a6de42bf68305f72412e5facfdb0c44167d9beef277b9f"
	pubKeyStrList := []string{pubKey1Str, pubKey2Str, pubKey3Str}
	redeemScript, _ := OMNIGetRedeemScriptByPubKeys(needCount, pubKeyStrList)
	fmt.Println(redeemScript)
}

func TestOMNIGetMultiSignAddressByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653ae"
	config.IsTestEnvironment = true
	addrStr, _ := OMNIGetMultiSignAddressByRedeemScript(redeemScriptStr)
	fmt.Println(addrStr)
}

func TestOMNIGetMultiSignScriptByRedeemScript(t *testing.T) {
	redeemScriptStr := "52210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653ae"
	p2shScript, _ := OMNIGetMultiSignScriptByRedeemScript(redeemScriptStr)
	fmt.Println(hex.EncodeToString(p2shScript))
}

func TestOMNIAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestOMNIAgent_BuildTrxFromUtxosRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	feeCost, inPuts, outPuts, _ := ag.BuildTrxFromUtxosRPC(utxos, "mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestOMNIAgent_CreateRawTransaction(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestOMNIAgent_CreateRawTransactionOpReturnRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	rawTrx := "0100000001f1238a114f4a2ca83a5dbbc04832dcf0e8d8e4c02eeaecab22f03f24bb9be5e50100000000ffffffff02a6889800000000001976a9142996789a3cdd905d8e44fb6033a08c667410ef6e88ac22020000000000001976a9149f2a66d7b349ba9a87e6bf9cf7da1df31697d63988ac00000000"
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	fmt.Println(rawTrx)
}

func TestOMNIUnPackRawTransaction(t *testing.T) {
	rawTrx := "0100000001f1238a114f4a2ca83a5dbbc04832dcf0e8d8e4c02eeaecab22f03f24bb9be5e50100000000ffffffff03a6889800000000001976a9142996789a3cdd905d8e44fb6033a08c667410ef6e88ac22020000000000001976a9149f2a66d7b349ba9a87e6bf9cf7da1df31697d63988ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"
	trx, _ := OMNIUnPackRawTransaction(rawTrx)
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

func TestOMNISignRawTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	keyIndex := uint16(5)

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

func TestOMNIAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	keyIndex := uint16(5)
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

func TestOMNIAgent_CreateRawTransaction2(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N6bnijfVT1FSUfDV65foDXrarbpThr7D5V", "n3K9Zmw3TKjfvmKVfT4Td2AgqTzkyb3vtz", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)

	fmt.Println(rawTrx)
}

func TestOMNIAgent_CreateRawTransactionOpReturn2(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	rawTrx := "0100000001acee2c0357fe21d668a79129072da073dd8c77b48e8d2dedd324eb5f8abf0f630000000000ffffffff02a68898000000000017a914927d471577d016ba43082d2f1e244b45640cf6478722020000000000001976a914ef172f2edb561b373a2677e8d97fbd4bf46b1dbd88ac00000000"
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")

	fmt.Println(rawTrx)
}

func TestOMNIAgent_MultiSignRawTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N6bnijfVT1FSUfDV65foDXrarbpThr7D5V", "n3K9Zmw3TKjfvmKVfT4Td2AgqTzkyb3vtz", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	redeemScriptStr := "52210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(60)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig1Hex:", trxSig1Hex)

	keyIndex2 := uint16(61)
	trxSig2Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex2, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSig2Hex:", trxSig2Hex)
}

func TestOMNIAgent_MultiVerifySignRawTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	trxSig1Hex := "0100000001acee2c0357fe21d668a79129072da073dd8c77b48e8d2dedd324eb5f8abf0f6300000000b40047304402200d04ac16f272858d4b3edf765a3def04557b3c14364fe19690fb6daed6e4231202206737f26d1da6b1d4656769b385e5f6988e0cbf8afff56661003164aec1dc84e4014c6952210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653aeffffffff03a68898000000000017a914927d471577d016ba43082d2f1e244b45640cf6478722020000000000001976a914ef172f2edb561b373a2677e8d97fbd4bf46b1dbd88ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"
	pubKey1 := "04947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c14b0d63474bd312a29c94af32ebcea6a5a11870e8f502aa7719ec7a63bf8fa8b4d"
	err := ag.MultiVerifySignRawTransactionRPC(trxSig1Hex, pubKey1, nil)
	if err == nil {
		fmt.Println("Verify trxSig1Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig1Hex fail")
	}

	trxSig2Hex := "0100000001acee2c0357fe21d668a79129072da073dd8c77b48e8d2dedd324eb5f8abf0f6300000000b400473044022042ded662a37181a0f29f260ec5c79ea55370edfaab1572e88e43fe3d60fa830702200c745992894efc54ad4c1e1dad772e7011fa459914671ae2bfa177decf3d26e6014c6952210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653aeffffffff03a68898000000000017a914927d471577d016ba43082d2f1e244b45640cf6478722020000000000001976a914ef172f2edb561b373a2677e8d97fbd4bf46b1dbd88ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"
	pubKey2 := "1d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c3430b5a987aa8eeca6f0fe1f973a65e0d4e2fcce2b8ed1386541e11d9ca3f20e29"
	err = ag.MultiVerifySignRawTransactionRPC(trxSig2Hex, pubKey2, nil)
	if err == nil {
		fmt.Println("Verify trxSig2Hex success")
	} else {
		fmt.Println(err)
		fmt.Println("Verify trxSig2Hex fail")
	}
}

func TestOMNIAgent_CombineRawTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	trxSig1Hex := "0100000001acee2c0357fe21d668a79129072da073dd8c77b48e8d2dedd324eb5f8abf0f6300000000b40047304402200d04ac16f272858d4b3edf765a3def04557b3c14364fe19690fb6daed6e4231202206737f26d1da6b1d4656769b385e5f6988e0cbf8afff56661003164aec1dc84e4014c6952210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653aeffffffff03a68898000000000017a914927d471577d016ba43082d2f1e244b45640cf6478722020000000000001976a914ef172f2edb561b373a2677e8d97fbd4bf46b1dbd88ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"
	trxSig2Hex := "0100000001acee2c0357fe21d668a79129072da073dd8c77b48e8d2dedd324eb5f8abf0f6300000000b400473044022042ded662a37181a0f29f260ec5c79ea55370edfaab1572e88e43fe3d60fa830702200c745992894efc54ad4c1e1dad772e7011fa459914671ae2bfa177decf3d26e6014c6952210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653aeffffffff03a68898000000000017a914927d471577d016ba43082d2f1e244b45640cf6478722020000000000001976a914ef172f2edb561b373a2677e8d97fbd4bf46b1dbd88ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"

	trxSigHex, _ := ag.CombineRawTransactionRPC([]string{trxSig1Hex, trxSig2Hex})
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestOMNIAgent_BroadcastTransactionRPC2(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")

	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("2N6bnijfVT1FSUfDV65foDXrarbpThr7D5V", "n3K9Zmw3TKjfvmKVfT4Td2AgqTzkyb3vtz", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _ = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	redeemScriptStr := "52210304947980b04247e81a852977ab83902c0df3ada5a0fbed64a32c17679dd14c1421031d54ab830decbb3c2fe240b2a222d8bc2d475e56ac661d38deff732958e35c342103155ef7103cffe296d3411b0d919c76455962fc4f99946a34413577eb118f9c2653ae"

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1

	keyIndex1 := uint16(60)
	trxSig1Hex, err := ag.MultiSignRawTransactionRPC(rawTrx, redeemScriptStr, keyIndex1, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	keyIndex2 := uint16(61)
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
