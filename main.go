package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"net/http"
)

var (
	redisPool      *redis.Pool
	redisAddress   = flag.String("redis-address", "redis:6379", "Address to the Redis server")
	maxConnections = flag.Int("max-connections", 10, "Max connections to Redis")
)

type Hitcount struct {
	RealIPAddr string `json:"User Address "`
	Count      int64  `json:"Hit "`
}

func main() {
	flag.Parse()

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", *redisAddress)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	redisConn := redisPool.Get()
	defer redisConn.Close()

	clientIPAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIPAddress = err.Error()
	}

	count, err := redisConn.Do("INCR", "count")
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"Could not increment count: %v\"}", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	hitcount := Hitcount{Count: count.(int64), RealIPAddr: string(clientIPAddress)}
	//hitcount := Hitcount{Count: count.(int64)}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(hitcount); err != nil {
		panic(err)
	}
}
