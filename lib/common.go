package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

// GetRoot 获取项目根路径
func GetRoot() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(fileName)
	}
	i := strings.LastIndex(abPath, "/")
	if i < 0 {
		return ""
	}
	return abPath[:i+1]
}

func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
