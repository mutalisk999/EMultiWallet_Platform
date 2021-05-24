package controller

import (
	"coin"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"session"
	"strconv"
	"utils"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/go-xorm/xorm"
	"config"
)

type EmptyResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

type RequestBase struct {
	Id      int    `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}

type ListWalletsParam struct {
	SessionId string `json:"sessionid"`
	AcctIds   []int  `json:"acctids"`
	CoinId    []int  `json:"coinid"`
	State     []int  `json:"state"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type ListWalletsRequest struct {
	RequestBase
	Params []ListWalletsParam `json:"params"`
}

type WalletResult struct {
	WalletId   int    `json:"walletid"`
	WalletUUid string `json:"walletuuid"`
	WalletName string `json:"walletname"`
	CoinId     int    `json:"coinid"`
	CoinSymbol string `json:"coinsymbol"`
	Balance    string `json:"balance"`
	FeeBalance string `json:"feebalance"`
	Address    string `json:"address"`
	State      int    `json:"state"`
}

type ListWalletResult struct {
	Total   int64          `json:"total"`
	Wallets []WalletResult `json:"wallets"`
}

type ListWalletResponse struct {
	Id     int               `json:"id"`
	Result *ListWalletResult `json:"result"`
	Error  *utils.Error      `json:"error"`
}

func BasicCheck(dbSession *xorm.Session, ctx iris.Context, sessionid string, role []int, seqtype int, mgmtid int, signature string, orgindata string) (bool, int, *utils.Error) {
	//SessionCheck
	sessionval, exist := session.GlobalSessionMgr.GetSessionValue(sessionid)
	if !exist {
		return false, 0, utils.MakeError(200004)
	}
	//defer session.GlobalSessionMgr.RefreshSessionValue(sessionid)
	//role Check
	if len(role) != 0 {
		inrole := false
		for erole := range role {
			if erole == sessionval.Role {
				inrole = true
				break
			}
		}
		if !inrole {
			return false, sessionval.AcctId, utils.MakeError(400009)
		}
	}
	//mgmtid Check
	if seqtype != 0 {
		vres, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, seqtype, mgmtid)
		if !vres || err != nil {
			return false, sessionval.AcctId, utils.MakeError(400010)
		}
	}
	//Signature Check
	if orgindata != "" {
		err := utils.RsaVerySignWithSha1Hex(orgindata, signature, sessionval.PubKey)
		if err != nil {
			return false, sessionval.AcctId, utils.MakeError(400002)
		}
	}
	return true, sessionval.AcctId, nil
}

func ListWalletController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req ListWalletsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ListWalletResponse
	res.Id = req.Id
	res.Result = nil
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	seval, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	cids := make([]int, 0)
	for _, coinid := range req.Params[0].CoinId {
		_, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, coinid)
		if err != nil {
			if err.Error() != "key not found" {
				res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "list wallets")
				ctx.JSON(res)
				return
			} else {
				res.Error = utils.MakeError(500001, coinid)
				ctx.JSON(res)
				return
			}
		}
		cids = append(cids, coinid)
	}
	var accids = make([]int, 0)
	if seval.Role == 0 {
		accids = append(accids, req.Params[0].AcctIds...)
	} else {
		accids = append(accids, seval.AcctId)
	}
	wallets, total, werr := model.GlobalDBMgr.WalletConfigMgr.ListWallets(dbSession, cids, req.Params[0].State, accids, req.Params[0].Offset, req.Params[0].Limit)
	if werr != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "list wallets")
		ctx.JSON(res)
		return
	}
	res.Result = new(ListWalletResult)
	res.Result.Total = total
	for _, wa := range wallets {
		var walletres WalletResult
		walletres.State = wa.State
		walletres.Address = wa.Address
		walletres.CoinId = wa.Coinid
		coincfg, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, wa.Coinid)
		if coincfg.State != 1 || err != nil {
			continue
		}

		ba, fee_balance, err := coin.GetBalance(coincfg.Coinsymbol, coincfg.Ip, coincfg.Rpcport, coincfg.Rpcuser, coincfg.Rpcpass, wa.Address)
		if err != nil {
			fmt.Println("no money")
			//continue
		}
		walletres.Balance = ba
		walletres.FeeBalance = fee_balance
		walletres.CoinSymbol = coincfg.Coinsymbol
		walletres.WalletId = wa.Walletid
		walletres.WalletUUid = wa.Walletuuid
		walletres.WalletName = wa.Walletname
		res.Result.Wallets = append(res.Result.Wallets, walletres)
	}
	ctx.JSON(res)
	return
}

type CreateWalletParam struct {
	SessionId    string `json:"sessionid"`
	MgmtId       int    `json:"mgmtid"`
	WalletName   string `json:"walletname"`
	CoinId       int    `json:"coinid"`
	KeySigServerId  []int `json:"keysigserverid"`
	NeedKeySigCount int   `json:"needkeysigcount"`
	DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	State        int    `json:"state"`
	Signature    string `json:"signature"`
}

type CreateWalletRequest struct {
	RequestBase
	Params []CreateWalletParam `json:"params"`
}

func CreateWalletController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req CreateWalletRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.NeedSigCount <= 0 {
		res.Error = utils.MakeError(500004)
		ctx.JSON(res)
		return
	}
	orgindata := "create_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + pa.WalletName + "," + strconv.Itoa(pa.CoinId) +
		"," + utils.IntArrayToString(pa.KeySigServerId) + "," + strconv.Itoa(pa.NeedKeySigCount) +
		"," + pa.DestAddress + "," + strconv.Itoa(pa.NeedSigCount) + "," + pa.Fee + "," + pa.GasPrice + "," + pa.GasLimit + "," +
		utils.IntArrayToString(pa.SigUserId) + "," + strconv.Itoa(pa.State)

	checkres, acctid, errres := BasicCheck(dbSession, ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		dbSession.Rollback()
		res.Error = errres
		ctx.JSON(res)
		return
	}

	var serverKeyInfos []model.ServerKeyInfo
	for _, serverId := range pa.KeySigServerId {
		isFound, serverInfo, err := model.GlobalDBMgr.ServerInfoMgr.GetServerInfoById(dbSession, serverId)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get server by id")
			ctx.JSON(res)
			return
		} else if !isFound {
			dbSession.Rollback()
			res.Error = utils.MakeError(200104, serverId)
			ctx.JSON(res)
		}
		var serverKeyInfo model.ServerKeyInfo
		serverKeyInfo.ServerId = serverId
		serverKeyInfo.StartIndex = serverInfo.Serverstartindex
		serverKeyInfos = append(serverKeyInfos, serverKeyInfo)
	}

	// choose unused server pubkey
	serverKeyInfos, err = model.GlobalDBMgr.PubKeyPoolMgr.GetUnusedServerKeys(dbSession, serverKeyInfos)
	if err != nil {
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(200100)
			ctx.JSON(res)
			return
		}
	}

	coinConfig, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, pa.CoinId)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(dbSession, serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get coin by id")
		ctx.JSON(res)
		return
	}

	var pubKeySlice []string
	for _, serverKeyInfo := range serverKeyInfos {
		pubKeySlice = append(pubKeySlice, serverKeyInfo.PubKey)
	}
	address, err := coin.GetMultiSignAddressByPubKeys(pa.NeedKeySigCount, pubKeySlice, coinConfig.Coinsymbol)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(200101)
		ctx.JSON(res)
		return
	}

	if pa.DestAddress != "" {
		dstAddrList := strings.Split(pa.DestAddress, ",")
		for _, dstAddress := range dstAddrList {
			valid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, dstAddress)
			if err != nil {
				// rollback when error
				//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
				dbSession.Rollback()
				res.Error = utils.MakeError(800000, err.Error())
				ctx.JSON(res)
				return
			}
			if !valid {
				// rollback when error
				//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
				dbSession.Rollback()
				res.Error = utils.MakeError(500003, dstAddress)
				ctx.JSON(res)
				return
			}
		}
	}

	err = coin.ImportAddress(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, address)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(800000, err.Error())
		ctx.JSON(res)
		return
	}

	// 生成钱包UUID
	u, _ := uuid.NewV4()
	walletUUid := u.String()

	var serverKeyStrSlice []string
	for _, serverKeyInfo := range serverKeyInfos {
		serverKeyStrSlice = append(serverKeyStrSlice, fmt.Sprintf("%d:%d", serverKeyInfo.ServerId, serverKeyInfo.KeyIndex))
	}
	serverKeysStr := strings.Join(serverKeyStrSlice, ",")

	err = model.GlobalDBMgr.WalletConfigMgr.InsertWallet(dbSession, pa.CoinId, walletUUid, pa.WalletName, serverKeysStr,
		config.LocalServerId, len(serverKeyStrSlice),
		pa.NeedKeySigCount, address, pa.DestAddress, pa.NeedSigCount, pa.Fee, pa.GasPrice, pa.GasLimit, 0)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Insert", "create wallets")
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletByName(dbSession, pa.WalletName)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get wallet by name")
		ctx.JSON(res)
		return
	}
	for _,uid := range pa.SigUserId {
		err = model.GlobalDBMgr.AcctWalletRelationMgr.InsertRelation(dbSession, uid, wa.Walletid)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Insert", "InsertRelation")
			ctx.JSON(res)
			return
		}
	}
	logmsg := "创建钱包,钱包名:" + pa.WalletName + ",CoinId:" + strconv.Itoa(pa.CoinId) + ",Address:" + address +
		",DestAddress:" + pa.DestAddress +
		",NeedKeySigCount:" + strconv.Itoa(pa.NeedKeySigCount) + ",Fee:" + pa.Fee + ",GasPrice:" + pa.GasPrice + ",GasLimit:" + pa.GasLimit +
		",SigUserId:" + utils.IntArrayToString(pa.SigUserId) + ",State:" + strconv.Itoa(pa.State)

	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acctid, 4, logmsg)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.OperationLogMgr.TableName, "insert", "new operator log")
		ctx.JSON(res)
		return
	}

	// web socket request
	var reqWS CreateWalletRequestWS
	reqWS.Id = 1
	reqWS.JsonRpc = "2.0"
	reqWS.Method = "wallet_create"
	var createWalletParamWS CreateWalletParamWS
	createWalletParamWS.CoinId = pa.CoinId
	createWalletParamWS.NeedSigCount = pa.NeedKeySigCount
	createWalletParamWS.TotalCount = len(pa.KeySigServerId)
	createWalletParamWS.WalletName = pa.WalletName
	createWalletParamWS.WalletUuid = walletUUid
	createWalletParamWS.Address = address
	createWalletParamWS.DestAddress = pa.DestAddress
	createWalletParamWS.Fee = pa.Fee
	createWalletParamWS.GasPrice = pa.GasPrice
	createWalletParamWS.GasLimit = pa.GasLimit

	for _, serverKeyInfo := range serverKeyInfos {
		var serverKey ServerKeyParamWS
		serverKey.ServerId = serverKeyInfo.ServerId
		serverKey.KeyIndex = serverKeyInfo.KeyIndex
		createWalletParamWS.KeyDetail = append(createWalletParamWS.KeyDetail, serverKey)
	}

	reqWS.Params = append(reqWS.Params, createWalletParamWS)
	msgBytes, err := json.Marshal(reqWS)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(900010)
		ctx.JSON(res)
		return
	}
	fmt.Println("CreateWalletController() WSRequest wallet_create")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "wallet", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(900020, err.Error())
		ctx.JSON(res)
		return
	}
	var resWS CreateWalletResponseWS
	err = json.Unmarshal([]byte(retString), &resWS)
	if err != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(900011)
		ctx.JSON(res)
		return
	}
	if resWS.Error != nil {
		// rollback when error
		//model.GlobalDBMgr.PubKeyPoolMgr.RollBackUsedServerKeys(serverKeyInfos)
		dbSession.Rollback()
		res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
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

type GetWalletParam struct {
	SessionId string `json:"sessionid"`
	WalletId  int    `json:"walletid"`
}

type GetWalletsRequest struct {
	RequestBase
	Params []GetWalletParam `json:"params"`
}

type GetWalletResult struct {
	WalletId     int    `json:"walletid"`
	WalletUuid   string `json:"walletuuid"`
	WalletName   string `json:"walletname"`
	CoinId       int    `json:"coinid"`
	ServerKeys   string `json:"serverkeys"`
	CreateServer int    `json:"createserver"`
	KeyCount 	 int    `json:"keycount"`
	NeedKeySigCount int `json:"needkeysigcount"`
	Address      string `json:"address"`
	DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	State        int    `json:"state"`
}

type GetWalletResponse struct {
	Id     int             `json:"id"`
	Result GetWalletResult `json:"result"`
	Error  *utils.Error    `json:"error"`
}

func GetWalletController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetWalletsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	var res GetWalletResponse
	res.Id = req.Id
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		ctx.JSON(res)
		return
	}

	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	// just check, no database transaction to rollback
	checkres, _, errres := BasicCheck(dbSession, ctx, pa.SessionId, []int{0, 1}, 0, 0, "", "")
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(dbSession, pa.WalletId)
	if err != nil {
		if err.Error() == "no find wallet" {
			res.Error = utils.MakeError(500000, pa.WalletId)
		} else {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "get wallet by id")
		}
	} else {
		res.Result.WalletId = wa.Walletid
		res.Result.WalletUuid = wa.Walletuuid
		res.Result.WalletName = wa.Walletname
		res.Result.CoinId = wa.Coinid
		res.Result.ServerKeys = wa.Serverkeys
		res.Result.CreateServer = wa.Createserver
		res.Result.KeyCount = wa.Keycount
		res.Result.NeedKeySigCount = wa.Needkeysigcount
		res.Result.Address = wa.Address
		res.Result.DestAddress = wa.Destaddress
		res.Result.NeedSigCount = wa.Needsigcount
		res.Result.Fee = wa.Fee
		res.Result.GasPrice = wa.Gasprice
		res.Result.GasLimit = wa.Gaslimit
		res.Result.State = wa.State
		uids := make([]int, 0)
		relations, err := model.GlobalDBMgr.AcctWalletRelationMgr.GetRelationsByWalletId(dbSession, wa.Walletid)
		if err != nil {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "query", "get relation by wallet id")
			ctx.JSON(res)
			return
		}
		for _, rel := range relations {
			uids = append(uids, rel.Acctid)
		}
		res.Result.SigUserId = uids
	}
	ctx.JSON(res)
	return
}

type ModifyWalletParam struct {
	SessionId    string `json:"sessionid"`
	MgmtId       int    `json:"mgmtid"`
	Walletid     int    `json:"walletid"`
	//WalletName   string `json:"walletname"`
	//DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	//Fee          string `json:"fee"`
	//GasPrice     string `json:"gasprice"`
	//GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	//State        int    `json:"state"`
	Signature    string `json:"signature"`
}

type ModifyWalletRequest struct {
	RequestBase
	Params []ModifyWalletParam `json:"params"`
}

func ModifyWalletController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ModifyWalletRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.NeedSigCount <= 0 {
		res.Error = utils.MakeError(500004)
		ctx.JSON(res)
		return
	}
	//if pa.State != 0 && pa.State != 1 && pa.State != 2 {
	//	res.Error = utils.MakeError(200001)
	//	ctx.JSON(res)
	//	return
	//}
	//orgindata := "modify_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid) + "," + pa.WalletName + "," +
	//	pa.DestAddress + "," + strconv.Itoa(pa.NeedSigCount) + "," + pa.Fee + "," + pa.GasPrice + "," + pa.GasLimit + "," +
	//	utils.IntArrayToString(pa.SigUserId) + "," + strconv.Itoa(pa.State)

	orgindata := "modify_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid) + "," +
		strconv.Itoa(pa.NeedSigCount) + "," + utils.IntArrayToString(pa.SigUserId)

	checkres, acctid, errres := BasicCheck(dbSession, ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		dbSession.Rollback()
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(dbSession, pa.Walletid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(500000, pa.Walletid)
		ctx.JSON(res)
		return
	}
	if wa.State == 3 {
		dbSession.Rollback()
		res.Error = utils.MakeError(500005, pa.Walletid)
		ctx.JSON(res)
		return
	}
	//coinConfig, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(wa.Coinid)
	_, err = model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, wa.Coinid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get coin by id")
		ctx.JSON(res)
		return
	}

	//if pa.DestAddress != "" {
	//	dstAddrList := strings.Split(pa.DestAddress, ",")
	//	for _, dstAddress := range dstAddrList {
	//		valid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, dstAddress)
	//		if err != nil {
	//			res.Error = utils.MakeError(800000, err.Error())
	//			ctx.JSON(res)
	//			return
	//		}
	//		if !valid {
	//			res.Error = utils.MakeError(500003, dstAddress)
	//			ctx.JSON(res)
	//			return
	//		}
	//	}
	//}

	//err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(pa.Walletid, pa.WalletName, pa.DestAddress, pa.NeedSigCount,
	//	pa.Fee, pa.GasPrice, pa.GasLimit, pa.State)
	err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(dbSession, pa.Walletid, pa.NeedSigCount)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "update wallet")
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctWalletRelationMgr.DeleteRelationByWalletId(dbSession, pa.Walletid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Delete", "Delete Relation by wallet id")
		ctx.JSON(res)
		return
	}
	for _, uid := range pa.SigUserId {
		err = model.GlobalDBMgr.AcctWalletRelationMgr.InsertRelation(dbSession, uid, pa.Walletid)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Insert", "InsertRelation")
			ctx.JSON(res)
			return
		}
	}

	//logmsg := "修改钱包,钱包id:" + strconv.Itoa(pa.Walletid) + ",修改后属性 钱包名:" + pa.WalletName + ",DestAddress:" + pa.DestAddress +
	//	",NeedSigCount:" + strconv.Itoa(pa.NeedSigCount) + ",Fee:" + pa.Fee + ",GasPrice:" + pa.GasPrice + ",GasLimit:" + pa.GasLimit +
	//	",SigUserId:" + utils.IntArrayToString(pa.SigUserId) + ",State:" + strconv.Itoa(pa.State)

	logmsg := "修改钱包,钱包id:" + strconv.Itoa(pa.Walletid) + ",修改后属性 NeedSigCount:" + strconv.Itoa(pa.NeedSigCount) +
		",SigUserId:" + utils.IntArrayToString(pa.SigUserId)

	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acctid, 4, logmsg)
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

	ctx.JSON(res)
	return
}

type DeleteWalletParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	Walletid  int    `json:"walletid"`
	Signature string `json:"signature"`
}

type DeleteWalletRequest struct {
	RequestBase
	Params []DeleteWalletParam `json:"params"`
}

func DeleteWalletController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req DeleteWalletRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	orgindata := "delete_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid)
	checkres, acctid, errres := BasicCheck(dbSession, ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		dbSession.Rollback()
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wal, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(dbSession, pa.Walletid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "select", "select wallet by id")
		ctx.JSON(res)
		return
	}
	wal.State = 3
	//err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(wal.Walletid, wal.Walletname, wal.Destaddress, wal.Needsigcount, wal.Fee, wal.Gasprice, wal.Gaslimit, wal.State)
	err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(dbSession, wal.Walletid, wal.Needsigcount)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "set wallet to delete")
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctWalletRelationMgr.DeleteRelationByWalletId(dbSession, pa.Walletid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "delete", "delete relation")
		ctx.JSON(res)
		return
	}

	logmsg := "删除钱包,钱包id:" + strconv.Itoa(pa.Walletid)
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acctid, 4, logmsg)
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

	ctx.JSON(res)
	return
}

type ChangeWalletStateParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	Walletid  int    `json:"walletid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}

type ChangeWalletStateRequest struct {
	RequestBase
	Params []ChangeWalletStateParam `json:"params"`
}

func ChangeWalletStateController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ChangeWalletStateRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.State != 0 && pa.State != 1 && pa.State != 2 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	orgindata := "change_wallet_state," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid) + "," + strconv.Itoa(pa.State)
	checkres, acctid, errres := BasicCheck(dbSession, ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		dbSession.Rollback()
		res.Error = errres
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.WalletConfigMgr.ChangeWalletState(dbSession, pa.Walletid, pa.State)
	if err != nil {
		if err.Error() == "no find wallet" {
			dbSession.Rollback()
			res.Error = utils.MakeError(500000, pa.Walletid)
			ctx.JSON(res)
			return
		} else {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "change wallet state")
			ctx.JSON(res)
			return
		}
	}

	logmsg := "更改钱包状态,钱包id:" + strconv.Itoa(pa.Walletid) + ",新状态:" + strconv.Itoa(pa.State)
	_, err = model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(dbSession, acctid, 4, logmsg)
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

	ctx.JSON(res)
	return
}

func WalletController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	var res utils.JsonRpcResponse
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		ctx.JSON(res)
		return
	}

	if funcName == "list_wallets" {
		ListWalletController(ctx, jsonRpcBody)
	} else if funcName == "create_wallet" {
		CreateWalletController(ctx, jsonRpcBody)
	} else if funcName == "get_wallet" {
		GetWalletController(ctx, jsonRpcBody)
	} else if funcName == "modify_wallet" {
		ModifyWalletController(ctx, jsonRpcBody)
	} else if funcName == "delete_wallet" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		//DeleteWalletController(ctx, jsonRpcBody)
	} else if funcName == "change_wallet_state" {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
		//ChangeWalletStateController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}

type ServerKeyParamWS struct {
	ServerId   int   `json:"serverid"`
	KeyIndex   int   `json:"keyindex"`
}

type CreateWalletParamWS struct {
	CoinId       int    `json:"coinid"`
	NeedSigCount int    `json:"needsigcount"`
	TotalCount 	 int    `json:"totalcount"`
	WalletName   string `json:"walletname"`
	WalletUuid   string `json:"walletuuid"`
	Address  	 string `json:"address"`
	DestAddress  string `json:"destaddress"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	KeyDetail    []ServerKeyParamWS  `json:"keydetail"`
}

type CreateWalletRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params []CreateWalletParamWS    `json:"params"`
}

type CreateWalletResponseWS struct {
	Id     int                      `json:"id"`
	Result bool      	     		`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type QueryWalletParamWS struct {
	CoinId       []int    			`json:"coinid"`
	WalletUuids  []string    		`json:"walletuuids"`
}

type QueryWalletRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params []QueryWalletParamWS     `json:"params"`
}

type QueryWalletResultWS struct {
	WalletId      int               `json:"walletid"`
	WalletUuid    string            `json:"walletuuid"`
	CoinId        int               `json:"coinid"`
	NeedSigCount  int               `json:"needsigcount"`
	TotalCount    int               `json:"totalcount"`
	KeyDetail     []ServerKeyParamWS   `json:"keydetail"`
	WalletName    string            `json:"walletname"`
	Address       string            `json:"address"`
	DestAddress       string        `json:"destaddress"`
	CreateServerId    int           `json:"createserverid"`
	Fee           string 			`json:"fee"`
	GasPrice      string 			`json:"gasprice"`
	GasLimit      string 			`json:"gaslimit"`
	State		  int           	`json:"state"`
}

type QueryWalletResponseWS struct {
	Id     int                   `json:"id"`
	Result []QueryWalletResultWS `json:"result"`
	Error  *utils.Error          `json:"error"`
}


