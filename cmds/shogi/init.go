package shogi

import (
	"fmt"

	"TomotakeYoshino/model"
)

// init players pieces
func initPlayerPieces(match *model.Match) {
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
func initBoard(match *model.Match) {
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
