package controller

import "github.com/kataras/iris"
import (
	"config"
	"encoding/json"
	"fmt"
	"model"
	"session"
	"strconv"
	"time"
	"utils"
)

func coin_convert_log_value(no_type int, args ...string) string {
	if no_type == 1 {
		//创建币种
		return "创建币种，操作ID:" + args[0] + "，币种标识:" + args[1] + "，币种状态:" + args[2] + "创建时间:" + time.Now().String()
	} else if no_type == 2 {
		//修改币种信息
		return "修改币种信息，操作ID:" + args[0] + "，币种ID:" + args[1] + "，币种标识:" + args[2] + "修改时间:" + time.Now().String()
	} else if no_type == 3 {
		//修改币种状态
		return "修改币种状态，操作ID:" + args[0] + "，币种ID:" + args[1] + "，币种标识:" + args[2] + "，旧币种状态:" + args[3] + "，新币种状态:" + args[4] + "修改时间:" + time.Now().String()
	} else if no_type == 4 {
		//用户注销
		return "用户注销： 真实姓名:" + args[0] + " ，手机号:" + args[1] + " ，注销时间:" + time.Now().String()
	} else if no_type == 5 {
		//用户修改账户状态
		return "用户修改账户状态： 真实姓名:" + args[0] + " ，手机号:" + args[1] + " ,原状态：" + args[2] + " ,新状态" + args[3] + " ，修改时间:" + time.Now().String()
	} else if no_type == 6 {
		//用户刷新币种
		return "用户刷新币种"
	}
	return ""
}

//list_coins

type ListCoinsParam struct {
	SessionId string `json:"sessionid"`
}

type ListCoinsRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []ListCoinsParam `json:"params"`
}

type ListCoinsResponse struct {
	Id     int                      `json:"id"`
	Result []map[string]interface{} `json:"result"`
	Error  *utils.Error             `json:"error"`
}

func ListCoinsController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req ListCoinsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ListCoinsResponse
	res.Id = req.Id
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

	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	all_coins, err := model.GlobalDBMgr.CoinConfigMgr.ListCoins(dbSession)
	if err != nil {
		fmt.Println(err.Error())
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	res_data := make([]map[string]interface{}, 0)
	for _, one_coin := range all_coins {
		one_rec := make(map[string]interface{})
		one_rec["coinid"] = one_coin.Coinid
		one_rec["coinsymbol"] = one_coin.Coinsymbol
		one_rec["state"] = one_coin.State
		if session_value.Role == 0 {
			one_rec["ip"] = one_coin.Ip
			one_rec["rpcport"] = one_coin.Rpcport
			one_rec["rpcuser"] = one_coin.Rpcuser
			one_rec["rpcpass"] = one_coin.Rpcpass
		} else {
			one_rec["ip"] = ""
			one_rec["rpcport"] = ""
			one_rec["rpcuser"] = ""
			one_rec["rpcpass"] = ""
		}
		res_data = append(res_data, one_rec)
	}

	res.Result = res_data
	ctx.JSON(res)
	return
}

//create_coin

type CreateCoinParam struct {
	SessionId  string `json:"sessionid"`
	MgmtId     int    `json:"mgmtid"`
	CoinSymbol string `json:"coinsymbol"`
	Ip         string `json:"ip"`
	RpcPort    int    `json:"rpcport"`
	RpcUser    string `json:"rpcuser"`
	RpcPass    string `json:"rpcpass"`
	State      int    `json:"state"`
	Signature  string `json:"signature"`
}

type CreateCoinRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []CreateCoinParam `json:"params"`
}

type CreateCoinResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func CreateCoinController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req CreateCoinRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res CreateCoinResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(req.Params[0].CoinSymbol) {
		res.Error = utils.MakeError(600001, req.Params[0].CoinSymbol)
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
	verify_res, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 6, req.Params[0].MgmtId)
	if !verify_res || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	sig_origin := "create_coin," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].MgmtId) + "," + req.Params[0].CoinSymbol + "," + req.Params[0].Ip + "," +
		strconv.Itoa(req.Params[0].RpcPort) + "," + req.Params[0].RpcUser + "," + req.Params[0].RpcPass + "," + strconv.Itoa(req.Params[0].State)
	err = utils.RsaVerySignWithSha1Hex(sig_origin, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	err = model.GlobalDBMgr.CoinConfigMgr.InsertCoin(dbSession, req.Params[0].CoinSymbol, req.Params[0].Ip, req.Params[0].RpcPort, req.Params[0].RpcUser, req.Params[0].RpcPass, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}

	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 6, coin_convert_log_value(1, strconv.Itoa(req.Params[0].MgmtId), req.Params[0].CoinSymbol, strconv.Itoa(req.Params[0].State)))
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

//modify_coin

type ModifyCoinParam struct {
	SessionId  string `json:"sessionid"`
	MgmtId     int    `json:"mgmtid"`
	CoinId     int    `json:"coinid"`
	CoinSymbol string `json:"coinsymbol"`
	Ip         string `json:"ip"`
	RpcPort    int    `json:"rpcport"`
	RpcUser    string `json:"rpcuser"`
	RpcPass    string `json:"rpcpass"`
	Signature  string `json:"signature"`
}

type ModifyCoinRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []ModifyCoinParam `json:"params"`
}

type ModifyCoinResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ModifyCoinController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ModifyCoinRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ModifyCoinResponse
	res.Id = req.Id
	res.Result = false
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(req.Params[0].CoinSymbol) {
		res.Error = utils.MakeError(600001, req.Params[0].CoinSymbol)
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
	verify_res, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 6, req.Params[0].MgmtId)
	if !verify_res || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	sig_origin := "modify_coin," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].MgmtId) + "," + strconv.Itoa(req.Params[0].CoinId) + "," + req.Params[0].CoinSymbol + "," + req.Params[0].Ip + "," +
		strconv.Itoa(req.Params[0].RpcPort) + "," + req.Params[0].RpcUser + "," + req.Params[0].RpcPass
	err = utils.RsaVerySignWithSha1Hex(sig_origin, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	err = model.GlobalDBMgr.CoinConfigMgr.UpdateCoin(dbSession, req.Params[0].CoinId, req.Params[0].CoinSymbol, req.Params[0].Ip, req.Params[0].RpcPort, req.Params[0].RpcUser, req.Params[0].RpcPass)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}

	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 6, coin_convert_log_value(2, strconv.Itoa(req.Params[0].MgmtId), strconv.Itoa(req.Params[0].CoinId), req.Params[0].CoinSymbol))
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

//change_coin_state

type ChangeCoinStateParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	CoinId    int    `json:"coinid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}

type ChangeCoinStateRequest struct {
	Id      int                    `json:"id"`
	JsonRpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  []ChangeCoinStateParam `json:"params"`
}

type ChangeCoinStateResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ChangeCoinStateController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ChangeCoinStateRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ChangeCoinStateResponse
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
	verify_res, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 6, req.Params[0].MgmtId)
	if !verify_res || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	sig_origin := "change_coin_state," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].MgmtId) + "," + strconv.Itoa(req.Params[0].CoinId) + "," + strconv.Itoa(req.Params[0].State)
	err = utils.RsaVerySignWithSha1Hex(sig_origin, req.Params[0].Signature, session_value.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	coin, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, req.Params[0].CoinId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	err = model.GlobalDBMgr.CoinConfigMgr.UpdateCoinState(dbSession, req.Params[0].CoinId, req.Params[0].State)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 6, coin_convert_log_value(3, strconv.Itoa(req.Params[0].MgmtId), strconv.Itoa(req.Params[0].CoinId), coin.Coinsymbol, strconv.Itoa(coin.State), strconv.Itoa(req.Params[0].State)))
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

//get_coin

type GetCoinParam struct {
	SessionId string `json:"sessionid"`
	CoinId    int    `json:"coinid"`
}

type GetCoinRequest struct {
	Id      int            `json:"id"`
	JsonRpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  []GetCoinParam `json:"params"`
}

type GetCoinResponse struct {
	Id     int                    `json:"id"`
	Result map[string]interface{} `json:"result"`
	Error  *utils.Error           `json:"error"`
}

func GetCoinController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetCoinRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res GetCoinResponse
	res.Id = req.Id
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

	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
	one_coin, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, req.Params[0].CoinId)
	if err != nil {
		fmt.Println(err.Error())
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	one_rec := make(map[string]interface{})
	one_rec["coinid"] = one_coin.Coinid
	one_rec["coinsymbol"] = one_coin.Coinsymbol
	one_rec["state"] = one_coin.State
	if session_value.Role == 0 {
		one_rec["ip"] = one_coin.Ip
		one_rec["rpcport"] = one_coin.Rpcport
		one_rec["rpcuser"] = one_coin.Rpcuser
		one_rec["rpcpass"] = one_coin.Rpcpass
	} else {
		one_rec["ip"] = ""
		one_rec["rpcport"] = ""
		one_rec["rpcuser"] = ""
		one_rec["rpcpass"] = ""
	}

	res.Result = one_rec
	ctx.JSON(res)
	return
}


type RefreshCoinsParam struct {
	SessionId  string `json:"sessionid"`
}

type RefreshCoinsRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []RefreshCoinsParam `json:"params"`
}

type RefreshCoinsResponse struct {
	Id     int                      	`json:"id"`
	Result []map[string]interface{}     `json:"result"`
	Error  *utils.Error             	`json:"error"`
}

func RefreshCoinsController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req RefreshCoinsRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res RefreshCoinsResponse
	res.Id = req.Id
	res.Result = nil
	res.Error = nil
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

	// web socket request
	var reqWS ListCoinsRequestWS
	reqWS.Id = 1
	reqWS.JsonRpc = "2.0"
	reqWS.Method = "list_coins"
	msgBytes, err := json.Marshal(reqWS)
	if err != nil {
		res.Error = utils.MakeError(900010)
		ctx.JSON(res)
		return
	}
	fmt.Println("RefreshCoinsController() WSRequest list_coins")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "coin", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		res.Error = utils.MakeError(900020, err.Error())
		ctx.JSON(res)
		return
	}
	var resWS ListCoinsResponseWS
	err = json.Unmarshal([]byte(retString), &resWS)
	if err != nil {
		res.Error = utils.MakeError(900011)
		ctx.JSON(res)
		return
	}
	if resWS.Error != nil {
		res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
		ctx.JSON(res)
		return
	}
	// update local coin config
	for _, coinResult := range resWS.Result {
		isFound, coinCfg, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin2(dbSession, coinResult.CoinId)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.CoinConfigMgr.TableName, "query", "get coin")
			ctx.JSON(res)
			return
		}
		if !isFound {
			err = model.GlobalDBMgr.CoinConfigMgr.InsertCoinWithCoinId(dbSession, coinCfg.Coinid, coinCfg.Coinsymbol, "", 0,
				"", "", coinCfg.State)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300001, model.GlobalDBMgr.CoinConfigMgr.TableName, "insert", "insert coin")
				ctx.JSON(res)
				return
			}
		} else {
			if coinResult.CoinSymbol != coinCfg.Coinsymbol {
				dbSession.Rollback()
				res.Error = utils.MakeError(900040)
				ctx.JSON(res)
				return
			}
			if coinResult.State != coinCfg.State {
				err = model.GlobalDBMgr.CoinConfigMgr.UpdateCoinState(dbSession, coinResult.CoinId, coinResult.State)
				if err != nil {
					dbSession.Rollback()
					res.Error = utils.MakeError(300001, model.GlobalDBMgr.CoinConfigMgr.TableName, "update", "update coin state")
					ctx.JSON(res)
					return
				}
			}
		}
	}

	all_coins, err := model.GlobalDBMgr.CoinConfigMgr.ListCoins(dbSession)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}
	res_data := make([]map[string]interface{}, 0)
	for _, one_coin := range all_coins {
		one_rec := make(map[string]interface{})
		one_rec["coinid"] = one_coin.Coinid
		one_rec["coinsymbol"] = one_coin.Coinsymbol
		one_rec["state"] = one_coin.State
		if session_value.Role == 0 {
			one_rec["ip"] = one_coin.Ip
			one_rec["rpcport"] = one_coin.Rpcport
			one_rec["rpcuser"] = one_coin.Rpcuser
			one_rec["rpcpass"] = one_coin.Rpcpass
		} else {
			one_rec["ip"] = ""
			one_rec["rpcport"] = ""
			one_rec["rpcuser"] = ""
			one_rec["rpcpass"] = ""
		}
		res_data = append(res_data, one_rec)
	}

	// add log
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, session_value.AcctId, 6, coin_convert_log_value(6))
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

	res.Result = res_data
	ctx.JSON(res)
	return
}


func CoinController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "list_coins" {
		ListCoinsController(ctx, jsonRpcBody)
	} else if funcName == "create_coin" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		// CreateCoinController(ctx, jsonRpcBody)
	} else if funcName == "modify_coin" {
		ModifyCoinController(ctx, jsonRpcBody)
	} else if funcName == "change_coin_state" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		// ChangeCoinStateController(ctx, jsonRpcBody)
	} else if funcName == "get_coin" {
		GetCoinController(ctx, jsonRpcBody)
	} else if funcName == "refresh_coins" {
		RefreshCoinsController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}

type ListCoinsRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []interface{} 			`json:"params"`
}

type ListCoinsResultWS struct {
	CoinId   		int    		`json:"coinid"`
	CoinSymbol		string    	`json:"coinsymbol"`
	Ip    			string 		`json:"ip"`
	RpcPort 		int 		`json:"rpcport"`
	RpcUser 		string 		`json:"rpcuser"`
	RpcPass 		string 		`json:"rpcpass"`
	State 			int 		`json:"state"`
}

type ListCoinsResponseWS struct {
	Id     int                      `json:"id"`
	Result []ListCoinsResultWS 		`json:"result"`
	Error  *utils.Error             `json:"error"`
}

