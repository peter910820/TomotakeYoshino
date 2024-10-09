package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")

	yoshinoBot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	yoshinoBot.AddHandler(guildMemberAdd)

	err = yoshinoBot.Open() // websocket connect
	if err != nil {
		log.Fatal(err)
	}

	log.Println("YoshinoBot is now running. Press CTRL+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	fmt.Println(interruptSignal)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeMessage := "Welcome " + m.User.Username + "!"
	s.ChannelMessageSend(m.GuildID, welcomeMessage)
}
