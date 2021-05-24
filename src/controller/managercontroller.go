package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"utils"
	"config"
)
//
////change_account_state
//
//type BackupDbParam struct {
//	SessionId string `json:"sessionid"`
//	Mgmtid    int    `json:"mgmtid"`
//
//	Signature string `json:"signature"`
//}
//
//type BackupDbRequest struct {
//	Id      int             `json:"id"`
//	JsonRpc string          `json:"jsonrpc"`
//	Method  string          `json:"method"`
//	Params  []BackupDbParam `json:"params"`
//}
//
//type BackupDbResponse struct {
//	Id     int          `json:"id"`
//	Result bool         `json:"result"`
//	Error  *utils.Error `json:"error"`
//}
//
//func BackupDbController(ctx iris.Context, jsonRpcBody []byte) {
//	var req BackupDbRequest
//	err := json.Unmarshal(jsonRpcBody, &req)
//	if err != nil {
//		utils.SetInternalError(ctx, err.Error())
//		return
//	}
//	var res BackupDbResponse
//	res.Id = req.Id
//	res.Result = false
//	if len(req.Params) != 1 {
//		res.Error = utils.MakeError(200001)
//		ctx.JSON(res)
//		return
//	}
//	session_value, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
//	if !exist || session_value.Role != 0 {
//		res.Error = utils.MakeError(200004)
//		ctx.JSON(res)
//		return
//	}
//	//session.GlobalSessionMgr.RefreshSessionValue(req.Params[0].SessionId)
//	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(7, req.Params[0].Mgmtid)
//	if !verify || err != nil {
//		res.Error = utils.MakeError(400005)
//		ctx.JSON(res)
//		return
//	}
//	//check signature
//	sig_origin_data := "backup_db," + req.Params[0].SessionId + "," + strconv.Itoa(req.Params[0].Mgmtid)
//	err = utils.RsaVerySignWithSha1Hex(sig_origin_data, req.Params[0].Signature, session_value.PubKey)
//	if err != nil {
//		res.Error = utils.MakeError(400002)
//		ctx.JSON(res)
//		return
//	}
//	model.GlobalDBMgr.DBEngine.Close()
//	err = ctx.SendFile(config.GlobalConfig.DbConfig.DbPath, "ewallet.db")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	err = model.InitDB(config.GlobalConfig.DbConfig.DbType, config.GlobalConfig.DbConfig.DbPath)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	return
//
//}

//server_init_status

type ServerInitStatusRequest struct {
	Id      int             `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  interface{} 	`json:"params"`
}

type ServerInitStatusResponse struct {
	Id     int          `json:"id"`
	Result bool         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ServerInitStatusController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req ServerInitStatusRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ServerInitStatusResponse
	res.Id = req.Id
	count ,err :=model.GlobalDBMgr.AcctConfigMgr.GetAccountCount(dbSession)
	if err!=nil{
		fmt.Println(err.Error())
		res.Error  = utils.MakeError(400006)
		ctx.JSON(res)
		return
	}

	res.Result = count>0
	ctx.JSON(res)
	return
}

//server_list_support_coins

type ServerListSupportCoinsRequest struct {
	Id      int             `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  interface{} 	`json:"params"`
}

type ServerListSupportCoinsResult struct {
	CoinName     string          `json:"coinname"`
	CoinSymbol     string          `json:"coinsymbol"`
	Precision	int		`json:"precision"`
	ConfirmCount	int		`json:"confirmcount"`
	IsErc20		bool		`json:"isErc20"`
	ContractAddress	string		`json:"contractaddress"`
	IsOmni		bool		`json:"isomni"`
	OmniPropertyId	int		`json:"omnipropertyid"`
}

type ServerListSupportCoinsResponse struct {
	Id     int          `json:"id"`
	Result []ServerListSupportCoinsResult         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ServerListSupportCoinsController(ctx iris.Context, jsonRpcBody []byte) {
	var req ServerListSupportCoinsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ServerListSupportCoinsResponse
	res.Id = req.Id
	for _,v := range config.GlobalSupportCoinMgr{

		res.Result = append(res.Result,ServerListSupportCoinsResult{v.CoinName,v.CoinSymbol,v.Precision,v.ConfirmCount,v.IsErc20,v.ContractAddress,v.IsOmni,v.OmniPropertyId})
	}

	ctx.JSON(res)
	return
}

func ManagerController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse

	if funcName == "server_init_status"{
		ServerInitStatusController(ctx, jsonRpcBody)
	}else if funcName == "server_list_support_coins"{
		ServerListSupportCoinsController(ctx, jsonRpcBody)
	}else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}
