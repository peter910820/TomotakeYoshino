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
			Name:        "gnncrawler",
			Description: "get gnn news",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "a number of queries, less than 20",
					Required:    false,
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
