package controller

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"session"
	"strconv"
	"strings"
	"utils"
	"coin"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"encoding/hex"
	"github.com/satori/go.uuid"
	"github.com/mutalisk999/bitcoin-lib/src/blob"
	"bytes"
	"io"
	"github.com/mutalisk999/bitcoin-lib/src/transaction"
	"github.com/go-xorm/xorm"
)

const (
	LogFormatTypeTrxTransfer = 1
	LogFormatTypeTrxConfirm  = 2
	LogFormatTypeTrxRevoke   = 3
)

func CreateTransactionString(trxUuid string, walletId int, coinId int, contractAddr string, acctId int, serverId int,
	fromAddr string, toDetails string, needConfirm int, fee string, gasPrice string, gasLimit string) (string) {
	walletIdStr := strconv.Itoa(walletId)
	coinIdStr := strconv.Itoa(coinId)
	acctIdStr := strconv.Itoa(acctId)
	serverIdStr := strconv.Itoa(serverId)
	needConfirmStr := strconv.Itoa(needConfirm)

	trxStr := strings.Join([]string{trxUuid, walletIdStr, coinIdStr, contractAddr, acctIdStr, serverIdStr, fromAddr, toDetails,
		needConfirmStr, fee, gasPrice, gasLimit}, ",")
	return trxStr
}

func GetTransactionLogFormat(fmtType int) string {
	if fmtType == LogFormatTypeTrxTransfer {
		return "用户[%s]发起了一笔从钱包[%s]的转账,交易ID:[%s],结果:[%s]"
	} else if fmtType == LogFormatTypeTrxConfirm {
		return "用户[%s]对交易[%s]进行了确认操作,确认结果:[%s]"
	} else if fmtType == LogFormatTypeTrxRevoke {
		return "用户[%s]对交易[%s]进行了撤销操作,撤销结果:[%s]"
	}
	return ""
}

type GetWalletTrxsParam struct {
	SessionId string    `json:"sessionid"`
	WalletId  []int     `json:"walletid"`
	CoinId    []int     `json:"coinid"`
	ServerId  []int     `json:"serverid"`
	AcctId    []int     `json:"acctid"`
	TrxTime   [2]string `json:"trxtime"`
	State     []int     `json:"state"`
	OffSet    int       `json:"offset"`
	Limit     int       `json:"limit"`
}

type GetWalletTrxsRequest struct {
	Id      int                  `json:"id"`
	JsonRpc string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  []GetWalletTrxsParam `json:"params"`
}

type WalletTrx struct {
	TrxId         int    `json:"trxid"`
	TrxUuid       string `json:"trxuuid"`
	RawTrxId      string `json:"rawtrxid"`
	WalletId      int    `json:"walletid"`
	CoinId        int    `json:"coinid"`
	ContractAddr  string `json:"contractaddr"`
	AcctId        int    `json:"acctid"`
	ServerId      int    `json:"serverid"`
	FromAddr      string `json:"fromaddr"`
	ToDetails     string `json:"todetails"`
	FeeCost       string `json:"feecost"`
	TrxTime       string `json:"trxtime"`
	NeedConfirm   int    `json:"needconfirm"`
	Confirmed     int    `json:"confirmed"`
	AcctConfirmed string `json:"acctconfirmed"`
	SignedServerIds string `json:"signedserverids"`
	State         int    `json:"state"`
}

type GetWalletTrxsResult struct {
	Total int         `json:"total"`
	Trxs  []WalletTrx `json:"trxs"`
}

type GetWalletTrxsResponse struct {
	Id     int                  `json:"id"`
	Result *GetWalletTrxsResult `json:"result"`
	Error  *utils.Error         `json:"error"`
}

func GetTransactionController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	var req GetWalletTrxsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetWalletTrxsResponse
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

	if sessionValue.Role != 0 && req.Params[0].WalletId != nil && len(req.Params[0].WalletId) != 0 {
		// 检查Wallet
		walletMgr := model.GlobalDBMgr.WalletConfigMgr
		walletConfigs, err := walletMgr.GetWalletsByIds(dbSession, req.Params[0].WalletId)
		if err != nil {
			res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
			ctx.JSON(res)
			return
		}
		for _, walletConfig := range walletConfigs {
			if walletConfig.State != 1 {
				res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
				ctx.JSON(res)
				return
			}
		}
	}

	if sessionValue.Role != 0 && req.Params[0].CoinId != nil && len(req.Params[0].CoinId) != 0 {
		// 检查Coin
		coinMgr := model.GlobalDBMgr.CoinConfigMgr
		coinConfigs, err := coinMgr.GetCoins(dbSession, req.Params[0].CoinId)
		if err != nil {
			res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
			ctx.JSON(res)
			return
		}
		for _, coinConfig := range coinConfigs {
			if coinConfig.State != 1 {
				res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
				ctx.JSON(res)
				return
			}
		}
	}

	if sessionValue.Role != 0 && req.Params[0].AcctId != nil && len(req.Params[0].AcctId) != 0 {
		// 检查Account
		acctMgr := model.GlobalDBMgr.AcctConfigMgr
		acctConfigs, err := acctMgr.GetAccountsByIds(dbSession, req.Params[0].AcctId)
		if err != nil {
			res.Error = utils.MakeError(300001, acctMgr.TableName, "query", "get account config")
			ctx.JSON(res)
			return
		}
		for _, acctConfig := range acctConfigs {
			if acctConfig.State != 1 {
				res.Error = utils.MakeError(200007, "account", acctConfig.Acctid)
				ctx.JSON(res)
				return
			}
		}
	}

	walletIdsArgs := req.Params[0].WalletId
	if sessionValue.Role != 0 {
		relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
		relations, err := relationMgr.GetRelationsByAcctId(dbSession, sessionValue.AcctId)
		if err != nil {
			res.Error = utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
			ctx.JSON(res)
			return
		}
		for _, walletId := range walletIdsArgs {
			hasRelation := false
			for _, relation := range relations {
				if walletId == relation.Walletid {
					hasRelation = true
					break
				}
			}
			if hasRelation == false {
				res.Error = utils.MakeError(200008)
				ctx.JSON(res)
				return
			}
		}
		if walletIdsArgs == nil {
			walletIdsArgs = make([]int, 0)
		}
		if len(walletIdsArgs) == 0 {
			for _, relation := range relations {
				walletIdsArgs = append(walletIdsArgs, relation.Walletid)
			}
		}
		// Acct不拥有任何钱包
		if len(walletIdsArgs) == 0 {
			ctx.JSON(res)
			return
		}
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	totalCount, transactions, err := trxMgr.GetTransactions(dbSession, walletIdsArgs, req.Params[0].CoinId,
		req.Params[0].ServerId, req.Params[0].AcctId,
		req.Params[0].State, req.Params[0].TrxTime, req.Params[0].OffSet, req.Params[0].Limit)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	res.Result = new(GetWalletTrxsResult)
	res.Result.Total = totalCount
	res.Result.Trxs = make([]WalletTrx, len(transactions), len(transactions))
	for i, trx := range transactions {
		res.Result.Trxs[i] = WalletTrx{
			trx.Trxid, trx.Trxuuid,trx.Rawtrxid, trx.Walletid, trx.Coinid, trx.Contractaddr,
			trx.Acctid, trx.Serverid,
			trx.Fromaddr, trx.Todetails, trx.Feecost,utils.TimeToFormatString(trx.Trxtime),
			trx.Needconfirm, trx.Confirmed, trx.Acctconfirmed, trx.Signedserverids,trx.State}
	}
	ctx.JSON(res)
}

func CheckAccountAvailable(acctId int) *utils.Error {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	// 获取Account的配置信息
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(dbSession, acctId)
	if err != nil {
		return utils.MakeError(300001, acctMgr.TableName, "query", "get account config")
	}
	if acctConfig.State != 1 {
		return utils.MakeError(200007, "account", acctConfig.Acctid)
	}
	return nil
}

func CheckRelationAvailable(acctId int, fromWalletId int) *utils.Error {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	// 获取Account关联的钱包信息
	relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
	relations, err := relationMgr.GetRelationsByAcctId(dbSession, acctId)
	if err != nil {
		return utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
	}
	hasRelation := false
	for _, relation := range relations {
		if fromWalletId == relation.Walletid {
			hasRelation = true
			break
		}
	}
	if hasRelation == false {
		return utils.MakeError(200008)
	}
	return nil
}

func CheckDestAddress(toAddr string, dstAddrConfig string, walletId int) *utils.Error {
	// 如果Destaddress非空  判断入账地址是否存在于Destaddress地址中
	if dstAddrConfig != "" || len(dstAddrConfig) != 0 {
		dstAddrList := strings.Split(dstAddrConfig, ",")
		inDstAddrs := false
		for _, dstAddr := range dstAddrList {
			if toAddr == dstAddr {
				inDstAddrs = true
				break
			}
		}
		if !inDstAddrs {
			return utils.MakeError(200010, toAddr, walletId)
		}
	}
	return nil
}

func TransferLog(dbSession *xorm.Session, isSuccQuit bool, acctId int, walletId int, trxId *int) error {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(dbSession, acctId)
	if err != nil {
		return err
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	trxIdStr := ""
	if trxId != nil {
		trxIdStr = strconv.Itoa(*trxId)
	}

	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(dbSession, walletId)
	if err != nil {
		return err
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxTransfer), acctConfig.Realname,
		walletConfig.Walletname, trxIdStr, resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(dbSession, acctId, 5, logContent)
	if err != nil {
		return err
	}
	return nil
}

type CreateTrxParam struct {
	SessionId    string `json:"sessionid"`
	OperateId    int    `json:"operateid"`
	FromWalletId int    `json:"fromwalletid"`
	ToAddr       string `json:"toaddr"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	Signature    string `json:"signature"`
}

type CreateTrxRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []CreateTrxParam `json:"params"`
}

type CreateTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func CreateTrxController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req CreateTrxRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res CreateTrxResponse
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
	if sessionValue.Role != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200009)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 4, req.Params[0].OperateId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "transfer"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	fromWalletIdStr := strconv.Itoa(req.Params[0].FromWalletId)
	toAddrStr := req.Params[0].ToAddr
	amountStr := req.Params[0].Amount
	feeStr := req.Params[0].Fee
	gasPriceStr := req.Params[0].GasPrice
	gasLimitStr := req.Params[0].GasLimit
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, fromWalletIdStr, toAddrStr, amountStr, feeStr, gasPriceStr, gasLimitStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	// 检查Account的有效性
	uErr := CheckAccountAvailable(sessionValue.AcctId)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 检查交易出账的Wallet是否属于Account关联的Wallet
	uErr = CheckRelationAvailable(sessionValue.AcctId, req.Params[0].FromWalletId)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 获取Wallet的配置信息
	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(dbSession, req.Params[0].FromWalletId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if walletConfig.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	if walletConfig.Needsigcount == 0 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200019, walletConfig.Walletid)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// parse serverkeys
	serverIdKeyPairMap := make(map[int]int)
	var serverKeyInfos []model.ServerKeyInfo
	serverIdKeyPairStrs := strings.Split(walletConfig.Serverkeys, ",")
	for _, serverIdKeyPairStr := range serverIdKeyPairStrs {
		tmpStrs := strings.Split(serverIdKeyPairStr, ":")
		if len(tmpStrs) != 2 {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}
		serverId, err := strconv.Atoi(tmpStrs[0])
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}
		keyIndex, err := strconv.Atoi(tmpStrs[1])
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}
		serverIdKeyPairMap[serverId] = keyIndex

		var serverKeyInfo model.ServerKeyInfo
		serverKeyInfo.ServerId = serverId
		serverKeyInfo.KeyIndex = keyIndex
		serverKeyInfos = append(serverKeyInfos, serverKeyInfo)
	}

	KeyIndex, ok := serverIdKeyPairMap[config.LocalServerId]
	if !ok {
		dbSession.Rollback()
		res.Error = utils.MakeError(200102, walletConfig.Walletid)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 获取钱包对应PubKey信息
	var pubKeyStrSlice []string
	for _, serverKeyInfo := range serverKeyInfos {
		keyPoolMgr := model.GlobalDBMgr.PubKeyPoolMgr
		pubkeyStr, err := keyPoolMgr.QueryPubKeyByKeyIndex(dbSession, serverKeyInfo.ServerId, serverKeyInfo.KeyIndex)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get key pool config")
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}
		pubKeyStrSlice = append(pubKeyStrSlice, pubkeyStr)
	}

	// 获取Coin配置信息
	coinMgr := model.GlobalDBMgr.CoinConfigMgr
	coinConfig, err := coinMgr.GetCoin(dbSession, walletConfig.Coinid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if coinConfig.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(coinConfig.Coinsymbol) {
		dbSession.Rollback()
		res.Error = utils.MakeError(600001, coinConfig.Coinsymbol)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	coinConfigDetail, _ := config.GlobalSupportCoinMgr[coinConfig.Coinsymbol]

	// 检查to address的合法性
	addrValid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser,
		coinConfig.Rpcpass, req.Params[0].ToAddr)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(800000, err.Error())
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if !addrValid {
		dbSession.Rollback()
		res.Error = utils.MakeError(500003, req.Params[0].ToAddr)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 检查存在Dest Address设置下  交易中to address的有效性
	uErr = CheckDestAddress(req.Params[0].ToAddr, walletConfig.Destaddress, walletConfig.Walletid)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 生成交易UUID
	u, _ := uuid.NewV4()
	trxUuid := u.String()

	// 生成交易防篡改签名
	trxStr := CreateTransactionString(trxUuid, req.Params[0].FromWalletId, walletConfig.Coinid,
		coinConfigDetail.ContractAddress, sessionValue.AcctId, config.LocalServerId, walletConfig.Address,
		fmt.Sprintf("%s~%s",req.Params[0].ToAddr, req.Params[0].Amount), walletConfig.Needsigcount,
		req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit)
	signature, err := coin.CoinSignTrx('2', utility.Sha256([]byte(trxStr)), uint16(KeyIndex))
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(900002)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	signatureHex := hex.EncodeToString(signature)

	trxMgr := model.GlobalDBMgr.TransactionMgr
	trxId, err := trxMgr.NewTransaction(dbSession, trxUuid, req.Params[0].FromWalletId, walletConfig.Coinid, coinConfigDetail.ContractAddress,
		sessionValue.AcctId, config.LocalServerId, walletConfig.Address, fmt.Sprintf("%s~%s",req.Params[0].ToAddr, req.Params[0].Amount),
		walletConfig.Needsigcount, req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit, signatureHex)

	_, trx, err := trxMgr.GetTransactionById(dbSession, trxId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// list unspent
	utxos, err := coin.ListUnSpent(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, walletConfig.Address)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(200103)
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// filter unspent
	utxosFilter := make([]coin.UTXODetail, 0)
	for _, utxo := range utxos {
		count ,err := model.GlobalDBMgr.PendingTransactionMgr.GetPendingTransactionByVinCount(dbSession, coinConfig.Coinid, utxo.TxId, utxo.Vout)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get pending transaction")
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}
		if count == 0 {
			utxosFilter = append(utxosFilter, utxo)
		}
	}

	// verify if balance is enough or not
	// TODO

	if walletConfig.Needsigcount == 1 {
		feeCostStr, trxSignedHex, err := coin.CreateMultiSignedTransaction(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, uint16(KeyIndex),
			walletConfig.Needkeysigcount, pubKeyStrSlice, utxosFilter, walletConfig.Address, req.Params[0].ToAddr, req.Params[0].Amount,
			req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit)

		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(800000, err.Error())
			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}

		if feeCostStr != ""{
			trx.Feecost = feeCostStr
		}
		trx.State = 1
		trx.Signedserverids = strconv.Itoa(config.LocalServerId)
		trx.Signedtrxs = trxSignedHex

		// add utxo to pending transaction
		Blob := new(blob.Byteblob)
		Blob.SetHex(trxSignedHex)
		bytesBuf := bytes.NewBuffer(Blob.GetData())
		bufReader := io.Reader(bytesBuf)
		t := new(transaction.Transaction)
		t.UnPack(bufReader)
		pendingTrxMgr := model.GlobalDBMgr.PendingTransactionMgr
		for i := 0; i < len(t.Vin); i++ {
			trxId := t.Vin[i].PrevOut.Hash.GetHex()
			vout := t.Vin[i].PrevOut.N
			_, err := pendingTrxMgr.NewPendingTransaction(dbSession, trxUuid, walletConfig.Coinid, trxId, int(vout), "", "")
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300001, pendingTrxMgr.TableName, "insert", "insert pending transaction")
				//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
				ctx.JSON(res)
				return
			}
		}

		// web socket request
		var reqWS CreateTransactionRequestWS
		reqWS.Id = 1
		reqWS.JsonRpc = "2.0"
		reqWS.Method = "create_transaction"
		var createTransactionParamWS CreateTransactionParamWS
		var trxDetail TransferDetailParamWS
		createTransactionParamWS.CoinId = walletConfig.Coinid
		createTransactionParamWS.WalletUuid = walletConfig.Walletuuid
		createTransactionParamWS.TrxUuid = trxUuid
		createTransactionParamWS.FromAddress = walletConfig.Address
		trxDetail.ToAddress = toAddrStr
		trxDetail.Value = amountStr
		createTransactionParamWS.ToDetail = append(createTransactionParamWS.ToDetail, trxDetail)
		createTransactionParamWS.TotalFee = feeCostStr
		createTransactionParamWS.SignedTrx = trxSignedHex
		reqWS.Params = append(reqWS.Params, createTransactionParamWS)

		msgBytes, err := json.Marshal(reqWS)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900010)
			ctx.JSON(res)
			return
		}
		fmt.Println("CreateTrxController() WSRequest create_transaction")
		fmt.Println("request:", string(msgBytes))
		retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900020, err.Error())
			ctx.JSON(res)
			return
		}
		var resWS CreateTransactionResponseWS
		err = json.Unmarshal([]byte(retString), &resWS)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900011)
			ctx.JSON(res)
			return
		}
		if resWS.Error != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
			ctx.JSON(res)
			return
		}
	}

	trx.Confirmed = 1
	trx.Acctconfirmed = strconv.Itoa(sessionValue.AcctId)

	err = trxMgr.UpdateTransaction(dbSession, trx)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction")
		//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 创建提醒
	//relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
	//relations, err := relationMgr.GetRelationsByWalletId(dbSession, req.Params[0].FromWalletId)
	//if err != nil {
	//	dbSession.Rollback()
	//	res.Error = utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
	//	//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
	//	ctx.JSON(res)
	//	return
	//}
	//notifyMgr := model.GlobalDBMgr.NotificationMgr
	//for _, relation := range relations {
	//	if sessionValue.AcctId != relation.Acctid {
	//		_, err := notifyMgr.NewNotification(dbSession, &relation.Acctid, &req.Params[0].FromWalletId, &trx.Trxid,
	//			1, fmt.Sprintf("有一笔新的转账交易产生, 需要您去处理, 交易ID: %d", trx.Trxid),
	//			0, "", "")
	//		if err != nil {
	//			dbSession.Rollback()
	//			res.Error = utils.MakeError(300001, trxMgr.TableName, "insert", "insert notification")
	//			//TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
	//			ctx.JSON(res)
	//			return
	//		}
	//	}
	//}

	err = TransferLog(dbSession, true, sessionValue.AcctId, req.Params[0].FromWalletId, &trx.Trxid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300011)
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

func RevokeLog(dbSession *xorm.Session, isSuccQuit bool, acctId int, trxId int) error {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(dbSession, acctId)
	if err != nil {
		return err
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxRevoke), acctConfig.Realname,
		strconv.Itoa(trxId), resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(dbSession, acctId, 5, logContent)
	if err != nil {
		return err
	}
	return nil
}

type RevokeTrxParam struct {
	SessionId string `json:"sessionid"`
	OperateId int    `json:"operateid"`
	TrxId     int    `json:"trxid"`
	Signature string `json:"signature"`
}

type RevokeTrxRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []RevokeTrxParam `json:"params"`
}

type RevokeTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func RevokeTrxController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req RevokeTrxRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res RevokeTrxResponse
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

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 4, req.Params[0].OperateId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "revoke"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	trxIdStr := strconv.Itoa(req.Params[0].TrxId)
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, trxIdStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	_, trx, err := trxMgr.GetTransactionById(dbSession, req.Params[0].TrxId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	if trx.State != 0 && trx.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "trx", trx.Trxid)
		ctx.JSON(res)
		return
	}

	if sessionValue.Role == 1 {
		if trx.Serverid != config.LocalServerId || sessionValue.AcctId != trx.Acctid {
			dbSession.Rollback()
			res.Error = utils.MakeError(200011)
			ctx.JSON(res)
			return
		}
	}

	// 删除提醒
	//notifyMgr := model.GlobalDBMgr.NotificationMgr
	//notifyMgr.DeleteNotification(dbSession, nil, nil, nil, &trx.Trxid, nil, nil, nil, nil)
	//if err != nil {
	//	dbSession.Rollback()
	//	res.Error = utils.MakeError(300001, trxMgr.TableName, "delete", "delete notification")
	//	//RevokeLog(false, sessionValue.AcctId, req.Params[0].TrxId)
	//	ctx.JSON(res)
	//	return
	//}

	err = RevokeLog(dbSession, true, sessionValue.AcctId, req.Params[0].TrxId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300012)
		ctx.JSON(res)
		return
	}

	if trx.State != 0 {
		// 不直接更改交易状态 等待服务端推送交易状态变更
		// web socket request
		var reqWS RevokeTransactionRequestWS
		reqWS.Id = 1
		reqWS.JsonRpc = "2.0"
		reqWS.Method = "revoke_transaction"
		var revokeTrxParamWS RevokeTransactionParamWS
		revokeTrxParamWS.TrxUuid = trx.Trxuuid
		reqWS.Params = append(reqWS.Params, revokeTrxParamWS)
		msgBytes, err := json.Marshal(reqWS)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900010)
			ctx.JSON(res)
			return
		}
		fmt.Println("RevokeTrxController() WSRequest revoke_transaction")
		fmt.Println("request:", string(msgBytes))
		retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900020, err.Error())
			ctx.JSON(res)
			return
		}
		var resWS RevokeTransactionResponseWS
		err = json.Unmarshal([]byte(retString), &resWS)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900011)
			ctx.JSON(res)
			return
		}
		if resWS.Error != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
			ctx.JSON(res)
			return
		}
	} else {
		err = trxMgr.UpdateTransactionState(dbSession, req.Params[0].TrxId, 4)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction state")
			//RevokeLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
	}

	err = dbSession.Commit()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	ctx.JSON(res)
	return
}

func ConfirmLog(dbSession *xorm.Session, isSuccQuit bool, acctId int, trxId int) error {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(dbSession, acctId)
	if err != nil {
		return err
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxConfirm), acctConfig.Realname,
		strconv.Itoa(trxId), resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(dbSession, acctId, 5, logContent)
	if err != nil {
		return err
	}
	return nil
}

type ConfirmTrxParam struct {
	SessionId string `json:"sessionid"`
	OperateId int    `json:"operateid"`
	TrxId     int    `json:"trxid"`
	Signature string `json:"signature"`
}

type ConfirmTrxRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []ConfirmTrxParam `json:"params"`
}

type ConfirmTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ConfirmTrxController(ctx iris.Context, jsonRpcBody []byte) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var req ConfirmTrxRequest
	err = json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res ConfirmTrxResponse
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
	if sessionValue.Role != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200012)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(dbSession, 4, req.Params[0].OperateId)
	if !verify || err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "confirm"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	trxIdStr := strconv.Itoa(req.Params[0].TrxId)
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, trxIdStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	_, trx, err := trxMgr.GetTransactionById(dbSession, req.Params[0].TrxId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	if trx.State != 0 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "trx", trx.Trxid)
		ctx.JSON(res)
		return
	}

	acctConfirmed := strings.Split(trx.Acctconfirmed, ",")
	for _, acctStr := range acctConfirmed {
		if strconv.Itoa(sessionValue.AcctId) == acctStr {
			dbSession.Rollback()
			res.Error = utils.MakeError(200014)
			ctx.JSON(res)
			return
		}
	}

	serverConfirmed := strings.Split(trx.Signedserverids, ",")
	for _, serverStr := range serverConfirmed {
		if strconv.Itoa(config.LocalServerId) == serverStr {
			dbSession.Rollback()
			res.Error = utils.MakeError(200017)
			ctx.JSON(res)
			return
		}
	}

	// 检查Account的有效性
	uErr := CheckAccountAvailable(sessionValue.AcctId)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 检查交易出账的Wallet是否属于Account关联的Wallet
	uErr = CheckRelationAvailable(sessionValue.AcctId, trx.Walletid)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 获取Wallet的配置信息
	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(dbSession, trx.Walletid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}
	if walletConfig.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	if walletConfig.Needsigcount == 0 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200019, walletConfig.Walletid)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	} else if trx.Needconfirm == 0 {
		trx.Needconfirm = walletConfig.Needsigcount
	}

	// 达到最大确认数
	if trx.Confirmed == trx.Needconfirm {
		dbSession.Rollback()
		res.Error = utils.MakeError(200013)
		ctx.JSON(res)
		return
	}

	// 达到最大服务器确认数
	if len(serverConfirmed) == walletConfig.Needkeysigcount {
		dbSession.Rollback()
		res.Error = utils.MakeError(200018)
		ctx.JSON(res)
		return
	}

	// parse serverkeys
	serverIdKeyPairMap := make(map[int]int)
	var serverKeyInfos []model.ServerKeyInfo
	serverIdKeyPairStrs := strings.Split(walletConfig.Serverkeys, ",")
	for _, serverIdKeyPairStr := range serverIdKeyPairStrs {
		tmpStrs := strings.Split(serverIdKeyPairStr, ":")
		if len(tmpStrs) != 2 {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
		serverId, err := strconv.Atoi(tmpStrs[0])
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
		keyIndex, err := strconv.Atoi(tmpStrs[1])
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(200102, walletConfig.Walletid)
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
		serverIdKeyPairMap[serverId] = keyIndex

		var serverKeyInfo model.ServerKeyInfo
		serverKeyInfo.ServerId = serverId
		serverKeyInfo.KeyIndex = keyIndex
		serverKeyInfos = append(serverKeyInfos, serverKeyInfo)
	}

	KeyIndex, ok := serverIdKeyPairMap[config.LocalServerId]
	if !ok {
		dbSession.Rollback()
		res.Error = utils.MakeError(200102, walletConfig.Walletid)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 获取钱包对应PubKey信息
	var pubKeyStrSlice []string
	for _, serverKeyInfo := range serverKeyInfos {
		keyPoolMgr := model.GlobalDBMgr.PubKeyPoolMgr
		pubkeyStr, err := keyPoolMgr.QueryPubKeyByKeyIndex(dbSession, serverKeyInfo.ServerId, serverKeyInfo.KeyIndex)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get key pool config")
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
		pubKeyStrSlice = append(pubKeyStrSlice, pubkeyStr)
	}

	// 获取Coin配置信息
	coinMgr := model.GlobalDBMgr.CoinConfigMgr
	coinConfig, err := coinMgr.GetCoin(dbSession, walletConfig.Coinid)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}
	if coinConfig.State != 1 {
		dbSession.Rollback()
		res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(coinConfig.Coinsymbol) {
		dbSession.Rollback()
		res.Error = utils.MakeError(600001, coinConfig.Coinsymbol)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 检查存在Dest Address设置下  交易中to address的有效性
	trxStrList := strings.Split(trx.Todetails, "~")
	uErr = CheckDestAddress(trxStrList[0], walletConfig.Destaddress, walletConfig.Walletid)
	if uErr != nil {
		dbSession.Rollback()
		res.Error = uErr
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// list unspent
	utxos, err := coin.ListUnSpent(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, walletConfig.Address)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(200103)
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// verify if balance is enough or not
	// TODO

	if trx.Serverid == config.LocalServerId {
		if trx.Confirmed +1 == trx.Needconfirm {
			// filter unspent
			utxosFilter := make([]coin.UTXODetail, 0)
			for _, utxo := range utxos {
				count, err := model.GlobalDBMgr.PendingTransactionMgr.GetPendingTransactionByVinCount(dbSession, coinConfig.Coinid, utxo.TxId, utxo.Vout)
				if err != nil {
					dbSession.Rollback()
					res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get pending transaction")
					//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
					ctx.JSON(res)
					return
				}
				if count == 0 {
					utxosFilter = append(utxosFilter, utxo)
				}
			}

			feeCostStr, trxSignedHex, err := coin.CreateMultiSignedTransaction(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, uint16(KeyIndex),
				walletConfig.Needkeysigcount, pubKeyStrSlice, utxosFilter, walletConfig.Address, trxStrList[0], trxStrList[1],
				trx.Fee, trx.Gasprice, trx.Gaslimit)

			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(800000, err.Error())
				//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
				ctx.JSON(res)
				return
			}

			if feeCostStr != "" {
				trx.Feecost = feeCostStr
			}
			trx.State = 1
			trx.Signedserverids = strconv.Itoa(config.LocalServerId)
			trx.Signedtrxs = trxSignedHex

			// web socket request
			var reqWS CreateTransactionRequestWS
			reqWS.Id = 1
			reqWS.JsonRpc = "2.0"
			reqWS.Method = "create_transaction"
			var createTransactionParamWS CreateTransactionParamWS
			var trxDetail TransferDetailParamWS
			createTransactionParamWS.CoinId = walletConfig.Coinid
			createTransactionParamWS.WalletUuid = walletConfig.Walletuuid
			createTransactionParamWS.TrxUuid = trx.Trxuuid
			createTransactionParamWS.FromAddress = walletConfig.Address
			createTransactionParamWS.TotalFee = feeCostStr
			trxDetail.ToAddress = trxStrList[0]
			trxDetail.Value = trxStrList[1]
			createTransactionParamWS.ToDetail = append(createTransactionParamWS.ToDetail, trxDetail)
			createTransactionParamWS.SignedTrx = trxSignedHex
			reqWS.Params = append(reqWS.Params, createTransactionParamWS)

			msgBytes, err := json.Marshal(reqWS)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900010)
				ctx.JSON(res)
				return
			}
			fmt.Println("ConfirmTrxController() WSRequest create_transaction")
			fmt.Println("request:", string(msgBytes))
			retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900020, err.Error())
				ctx.JSON(res)
				return
			}
			var resWS CreateTransactionResponseWS
			err = json.Unmarshal([]byte(retString), &resWS)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900011)
				ctx.JSON(res)
				return
			}
			if resWS.Error != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
				ctx.JSON(res)
				return
			}

			// add utxo to pending transaction
			Blob := new(blob.Byteblob)
			Blob.SetHex(trxSignedHex)
			bytesBuf := bytes.NewBuffer(Blob.GetData())
			bufReader := io.Reader(bytesBuf)
			t := new(transaction.Transaction)
			t.UnPack(bufReader)
			pendingTrxMgr := model.GlobalDBMgr.PendingTransactionMgr
			for i := 0; i < len(t.Vin); i++ {
				trxId := t.Vin[i].PrevOut.Hash.GetHex()
				vout := t.Vin[i].PrevOut.N
				_, err := pendingTrxMgr.NewPendingTransaction(dbSession, trx.Trxuuid, walletConfig.Coinid, trxId, int(vout), "", "")
				if err != nil {
					dbSession.Rollback()
					res.Error = utils.MakeError(300001, pendingTrxMgr.TableName, "insert", "insert pending transaction")
					//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
					ctx.JSON(res)
					return
				}
			}
		}

		// 交易防篡改签名验签
		trxStr := CreateTransactionString(trx.Trxuuid, walletConfig.Walletid, walletConfig.Coinid,
			trx.Contractaddr, trx.Acctid, trx.Serverid, walletConfig.Address,
			trx.Todetails, walletConfig.Needsigcount,
			trx.Fee, trx.Gasprice, trx.Gaslimit)

		signature, err := hex.DecodeString(trx.Signature)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900004)
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}
		_, err = coin.CoinVerifyTrx('2', uint16(KeyIndex), utility.Sha256([]byte(trxStr)), signature)
		if err != nil {
			dbSession.Rollback()
			res.Error = utils.MakeError(900003)
			//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}

		trx.Trxid = req.Params[0].TrxId
		trx.Confirmed = trx.Confirmed + 1
		trx.Acctconfirmed = trx.Acctconfirmed + "," + strconv.Itoa(sessionValue.AcctId)
	} else {
		if trx.Confirmed +1 == trx.Needconfirm {
			// TODO
			//先反序列化表中第一个signedtrx，将vin中的sigscript置为nil
			Blob := new(blob.Byteblob)
			signedtrx1 := strings.Split(trx.Signedtrxs,",")[0]
			Blob.SetHex(signedtrx1)
			bytesBuf := bytes.NewBuffer(Blob.GetData())
			bufReader := io.Reader(bytesBuf)
			transaction := new(transaction.Transaction)
			transaction.UnPack(bufReader)
			for k,_ := range transaction.Vin {
				transaction.Vin[k].ScriptSig.SetScriptBytes(nil)
				transaction.Vin[k].ScriptWitness.SetScriptWitnessBytes(nil)
			}

			//再序列化交易
			bytesBuf_db := bytes.NewBuffer([]byte{})
			bufWriter := io.Writer(bytesBuf_db)
			err = transaction.Pack(bufWriter)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300014)
				return
			}
			encodetrx := hex.EncodeToString(bytesBuf_db.Bytes())
			//生成redeemscript
			redeemScript, err := coin.GetMultiSignRedeemScript(walletConfig.Needkeysigcount, pubKeyStrSlice, coinConfig.Coinsymbol)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300015)
				//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
				ctx.JSON(res)
				return
			}

			// sign transaction
			trxSigHex, err := coin.MultiSignRawTransaction(coinConfig.Coinsymbol,coinConfig.Ip,coinConfig.Rpcport,coinConfig.Rpcuser,coinConfig.Rpcpass,utxos,uint16(KeyIndex), redeemScript, encodetrx)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(300016)
				//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
				ctx.JSON(res)
				return
			}
			//fmt.Println("trxSigHex:", trxSigHex)
			trx.State = 1
			trx.Signedserverids = trx.Signedserverids + "," + strconv.Itoa(config.LocalServerId)
			trx.Signedtrxs = trx.Signedtrxs + "," + trxSigHex

			// web socket request
			var reqWS ConfirmTransactionRequestWS
			reqWS.Id = 1
			reqWS.JsonRpc = "2.0"
			reqWS.Method = "confirm_transaction"
			var confirmTransactionParamWS ConfirmTransactionParamWS
			confirmTransactionParamWS.TrxUuid = trx.Trxuuid
			confirmTransactionParamWS.SignedTrx = trxSigHex
			reqWS.Params = append(reqWS.Params, confirmTransactionParamWS)

			msgBytes, err := json.Marshal(reqWS)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900010)
				ctx.JSON(res)
				return
			}
			fmt.Println("ConfirmTrxController() WSRequest confirm_transaction")
			fmt.Println("request:", string(msgBytes))
			retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900020, err.Error())
				ctx.JSON(res)
				return
			}
			var resWS ConfirmTransactionResponseWS
			err = json.Unmarshal([]byte(retString), &resWS)
			if err != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900011)
				ctx.JSON(res)
				return
			}
			if resWS.Error != nil {
				dbSession.Rollback()
				res.Error = utils.MakeError(900021, resWS.Error.ErrMsg)
				ctx.JSON(res)
				return
			}
		}

		trx.Trxid = req.Params[0].TrxId
		trx.Confirmed = trx.Confirmed + 1
		if trx.Acctconfirmed == "" {
			trx.Acctconfirmed = strconv.Itoa(sessionValue.AcctId)
		} else {
			trx.Acctconfirmed = trx.Acctconfirmed + "," + strconv.Itoa(sessionValue.AcctId)
		}
	}

	err = trxMgr.UpdateTransaction(dbSession, trx)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction")
		//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 删除提醒
	//notifyMgr := model.GlobalDBMgr.NotificationMgr
	//var pAcctId *int
	//pAcctId = nil
	//if trx.State == 0 {
	//	pAcctId = &sessionValue.AcctId
	//}
	//notifyMgr.DeleteNotification(dbSession, nil, pAcctId, nil, &trx.Trxid, nil, nil, nil, nil)
	//if err != nil {
	//	dbSession.Rollback()
	//	res.Error = utils.MakeError(300001, trxMgr.TableName, "delete", "delete notification")
	//	//ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
	//	ctx.JSON(res)
	//	return
	//}

	err = ConfirmLog(dbSession, true, sessionValue.AcctId, req.Params[0].TrxId)
	if err != nil {
		dbSession.Rollback()
		res.Error = utils.MakeError(300013)
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

func TransactionController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "get_wallet_trxs" {
		GetTransactionController(ctx, jsonRpcBody)
	} else if funcName == "transfer" {
		CreateTrxController(ctx, jsonRpcBody)
	} else if funcName == "confirm" {
		ConfirmTrxController(ctx, jsonRpcBody)
	} else if funcName == "revoke" {
		RevokeTrxController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}

type TransferDetailParamWS struct {
	ToAddress 		string                  `json:"toaddress"`
	Value 			string                  `json:"value"`
}

type CreateTransactionParamWS struct {
	CoinId          int                     `json:"coinid"`
	WalletUuid 		string                  `json:"walletuuid"`
	TrxUuid 		string                  `json:"trxuuid"`
	FromAddress 	string                  `json:"fromaddress"`
	ToDetail 		[]TransferDetailParamWS                  `json:"todetail"`
	TotalFee 		string                  `json:"totalfee"`
	SignedTrx 		string                  `json:"signedtrx"`
}

type CreateTransactionRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []CreateTransactionParamWS 		`json:"params"`
}

type CreateTransactionResponseWS struct {
	Id     int                      `json:"id"`
	Result bool					 	`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type ConfirmTransactionParamWS struct {
	TrxUuid 		string                  `json:"trxuuid"`
	SignedTrx 		string                  `json:"signedtrx"`
}

type ConfirmTransactionRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []ConfirmTransactionParamWS 		`json:"params"`
}

type ConfirmTransactionResponseWS struct {
	Id     int                      `json:"id"`
	Result bool					 	`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type RevokeTransactionParamWS struct {
	TrxUuid   string             	`json:"trxuuid"`
}

type RevokeTransactionRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []RevokeTransactionParamWS 		`json:"params"`
}

type RevokeTransactionResponseWS struct {
	Id     int                      `json:"id"`
	Result bool					 	`json:"result"`
	Error  *utils.Error             `json:"error"`
}

type QueryTransactionParamWS struct {
	TrxUuid   string             	`json:"trxuuid"`
}

type QueryTransactionRequestWS struct {
	Id      int                     `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  []QueryTransactionParamWS 		`json:"params"`
}

type QueryTransactionResultWS struct {
	TrxId      		int                     `json:"trxid"`
	TrxUuid     	string                  `json:"trxuuid"`
	RawTrxId      	string                  `json:"rawtrxid"`
	CoinId     		int                 	`json:"coinid"`
	WalletUuid     	string                  `json:"walletuuid"`
	FromAddress     string                  `json:"fromaddress"`
	ToDetail 		[]TransferDetailParamWS `json:"todetail"`
	TotalFee 		string                  `json:"totalfee"`
	SignedTrx 		[]string                `json:"signedtrx"`
	CreateServerId  int                 	`json:"createserverid"`
	State			int                     `json:"state"`
	SignedServerIds []int                   `json:"signedserverids"`
}

type QueryTransactionResponseWS struct {
	Id     int                      `json:"id"`
	Result QueryTransactionResultWS		`json:"result"`
	Error  *utils.Error             `json:"error"`
}

