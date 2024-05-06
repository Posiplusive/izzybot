package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func waifuHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.State.Channel(m.ChannelID)
	if ch == nil {
		ch, _ = s.Channel(m.ChannelID)
	}

	fmt.Printf("Message from @%s in #%+v: %s\n", m.Author.Username, ch.Name, m.Content)

	if m.Content == "!waifu" {
		s.ChannelMessageSend(m.ChannelID, "Me! (Posi+ive is simping hard for me...)")
	}
}

func main() {
	token, err := os.ReadFile("./token.txt")
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	izzy, err := discordgo.New("Bot " + string(token))

	if err != nil {
		x := fmt.Sprintf("I can't create Izzy! (%d)", err)
		fmt.Println(x)
	}

	izzy.AddHandler(waifuHandler)

	izzy.Identify.Intents = discordgo.IntentsGuildMessages

	izzy.Open()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	izzy.Close()
}
