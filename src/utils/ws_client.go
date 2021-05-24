package utils

import (
	"sync"
	"github.com/kataras/iris/websocket"
	"fmt"
	xwebsocket "golang.org/x/net/websocket"
	"encoding/json"
	"time"
	"github.com/kataras/iris/core/errors"
	"strings"
)

var GlobalReqMap sync.Map
var GlobalWsConn *xwebsocket.Conn
var GlobalWsTimeOut uint
var GlobalIsLogin bool
var JsonId int
var GlobalWsConnEstablish chan bool
var GlobalWsConnReconnect chan bool
var GlobalPushTasks chan TaskPushParamWS


func GetJsonId() (int){
	JsonId++
	return JsonId
}

func WSRequestFuture(index int) *chan interface{}{
	future := make(chan interface{},1)
	GlobalReqMap.Store(index,&future)
	return &future
}

func WSRequest(ws *xwebsocket.Conn, eventName string, message string, timeout uint)(string,error){
	var req JsonRpcRequest
	err := json.Unmarshal([]byte(message),&req)
	if err!=nil{
		return "",err
	}
	req.Id = GetJsonId()
	data,err := json.Marshal(req)
	if err!= nil {
		return "",err
	}
	err = SendMessage(ws,eventName,string(data))
	if err!=nil{
		return "",err
	}

	timer := time.NewTimer(time.Duration(timeout)*time.Second)
	go func() {
		<- timer.C
	}()
	future := WSRequestFuture(req.Id)
	select {
	case resData:=<- *future:
		return resData.(string),nil
	case <- timer.C:
		return "",errors.New("request timeout!")
	}
}

func SendMessage(ws *xwebsocket.Conn, eventName, message string) error {
	buffer := []byte(message)
	return sendBytes(ws,eventName, buffer)
}

func sendBytes(ws *xwebsocket.Conn, eventName string, message []byte) error {
	buffer := []byte(fmt.Sprintf("%s%v;0;", websocket.DefaultEvtMessageKey, eventName))
	buffer = append(buffer, message...)
	lenLeft := len(buffer)
	for {
		n, err := ws.Write(buffer)
		if err != nil {
			return err
		}
		lenLeft -= n
		buffer = buffer[n:]
		if lenLeft <= 0 {
			break
		}
	}
	return nil
}

func ConnectWebSocket(url string, origin string) (*xwebsocket.Conn, error) {
	ws, err := xwebsocket.Dial(url, "", origin)
	return ws, err
}

func CloseWebSocket(ws *xwebsocket.Conn) error {
	if ws != nil {
		return ws.Close()
	}
	return nil
}

func WSReadAndDispose() error {
	// read from web socket
	fmt.Println("WSReadAndDispose...")

	for {
		var buffer = make([]byte, 0)
		var bufferLen = 0
		for {
			var readBuf = make([]byte, 10*1024*1024)
			n, err := GlobalWsConn.Read(readBuf)
			if err != nil {
				CloseWebSocket(GlobalWsConn)
				return err
			}
			bufferLen += n
			buffer = append(buffer, readBuf[0:n]...)

			// 完整的json返回数据
			if strings.Count(string(buffer), "{") == strings.Count(string(buffer), "}") {
				break
			}

			if bufferLen > 4*1024*1024 {
				CloseWebSocket(GlobalWsConn)
				return errors.New("size of bytes read from websocket is too big")
			}
		}

		fmt.Println("WSReadAndDispose() read from server:", string(buffer[:]))
		var reqPush TaskPushRequestWS
		err := json.Unmarshal(buffer[:], &reqPush)
		if err != nil {
			CloseWebSocket(GlobalWsConn)
			return err
		} else if reqPush.Method == "" {
			var res JsonRpcResponse
			err = json.Unmarshal(buffer[:], &res)
			if err != nil {
				CloseWebSocket(GlobalWsConn)
				return err
			}
			// deal response message
			chanObj, isExist := GlobalReqMap.Load(res.Id)
			if isExist {
				*chanObj.(*chan interface{}) <- string(buffer[:])
			}
		} else {
			// deal push message
			for _, task := range reqPush.Params {
				GlobalPushTasks <- task
			}

			// send task-push response
			//var reqPushRes TaskPushResponseWS
			//reqPushRes.Id = reqPush.Id
			//reqPushRes.Error = nil
			//reqPushRes.Result = true
			//bytesSend , err := json.Marshal(reqPushRes)
			//if err != nil {
			//	CloseWebSocket(ws)
			//	return err
			//}
			//
			//_, err = ws.Write(bytesSend)
			//if err != nil {
			//	CloseWebSocket(ws)
			//	return err
			//}
		}
	}
}