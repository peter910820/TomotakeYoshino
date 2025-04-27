package cmds

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

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
