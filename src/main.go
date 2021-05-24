package main

import (
	"bufio"
	"bytes"
	"coin"
	"config"
	"controller"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"github.com/mutalisk999/bitcoin-lib/src/blob"
	"github.com/mutalisk999/bitcoin-lib/src/script"
	"github.com/mutalisk999/bitcoin-lib/src/serialize"
	"github.com/mutalisk999/bitcoin-lib/src/transaction"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"io"
	"model"
	"os"
	"session"
	"strconv"
	"strings"
	"time"
	"utils"
	"math"
)

var app *iris.Application
var goroutineMgr *goroutine_mgr.GoroutineManager

func DoSessionMaintain(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	mgr := session.GlobalSessionMgr
	sessionsOvertime := make([]string, 0)
	for {
		mgr.Mutex.Lock()
		for sid, sessionValue := range mgr.SessionStore {
			if time.Now().Unix()-sessionValue.UpdateTime.Unix() > 30*60 {
				sessionsOvertime = append(sessionsOvertime, sid)
			}
		}
		mgr.Mutex.Unlock()

		dbSession := model.GetDBEngine().NewSession()
		dbSession.Begin()

		for _, sid := range sessionsOvertime {
			session.GlobalSessionMgr.DeleteSessionValue(dbSession, sid)
		}

		dbSession.Commit()
		dbSession.Close()

		time.Sleep(30 * time.Second)
	}
}

func DoTransactionMaintain(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	for {
		dbSession := model.GetDBEngine().NewSession()
		dbSession.Begin()

		trxMgr := model.GlobalDBMgr.TransactionMgr
		trxs, err := trxMgr.GetUnComfirmedTransactions(dbSession)

		if err != nil {
			fmt.Println("DoTransactionMaintain | GetTransactionsByState: " + err.Error())
		}

		for _, trx := range trxs {
			coinConfigMgr := model.GlobalDBMgr.CoinConfigMgr
			coinConfig, err := coinConfigMgr.GetCoin(dbSession, trx.Coinid)
			if err != nil {
				fmt.Println("DoTransactionMaintain | GetCoin: " + err.Error())
				continue
			}

			isConfirmed, err := coin.IsTrxConfirmed(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport,
				coinConfig.Rpcuser, coinConfig.Rpcpass, trx.Rawtrxid)
			if !isConfirmed {
				if err != nil {
					fmt.Println("DoTransactionMaintain | IsTrxConfirmed: " + err.Error())
				}
				continue
			}

			trxState := 3
			if err != nil {
				trxState = 4
			}

			if isConfirmed {
				err := coin.ConfirmTransaction(dbSession, coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport,
					coinConfig.Rpcuser, coinConfig.Rpcpass, trx.Trxid, trx.Rawtrxid, trxState)
				if err != nil {
					dbSession.Rollback()
					fmt.Println("DoTransactionMaintain | ConfirmTransaction: " + err.Error())
					continue
				}
				err = model.GlobalDBMgr.PendingTransactionMgr.DeletePendingTransactionByTrxUuid(dbSession, trx.Trxuuid)
				if err != nil {
					dbSession.Rollback()
					fmt.Println("DoTransactionMaintain | DeletePendingTransactionByTrxUuid: " + err.Error())
					continue
				}
			}

			time.Sleep(1 * time.Second)
		}

		dbSession.Commit()
		dbSession.Close()

		time.Sleep(30 * time.Second)
	}
}

func RequestToUpdateTaskState(dbSession *xorm.Session, taskUuid string, state int) error {
	// accept_task
	var req utils.TaskAcceptRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "accept_task"
	var param utils.TaskAcceptParamWS
	param.TaskUuid = taskUuid
	param.State = state
	req.Params = append(req.Params, param)

	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("RequestToUpdateTaskState() WSRequest accept_task")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "task", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res utils.TaskAcceptResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return err
	}
	return nil
}

func DoTaskPersistenceProcess(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	for {
		var task utils.TaskPushParamWS
		task = <- utils.GlobalPushTasks

		dbSession := model.GetDBEngine().NewSession()
		dbSession.Begin()

		isFound, _, err := model.GlobalDBMgr.TaskPersistenceMgr.GetTaskByTaskUuid(dbSession, task.TaskUuid)
		if err != nil {
			fmt.Println("DoTaskPersistenceProcess | GetTaskByTaskUuid: " + err.Error())
			continue
		}

		if isFound {
			continue
		}

		err = model.GlobalDBMgr.TaskPersistenceMgr.InsertTask(dbSession, task.TaskUuid, task.WalletUuid, task.TrxUuid, task.PushType, 1)
		if err != nil {
			fmt.Println("DoTaskPersistenceProcess | InsertTask: " + err.Error())
			dbSession.Rollback()
			dbSession.Close()
			continue
		}

		err = RequestToUpdateTaskState(dbSession, task.TaskUuid, 1)
		if err != nil {
			fmt.Println("DoTaskPersistenceProcess | RequestToUpdateTaskState: " + err.Error())
			dbSession.Rollback()
			dbSession.Close()
			continue
		}

		dbSession.Commit()
		dbSession.Close()
	}
}

func DoTaskProcess(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	for {
		dbSession := model.GetDBEngine().NewSession()
		dbSession.Begin()

		tasks, err := model.GlobalDBMgr.TaskPersistenceMgr.GetTasksByState(dbSession, 1)
		if err != nil {
			fmt.Println("DoTaskProcess | GetTasksByState: " + err.Error())
		} else if len(tasks) == 0 {
			time.Sleep(5 * time.Second)
		} else {
			for _, task := range tasks {
				state := 2
				fmt.Println("push type:", task.Pushtype)
				if task.Pushtype == 1 {
					err = RequestAndUpdateLocalPubKeys(dbSession)
					if err != nil {
						dbSession.Rollback()
						fmt.Println("DoTaskProcess | RequestAndUpdateLocalPubKeys:" + err.Error())
						state = 3
					}
					dbSession.Commit()
				} else if task.Pushtype == 2 {
					err = RequestAndCreateLocalWallet(dbSession, task.Walletuuid)
					if err != nil {
						dbSession.Rollback()
						fmt.Println("DoTaskProcess | RequestAndCreateLocalWallet:" + err.Error())
						state = 3
					}
					dbSession.Commit()
				} else if task.Pushtype == 3 {
					err = RequestAndUpdateLocalWalletResult(dbSession, task.Walletuuid)
					if err != nil {
						dbSession.Rollback()
						fmt.Println("DoTaskProcess | RequestAndUpdateLocalWalletResult:" + err.Error())
						state = 3
					}
					dbSession.Commit()
				} else if task.Pushtype == 4 {
					err = RequestAndCreateLocalTrx(dbSession, task.Walletuuid, task.Trxuuid)
					if err != nil {
						dbSession.Rollback()
						fmt.Println("DoTaskProcess | RequestAndCreateLocalTrx:" + err.Error())
						state = 3
					}
					dbSession.Commit()
				} else if task.Pushtype == 5 {
					err = RequestAndUpdateLocalTrxResult(dbSession, task.Walletuuid, task.Trxuuid)
					if err != nil {
						dbSession.Rollback()
						fmt.Println("DoTaskProcess | RequestAndUpdateLocalTrxResult:" + err.Error())
						state = 3
					}
					dbSession.Commit()
				} else {
					fmt.Println("DoTaskProcess | push type: ", task.Pushtype, "not support")
				}

				err := model.GlobalDBMgr.TaskPersistenceMgr.UpdateTaskState(dbSession, task.Id, state)
				if err != nil {
					fmt.Println("DoTaskProcess | UpdateTaskState: " + err.Error())
					dbSession.Rollback()
				}

				err = RequestToUpdateTaskState(dbSession, task.Taskuuid, state)
				if err != nil {
					fmt.Println("DoTaskPersistenceProcess | RequestToUpdateTaskState: " + err.Error())
					dbSession.Rollback()
				}
				dbSession.Commit()
			}
		}
		dbSession.Close()
	}
}

func StartWebSocketConnect() uint64 {
	fmt.Println("start goroutine WebSocketConnector...")
	return goroutineMgr.GoroutineCreatePn("WebSocketConnector", DoWebSocketConnect, nil)
}

func StartWebSocketConnMaintainer() uint64 {
	fmt.Println("start goroutine WebSocketConnMaintainer...")
	return goroutineMgr.GoroutineCreatePn("WebSocketConnMaintainer", DoWebSocketConnMaintain, nil)
}

func StartSessionMaintainer() uint64 {
	fmt.Println("start goroutine SessionMaintainer...")
	return goroutineMgr.GoroutineCreatePn("SessionMaintainer", DoSessionMaintain, nil)
}

func StartTransactionMaintainer() uint64 {
	fmt.Println("start goroutine TransactionMaintainer...")
	return goroutineMgr.GoroutineCreatePn("TransactionMaintainer", DoTransactionMaintain, nil)
}

func StartTaskPersistenceProcessor() uint64 {
	fmt.Println("start goroutine TaskPersistenceProcessor...")
	return goroutineMgr.GoroutineCreatePn("TaskPersistenceProcessor", DoTaskPersistenceProcess, nil)
}

func StartTaskProcessor() uint64 {
	fmt.Println("start goroutine TaskProcessor...")
	return goroutineMgr.GoroutineCreatePn("TaskProcessor", DoTaskProcess, nil)
}


func LoadConf() error {
	// init config
	jsonParser := new(config.JsonStruct)
	//err := jsonParser.Load("D:/EMultiWallet_Platform/src/config.json", &config.GlobalConfig)
	err := jsonParser.Load("config.json", &config.GlobalConfig)
	if err != nil {
		fmt.Println("Load config.json", err)
		return err
	}
	return nil
}

func Init() error {
	utils.JsonId = 0
	utils.GlobalIsLogin = false
	utils.GlobalWsConnEstablish = make(chan bool)
	utils.GlobalWsConnReconnect = make(chan bool)
	utils.GlobalPushTasks = make(chan utils.TaskPushParamWS, 1024)

	err := LoadConf()
	if err != nil {
		return err
	}

	utils.InitGlobalError()
	session.InitSessionMgr()

	goroutineMgr = new(goroutine_mgr.GoroutineManager)
	goroutineMgr.Initialise("MainGoroutineManager")

	fmt.Println("db path:", config.GlobalConfig.DbConfig.DbSource)
	err = model.InitDB(config.GlobalConfig.DbConfig.DbType, config.GlobalConfig.DbConfig.DbSource)
	if err != nil {
		return err
	}

	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	// load persistent sessions from db
	sessionPersistences, err := model.GlobalDBMgr.SessionPersistenceMgr.LoadAllSessionValue(dbSession)
	if err != nil {
		return err
	}
	for _, sessionPersistence := range sessionPersistences {
		var sessionValue session.SessionValue
		err := json.Unmarshal([]byte(sessionPersistence.Sessionvalue), &sessionValue)
		if err != nil {
			return err
		}
		sessionValue.UpdateTime = time.Now()
		session.GlobalSessionMgr.SetSessionValue(sessionPersistence.Sessionid, sessionValue)
	}

	app = iris.New()
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	return nil
}

func Run(endpoint string, charset string) {
	app.Run(iris.Addr(endpoint), iris.WithCharset(charset))
}

func RunTLS(endpoint string, certFile string, keyFile string, charset string) {
	app.Run(iris.TLS(endpoint, certFile, keyFile), iris.WithCharset(charset))
}

func RegisterUrlRouter() {
	app.Post("/apis/authcode", controller.AuthCodeController)
	app.Post("/apis/identity", controller.IdentityController)
	app.Post("/apis/user", controller.UserController)
	app.Post("/apis/account", controller.AccountController)
	app.Post("/apis/wallet", controller.WalletController)
	app.Post("/apis/notification", controller.NotificationController)
	app.Post("/apis/transaction", controller.TransactionController)
	app.Post("/apis/log", controller.LogController)
	app.Post("/apis/coin", controller.CoinController)
	app.Post("/apis/manager", controller.ManagerController)
}

func InitLocalServerId(dbSession *xorm.Session) error {
	isFound, serverInfo, err := model.GlobalDBMgr.ServerInfoMgr.GetLocalServerInfo(dbSession)
	if err != nil {
		return err
	}
	if !isFound {
		config.LocalServerId = -1
		return nil
	} else {
		config.LocalServerId = serverInfo.Serverid
		return nil
	}
	return nil
}

func RequestAndUpdateLocalPubKeys(dbSession *xorm.Session) error {
	// web socket request for query_pubkeys
	var req3 controller.QueryPubKeyRequestWS
	req3.Id = 1
	req3.JsonRpc = "2.0"
	req3.Method = "query_pubkeys"
	req3.Params = make([]interface{}, 0)
	msgBytes, err := json.Marshal(req3)
	if err != nil {
		return err
	}
	fmt.Println("RequestAndUpdateLocalPubKeys() WSRequest query_pubkeys")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "pubkey", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res3 controller.QueryPubKeyResponseWS
	err = json.Unmarshal([]byte(retString), &res3)
	if err != nil {
		return err
	}
	if res3.Error != nil {
		return errors.New(res3.Error.ErrMsg)
	}
	for _, queryPubKeyResult := range res3.Result {
		isFound, serverInfo, err := model.GlobalDBMgr.ServerInfoMgr.GetServerInfoCountById(dbSession, queryPubKeyResult.ServerId)
		if err != nil {
			return err
		}
		if !isFound {
			isLocalServer := false
			if queryPubKeyResult.ServerId == config.LocalServerId {
				isLocalServer = true
			}
			if isLocalServer {
				if config.GlobalConfig.ServerInfoConfig.ServerName != queryPubKeyResult.ServerName {
					return errors.New("different Server name from Local to Server name from Remote")
				}
				if config.GlobalConfig.ServerInfoConfig.StartIndex != queryPubKeyResult.StartIndex {
					return errors.New("different Server start index from Local to Server start index from Remote")
				}
			}
			err = model.GlobalDBMgr.ServerInfoMgr.InsertServerInfo(dbSession, queryPubKeyResult.ServerId, queryPubKeyResult.ServerName,
				isLocalServer, "", queryPubKeyResult.StartIndex, 1)
			if err != nil {
				return err
			}
		} else {
			isLocalServer := false
			if queryPubKeyResult.ServerId == config.LocalServerId {
				isLocalServer = true
			}

			if serverInfo.Islocalserver != isLocalServer {
				return errors.New("different islocalserver flag from Local DB to islocalserver flag from Remote")
			}
			if serverInfo.Servername != queryPubKeyResult.ServerName {
				return errors.New("different Server name from Local DB to Server name from Remote")
			}
			if serverInfo.Serverstartindex != queryPubKeyResult.StartIndex {
				return errors.New("different Server start index from Local DB to Server start index from Remote")
			}
		}

		serverKeyInfos, err := model.GlobalDBMgr.PubKeyPoolMgr.LoadPubKeysByServerId(dbSession, queryPubKeyResult.ServerId)
		if err != nil {
			return err
		}
		if len(serverKeyInfos) != 0 {
			if len(serverKeyInfos) != len(queryPubKeyResult.Keys) {
				return errors.New(fmt.Sprintf("different PubKey count from Local DB to PubKey count from Remote, ServerId: %d",
					queryPubKeyResult.ServerId))
			}

			for _, serverKeyInfo := range serverKeyInfos {
				pubKeyStr, ok := queryPubKeyResult.Keys[strconv.Itoa(serverKeyInfo.Keyindex)]
				if !ok {
					return errors.New(fmt.Sprintf("different PubKey count from Local DB to PubKey count from Remote, ServerId: %d, KeyIndex: %d",
						queryPubKeyResult.ServerId, serverKeyInfo.Keyindex))
				} else {
					if serverKeyInfo.Pubkey != pubKeyStr {
						return errors.New(fmt.Sprintf("different PubKey count from Local DB to PubKey count from Remote, ServerId: %d, KeyIndex: %d",
							queryPubKeyResult.ServerId, serverKeyInfo.Keyindex))
					}
				}
			}
		} else {
			for k,v := range queryPubKeyResult.Keys {
				keyIndex, err := strconv.Atoi(k)
				if err != nil {
					return err
				}
				err = model.GlobalDBMgr.PubKeyPoolMgr.InsertPubkey(dbSession, queryPubKeyResult.ServerId, keyIndex, v)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func RequestAndCreateLocalWallet(dbSession *xorm.Session, walletUuid string) error {
	// web socket request for query_wallet
	var req controller.QueryWalletRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "query_wallet"
	var param controller.QueryWalletParamWS
	param.WalletUuids = append(param.WalletUuids, walletUuid)
	req.Params = append(req.Params, param)
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("RequestAndCreateLocalWallet() WSRequest query_wallet")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "wallet", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.QueryWalletResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}

	if len(res.Result) >= 1 {
		walletCfg := res.Result[0]
		coinConfig, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession, walletCfg.CoinId)
		if err != nil {
			return errors.New("invalid coin id")
		}
		if walletCfg.CreateServerId == config.LocalServerId {
			return errors.New("CreateServerId is same with LocalServerId")
		}
		if walletCfg.NeedSigCount > walletCfg.TotalCount {
			return errors.New("NeedSigCount is bigger than TotalCount")
		}
		if walletCfg.TotalCount != len(walletCfg.KeyDetail) {
			return errors.New("TotalCount not same with the size of KeyDetail")
		}
		var pubKeySlice []string
		var serverKeyStrSlice []string
		isLocalWallet := false
		for _, serverKey := range walletCfg.KeyDetail {
			pubKey, err := model.GlobalDBMgr.PubKeyPoolMgr.UsePubkey(dbSession, serverKey.ServerId, serverKey.KeyIndex)
			if err != nil {
				return err
			}
			pubKeySlice = append(pubKeySlice, pubKey)
			serverKeyStrSlice = append(serverKeyStrSlice, fmt.Sprintf("%d:%d", serverKey.ServerId, serverKey.KeyIndex))

			if serverKey.ServerId == config.LocalServerId {
				isLocalWallet = true
			}
		}
		serverKeysStr := strings.Join(serverKeyStrSlice, ",")

		address, err := coin.GetMultiSignAddressByPubKeys(walletCfg.NeedSigCount, pubKeySlice, coinConfig.Coinsymbol)
		if err != nil {
			return err
		}

		if address != walletCfg.Address {
			return errors.New("Address is not match to the pubkeys")
		}

		if walletCfg.DestAddress != "" {
			dstAddrList := strings.Split(walletCfg.DestAddress, ",")
			for _, dstAddress := range dstAddrList {
				valid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, dstAddress)
				if err != nil {
					return errors.New(utils.MakeError(800000, err.Error()).ErrMsg)
				}
				if !valid {
					return errors.New(utils.MakeError(500003, dstAddress).ErrMsg)
				}
			}
		}

		if !isLocalWallet {
			return nil
		}

		err = coin.ImportAddress(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, address)
		if err != nil {
			return errors.New(utils.MakeError(800000, err.Error()).ErrMsg)
		}

		err = model.GlobalDBMgr.WalletConfigMgr.InsertWallet(dbSession, walletCfg.CoinId, walletCfg.WalletUuid, walletCfg.WalletName, serverKeysStr,
			walletCfg.CreateServerId, walletCfg.TotalCount,
			walletCfg.NeedSigCount, address, walletCfg.DestAddress, 0, walletCfg.Fee, walletCfg.GasPrice, walletCfg.GasLimit, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func RequestAndUpdateLocalWalletResult(dbSession *xorm.Session, walletUuid string) error {
	// web socket request for query_wallet
	var req controller.QueryWalletRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "query_wallet"
	var param controller.QueryWalletParamWS
	param.WalletUuids = append(param.WalletUuids, walletUuid)
	req.Params = append(req.Params, param)
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("RequestAndUpdateLocalWalletResult() WSRequest query_wallet")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "wallet", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.QueryWalletResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}

	if len(res.Result) >= 1 {
		walletCfg := res.Result[0]
		err = model.GlobalDBMgr.WalletConfigMgr.ChangeWalletStateByUuid(dbSession, walletCfg.WalletUuid, walletCfg.State)
		if err != nil {
			return err
		}
	}

	return nil
}

func RequestAndCreateLocalTrx(dbSession *xorm.Session, walletUuid string, trxUuid string) error {
	// web socket request for query_transaction
	var req controller.QueryTransactionRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "query_transaction"
	var param controller.QueryTransactionParamWS
	param.TrxUuid = trxUuid
	req.Params = append(req.Params, param)
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("RequestAndCreateLocalTrx() WSRequest query_transaction")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.QueryTransactionResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}
	trxWS := res.Result

	if trxWS.WalletUuid != walletUuid {
		return errors.New("wallet uuid is not match")
	}
	if trxWS.CreateServerId == config.LocalServerId {
		return errors.New("CreateServerId is same with LocalServerId")
	}

	walletCfg, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletByUUId(dbSession, walletUuid)
	if walletCfg.State != 1 {
		return errors.New("wallet is not active")
	}
	isFound, _, err := model.GlobalDBMgr.TransactionMgr.GetTransactionByUuId(dbSession, trxUuid)
	if err != nil {
		return err
	} else if isFound {
		return errors.New("duplicate trxUuid")
	}

	coinInfo,err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(dbSession,trxWS.CoinId)
	if coinInfo.State != 1 {
		return errors.New("coin is not active")
	}

	coinObj,exist := config.GlobalSupportCoinMgr[coinInfo.Coinsymbol]
	if !exist{
		dbSession.Rollback()
		return errors.New("coin is not exist in config file.")
	}

	//生成redeemscript
	pubkey_redeemscript,err := getRedeemScriptByPubkeys(walletCfg.Needkeysigcount,walletCfg.Serverkeys,coinInfo.Coinsymbol)
	if err != nil {
		dbSession.Rollback()
		return errors.New("getRedeemScriptByPubkeys failed.%s")
	}

	// list unspent
	utxos, err := coin.ListUnSpent(coinInfo.Coinsymbol, coinInfo.Ip, coinInfo.Rpcport, coinInfo.Rpcuser, coinInfo.Rpcpass, walletCfg.Address)
	if err != nil {
		return errors.New("Listunspent failed.")
	}

	// filter unspent
	utxosFilter := make([]coin.UTXODetail, 0)
	for _, utxo := range utxos {
		count ,err := model.GlobalDBMgr.PendingTransactionMgr.GetPendingTransactionByVinCount(dbSession, coinInfo.Coinid, utxo.TxId, utxo.Vout)
		if err != nil {
			return errors.New("Get utxo failed")
		}
		if count == 0 {
			utxosFilter = append(utxosFilter, utxo)
		}
	}

	var detailSlice []string
	keyID_detail := make(map[string]int64,len(trxWS.ToDetail))
	for _, detail := range trxWS.ToDetail {
		voutvalue,_ := coin.ToPrecisionAmount(detail.Value,coinObj.Precision)
		if strings.Contains(walletCfg.Destaddress,detail.ToAddress) {
			return errors.New("To detail address is not in destaddress of walletconfig table")
		}
		detailStr := fmt.Sprintf("%s:%s", detail.ToAddress, detail.Value)
		keyid,err := coin.GetKeyID(coinInfo.Coinsymbol,detail.ToAddress)
		if err != nil{
			return errors.New("Get keyid from address failed.")
		}
		keyID_detail[keyid] = int64(voutvalue)
		detailSlice = append(detailSlice, detailStr)
	}
	detailsStr := strings.Join(detailSlice, ",")

	//如果coinObj是omni，则todetail长度只能是1
	//判断req_todetail的len是否和req中的todetail的len相等，如果不相等，说明todetail里存在重复的地址
	if (coinObj.IsOmni == true && len(keyID_detail) != 1 && 1 != len(trxWS.ToDetail)) || (len(keyID_detail) != len(trxWS.ToDetail)) {
		return errors.New("Exist same address in todetail of request.")
	}
	// TODO
	if len(trxWS.SignedTrx) != len(trxWS.SignedServerIds){
		return errors.New("Signed serverid is not consistent with signedtrx in transactiono table.")
	}

	db_pending_vin :=make(map[string]int)
	for i := 0;i < len(trxWS.SignedTrx);i++{
		Blob := new(blob.Byteblob)
		Blob.SetHex(trxWS.SignedTrx[i])
		bytesBuf := bytes.NewBuffer(Blob.GetData())
		bufReader := io.Reader(bytesBuf)
		trx := new(transaction.Transaction)
		trx.UnPack(bufReader)

		var totalvinvalue int64 = 0
		var totalvoutvalue int64 = 0

		//agent := coin.AgentFactory(coinSymbol)
		//url := fmt.Sprintf("http://%s:%s@%s:%d", coinInfo.Rpcuser, coinInfo.Rpcpass, coinInfo.Ip, coinInfo.Rpcport)
		//agent.Init(url)
		//验证签名
		keys := strings.Split(walletCfg.Serverkeys,",")
		for i := 0;i < len(trxWS.SignedServerIds);i++{
			var onekey []string = strings.Split(keys[i],":")
			servkey,_ := strconv.Atoi(onekey[0])
			if servkey == trxWS.SignedServerIds[i] {
				//根据serverID、keyindex查找serverpubkey表找到对应的pubkey
				keyindex,_ := strconv.Atoi(onekey[1])
				pubkeyinfo, err := model.GlobalDBMgr.PubKeyPoolMgr.QueryPubKeyByKeyIndex(dbSession,trxWS.SignedServerIds[i],keyindex)
				if err != nil {
					dbSession.Rollback()
					return errors.New("Get unused pubkey failed")
				}

				err = coin.MultiVerifySignRawTransaction(coinInfo.Coinsymbol,coinInfo.Ip, coinInfo.Rpcport,coinInfo.Rpcuser, coinInfo.Rpcpass,utxosFilter, pubkeyinfo, trxWS.SignedTrx[i])
				if err == nil {
					fmt.Println("Verify trxSig1Hex success")
				} else {
					dbSession.Rollback()
					return errors.New("Verify trxsig failed")
				}
			}
		}

		keyid,err := coin.GetKeyID(coinInfo.Coinsymbol,trxWS.FromAddress)
		keyid_fromaddress := keyid
		//keyIDBase58, err := soluKeyID.ToBase58Address(version)
		if err != nil{
			dbSession.Rollback()
			return errors.New("Get from address keyid failed")
		}
		//db_serverpengd_vin :=make(map[string]int,len(trx.Vin))
		for i:=0;i<len(trx.Vin);i++{
			//在vin的sigscript中提取redeemscript，并和生成的redeemscript进行比较
			sig_redeemscript,err := getRedeemScriptBySigScript(trx.Vin[i].ScriptSig)
			if err != nil {
				dbSession.Rollback()
				return errors.New("getRedeemScriptBySigScript failed")
			}
			if pubkey_redeemscript != sig_redeemscript {
				dbSession.Rollback()
				return errors.New("redeemscript is error.")
			}
			//根据vin里的txid和vout查找serverpending表里是否有该记录，有则提示该utxo已经使用，没有则插入新纪录
			serverPending,err := model.GlobalDBMgr.PendingTransactionMgr.GetPendingTransactionByVinCount(dbSession,coinInfo.Coinid,trx.Vin[i].PrevOut.Hash.GetHex(),int(trx.Vin[i].PrevOut.N))
			if err != nil {
				dbSession.Rollback()
				return errors.New("Query pending table failed")
			}
			if serverPending != 0 {
				dbSession.Rollback()
				return errors.New("Vin has in pending table")
			}
			//根据vin里的txid调用getrantransaction接口
			rawtx, err := coin.GetRawTransaction(coinInfo.Coinsymbol,coinInfo.Ip, coinInfo.Rpcport,coinInfo.Rpcuser, coinInfo.Rpcpass,trx.Vin[i].PrevOut.Hash.GetHex())
			if err == nil {
				fmt.Println("GetRawTransaction success")
			} else {
				dbSession.Rollback()
				return errors.New("Getrawtransaction failed")
			}

			//将返回的二进制流unpack为transaction对象
			Blob := new(blob.Byteblob)
			Blob.SetHex(rawtx)
			bytesBuf := bytes.NewBuffer(Blob.GetData())
			bufReader := io.Reader(bytesBuf)
			prevout := new(transaction.Transaction)
			prevout.UnPack(bufReader)

			//ok,_,addr := script.ExtractDestination(prevout.Vout[trx.Vin[i].PrevOut.N].ScriptPubKey)
			success,scriptHash := coin.ExtractDestination(prevout.Vout[trx.Vin[i].PrevOut.N].ScriptPubKey)
			if success == false{
				dbSession.Rollback()
				return errors.New("ExtractDestination failed")
			}
			//var soluKeyID keyid.KeyID
			//soluKeyID.SetKeyIDData(scriptHash)
			//keyIDBase58, err := soluKeyID.ToBase58Address(version)
			//if err != nil{
			//	dbSession.Rollback()
			//	return errors.New("coin is not active")
			//}

			//比较vin里的地址和请求里的fromaddress是否一样
			if string(scriptHash) != keyid_fromaddress{
				dbSession.Rollback()
				return errors.New("Address in vin is different from Fromaddress of transaction table")
			}

			db_pending_vin[trx.Vin[i].PrevOut.Hash.GetHex()] = int(trx.Vin[i].PrevOut.N)
			//计算，vin里的utxo的总额，比较vin的value总数和vout+totalfee是否相等
			totalvinvalue = totalvinvalue + prevout.Vout[trx.Vin[i].PrevOut.N].Value
		}

		fromAddr_found := false
		//omni的vout特殊处理
		if coinObj.IsOmni == true {
			//判断是否存在value=546的vout
			found_546 := false
			found_opreturn := false
			for i := 0; i < len(trx.Vout); i++ {
				//获取vout里的目的地址
				//先判断是否存在已6a开头的scriptpubkey，如果存在，则改utxo是omni的opreturn脚本，解析脚本
				//格式6a+14+6f6d6e69+00000000+propertyid_4byte+amount_8byte
				if trx.Vout[i].ScriptPubKey.GetScriptBytes()[0] == 0x6a {
					found_opreturn = true
					//判断第二个字节是否是长度20
					if trx.Vout[i].ScriptPubKey.GetScriptBytes()[1] != 0x14 {
						return errors.New("This trx is a omni trx,but second byte of opreturn is not 0x14.")
					}
					//判断6a之后的脚本是否已“6f6d6e69”开头，后8个字节是amount，amount前4个字节是propertyID
					if hex.EncodeToString(trx.Vout[i].ScriptPubKey.GetScriptBytes()[2:6]) != "6f6d6e69" {
						return errors.New("This trx is a omni trx,but begin of opreturn is not 0x6f6d6e69")
					}

					//10-13字节表示propertyID
					propertyID_str := hex.EncodeToString(trx.Vout[i].ScriptPubKey.GetScriptBytes()[10:14])
					propertyID,err := strconv.Atoi(propertyID_str)
					if err != nil {
						return errors.New("String convert to int64 failed")
					}
					if propertyID != coinObj.OmniPropertyId {
						return errors.New("Propertyid in script is different from config.")
					}

					//获取amount
					amount_str := hex.EncodeToString(trx.Vout[i].ScriptPubKey.GetScriptBytes()[14:])
					amount,err := strconv.ParseInt(amount_str,16,64)
					if err != nil {
						return errors.New("String convert to int64 failed")
					}
					for _, v := range keyID_detail {
						if v != amount {
							return errors.New("OMNI Value in vout is not same with value in todetail of request.")
						}
					}

					//获取omni资产的余额，和amount比较
					balance,_,err := coin.GetBalance(coinInfo.Coinsymbol,coinInfo.Ip,coinInfo.Rpcport,coinInfo.Rpcuser,coinInfo.Rpcpass,trxWS.FromAddress)
					if err != nil {
						return errors.New("OMNI getbalance faild.")
					}

					omni_balance,_ := coin.ToPrecisionAmount(balance,coinObj.Precision)
					if amount > omni_balance {
						return errors.New("OMNI value of vout is bigger than balance of this address.")
					}

					//由于opreturn的vout其他都是null，直接continue
					continue
				}
				success, scriptHash := coin.ExtractDestination(trx.Vout[i].ScriptPubKey)
				if success == false {
					return errors.New("ExtractScript failed.")
				}

				if trx.Vout[i].Value == 546 {
					if found_546 == true {
						return errors.New("OMNI exist more than one vout value is 546.")
					}
					//546的vout里是目的地址
					for k, _ := range keyID_detail {
						if string(scriptHash) != k {
							return errors.New("OMNI address in 546vout is different with address in todetail of request.")
						}
					}
					found_546 = true
					totalvoutvalue = totalvoutvalue + trx.Vout[i].Value

					continue
				}
				//判断vout里的keyid是否和req里的fromaddress的keyid一样
				if string(scriptHash) == keyid_fromaddress {
					fromAddr_found = true
				} else {
					return errors.New("OMNI vout of change is wrong.")
				}

				totalvoutvalue = totalvoutvalue + trx.Vout[i].Value
			}
			if found_546 == false || found_opreturn == false || (fromAddr_found == false && len(trx.Vout) != 2) || (fromAddr_found == true && len(trx.Vout) != 3){
				return errors.New("OMNI vout is wrong.")
			}
		} else {
			for i := 0; i < len(trx.Vout); i++ {
				//获取vout里的目的地址
				success, scriptHash := coin.ExtractDestination(trx.Vout[i].ScriptPubKey)
				if success == false {
					return errors.New("ExtractScript failed.")
				}
				//scriptHash := trx.Vout[i].ScriptPubKey.GetScriptBytes()[2:22]
				//var soluKeyID keyid.KeyID
				//soluKeyID.SetKeyIDData(scriptHash)
				//keyIDBase58, err := soluKeyID.ToBase58Address(version)
				//if err != nil{
				//	dbSession.Rollback()
				//	c.EmitMessage(error_return(res,600004,err.Error()))
				//	return
				//}

				//判断vout里的keyid是否和req里的fromaddress的keyid一样
				if string(scriptHash) == keyid_fromaddress {
					fromAddr_found = true
				}
				//判断目的地址是否在请求的todetail里
				found := false
				for k, v := range keyID_detail {
					if string(scriptHash) == k {
						if v != trx.Vout[i].Value {
							return errors.New("Value in vout is not same with value in todetail of request.")
						}
						found = true
					}
				}

				if fromAddr_found == false && found == false {
					return errors.New("Address in request_todetail is not consistent with address in trxvout.")
				}
				totalvoutvalue = totalvoutvalue + trx.Vout[i].Value
			}
		}

		if coinObj.IsOmni == false {
			if (fromAddr_found == false && len(trx.Vout) != len(trxWS.ToDetail)) || (fromAddr_found == true && (len(trx.Vout) != (len(trxWS.ToDetail) + 1))) {
				return errors.New("Vout of trx is wrong,please check.")
			}
		}
		//判断签名交易的vout+请求里的totalfee和交易里的vin的utxo总额是否相等
		//totalfee,_ := strconv.Atoi(trxWS.TotalFee)
		if totalvoutvalue >= totalvinvalue {
			dbSession.Rollback()
			return errors.New("Total value of vouts >= total value of vins")
		}

		//插入servertransaction表
		isFound, _, err := model.GlobalDBMgr.TransactionMgr.GetTransactionByUuId(dbSession,req.Params[0].TrxUuid)
		if err != nil {
			dbSession.Rollback()
			return errors.New("Query transaction table failed")
		}
		if isFound {
			dbSession.Rollback()
			return errors.New("This transaction has exist")
		}
	}
	//插入serverpending表
	for k,v := range db_pending_vin {
		_,err := model.GlobalDBMgr.PendingTransactionMgr.NewPendingTransaction(dbSession,trxWS.TrxUuid,coinInfo.Coinid,k,v,"","")
		if err != nil {
			dbSession.Rollback()
			return errors.New("Insert pending table failed")
		}
	}

	_, err = model.GlobalDBMgr.TransactionMgr.NewTransaction(dbSession, trxUuid, walletCfg.Walletid, trxWS.CoinId, "",
		-1, trxWS.CreateServerId, trxWS.FromAddress, detailsStr, walletCfg.Needsigcount,
		walletCfg.Fee, walletCfg.Gasprice, walletCfg.Gaslimit, "")
	if err != nil {
		dbSession.Rollback()
		return err
	}

	isFound, trx, err := model.GlobalDBMgr.TransactionMgr.GetTransactionByUuId(dbSession,trxUuid)
	if !isFound || err != nil {
		dbSession.Rollback()
		return err
	}
	trx.Signedtrxs = strings.Join(trxWS.SignedTrx,",")
	arrayStr := make([]string, len(trxWS.SignedServerIds))
	for i := 0; i < len(trxWS.SignedServerIds); i++ {
		arrayStr[i] = strconv.Itoa(trxWS.SignedServerIds[i])
	}
	trx.Signedserverids = strings.Join(arrayStr, ",")
	trx.Feecost = trxWS.TotalFee
	trx.Rawtrxid = trxWS.RawTrxId
	trx.State = 0

	err = model.GlobalDBMgr.TransactionMgr.UpdateTransaction(dbSession, trx)
	if err != nil {
		dbSession.Rollback()
		return err
	}
	dbSession.Commit()
	return nil
}


func getRedeemScriptBySigScript(sigscript script.Script) (string,error) {
	bytesBuf := bytes.NewBuffer(sigscript.GetScriptBytes())
	bufReader := io.Reader(bytesBuf)
	u8, err := serialize.UnPackUint8(bufReader)
	if err != nil {
		return "",err
	}
	if u8 != 0 {
		return "",errors.New("invalid multisig script, not started with 0x0")
	}

	signatureScript := new(script.Script)
	err = signatureScript.UnPack(bufReader)
	if err != nil {
		return "",err
	}

	signatureScriptBytes := signatureScript.GetScriptBytes()
	//if signatureScriptBytes[len(signatureScriptBytes)-1] != 0x1 {
	//	return "",errors.New("invalid signature, not ended with 0x1[ALL]")
	//}
	signatureScriptBytes = signatureScriptBytes[:len(signatureScriptBytes)-1]

	// skip oppushdata
	tmpBufReader := bufio.NewReader(bufReader)
	opPushDataBytes, err := tmpBufReader.Peek(1)
	if err != nil {
		return "",err
	}

	bufReader = io.Reader(tmpBufReader)
	if opPushDataBytes[0] == script.OP_PUSHDATA1 || opPushDataBytes[0] == script.OP_PUSHDATA2 || opPushDataBytes[0] == script.OP_PUSHDATA4 {
		_, err = serialize.UnPackUint8(bufReader)
		if err != nil {
			return "",err
		}
	}

	redeemScript := new(script.Script)
	err = redeemScript.UnPack(bufReader)
	if err != nil {
		return "",err
	}
	return hex.EncodeToString(redeemScript.GetScriptBytes()),nil
}

func getRedeemScriptByPubkeys(needCount int,keyDetail string,coinsymbol string)(string,error) {
	dbSession := model.GetDBEngine().NewSession()
	defer dbSession.Close()

	err := dbSession.Begin()
	if err != nil {
		dbSession.Rollback()
		return "",errors.New("dbsession begin error.")
	}

	keydetail := strings.Split(keyDetail,",")
	pubkeys := make([]string, 0, len(keydetail))
	for _, onedetail := range keydetail {
		//create keys string check key state
		serverid,_ := strconv.Atoi(strings.Split(onedetail,":")[0])
		keyindex,_ := strconv.Atoi(strings.Split(onedetail,":")[1])
		//根据serverid和keyindex获取pubkey
		key, err := model.GlobalDBMgr.PubKeyPoolMgr.GetPubkeyByIdIndex(dbSession, keyindex, serverid)
		if err != nil {
			dbSession.Rollback()
			return "",errors.New("GetPubkeyByIdIndex failed")
		}
		if key == nil {
			dbSession.Rollback()
			return "",errors.New("GetPubkeyByIdIndex failed,can not find keyinfo")
		}

		pubkeys = append(pubkeys, key.Pubkey)
	}
	//Todo: check address and get redeemscript
	redeemscript, err := coin.GetMultiSignRedeemScript(needCount, pubkeys, coinsymbol)
	if err != nil {
		dbSession.Rollback()
		return "",errors.New("GetMultiSignRedeemScript failed.")
	}
	return redeemscript,nil
}

func RequestAndUpdateLocalTrxResult(dbSession *xorm.Session, walletUuid string, trxUuid string) error {
	// web socket request for query_transaction
	var req controller.QueryTransactionRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "query_transaction"
	var param controller.QueryTransactionParamWS
	param.TrxUuid = trxUuid
	req.Params = append(req.Params, param)
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("RequestAndUpdateLocalTrxResult() WSRequest query_transaction")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "transaction", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.QueryTransactionResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}
	trxWS := res.Result

	if trxWS.WalletUuid != walletUuid {
		return errors.New("wallet uuid is not match")
	}

	walletCfg, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletByUUId(dbSession, walletUuid)
	if walletCfg.State != 1 {
		return errors.New("wallet is not active")
	}
	isFound, trx, err := model.GlobalDBMgr.TransactionMgr.GetTransactionByUuId(dbSession, trxUuid)
	if err != nil {
		return err
	} else if !isFound {
		return errors.New("can not get transaction by trx uuid")
	}

	// TODO
	var serverIdSlice []string
	for _, serverId := range trxWS.SignedServerIds {
		serverIdStr := strconv.Itoa(serverId)
		serverIdSlice = append(serverIdSlice, serverIdStr)
	}
	serverIdsStr := strings.Join(serverIdSlice, ",")

	trx.State = trxWS.State
	trx.Rawtrxid = trxWS.RawTrxId
	trx.Signedtrxs = strings.Join(trxWS.SignedTrx, ",")
	trx.Signedserverids = serverIdsStr

	err = model.GlobalDBMgr.TransactionMgr.UpdateTransaction(dbSession, trx)
	if err != nil {
		return err
	}

	if trx.State == 4 || trx.State == 5 {
		err = model.GlobalDBMgr.PendingTransactionMgr.DeletePendingTransactionByTrxUuid(dbSession, trx.Trxuuid)
		if err != nil {
			return err
		}
	}

	return nil
}

func WebSocketDataInit(dbSession *xorm.Session) error {
	// web socket request for list_coins
	var req controller.ListCoinsRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "list_coins"
	req.Params = make([]interface{}, 0)
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("WebSocketDataInit() WSRequest list_coins")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "coin", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.ListCoinsResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}
	// update local coin config
	for _, coinResult := range res.Result {
		isFound, coinCfg, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin2(dbSession, coinResult.CoinId)
		if err != nil {
			return err
		}
		if !isFound {
			err = model.GlobalDBMgr.CoinConfigMgr.InsertCoinWithCoinId(dbSession, coinResult.CoinId, coinResult.CoinSymbol, "", 0,
				"", "", coinResult.State)
			if err != nil {
				return err
			}
		} else {
			if coinResult.CoinSymbol != coinCfg.Coinsymbol {
				return errors.New("different CoinSymbols with the same CoinId")
			}
			if coinResult.State != coinCfg.State {
				err = model.GlobalDBMgr.CoinConfigMgr.UpdateCoinState(dbSession, coinResult.CoinId, coinResult.State)
				if err != nil {
					return err
				}
			}
		}
	}

	type ServerInfoTuple struct {
		ServerId int
		ServerName string
		ServerStartIndex int
	}

	var serverInfoTuples []ServerInfoTuple
	// first time start
	if config.LocalServerId == -1 {
		serverInfoTuple := ServerInfoTuple{config.LocalServerId,config.GlobalConfig.ServerInfoConfig.ServerName,config.GlobalConfig.ServerInfoConfig.StartIndex}
		serverInfoTuples = append(serverInfoTuples, serverInfoTuple)
	} else {
		serverInfos, err := model.GlobalDBMgr.ServerInfoMgr.GetAllServerInfo(dbSession)
		if err != nil {
			return err
		}
		for _, serverInfo := range serverInfos {
			serverInfoTuple := ServerInfoTuple{serverInfo.Serverid, serverInfo.Servername, serverInfo.Serverstartindex}
			serverInfoTuples = append(serverInfoTuples, serverInfoTuple)
		}
	}

	// web socket request for init_pubkey
	var req2 controller.InitPubKeyRequestWS
	req2.Id = 1
	req2.JsonRpc = "2.0"
	req2.Method = "init_pubkey"
	req2.Params = make([]controller.InitPubKeyParamWS, 0)

	for _, serverInfoTuple := range serverInfoTuples {
		serverKeyInfos, err := model.GlobalDBMgr.PubKeyPoolMgr.LoadPubKeysByServerId(dbSession, serverInfoTuple.ServerId)
		if err != nil {
			return err
		}
		var initPubKeyParam controller.InitPubKeyParamWS
		initPubKeyParam.ServerName = serverInfoTuple.ServerName
		if serverInfoTuple.ServerId == config.LocalServerId {
			initPubKeyParam.IsLocalServer = 1
		} else {
			initPubKeyParam.IsLocalServer = 0
		}
		initPubKeyParam.ServerPubkey = ""
		initPubKeyParam.StartIndex = serverInfoTuple.ServerStartIndex
		initPubKeyParam.Keys = make(map[string]string)
		for _, serverKeyInfo := range serverKeyInfos {
			initPubKeyParam.Keys[strconv.Itoa(serverKeyInfo.Keyindex)] = serverKeyInfo.Pubkey
		}
		req2.Params = append(req2.Params, initPubKeyParam)
	}

	msgBytes, err = json.Marshal(req2)
	if err != nil {
		return err
	}
	fmt.Println("WebSocketDataInit() WSRequest init_pubkey")
	fmt.Println("request:", string(msgBytes))
	retString, err = utils.WSRequest(utils.GlobalWsConn, "pubkey", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res2 controller.InitPubKeyResponseWS
	err = json.Unmarshal([]byte(retString), &res2)
	if err != nil {
		return err
	}
	if res2.Error != nil {
		return errors.New(res2.Error.ErrMsg)
	}
	serverId := res2.Result.Serverid

	if config.LocalServerId != -1 {
		if config.LocalServerId != serverId {
			return errors.New("different Server id from Local to Server id from Remote")
		}
	} else {
		config.LocalServerId = serverId
		fmt.Println("LocalServerId:", config.LocalServerId)
	}
	err = model.GlobalDBMgr.PubKeyPoolMgr.UpdateServerIdByServerId(dbSession, -1, serverId)
	if err != nil {
		return err
	}

	// request and update local pubkeys from web socket response of query_pubkeys
	err = RequestAndUpdateLocalPubKeys(dbSession)
	if err != nil {
		return err
	}

	return nil
}

func WebSocketLogin() error {
	// query message
	var req controller.QueryMessageRequestWS
	req.Id = 1
	req.JsonRpc = "2.0"
	req.Method = "query_message"
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("WebSocketLogin() WSRequest query_message")
	fmt.Println("request:", string(msgBytes))
	retString, err := utils.WSRequest(utils.GlobalWsConn, "login", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res controller.QueryMessageResponseWS
	err = json.Unmarshal([]byte(retString), &res)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(res.Error.ErrMsg)
	}

	randomStr := res.Result
	if len(randomStr) != 10 {
		return errors.New("query_message: get invalid random string")
	}

	// sign message
	var rsBytes []byte
	for {
		rsBytes, err = coin.CoinSignTrx('1', utility.Sha256([]byte(randomStr)), 1)
		if err != nil {
			return err
		}
		if len(rsBytes) != 64 {
			return errors.New("invalid r/s lens")
		}
		if rsBytes[32] < 128 {
			break
		}
	}

	verifyOk, err := coin.CoinVerifyTrx('1', 1, utility.Sha256([]byte(randomStr)), rsBytes)
	if err != nil {
		return err
	}
	if !verifyOk {
		return errors.New("verify signature error")
	}
	rBytes := rsBytes[0:32]
	sBytes := rsBytes[32:64]

	// serialize r,s to der encoding
	signedData, err := coin.SerializeDerEncoding(rBytes, sBytes)
	if err != nil {
		return err
	}
	signedDataHex := hex.EncodeToString(signedData)

	var req2 controller.VerifyMessageRequestWS
	req2.Id = 1
	req2.JsonRpc = "2.0"
	req2.Method = "verify_message"
	var verifyMessage controller.VerifyMessageParam
	verifyMessage.ServerId = config.LocalServerId
	verifyMessage.SignedMessage = signedDataHex
	req2.Params = append(req2.Params, verifyMessage)
	msgBytes, err = json.Marshal(req2)
	if err != nil {
		return err
	}
	fmt.Println("WebSocketLogin() WSRequest verify_message")
	fmt.Println("request:", string(msgBytes))
	retString, err = utils.WSRequest(utils.GlobalWsConn, "login", string(msgBytes), utils.GlobalWsTimeOut)
	if err != nil {
		return err
	}
	var res2 controller.VerifyMessageResponseWS
	err = json.Unmarshal([]byte(retString), &res2)
	if err != nil {
		return err
	}
	if res2.Error != nil {
		return errors.New(res2.Error.ErrMsg)
	}
	if !res2.Result {
		return errors.New("login fail")
	}
	return nil
}

func DoWebSocketConnect(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()

	url, origin := "", ""
	timeout := uint(0)
	if !config.GlobalConfig.IsWss {
		url = config.GlobalConfig.WsConfig.Url
		origin = config.GlobalConfig.WsConfig.Origin
		timeout = config.GlobalConfig.WsConfig.TimeOut
	} else {
		url = config.GlobalConfig.WssConfig.Url
		origin = config.GlobalConfig.WssConfig.Origin
		timeout = config.GlobalConfig.WssConfig.TimeOut
	}
	fmt.Println("[WebSocketConnector] DoWebSocketConnect")

	utils.GlobalWsTimeOut = timeout
	// connect to web socket
	for retries :=5; retries >= 0; retries-- {
		var err error
		fmt.Println("[WebSocketConnector] ConnectWebSocket")
		utils.GlobalWsConn, err = utils.ConnectWebSocket(url, origin)
		if err == nil {
			break
		} else if retries == 0 {
			fmt.Println("Can not connect to web socket:", url)
			os.Exit(0)
		}
		waitSeconds := int(math.Pow(float64(2), float64(5-retries)))
		fmt.Println("[WebSocketConnector] Wait", waitSeconds, "Seconds...")
		time.Sleep(time.Duration(waitSeconds) * time.Duration(time.Second))
	}

	fmt.Println("[WebSocketConnector] utils.GlobalWsConnEstablish <- true")
	utils.GlobalWsConnEstablish <- true
	err := utils.WSReadAndDispose()
	// while web socket error, build one new web socket connection
	if err != nil {
		utils.GlobalWsConnReconnect <- true
		utils.GlobalIsLogin = false
	}
}

func DoWebSocketConnMaintain(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()

	for {
		<- utils.GlobalWsConnReconnect

		StartWebSocketConnect()
		// block until web socket established
		<- utils.GlobalWsConnEstablish

		// web socket data init
		dbSession := model.GetDBEngine().NewSession()
		err := dbSession.Begin()
		if err != nil {
			dbSession.Close()
			return
		}
		err = WebSocketDataInit(dbSession)
		if err != nil {
			fmt.Println("WebSocketDataInit()", err)
			dbSession.Rollback()
			dbSession.Close()
			return
		}
		err = dbSession.Commit()
		if err != nil {
			dbSession.Close()
			return
		}

		// web socket login
		err = WebSocketLogin()
		if err != nil {
			fmt.Println("WebSocketLogin()", err)
			return
		}
		utils.GlobalIsLogin = true
	}
}

func main() {
	//if len(os.Args) > 1 && os.Args[1] == "test" {
	//	config.IsTestEnvironment = true
	//	fmt.Println("[main] Run Test Environment")
	//} else {
	//	config.IsTestEnvironment = false
	//	fmt.Println("[main] Run Product Environment")
	//}

	config.IsTestEnvironment = true
	fmt.Println("[main] Run Test Environment")

	fmt.Println("[main] Init()")
	err := Init()
	if err != nil {
		fmt.Println("[main] Init()", err)
		return
	}

	// init local server id
	dbSession := model.GetDBEngine().NewSession()
	fmt.Println("[main] InitLocalServerId()")
	err = InitLocalServerId(dbSession)
	if err != nil {
		fmt.Println("[main] InitLocalServerId()", err)
		dbSession.Close()
		return
	}
	fmt.Println("[main] LocalServerId =", config.LocalServerId)
	dbSession.Close()

	// web socket connection maintainer
	StartWebSocketConnMaintainer()
	StartWebSocketConnect()
	// block until web socket established
	fmt.Println("[main] <- utils.GlobalWsConnEstablish")
	<- utils.GlobalWsConnEstablish

	// web socket data init
	dbSession = model.GetDBEngine().NewSession()
	err = dbSession.Begin()
	if err != nil {
		dbSession.Close()
		return
	}

	fmt.Println("[main] WebSocketDataInit()")
	err = WebSocketDataInit(dbSession)
	if err != nil {
		fmt.Println("[main] WebSocketDataInit()", err)
		dbSession.Rollback()
		dbSession.Close()
		return
	}
	err = dbSession.Commit()
	if err != nil {
		dbSession.Close()
		return
	}

	// web socket login
	fmt.Println("[main] WebSocketLogin()")
	err = WebSocketLogin()
	if err != nil {
		fmt.Println("[main] WebSocketLogin()", err)
		return
	}

	RegisterUrlRouter()
	StartSessionMaintainer()
	StartTransactionMaintainer()
	StartTaskPersistenceProcessor()
	StartTaskProcessor()

	// http
	if !config.GlobalConfig.IsHttps {
		Run(config.GlobalConfig.HttpConfig.EndPoint, config.GlobalConfig.HttpConfig.CharSet)
	} else {
	// https
		RunTLS(config.GlobalConfig.HttpsConfig.EndPoint, config.GlobalConfig.HttpsConfig.CertFile,
			config.GlobalConfig.HttpsConfig.KeyFile, config.GlobalConfig.HttpsConfig.CharSet)
	}
}
