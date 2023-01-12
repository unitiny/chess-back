package chat

import (
	chess "Chess/chess/handler"
	"Chess/redispool"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	Name    string
	Room    string
	PubName string
	Conn    *websocket.Conn
	Close   chan bool
}

type Message struct {
	Conn      *websocket.Conn `json:"-"`
	EventType int             `json:"type"`    // 0 发送消息 1 发送消息给某个用户 2 用户加入 3 用户退出
	Content   string          `json:"content"` // 消息内容
	Name      string          `json:"name"`    // 用户名称
}

type JSONResult struct {
	Status int
	Error  string
	Data   interface{}
}

func Get(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func SetHeader(w http.ResponseWriter, key, val string) {
	w.Header().Set(key, val)
}

func ReturnMsg(w http.ResponseWriter, res string, obj interface{}) {
	SetHeader(w, "Content-Type", "application/json")

	msg := &JSONResult{Status: 200, Error: res, Data: obj}
	if res != "" {
		log.Println(res)
		msg.Status = 404
		msg.Error = res
	}

	str, err := json.Marshal(msg)
	if err != nil {
		log.Println("json.Marshal 错误", err)
	}
	_, _ = w.Write(str)
}

func NewClient(conn *websocket.Conn, name, room, pub string) *Client {
	return &Client{Name: name, Conn: conn, Room: room, PubName: pub, Close: make(chan bool, 1)}
}

func NewMsg(eventType int, content, name string, conn *websocket.Conn) *Message {
	return &Message{EventType: eventType, Content: content, Name: name, Conn: conn}
}

// LeaveRoom 没什么意义，每次调用都会该redisConn取消订阅，本来就没订阅故没有效果
func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()

	room := Get(r, "room")
	if room == "" {
		ReturnMsg(w, "room为空", nil)
		return
	}

	reply, _ := redis.Int(conn.Do("get", room))
	_, _ = conn.Do("set", reply-1)
	ReturnMsg(w, "", fmt.Sprintf("当前房间剩余人数 %d", reply-1))
}

func HaveRoom(w http.ResponseWriter, r *http.Request) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()

	room := Get(r, "room")
	reply, err := redis.Int(conn.Do("get", room))
	if err != nil {
		fmt.Println(err)
	}

	ReturnMsg(w, "", reply)
}

func RoomNum(w http.ResponseWriter, r *http.Request) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()

	room := Get(r, "room")
	reply, _ := redis.String(conn.Do("get", room))
	if reply != "" {
		ReturnMsg(w, "", reply)
	}
	ReturnMsg(w, "房间不存在", 0)
}

func Room(w http.ResponseWriter, r *http.Request) {
	name := Get(r, "name")
	room := Get(r, "room")

	// 0 创建房间  1 加入房间  2 机器房间
	action := Get(r, "action")

	// 参数判断
	if checkParams(name, room, action) != nil {
		ReturnMsg(w, "所有参数不得为空", nil)
		return
	}

	// 连接判断
	conn, err := connectWS(w, r)
	if err != nil {
		ReturnMsg(w, "连接错误", nil)
		return
	}

	// 创建房间及订阅频道
	eventType, _ := strconv.Atoi(action)
	redisConn := redispool.RedisPool.Get()
	redispool.GetsSetRoom(redisConn, room)

	pub := "chat-" + room
	if eventType == 2 {
		pub = "machine-" + room
	}
	client := NewClient(conn, name, room, pub)
	newPub := client.JoinRoom() // 加入房间频道

	// 消息往来
	go client.broadCast(newPub)
	go client.readMessage(eventType)

	// 监听停止
	for {
		select {
		case <-client.Close:
			close(client.Close)

			// 关闭房间
			_ = closeRoom(redisConn, room)

			// 退订频道
			log.Println(room + "退订频道")
			_ = newPub.Unsubscribe(room)
			_ = newPub.Conn.Close()
			_ = newPub.Close()
			return
		}
	}
}

// 参数校验
func checkParams(params ...string) error {
	for i := 0; i < len(params); i++ {
		if params[i] == "" {
			return errors.New("参数为空")
		}
	}
	return nil
}

// 连接ws
func connectWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := ws.Upgrade(w, r, nil)
	return conn, err
}

// 关闭房间
func closeRoom(conn redis.Conn, room string) error {
	defer conn.Close()
	// 还有人存在则不关闭
	reply, _ := redis.Int(conn.Do("get", room))
	reply--
	if reply <= 0 {
		_, _ = conn.Do("del", room)
	} else {
		_, _ = conn.Do("set", room, reply)
	}
	return nil
}

// LeaveRoom 退出频道
func (c *Client) LeaveRoom() {

}

// JoinRoom 订阅频道
func (c *Client) JoinRoom() redis.PubSubConn {
	conn := redispool.RedisPool.Get()
	newPub := redis.PubSubConn{
		Conn: conn,
	}

	//订阅一个频道
	err := newPub.Subscribe(c.PubName)
	if err != nil {
		fmt.Println("Subscribe err:", err)
		return newPub
	}
	log.Println(c.Name, "join room", c.Room)
	return newPub
}

// 监听订阅消息
func (c *Client) broadCast(Pub redis.PubSubConn) {
	for {
		//Receive()返回的是空接口interface{}的类型,所以需要断言
		switch v := Pub.Receive().(type) {
		//Redis.Message结构体
		//type Message struct {
		//	Channel string
		//	Pattern string
		//	Data    []byte
		//}

		case redis.Message:
			_ = c.Conn.WriteMessage(websocket.TextMessage, v.Data)
			fmt.Printf("channel:%s,\tdata:%s\n", v.Channel, string(v.Data))

		//订阅或者取消订阅
		case redis.Subscription:
			fmt.Printf("Channel:%s\tCount:%d \tKind:%s\n", v.Channel, v.Count, v.Kind)
		}

	}
}

// 读取消息,并作为发布者写入信息
func (c *Client) readMessage(eventType int) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("readMessage:", err)
			c.Close <- true
			break
		}

		// 如果要人机下
		var message = new(chess.MachineMsg)
		err = json.Unmarshal(data, message)
		if err == nil && message.Room != "" {
			log.Println("message ", message)

			data = []byte(chess.MachineActionFuncs[message.Action](message))
			log.Println("result ", string(data))
		}

		msg := NewMsg(eventType, string(data), c.Name, c.Conn)
		c.writeMessage(msg)
	}
}

// 写入消息
func (c *Client) writeMessage(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("msg error:", err)
	}

	sender := redispool.RedisPool.Get()
	defer sender.Close()

	_, err = sender.Do("publish", c.PubName, data)
	if err != nil {
		log.Println("send msg error:", err)
	}
}
