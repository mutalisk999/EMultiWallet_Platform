package controller

import (
	"testing"
	"fmt"
)

func TestGetOpLogsController(t *testing.T) {
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

	var param GetOpLogsParam
	param.SessionId = sid
	param.AcctId = []int{}
	param.OpType = []int{}
	param.OpTime = [2]string{"",""}
	param.OffSet = 0
	param.Limit = 100
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/log", "get_op_logs", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}

	result := res.Result.(map[string]interface{})
	total, _ := result["total"]
	fmt.Println("total:", total)
	logs, _ := result["logs"].([]interface{})
	for _, log := range logs {
		res := log.(map[string]interface{})
		logid, _ := res["logid"]
		acctid, _ := res["acctid"]
		optype, _ := res["optype"]
		optime, _ := res["optime"]
		content, _ := res["content"]
		fmt.Println("logid:", logid)
		fmt.Println("acctid:", acctid)
		fmt.Println("optype:", optype)
		fmt.Println("optime:", optime)
		fmt.Println("content:", content)
	}
}
