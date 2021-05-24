package controller

import (
	"encoding/json"
	"github.com/kataras/iris"
	"model"
	"session"
	"utils"
)

type GetOpLogsParam struct {
	SessionId string    `json:"sessionid"`
	AcctId    []int     `json:"acctid"`
	OpType    []int     `json:"optype"`
	OpTime    [2]string `json:"optime"`
	OffSet    int       `json:"offset"`
	Limit     int       `json:"limit"`
}

type GetOpLogsRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []GetOpLogsParam `json:"params"`
}

type OperatorLog struct {
	LogId   int    `json:"logid"`
	AcctId  int    `json:"acctid"`
	OpType  int    `json:"optype"`
	OpTime  string `json:"optime"`
	Content string `json:"content"`
}

type GetOpLogsResult struct {
	Total int           `json:"total"`
	Logs  []OperatorLog `json:"logs"`
}

type GetOpLogsResponse struct {
	Id     int              `json:"id"`
	Result *GetOpLogsResult `json:"result"`
	Error  *utils.Error     `json:"error"`
}

func GetOpLogsController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetOpLogsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetOpLogsResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	if req.Params[0].OpTime[0] != "" {
		res.Error = utils.CheckDateTimeString("starttime", req.Params[0].OpTime[0])
		if res.Error != nil {
			ctx.JSON(res)
			return
		}
	}
	if req.Params[0].OpTime[1] != "" {
		res.Error = utils.CheckDateTimeString("endtime", req.Params[0].OpTime[1])
		if res.Error != nil {
			ctx.JSON(res)
			return
		}
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	if sessionValue.Role == 1 {
		if len(req.Params[0].AcctId) != 1 || req.Params[0].AcctId[0] != sessionValue.AcctId {
			res.Error = utils.MakeError(200005)
			ctx.JSON(res)
			return
		}
	}

	mgr := model.GlobalDBMgr.OperationLogMgr
	totalCount, operatorlogs, err := mgr.GetOperatorLogs(dbSession, req.Params[0].AcctId, req.Params[0].OpType, req.Params[0].OpTime, req.Params[0].OffSet,
		req.Params[0].Limit)
	if err != nil {
		res.Error = utils.MakeError(300001, mgr.TableName, "query", "query operator logs")
		ctx.JSON(res)
		return
	}
	res.Result = new(GetOpLogsResult)
	res.Result.Total = totalCount
	res.Result.Logs = make([]OperatorLog, len(operatorlogs), len(operatorlogs))
	for i, operatorlog := range operatorlogs {
		res.Result.Logs[i] = OperatorLog{
			operatorlog.Logid, operatorlog.Acctid, operatorlog.Optype,
			utils.TimeToFormatString(operatorlog.Createtime),
			operatorlog.Content}
	}
	ctx.JSON(res)
}

func LogController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "get_op_logs" {
		GetOpLogsController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
