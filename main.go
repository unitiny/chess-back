package main

import (
	"Chess/chat"
	chess "Chess/chess/handler"
	_ "Chess/redispool"
	"net/http"
)

func main() {
	http.HandleFunc("/joinRoom", chat.Room)
	http.HandleFunc("/haveRoom", chat.HaveRoom)
	http.HandleFunc("/roomNum", chat.RoomNum)

	http.Handle("/chess", &chess.MachineMsg{})
	_ = http.ListenAndServe("0.0.0.0:9000", nil)
}
