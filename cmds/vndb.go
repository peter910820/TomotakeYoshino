package cmds

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/utils"
)

func VndbSearch(s *discordgo.Session, i *discordgo.InteractionCreate, brand string) {
	data := utils.VndbRequestData(brand)

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

	logrus.Println("狀態碼：", resp.StatusCode)
	logrus.Println("回傳內容：", string(body))
}
