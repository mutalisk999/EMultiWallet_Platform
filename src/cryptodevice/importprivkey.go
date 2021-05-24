package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/mutalisk999/go-cryptocard/src/cryptocard"
	"github.com/mutalisk999/go-lib/src/net/buffer_tcp"
	"io"
	"os"
	"strconv"
	"strings"
)

func importUsage() {
	fmt.Println("./main device_ip device_port privkey_file")
	fmt.Println("example: ")
	fmt.Println("./main 192.168.1.188 1818 /home/test/privkey.bak.txt")
}

func main() {
	if len(os.Args) != 6 {
		importUsage()
	}
	//deviceIp := "192.168.1.188"
	//devicePort := 1818
	//pkFile := "D:/EMultiWallet_Platform/src/cryptodevice/privkey.txt"

	deviceIp := os.Args[1]
	devicePort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("args[2] error:" + err.Error())
		return
	}
	pkFile := os.Args[3]

	conn := new(buffer_tcp.BufferTcpConn)
	err = conn.TCPConnect(deviceIp, uint16(devicePort), 1)
	if err != nil {
		fmt.Println("connect error:" + err.Error())
		return
	}

	file, err := os.OpenFile(pkFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("open pkFile error:" + err.Error())
		return
	}

	bufReader := bufio.NewReader(file)
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("read pkFile error:" + err.Error())
			}
		}
		line = strings.TrimSpace(line)
		data := strings.Split(line, ":")

		var l8req cryptocard.L8Request
		keyIndex, err := strconv.Atoi(data[0])
		if err != nil {
			fmt.Println("keyIndex error:" + data[0])
		}
		privKey, err := hex.DecodeString(data[1])
		if err != nil {
			fmt.Println("privKey error:" + data[1])
		}
		l8req.Set(uint16(keyIndex), privKey)
		err = l8req.Pack(conn)
		if err != nil {
			fmt.Println("send error:" + err.Error())
		}

		var l9resp cryptocard.L9Response
		err = l9resp.UnPack(conn)
		if err != nil {
			fmt.Println("recv error:" + err.Error())
		}

		fmt.Println("import privkey at keyindex:", data[0])
	}

	file.Close()
}
