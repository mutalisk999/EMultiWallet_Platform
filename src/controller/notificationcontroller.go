package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"session"
	"strconv"
	"strings"
	"utils"
	"github.com/go-xorm/xorm"
)

const (
	LogFormatTypeChangeNotifyState = 1
)

func GetNotificationLogFormat(fmtType int) string {
	if fmtType == LogFormatTypeChangeNotifyState {
		return "用户[%s]设定通知状态为[%s],通知ID:[%s],修改结果:[%s]"
	}
	return ""
}

type ListNotificationParam struct {
	SessionId string `json:"sessionid"`
}

type ListNotificationRequest struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []ListNotificationParam `json:"params"`
}

type ListNotificationResult struct {
	NotifyId   int    `json:"notifyid"`
	NotifyType int    `json:"notifytype"`
	Content    string `json:"content"`
	CreateTime string `json:"createtime"`
}

type ListNotificationResponse struct {
	Id     int                      `json:"id"`
	Result []ListNotificationResult `json:"result"`
	Error  *utils.Error             `json:"error"`
}

func ListNotificationController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req ListNotificationRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res ListNotificationResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	// 列出所有未处理过的提醒
	mgr := model.GlobalDBMgr.NotificationMgr
	notifications, err := mgr.ListNotifications(dbSession, sessionValue.AcctId)
	if err != nil {
		res.Error = utils.MakeError(300001, mgr.TableName, "query", "list notifications")
		ctx.JSON(res)
		return
	}
	res.Result = make([]ListNotificationResult, len(notifications), len(notifications))
	for i, notification := range notifications {
		res.Result[i] = ListNotificationResult{
			notification.Notifyid,
			notification.Notifytype,
			notification.Notification,
			utils.TimeToFormatString(notification.Createtime)}
	}
	ctx.JSON(res)
}

func ChangeNotificationStateLog(dbSession *xorm.Session, isSuccQuit bool, acctId int, notifyId []int, state int) error {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(dbSession, acctId)
	if err != nil {
		return err
	}

	stateStr := ""
	if state == 0 {
		stateStr = "未处理"
	} else if state == 1 {
		stateStr = "已处理"
	} else if state == 2 {
		stateStr = "已忽略"
	} else {
		stateStr = "未知状态"
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	notifyIdStr := utils.IntArrayToString(notifyId)
	logContent := fmt.Sprintf(GetNotificationLogFormat(LogFormatTypeChangeNotifyState), acctConfig.Realname,
		stateStr, notifyIdStr, resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(dbSession, acctId, 7, logContent)
	if err != nil {
		return err
	}
	return nil
}

type GetNotificationCountParam struct {
	SessionId string `json:"sessionid"`
	State     int    `json:"state"`
}

type GetNotificationCountRequest struct {
	Id      int                         `json:"id"`
	JsonRpc string                      `json:"jsonrpc"`
	Method  string                      `json:"method"`
	Params  []GetNotificationCountParam `json:"params"`
}

type GetNotificationCountResponse struct {
	Id     int          `json:"id"`
	Result int          `json:"result"`
	Error  *utils.Error `json:"error"`
}

func GetNotificationCountController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetNotificationCountRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetNotificationCountResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	mgr := model.GlobalDBMgr.NotificationMgr
	res.Result, err = mgr.GetNotificationCount(dbSession, sessionValue.AcctId, req.Params[0].State)
	if err != nil {
		res.Error = utils.MakeError(300001, mgr.TableName, "query", "get notification count")
		ctx.JSON(res)
		return
	}

	ctx.JSON(res)
	return
}

type ChangeNotificationStateParam struct {
	SessionId string `json:"sessionid"`
	OperateId int    `json:"operateid"`
	NotifyId  []int  `json:"notifyid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}

type ChangeNotificationStateRequest struct {
	Id      int                            `json:"id"`
	JsonRpc string                         `json:"jsonrpc"`
	Method  string                         `json:"method"`
	Params  []ChangeNotificationStateParam `json:"params"`
}

type ChangeNotificationStateResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ChangeNotificationStateController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ChangeNotificationStateRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res ChangeNotificationStateResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		dbSession.Rollback()
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 5, req.Params[0].OperateId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	funcNameStr := "change_state"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	notifyIdStr := utils.IntArrayToString(req.Params[0].NotifyId)
	stateStr := strconv.Itoa(req.Params[0].State)
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, notifyIdStr, stateStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	mgr := model.GlobalDBMgr.NotificationMgr

	if sessionValue.Role != 0 {
		notifications, err := mgr.ListNotifications(dbSession, sessionValue.AcctId)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, mgr.TableName, "query", "list notifications")
			//ChangeNotificationStateLog(false, sessionValue.AcctId, req.Params[0].NotifyId, req.Params[0].State)
			ctx.JSON(res)
			return
		}

		for _, idArg := range req.Params[0].NotifyId {
			isAcctNotify := false
			for _, notify := range notifications {
				if idArg == notify.Notifyid {
					isAcctNotify = true
					break
				}
			}
			if !isAcctNotify {
				dbSession.Rollback()
				res.Error = utils.MakeError(200015)
				//ChangeNotificationStateLog(false, sessionValue.AcctId, req.Params[0].NotifyId, req.Params[0].State)
				ctx.JSON(res)
				return
			}
		}
	}

	err = mgr.UpdateNotificationsState(dbSession, req.Params[0].NotifyId, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, mgr.TableName, "update", "update notification state")
		//ChangeNotificationStateLog(false, sessionValue.AcctId, req.Params[0].NotifyId, req.Params[0].State)
		ctx.JSON(res)
		return
	}

	err = ChangeNotificationStateLog(dbSession, true, sessionValue.AcctId, req.Params[0].NotifyId, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300010)
		ctx.JSON(res)
		return
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	ctx.JSON(res)
	return
}

func NotificationController(ctx iris.Context) {
	//id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	id, funcName, _, err := utils.ReadJsonRpcBody(ctx)

	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "list_notifications" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		//ListNotificationController(ctx, jsonRpcBody)
	} else if funcName == "get_notification_count" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		//GetNotificationCountController(ctx, jsonRpcBody)
	} else if funcName == "change_state" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		//ChangeNotificationStateController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
