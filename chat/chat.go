package chat

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"sync"
	"twitchdata/db"
)

func Connect(channel string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := twitch.NewAnonymousClient()

	client.OnConnect(func() {
		log.Printf("Connected to %s", channel)
	})
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msg := db.ChatMessage{
			UserID:          message.User.ID,
			UserName:        message.User.Name,
			UserDisplayName: message.User.DisplayName,
			Message:         message.Message,
			Channel:         message.Channel,
			Time:            message.Time,
			Bits:            message.Bits,
		}
		err := db.InsertMessage(msg)
		if err != nil {
			log.Fatal(err)
		}

		// log message
		//log.Printf("%+v", msg)
	})

	client.Join(channel)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
