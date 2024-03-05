package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
)

func hello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	io.WriteString(w, fmt.Sprintf("[v3] Hello, Kubernetes!, From host: %s\r\n", host))
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/redis", redisHandler)
	http.ListenAndServe(":3000", nil)
}

func redisHandler(w http.ResponseWriter, r *http.Request) {
	pool := &redis.Pool{
		MaxIdle: 100,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", "10.105.184.56:6379", redis.DialPassword("zBaMX2k6"))
			return
		},
	}
	conn := pool.Get()
	defer conn.Close()
	countKey := "count_key"
	n, err := conn.Do("INCR", countKey)
	if err != nil {
		fmt.Println("INCR RES:", n, err)
	}
	host, _ := os.Hostname()
	io.WriteString(w, fmt.Sprintf("[v3] Hello, Kubernetes!, From host: %s redis_key: %v redis_reply: %v\r\n", host, countKey, n))
}