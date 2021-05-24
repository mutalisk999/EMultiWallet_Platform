package controller

import (
	"encoding/json"
	"github.com/kataras/iris"
	"model"
	"session"
	"utils"
)

type GetIdentityParam struct {
	SessionId string `json:"sessionid"`
	IdType    int    `json:"idtype"`
}

type GetIdentityRequest struct {
	Id      int                `json:"id"`
	JsonRpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  []GetIdentityParam `json:"params"`
}

type GetIdentityResponse struct {
	Id     int          `json:"id"`
	Result int          `json:"result"`
	Error  *utils.Error `json:"error"`
}

func GetIdentityController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req GetIdentityRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetIdentityResponse
	res.Id = req.Id
	res.Result = 0
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	if req.Params[0].IdType == 1 && len(req.Params[0].SessionId) != 0 {
		res.Error = utils.MakeError(200002)
		ctx.JSON(res)
		return
	} else if req.Params[0].IdType != 1 && len(req.Params[0].SessionId) == 0 {
		res.Error = utils.MakeError(200003)
		ctx.JSON(res)
		return
	}
	if req.Params[0].IdType != 1 {
		hasSessionId := session.GlobalSessionMgr.HasSessionId(req.Params[0].SessionId)
		if !hasSessionId {
			res.Error = utils.MakeError(200004)
			ctx.JSON(res)
			return
		}
	}
	mgr := model.GlobalDBMgr.SequenceMgr
	res.Result, err = mgr.NewSequence(dbSession, req.Params[0].IdType)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, mgr.TableName, "insert", "create next sequence")
		ctx.JSON(res)
		return
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	ctx.JSON(res)
}

func IdentityController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "get_auto_inc_id" {
		GetIdentityController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
