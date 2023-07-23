package lib

type Status interface {
	Reverse(target ...int)
	Reset(target ...int)
	Set(target ...int)
	Has(target int) bool
	IsSame(target int, offset int) bool
}

type BasicStatus struct {
	Status int
}

// Reverse 翻转状态
func (c *BasicStatus) Reverse(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res ^= target[i]
	}
	c.Status = res
}

// Reset 重置状态
func (c *BasicStatus) Reset(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res &= ^target[i]
	}
	c.Status = res
}

// Set 设置状态
func (c *BasicStatus) Set(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res |= target[i]
	}
	c.Status = res
}

// Has 目标状态位是否为1
func (c *BasicStatus) Has(target int) bool {
	if c.Status&target == 0 {
		return false
	}
	return true
}

// IsSame 判断子状态是否相同
func (c *BasicStatus) IsSame(target int, offset int) bool {
	if target == 0 {
		return c.Status&1 == target
	}

	status := c.Status
	status = status >> offset

	i := 1
	for target >= i {
		s1 := status & i
		s2 := target & i
		if s1^s2 != 0 {
			return false
		}
		i = i << 1
	}
	return true
}
