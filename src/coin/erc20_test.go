package coin

import (
	"config"
	"fmt"
	"testing"
)

func TestERC20Agent_CreateERC20Transaction(t *testing.T) {
	agent := ERC20Agent{}
	agent.Init("http://192.168.1.164:28000")
	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	raw_trx, err := agent.CreateERC20Transaction("TTT", "0xDB3884C583eB79F72B340F0968FCf8Bee9bd241A", "0x5b8455c9613baaba20296c52ddceb607f47222b6", "4000", "20000000000", "0.001", "", 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(raw_trx)
	trxId, err := agent.BroadcastTransaction(raw_trx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(trxId)

}

func TestERC20Agent_GetERC20BalanceByAddress(t *testing.T) {
	agent := ERC20Agent{}
	agent.Init("http://192.168.1.164:28000")
	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	_, balance, err := agent.GetERC20BalanceByAddress("NULS", "0xb2b4044E4C71bdfBcE06db946c6f6477f8B21351")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)

}