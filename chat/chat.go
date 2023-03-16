package chat

import (
	chess "Chess/chess/handler"
	"Chess/lib"
	"Chess/redispool"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	Name  string
	Room  string
	Conn  *websocket.Conn
	Pub   redis.PubSubConn
	Close chan bool
}

type Message struct {
	Conn    *websocket.Conn `json:"-"`
	Content string          `json:"content"` // 消息内容
	Name    string          `json:"name"`    // 用户名称
}

type JSONResult struct {
	Status int
	Error  string
	Data   interface{}
}

func NewClient(conn *websocket.Conn, name, room string) *Client {
	return &Client{Name: name, Conn: conn, Room: room, Close: make(chan bool, 1)}
}

func NewMsg(content, name string, conn *websocket.Conn) *Message {
	return &Message{Content: content, Name: name, Conn: conn}
}

func HaveRoom(w http.ResponseWriter, r *http.Request) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()

	room := lib.Get(r, "room")
	reply, err := redis.Int(conn.Do("get", room))
	if err != nil {
		fmt.Println(err)
	}

	lib.ReturnMsg(w, "", reply)
}

func RoomNum(w http.ResponseWriter, r *http.Request) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()

	room := lib.Get(r, "room")
	reply, _ := redis.String(conn.Do("get", room))
	if reply != "" {
		lib.ReturnMsg(w, "", reply)
	}
	lib.ReturnMsg(w, "房间不存在", 0)
}

// Room 加入房间，建立连接
func Room(w http.ResponseWriter, r *http.Request) {
	name := lib.Get(r, "name")
	room := lib.Get(r, "room")

	// 参数判断
	if checkParams(name, room) != nil {
		lib.ReturnMsg(w, "所有参数不得为空", nil)
		return
	}

	// 连接判断
	conn, err := connectWS(w, r)
	if err != nil {
		lib.ReturnMsg(w, "连接错误", nil)
		return
	}

	// 创建房间及订阅频道
	redisConn := redispool.RedisPool.Get()
	redispool.GetSetRoom(redisConn, room)

	client := NewClient(conn, name, room)
	client.joinRoom() // 加入房间频道

	// 消息往来
	go client.broadCast()
	go client.readMessage()

	// 加入成功，广播房间
	msg := NewMsg(fmt.Sprintf("用户[%s]加入房间[%s]", client.Name, client.Room),
		client.Name, client.Conn)
	client.writeMessage(msg)

	// 监听停止
	for {
		select {
		case <-client.Close:
			// 关闭房间
			_ = closeRoom(redisConn, room)

			// 广播退出，退订频道
			client.leaveRoom()
			close(client.Close)
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

// leaveRoom 退出频道
func (c *Client) leaveRoom() {
	content := fmt.Sprintf("用户[%s]退出房间[%s]", c.Name, c.Room)
	msg := NewMsg(content, c.Name, c.Conn)
	c.writeMessage(msg) // 广播退出

	_ = c.Pub.Unsubscribe(c.Room)
	_ = c.Pub.Conn.Close()
	_ = c.Pub.Close()
}

// joinRoom 订阅频道
func (c *Client) joinRoom() {
	conn := redispool.RedisPool.Get()
	newPub := redis.PubSubConn{
		Conn: conn,
	}

	//订阅一个频道
	err := newPub.Subscribe(c.Room)
	if err != nil {
		fmt.Println("Subscribe err:", err)
		return
	}
	c.Pub = newPub
}

// 监听订阅消息
func (c *Client) broadCast() {
	for {
		//Receive()返回的是空接口interface{}的类型,所以需要断言
		switch v := c.Pub.Receive().(type) {
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
func (c *Client) readMessage() {
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

		msg := NewMsg(string(data), c.Name, c.Conn)
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

	_, err = sender.Do("publish", c.Room, data)
	if err != nil {
		log.Println("send msg error:", err)
	}
}
