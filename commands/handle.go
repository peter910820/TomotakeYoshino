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

func GnnCrawler(s *discordgo.Session, i *discordgo.InteractionCreate, amount int64) {
	if amount > 20 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintln("資料筆數太大, 請輸入20以下的數字!"),
			},
		})
		return
	}
	resultData := "巴哈新聞:\n"
	c := colly.NewCollector()
	c.OnHTML(`h1.GN-lbox2D > a`, func(e *colly.HTMLElement) {
		if amount > 0 {
			resultData += fmt.Sprintf("* %s: https://%s\n", e.Text, e.Attr("href"))
			amount--
		}
	})
	err := c.Visit("https://gnn.gamer.com.tw/")
	if err != nil {
		log.Println("[ERROR]: ", err)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintln("```", resultData, "```"),
		},
	})
}
