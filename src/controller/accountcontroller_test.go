package controller

import (
	"testing"
	"fmt"
	"strconv"
	"utils"
)

func TestListAccountsController(t *testing.T) {
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

	var param ListAccountsParam
	param.SessionId = sid
	param.State = []int{0,1,2}
	param.Offset = 0
	param.Limit = 10
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/account", "list_accounts", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	result := res.Result.(map[string]interface{})
	total, _ := result["total"]
	fmt.Println("total:", total)
	accts, _ := result["accts"].([]interface{})
	for _, acct := range accts {
		res := acct.(map[string]interface{})
		acctid, _ := res["acctid"]
		cellnumber, _ := res["cellnumber"]
		realname, _ := res["realname"]
		idcard, _ := res["idcard"]
		regtime, _ := res["regtime"]
		state, _ := res["state"]
		fmt.Println("acctid:", acctid)
		fmt.Println("cellnumber:", cellnumber)
		fmt.Println("realname:", realname)
		fmt.Println("idcard:", idcard)
		fmt.Println("regtime:", regtime)
		fmt.Println("state:", state)
	}
}

func TestGetAccountController(t *testing.T) {
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

	var param GetAccountParam
	param.SessionId = sid
	param.AcctId = []int{2}
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/account", "get_account", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	results := res.Result.([]interface{})
	for _, result := range results {
		res := result.(map[string]interface{})
		acctid, _ := res["acctid"]
		cellnumber, _ := res["cellnumber"]
		realname, _ := res["realname"]
		idcard, _ := res["idcard"]
		regtime, _ := res["regtime"]
		state, _ := res["state"]
		fmt.Println("acctid:", acctid)
		fmt.Println("cellnumber:", cellnumber)
		fmt.Println("realname:", realname)
		fmt.Println("idcard:", idcard)
		fmt.Println("regtime:", regtime)
		fmt.Println("state:", state)
	}
}

func TestModifyAcctController(t *testing.T) {
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

	var param ModifyAcctParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 2)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.MgmtId = int(mgmtid)
	param.AcctId = 2
	param.Walletid = []int{1}
	param.State = 0
	sigData := "modify_acct," + param.SessionId + "," + strconv.Itoa(param.MgmtId) + "," + strconv.Itoa(param.AcctId) + "," +
		utils.IntArrayToString(param.Walletid)  + "," + strconv.Itoa(param.State)
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, adminAcc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/account", "modify_acct", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestChangeAccountStateController(t *testing.T) {
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

	var param ChangeAccountStateParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 2)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.MgmtId = int(mgmtid)
	param.AcctId = 2
	param.State = 1
	sigData := "change_acct_state," + param.SessionId + "," + strconv.Itoa(param.MgmtId) + "," + strconv.Itoa(param.AcctId) + "," + strconv.Itoa(param.State)
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, adminAcc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/account", "change_acct_state", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestGetAccountWalletsController(t *testing.T) {
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

	var param GetAccountWalletsParam
	param.SessionId = sid
	param.AcctId = 2
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/account", "get_acct_wallets", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	results := res.Result.([]interface{})
	for _, result := range results {
		res := result.(map[string]interface{})
		walletid, _ := res["walletid"]
		walletname, _ := res["walletname"]
		coinid, _ := res["coinid"]
		address, _ := res["address"]
		state, _ := res["state"]
		fmt.Println("walletid:", walletid)
		fmt.Println("walletname:", walletname)
		fmt.Println("coinid:", coinid)
		fmt.Println("address:", address)
		fmt.Println("state:", state)
	}
}

