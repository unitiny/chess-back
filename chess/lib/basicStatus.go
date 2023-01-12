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

func (c *BasicStatus) Reverse(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res ^= target[i]
	}
	c.Status = res
}

func (c *BasicStatus) Reset(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res &= ^target[i]
	}
	c.Status = res
}

func (c *BasicStatus) Set(target ...int) {
	res := c.Status
	for i := 0; i < len(target); i++ {
		res |= target[i]
	}
	c.Status = res
}

func (c *BasicStatus) Has(target int) bool {
	if c.Status&target == 0 {
		return false
	}
	return true
}

func (c *BasicStatus) IsSame(target int, offset int) bool {
	if target == 0 {
		return c.Status&1 == 0
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
