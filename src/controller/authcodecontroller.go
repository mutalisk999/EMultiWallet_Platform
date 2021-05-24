package controller

import (
	"authcode"
	"encoding/json"
	"github.com/kataras/iris"
	"utils"
)

type GetAuthCodeParam struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Len    int `json:"len"`
}

type GetAuthCodeRequest struct {
	Id      int                `json:"id"`
	JsonRpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  []GetAuthCodeParam `json:"params"`
}

type GetAuthCodeResult struct {
	AuthCodeId     string `json:"authcodeid"`
	AuthCodeStream string `json:"authcodestream"`
}

type GetAuthCodeResponse struct {
	Id     int                `json:"id"`
	Result *GetAuthCodeResult `json:"result"`
	Error  *utils.Error       `json:"error"`
}

func GetAuthCodeController(ctx iris.Context, jsonRpcBody []byte) {
	var req GetAuthCodeRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetAuthCodeResponse
	res.Id = req.Id
	res.Result = nil
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	res.Error = utils.CheckInteger("height", req.Params[0].Height, 20, 200)
	if res.Error != nil {
		ctx.JSON(res)
		return
	}
	res.Error = utils.CheckInteger("width", req.Params[0].Width, 40, 400)
	if res.Error != nil {
		ctx.JSON(res)
		return
	}
	res.Error = utils.CheckInteger("len", req.Params[0].Len, 1, 10)
	if res.Error != nil {
		ctx.JSON(res)
		return
	}

	code := authcode.CreateAuthCode(req.Params[0].Height, req.Params[0].Width, req.Params[0].Len)
	res.Result = new(GetAuthCodeResult)
	res.Result.AuthCodeId = code.AuthCodeId
	res.Result.AuthCodeStream = code.Base64PicData
	res.Error = nil
	ctx.JSON(res)
}

func AuthCodeController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "get_authcode" {
		GetAuthCodeController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
