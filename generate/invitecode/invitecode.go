package invitecode

import (
	"math/rand"
	"strings"
	"time"
)

var (
	baseChar    = "123456789ABCDEGHJKLMNPQRSTUVWXYZ"
	baseDecimal = uint64(32)
	baseLen     = 6
	basePad     = "F"
)

func Encode(id uint64) string {
	res := ""
	for id != 0 {
		mod := id % baseDecimal
		id = id / baseDecimal
		res += string(baseChar[mod])
	}
	if len(res) < baseLen {
		res += basePad
		l := len(res)
		for i := 0; i < baseLen-l; i++ {
			rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100000)))
			res += string(baseChar[rand.Intn(int(baseDecimal))])
		}
	}
	return res
}

func Decode(code string) uint64 {
	res := uint64(0)
	lenCode := len(code)

	baseArr := []byte(baseChar)   // 字符串进制转换为byte数组
	baseRev := make(map[byte]int) // 进制数据键值转换为map
	for k, v := range baseArr {
		baseRev[v] = k
	}

	// 查找补位字符的位置
	isPad := strings.Index(code, basePad)
	if isPad != -1 {
		lenCode = isPad
	}
	r := 0
	for i := 0; i < lenCode; i++ {
		// 补充字符直接跳过
		if string(code[i]) == basePad {
			continue
		}
		index := baseRev[code[i]]
		b := uint64(1)
		for j := 0; j < r; j++ {
			b *= baseDecimal
		}
		res += uint64(index) * b
		r++
	}

	return res
}
