package lib

import (
	"encoding/json"
	"log"
	"net/http"
)

func Get(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func SetHeader(w http.ResponseWriter, key, val string) {
	w.Header().Set(key, val)
}

func ReturnMsg(w http.ResponseWriter, res string, obj interface{}) {
	SetHeader(w, "Content-Type", "application/json")

	msg := map[string]interface{}{
		"status": http.StatusOK,
		"error":  res,
		"data":   obj,
	}
	if res != "" {
		log.Println(res)
		msg["status"] = http.StatusBadRequest
	}

	str, err := json.Marshal(msg)
	if err != nil {
		log.Println("json.Marshal 错误", err)
	}
	_, _ = w.Write(str)
}
