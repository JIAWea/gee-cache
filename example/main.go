package main

import (
	"flag"
	"fmt"
	"gee-cache"
	"log"
	"net/http"
)

var mockDB = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key ", key)
			if value, ok := mockDB[key]; ok {
				return []byte(value), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrList []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrList...)
	gee.RegisterPeers(peers)
	log.Println("gee-cache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(view.ByteSlice())
		}))
	log.Println("frontend server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gee-cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "0.0.0.0:9999"
	addrMap := map[int]string{
		8001: "0.0.0.0:8001",
		8002: "0.0.0.0:8002",
		8003: "0.0.0.0:8003",
	}

	var addrList []string
	for _, v := range addrMap {
		addrList = append(addrList, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrList, gee)
}
