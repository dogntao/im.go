package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var homeTempl = template.Must(template.ParseFiles("home.html"))

// 根路由(解析到html页面)
func serverHome(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/favicon.ico" {
		homeTempl.Execute(w, r.Host)
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

	// 获取消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("message error:", err)
			return
		}
		fmt.Println(conn)
		fmt.Println(messageType)
		fmt.Println(string(p))
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
