package model

// 用map[string]Match 結構來代表多筆對局

// 每一場對局
type Match struct {
	FirstPlayerID      string
	FirstPlayerName    string
	SecondPlayerID     string
	SecondPlayerName   string
	FirstPlayerPieces  map[string]Position
	SecondPlayerPieces map[string]Position
	TurnID             bool
	Board              [][]string // 詳細描述盤面狀態，每次更新時只需更新移動的棋子的目標跟原位子兩個點就好
}

// 每一顆棋子
// type Pieces struct {
// 	Oushou   PiecesStatus // Gyokushou
// 	Kinshou  PiecesStatus
// 	Kinshou2 PiecesStatus
// 	Ginshou  PiecesStatus
// 	Ginshou2 PiecesStatus
// 	Keima    PiecesStatus
// 	Keima2   PiecesStatus
// 	Kyousha  PiecesStatus
// 	Kyousha2 PiecesStatus
// 	Kakugyou PiecesStatus
// 	Hisha    PiecesStatus
// 	Fuhyou   PiecesStatus
// 	Fuhyou2  PiecesStatus
// 	Fuhyou3  PiecesStatus
// 	Fuhyou4  PiecesStatus
// 	Fuhyou5  PiecesStatus
// 	Fuhyou6  PiecesStatus
// 	Fuhyou7  PiecesStatus
// 	Fuhyou8  PiecesStatus
// 	Fuhyou9  PiecesStatus
// 	Captured []string // 吃掉對手的棋子
// }

// 棋子位置
type Position struct {
	X int
	Y int
}
