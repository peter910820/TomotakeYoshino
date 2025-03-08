package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
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
		logrus.Error("[ERROR]: ", err)
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

func GnnCrawler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	amount := 10
	resultData := ""
	c := colly.NewCollector()
	c.OnHTML(`h1.GN-lbox2D > a`, func(e *colly.HTMLElement) {
		if amount > 0 {
			resultData += fmt.Sprintf("* %s: https://%s\n\n", e.Text, e.Attr("href"))
			amount--
		}
	})
	err := c.Visit("https://gnn.gamer.com.tw/")
	if err != nil {
		logrus.Error("[ERROR]: ", err)
	}
	resultData = strings.TrimSuffix(resultData, "\n")
	resultData = strings.TrimSuffix(resultData, "\n")
	embed := &discordgo.MessageEmbed{
		Title: "巴哈姆特最新新聞列表",
		Color: 0x51a1b4,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "訊息內容",
				Value:  resultData,
				Inline: false,
			},
		},
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
