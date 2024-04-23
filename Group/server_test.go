package Group

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

//var db = map[string]string{
//	"Tom":  "630",
//	"Jack": "589",
//	"Sam":  "567",
//}

func TestServer(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	NewGroup("scores", 1024, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	node := NewhttpNode("127.0.0.1:9090")
	s := http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: node,
	}

	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
