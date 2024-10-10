package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
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
		log.Println("[ERROR]: ", err)
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

func SearchGalgame(s *discordgo.Session, i *discordgo.InteractionCreate, galgame string) {
	c := colly.NewCollector()
	postingData := map[string]string{
		"mod":          "curforum",
		"formhash":     "",
		"srchtype":     "title",
		"srchfrom":     "0",
		"cid":          "",
		"srhfid":       "",
		"srhlocality":  "forum::forumdisplay",
		"srchtxt":      "",
		"searchsubmit": "",
	}
	c.OnHTML(`input[type="hidden"]`, func(e *colly.HTMLElement) {
		name := e.Attr("name")
		_, ok := postingData[name]
		if ok {
			fmt.Println(name)
		}
	})
	c.Visit("https://www.eyny.com/forum-3691-1.html")
}
