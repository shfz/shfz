package model

type Api struct {
	Name      string      `json:"name" binding:"required"` // API名
	UsedFuzzs []*UsedFuzz // すでに使用したFuzz情報
}

type UsedFuzz struct {
	ID string `json:"id"`
	FuzzTexts
	IsClientFeedbacked bool // client側からのデータが反映されているか確認するためのフラグ
	ClientFeedback          // ClientFeedback
	IsServerFeedbacked bool // server側からのデータが反映されているか確認するためのフラグ
	ServerFeedback          // ServerFeedback
}

// filibから使用したFuzzを登録するときの構造体
type FuzzTexts struct {
	FuzzTexts []*FuzzText `json:"fuzz"`
}

type FuzzText struct {
	Name string `json:"name"` // fuzz name ex. username password
	Text string `json:"text"` // Fuzz文字列
}

type ClientFeedback struct {
	IsClientError *bool  `json:"isClientError" binding:"required"` // filib側でエラーが発生したか
	ClientError   string `json:"clientError"`                      // filib側で取得したエラー
}

type ServerFeedback struct {
	IsServerError     *bool    `json:"isServerError" binding:"required"` // webアプリ側でエラー(Exception)が発生したか
	ServerError       string   `json:"serverError"`                      // webアプリ側で取得したエラー(Exception)
	ServerErrorFile   string   `json:"serverErrorFile"`                  // webアプリ側で取得したエラー(Exception)の発生箇所のファイル
	ServerErrorLineNo int      `json:"serverErrorLineNo"`                // webアプリ側で取得したエラー(Exception)の発生箇所の行数
	ServerErrorFunc   string   `json:"serverErrorFunc"`                  // webアプリ側で取得したエラー(Exception)の発生箇所の関数名
	Frames            []*Frame `json:"frames" binding:"required"`        // webアプリ側で取得したFrame情報
	Framelen          int      `json:"framelen" binding:"required"`      // webアプリ側で取得したFrame情報
}

type Frame struct {
	Name string `json:"name" binding:"required"` // 関数名
	File string `json:"file" binding:"required"` // 関数が含まれるファイル
}

// Fuzz取得リクエスト
type FuzzInfo struct {
	Name      string `json:"name" binding:"required"`      // fuzz name ex. username password
	Charset   string `json:"charset" binding:"required"`   // Fuzzで使用可能な文字セット
	IsGenetic *bool  `json:"isGenetic" binding:"required"` // 遺伝的アルゴリズムを使用するか
	MaxLen    int    `json:"maxLen" binding:"required"`    // Fuzzの最大文字数
	MinLen    int    `json:"minLen" binding:"required"`    // Fuzzの最小文字数
}

// Fuzz情報登録リクエスト
type ApiParam struct {
	Name string `json:"name" binding:"required"`
	ID   string `json:"id" binding:"required"`
	FuzzTexts
}

// Markdown出力
type ReportClientError struct {
	ClientError string
	Fuzz        []*FuzzTexts
}

// Markdown出力
type ReportServerError struct {
	ServerError       string
	ServerErrorFile   string
	ServerErrorLineNo int
	ServerErrorFunc   string
	Fuzz              []*FuzzTexts
}

// Markdown出力リクエスト
type ReportReq struct {
	Hash   string `form:"hash"`
	Repo   string `form:"repo"`
	RunID  string `form:"id"`
	Job    string `form:"job"`
	Number string `form:"number"`
}
