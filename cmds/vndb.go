package cmds

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/utils"
)

func VndbSearch(s *discordgo.Session, i *discordgo.InteractionCreate, brand string) {
	// id, err := getProducer(brand)
	// if err != nil {
	// 	logrus.Error(err)
	// 	utils.SlashCommandError(s, i, err.Error())
	// 	return
	// }
	// logrus.Debug(id)
	// data := utils.VndbRequestData(brand)

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	logrus.Error(err)
	// 	utils.SlashCommandError(s, i, err.Error())
	// 	return
	// }

	// resp, body, err := utils.JsonRequest("https://api.vndb.org/kana/vn", "POST", jsonData)
	// if err != nil {
	// 	logrus.Error(err)
	// 	utils.SlashCommandError(s, i, err.Error())
	// 	return
	// }

	// if resp.StatusCode != 200 {
	// 	logrus.Errorf("the server returned an error status code %d", resp.StatusCode)
	// 	utils.SlashCommandError(s, i, fmt.Sprintf("the server returned an error status code %d", resp.StatusCode))
	// 	return
	// }

	// logrus.Infof("狀態碼: %d", resp.StatusCode)
	// logrus.Infof("回傳內容: \n%s", string(body))

	// vndbResponse(body)

	// embed := &discordgo.MessageEmbed{
	// 	Title: "Vndb查詢結果",
	// 	Color: 0x04108e,
	// 	Fields: []*discordgo.MessageEmbedField{
	// 		{
	// 			Name:   "訊息內容",
	// 			Value:  string(body),
	// 			Inline: false,
	// 		},
	// 	},
	// }
	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Embeds: []*discordgo.MessageEmbed{embed},
	// 	},
	// })
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

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	logrus.Debug(result["results"])

	embed := &discordgo.MessageEmbed{
		Title: "Vndb查詢結果",
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

func vndbResponse(body []byte) {
	var result map[string]interface{}

	err := json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
}
