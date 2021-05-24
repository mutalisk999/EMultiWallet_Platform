package main

import (
	"encoding/hex"
	"fmt"
	"github.com/mutalisk999/go-cryptocard/src/cryptocard"
	"github.com/mutalisk999/go-lib/src/net/buffer_tcp"
	"os"
	"time"
	"utils"
	"strconv"
)

func initUsage() {
	fmt.Println("./main device_ip device_port start_index end_index sql_file privkey_file")
	fmt.Println("example: ")
	fmt.Println("./main 192.168.1.188 1818 1 1024 /home/test/device_init.sql /home/test/privkey.txt")
}

func main() {
	if len(os.Args) != 7 {
		initUsage()
	}
	//deviceIp := "192.168.1.188"
	//devicePort := 1818
	//startIndex := 1
	//endIndex := 1024
	//sqlFile := "D:/EMultiWallet_Platform/src/cryptodevice/device_init.sql"
	//pkFile := "D:/EMultiWallet_Platform/src/cryptodevice/privkey.txt"

	deviceIp := os.Args[1]
	devicePort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("args[2] error:" + err.Error())
		return
	}
	startIndex, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("args[3] error:" + err.Error())
		return
	}
	endIndex, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("args[4] error:" + err.Error())
		return
	}
	sqlFile := os.Args[5]
	pkFile := os.Args[6]

	conn := new(buffer_tcp.BufferTcpConn)
	err = conn.TCPConnect(deviceIp, uint16(devicePort), 1)
	if err != nil {
		fmt.Println("connect error:" + err.Error())
		return
	}

	file, err := os.OpenFile(sqlFile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("open sqlFile error:" + err.Error())
		return
	}
	file2, err := os.OpenFile(pkFile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("open pkFile error:" + err.Error())
		return
	}

	for i := startIndex; i <= endIndex; i++ {
		var l1req cryptocard.L1Request
		l1req.Set('4', uint16(i))
		err = l1req.Pack(conn)
		if err != nil {
			fmt.Println("send error:" + err.Error())
		}

		var l2resp cryptocard.L2Response
		err = l2resp.UnPack(conn)
		if err != nil {
			fmt.Println("recv error:" + err.Error())
		}

		file.WriteString(fmt.Sprintf("insert into tbl_pubkey_pool(keyindex, pubkey, isused, createtime) values(%d, \"%s\", false, \"%s\")",
			i, hex.EncodeToString(l2resp.PubKey), utils.TimeToFormatString(time.Now())) + ";\n")

		file2.WriteString(fmt.Sprintf("%d:%s", i, hex.EncodeToString(l2resp.PrivKey)) + "\n")

		fmt.Println("create privkey at keyindex:", i)
	}

	file.Close()
	file2.Close()

	conn.TCPDisConnect()
}
