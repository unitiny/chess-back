package lib

import (
	"fmt"
	"strconv"
	"strings"
)

type T interface {
	string | ~int | ~int64
}

// HasElement 切片中是否存在某值
func HasElement[t T](arr []t, value t) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == value {
			return true
		}
	}
	return false
}

// IsSameCamp 是否同阵营
func IsSameCamp(c1, c2 string) bool {
	return SubSameStatus(c1, c2, 0)
}

// Abs 绝对值
func Abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

// StringsToIntArr 字符串数组转化为整形数组
func StringsToIntArr(str string) []int {
	strArr := strings.Split(str, " ")

	arr := make([]int, len(strArr))
	for i := 0; i < len(arr); i++ {
		arr[i], _ = strconv.Atoi(strArr[i])
	}
	return arr
}

// InterfaceArrToString 字符态数组转变为字符串
func InterfaceArrToString(arr ...[]interface{}) string {
	res := ""
	for _, v := range arr {
		str := ""
		for _, val := range v {
			str += fmt.Sprintf("%v", val) + " "
		}
		res += str
	}
	return res[:len(res)-1]
}

// ReverseStatus 某位置取反
func ReverseStatus(status *string, target string) {
	statusNum, _ := strconv.ParseInt(*status, 2, 64)
	targetNum, _ := strconv.ParseInt(target, 2, 64)

	*status = fmt.Sprintf("%s", strconv.FormatInt(int64(statusNum^targetNum), 2))
}

// ResetStatus 某位置归零
func ResetStatus(status *string, target string) {
	statusNum, _ := strconv.ParseInt(*status, 2, 64)
	targetNum, _ := strconv.ParseInt(target, 2, 64)

	*status = fmt.Sprintf("%s", strconv.FormatInt(int64(statusNum & ^targetNum), 2))
}

// SetStatus 某位置归1
func SetStatus(status *string, target string) {
	statusNum, _ := strconv.ParseInt(*status, 2, 64)
	targetNum, _ := strconv.ParseInt(target, 2, 64)

	*status = fmt.Sprintf("%s", strconv.FormatInt(int64(statusNum|targetNum), 2))
}

// HasStatus 判断某些位置是否为1
func HasStatus(status *string, target string) bool {
	statusNum, _ := strconv.ParseInt(*status, 2, 64)
	targetNum, _ := strconv.ParseInt(target, 2, 64)

	return (statusNum & targetNum) == targetNum
}

// IsSameStatus 判断状态是否相同
func IsSameStatus(status *string, target string, offset int) bool {
	if len(*status) < offset+len(target) {
		fmt.Println("IsSameStatus 长度越界")
	}

	start := len(*status) - len(target) - offset
	end := len(*status) - offset
	return (*status)[start:end] == target
}

// SubSameStatus 两字符串某位置状态是否相同
func SubSameStatus(s1, s2 string, offset int) bool {
	return s1[len(s1)-1-offset] == s2[len(s2)-1-offset]
}

// IsAlive 是否存活
func IsAlive(status string) bool {
	return HasStatus(&status, "10")
}
