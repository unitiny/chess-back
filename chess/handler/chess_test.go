package handler

import (
	"Chess/chess/config"
	"Chess/chess/lib"
	"Chess/redispool"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func isEqualEngine(e1, e2 *Engine) {
	if !reflect.DeepEqual(e1.Chesses, e2.Chesses) {
		//fmt.Printf("%+v\n", e1)
		e1.DrawGame()
		e2.DrawGame()
		panic("chesses not equal")
	}

	if !reflect.DeepEqual(e1.PosDict, e2.PosDict) {
		//fmt.Printf("%+v\n", e1)
		e1.DrawGame()
		e2.DrawGame()
		panic("PosDict not equal")
	}
}

func copyPosDic(p1 map[Pos]int) map[Pos]int {
	p2 := make(map[Pos]int)
	for k, v := range p1 {
		p2[k] = v
	}
	return p2
}

func copyEngine(e *Engine) *Engine {
	e1 := *e
	e1.Chesses = make([]*Chess, len(e.Chesses))
	for i := 0; i < len(e.Chesses); i++ {
		chess := *(e.Chesses[i])
		e1.Chesses[i] = &chess
	}
	e1.PosDict = copyPosDic(e.PosDict)
	return &e1
}

func getEngine() *Engine {
	e := NewEngine("1", "abcd")
	e.initChessGame([]Step{})
	return e
}

func TestCommon(t *testing.T) {
	//fmt.Println(subSameStatus("10110", "11011", 2))
	//fmt.Println(subSameStatus("10110", "11011", 1))
	//fmt.Println(subSameStatus("10110", "11011", 0))
}

func TestDeepEqual(t *testing.T) {
	e1 := getEngine()
	e2 := getEngine()
	fmt.Println(reflect.DeepEqual(e1.Chesses, e2.Chesses))

	step1 := *NewStep(0, Pos{1, 7}, Pos{1, 0})
	e1.ahead(step1)
	fmt.Println(reflect.DeepEqual(e1.Chesses, e2.Chesses))

	e2.ahead(step1)
	fmt.Println(reflect.DeepEqual(e1.Chesses, e2.Chesses))
}

func TestHasElement(t *testing.T) {
	arr := []int{1, 2, 3}
	fmt.Println(lib.HasElement[int](arr, 3))
	fmt.Println(lib.HasElement[int](arr, 0))
}

func TestIsCheckmate(t *testing.T) {
	e := getEngine()
	step1 := NewStep(0, Pos{2, 9}, Pos{4, 7})
	step2 := NewStep(0, Pos{1, 2}, Pos{1, 9})
	step3 := NewStep(0, Pos{0, 9}, Pos{1, 9})
	step4 := NewStep(0, Pos{7, 2}, Pos{4, 2})
	step5 := NewStep(0, Pos{7, 7}, Pos{7, 4})
	step6 := NewStep(0, Pos{4, 2}, Pos{4, 6})
	step7 := NewStep(0, Pos{7, 4}, Pos{4, 3})

	steps := []*Step{step1, step2, step3, step4, step5, step6, step7}
	for _, step := range steps {
		e.ahead(*step)
		e.changeCurCamp()
		fmt.Println(e.isCheckmate())
	}
}

func TestEvaluateValue(t *testing.T) {
	e := getEngine()
	step1 := *NewStep(0, Pos{2, 9}, Pos{4, 7})
	step2 := *NewStep(0, Pos{1, 2}, Pos{1, 9})
	step3 := *NewStep(0, Pos{0, 9}, Pos{1, 9})
	step4 := *NewStep(0, Pos{7, 2}, Pos{4, 2})
	step5 := *NewStep(0, Pos{7, 7}, Pos{7, 4})
	step6 := *NewStep(0, Pos{4, 2}, Pos{4, 6})
	step7 := *NewStep(0, Pos{7, 4}, Pos{4, 3})

	steps := []Step{step1, step2, step3, step4, step5, step6, step7}
	for _, step := range steps {

		fmt.Println(e.evaluateScore())
		e.ahead(step)
		fmt.Println(e.evaluateScore())
		//fmt.Println(e.isCheckmate())
	}
}

func TestGetAllSteps(t *testing.T) {
	e := getEngine()
	steps := e.getAllStep()

	sort.Slice(steps, func(i, j int) bool {
		chessI := e.Chesses[e.PosDict[steps[i].Start]]
		chessJ := e.Chesses[e.PosDict[steps[j].Start]]
		return chessI.name < chessJ.name
	})

	for i := 0; i < len(steps); i++ {
		fmt.Println(*steps[i])
	}
	fmt.Println(len(steps))
}

func TestGetNextStep(t *testing.T) {
	redispool.RedisPool = redispool.StartRedisPool()

	m := new(MachineMsg)
	m.Action = 0
	m.Room = "abcd"
	m.Start = map[string]interface{}{"x": 1.0, "y": 7.0}
	m.End = map[string]interface{}{"x": 1.0, "y": 0.0}

	res := GetChessStep(m)
	fmt.Println(res)

	defer DelRoom(m)
}

func TestDrawGame(t *testing.T) {
	e := getEngine()
	e.DrawGame()
}

func TestStepAndDraw(t *testing.T) {
	e := getEngine()
	step1 := *NewStep(0, Pos{1, 7}, Pos{1, 0})
	step2 := *NewStep(0, Pos{0, 0}, Pos{1, 0})
	step3 := *NewStep(0, Pos{7, 7}, Pos{7, 0})
	//step4 := *NewStep(0, Pos{7, 2}, Pos{4, 2})
	//step5 := *NewStep(0, Pos{7, 7}, Pos{7, 4})
	//step6 := *NewStep(0, Pos{4, 2}, Pos{4, 6})
	//step7 := *NewStep(0, Pos{7, 4}, Pos{4, 3})

	steps := []Step{step1, step2, step3}
	for _, step := range steps {
		e.DrawGame()
		e.ahead(step)
	}
	e.DrawGame()
}

func getLatestGame() *Engine {
	log.Printf("getLatestGame...\n\n")

	msg := &MachineMsg{Action: GO, Room: "abc"}
	steps := GetRoomSteps(msg.Room)
	log.Println("GetChessSteps steps: ", steps)
	//steps = steps[:len(steps)-1]

	camp := strconv.Itoa(len(steps) % 2) // 黑方0 红方1
	engine := NewEngine(camp, msg.Room)
	//lib.ReverseStatus(&engine.camp, "1")
	//lib.ReverseStatus(&engine.curCamp, "1")

	err := engine.initChessGame(steps)
	if err != nil {
		fmt.Println(err)
	}

	//recordTree.StoreInFile()
	return engine
}

func TestGame(t *testing.T) {
	fmt.Println("TestGame....")

	redispool.RedisPool = redispool.StartRedisPool()
	e := getLatestGame()
	e.DrawGame()
	fmt.Println(e.isCheckmate())
	e.getNextStep(config.MAX_DEPTH, -MAX_SCORE, MAX_SCORE)

	fmt.Printf("%+v\n", e.record.bestStep)

	//for i := 0; i < len(e.Chesses); i++ {
	//	chess := e.Chesses[i]
	//	if chess.name == HORSE && isSameCamp(chess.status, e.curCamp) && isAlive(chess.status) {
	//		fmt.Println(chessStepsFuncs[chess.name](e.Chesses, e.PosDict, i))
	//	}
	//}
}

func TestZobristKey(t *testing.T) {
	redispool.RedisPool = redispool.StartRedisPool()
	msg := &MachineMsg{Action: GO, Room: "abc"}
	steps := GetRoomSteps(msg.Room)
	log.Println("GetChessSteps steps: ", steps)
	//steps = steps[:len(steps)-1]

	step1 := *NewStep(0, Pos{1, 7}, Pos{4, 7})
	step2 := *NewStep(0, Pos{0, 0}, Pos{0, 1})
	step3 := *NewStep(0, Pos{4, 7}, Pos{1, 7})
	step4 := *NewStep(0, Pos{0, 1}, Pos{0, 0})

	steps = []Step{
		step1, step2, step3, step4,
	}

	camp := strconv.Itoa(len(steps) % 2) // 黑方0 红方1
	engine := NewEngine(camp, msg.Room)
	//lib.ReverseStatus(&engine.camp, "1")
	//lib.ReverseStatus(&engine.curCamp, "1")

	engine.initChessGame([]Step{})
	engine.DrawGame()

	z := make([]int64, 0)
	z = append(z, engine.ZobristKey(zobrist))
	for _, step := range steps {
		engine.ahead(step)
		key := engine.ZobristKey(zobrist)
		//engine.DrawGame()
		fmt.Println(key)
		if lib.HasElement[int64](z, key) {
			panic("")
		}
		z = append(z, key)
	}
}
