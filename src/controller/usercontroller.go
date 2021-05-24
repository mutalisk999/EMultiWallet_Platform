package controller

import (
	"github.com/kataras/iris"
	"utils"
	"authcode"
	"encoding/json"
	"model"
	"session"
	"strconv"
	"time"
)

type UserRegisterParam struct {
	CellNumber   string `json:"cellnumber"`
	RealName     string `json:"realname"`
	IdCard       string `json:"idcard"`
	VerifyCodeId string `json:"verifycodeid"`
	VerifyCode   string `json:"verifycode"`
	Pubkey       string `json:"pubkey"`
}

type UserRegisterRequest struct {
	Id      int                 `json:"id"`
	JsonRpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  []UserRegisterParam `json:"params"`
}

type UserRegisterResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func user_convert_notification_value(no_type int, args ...string) string {
	if no_type == 1 {
		//用户注册
		return "新用户注册，手机号:" + args[0] + "，真实姓名:" + args[1] + "，身份证号码:" + args[2] + "，注册时间:" + utils.TimeToFormatString(time.Now())
	}
	return ""
}

func user_convert_log_value(no_type int, args ...string) string {
	if no_type == 1 {
		//用户注册
		return "新用户注册，手机号:" + args[0] + "，真实姓名:" + args[1] + "，身份证号码:" + args[2] + "，注册时间:" + utils.TimeToFormatString(time.Now())
	} else if no_type == 2 {
		//用户注册
		return "新用户注册，手机号:" + args[0] + "，真实姓名:" + args[1] + "，身份证号码:" + args[2] + "，注册时间:" + utils.TimeToFormatString(time.Now()) + "，未知异常，插入后数据库查询失败"
	} else if no_type == 3 {
		//用户登录
		return "用户登录，登录ID:" + args[0] + " ,真实姓名:" + args[1] + " ，手机号:" + args[2] + " ，登录时间:" + utils.TimeToFormatString(time.Now())
	} else if no_type == 4 {
		//用户注销
		return "用户注销： 真实姓名:" + args[0] + " ，手机号:" + args[1] + " ，注销时间:" + utils.TimeToFormatString(time.Now())
	}
	return ""
}

//user_register
func UserRegisterController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req UserRegisterRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res UserRegisterResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	VerifyRes := authcode.VerifyAuthCode(req.Params[0].VerifyCodeId, req.Params[0].VerifyCode)
	if !VerifyRes {
		res.Error = utils.MakeError(400001)
		ctx.JSON(res)
		return
	}
	convert_pubkey, rest := utils.RsaReadPEMPublicKey(req.Params[0].Pubkey)
	if rest != "" {
		res.Error = utils.MakeError(400013)
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctConfigMgr.VerifyUnique(dbSession, req.Params[0].CellNumber, req.Params[0].RealName, req.Params[0].IdCard, convert_pubkey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400011)
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctConfigMgr.InsertAcct(dbSession, req.Params[0].CellNumber, req.Params[0].RealName, req.Params[0].IdCard, convert_pubkey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.FormatSysError(err)
		ctx.JSON(res)
		return
	}
	//add notification
	//account_count,err := model.GlobalDBMgr.AcctConfigMgr.GetAccountCount(dbSession)
	//if account_count > 1 && err==nil{
	//	adminId := model.GlobalDBMgr.AcctConfigMgr.GetAdminId(dbSession)
	//	acctId, _ := model.GlobalDBMgr.AcctConfigMgr.GetAccountIdByPubkey(dbSession, convert_pubkey)
	//	_, err = model.GlobalDBMgr.NotificationMgr.NewNotification(dbSession, &adminId, nil, nil, 0, user_convert_notification_value(1, req.Params[0].CellNumber, req.Params[0].RealName, req.Params[0].IdCard), 0, strconv.Itoa(acctId), "")
	//	if err != nil {
	//		dbSession.Rollback()
	//		res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "insert", "new notification")
	//		ctx.JSON(res)
	//		return
	//	}
	//}
	//add log
	acct_id, err := model.GlobalDBMgr.AcctConfigMgr.GetAccountIdByPubkey(dbSession, convert_pubkey)
	if err == nil {
		_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acct_id, 0, user_convert_log_value(1, req.Params[0].CellNumber, req.Params[0].RealName, req.Params[0].IdCard))
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "insert", "new operator log")
			ctx.JSON(res)
			return
		}
	} else {
		_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acct_id, 0, user_convert_log_value(2, req.Params[0].CellNumber, req.Params[0].RealName, req.Params[0].IdCard))
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "insert", "new operator log")
			ctx.JSON(res)
			return
		}
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	res.Result = true
	ctx.JSON(res)
	return
}

//user_login

type UserLoginParam struct {
	LoginId   int    `json:"loginid"`
	Pubkey    string `json:"pubkey"`
	Signature string `json:"signature"`
}

type UserLoginRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []UserLoginParam `json:"params"`
}

type UserLoginResponse struct {
	Id     int                    `json:"id"`
	Result map[string]interface{} `json:"result"`
	Error  *utils.Error           `json:"error"`
}

//user_login
func UserLoginController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req UserLoginRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res UserLoginResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	verify_res, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 1, req.Params[0].LoginId)
	if !verify_res || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}
	sig_origin := "user_login," + strconv.Itoa(req.Params[0].LoginId)
	convert_pubkey, rest := utils.RsaReadPEMPublicKey(req.Params[0].Pubkey)
	if rest != "" {
		dbSession.Rollback()
		res.Error = utils.MakeError(400013)
		ctx.JSON(res)
		return
	}
	err = utils.RsaVerySignWithSha1Hex(sig_origin, req.Params[0].Signature, convert_pubkey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}
	acct, err := model.GlobalDBMgr.AcctConfigMgr.FindAccountByPubkey(dbSession, convert_pubkey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400003)
		ctx.JSON(res)
		return
	}
	if acct.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(400012)
		ctx.JSON(res)
		return
	}
	session.GlobalSessionMgr.DeleteSessionValueByAcctId(dbSession, acct.Acctid)
	var sessionValue session.SessionValue
	sessionValue.CellNumber = acct.Cellphone
	sessionValue.IdCard = acct.Idcard
	sessionValue.PubKey = acct.Pubkey
	sessionValue.RealName = acct.Realname
	sessionValue.AcctId = acct.Acctid
	sessionValue.Role = acct.Role
	sessionId, err := session.GlobalSessionMgr.NewSessionValue(dbSession, sessionValue)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400004)
		ctx.JSON(res)
		return
	}
	res_data := make(map[string]interface{})
	res_data["sessionid"] = sessionId
	res_data["usertype"] = acct.Role
	res_data["acctid"] = acct.Acctid
	res.Result = res_data

	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acct.Acctid, 1, user_convert_log_value(3, strconv.Itoa(req.Params[0].LoginId), acct.Realname, acct.Cellphone))
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "insert", "new operator log")
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

//user_logout

type UserLogoutParam struct {
	SessionId string `json:"sessionid"`
	Signature string `json:"signature"`
}

type UserLogoutRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []UserLogoutParam `json:"params"`
}

type UserLogoutResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func UserLogoutController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req UserLogoutRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res UserLogoutResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sig_origin := "user_logout," + req.Params[0].SessionId
	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	err = utils.RsaVerySignWithSha1Hex(sig_origin, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	session.GlobalSessionMgr.DeleteSessionValue(dbSession, req.Params[0].SessionId)

	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 2, user_convert_log_value(4, session_value.RealName, session_value.CellNumber))
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.OperationLogMgr.TableName, "insert", "new operator log")
		ctx.JSON(res)
		return
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	res.Result = true
	ctx.JSON(res)
	return
}

//user_getinfo

type UserGetInfoParam struct {
	Pubkey string `json:"pubkey"`
}

type UserGetInfoRequest struct {
	Id      int                `json:"id"`
	JsonRpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  []UserGetInfoParam `json:"params"`
}

type UserGetInfoResponse struct {
	Id     int                    `json:"id"`
	Result map[string]interface{} `json:"result"`
	Error  *utils.Error           `json:"error"`
}

func UserGetInfoController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req UserGetInfoRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res UserGetInfoResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	convert_pubkey, rest := utils.RsaReadPEMPublicKey(req.Params[0].Pubkey)
	if rest != "" {
		res.Error = utils.MakeError(400013)
		ctx.JSON(res)
		return
	}
	tbl_acct_config, err := model.GlobalDBMgr.AcctConfigMgr.FindAccountByPubkey(dbSession, convert_pubkey)
	if err != nil {
		res.Error = utils.MakeError(400003)
		ctx.JSON(res)
		return
	}
	res_data := make(map[string]interface{})
	res_data["id"] = tbl_acct_config.Acctid
	res_data["cellnumber"] = tbl_acct_config.Cellphone
	res_data["realname"] = tbl_acct_config.Realname
	res_data["idcard"] = tbl_acct_config.Idcard
	res_data["state"] = tbl_acct_config.State
	res_data["regtime"] = tbl_acct_config.Createtime
	res.Result = res_data
	ctx.JSON(res)
	return
}

func UserController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res utils.JsonRpcResponse

	if funcName == "user_register" {
		UserRegisterController(ctx, jsonRpcBody)
	} else if funcName == "user_login" {
		UserLoginController(ctx, jsonRpcBody)
	} else if funcName == "user_logout" {
		UserLogoutController(ctx, jsonRpcBody)
	} else if funcName == "user_getinfo" {
		UserGetInfoController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
