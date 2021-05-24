package controller

import "utils"


type InitPubKeyParamWS struct {
	ServerName 		string                  `json:"servername"`
	IsLocalServer   int                     `json:"islocalserver"`
	ServerPubkey 	string                  `json:"serverpubkey"`
	StartIndex      int                     `json:"startindex"`
	Keys            map[string]string       	`json:"keys"`
}

type InitPubKeyRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []InitPubKeyParamWS 		`json:"params"`
}

type InitPubKeyResultWS struct {
	Serverid   		int    		`json:"serverid"`
}

type InitPubKeyResponseWS struct {
	Id     int                      `json:"id"`
	Result InitPubKeyResultWS 		`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type QueryPubKeyRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []interface{} 			`json:"params"`
}

type QueryPubKeyResultWS struct {
	ServerId        int                     `json:"serverid"`
	ServerName 		string                  `json:"servername"`
	StartIndex      int                     `json:"startindex"`
	Keys            map[string]string       `json:"keys"`
}

type QueryPubKeyResponseWS struct {
	Id     int                      `json:"id"`
	Result []QueryPubKeyResultWS 	`json:"result"`
	Error  *utils.Error             `json:"error"`
}
