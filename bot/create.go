package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func BasicCommand(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "return bot heartbeatlatency",
		},
		{
			Name:        "guild",
			Description: "get guild id",
		},
		{
			Name:        "index",
			Description: "get all command",
		},
		{
			Name:        "gnncrawler",
			Description: "get gnn news",
		},
		{
			Name:        "vndbsearch",
			Description: "search galgame data for vndb",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "brand",
					Description: "搜尋的品牌名稱",
					Required:    true,
				},
			},
		},
	}
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			logrus.Errorf("Register slash command faild: %s", err)
			return
		}
	}
}
