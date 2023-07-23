package main

import (
	"Chess/chat"
	chess "Chess/chess/handler"
	"Chess/redispool"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func EditHandler(root string, prefix string) {
	dir, err := os.Open(root)
	defer dir.Close()
	if err != nil {
		fmt.Println(err)
	}

	dirEntries, err := dir.ReadDir(-1)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			dirPath := getPath(root, entry.Name())
			fmt.Printf("%s %s  %s\n", prefix, entry.Name(), dirPath)
			EditHandler(dirPath, prefix+"  ")
		} else {
			if len(entry.Name()) < 20 {
				fmt.Printf("%s %s\n", prefix, entry.Name())
			}
		}
	}
}

func getPath(prefix string, name ...string) string {
	res := strings.Builder{}
	res.WriteString(prefix)
	for _, s := range name {
		res.WriteString("/")
		res.WriteString(s)
	}
	return res.String()
}

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
