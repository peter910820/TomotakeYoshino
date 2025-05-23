package shogi

import (
	"TomotakeYoshino/model"
	"bytes"
	"errors"
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
	pos, err := conversionToPosition(runes)
	if err != nil {
		logrus.Error(err)
		_, err = s.ChannelMessageSend(match.ChannleID, err.Error())
		if err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	piecePos, err := judgeMove(pos, string(runes[2]), match)
	if err != nil {
		logrus.Error(err)
		_, err = s.ChannelMessageSend(match.ChannleID, err.Error())
		if err != nil {
			logrus.Error(err)
			return
		}
		return
	}
	// 2.處理盤面座標
	refreshBoard(*pos, piecePos, string(runes[2]), match)
	// switch turn
	match.Turn = !match.Turn

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
		_, err = s.ChannelMessageSend(match.ChannleID, err.Error())
		if err != nil {
			logrus.Error(err)
			return
		}
		return
	}
}

// conversion command "shogimove" position parameter to model.Position struct
func conversionToPosition(position []rune) (*model.Position, error) {
	var pos model.Position
	var err error
	pos.X, err = strconv.Atoi(string(position[0]))
	if err != nil {
		return &pos, err
	}
	pos.Y, err = strconv.Atoi(string(position[1]))
	if err != nil {
		return &pos, err
	}
	return &pos, nil
}

// handle move
func judgeMove(pos *model.Position, pieceName string, match *model.Match) (model.Position, error) {
	var piecePos model.Position = model.Position{}
	var playerPieces map[string]model.Position
	matchPieces := []string{}

	if match.Turn {
		playerPieces = match.FirstPlayerPieces
	} else {
		playerPieces = match.SecondPlayerPieces
	}
	for k, v := range playerPieces {
		// exclude the situation where the target position has its own chess piece
		if v == *pos {
			return piecePos, errors.New("目標位置上有自己的棋子")
		}
		if strings.HasPrefix(k, CorrespondMap[pieceName]) {
			matchPieces = append(matchPieces, k) // pass the matched key itself into the slice
		}
	}

	var finallyMovePiece string = ""
	for _, v := range matchPieces {
		if piecesRules(v, playerPieces[v], *pos, match.Turn) {
			if finallyMovePiece != "" {
				return piecePos, errors.New("模稜兩可的操作")
			}
			piecePos = playerPieces[v]
			finallyMovePiece = v
		}
	}
	// eat pieces
	if match.Turn {
		match.FirstPlayerPieces[finallyMovePiece] = *pos
		for k, v := range match.SecondPlayerPieces {
			if v == *pos {
				match.FirstPlayerCapture = append(match.FirstPlayerCapture, k)
				delete(match.SecondPlayerPieces, k)
				break
			}
		}
	} else {
		match.SecondPlayerPieces[finallyMovePiece] = *pos
		for k, v := range match.FirstPlayerPieces {
			if v == *pos {
				match.FirstPlayerCapture = append(match.SecondPlayerCapture, k)
				delete(match.FirstPlayerPieces, k)
				break
			}
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
