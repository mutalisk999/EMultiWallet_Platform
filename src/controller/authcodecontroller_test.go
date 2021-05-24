package controller

import (
	"testing"
	"fmt"
)

func TestGetAuthCodeController1(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param GetAuthCodeParam
	param.Height = 80
	param.Width = 240
	param.Len = 6
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/authcode", "get_authcode", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}
	authcodeid, _ := res.Result.(map[string]interface{})["authcodeid"]
	authcodestream, _ := res.Result.(map[string]interface{})["authcodestream"]

	fmt.Println("authcodeid:", authcodeid)
	fmt.Println("authcodestream:", authcodestream)
}

func TestGetAuthCodeController2(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param GetAuthCodeParam
	param.Height = 80
	param.Width = 240
	param.Len = 6
	params = append(params, param)
	res, err := tc2.doHttpJsonRpcCallType1("/apis/authcode", "get_authcode", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}
	authcodeid, _ := res.Result.(map[string]interface{})["authcodeid"]
	authcodestream, _ := res.Result.(map[string]interface{})["authcodestream"]

	fmt.Println("authcodeid:", authcodeid)
	fmt.Println("authcodestream:", authcodestream)
}

func TestGetAuthCodeController3(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param GetAuthCodeParam
	param.Height = 80
	param.Width = 240
	param.Len = 6
	params = append(params, param)
	res, err := tc3.doHttpJsonRpcCallType1("/apis/authcode", "get_authcode", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}
	authcodeid, _ := res.Result.(map[string]interface{})["authcodeid"]
	authcodestream, _ := res.Result.(map[string]interface{})["authcodestream"]

	fmt.Println("authcodeid:", authcodeid)
	fmt.Println("authcodestream:", authcodestream)
}

