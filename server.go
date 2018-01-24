package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// 连接map，用于发送广播
var connMap map[string]*net.TCPConn

// 用户和连接关联关系，用于私聊
var userMap map[string]*net.TCPConn

// 连接和用户关联关系，用于连接断开提示
var conUserMap map[*net.TCPConn]string

func say(tcpConn *net.TCPConn) {
	// 一直监听用户输入
	for {
		// 获取用户消息
		data := make([]byte, 128)
		total, err := tcpConn.Read(data)
		if err != nil {
			fmt.Println("用户:" + conUserMap[tcpConn] + "退出")
			break
		}

		// 处理用户输入
		dataInfo := strings.Split(string(data[:total]), ":")
		toUser := ""
		fmt.Println(string(data[:total]))
		// fmt.Println(dataInfo)

		if dataInfo[0] == "login" {
			// 处理用户登录(以login开头)
			dataInfo[1] = strings.TrimRight(dataInfo[1], "\r\n")
			userMap[dataInfo[1]] = tcpConn
			conUserMap[tcpConn] = dataInfo[1]
			data = []byte("欢迎" + dataInfo[1] + "加入")
		} else {
			// 处理私聊(以@开头例:@dt hello world)
			reg := regexp.MustCompile(`^@.*? `)
			toUser = reg.FindString(dataInfo[1])
			data = []byte(dataInfo[0] + ":" + strings.Replace(dataInfo[1], toUser, "", 1))
			toUser = strings.Replace(toUser, "@", "", 1)
			toUser = strings.Replace(toUser, " ", "", 1)
		}
		// fmt.Println(toUser)
		// fmt.Println(string(data))
		if toUser == "" {
			// 发送广播
			for _, conn := range connMap {
				// 自己不给自己发送信息
				if conn == tcpConn {
					continue
				}
				conn.Write(data)
			}
		} else {
			// 私聊
			if _, ok := userMap[toUser]; ok {
				userMap[toUser].Write(data)
			}
		}
	}
}

func main() {
	// 创建socket监听
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	connMap = make(map[string]*net.TCPConn)
	userMap = make(map[string]*net.TCPConn)
	conUserMap = make(map[*net.TCPConn]string)
	// 一直监听用户链接
	for {
		// 用户连接
		tcpConn, _ := tcpListener.AcceptTCP()
		defer tcpConn.Close()
		// 用户地址存入map(用于广播)
		connMap[tcpConn.RemoteAddr().String()] = tcpConn
		// 处理用户输入
		go say(tcpConn)
	}
}
