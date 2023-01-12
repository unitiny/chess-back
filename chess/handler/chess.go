package handler

import (
	"Chess/chess/config"
	"Chess/chess/lib"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
)

type Type int

type Chess struct {
	id        int
	x         int
	y         int
	status    *ChessStatus
	name      Type
	dieStepId int // 被吃掉时所处步数
}

// Pos 位置
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Step 动作
type Step struct {
	Id    int `json:"id"`
	Start Pos `json:"start"`
	End   Pos `json:"end"`
}

// Record 记录
type Record struct {
	maxValue int    // 最优走法价值
	curStep  []Step // 当前递归走法
	bestStep Step   // 最优走法
}

// Zobrist 储存棋局 每个位置存储所有状态
type Zobrist map[Pos][POS_STATE]int

var zobrist Zobrist

var chessStepsFuncs map[Type]ChessSteps

type HashGame struct {
	depth int
	value int
	flag  int
}

// Engine 全局引擎
type Engine struct {
	depth       int                 // 当前递归深度
	gameSteps   int                 // 棋局步数
	redKingId   int                 // 红方将军id
	blockKingId int                 // 黑方将军id
	camp        int                 // 机器所属阵营
	curCamp     int                 // 当前计算阵营
	room        string              // 所属房间
	record      *Record             // 记录结果
	Chesses     []*Chess            // 所有棋子
	PosDict     map[Pos]int         // 记录棋子位置
	posValues   map[Type][][]int    // 记录棋子得分
	diePosDict  map[Pos][]int       // 记录被吃棋子位置
	hashTable   map[int64]*HashGame // 历史棋局散列表
}

// 初始化
func init() {

	chessStepsFuncs = make(map[Type]ChessSteps)
	chessStepsFuncs[SOLDIER] = new(soldier)
	chessStepsFuncs[CANNON] = new(cannon)
	chessStepsFuncs[CAR] = new(car)
	chessStepsFuncs[HORSE] = new(horse)
	chessStepsFuncs[ELEPHANT] = new(elephant)
	chessStepsFuncs[SCHOLAR] = new(scholar)
	chessStepsFuncs[KING] = new(king)

	// 给每个Pos的每个状态赋值一个随机数
	zobrist = make(map[Pos][POS_STATE]int, 0)
	for i := 0; i < COLUMNS; i++ {
		for j := 0; j < ROWS; j++ {
			var states [POS_STATE]int
			for k := 0; k < POS_STATE; k++ {
				states[k] = rand.Int()
			}
			pos := *NewPos(j, i)
			zobrist[pos] = states
		}
	}
}

func NewEngine(camp string, room string) *Engine {
	record := new(Record)
	record.maxValue = math.MinInt

	posValues := make(map[Type][][]int, 0)
	posValues[SOLDIER] = chessValueB
	posValues[CANNON] = chessValueP
	posValues[CAR] = chessValueC
	posValues[HORSE] = chessValueM
	posValues[ELEPHANT] = chessValueX
	posValues[SCHOLAR] = chessValueS
	posValues[KING] = chessValueJ

	c, _ := strconv.Atoi(camp)
	return &Engine{
		record:  record,
		camp:    c,
		curCamp: c,
		room:    room,
		Chesses: []*Chess{
			// 红方
			{x: 0, y: 6, name: SOLDIER},
			{x: 2, y: 6, name: SOLDIER},
			{x: 4, y: 6, name: SOLDIER},
			{x: 6, y: 6, name: SOLDIER},
			{x: 8, y: 6, name: SOLDIER},
			{x: 1, y: 7, name: CANNON},
			{x: 7, y: 7, name: CANNON},
			{x: 0, y: 9, name: CAR},
			{x: 8, y: 9, name: CAR},
			{x: 1, y: 9, name: HORSE},
			{x: 7, y: 9, name: HORSE},
			{x: 2, y: 9, name: ELEPHANT},
			{x: 6, y: 9, name: ELEPHANT},
			{x: 3, y: 9, name: SCHOLAR},
			{x: 5, y: 9, name: SCHOLAR},
			{x: 4, y: 9, name: KING},

			// 黑方
			{x: 4, y: 0, name: KING},
			{x: 3, y: 0, name: SCHOLAR},
			{x: 5, y: 0, name: SCHOLAR},
			{x: 2, y: 0, name: ELEPHANT},
			{x: 6, y: 0, name: ELEPHANT},
			{x: 1, y: 0, name: HORSE},
			{x: 7, y: 0, name: HORSE},
			{x: 0, y: 0, name: CAR},
			{x: 8, y: 0, name: CAR},
			{x: 1, y: 2, name: CANNON},
			{x: 7, y: 2, name: CANNON},
			{x: 0, y: 3, name: SOLDIER},
			{x: 2, y: 3, name: SOLDIER},
			{x: 4, y: 3, name: SOLDIER},
			{x: 6, y: 3, name: SOLDIER},
			{x: 8, y: 3, name: SOLDIER},
		},
		posValues:  posValues,
		PosDict:    make(map[Pos]int),
		diePosDict: make(map[Pos][]int),
		hashTable:  make(map[int64]*HashGame),
	}
}

// GetChessStep 得到最优结果并返回
func GetChessStep(msg *MachineMsg) string {
	log.Printf("GetChessStep...\n\n")

	steps := GetRoomSteps(msg.Room)
	if len(msg.Start) != 0 && len(msg.End) != 0 {
		steps = append(steps, Step{Id: len(steps), Start: *mapToPos(msg.Start), End: *mapToPos(msg.End)})
	}
	log.Println("GetChessSteps steps: ", steps)

	camp := strconv.Itoa((len(steps) % 2) ^ 1) // 黑方0 红方1
	engine := NewEngine(camp, msg.Room)

	err := engine.initChessGame(steps)
	if err != nil {
		return err.Error()
	}

	engine.DrawGame()
	engine.getNextStep(config.MAX_DEPTH, -MAX_SCORE, MAX_SCORE)

	if len(steps) != 0 {
		RecordRoomStep(msg.Room, steps[len(steps)-1], engine.record.bestStep)
	} else {
		RecordRoomStep(msg.Room, engine.record.bestStep)
	}

	result, _ := json.Marshal(engine.record.bestStep)
	return string(result)
}

// 初始化棋局
func (e *Engine) initChessGame(steps []Step) error {
	// 字典储存位置，并赋值id, step有index可改进
	for k, v := range e.Chesses {
		pos := Pos{v.x, v.y}
		e.PosDict[pos] = k
		e.Chesses[k].id = k

		if e.Chesses[k].status == nil {
			e.Chesses[k].status = new(ChessStatus)
		}
		e.Chesses[k].status.Set(RED, ALIVE)

		// 记录kingId
		if v.name == KING {
			if k > 15 {
				e.blockKingId = k
			} else {
				e.redKingId = k
			}
		}

		if k > 15 {
			e.Chesses[k].status.Reset(RED)
		}
	}

	// 移动棋子
	for _, step := range steps {
		// 改变位置
		Start, has := e.PosDict[step.Start]
		if !has {
			return errors.New("init fail: none start pos")
		}
		e.Chesses[Start].x = step.End.X
		e.Chesses[Start].y = step.End.Y

		// 若目标位置有棋子，则被吃
		if end, ok := e.PosDict[step.End]; ok {
			e.Chesses[end].status.Reset(ALIVE)
		}

		// 更新字典
		delete(e.PosDict, step.Start)
		e.PosDict[step.End] = Start
	}
	e.gameSteps = len(steps)
	return nil
}

/* dfs + 回溯，得到人机结果
*  先深入到叶子节点，然后回溯，得出每层的最优解，
*  到根节点自然得到最优走法了
 */
func (e *Engine) getNextStep(depth, Alpha, Beta int) int {
	// 0 查找历史表
	flag := HASH_ALPHA
	if val := e.probeHashGame(depth, Alpha, Beta); val != VAL_UNKNOWN {
		return val
	}

	// 1 终止条件
	if depth == 0 {
		score := e.evaluateScore()
		e.recordHashGame(depth, score, HASH_EXACT) // 记录分数
		return score                               // 返回当前阵营局势分数
	}

	// 2 遍历所有走法，得出最优几种
	e.depth = depth
	steps := e.getAllStep()

	// 3 逐个尝试走法
	for i := 0; i < len(steps); i++ {
		// 3.1 移动棋局并继续深入
		step := *steps[i]
		e.ahead(step)
		value := -e.getNextStep(depth-1, -Beta, -Alpha)

		// 3.2 回溯复原
		e.depth = depth
		e.back(step)

		// 3.3 对该走法评估
		if value >= Beta {
			e.recordHashGame(depth, Beta, HASH_BETA)
			return Beta
		}

		// 更新机器最优走法分数
		if value > Alpha {
			Alpha = value
			flag = HASH_EXACT
			// 若为回溯到栈底，则记录最佳走法
			if depth == config.MAX_DEPTH {
				//fmt.Println("bestStep", step)
				e.record.bestStep = step
			}
		}
	}

	e.recordHashGame(depth, Alpha, flag) // 记录到历史表中
	return Alpha
}

// 移动棋子
func (e *Engine) goStep(step Step, isBack bool) {
	// 移动
	i := e.PosDict[step.Start]
	e.Chesses[i].x = step.End.X
	e.Chesses[i].y = step.End.Y

	// 终点位置有对方阵营棋子，则吃掉
	k, has := e.PosDict[step.End]
	if has && !isSameCamp(e.Chesses[i].status, e.Chesses[k].status) {
		e.Chesses[k].status.Reset(ALIVE)
		e.Chesses[k].dieStepId = e.getStepId()
		e.diePosDict[step.End] = append(e.diePosDict[step.End], k)
	}

	e.PosDict[step.End] = i // 之前值被覆盖 相对于被吃掉
	delete(e.PosDict, step.Start)

	// 悔棋
	if isBack {
		// 肯定是复活最后添加的棋子
		l := len(e.diePosDict[step.Start])
		if l >= 1 {
			v := e.diePosDict[step.Start][l-1]
			//fmt.Println(step, e.Chesses[v].dieStepId, e.getStepId())
			if e.Chesses[v].dieStepId != e.getStepId() {
				return // 若不是上一步，则无法悔棋
			}

			// 复活棋子
			e.Chesses[v].status.Set(ALIVE)
			e.Chesses[v].dieStepId = 0
			e.PosDict[step.Start] = v

			// 删除索引
			e.diePosDict[step.Start] = e.diePosDict[step.Start][:l-1]
			if l == 1 {
				delete(e.diePosDict, step.Start)
			}
		}
	}
}

// 回退上一步
func (e *Engine) back(step Step) {
	backStep := reverseStep(step)
	e.goStep(backStep, true)
	e.changeCurCamp()
}

// 前进
func (e *Engine) ahead(step Step) {
	e.goStep(step, false)
	e.changeCurCamp()
}

// 获取所有走法
func (e *Engine) getAllStep() []*Step {
	steps := make([]*Step, 0)
	for i := 0; i < len(e.Chesses); i++ {
		chess := e.Chesses[i]
		if !isSameCamp2(chess.status, e.curCamp) || !isAlive(chess.status) {
			continue
		}

		startPos := NewPos(chess.x, chess.y)
		allPos := e.getChessSteps(i) // 获取该棋子所有走
		for _, pos := range allPos {
			steps = append(steps, NewStep(0, *startPos, pos))
		}
	}

	// 将走法按分值排序
	type valueStep struct {
		value int
		step  Step
	}
	valueSteps := make([]*valueStep, 0)
	for i := 0; i < len(steps); i++ {
		e.ahead(*steps[i])

		if e.isCheckmate() {
			e.back(*steps[i])
			continue // 过滤被将军走法
		}
		score := -e.evaluateScore()
		vs := &valueStep{score, *steps[i]}
		valueSteps = append(valueSteps, vs)

		e.back(*steps[i])
	}

	sort.Slice(valueSteps, func(i, j int) bool {
		return valueSteps[i].value > valueSteps[j].value
	})

	resSteps := make([]*Step, len(valueSteps))
	for i := 0; i < len(valueSteps); i++ {
		resSteps[i] = &valueSteps[i].step
	}
	return resSteps
}

// 获取棋子可走位置
func (e *Engine) getChessSteps(index int) []Pos {
	c := e.Chesses[index]
	res := chessStepsFuncs[c.name].Steps(e.Chesses, e.PosDict, index)
	return res
}

// 获取棋子得分
func (e *Engine) getScore(c Chess) int {
	pos := *NewPos(c.x, c.y)
	if !isSameCamp2(c.status, e.camp) {
		pos = reversePos(pos) // 非同阵营翻转位置
	}

	return e.posValues[c.name][pos.Y][pos.X]
}

// 统计当前阵营棋局分数
func (e *Engine) evaluateScore() int {
	score := 0
	for i := 0; i < len(e.Chesses); i++ {
		c := e.Chesses[i]
		if !isAlive(c.status) {
			continue
		}

		value := e.getScore(*c)
		if !isSameCamp2(c.status, e.curCamp) {
			value *= -1
		}
		score += value
	}
	return score
}

// 改变阵营
func (e *Engine) changeCurCamp() {
	e.curCamp ^= 1
}

// 获取敌方阵营将军索引
func (e *Engine) getEnemyKingId() int {
	if e.curCamp == BLACK {
		return e.redKingId
	}
	return e.blockKingId
}

// 判断当前阵营是否将军
func (e *Engine) isCheckmate() bool {
	// 只需遍历己方的车 马 炮 将军 兵
	someChessType := []Type{CAR, HORSE, CANNON, KING, SOLDIER}
	king := e.Chesses[e.getEnemyKingId()]
	kingPos := NewPos(king.x, king.y)

	for i := 0; i < len(e.Chesses); i++ {
		chess := e.Chesses[i]
		if !isSameCamp2(chess.status, e.curCamp) || !isAlive(chess.status) {
			continue // 不同阵营或已被吃掉棋子跳过
		}

		if lib.HasElement[Type](someChessType, chess.name) {
			pos := chessStepsFuncs[chess.name].Steps(e.Chesses, e.PosDict, chess.id)
			// 若有杀到将军走法，则将军了
			for _, p := range pos {
				if reflect.DeepEqual(p, *kingPos) {
					return true
				}
			}
		}
	}
	return false
}

// 获得当前棋局步数id
func (e *Engine) getStepId() int {
	return config.MAX_DEPTH - e.depth
}

// DrawGame 绘制当前棋局
func (e *Engine) DrawGame() {
	name := make(map[Type]string)
	name[SOLDIER] = "兵"
	name[CANNON] = "炮"
	name[CAR] = "车"
	name[HORSE] = "马"
	name[ELEPHANT] = "象"
	name[SCHOLAR] = "士"
	name[KING] = "将"

	for i := 0; i < COLUMNS; i++ {
		for j := 0; j < ROWS; j++ {
			pos := *NewPos(j, i)
			if index, ok := e.PosDict[pos]; ok {
				chess := e.Chesses[index]
				if !isAlive(chess.status) {
					fmt.Printf("%-4s ", "*")
					continue
				}
				if isSameCamp2(chess.status, RED) {
					fmt.Printf("%c[%d;%d;%dm%-3s %c[0m", 0x1B, 1, 0, 31, name[chess.name], 0x1B)
				} else {
					fmt.Printf("%-3s ", name[chess.name])
				}
			} else {
				fmt.Printf("%-4s ", "*")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// ZobristKey 返回棋局唯一key
func (e *Engine) ZobristKey(zobrist Zobrist) int64 {
	// 从zobrist中异或每个格子状态，得出唯一key
	key := 0
	for i := 0; i < COLUMNS; i++ {
		for j := 0; j < ROWS; j++ {
			pos := *NewPos(j, i)
			if v, has := e.PosDict[pos]; has {
				chess := e.Chesses[v]
				index := chess.name
				if isSameCamp2(chess.status, RED) {
					index += CHESS_CATE
				}
				key ^= zobrist[pos][index]
			} else {
				key ^= zobrist[pos][POS_STATE-1] // 异或无子空格
			}
		}
	}
	return int64(key)
}

// probeHashGame 获取历史评分
func (e *Engine) probeHashGame(depth, alpha, beta int) int {
	key := e.ZobristKey(zobrist)

	// 验证历史表是否存在该局面，以及是否同深度
	if hg, has := e.hashTable[key]; has {
		if hg.depth == depth {
			// 根据标签返回不同值
			if hg.flag == HASH_EXACT {
				return hg.value
			} else if hg.flag == HASH_ALPHA && hg.value <= alpha {
				return alpha
			} else if hg.flag == HASH_BETA && hg.value >= beta {
				return beta
			}
		}
	}
	return VAL_UNKNOWN
}

// recordHashGame 记录历史评分
func (e *Engine) recordHashGame(depth, val, flag int) {
	key := e.ZobristKey(zobrist)
	hashGame := &HashGame{
		depth: depth,
		value: val,
		flag:  flag,
	}
	e.hashTable[key] = hashGame
}
