package main

import (
	"Chess/chat"
	chess "Chess/chess/handler"
	"Chess/redispool"
	"log"
	"net/http"
)

func main() {
	redispool.Start()
	http.HandleFunc("/joinRoom", chat.Room)
	http.HandleFunc("/haveRoom", chat.HaveRoom)
	http.HandleFunc("/roomNum", chat.RoomNum)

	http.Handle("/chess", &chess.MachineMsg{})
	err := http.ListenAndServe("0.0.0.0:9000", nil)
	if err != nil {
		log.Println(err)
	}
}
