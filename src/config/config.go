package config

import (

	"io/ioutil"
	"encoding/json"
)

type SupportCoin struct {
	CoinName        string
	CoinSymbol      string
	Precision       int
	ConfirmCount    int
	IsErc20         bool
	ContractAddress string
	IsOmni          bool
	OmniPropertyId  int
}

func NewSupportCoin(CoinName string, CoinSymbol string, Precision int, ConfirmCount int,
	IsErc20 bool, ContractAddress string, IsOmni bool, OmniPropertyId int) SupportCoin {
	var supportCoin SupportCoin
	supportCoin.CoinName = CoinName
	supportCoin.CoinSymbol = CoinSymbol
	supportCoin.Precision = Precision
	supportCoin.ConfirmCount = ConfirmCount
	supportCoin.IsErc20 = IsErc20
	supportCoin.ContractAddress = ContractAddress
	supportCoin.IsOmni = IsOmni
	supportCoin.OmniPropertyId = OmniPropertyId
	return supportCoin
}

var GlobalSupportCoinMgr = map[string]SupportCoin{
	"BTC": NewSupportCoin("BitCoin", "BTC", 8, 6, false, "", false, 0),
	"LTC": NewSupportCoin("LiteCoin", "LTC", 8, 6, false, "", false, 0),
	"BCH": NewSupportCoin("BitCoinCash", "BCH", 8, 6, false, "", false, 0),
	"ETH": NewSupportCoin("Ethereum", "ETH", 18, 30, false, "", false, 0),
	"1ST": NewSupportCoin("FirstBlood", "1ST", 18, 30, true, "0x158b477fde01f7aaaefb089d7545a4a9513e3f66", false, 0),
	//  0xDf7b27c61A475413bF07c9636a5B58CAE3CFFCf9
	"NULS": NewSupportCoin("NULS", "NULS", 18, 30, true, "0x19854b9A782B14AaBF2425e6fC2d18EC297F6839", false, 0),
	"USDT": NewSupportCoin("TetherUSD", "USDT", 8, 6, false, "", true, 31),

	// for test
	"OMNI": NewSupportCoin("Omni", "OMNI", 8, 6, false, "", true, 1),
	"TOMNI": NewSupportCoin("Test Omni", "TOMNI", 8, 6, false, "", true, 2),

	"UB": NewSupportCoin("UnitedBitCoin", "UB", 8, 6, false, "", false, 0),

}

func IsSupportCoin(coinSymbol string) bool {
	_, exist := GlobalSupportCoinMgr[coinSymbol]
	return exist
}


var IsTestEnvironment bool
var LocalServerId int

type DbConfig struct {
	DbType string   `json:"dbType"`
	DbSource string   `json:"dbSource"`
}

type CryptoDeviceConfig struct {
	DeviceIp   string `json:"deviceIp"`
	DevicePort uint16 `json:"devicePort"`
	TimeOut    uint16 `json:"timeOut"`
}

type ServerInfoConfig struct {
	ServerName   string `json:"serverName"`
	StartIndex 	 int    `json:"startIndex"`
}

type HttpConfig struct {
	EndPoint   string `json:"endPoint"`
	CharSet    string `json:"charSet"`
}

type HttpsConfig struct {
	EndPoint   string `json:"endPoint"`
	CertFile   string `json:"certFile"`
	KeyFile    string `json:"keyFile"`
	CharSet    string `json:"charSet"`
}

type WsConfig struct {
	Url   		string `json:"url"`
	Origin    	string `json:"origin"`
	TimeOut     uint    `json:"timeout"`
}

type WssConfig struct {
	Url   		string `json:"url"`
	Origin    	string `json:"origin"`
	TimeOut     uint    `json:"timeout"`
}

type Config struct {
	ServerInfoConfig          ServerInfoConfig          `json:"serverInfoConfig"`
	CryptoDeviceConfig        CryptoDeviceConfig        `json:"cryptoDeviceConfig"`
	DbConfig                  DbConfig                  `json:"dbConfig"`
	IsHttps                   bool                      `json:"isHttps"`
	HttpConfig                HttpConfig                `json:"httpConfig"`
	HttpsConfig               HttpsConfig                `json:"httpsConfig"`
	IsWss                     bool                       `json:"isWss"`
	WsConfig                  WsConfig                   `json:"wsConfig"`
	WssConfig                 WssConfig                  `json:"wssConfig"`
}

var GlobalConfig Config

type JsonStruct struct {
}

func (j *JsonStruct) Load(configFile string, config interface{}) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}
