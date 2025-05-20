package cmds

import (
	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func ShogiStart(s *discordgo.Session, i *discordgo.InteractionCreate, shogi *map[string]*model.Match, opponent string) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		logrus.Error(err)
		return
	}

	userID, err := utils.GetUserID(i)
	if err != nil {
		logrus.Error(err)
		return
	}

	(*shogi)[channel.ID] = &model.Match{
		FirstPlayerID:  userID,
		SecondPlayerID: opponent,
	}
	_, err = s.ChannelMessageSend(channel.ID, "")
	if err != nil {
		logrus.Error(err)
		return
	}
}
