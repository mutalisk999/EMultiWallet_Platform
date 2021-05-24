package controller

import (
	"testing"
	"fmt"
	"strconv"
	"utils"
)

// 加密机1 管理员注册
func TestUserRegisterControllerDevice1Admin(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678901"
	param.RealName = "加密机1管理员"
	param.IdCard = "110105196812272168"
	param.VerifyCodeId = "8ncABQkdr2bntZXb4UJd"
	param.VerifyCode = "682199"
	param.Pubkey = adminAcc.pubkeypem
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机1 用户a注册
func TestUserRegisterControllerDevice1Usera(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678902"
	param.RealName = "加密机1用户a"
	param.IdCard = "110105196812272169"
	param.VerifyCodeId = "7HBhcr0YZ97kDcAqUQgY"
	param.VerifyCode = "189229"
	param.Pubkey = acc.pubkeypem
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机1 用户b注册
func TestUserRegisterControllerDevice1Userb(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678903"
	param.RealName = "加密机1用户b"
	param.IdCard = "110105196812272170"
	param.VerifyCodeId = "hsQhzm0M0vkuCKRboOVh"
	param.VerifyCode = "920724"
	param.Pubkey = acc2.pubkeypem
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机2 管理员注册
func TestUserRegisterControllerDevice2Admin(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678901"
	param.RealName = "加密机2管理员"
	param.IdCard = "110105196812272168"
	param.VerifyCodeId = "mHUcQsDKTJ41Fq7cCExe"
	param.VerifyCode = "233338"
	param.Pubkey = adminAcc2.pubkeypem
	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机2 用户a注册
func TestUserRegisterControllerDevice2Usera(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678902"
	param.RealName = "加密机2用户a"
	param.IdCard = "110105196812272169"
	param.VerifyCodeId = "Zug74DeqqJE0FuyRBTR4"
	param.VerifyCode = "189389"
	param.Pubkey = acc3.pubkeypem
	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机2 用户b注册
func TestUserRegisterControllerDevice2Userb(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678903"
	param.RealName = "加密机2用户b"
	param.IdCard = "110105196812272170"
	param.VerifyCodeId = "fN0yqCSdzKXm9StMEno9"
	param.VerifyCode = "472497"
	param.Pubkey = acc4.pubkeypem
	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机3 管理员注册
func TestUserRegisterControllerDevice3Admin(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678901"
	param.RealName = "加密机3管理员"
	param.IdCard = "110105196812272168"
	param.VerifyCodeId = "r0DcaOCcNZYHHMHMKTFV"
	param.VerifyCode = "774837"
	param.Pubkey = adminAcc3.pubkeypem
	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机3 用户a注册
func TestUserRegisterControllerDevice3Usera(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678902"
	param.RealName = "加密机3用户a"
	param.IdCard = "110105196812272169"
	param.VerifyCodeId = "xjo7D1HNEA97N7YyCYc8"
	param.VerifyCode = "931700"
	param.Pubkey = acc5.pubkeypem
	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

// 加密机3 用户b注册
func TestUserRegisterControllerDevice3Userb(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserRegisterParam
	param.CellNumber = "12345678903"
	param.RealName = "加密机3用户b"
	param.IdCard = "110105196812272170"
	param.VerifyCodeId = "aNj4OCuYF33OxCVmcF2z"
	param.VerifyCode = "272113"
	param.Pubkey = acc6.pubkeypem
	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/user", "user_register", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestUserLoginController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserLoginParam
	loginid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get loginid error:", err.Error())
	}
	param.LoginId = int(loginid)
	param.Pubkey = adminAcc.pubkeypem
	sigData := "user_login," + strconv.Itoa(param.LoginId)
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, adminAcc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes
	params = append(params, param)
	fmt.Println(sigRes)

	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_login", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	sessionid, _ := res.Result.(map[string]interface{})["sessionid"]
	acctid, _ := res.Result.(map[string]interface{})["acctid"]
	usertype, _ := res.Result.(map[string]interface{})["usertype"]

	fmt.Println("sessionid:", sessionid)
	fmt.Println("acctid:", acctid)
	fmt.Println("usertype:", usertype)
}

func TestUserLoginController2(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserLoginParam
	loginid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get loginid error:", err.Error())
	}
	param.LoginId = int(loginid)
	param.Pubkey = acc.pubkeypem
	sigData := "user_login," + strconv.Itoa(param.LoginId)
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes
	params = append(params, param)
	fmt.Println(sigRes)

	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_login", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	sessionid, _ := res.Result.(map[string]interface{})["sessionid"]
	acctid, _ := res.Result.(map[string]interface{})["acctid"]
	usertype, _ := res.Result.(map[string]interface{})["usertype"]

	fmt.Println("sessionid:", sessionid)
	fmt.Println("acctid:", acctid)
	fmt.Println("usertype:", usertype)
}

func TestUserLogoutController(t *testing.T) {
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

	var param UserLogoutParam
	param.SessionId = sid
	sigData := "user_logout," + param.SessionId
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, adminAcc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes
	params = append(params, param)
	fmt.Println(sigRes)

	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_logout", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestUserGetInfoController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserGetInfoParam
	param.Pubkey = adminAcc.pubkeypem
	params = append(params, param)

	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_getinfo", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	cellnumber, _ := res.Result.(map[string]interface{})["cellnumber"]
	realname, _ := res.Result.(map[string]interface{})["realname"]
	idcard, _ := res.Result.(map[string]interface{})["idcard"]
	state, _ := res.Result.(map[string]interface{})["state"]
	regtime, _ := res.Result.(map[string]interface{})["regtime"]

	fmt.Println("cellnumber:", cellnumber)
	fmt.Println("realname:", realname)
	fmt.Println("idcard:", idcard)
	fmt.Println("state:", state)
	fmt.Println("regtime:", regtime)
}

func TestUserGetInfoController2(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param UserGetInfoParam
	param.Pubkey = acc.pubkeypem
	params = append(params, param)

	res, err := tc.doHttpJsonRpcCallType1("/apis/user", "user_getinfo", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	cellnumber, _ := res.Result.(map[string]interface{})["cellnumber"]
	realname, _ := res.Result.(map[string]interface{})["realname"]
	idcard, _ := res.Result.(map[string]interface{})["idcard"]
	state, _ := res.Result.(map[string]interface{})["state"]
	regtime, _ := res.Result.(map[string]interface{})["regtime"]

	fmt.Println("cellnumber:", cellnumber)
	fmt.Println("realname:", realname)
	fmt.Println("idcard:", idcard)
	fmt.Println("state:", state)
	fmt.Println("regtime:", regtime)
}

