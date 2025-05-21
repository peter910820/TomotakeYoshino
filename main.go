package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/bot"
	"TomotakeYoshino/cmds"
	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"
)

var (
	appId string

	shogi map[string]*model.Match = make(map[string]*model.Match)
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal(err)
	}
	appId = os.Getenv("APP_ID")
	if appId == "" {
		logrus.Fatal("appId not set")
	}

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
	yoshinoBot.AddHandler(onInteraction)

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

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logrus.Infof("InteractionCommand: %+v\n", i.ApplicationCommandData())
	switch i.ApplicationCommandData().Name {
	case "ping":
		go cmds.Ping(s, i, s.HeartbeatLatency())
	case "guild":
		go cmds.Guild(s, i)
	case "index":
		go cmds.Index(s, i, appId)
	case "gnncrawler":
		go cmds.GnnCrawler(s, i)
	case "vndbsearchbrand":
		value, err := utils.GetOptions(i, "brand")
		if err != nil {
			logrus.Error(err)
			return
		}
		go cmds.VndbSearchProducer(s, i, value)
	case "vndbsearch":
		value, err := utils.GetOptions(i, "brandid")
		if err != nil {
			logrus.Error(err)
			return
		}
		go cmds.VndbSearchVn(s, i, value)
	case "shogistart":
		value, err := utils.GetOptions(i, "opponent")
		if err != nil {
			logrus.Error(err)
			return
		}

		channel, err := s.Channel(i.ChannelID)
		if err != nil {
			logrus.Error(err)
			return
		}
		_, ok := shogi[channel.ID]
		if ok {
			logrus.Error("該頻道已有未結束的對局")
			utils.SlashCommandRespond(s, i, "該頻道已有未結束的對局")
			return
		}
		go cmds.ShogiStart(s, i, &shogi, value)
		utils.SlashCommandRespond(s, i, "正在開始創建對局，請稍後")
	}
}
