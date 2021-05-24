package controller

import "github.com/kataras/iris"
import (
	"encoding/json"
	"model"
	"session"
	"strconv"
	"time"
	"utils"
)

func account_convert_log_value(no_type int, args ...string) string {
	if no_type == 1 {
		//用户修改账户状态
		return "用户修改账户状态： 真实姓名:" + args[0] + " ，手机号:" + args[1] + " ,原状态：" + args[2] + " ,新状态" + args[3] + " ，修改时间:" + time.Now().String()
	}else if no_type ==2 {
		//用户修改账户设置
		return "用户修改账户设置： 真实姓名:" + args[0] + " ，手机号:" + args[1] + " ,原状态：" + args[2] + " ,新状态" + args[3] +",钱包配置"+args[4]+ " ，修改时间:" + time.Now().String()
	}
	return ""
}

//list_accounts

type ListAccountsParam struct {
	SessionId string `json:"sessionid"`
	State     []int  `json:"state"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type ListAccountsRequest struct {
	Id      int                 `json:"id"`
	JsonRpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  []ListAccountsParam `json:"params"`
}

type ListAccountsResponse struct {
	Id     int                    `json:"id"`
	Result map[string]interface{} `json:"result"`
	Error  *utils.Error           `json:"error"`
}

func ListAccountsController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ListAccountsRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ListAccountsResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)

	if !exist || session_value.Role != 0 {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)

	user_accounts, err := model.GlobalDBMgr.AcctConfigMgr.ListNormalAccount(dbSession, req.Params[0].State, req.Params[0].Limit, req.Params[0].Offset)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	total_count, err := model.GlobalDBMgr.AcctConfigMgr.GetNormalAccountCount(dbSession, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	all_users_data := make([]map[string]interface{}, 0, len(user_accounts))
	res_data := make(map[string]interface{})
	for _, tbl_acct_config := range user_accounts {
		one_res_data := make(map[string]interface{})
		one_res_data["acctid"] = tbl_acct_config.Acctid
		one_res_data["cellnumber"] = tbl_acct_config.Cellphone
		one_res_data["realname"] = tbl_acct_config.Realname
		one_res_data["idcard"] = tbl_acct_config.Idcard
		one_res_data["state"] = tbl_acct_config.State
		one_res_data["regtime"] = tbl_acct_config.Createtime
		all_users_data = append(all_users_data, one_res_data)
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	res_data["total"] = total_count
	res_data["accts"] = all_users_data

	res.Result = res_data
	ctx.JSON(res)
	return
}

//get_account

type GetAccountParam struct {
	SessionId string `json:"sessionid"`
	AcctId    []int    `json:"acctid"`
}

type GetAccountRequest struct {
	AcctId      int           `json:"acctid"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []GetAccountParam `json:"params"`
}

type GetAccountResponse struct {
	Id     int                    `json:"id"`
	Result []GetAccountResult `json:"result"`
	Error  *utils.Error           `json:"error"`
}
type GetAccountResult struct {
	AcctId      int               `json:"acctid"`
	CellNumber string            `json:"cellnumber"`
	RealName  string            `json:"realname"`
	IdCard  string 				`json:"idcard"`
	State	int					`json:"state"`
	RegTime	time.Time		`json:"regtime"`
	PubKey string 		`json:"pubkey"`
}

func GetAccountController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetAccountRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res GetAccountResponse
	res.Id = req.AcctId
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	_, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	tbl_acct_configs, err := model.GlobalDBMgr.AcctConfigMgr.GetAccountsByIds(dbSession, req.Params[0].AcctId)
	if err != nil {
		res.Error = utils.MakeError(400007)
		ctx.JSON(res)
		return
	}
	for _,tbl_acct_config := range tbl_acct_configs{
		res.Result = append(res.Result,GetAccountResult{tbl_acct_config.Acctid,tbl_acct_config.Cellphone,tbl_acct_config.Realname,tbl_acct_config.Idcard[:4] +"********"+ tbl_acct_config.Idcard[len(tbl_acct_config.Idcard)-3:],tbl_acct_config.State,tbl_acct_config.Createtime,tbl_acct_config.Pubkey})
	}
	ctx.JSON(res)
	return

}

//change_account_state

type ChangeAccountStateParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	AcctId    int    `json:"acctid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}

type ChangeAccountStateRequest struct {
	Id      int                       `json:"id"`
	JsonRpc string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  []ChangeAccountStateParam `json:"params"`
}

type ChangeAccountStateResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ChangeAccountStateController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ChangeAccountStateRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ChangeAccountStateResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist || session_value.Role != 0 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 2, req.Params[0].MgmtId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}
	//check signature
	sig_origin_data := "change_acct_state," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].MgmtId) + "," + strconv.Itoa(req.Params[0].AcctId) + "," + strconv.Itoa(req.Params[0].State)
	err = utils.RsaVerySignWithSha1Hex(sig_origin_data, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}
	acct_old, _ := model.GlobalDBMgr.AcctConfigMgr.GetAccountById(dbSession, req.Params[0].AcctId)
	err = model.GlobalDBMgr.AcctConfigMgr.UpdateAcctState(dbSession, req.Params[0].AcctId, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400007)
		ctx.JSON(res)
		return
	}
	if req.Params[0].State == 2 {
		session.GlobalSessionMgr.DeleteSessionValueByAcctId(dbSession, req.Params[0].AcctId)
	}
	if req.Params[0].State == 1 {
		err = model.GlobalDBMgr.NotificationMgr.DeleteRegisterNotification(dbSession, session_value.AcctId, 0, strconv.Itoa(req.Params[0].AcctId))
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "delete", "delete register notification")
			ctx.JSON(res)
			return
		}
	}
	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 3, account_convert_log_value(1, acct_old.Realname, acct_old.Cellphone, strconv.Itoa(acct_old.State), strconv.Itoa(req.Params[0].State)))
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


//modify_acct

type ModifyAcctParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	AcctId    int    `json:"acctid"`
	Walletid  []int	 `json:"walletid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}

type ModifyAcctRequest struct {
	Id      int                       `json:"id"`
	JsonRpc string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  []ModifyAcctParam `json:"params"`
}

type ModifyAcctResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ModifyAcctController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ModifyAcctRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ModifyAcctResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist || session_value.Role != 0 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession,2, req.Params[0].MgmtId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}
	//check signature
	sig_origin_data := "modify_acct," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].MgmtId) + "," + strconv.Itoa(req.Params[0].AcctId)+ "," +utils.IntArrayToString(req.Params[0].Walletid) + "," + strconv.Itoa(req.Params[0].State)
	err = utils.RsaVerySignWithSha1Hex(sig_origin_data, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}
	acct_old, err := model.GlobalDBMgr.AcctConfigMgr.GetAccountById(dbSession, req.Params[0].AcctId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctConfigMgr.TableName, "query", "query account by id")
		ctx.JSON(res)
		return
	}

	err = model.GlobalDBMgr.AcctConfigMgr.UpdateAcctState(dbSession, req.Params[0].AcctId, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400007)
		ctx.JSON(res)
		return
	}
	if req.Params[0].State == 2 {
		session.GlobalSessionMgr.DeleteSessionValueByAcctId(dbSession, req.Params[0].AcctId)
	}
	if req.Params[0].State == 1 {
		err = model.GlobalDBMgr.NotificationMgr.DeleteRegisterNotification(dbSession, session_value.AcctId, 0, strconv.Itoa(req.Params[0].AcctId))
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.NotificationMgr.TableName, "delete", "delete register notification")
			ctx.JSON(res)
			return
		}
	}
	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 3, account_convert_log_value(2, acct_old.Realname, acct_old.Cellphone, strconv.Itoa(acct_old.State), strconv.Itoa(req.Params[0].State),utils.IntArrayToString(req.Params[0].Walletid)))
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.OperationLogMgr.TableName, "insert", "new operator log")
		ctx.JSON(res)
		return
	}

	//remove old relation
	old_relations, err :=model.GlobalDBMgr.AcctWalletRelationMgr.GetRelationsByAcctId(dbSession, req.Params[0].AcctId)
	if err !=nil{
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.OperationLogMgr.TableName, "delete", "delete relation")
		ctx.JSON(res)
		return
	}
	tmp_relation_map := make(map[int]int)
	for _,one_wallet_id := range req.Params[0].Walletid{
		tmp_relation_map[one_wallet_id] =0
	}

	for _,one_wallet_relation := range old_relations{
		_,exist := tmp_relation_map[one_wallet_relation.Walletid]
		if !exist{
			err = model.GlobalDBMgr.AcctWalletRelationMgr.DeleteRelation(dbSession, one_wallet_relation.Relationid)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "delete", "delete relation")
				ctx.JSON(res)
				return
			}
		}

	}
	for _,one_wallet_id := range req.Params[0].Walletid{
		err = model.GlobalDBMgr.AcctWalletRelationMgr.InsertRelation(dbSession, req.Params[0].AcctId,one_wallet_id)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "insert", "insert relation")
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


//get_account_wallets

type GetAccountWalletsParam struct {
	SessionId string `json:"sessionid"`
	AcctId    int    `json:"acctid"`
}

type GetAccountWalletsRequest struct {
	Id      int                       `json:"id"`
	JsonRpc string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  []GetAccountWalletsParam `json:"params"`
}

type GetAccountWalletRecord struct{
	Walletid int `json:"walletid"`
	Walletname string `json:"walletname"`
	Coinid int 	`json:"coinid"`
	Address string `json:"address"`
	State int  `json:"state"`

}

type GetAccountWalletsResponse struct {
	Id     int          `json:"id"`
	Result []GetAccountWalletRecord         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func GetAccountWalletsController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetAccountWalletsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res GetAccountWalletsResponse
	res.Id = req.Id
	res.Result = make([]GetAccountWalletRecord,0)
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	if session_value.Role == 1 && session_value.AcctId != req.Params[0].AcctId{
		res.Error = utils.MakeError(400008)
		ctx.JSON(res)
		return
	}

	relations,err := model.GlobalDBMgr.AcctWalletRelationMgr.GetRelationsByAcctId(dbSession, req.Params[0].AcctId)
	if err != nil {
		res.Error = utils.MakeError(400007)
		ctx.JSON(res)
		return
	}
	wallet_ids := make([]int,0,len(relations))
	for _,one_relation := range relations{
		wallet_ids = append(wallet_ids,one_relation.Walletid)
	}
	if len(wallet_ids)>0{
		wallets,err :=model.GlobalDBMgr.WalletConfigMgr.GetWalletsByIds(dbSession, wallet_ids)
		if err != nil {
			res.Error = utils.MakeError(400007)
			ctx.JSON(res)
			return
		}
		for _,one_wallet := range wallets{
			res.Result = append(res.Result,GetAccountWalletRecord{one_wallet.Walletid,one_wallet.Walletname,one_wallet.Coinid,one_wallet.Address,one_wallet.State})
		}
	}
	ctx.JSON(res)
	return
}

func AccountController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse

	if funcName == "list_accounts" {
		ListAccountsController(ctx, jsonRpcBody)
	} else if funcName == "get_account" {
		GetAccountController(ctx, jsonRpcBody)
	} else if funcName == "change_acct_state" {
		ChangeAccountStateController(ctx, jsonRpcBody)
	} else if funcName == "get_acct_wallets" {
		GetAccountWalletsController(ctx, jsonRpcBody)
	} else if funcName == "modify_acct" {
		ModifyAcctController(ctx, jsonRpcBody)
	}else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
