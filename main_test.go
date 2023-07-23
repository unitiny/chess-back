package main

import (
	"Chess/chess/config"
	chess "Chess/chess/handler"
	"path"
	"runtime"
	"testing"
)

func BenchmarkGetChessStepStep(b *testing.B) {
	// 0.914579 -> 0.989741 ns/op
	for i := 0; i < b.N; i++ {
		e := chess.GetEngine()
		e.GetNextStep(config.MAX_DEPTH, -chess.MAX_SCORE, chess.MAX_SCORE)
	}
}

// GetRoot 获取项目根路径
func GetRoot() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
