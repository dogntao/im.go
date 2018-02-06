package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

// 定义用户信息
type User struct {
	// 用户名,用于显示谁发送和私聊
	UserName string
	// websocket连接
	WbConn *websocket.Conn
}

// 用户列表
var UserList []*User

// 定义消息信息
type MsgInfo struct {
	// 发送者
	FromUser string
	// 接收者
	ToUser string
	// 内容类型
	MsgType int
	// 内容
	Msg string
}

// 消息通道
var MsgChannel chan *MsgInfo

// 定义websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 页面模板
var homeTempl = template.Must(template.ParseFiles("home.html"))

// 根路由(解析到html页面)
func serverHome(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/favicon.ico" {
		homeTempl.Execute(w, r.Host)
	}
}

// 输出msg到客户端
func writeMsg() {
	for {
		select {
		case msg := <-MsgChannel:
			for _, v := range UserList {
				// 广播
				if msg.ToUser == "all" {
					v.WbConn.WriteMessage(msg.MsgType, []byte(msg.FromUser+":"+msg.Msg))
				} else {
					// 私聊
					if v.UserName == msg.FromUser || v.UserName == msg.ToUser {
						v.WbConn.WriteMessage(msg.MsgType, []byte(msg.FromUser+":"+msg.Msg))
					}
				}
			}
			fmt.Println(msg)
		}
	}
}

// ws路由(处理websocket请求)
func serverWs(w http.ResponseWriter, r *http.Request) {
	// 获取客户端链接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Websocket server error:", err)
		return
	}
	user := &User{WbConn: conn}
	UserList = append(UserList, user)
	msgInfo := &MsgInfo{FromUser: "system", ToUser: "all"}

	reg := regexp.MustCompile(`^@.*? `)
	var msg string
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("message error:", err)
			return
		}
		msgInfo.MsgType = messageType
		msg = string(p)

		// 登录过，增加前缀
		if user.UserName != "" {
			msg = user.UserName + ":" + msg
			msgInfo.FromUser = user.UserName
		}

		// 处理用户输入
		dataInfo := strings.Split(msg, ":")
		if dataInfo[0] == "login" {
			// 处理用户登录(以login开头)
			user.UserName = dataInfo[1]
			msgInfo.Msg = "欢迎" + user.UserName + "加入"
		} else {
			// 私聊
			toUser := reg.FindString(dataInfo[1])
			msgInfo.Msg = strings.Replace(dataInfo[1], toUser, "", 1)
			if toUser != "" {
				toUser = strings.Replace(toUser, "@", "", 1)
				toUser = strings.Replace(toUser, " ", "", 1)
				msgInfo.ToUser = toUser
			}
		}
		// 放入消息通道
		MsgChannel <- msgInfo
	}
}

func main() {
	MsgChannel = make(chan *MsgInfo, 1)
	go writeMsg()

	// 创建根路由
	http.HandleFunc("/", serverHome)
	http.HandleFunc("/ws", serverWs)
	// 创建http服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
		return
	}
}
