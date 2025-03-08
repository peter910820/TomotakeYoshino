package main

import (
	"TomotakeYoshino/commands"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	yoshinoBot *discordgo.Session
	token      string
	appId      string
	err        error
)

func main() {
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("[ERROR]: ", err)
	}
	token = os.Getenv("DISCORD_BOT_TOKEN")
	appId = os.Getenv("DISCORD_Application_ID")

	yoshinoBot, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("[ERROR]: ", err)
	}

	yoshinoBot.AddHandler(ready)
	yoshinoBot.AddHandler(guildMemberAdd)
	yoshinoBot.AddHandler(onInteraction)

	err = yoshinoBot.Open() // websocket connect
	if err != nil {
		log.Fatal("[ERROR]: ", err)
	}

	log.Println("YoshinoBot is now running. Press CTRL+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	yoshinoBot.Close() // websocket disconnect
	log.Println(interruptSignal)
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "クナド国記")
	commands.BasicCommand(s)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeMessage := "Welcome " + m.User.Username + "!"
	s.ChannelMessageSend(m.GuildID, welcomeMessage)
}

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("InteractionCommand: %+v\n", i.ApplicationCommandData())
	switch i.ApplicationCommandData().Name {
	case "ping":
		delay := yoshinoBot.HeartbeatLatency()
		go commands.Ping(s, i, delay)
	case "guild":
		go commands.Guild(s, i)
	case "index":
		go commands.Index(s, i, appId)
	case "gnncrawler":
		go commands.GnnCrawler(s, i)
	}
}
