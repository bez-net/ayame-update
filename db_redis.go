package main

import (
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

func ExampleOps() {
	conn, err := connectRedis()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err = conn.Do("SET", "key", "Hello World!"); err != nil {
		log.Fatal(err)
	}

	str, err := redis.String(conn.Do("GET", "key"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
}

func connectRedis() (conn redis.Conn, err error) {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		addr = "localhost:6379" // default
	}
	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		conn, err = redis.Dial("tcp", addr) // default
	} else {
		conn, err = redis.Dial("tcp", addr, redis.DialPassword("stoney"))
	}
	return
}
