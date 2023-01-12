package main

import (
	"Chess/chat"
	"net/http"
)

func main() {
	http.HandleFunc("/", chat.Room)
	http.HandleFunc("/leaveRoom", chat.LeaveRoom)
	http.HandleFunc("/haveRoom", chat.HaveRoom)
	http.HandleFunc("/roomNum", chat.RoomNum)
	_ = http.ListenAndServe("0.0.0.0:9000", nil)
}
