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
	ShogiMatch    map[string]*model.Match = make(map[string]*model.Match)
	CorrespondMap map[string]string       = map[string]string{
		"步": "fuhyou",
		"桂": "keima",
	}
)

// 輔助字有兩種 一種是標示位置的輔助字(上下左右) 一種是升變(成)
func ShogiMove(s *discordgo.Session, i *discordgo.InteractionCreate, match *model.Match, position string) {
	// support word
	_ = ""
	if len(position) > 3 {
		// 有輔助字的狀況
		_ = position[3:len(position)]
	}

	// 1.處理棋子座標

	runes := []rune(position)
	// 先處理沒有輔助字的狀況
	// 轉換座標資料格式
	pos, err := conversionToPosition(position)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}

	fmt.Println(string(runes[2]))

	piecePos, err := judgeMove(pos, string(runes[2]), match)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandError(s, i, err.Error())
		return
	}
	// 2.處理盤面座標
	refreshBoard(*pos, piecePos, string(runes[2]), match)

	var buf bytes.Buffer
	buf.WriteString("```")
	for i := range 10 {
		for j := 10 - 1; j >= 0; j-- {
			buf.WriteString(match.Board[j][i])
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

// conversion command "shogimove" position parameter to model.Position struct
func conversionToPosition(position string) (*model.Position, error) {
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
	var piecePos model.Position = model.Position{}
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
		if piecesRules(v, match.FirstPlayerPieces[v], *pos, match.Turn) {
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

func piecesRules(pieceName string, piecePos model.Position, targetPos model.Position, turn bool) bool {
	switch {
	case strings.HasPrefix(pieceName, "fuhyou"):
		return fuhyouRule(piecePos, targetPos, turn)
	case strings.HasPrefix(pieceName, "keima"):
		return keimaRule(piecePos, targetPos, turn)
	default:
		return false
	}
}

// refresh match board status
func refreshBoard(targetPos model.Position, piecePos model.Position, pieceName string, match *model.Match) {
	match.Board[piecePos.X][piecePos.Y] = "＿"
	match.Board[targetPos.X][targetPos.Y] = pieceName
}

// 以後只會傳match的指標，將選擇match的判斷留給進goruntine之前
