package controller

import (
	"testing"
	"fmt"
	"reflect"
	"encoding/json"
)

func TestGetIdentityController(t *testing.T) {
	InitTest()
	params := make([]interface{}, 0)
	var param GetIdentityParam
	param.SessionId = ""
	param.IdType = 1
	params = append(params, param)
	res, err := tc.doHttpJsonRpcCallType1("/apis/identity", "get_auto_inc_id", params)
	if err != nil {
		fmt.Println(err.Error())
	}
	if res.Error != nil {
		fmt.Println(res.Error.Message)
	}
	fmt.Println(reflect.TypeOf(res.Result))
	identityid, _ := res.Result.(json.Number).Int64()

	fmt.Println("identityid:", identityid)
}