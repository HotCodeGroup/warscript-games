package jmodels

// BasicUser базовые поля
type BasicUser struct {
	Username  string `json:"username"`
	PhotoUUID string `json:"photo_uuid"`
}

// InfoUser BasicUser, расширенный служебной инфой
type InfoUser struct {
	BasicUser
	ID     int64 `json:"id"`
	Active bool  `json:"active"`
}

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

// GameFull полная инфа об игре
type GameFull struct {
	Game
	Description string `json:"description"`
	Rules       string `json:"rules"`
	CodeExample string `json:"code_example"`
	BotCode     string `json:"bot_code"`
	LogoUUID    string `json:"logo_uuid"`
}
