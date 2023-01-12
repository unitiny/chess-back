package handler

import (
	"Chess/chess/config"
	"Chess/redispool"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

type MachineActionFunc func(msg *MachineMsg) string

// MachineMsg 前端消息格式
type MachineMsg struct {
	Action Type                   `json:"action"`
	Room   string                 `json:"room"`
	Start  map[string]interface{} `json:"start"`
	End    map[string]interface{} `json:"end"`
}

var MachineActionFuncs map[Type]MachineActionFunc

func init() {
	MachineActionFuncs = make(map[Type]MachineActionFunc)
	MachineActionFuncs[CLEAR] = DelRoom
	MachineActionFuncs[BACK] = BackRoomStep
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
	_, _ = redis.Int(conn.Do("llen", key))

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
