package handler

func NewStep(id int, Start Pos, end Pos) *Step {
	return &Step{Id: id, Start: Start, End: end}
}

func NewPos(x int, y int) *Pos {
	return &Pos{X: x, Y: y}
}

// 翻转位置
func reversePos(pos Pos) Pos {
	pos.X = 8 - pos.X
	pos.Y = 9 - pos.Y
	return pos
}

// 反转step
func reverseStep(step Step) Step {
	newStep := &Step{Id: step.Id, Start: step.End, End: step.Start}
	return *newStep
}

// map转化为Pos
func mapToPos(m map[string]interface{}) *Pos {
	return &Pos{X: int(m["x"].(float64)), Y: int(m["y"].(float64))}
}

// 获取阵营底线
func getBottom(status *ChessStatus) int {
	if status.Has(RED) {
		return RED_BOTTOM // 红方阵营底线
	}
	return BLOCK_BOTTOM
}
