package main

// ScoredUser инфа о юзере расширенная его баллами
type ScoredUser struct {
	InfoUser
	Score int32 `json:"score"`
}

// Game схема объекта игры для карусельки
type Game struct {
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	BackgroundUUID string `json:"background_uuid"`
}

type GameFull struct {
	Game
	ID          int64
	Description string `json:"description"`
	Rules       string `json:"rules"`
	CodeExample string `json:"code_example"`
	BotCode     string `json:"bot_code"`
	LogoUUID    string `json:"logo_uuid"`
}
