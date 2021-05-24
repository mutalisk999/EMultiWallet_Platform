package controller

import "utils"

type QueryMessageRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []interface{} 			`json:"params"`
}

type QueryMessageResponseWS struct {
	Id     int                      `json:"id"`
	Result string               	`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type VerifyMessageParam struct {
	ServerId     int                `json:"serverid"`
	SignedMessage string            `json:"signedmessage"`
}

type VerifyMessageRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []VerifyMessageParam 	`json:"params"`
}

type VerifyMessageResponseWS struct {
	Id     int                      `json:"id"`
	Result bool                 	`json:"result"`
	Error  *utils.Error             `json:"error"`
}

