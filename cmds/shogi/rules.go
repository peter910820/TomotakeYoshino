package shogi

import (
	"TomotakeYoshino/model"
	"TomotakeYoshino/utils"
)

func keimaRule(piecePos model.Position, targetPos model.Position, turn bool) bool {
	if turn {
		if (utils.Abs(targetPos.X-piecePos.X) == 1) && (piecePos.Y-targetPos.Y == 2) {
			return true
		} else {
			return false
		}
	} else {
		if (utils.Abs(targetPos.X-piecePos.X) == 1) && (targetPos.Y-piecePos.Y == 2) {
			return true
		} else {
			return false
		}
	}
}
