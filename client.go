package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var username string

// 用户登录
func Login(tcpConn *net.TCPConn) {
	// 提示输入名字
	fmt.Println("请输入你的名字:")
	// 获取输入名字
	fmt.Scanln(&username)
	// 提示登录成功
	fmt.Println("恭喜" + username + "登录成功")
	// 输入给服务器
	tcpConn.Write([]byte("login:" + username))
}

// 显示用户输入
func Read(tcpConn *net.TCPConn) {
	buffer := make([]byte, 128)
	for {
		total, err := tcpConn.Read(buffer)
		if err != nil {
			fmt.Println("与服务器链接断开")
			break
		}
		fmt.Println(string(buffer[:total]))
	}
}

func main() {
	// 监听服务器socket
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	tcpConn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer tcpConn.Close()

	// 显示用户输入(一定要放在下边的for前边，避免不执行)
	go Read(tcpConn)

	// 登录
	Login(tcpConn)

	// 一直监听用户输入
	for {
		var msg string
		// 每行作为输入
		reader := bufio.NewReader(os.Stdin)
		msg, _ = reader.ReadString('\n')
		msg = strings.TrimRight(msg, "\r\n")
		// 输入给服务器
		tcpConn.Write([]byte(username + ":" + msg))
	}
}
