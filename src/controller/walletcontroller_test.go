package controller

import (
	"fmt"
	"testing"
	"math/rand"
	"strconv"
	"utils"
)

//func (c *testClient) modifywalletcfg(acc account, sid string, mgmtid int, wid int, name string, daddr string, nsc int, fee string, gasprice string, gaslimit string, sids []int, sta int) string {
//	sigdata := "modify_wallet," + sid + "," + strconv.Itoa(mgmtid) + "," + strconv.Itoa(wid) + "," + name + "," + daddr + "," + strconv.Itoa(nsc) + "," + fee + "," + gasprice + "," + gaslimit + "," + utils.IntArrayToString(sids) + "," + strconv.Itoa(sta)
//	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikeyhex)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(sig_res)
//	pa := make([]map[string]interface{}, 0)
//	var pa0 map[string]interface{}
//	pa0 = make(map[string]interface{})
//	pa0["sessionid"] = sid
//	pa0["mgmtid"] = mgmtid
//	pa0["walletid"] = wid
//	pa0["walletname"] = name
//	pa0["destaddress"] = daddr
//	pa0["needsigcount"] = nsc
//	pa0["fee"] = fee
//	pa0["gasprice"] = gasprice
//	pa0["gaslimit"] = gaslimit
//	pa0["siguserid"] = sids
//	pa0["state"] = 0
//	pa0["signature"] = sig_res
//	pa = append(pa, pa0)
//	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "modify_wallet", pa)
//	if err != nil {
//		return err.Error()
//	}
//	if res.Error != nil {
//		return res.Error.Message
//	}
//	return "Modify" + name + "Success"
//}

func (c *testClient) modifywalletcfg(acc account, sid string, mgmtid int, wid int, nsc int, sids []int) string {
	sigdata := "modify_wallet," + sid + "," + strconv.Itoa(mgmtid) + "," + strconv.Itoa(wid) + "," + strconv.Itoa(nsc) + "," + utils.IntArrayToString(sids)
	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sig_res)
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["mgmtid"] = mgmtid
	pa0["walletid"] = wid
	pa0["needsigcount"] = nsc
	pa0["siguserid"] = sids
	pa0["signature"] = sig_res
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "modify_wallet", pa)
	if err != nil {
		return err.Error()
	}
	if res.Error != nil {
		return res.Error.Message
	}
	return "Modify" + strconv.Itoa(wid) + "Success"
}

func (c *testClient) createwalletcfg(acc account, sid string, mgmtid int, name string, coinid int, keysigserverid []int,
	needkeysigcount int, daddr string, nsc int, fee string, gasprice string, gaslimit string, sids []int, sta int) string {
	sigdata := "create_wallet," + sid + "," + strconv.Itoa(mgmtid) + "," + name + "," + strconv.Itoa(coinid) + "," + utils.IntArrayToString(keysigserverid) + "," + strconv.Itoa(needkeysigcount) + "," + daddr + "," + strconv.Itoa(nsc) + "," + fee + "," + gasprice + "," + gaslimit + "," + utils.IntArrayToString(sids) + "," + strconv.Itoa(sta)

	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sig_res)
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["mgmtid"] = mgmtid
	pa0["walletname"] = name
	pa0["coinid"] = coinid
	tmpSlice := make([]int,0)
	tmpSlice = append(tmpSlice, keysigserverid...)
	pa0["keysigserverid"] = tmpSlice
	pa0["needkeysigcount"] = needkeysigcount
	pa0["destaddress"] = daddr
	pa0["needsigcount"] = nsc
	pa0["fee"] = fee
	pa0["gasprice"] = gasprice
	pa0["gaslimit"] = gaslimit
	pa0["siguserid"] = sids
	pa0["state"] = 0
	pa0["signature"] = sig_res
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "create_wallet", pa)
	if err != nil {
		return err.Error()
	}
	if res.Error != nil {
		return res.Error.Message
	}
	return "Create" + name + "Success"
}

func (c *testClient) GetWalletcfg(wid int, sid string) error {
	var pa0 map[string]interface{}

	pa := make([]map[string]interface{}, 0)
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["walletid"] = wid
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "get_wallet", pa)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
		return err
	}

	result := res.Result.(map[string]interface{})
	walletid, _ := result["walletid"]
	walletuuid, _ := result["walletuuid"]
	walletname, _ := result["walletname"]
	coinid, _ := result["coinid"]
	serverkeys, _ := result["serverkeys"]
	createserver, _ := result["createserver"]
	keycount, _ := result["keycount"]
	needkeysigcount, _ := result["needkeysigcount"]
	address, _ := result["address"]
	destaddress, _ := result["destaddress"]
	needsigcount, _ := result["needsigcount"]
	fee, _ := result["fee"]
	gasprice, _ := result["gasprice"]
	gaslimit, _ := result["gaslimit"]
	siguserid, _ := result["siguserid"]
	state, _ := result["state"]
	fmt.Println("walletid:", walletid)
	fmt.Println("walletuuid:", walletuuid)
	fmt.Println("walletname:", walletname)
	fmt.Println("coinid:", coinid)
	fmt.Println("serverkeys:", serverkeys)
	fmt.Println("createserver:", createserver)
	fmt.Println("keycount:", keycount)
	fmt.Println("needkeysigcount:", needkeysigcount)
	fmt.Println("address:", address)
	fmt.Println("destaddress:", destaddress)
	fmt.Println("needsigcount:", needsigcount)
	fmt.Println("fee:", fee)
	fmt.Println("gasprice:", gasprice)
	fmt.Println("gaslimit:", gaslimit)
	fmt.Println("siguserid:", siguserid)
	fmt.Println("state:", state)

	return nil
}

func (c *testClient) listwalletcfg(acc account, sid string, coinid []int, state []int, offset int, limit int) {
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["coinid"] = coinid
	pa0["state"] = state
	pa0["offset"] = offset
	pa0["limit"] = limit
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "list_wallets", pa)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
		return
	}
	re := res.Result.(map[string]interface{})
	fmt.Println("total:", re["total"])
	wallets := re["wallets"]
	//fmt.Println(reflect.TypeOf(wallets))
	if wallets != nil {
		ress := wallets.([]interface{})
		//fmt.Println(ress)
		for _, va := range ress {
			res := va.(map[string]interface{})
			walletid, _ := res["walletid"]
			walletuuid, _ := res["walletuuid"]
			walletname, _ := res["walletname"]
			coinid, _ := res["coinid"]
			coinsymbol, _ := res["coinsymbol"]
			balance, _ := res["balance"]
			feebalance, _ := res["feebalance"]
			address, _ := res["address"]
			state, _ := res["state"]
			fmt.Println("walletid:", walletid)
			fmt.Println("walletuuid:", walletuuid)
			fmt.Println("walletname:", walletname)
			fmt.Println("coinid:", coinid)
			fmt.Println("coinsymbol:", coinsymbol)
			fmt.Println("balance:", balance)
			fmt.Println("feebalance:", feebalance)
			fmt.Println("address:", address)
			fmt.Println("state:", state)
		}
	}
}

func TestModifyWalletController2(t *testing.T) {
	InitTest()
	auid, err := tc2.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc2.login(auid, adminAcc2)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	wid := 15
	tc2.GetWalletcfg(wid, sid)
	mgmtid, err := tc2.get_mgmtid(sid, 3)
	mid := int(mgmtid)
	//addname := rand.Int63()
	//walletna := "tw" + strconv.FormatInt(addname, 10)
	fmt.Println(tc2.modifywalletcfg(adminAcc2, sid, mid, wid,2, []int{2,3}))
	tc2.GetWalletcfg(wid, sid)
}

func TestModifyWalletController3(t *testing.T) {
	InitTest()
	auid, err := tc3.get_authid()
	if err != nil {
		fmt.Println("get auid error:", err.Error())
	}
	fmt.Println(auid)
	sid, err := tc3.login(auid, adminAcc3)
	if err != nil {
		fmt.Println("login error:", err.Error())
	}
	fmt.Println(sid)

	wid := 16
	tc3.GetWalletcfg(wid, sid)
	mgmtid, err := tc3.get_mgmtid(sid, 3)
	mid := int(mgmtid)
	//addname := rand.Int63()
	//walletna := "tw" + strconv.FormatInt(addname, 10)
	fmt.Println(tc3.modifywalletcfg(adminAcc3, sid, mid, wid,2, []int{2,3}))
	tc3.GetWalletcfg(wid, sid)
}

func TestGetWalletController(t *testing.T) {
	InitTest()
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
	tc.GetWalletcfg(1, sid)
}

func TestListWalletController(t *testing.T) {
	InitTest()
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
	tc.listwalletcfg(acc, sid, []int{1}, []int{0, 1, 2}, 0, 100)
}

//func TestWalletControllerListWallet(t *testing.T) {
//	InitTest()
//	auid, err := tc.get_authid()
//	if err != nil {
//		fmt.Println("get auid error:", err.Error())
//	}
//	fmt.Println(auid)
//	sid, err := tc.login(auid, adminAcc)
//	if err != nil {
//		fmt.Println("login error:", err.Error())
//	}
//	fmt.Println(sid)
//	mgmtid, err := tc.get_mgmtid(sid, 3)
//	mid := int(mgmtid)
//	addname := rand.Int63()
//	walletna := "tw" + strconv.FormatInt(addname, 10)
//	fmt.Println(tc.createwalletcfg(adminAcc, sid, mid, walletna, 1, []int{1,2,3}, 2, walletna, 1,  "0.001", "", "", []int{1}, 0))
//}

func TestCreateWalletController(t *testing.T) {
	InitTest()
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
	mgmtid, err := tc.get_mgmtid(sid, 3)
	mid := int(mgmtid)
	addname := rand.Int63()
	walletna := "tw" + strconv.FormatInt(addname, 10)
	fmt.Println(tc.createwalletcfg(adminAcc, sid, mid, walletna, 4, []int{37,38,39}, 2, "", 2, "0.001", "", "", []int{2,3}, 0))
}

//func TestWalletControllerLogin(t *testing.T) {
//	InitTest()
//	//var tc testClient
//	//tc.url = "http://localhost:18080"
//	//var acc account
//	//acc.prikeyhex = "30820277020100300d06092a864886f70d0101010500048202613082025d02010002818100ece10d27f9eaecbdd5268e865715521c12634c7a1fe70fd4128147f3a5cfc86fd50660266daba2affb760d815cc132a2cafb904bde1c909e12cfa89eaf841abcce1936524c783ce5c6cd8b8b60b47050e14d17a9232e061b21d7b60746780e5bb7142b437ad1867e439deebfb6eea55892d98d95e29c8e2cbfe96a8141962665020301000102818100c3ce000f0471e1c1c558bac56764935beb0333eb5b45a77ad8d50ec1e3550f4d09dcdc4bc7a9f1afe07fa40843c0db775fac7489920f30a7c9cae78a4c71399b32e5c3405bef3ba8ecc09f39ee9444b488ab46bd1db6d35f66b53033d067934a3aa9afe4fe820f7e2877dc570b12139bf67b7ae130012d2166b0d0162cfa0081024100f9da26c8b25c8f508da7fa4961bfe917a9b5f296352f75cb10b52149f5ee4366ba93284a4f62797d748ba968da0bfdff575d3993d8b4888d5e476d2d001134a1024100f2b52e5de3a8c81dccf4cc9f76e9e3b92e1082b0e3e42e26c8a09cfc55151adf5f258440cb2fed1bea4461cd53b5c57e347d5b4be1a9513bf4a0562f45e89745024049cfc1bf6eb9db78dfb06d477a8238794e20bb5ed6ad83ae0eec83be16d261c51a7f58dac323036325a810cf320372f1193d28425e35cf557c31bbb61c9bf8a1024100d31b5cd535ac1d442e7dcf6df4d4c3cad29271dac29e3ab66953aab8a098d6e2ebb38a31818a9246f1a1c120036bdbf84657b30d719f009591d5d4a111c199b1024010dd12bb1707813c3384d18195f346e73640511aa9930945ec49f596a3b5523b913eaff3b87f9b1999a880b777912ebd82fa726cda07e316c4d2a9637b22a8cc"
//	//acc.pubkeypem = "-----BEGIN 公钥-----\nMIGJAoGBAOzhDSf56uy91SaOhlcVUhwSY0x6H+cP1BKBR/Olz8hv1QZgJm2roq/7\ndg2BXMEyosr7kEveHJCeEs+onq+EGrzOGTZSTHg85cbNi4tgtHBQ4U0XqSMuBhsh\n17YHRngOW7cUK0N60YZ+Q53uv7bupViS2Y2V4pyOLL/paoFBliZlAgMBAAE=\n-----END 公钥-----"
//	auid, err := tc.get_authid()
//	if err != nil {
//		fmt.Println("get auid error:", err.Error())
//	}
//	fmt.Println(auid)
//	sid, err := tc.login(auid, acc)
//	if err != nil {
//		fmt.Println("login error:", err.Error())
//	}
//	fmt.Println(sid)
//}

func (c *testClient) ChangeWalletState(acc account, sid string, mid int, wid int, sta int) string {
	sigdata := "change_wallet_state," + sid + "," + strconv.Itoa(mid) + "," + strconv.Itoa(wid) + "," + strconv.Itoa(sta)
	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikeyhex)
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(sig_res)
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["mgmtid"] = mid
	pa0["walletid"] = wid
	pa0["state"] = sta
	pa0["signature"] = sig_res
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "change_wallet_state", pa)
	if err != nil {
		return err.Error()
	}
	if res.Error != nil {
		return res.Error.Message
	}
	return "change_wallet_state " + strconv.Itoa(wid) + " Success"
}

//func TestWalletControllerChangeState(t *testing.T) {
//	InitTest()
//	auid, err := tc.get_authid()
//	if err != nil {
//		fmt.Println("get auid error:", err.Error())
//	}
//	fmt.Println(auid)
//	sid, err := tc.login(auid, acc)
//	if err != nil {
//		fmt.Println("login error:", err.Error())
//	}
//	fmt.Println(sid)
//	mgmtid, err := tc.get_mgmtid(sid, 3)
//	mid := int(mgmtid)
//	tc.GetWalletcfg(4, sid)
//	fmt.Println(tc.ChangeWalletState(acc, sid, mid, 4, 0))
//	tc.GetWalletcfg(4, sid)
//	//mgmtid, err = tc.get_mgmtid(sid, 3)
//	//mid = int(mgmtid)
//	//fmt.Println(tc.ChangeWalletState(acc, sid, mid, 4,1))
//	tc.GetWalletcfg(4, sid)
//}

func (c *testClient) DeleteWalletCfg(acc account, sid string, mid int, wid int) string {
	sigdata := "delete_wallet," + sid + "," + strconv.Itoa(mid) + "," + strconv.Itoa(wid)
	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikeyhex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sig_res)
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["mgmtid"] = mid
	pa0["walletid"] = wid
	pa0["signature"] = sig_res
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/wallet", "delete_wallet", pa)
	if err != nil {
		return err.Error()
	}
	if res.Error != nil {
		return res.Error.Message
	}
	return "Create " + strconv.Itoa(mid) + " Success"
}

//func TestWalletControllerDelete(t *testing.T) {
//	InitTest()
//	auid, err := tc.get_authid()
//	if err != nil {
//		fmt.Println("get auid error:", err.Error())
//	}
//	fmt.Println(auid)
//	sid, err := tc.login(auid, acc)
//	if err != nil {
//		fmt.Println("login error:", err.Error())
//	}
//	fmt.Println(sid)
//	mgmtid, err := tc.get_mgmtid(sid, 3)
//	mid := int(mgmtid)
//	fmt.Println(tc.DeleteWalletCfg(acc, sid, mid, 3))
//}
