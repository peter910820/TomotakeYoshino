package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/bot"
)

func init() {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	token := os.Getenv("TOKEN")
	yoshinoBot, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatal(err)
	}

	// add handler
	yoshinoBot.AddHandler(bot.Ready)
	yoshinoBot.AddHandler(bot.GuildMemberAdd)
	yoshinoBot.AddHandler(bot.OnInteraction)

	err = yoshinoBot.Open() // websocket connect
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("YoshinoBot is now running. Press CTRL+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	yoshinoBot.Close() // websocket disconnect
	logrus.Debug(interruptSignal)
}
