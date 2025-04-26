package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/bot"
)

var (
	yoshinoBot *discordgo.Session
	token      string
	appId      string
	err        error
)

func main() {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	err = godotenv.Load(".env")
	if err != nil {
		logrus.Fatal("[ERROR]: ", err)
	}
	token = os.Getenv("DISCORD_BOT_TOKEN")
	appId = os.Getenv("DISCORD_Application_ID")

	yoshinoBot, err = discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatal("[ERROR]: ", err)
	}

	yoshinoBot.AddHandler(ready)
	yoshinoBot.AddHandler(guildMemberAdd)
	yoshinoBot.AddHandler(onInteraction)

	err = yoshinoBot.Open() // websocket connect
	if err != nil {
		logrus.Fatal("[ERROR]: ", err)
	}

	logrus.Info("YoshinoBot is now running. Press CTRL+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	yoshinoBot.Close() // websocket disconnect
	logrus.Debug(interruptSignal)
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "クナド国記")
	bot.BasicCommand(s)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeMessage := "Welcome " + m.User.Username + "!"
	s.ChannelMessageSend(m.GuildID, welcomeMessage)
}

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logrus.Infof("InteractionCommand: %+v\n", i.ApplicationCommandData())
	switch i.ApplicationCommandData().Name {
	case "ping":
		delay := yoshinoBot.HeartbeatLatency()
		go bot.Ping(s, i, delay)
	case "guild":
		go bot.Guild(s, i)
	case "index":
		go bot.Index(s, i, appId)
	case "gnncrawler":
		go bot.GnnCrawler(s, i)
	}
}
