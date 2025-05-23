package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/cmds"
	"TomotakeYoshino/cmds/shogi"
	"TomotakeYoshino/utils"
)

var (
	AppID string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal(err)
	}
	AppID = os.Getenv("APP_ID")
	if AppID == "" {
		logrus.Fatal("appId not set")
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "サクラノ詩-櫻の森の上を舞う-")
	RegisterCommand(s)
}

func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeMessage := "Welcome " + m.User.Username + "!"
	s.ChannelMessageSend(m.GuildID, welcomeMessage)
}

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logrus.Infof("InteractionCommand: %+v\n", i.ApplicationCommandData())
	switch i.ApplicationCommandData().Name {
	case "ping":
		go cmds.Ping(s, i, s.HeartbeatLatency())
	case "guild":
		go cmds.Guild(s, i)
	case "index":
		go cmds.Index(s, i, AppID)
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
		opponentID, err := utils.GetOptions(i, "opponentid")
		if err != nil {
			logrus.Error(err)
			return
		}
		// check if channel has other match
		channel, err := s.Channel(i.ChannelID)
		if err != nil {
			logrus.Error(err)
			return
		}
		_, ok := shogi.ShogiMatch[channel.ID]
		if ok {
			logrus.Error("該頻道已有未結束的對局")
			utils.SlashCommandRespond(s, i, "該頻道已有未結束的對局")
			return
		}
		// start a shogi match
		go shogi.ShogiStart(s, i, &shogi.ShogiMatch, opponentID)
		utils.SlashCommandRespond(s, i, "正在開始創建對局，請稍後")
	case "shogimove":
		behavior, err := utils.GetOptions(i, "behavior")
		if err != nil {
			logrus.Error(err)
			utils.SlashCommandError(s, i, err.Error())
			return
		}
		position, err := utils.GetOptions(i, "position")
		if err != nil {
			logrus.Error(err)
			utils.SlashCommandError(s, i, err.Error())
			return
		}
		// check if channel has other match
		channel, err := s.Channel(i.ChannelID)
		if err != nil {
			logrus.Error(err)
			utils.SlashCommandError(s, i, err.Error())
			return
		}
		utils.SlashCommandRespond(s, i, "正在移動棋子")
		if behavior == "move" {
			go shogi.ShogiMove(s, i, shogi.ShogiMatch[channel.ID], position)
		} else {

		}
	case "shogidebug":
		channelID, err := utils.GetOptions(i, "channleid")
		if err != nil {
			logrus.Error(err)
			utils.SlashCommandError(s, i, err.Error())
			return
		}
		_, ok := shogi.ShogiMatch[channelID]
		if ok {
			userID, err := utils.GetUserID(i)
			if err != nil {
				utils.SlashCommandRespond(s, i, "找不到使用者") // 基本上應該不會發生這種狀況
				return
			}
			// check if is user turn
			if shogi.ShogiMatch[channelID].Turn {
				if shogi.ShogiMatch[channelID].FirstPlayerID != userID {
					utils.SlashCommandRespond(s, i, "現在不是你的回合")
					return
				}
			} else {
				if shogi.ShogiMatch[channelID].SecondPlayerID != userID {
					utils.SlashCommandRespond(s, i, "現在不是你的回合")
					return
				}
			}
			go shogi.GetShogiPiecesData(s, i, shogi.ShogiMatch[channelID])
		} else {
			utils.SlashCommandRespond(s, i, "該對局不存在")
		}
		// value, err := utils.GetOptions(i, "support")
		// if err != nil {
		// 	logrus.Error(err)
		// 	return
		// }
	}
}
