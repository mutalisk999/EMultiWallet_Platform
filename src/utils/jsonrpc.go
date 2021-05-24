package utils

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
)

type JsonRpcRequest struct {
	Id      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonRpcResponse struct {
	Id     int          `json:"id"`
	Result *interface{} `json:"result"`
	Error  *Error       `json:"error"`
}

func ReadJsonRpcBody(ctx iris.Context) (int, string, []byte, error) {
	reader := ctx.Request().Body
	bodyBytes := make([]byte, ctx.Request().ContentLength, ctx.Request().ContentLength)
	readCount, err := reader.Read(bodyBytes)
	if err.Error() != "EOF" || int64(readCount) != ctx.Request().ContentLength {
		return 0, "", nil, errors.New("read http jsonrpc body error")
	}
	var jsonRpcRequest JsonRpcRequest
	err = json.Unmarshal(bodyBytes, &jsonRpcRequest)
	if err != nil {
		return 0, "", nil, err
	}
	return jsonRpcRequest.Id, jsonRpcRequest.Method, bodyBytes, nil
}

type TaskPushParamWS struct {
	TaskUuid 		string                  `json:"taskuuid"`
	WalletUuid 		string                  `json:"walletuuid"`
	TrxUuid 		string                  `json:"trxuuid"`
	PushType		int                     `json:"pushtype"`
}

type TaskPushRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []TaskPushParamWS 		`json:"params"`
}

type TaskPushResponseWS struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *Error		`json:"error"`
}

type TaskAcceptParamWS struct {
	TaskUuid 		string                  `json:"taskuuid"`
	State			int                     `json:"state"`
}

type TaskAcceptRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []TaskAcceptParamWS 		`json:"params"`
}

type TaskAcceptResponseWS struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *Error		`json:"error"`
}

