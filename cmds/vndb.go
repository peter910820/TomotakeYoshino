package cmds

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"
)

func VndbSearchVn(s *discordgo.Session, i *discordgo.InteractionCreate, brandId string) {
	data := utils.VndbVnRequest(brandId)

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	resp, body, err := utils.JsonRequest("https://api.vndb.org/kana/vn", "POST", jsonData)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("the server returned an error status code %d", resp.StatusCode)
		utils.SlashCommandError(s, i, fmt.Sprintf("the server returned an error status code %d", resp.StatusCode))
		return
	}

	logrus.Infof("狀態碼: %d", resp.StatusCode)
	logrus.Infof("回傳內容: \n%s", string(body))

	displayData, err := vndbResponseHandle(&body)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Vndb遊戲查詢結果",
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "遊戲名稱/Vndb評分",
				Value:  displayData,
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

func VndbSearchProducer(s *discordgo.Session, i *discordgo.InteractionCreate, brand string) {
	data := utils.VndbProducerRequest(brand)

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	resp, body, err := utils.JsonRequest("https://api.vndb.org/kana/producer", "POST", jsonData)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("the server returned an error status code %d", resp.StatusCode)
		utils.SlashCommandError(s, i, fmt.Sprintf("the server returned an error status code %d", resp.StatusCode))
		return
	}

	// var result map[string]interface{}
	// json.Unmarshal(body, &result)

	logrus.Infof("狀態碼: %d", resp.StatusCode)
	logrus.Infof("回傳內容: \n%s", string(body))

	embed := &discordgo.MessageEmbed{
		Title: "Vndb品牌查詢結果",
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "訊息內容",
				Value:  string(body),
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

func vndbResponseHandle(body *[]byte) (string, error) {
	var result model.VndbVnResponse
	var displayData string = ""

	err := json.Unmarshal(*body, &result)
	if err != nil {
		return displayData, err
	}

	for index, value := range result.Results {
		tmpStr := fmt.Sprintf("%d. ", index+1)
		if value.AltTitle == nil {
			tmpStr += fmt.Sprint(value.Title + ": ")
		} else {
			tmpStr += fmt.Sprint(*value.AltTitle + ": ")
		}
		tmpStr += fmt.Sprintf("%f\n", math.Round(value.Rating*10)/10)

		displayData += tmpStr
	}
	if result.More {
		displayData += "還有更多筆資料..."
	}
	return displayData, nil
}
