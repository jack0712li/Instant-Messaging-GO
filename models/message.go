package models

import (
	// "encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   int64 //sender
	TargetId int64 //receiver
	Type     int //message sending type group single or broadcast
	Media    int    //message media type text image voice
	Content  string //message content
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}


type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}


var clientMap map[int64]*Node = make(map[int64]*Node, 0)


// lock
var rwLocker sync.RWMutex


func Chat(writer http.ResponseWriter, request *http.Request) {
	//1. 获取参数和校验token etc
	// token := query.Get("token")

	query := request.URL.Query()
	Id := query.Get("userId")
	userId, err := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalid := true //check token
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2. get connection
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//3. user contact

	//4. userId combine node and lock
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5. send
	go sendProc(node)

	//6. recv
	go recvProc(node)
	sendMsg(userId, []byte("Welcome to ginchat"))



}


func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}


func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(172, 19, 80, 1),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := Message{}
	fmt.Println(data)
	fmt.Println(msg)
	// err := json.Unmarshal(data, &msg)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		// case 2://群发
		// 	sendGroupMsg()
		// case 3://广播
		// 	sendAllMsg()
		// case 4:

	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}