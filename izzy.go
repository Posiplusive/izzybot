package main

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var ostime = time.Now()

func waifuHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.State.Channel(m.ChannelID)
	if ch == nil {
		ch, _ = s.Channel(m.ChannelID)
	}

	sv, _ := s.State.Guild(m.GuildID)
	if sv == nil {
		sv, _ = s.Guild(m.GuildID)
	}

	fmt.Printf("Message in %s from @%s in #%+v: %s\n", sv.Name, m.Author.Username, ch.Name, m.Content)

	file, _ := os.ReadFile("./waifu.yaml")

	if m.Content == "!waifu" { // Honestly this part is just fucking hell

		type UserData struct { // Declares the struct for the yaml
			UserId []string `yaml:"userid"`
			Quote  []string `yaml:"quote"`
		}

		var waifu UserData

		err := yaml.Unmarshal(file, &waifu) // Turns the yaml to a slice and put it in waifu
		if err != nil {
			fmt.Print(err)
		}

		id := m.Author.ID
		x := slices.Index(waifu.UserId, id) // Returns the index of the matching user ID
		var y string                        // Declare y var

		if x == -1 { // Puts a random string here so that the
			y = "undefined" // switch case below goes to the default part
		} else {
			y = waifu.UserId[x] // Puts the user ID in a variable for comparison
		}

		switch id {
		case y: // if y matches the user,
			s.ChannelMessageSend(m.ChannelID, waifu.Quote[x]) // then send a message to the user
		default: // Otherwise just tell them they're not in the list
			s.ChannelMessageSend(m.ChannelID, "You don't have a waifu yet! Ask my husband to add her for you.")
		}
	}

	if m.Content == "!time" {
		ostime := time.Now()
		mytime := ostime.Format(time.UnixDate)
		var timemsg string = fmt.Sprintf("The time at my husband's place is currently %s!", mytime)
		s.ChannelMessageSend(m.ChannelID, timemsg)
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

	izzy.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	izzy.Open()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Printf("Bot was started in %s\n", ostime.Format(time.UnixDate))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	izzy.Close()
}
