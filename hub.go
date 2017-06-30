package main

import (
	"log"

	"github.com/dustin/go-broadcast"
	"github.com/garyburd/redigo/redis"
)

type Hub struct {
	broaddcaster broadcast.Broadcaster

	//utilize redis pubsub to listen event from remote server
	psc redis.PubSubConn
}

func newHub() (*Hub, error) {

	prc, err := redis.Dial("tcp", redisAddr)

	if err != nil {
		log.Println("Connect to redis error", err)
		return nil, err
	}

	psc := redis.PubSubConn{Conn: prc}

	psc.Subscribe("messagebox:sse")

	hub := &Hub{
		broaddcaster: broadcast.NewBroadcaster(1024),
		psc:          psc,
	}
	//go hub.start()

	return hub, nil
}

func (hub *Hub) start() {
	for {
		switch v := hub.psc.Receive().(type) {
		case redis.Message:
			log.Println("received: " + string(v.Data))
			hub.broaddcaster.Submit(v.Data)
		case redis.Subscription:
			break
		case error:
			log.Println(v)
			return
		}
	}
}

func (h *Hub) registerListener() chan interface{} {
	listener := make(chan interface{})
	h.broaddcaster.Register(listener)
	return listener
}

func (h *Hub) closeListener(listener chan interface{}) {
	h.broaddcaster.Unregister(listener)
	close(listener)
}
