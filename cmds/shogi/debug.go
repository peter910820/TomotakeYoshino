package shogi

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"TomotakeYoshino/model"
)

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
