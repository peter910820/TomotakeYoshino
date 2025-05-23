package shogi

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"
)

// start a shogi match
func ShogiStart(s *discordgo.Session, i *discordgo.InteractionCreate, shogi *map[string]*model.Match, opponent string) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		logrus.Error(err)
		return
	}

	userID, err := utils.GetUserID(i)
	if err != nil {
		logrus.Error(err)
		return
	}

	userName, err := utils.GetUserName(s, userID)
	if err != nil {
		logrus.Error(err)
		return
	}
	opponentName, err := utils.GetUserName(s, opponent)
	if err != nil {
		logrus.Error(err)
		return
	}

	(*shogi)[channel.ID] = &model.Match{
		ChannleID:        channel.ID,
		FirstPlayerID:    userID,
		FirstPlayerName:  userName,
		SecondPlayerID:   opponent,
		SecondPlayerName: opponentName,
		Turn:             true,
	}

	initPlayerPieces((*shogi)[channel.ID])
	_, err = s.ChannelMessageSend(channel.ID, "棋子初始化成功，正在初始化盤面狀態")
	if err != nil {
		logrus.Error(err)
		return
	}

	initBoard((*shogi)[channel.ID])
	_, err = s.ChannelMessageSend(channel.ID, "盤面狀態初始化成功")
	if err != nil {
		logrus.Error(err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("```")
	for i := range 10 {
		for j := 10 - 1; j >= 0; j-- {
			buf.WriteString((*shogi)[channel.ID].Board[j][i])
		}
		buf.WriteString("\n")
	}
	buf.WriteString("```")

	_, err = s.ChannelMessageSend(channel.ID, buf.String())
	if err != nil {
		logrus.Error(err)
		return
	}
}

// init players pieces
func initPlayerPieces(match *model.Match) {
	// 這邊用的都是陣列座標 非實際看到的座標

	match.SecondPlayerPieces = map[string]model.Position{
		"oushou":    {X: 5, Y: 1},
		"kinshou1":  {X: 4, Y: 1},
		"kinshou2":  {X: 6, Y: 1},
		"ginshou1":  {X: 3, Y: 1},
		"ginshou2":  {X: 7, Y: 1},
		"keima1":    {X: 2, Y: 1},
		"keima2":    {X: 8, Y: 1},
		"kyousha1":  {X: 1, Y: 1},
		"kyousha2":  {X: 9, Y: 1},
		"kakugyou1": {X: 3, Y: 2},
		"Hisha1":    {X: 7, Y: 2},
	}
	match.FirstPlayerPieces = map[string]model.Position{
		"gyokushou": {X: 5, Y: 9},
		"kinshou1":  {X: 4, Y: 9},
		"kinshou2":  {X: 6, Y: 9},
		"ginshou1":  {X: 3, Y: 9},
		"ginshou2":  {X: 7, Y: 9},
		"keima1":    {X: 2, Y: 9},
		"keima2":    {X: 8, Y: 9},
		"kyousha1":  {X: 1, Y: 9},
		"kyousha2":  {X: 9, Y: 9},
		"kakugyou1": {X: 3, Y: 8},
		"Hisha1":    {X: 7, Y: 8},
	}

	for i := 1; i <= 9; i++ {
		match.SecondPlayerPieces[fmt.Sprintf("fuhyou%d", i)] = model.Position{X: i, Y: 3}
		match.FirstPlayerPieces[fmt.Sprintf("fuhyou%d", i)] = model.Position{X: i, Y: 7}
	}
}

// init board pieces
func initBoard(match *model.Match) {
	// 這邊用的都是陣列座標 非實際看到的座標

	for i := range match.Board {
		for j := range match.Board[i] {
			match.Board[i][j] = "＿"
		}
	}

	fullWidthDigits := []string{"１", "２", "３", "４", "５", "６", "７", "８", "９"}
	for i := 1; i <= 9; i++ {
		match.Board[i][0] = fullWidthDigits[i-1]
		match.Board[0][i] = fullWidthDigits[i-1]
	}

	pieces := []string{"香", "桂", "銀", "金", "王", "金", "銀", "桂", "香"}
	for i := 1; i <= 9; i++ {
		match.Board[i][1] = pieces[i-1]
		match.Board[i][9] = pieces[i-1]
	}
	pieces2 := []string{"＿", "角", "＿", "＿", "＿", "＿", "＿", "飛", "＿"}
	for i := 1; i <= 9; i++ {
		match.Board[i][2] = pieces2[i-1]
	}
	pieces3 := []string{"＿", "飛", "＿", "＿", "＿", "＿", "＿", "角", "＿"}
	for i := 1; i <= 9; i++ {
		match.Board[i][8] = pieces3[i-1]
	}
	pieces4 := []string{"步", "步", "步", "步", "步", "步", "步", "步", "步"}
	for i := 1; i <= 9; i++ {
		match.Board[i][3] = pieces4[i-1]
		match.Board[i][7] = pieces4[i-1]
	}
}
