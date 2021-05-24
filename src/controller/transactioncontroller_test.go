package controller

import (
	"testing"
	"fmt"
	"strconv"
	"strings"
	"utils"
)

func TestCreateTrxController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, acc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param CreateTrxParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.FromWalletId = 23

	// btc test p2pkh address
	//param.ToAddr = "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV"
	// btc test p2sh address
	//param.ToAddr = "2NBa9HHHrRM116kZFfmniu24Ymh5a2z6wCe"

	// ltc test p2pkh address
	//param.ToAddr = "myBBLuSWz6XVbMGp4wXLiFHtYcoMdEHVvU"
	// ltc test p2sh address
	//param.ToAddr = "Qd474mK5s7fUYMA1DxJWZSQQsbSq6maeug"

	// bch test p2pkh address
	//param.ToAddr = "bchtest:pz9er9lrgstcc8lqmjpxls8ldwk20d5zhvjltelweu"
	// bch test p2sh address
	//param.ToAddr = ""

	// ub test p2pkh address
	//param.ToAddr = "mmCYpBtkQ5fbEohtFym6BHVmnGP7VC23ot"
	// ub test p2sh address
	//param.ToAddr = ""

	// omni test p2pkh address
	param.ToAddr = "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV"
	// omni test p2sh address
	//param.ToAddr = ""

	param.Amount = "0.001"
	param.Fee = "0.001"
	param.GasPrice = ""
	param.GasLimit = ""
	funcNameStr := "transfer"
	sessionIdStr := param.SessionId
	operatorIdStr := strconv.Itoa(param.OperateId)
	fromWalletIdStr := strconv.Itoa(param.FromWalletId)
	toAddrStr := param.ToAddr
	amountStr := param.Amount
	feeStr := param.Fee
	gasPriceStr := param.GasPrice
	gasLimitStr := param.GasLimit
	sigData := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, fromWalletIdStr, toAddrStr, amountStr, feeStr, gasPriceStr, gasLimitStr}, ",")
	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/transaction", "transfer", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestConfirmTrxControllerDevice1(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, acc2)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ConfirmTrxParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 42
	sigData := "confirm," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc2.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/transaction", "confirm", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestConfirmTrxControllerDevice2Acc3(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc2.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc2.login(auid, acc3)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ConfirmTrxParam
	param.SessionId = sid
	mgmtid, err := tc2.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 30
	sigData := "confirm," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc3.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/transaction", "confirm", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestConfirmTrxControllerDevice2Acc4(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc2.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc2.login(auid, acc4)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ConfirmTrxParam
	param.SessionId = sid
	mgmtid, err := tc2.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 30
	sigData := "confirm," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc4.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/transaction", "confirm", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestConfirmTrxControllerDevice3Acc5(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc3.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc3.login(auid, acc5)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ConfirmTrxParam
	param.SessionId = sid
	mgmtid, err := tc3.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 7
	sigData := "confirm," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc5.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/transaction", "confirm", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestConfirmTrxControllerDevice3Acc6(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc3.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc3.login(auid, acc6)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param ConfirmTrxParam
	param.SessionId = sid
	mgmtid, err := tc3.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 7
	sigData := "confirm," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc6.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/transaction", "confirm", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestRevokeTrxController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, acc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	var param RevokeTrxParam
	param.SessionId = sid
	mgmtid, err := tc.get_mgmtid(param.SessionId, 4)
	if err != nil {
		fmt.Println("get mgmtid error:", err.Error())
	}
	param.OperateId = int(mgmtid)
	param.TrxId = 14
	sigData := "revoke," + sid + "," + strconv.Itoa(param.OperateId) + "," + strconv.Itoa(param.TrxId)

	sigRes, err := utils.RsaSignWithSha1Hex(sigData, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sigRes

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/transaction", "revoke", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	fmt.Println("result:", res.Result)
}

func TestGetTransactionController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)

	var param GetWalletTrxsParam

	auid, err := tc.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc.login(auid, acc)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	param.SessionId = sid
	param.WalletId = []int{}
	param.CoinId = []int{}
	param.ServerId = []int{}
	param.AcctId = []int{}
	param.TrxTime = [2]string{}
	param.State = []int{}
	param.OffSet = 0
	param.Limit = 100

	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/transaction", "get_wallet_trxs", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	result := res.Result.(map[string]interface{})
	total, _ := result["total"]
	fmt.Println("total:", total)
	trxs, _ := result["trxs"].([]interface{})
	for _, trx := range trxs {
		res := trx.(map[string]interface{})
		trxid, _ := res["trxid"]
		trxuuid, _ := res["trxuuid"]
		rawtrxid, _ := res["rawtrxid"]
		walletid, _ := res["walletid"]
		coinid, _ := res["coinid"]
		contractaddr, _ := res["contractaddr"]
		acctid, _ := res["acctid"]
		serverid, _ := res["serverid"]
		fromaddr, _ := res["fromaddr"]
		todetails, _ := res["todetails"]
		feecost, _ := res["feecost"]
		trxtime, _ := res["trxtime"]
		needconfirm, _ := res["needconfirm"]
		confirmed, _ := res["confirmed"]
		acctconfirmed, _ := res["acctconfirmed"]
		state, _ := res["state"]

		fmt.Println("trxid:", trxid)
		fmt.Println("trxuuid:", trxuuid)
		fmt.Println("rawtrxid:", rawtrxid)
		fmt.Println("walletid:", walletid)
		fmt.Println("coinid:", coinid)
		fmt.Println("contractaddr:", contractaddr)
		fmt.Println("acctid:", acctid)
		fmt.Println("serverid:", serverid)
		fmt.Println("fromaddr:", fromaddr)
		fmt.Println("todetails:", todetails)
		fmt.Println("feecost:", feecost)
		fmt.Println("trxtime:", trxtime)
		fmt.Println("needconfirm:", needconfirm)
		fmt.Println("confirmed:", confirmed)
		fmt.Println("acctconfirmed:", acctconfirmed)
		fmt.Println("state:", state)
	}
}


