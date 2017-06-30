// MessageBox project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	redisAddr = os.Getenv("REDISADDR")
	redisPs   = ""
	redisDB   = os.Getenv("REDISDB")
)

type message struct {
	Content string `json:"content"`
}

func testSend(hub *Hub) {
	c, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	for i := 0; ; i++ {
		j := message{fmt.Sprintf("%d - the time is %v", i, time.Now())}
		js, _ := json.Marshal(j)
		c.Do("PUBLISH", "messagebox:sse", js)
		log.Printf("Sent message %d ", i)
		time.Sleep(3 * 1e9)
	}
}

func main() {
	hub, err := newHub()

	go hub.start()

	if err != nil {
		log.Println(err)
		return
	}

	//go testSend(hub)

	router := gin.Default()

	router.GET("/messagebox", func(c *gin.Context) {
		listener := hub.registerListener()

		defer hub.closeListener(listener)
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
		c.Stream(func(w io.Writer) bool {
			message := <-listener
			c.SSEvent("message", string(message.([]byte)))
			return true
		})
	})

	router.Run(":9977")
}
