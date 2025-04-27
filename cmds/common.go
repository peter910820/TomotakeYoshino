package cmds

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate, delay time.Duration) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Delay: %v", delay),
		},
	})
}

func Guild(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "guildID: " + i.GuildID,
		},
	})
}

func Index(s *discordgo.Session, i *discordgo.InteractionCreate, appId string) {
	commands, err := s.ApplicationCommands(appId, "")
	if err != nil {
		logrus.Error(err)
		return
	}
	resultData := "ALL commands:\n"
	for _, command := range commands {
		resultData += fmt.Sprintf("  * %s: %s\n", command.Name, command.Description)
		for _, option := range command.Options {
			resultData += "    Option: "
			resultData += fmt.Sprintf("%s: %s 必選: %t\n", option.Name, option.Description, option.Required)
		}
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintln("```", resultData, "```"),
		},
	})
}
