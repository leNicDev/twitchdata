package main

import (
	"sync"
	"twitchdata/chat"
	"twitchdata/db"
	"twitchdata/twitch"
)

func main() {
	db.Connect()
	defer db.Disconnect()

	// get top channels
	channels, err := twitch.GetTopStreamNames(100)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < len(channels); i++ {
		wg.Add(i)
		go chat.Connect(channels[i], &wg)
	}

	// wait for all connected chats
	wg.Wait()
}
