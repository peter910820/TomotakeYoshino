package shogi

import (
	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	CorrespondMap map[string]string = map[string]string{"桂": "keima"}
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

	InitPlayerPieces((*shogi)[channel.ID])
	_, err = s.ChannelMessageSend(channel.ID, "棋子初始化成功，正在初始化盤面狀態")
	if err != nil {
		logrus.Error(err)
		return
	}

	InitBoard((*shogi)[channel.ID])
	_, err = s.ChannelMessageSend(channel.ID, "盤面狀態初始化成功")
	if err != nil {
		logrus.Error(err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("```")
	for i := 0; i < 10; i++ {
		for j := 10 - 1; j >= 0; j-- {
			buf.WriteString((*shogi)[channel.ID].Board[i][j])
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

func InitPlayerPieces(match *model.Match) {
	match.SecondPlayerPieces = map[string]model.Position{
		"oushou":    {X: 1, Y: 5},
		"kinshou1":  {X: 1, Y: 4},
		"kinshou2":  {X: 1, Y: 6},
		"ginshou1":  {X: 1, Y: 3},
		"ginshou2":  {X: 1, Y: 7},
		"keima1":    {X: 1, Y: 2},
		"keima2":    {X: 1, Y: 8},
		"kyousha1":  {X: 1, Y: 1},
		"kyousha2":  {X: 1, Y: 9},
		"kakugyou1": {X: 2, Y: 3},
		"Hisha1":    {X: 2, Y: 7},
	}
	match.FirstPlayerPieces = map[string]model.Position{
		"gyokushou": {X: 9, Y: 5},
		"kinshou1":  {X: 9, Y: 4},
		"kinshou2":  {X: 9, Y: 6},
		"ginshou1":  {X: 9, Y: 3},
		"ginshou2":  {X: 9, Y: 7},
		"keima1":    {X: 9, Y: 2},
		"keima2":    {X: 9, Y: 8},
		"kyousha1":  {X: 9, Y: 1},
		"kyousha2":  {X: 9, Y: 9},
		"kakugyou1": {X: 8, Y: 3},
		"Hisha1":    {X: 8, Y: 7},
	}

	for i := 1; i <= 9; i++ {
		match.SecondPlayerPieces[fmt.Sprintf("fuhyou%d", i)] = model.Position{X: 3, Y: i}
		match.FirstPlayerPieces[fmt.Sprintf("fuhyou%d", i)] = model.Position{X: 7, Y: i}
	}
}

// init board pieces
func InitBoard(match *model.Match) {
	for i := range match.Board {
		for j := range match.Board[i] {
			match.Board[i][j] = "＿"
		}
	}

	fullWidthDigits := []string{"１", "２", "３", "４", "５", "６", "７", "８", "９"}
	for i := 1; i <= 9; i++ {
		match.Board[0][i] = fullWidthDigits[i-1]
		match.Board[i][0] = fullWidthDigits[i-1]
	}

	pieces := []string{"香", "桂", "銀", "金", "王", "金", "銀", "桂", "香"}
	for i := 1; i <= 9; i++ {
		match.Board[1][i] = pieces[i-1]
		match.Board[9][i] = pieces[i-1]
	}
	pieces2 := []string{"＿", "飛", "＿", "＿", "＿", "＿", "＿", "角", "＿"}
	for i := 1; i <= 9; i++ {
		match.Board[2][i] = pieces2[i-1]
	}
	pieces3 := []string{"＿", "角", "＿", "＿", "＿", "＿", "＿", "飛", "＿"}
	for i := 1; i <= 9; i++ {
		match.Board[8][i] = pieces3[i-1]
	}
	pieces4 := []string{"步", "步", "步", "步", "步", "步", "步", "步", "步"}
	for i := 1; i <= 9; i++ {
		match.Board[3][i] = pieces4[i-1]
		match.Board[7][i] = pieces4[i-1]
	}
}

// 輔助字有兩種 一種是標示位置的輔助字(上下左右) 一種是升變(成)
func ShogiMove(s *discordgo.Session, i *discordgo.InteractionCreate, match *model.Match, position string) {
	// support word
	_ = ""
	if len(position) > 3 {
		// 有輔助字的狀況
		_ = position[3:len(position)]
	}

	// 1.處理棋子座標

	// 先處理沒有輔助字的狀況
	// 轉換座標資料格式
	pos, err := translateToPosition(position)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	piecePos, err := judgeMove(pos, position[2:3], match)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}
	// 2.處理盤面座標
	refreshBoard(*pos, piecePos, position[2:3], match)

	var buf bytes.Buffer
	buf.WriteString("```")
	for i := range 10 {
		for j := 10 - 1; j >= 0; j-- {
			buf.WriteString(match.Board[i][j])
		}
		buf.WriteString("\n")
	}
	buf.WriteString("```")

	_, err = s.ChannelMessageSend(match.ChannleID, buf.String())
	if err != nil {
		logrus.Error(err)
		return
	}
}

// translate command "shogimove" position parameter to model.Position struct
func translateToPosition(position string) (*model.Position, error) {
	var pos model.Position
	var err error
	pos.X, err = strconv.Atoi(position[:1])
	if err != nil {
		return &pos, err
	}
	pos.Y, err = strconv.Atoi(position[1:2])
	if err != nil {
		return &pos, err
	}
	return &pos, nil
}

// handle move
func judgeMove(pos *model.Position, pieceName string, match *model.Match) (model.Position, error) {
	var piecePos model.Position
	matchPieces := []string{}
	if match.Turn {
		for k, v := range match.FirstPlayerPieces {
			// 排除目標位置有自己的棋子的狀況
			if v == *pos {
				return piecePos, errors.New("目標位置上有自己的棋子")
			}
			if strings.HasPrefix(k, CorrespondMap[pieceName]) {
				matchPieces = append(matchPieces, k) // 將匹配到的鍵本身傳入切片中
			}
		}
	} else {
		for k, v := range match.SecondPlayerPieces {
			// 排除目標位置有自己的棋子的狀況
			if v == *pos {
				return piecePos, errors.New("目標位置上有自己的棋子")
			}
			if strings.HasPrefix(k, CorrespondMap[pieceName]) {
				matchPieces = append(matchPieces, k) // 將匹配到的鍵本身傳入切片中
			}
		}
	}
	// 目前只有先手的狀況 因為邏輯較為複雜所以之後這段會重構 要代碼複用的代碼複用 要省略的省略
	var finallyMovePiece string = ""
	for _, v := range matchPieces {
		judgeFunc := piecesRules(v)
		if judgeFunc(match.FirstPlayerPieces[v], *pos, match.Turn) {
			if finallyMovePiece != "" {
				return piecePos, errors.New("模稜兩可的操作")
			}
			piecePos = match.FirstPlayerPieces[v]
			finallyMovePiece = v
		}
	}
	match.FirstPlayerPieces[finallyMovePiece] = *pos

	// 這邊還要去撈目標位置是否有敵方棋子，有的話刪除他並加到自己的capture陣列中
	for k, v := range match.SecondPlayerPieces {
		if v == *pos {
			match.FirstPlayerCapture = append(match.FirstPlayerCapture, k)
			delete(match.SecondPlayerPieces, k)
			break
		}
	}

	return piecePos, nil
}

func piecesRules(pieceName string) func(model.Position, model.Position, bool) bool {
	switch pieceName {
	case "keima":
		return keimaRule
	}
	return keimaRule
}

// refresh match board status
func refreshBoard(piecePos model.Position, targetPos model.Position, pieceName string, match *model.Match) {
	match.Board[piecePos.X][piecePos.Y] = "＿"
	match.Board[targetPos.X][targetPos.Y] = pieceName
}

// get the shogi pieces data for test
func GetShogiPiecesData(s *discordgo.Session, i *discordgo.InteractionCreate, match *model.Match) {
	var buf bytes.Buffer
	buf.WriteString("```")
	buf.WriteString("First Player Pieces:\n")
	for k, v := range match.FirstPlayerPieces {
		buf.WriteString(fmt.Sprintf("  %s: {%d %d}\n", k, v.X, v.Y))
	}
	buf.WriteString("Second Player Pieces:\n")
	for k, v := range match.SecondPlayerPieces {
		buf.WriteString(fmt.Sprintf("  %s: {%d %d}\n", k, v.X, v.Y))
	}
	buf.WriteString("```")

	_, err := s.ChannelMessageSend(match.ChannleID, buf.String())
	if err != nil {
		logrus.Error(err)
		return
	}
}

// 以後只會傳match的指標，將選擇match的判斷留給進goruntine之前
