package handler

import (
	"Chess/chess/lib"
	"math"
)

// ChessSteps 查找所有走法
type ChessSteps interface {
	Steps(c []*Chess, posDic map[Pos]int, index int) []Pos
}

type king struct{}

func (s *king) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	bottom := getBottom(chess.status)

	res := make([]Pos, 0)
	offsets := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	for _, offset := range offsets {
		end := NewPos(chess.x+offset[0], chess.y+offset[1])

		// 越界判断
		if end.X < 3 || end.X > 5 ||
			math.Abs(float64(end.Y-bottom)) > 3 || end.Y < 0 || end.Y > 9 {
			continue
		}

		// 是否有己方棋子占位
		if v, has := posDic[*end]; has {
			if isSameCamp(c[v].status, chess.status) {
				continue
			}
		}

		res = append(res, *end)
	}

	// 攻击对方将军
	for i := 1; i <= RED_BOTTOM; i++ {
		pos := Pos{}
		if bottom == RED_BOTTOM {
			pos = Pos{X: chess.x, Y: chess.y - i}
		} else {
			pos = Pos{X: chess.x, Y: chess.y + i}
		}
		if index, ok := posDic[pos]; ok {
			if !c[index].status.Has(ALIVE) {
				continue
			}
			if c[index].name == KING {
				res = append(res, pos)
			}
			break
		}
	}

	return res
}

type scholar struct{}

func (s *scholar) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)
	bottom := getBottom(chess.status)

	offsets := [][]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, offset := range offsets {
		end := Pos{chess.x + offset[0], chess.y + offset[1]}

		// 越界判断
		if end.X < 3 || end.X > 5 ||
			end.Y < 0 || end.Y > 9 || math.Abs(float64(end.Y-bottom)) > 2 {
			continue
		}

		// 被己方占位判断
		if v, has := posDic[end]; has {
			if isSameCamp(c[v].status, chess.status) {
				continue
			}
		}

		res = append(res, end)
	}

	return res
}

type elephant struct{}

func (s *elephant) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)
	bottom := getBottom(chess.status)

	offsets := [][]int{{2, 2}, {2, -2}, {-2, 2}, {-2, -2}}
	for _, offset := range offsets {
		end := Pos{chess.x + offset[0], chess.y + offset[1]}

		// 越界判断
		if end.X < 0 || end.X > 8 ||
			end.Y < 0 || end.Y > 9 || math.Abs(float64(end.Y-bottom)) > 4 {
			continue
		}

		// 被己方占位判断
		if v, has := posDic[end]; has {
			if isSameCamp(c[v].status, chess.status) {
				continue
			}
		}

		// 被挡象脚判断
		resistPos := Pos{chess.x + offset[0]/2, chess.y + offset[1]/2}
		if _, has := posDic[resistPos]; has {
			continue
		}

		res = append(res, end)
	}

	return res
}

type horse struct{}

func (s *horse) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)

	offsets := [][]int{{1, 2}, {1, -2}, {-1, 2}, {-1, -2}, {2, 1}, {2, -1}, {-2, 1}, {-2, -1}}
	for _, offset := range offsets {
		end := Pos{chess.x + offset[0], chess.y + offset[1]}

		// 越界判断
		if end.X < 0 || end.X > 8 ||
			end.Y < 0 || end.Y > 9 {
			continue
		}

		// 被己方占位判断
		if v, has := posDic[end]; has {
			if isSameCamp(c[v].status, chess.status) {
				continue
			}
		}

		// 被挡马脚判断
		resistPos := Pos{chess.x + offset[0]/2, chess.y + offset[1]/2}
		if _, has := posDic[resistPos]; has {
			continue
		}

		res = append(res, end)
	}

	return res
}

type car struct{}

func (s *car) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)

	// 上检索
	for i := chess.y - 1; i >= 0; i-- {
		pos := NewPos(chess.x, i)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			res = append(res, *pos)
		} else {
			// 被遮挡，遮挡棋子若为敌方则添加，并终止循环
			endChess := c[v]
			if !isSameCamp(chess.status, endChess.status) {
				res = append(res, *pos)
			}
			break
		}
	}

	// 下检索
	for i := chess.y + 1; i <= 9; i++ {
		pos := NewPos(chess.x, i)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			res = append(res, *pos)
		} else {
			// 被遮挡，遮挡棋子若为敌方则添加，并终止循环
			endChess := c[v]
			if !isSameCamp(chess.status, endChess.status) {
				res = append(res, *pos)
			}
			break
		}
	}

	// 右检索
	for i := chess.x + 1; i <= 8; i++ {
		pos := NewPos(i, chess.y)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			res = append(res, *pos)
		} else {
			// 被遮挡，遮挡棋子若为敌方则添加，并终止循环
			endChess := c[v]
			if !isSameCamp(chess.status, endChess.status) {
				res = append(res, *pos)
			}
			break
		}
	}

	// 左检索
	for i := chess.x - 1; i >= 0; i-- {
		pos := NewPos(i, chess.y)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			res = append(res, *pos)
		} else {
			// 被遮挡，遮挡棋子若为敌方则添加，并终止循环
			endChess := c[v]
			if !isSameCamp(chess.status, endChess.status) {
				res = append(res, *pos)
			}
			break
		}
	}

	return res
}

type cannon struct{}

func (s *cannon) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)

	// 上检索
	hasResist := false
	for i := chess.y - 1; i >= 0; i-- {
		pos := NewPos(chess.x, i)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			if !hasResist {
				res = append(res, *pos)
			}
		} else {
			endChess := c[v]
			// 如果已被遮挡，则炮可打有敌方棋子的位置，并终止循环
			if hasResist {
				if !isSameCamp(chess.status, endChess.status) {
					res = append(res, *pos)
				}
				break
			}
			hasResist = true
		}
	}

	// 下检索
	hasResist = false
	for i := chess.y + 1; i <= 9; i++ {
		pos := NewPos(chess.x, i)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			if !hasResist {
				res = append(res, *pos)
			}
		} else {
			endChess := c[v]
			// 如果已被遮挡，则炮可打有棋子的位置，并终止循环
			if hasResist {
				if !isSameCamp(chess.status, endChess.status) {
					res = append(res, *pos)
				}
				break
			}
			hasResist = true
		}
	}

	// 右检索
	hasResist = false
	for i := chess.x + 1; i <= 8; i++ {
		pos := NewPos(i, chess.y)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			if !hasResist {
				res = append(res, *pos)
			}
		} else {
			endChess := c[v]
			// 如果已被遮挡，则炮可打有棋子的位置，并终止循环
			if hasResist {
				if !isSameCamp(chess.status, endChess.status) {
					res = append(res, *pos)
				}
				break
			}
			hasResist = true
		}
	}

	// 左检索
	hasResist = false
	for i := chess.x - 1; i >= 0; i-- {
		pos := NewPos(i, chess.y)
		// 空位且未被遮挡
		if v, has := posDic[*pos]; !has {
			if !hasResist {
				res = append(res, *pos)
			}
		} else {
			endChess := c[v]
			// 如果已被遮挡，则炮可打有棋子的位置，并终止循环
			if hasResist {
				if !isSameCamp(chess.status, endChess.status) {
					res = append(res, *pos)
				}
				break
			}
			hasResist = true
		}
	}

	return res
}

// cannonAttacks 查找所有炮能攻击到的位置
func (s *cannon) cannonAttacks(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)

	// 上检索
	hasResist := false
	for i := chess.y - 1; i >= 0; i-- {
		pos := NewPos(chess.x, i)
		// 空位且被遮挡
		if _, has := posDic[*pos]; !has {
			if hasResist {
				res = append(res, *pos)
			}
		} else {
			// 如果第二次被遮挡，加入打击位置中，并终止循环
			if hasResist {
				res = append(res, *pos)
				break
			}
			hasResist = true // 第一次被遮挡
		}
	}

	// 下检索
	hasResist = false
	for i := chess.y + 1; i <= 9; i++ {
		pos := NewPos(chess.x, i)
		// 空位且被遮挡
		if _, has := posDic[*pos]; !has {
			if hasResist {
				res = append(res, *pos)
			}
		} else {
			// 如果第二次被遮挡，加入打击位置中，并终止循环
			if hasResist {
				res = append(res, *pos)
				break
			}
			hasResist = true // 第一次被遮挡
		}
	}

	// 右检索
	hasResist = false
	for i := chess.x + 1; i <= 8; i++ {
		pos := NewPos(i, chess.y)
		// 空位且被遮挡
		if _, has := posDic[*pos]; !has {
			if hasResist {
				res = append(res, *pos)
			}
		} else {
			// 如果第二次被遮挡，加入打击位置中，并终止循环
			if hasResist {
				res = append(res, *pos)
				break
			}
			hasResist = true // 第一次被遮挡
		}
	}

	// 左检索
	hasResist = false
	for i := chess.x - 1; i >= 0; i-- {
		pos := NewPos(i, chess.y)
		// 空位且被遮挡
		if _, has := posDic[*pos]; !has {
			if hasResist {
				res = append(res, *pos)
			}
		} else {
			// 如果第二次被遮挡，加入打击位置中，并终止循环
			if hasResist {
				res = append(res, *pos)
				break
			}
			hasResist = true // 第一次被遮挡
		}
	}

	return res
}

type soldier struct{}

func (s *soldier) Steps(c []*Chess, posDic map[Pos]int, index int) []Pos {
	chess := c[index]
	res := make([]Pos, 0)
	bottom := getBottom(chess.status)

	// 处于己方地盘只能前进
	if lib.Abs(chess.y-bottom) <= 4 {
		y := lib.Abs(bottom - (lib.Abs(chess.y-bottom) + 1))
		end := Pos{chess.x, y}
		// 被己方占位判断
		if v, has := posDic[end]; has {
			if isSameCamp(c[v].status, chess.status) {
				return res
			}
		}
		res = append(res, end)
		return res
	}

	offsets := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, offset := range offsets {
		end := Pos{chess.x + offset[0], chess.y + offset[1]}

		// 越界判断
		if end.X < 0 || end.X > 8 ||
			end.Y < 0 || end.Y > 9 ||
			math.Abs(float64(end.Y-bottom)) <= 4 ||
			math.Abs(float64(end.Y-bottom)) < math.Abs(float64(chess.y-bottom)) {
			continue
		}

		// 被己方占位判断
		if v, has := posDic[end]; has {
			if isSameCamp(c[v].status, chess.status) {
				continue
			}
		}

		res = append(res, end)
	}

	return res
}
