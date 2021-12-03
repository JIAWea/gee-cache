package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 每个节点有 2 个虚拟节点，通过上面的 hash 计算得出
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("2", "6", "4")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"13": "4",
		"14": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s should have yieled %s", k, v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s should have yieled %s", k, v)
		}
	}

}
