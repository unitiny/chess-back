package handler

import (
	"Chess/chess/config"
	"Chess/lib"
	"Chess/redispool"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

type MachineActionFunc func(msg *MachineMsg) string

// MachineMsg 前端消息格式
type MachineMsg struct {
	Action Type                   `json:"action"`
	ID     int                    `json:"id"`
	Room   string                 `json:"room"`
	Start  map[string]interface{} `json:"start"`
	End    map[string]interface{} `json:"end"`
}

var MachineActionFuncs map[Type]MachineActionFunc

func init() {
	MachineActionFuncs = make(map[Type]MachineActionFunc)
	MachineActionFuncs[CLEAR] = DelRoom
	//MachineActionFuncs[BACK] = BackRoomStep // 取消悔棋接口
	MachineActionFuncs[GO] = GetChessStep
}

// GetRoomSteps 根据房间获取该棋局历史记录
func GetRoomSteps(room string) []Step {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := config.RoomPrefix + room

	reply, err := redis.Values(conn.Do("lrange", key, 0, -1))
	if err != nil {
		log.Println("GetRoomSteps err", err)
	}

	steps := make([]Step, len(reply))
	for i := 0; i < len(reply); i++ {
		step := new(Step)
		_ = json.Unmarshal(reply[i].([]byte), step)
		steps[i] = *step
	}
	return steps
}

// RecordRoomStep 向房间写入下棋记录
func RecordRoomStep(room string, steps ...Step) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := config.RoomPrefix + room
	//_, _ = redis.Int(conn.Do("llen", key))

	for i := 0; i < len(steps); i++ {
		if i != 0 {
			steps[i].Id = steps[i-1].Id + 1
		}

		str, _ := json.Marshal(steps[i])
		reply, err := conn.Do("rpush", key, str)
		if err != nil {
			log.Println("RecordRoomStep :", reply, err)
			return
		}
	}
}

// RecordRoomStep1 向房间写入下棋记录
func RecordRoomStep1(room string, steps []Step) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := config.RoomPrefix + room
	length, _ := redis.Int(conn.Do("llen", key))
	if length > len(steps) {
		// 删除记录，重新写入
		_, _ = conn.Do("del", key)
	} else {
		// 把新记录加入即可
		steps = steps[length:]
	}

	for i := 0; i < len(steps); i++ {
		if i != 0 {
			steps[i].Id = steps[i-1].Id + 1
		}

		str, _ := json.Marshal(steps[i])
		reply, err := conn.Do("rpush", key, str)
		if err != nil {
			log.Println("RecordRoomStep1 :", reply, err)
			return
		}
	}
}

// DelRoom 删除房间记录
func DelRoom(msg *MachineMsg) string {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := config.RoomPrefix + msg.Room
	length, _ := redis.Int(conn.Do("llen", key))
	res := "DelRoom success"
	defer func() {
		log.Println(res)
	}()

	if length <= 0 {
		res = "DelRoom room not exist"
		return res
	}

	_, err := conn.Do("del", key)
	if err != nil {
		res = fmt.Sprintf("DelRoom : %s", err)
		return res
	}
	return res
}

// BackRoomStep 回退房间棋步
func BackRoomStep(msg *MachineMsg) string {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := config.RoomPrefix + msg.Room
	result, _ := conn.Do("rpop", key)

	res := "BackRoomStep success"
	if result == nil {
		res = "BackRoomStep pop fail"
	}
	return res
}

func (m *MachineMsg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		msg := new(MachineMsg)
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(msg)
		fmt.Printf("%+v\n", msg)
		if err != nil {
			lib.ReturnMsg(w, err.Error(), nil)
			return
		}

		res := MachineActionFuncs[msg.Action](msg)
		lib.ReturnMsg(w, "", res)
	default:
		lib.ReturnMsg(w, "only support post method", nil)
	}
}
