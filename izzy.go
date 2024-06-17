package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

// Variable initialization
var re = regexp.MustCompile(`(?:\A|\s|^)(cry|crying|sob|sobbing|😭)(?:\s|\z|$)`)

//var prefix = "!"

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		var systz string
		if runtime.GOOS == "android" {
			systz = "Asia/Kuala_Lumpur" //TODO: do not hardcode timezone info in bot
		} else {
			systz = "Local"
		}
		tz, _ := time.LoadLocation(systz)
		ostime := time.Now()
		mytime := ostime.In(tz).Format(time.UnixDate)
		var timemsg string = fmt.Sprintf("The time at my husband's place is currently %s!", mytime)
		s.ChannelMessageSend(m.ChannelID, timemsg)
	}

	if m.Content == "!boop" {
		s.ChannelMessageSend(m.ChannelID, "<a:IzzyCloseBoop:1237954801459794001>")
	}

	if m.Content == "!rizz" {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: "I will rizz you up my dear~", StickerIDs: []string{"1232131198231380049"}})
	}

	// TODO: Implement image grabbing command
	if strings.HasPrefix(m.Content, ">>>") {

		postID := strings.TrimPrefix(m.Content, ">>>")
		//derpiLink := fmt.Sprintf("https://derpibooru.org/api/v1/json/images/%s", postID)
		derpiLink := fmt.Sprintf("https://derpibooru.org/images/%s", postID)
		s.ChannelMessageSend(m.ChannelID, derpiLink)

		//reqUrl, err := http.Get(derpiLink)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//defer reqUrl.Body.Close()
		//
		//responsedata, err := io.ReadAll(reqUrl.Body)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//type Post struct {
		//	Image string `json:"image"`
		//	ID    string `json:"id"`
		//}
		//
		//var unmarshaled Post
		//
		//unmarshalErr := json.Unmarshal(responsedata, &unmarshaled)
		//if err != nil {
		//	log.Fatal(unmarshalErr)
		//}
		//
		//fmt.Printf(unmarshaled.ID)
		// derpiLink := fmt.Sprintf("https://derpibooru.org/images/%s", )
		// s.ChannelMessageSend(m.ChannelID, derpiLink)
	}

	if strings.Contains(m.Content, "hello") {
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: "Hi friend!!!", StickerIDs: []string{"1248585770902491187"}})
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(m.Content, "chad") {
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{StickerIDs: []string{"1246019517549645844"}})
		if err != nil {
			fmt.Println(err)
		}
	}

	if re.MatchString(m.Content) {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":IzzySob:1246619614285135952")
	}

	if strings.Contains(m.Content, "kek") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":IzzyKek:1239569282791247963")
	}

	if strings.Contains(m.Content, "good morning izzy") {
		s.ChannelMessageSend(m.ChannelID, "Good morning honey~")
	}
}

func messageEdit(s *discordgo.Session, e *discordgo.MessageUpdate) {
	var beforeEdit string
	var afterEdit string
	var authorID string

	afterEdit = e.Content

	if afterEdit == "" {
		return
	}

	if e.BeforeUpdate == nil {
		beforeEdit = "`Didn't catch that.`"
		afterEdit = "`Didn't catch that.`"
	} else {
		beforeEdit = e.BeforeUpdate.Content
		afterEdit = e.Content
	}

	if e.Author == nil {
		authorID = "`Didn't catch that.`"
	} else {
		authorID = e.Author.ID
	}

	editEmbed := fmt.Sprintf("Author: <@%s>\nIn: <#%s>\nEdited: %s\nBefore edit: %s", authorID, e.ChannelID, afterEdit, beforeEdit)
	s.ChannelMessageSendEmbed("1251510834832736300", &discordgo.MessageEmbed{Title: "Message edited!", Description: editEmbed})
}

func messageDelete(s *discordgo.Session, d *discordgo.MessageDelete) {
	var beforeDeleteAuthorID string
	var beforeDeleteContent string
	if d.BeforeDelete == nil {
		beforeDeleteAuthorID = "None"
		beforeDeleteContent = "None"
	} else {
		beforeDeleteAuthorID = d.BeforeDelete.Author.ID
		beforeDeleteContent = d.BeforeDelete.Content
	}
	deleteEmbed := fmt.Sprintf("User: <@%s>\nIn: <#%s>\nMessage: %s", beforeDeleteAuthorID, d.ChannelID, beforeDeleteContent)
	s.ChannelMessageSendEmbed("1251510834832736300", &discordgo.MessageEmbed{Title: "Deleted message!", Description: deleteEmbed, Color: 16711680})
}

func invHandler(s *discordgo.Session, i *discordgo.InviteCreate) {
	content := fmt.Sprintf("Created by: <@%s>\nDestination: <#%s>\nCode: %s", i.Inviter.ID, i.ChannelID, i.Code)
	s.ChannelMessageSendEmbed("1239034730855530639", &discordgo.MessageEmbed{Title: "Invite created!", Description: content})
}

func embedSend(s *discordgo.Session, m *discordgo.MessageCreate) { // Just ignore this lmao this is just an embed testcode
	if strings.Contains(m.Content, "!embed") {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Title: "Skibidi", Description: "I watch skibidi toilet at 3AM", URL: "https://google.com", Color: 2047602})
	}
}

func main() {
	var ostime = time.Now()
	var systz string

	if runtime.GOOS == "android" {
		systz = "Asia/Kuala_Lumpur"
	} else {
		systz = "Local"
	}

	tz, _ := time.LoadLocation(systz)

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

	izzy.AddHandler(messageHandler)
	izzy.AddHandler(invHandler)
	izzy.AddHandler(messageEdit)
	izzy.AddHandler(embedSend)
	izzy.AddHandler(messageDelete)
	izzy.StateEnabled = true
	izzy.State.TrackChannels = true
	izzy.State.MaxMessageCount = 100
	izzy.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	izzy.Open()
	izzy.UpdateGameStatus(0, "with Posi+ive!")

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Printf("Bot was started in %s\n", ostime.In(tz).Format(time.UnixDate))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	izzy.Close()
}
