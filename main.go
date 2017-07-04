// MessageBox project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	redisAddr = os.Getenv("REDISADDR")
	redisPs   = os.Getenv("REDISPS")
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
		log.Infoln(fmt.Sprintf("message: %d", i))
		time.Sleep(3 * 1e9)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	file, err := os.OpenFile("./logs/messagebox.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetLevel(log.WarnLevel)
}

func main() {
	hub, err := newHub()

	go hub.start()

	if err != nil {
		log.Fatalln(err)
		return
	}

	go testSend(hub)

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
