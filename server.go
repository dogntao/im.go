package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

// 定义用户信息
type User struct {
	// 用户名,用于显示谁发送和私聊
	UserName string
	// 接收到的消息
	ReceiveMsg chan []byte
	// 接收到的消息
	ReceiveMsgType chan []byte
	// websocket连接
	WbConn *websocket.Conn
}

// 用户列表
var UserList []User

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

// 处理用户输入(每个用户一个goroutine)
func say(user *User) {
	// 获取消息
	// for {
	// 	messageType, p, err := user.WbConn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("message error:", err)
	// 		return
	// 	}
	// 	fmt.Println(conn)
	// 	fmt.Println(messageType)
	// 	fmt.Println(string(p))
	// }
}

// ws路由(处理websocket请求)
func serverWs(w http.ResponseWriter, r *http.Request) {
	// 获取客户端链接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Websocket server error:", err)
		return
	}
	var msg string
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("message error:", err)
			return
		}
		// 登录
		msg = string(p)
		conn.WriteMessage(messageType, p)
	}
}

func main() {
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
