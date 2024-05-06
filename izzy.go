package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!waifu" {
		s.ChannelMessageSend(m.ChannelID, "Me! (Posi+ive is simping hard for me...)")
	}

	// msg := fmt.Sprintf("Message from @%s in #%s: %s", m.Author.Username, discordgo.Channel.Name, m.Content)
	// fmt.Println(msg)
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

	izzy.AddHandler(messageCreate)

	izzy.Identify.Intents = discordgo.IntentsGuildMessages

	izzy.Open()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	izzy.Close()
}
