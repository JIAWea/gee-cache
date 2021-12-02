package geecache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestHTTPServe(t *testing.T) {
	var mockDB = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key: " + key)
		if v, ok := mockDB[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := "localhost:6889"
	peers := NewHTTPPool(addr)
	log.Println("gee-cache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
