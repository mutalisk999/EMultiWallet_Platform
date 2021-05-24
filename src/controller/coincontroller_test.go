package controller

import (
	"testing"
	"fmt"
	"strconv"
	"utils"
)

func TestListCoinsController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, adminAcc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ListCoinsParam
	param.SessionId = sid

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/coin", "list_coins", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	results := res.Result.([]interface{})
	for _, result := range results {
		res := result.(map[string]interface{})
		coinid, _ := res["coinid"]
		coinsymbol, _ := res["coinsymbol"]
		ip, _ := res["ip"]
		rpcport, _ := res["rpcport"]
		rpcuser, _ := res["rpcuser"]
		rpcpass, _ := res["rpcpass"]
		state, _ := res["state"]
		fmt.Println("coinid:", coinid)
		fmt.Println("coinsymbol:", coinsymbol)
		fmt.Println("ip:", ip)
		fmt.Println("rpcport:", rpcport)
		fmt.Println("rpcuser:", rpcuser)
		fmt.Println("rpcpass:", rpcpass)
		fmt.Println("state:", state)
	}
}

func TestRefreshCoinsController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, adminAcc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param RefreshCoinsParam
	param.SessionId = sid

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/coin", "refresh_coins", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	results := res.Result.([]interface{})
	for _, result := range results {
		res := result.(map[string]interface{})
		coinid, _ := res["coinid"]
		coinsymbol, _ := res["coinsymbol"]
		ip, _ := res["ip"]
		rpcport, _ := res["rpcport"]
		rpcuser, _ := res["rpcuser"]
		rpcpass, _ := res["rpcpass"]
		state, _ := res["state"]
		fmt.Println("coinid:", coinid)
		fmt.Println("coinsymbol:", coinsymbol)
		fmt.Println("ip:", ip)
		fmt.Println("rpcport:", rpcport)
		fmt.Println("rpcuser:", rpcuser)
		fmt.Println("rpcpass:", rpcpass)
		fmt.Println("state:", state)
	}
}

func TestModifyCoinController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, adminAcc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ModifyCoinParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 6)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.MgmtId = int(mgmtid)
	param.CoinId = 1
	param.CoinSymbol = "BTC"
	param.Ip = "192.168.1.124"
	param.RpcPort = 10001
	param.RpcUser = "test"
	param.RpcPass = "test"

	sigData := "modify_coin," + param.SessionId + "," + strconv.Itoa(param.MgmtId) + "," + strconv.Itoa(param.CoinId) + "," + param.CoinSymbol + "," + param.Ip + "," +
		strconv.Itoa(param.RpcPort) + "," + param.RpcUser + "," + param.RpcPass
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, adminAcc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes
	params = append(params, param)
	fmt.Println(sigRes)

	res, err := tc.doHttpJsonRpcCallType1("/apis/coin", "modify_coin", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestGetCoinController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, adminAcc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param GetCoinParam
	param.SessionId = sid
	param.CoinId = 1

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/coin", "get_coin", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	result := res.Result.(map[string]interface{})
	coinid, _ := result["coinid"]
	coinsymbol, _ := result["coinsymbol"]
	ip, _ := result["ip"]
	rpcport, _ := result["rpcport"]
	rpcuser, _ := result["rpcuser"]
	rpcpass, _ := result["rpcpass"]
	state, _ := result["state"]
	fmt.Println("coinid:", coinid)
	fmt.Println("coinsymbol:", coinsymbol)
	fmt.Println("ip:", ip)
	fmt.Println("rpcport:", rpcport)
	fmt.Println("rpcuser:", rpcuser)
	fmt.Println("rpcpass:", rpcpass)
	fmt.Println("state:", state)
}


