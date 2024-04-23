package ConsistentHash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewNodeMap(func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}, 3)

	// Given the above hash function, this will give replicas with "hashes":
	// 02,04,06,12,14,16,22,24,26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"27": "2",
		"23": "4",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
