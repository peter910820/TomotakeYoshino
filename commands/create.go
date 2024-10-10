package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
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
			Name:        "searchgalgame",
			Description: "在eyny-GalGame遊戲下載區(上傳空間)搜尋該galgame有無資料",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "galgame",
					Description: "a galgame to search",
					Required:    true,
				},
			},
		},
	}
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Println("[ERROR]: ", err)
			log.Println("Register slash command faild!")
			return
		}
	}
}
