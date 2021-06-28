package server

import (
	"sync"
	"time"

	"github.com/tapvanvn/go-wspubsub/utility"
)

var __mux sync.Mutex
var __msgmap map[string]*Message = make(map[string]*Message)
var __timemap map[string]int64 = map[string]int64{}
var __tier1_checking bool = false

var __tier2_checking bool = false
var __tier2mux sync.Mutex
var __tier2map map[*Message]*Client = make(map[*Message]*Client)

type Observer struct {
	client  *Client
	message *Message
}

func responseTier1(code string) {
	__mux.Lock()
	delete(__timemap, code)
	delete(__msgmap, code)
	__mux.Unlock()
}

func observeTie1Check() {
	__mux.Lock()
	remain := map[string]int64{}
	now := time.Now().Unix()
	for code, lastTime := range __timemap {
		if now > lastTime {
			msg, _ := __msgmap[code]
			topic := GetTopic(msg.Topic)
			go func() { topic.pick <- msg }()
			delete(__msgmap, code)
		} else {
			remain[code] = lastTime
		}
	}
	__timemap = remain
	__mux.Unlock()
}

func ObserveTier1Check() {

	if __tier1_checking {
		return
	}
	__tier1_checking = true
	go utility.Schedule(observeTie1Check, time.Second)
}

func ObserveTier1(message *Message) {
	__mux.Lock()
	code := utility.GenCode(5)
	message.Attributes["raycode"] = code
	__msgmap[code] = message
	__timemap[code] = time.Now().Unix() + 2
	__mux.Unlock()
}

func observerTier2Check() {
	__tier2mux.Lock()
	for msg, client := range __tier2map {
		if !client.live {
			delete(__tier2map, msg)
			topic := GetTopic(msg.Topic)
			go func() { topic.pick <- msg }()
			break
		}
	}
	__tier2mux.Unlock()
}

func ObserverTier2Check() {
	if __tier2_checking {
		return
	}
	__tier2_checking = true
	utility.Schedule(observerTier2Check, time.Second*2)
}

func ObserveTier2(message *Message, client *Client) {
	__tier2mux.Lock()
	__tier2map[message] = client
	__tier2mux.Unlock()
}
