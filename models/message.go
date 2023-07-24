package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 调度
func init() {
	fmt.Println("init...")
	go udpSendProc()
	go udpRecvProc()
}

type Message struct {
	gorm.Model
	FromId   int64  // 发送者
	TargetId int64  // 接收者
	Type     int    // 发送类型: 群聊,私聊,广播等
	Media    int    // 文字,图片,音频
	Content  string // 消息内容
	Pic      string // 图片
	Url      string // 链接
	Desc     string // 描述
	Amount   int    // 其它统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// Node 初始化

// 映射关系
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {

	// 从前端获取的请求参数都是字符串类型的

	// 1.获取参数 以及 检验 token 以及其它合法性
	// token := query.Get("token") //暂时不校验
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")

	isVALID := true // 待完成
	conn, err := (&websocket.Upgrader{
		// token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isVALID
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2.获取 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),   // 有可能有多个人给一个人发送消息,管道容量暂设定为 50
		GroupSets: set.New(set.ThreadSafe), // 线程安全群集合
	}
	// 3.获取关系
	// 4.userId 与 node绑定,并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5.多协程发送逻辑:一个人有可能同时给多个人发送多个大文件
	go sendProc(node)
	// 6.多协程完成接收逻辑:一个人有可能同时接收多个人发送的大文件
	go recvProc(node)
	// 仅供测试
	sendMsg(userId, []byte("欢迎来到聊天室233!"))

}
func sendProc(node *Node) {
	fmt.Println("sendProc...")
	// 一直循环等待处理 Node_用户 所发的消息
	for {
		select {
		case data := <-node.DataQueue:
			// 将消息写进 Node_用户 的conn
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}
func recvProc(node *Node) {
	fmt.Println("recvProc...")
	// 一直循环等待处理 Node_用户 所接受的消息
	for {
		// 将 Node_用户 的conn 中的消息读出来
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 将读到的数据传送进 udpSendChan
		broadMsg(data)
		// 后端显示
		fmt.Println("[ws] <<<<< ", data)
	}

}

// 所有用户的数据来后,都将存储到 udpSendChan
var udpSendChan = make(chan []byte, 1024*32)

func broadMsg(data []byte) {
	udpSendChan <- data

}

// 完成数据发送协程,发送到 UDPconn 中
func udpSendProc() {
	fmt.Println("udpSendProc...")
	UDPconn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 30, 0, 159),
		Port: 3000,
	})
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(UDPconn)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpSendChan:
			_, err := UDPconn.Write(data)
			if err != nil {
				fmt.Println(err)
				return

			}
		}
	}

}

// 完成数据接收协程
func udpRecvProc() {
	fmt.Println("udpRecvProc...")
	UDPconn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(UDPconn)
	if err != nil {
		fmt.Println(err)
	}
	for {
		var buf [512]byte
		n, err := UDPconn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		disPatch(buf[0:n])

	}

}

// 后端调度逻辑,这一块儿还需要完善
func disPatch(data []byte) {
	fmt.Println("disPatch1...")
	// 初始化 message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(data))
		return
	}
	switch msg.Type {
	case 1:
		fmt.Println(msg)
		sendMsg(msg.TargetId, data)
		//case 2:
		//	sendGroupMsg()
		//case 3:
		//	sendAllMsg()
		//case 4:
		//	sendMsg()

	}
	fmt.Println("disPatch2...")
}
func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg...")
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}

}
