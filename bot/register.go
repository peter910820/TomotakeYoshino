package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func RegisterCommand(s *discordgo.Session) {
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
			Name:        "vndbsearchbrand",
			Description: "search galgame brand data for vndb",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "brand",
					Description: "搜尋的品牌名稱",
					Required:    true,
				},
			},
		},
		{
			Name:        "vndbsearch",
			Description: "use vndb barnd id search galgame data for vndb",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "brandid",
					Description: "搜尋的品牌ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "shogistart",
			Description: "開始一場將棋對弈",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "opponent",
					Description: "要對弈的對手ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "shogimove",
			Description: "操作將棋棋子",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "behavior",
					Description: "選擇進行移動棋子或是打入",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "移動",
							Value: "move",
						},
						{
							Name:  "打入",
							Value: "put",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "position",
					Description: "選擇棋子移動/打入位置",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "support",
					Description: "選擇輔助提示詞(開發中)",
					Required:    false,
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
