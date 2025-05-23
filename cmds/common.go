package cmds

import (
	"TomotakeYoshino/utils"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate, delay time.Duration) {
	utils.SlashCommandRespond(s, i, fmt.Sprintf("Delay: %v", delay))
}

func Guild(s *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.SlashCommandRespond(s, i, "guildID: "+i.GuildID)
}

func Index(s *discordgo.Session, i *discordgo.InteractionCreate, appID string) {
	commands, err := s.ApplicationCommands(appID, "")
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
	utils.SlashCommandRespond(s, i, fmt.Sprintln("```", resultData, "```"))
}
