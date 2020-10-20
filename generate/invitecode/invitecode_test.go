package invitecode

import (
	"fmt"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	uid := uint64(18739)
	fmt.Println(Encode(uid))
}

func TestDecode(t *testing.T) {
	code := "MALF73"
	fmt.Println(Decode(code))
}

func TestEncodeAndDecode(t *testing.T) {
	s := time.Now()
	count := 1000000
	existsMap := map[string]bool{}
	for count > 0 {
		count--
		uid := uint64(count)
		code := Encode(uid)
		rawUid := Decode(code)
		if uid != rawUid {
			t.Fatal(uid, code, rawUid)
		}
		if _, ok := existsMap[code]; ok {
			t.Fatal("exist")
		} else {
			existsMap[code] = true
		}
	}
	fmt.Println("time: ", time.Now().Sub(s))
}
