package main

import (
	"encoding/hex"
	"fmt"
	"github.com/mutalisk999/go-cryptocard/src/cryptocard"
	"github.com/mutalisk999/go-lib/src/net/buffer_tcp"
	"os"
	"strconv"
)

func exportUsage() {
	fmt.Println("./main device_ip device_port start_index end_index privkey_file")
	fmt.Println("example: ")
	fmt.Println("./main 192.168.1.188 1818 1 1024 /home/test/privkey.bak.txt")
}

func main() {
	if len(os.Args) != 6 {
		exportUsage()
	}
	//deviceIp := "192.168.1.188"
	//devicePort := 1818
	//startIndex := 1
	//endIndex := 1024
	//pkBakFile := "D:/EMultiWallet_Platform/src/cryptodevice/privkey.bak.txt"

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
	pkBakFile := os.Args[5]

	conn := new(buffer_tcp.BufferTcpConn)
	err = conn.TCPConnect(deviceIp, uint16(devicePort), 1)
	if err != nil {
		fmt.Println("connect error:" + err.Error())
		return
	}

	file, err := os.OpenFile(pkBakFile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("open pkBakFile error:" + err.Error())
		return
	}

	for i := startIndex; i <= endIndex; i++ {
		var l5req cryptocard.L5Request
		l5req.Set(uint16(i), 0)
		err = l5req.Pack(conn)
		if err != nil {
			fmt.Println("send error:" + err.Error())
		}

		var l6resp cryptocard.L6Response
		err = l6resp.UnPack(conn)
		if err != nil {
			fmt.Println("recv error:" + err.Error())
		}

		file.WriteString(fmt.Sprintf("%d:%s\n", i, hex.EncodeToString(l6resp.PrivKey)))
		fmt.Println("export privkey at keyindex:", i)
	}

	file.Close()
}
