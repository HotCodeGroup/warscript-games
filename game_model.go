package main

import (
	"context"
	"database/sql"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/HotCodeGroup/warscript-utils/postgresql"
	"github.com/HotCodeGroup/warscript-utils/utils"

	"github.com/pkg/errors"

	// драйвер Database
	_ "github.com/lib/pq"
)

var pqConn *sql.DB

// GameAccessObject DAO for User model
type GameAccessObject interface {
	GetGameBySlug(slug string) (*GameModel, error)
	GetGameTotalPlayersBySlug(slug string) (int64, error)
	GetGameList() ([]*GameModel, error)
	GetGameLeaderboardBySlug(slug string, limit, offset int) ([]*ScoredUserModel, error)
}

// AccessObject implementation of GameAccessObject
type AccessObject struct{}

// Games interface variable for models methods
var Games GameAccessObject

func init() {
	Games = &AccessObject{}
}

// GameModel модель для таблицы games
type GameModel struct {
	ID             int64
	Slug           string
	Title          string
	Description    string
	Rules          string
	CodeExample    string
	BotCode        string
	LogoUUID       sql.NullString
	BackgroundUUID sql.NullString
}

// GetLogoUUID возвращает LogoUUID или пустую строку, если его нет в базе
func (u *GameModel) GetLogoUUID() string {
	if u.LogoUUID.Valid {
		return u.LogoUUID.String
	}

	return ""
}

// GetBackgroundUUID возвращает BackgroundUUID или пустую строку, если его нет в базе
func (u *GameModel) GetBackgroundUUID() string {
	if u.BackgroundUUID.Valid {
		return u.BackgroundUUID.String
	}

	return ""
}

// ScoredUserModel User with score
type ScoredUserModel struct {
	ID        int64
	Username  string
	PhotoUUID sql.NullString
	Active    bool
	Score     int32
}

// GetPhotoUUID возвращает photoUUID или пустую строку, если его нет в базе
func (u *ScoredUserModel) GetPhotoUUID() string {
	if u.PhotoUUID.Valid {
		return u.PhotoUUID.String
	}

	return ""
}

// GetGameBySlug получает информацию об игре по slug
func (gs *AccessObject) GetGameBySlug(slug string) (*GameModel, error) {
	g, err := gs.getGameImpl(pqConn, "slug", slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotExists
		}

		return nil, errors.Wrapf(utils.ErrInternal, "get game by slug error: %s", err.Error())
	}

	return g, nil
}

// GetGameTotalPlayersBySlug получение общего количества игроков
func (gs *AccessObject) GetGameTotalPlayersBySlug(slug string) (int64, error) {
	tx, err := pqConn.Begin()
	if err != nil {
		return 0, errors.Wrapf(utils.ErrInternal, "can not open GetGameTotalPlayersByID transaction: %v", err)
	}

	//nolint: errcheck
	defer tx.Rollback()

	g, err := gs.getGameImpl(tx, "slug", slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, utils.ErrNotExists
		}

		return 0, errors.Wrapf(utils.ErrInternal, "GetGameTotalPlayersByID can not get game by id: %v", err)
	}

	var totalPlayers int64
	row := tx.QueryRow(`SELECT count(*) FROM users_games WHERE game_id = $1;`, &g.ID)
	if err = row.Scan(&totalPlayers); err != nil {
		return 0, errors.Wrapf(utils.ErrInternal, "get game total players error: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Wrapf(utils.ErrInternal, "can not commit GetGameTotalPlayersByID transaction: %v", err)
	}

	return totalPlayers, nil
}

// GetGameLeaderboardBySlug получаем leaderboard по slug
func (gs *AccessObject) GetGameLeaderboardBySlug(slug string, limit, offset int) ([]*ScoredUserModel, error) {
	// узнаём количество

	rows, err := pqConn.Query(`SELECT ug.user_id, ug.score FROM users_games ug
					RIGHT JOIN games g on ug.game_id = g.id
					WHERE g.slug = $1 ORDER BY ug.score DESC OFFSET $2 LIMIT $3;`, slug, offset, limit)
	if err != nil {
		return nil, errors.Wrapf(utils.ErrInternal, "get leaderboard error: %v", err)
	}
	defer rows.Close()

	IDs := make([]*models.UserID, 0)
	leaderboard := make([]*ScoredUserModel, 0)
	for rows.Next() {
		scoredUser := &ScoredUserModel{}
		err = rows.Scan(&scoredUser.ID, &scoredUser.Score)
		if err != nil {
			return nil, errors.Wrapf(utils.ErrInternal, "get leaderboard scan user error: %v", err)
		}
		leaderboard = append(leaderboard, scoredUser)
		IDs = append(IDs, &models.UserID{
			ID: scoredUser.ID,
		})
	}

	if len(leaderboard) == 0 {
		return nil, utils.ErrNotExists
	}

	users, err := authGPRC.GetUsersByIDs(context.Background(), &models.UserIDs{
		IDs: IDs,
	})

	if err != nil {
		return nil, errors.Wrapf(utils.ErrInternal, "can't connect to auth service to get users error: %v", err)
	}

	for i := 0; i < len(leaderboard); i++ {
		leaderboard[i].Username = users.Users[i].Username
		leaderboard[i].Active = users.Users[i].Active

		if users.Users[i].PhotoUUID == "" {
			leaderboard[i].PhotoUUID.Valid = false
		} else {
			leaderboard[i].PhotoUUID.String = users.Users[i].PhotoUUID
			leaderboard[i].PhotoUUID.Valid = true
		}
	}

	return leaderboard, nil
}

// GetGameList returns full list of active games
func (gs *AccessObject) GetGameList() ([]*GameModel, error) {
	rows, err := pqConn.Query(`SELECT g.id, g.slug, g.title, g.description,
								g.rules, g.code_example, g.bot_code, g.logo_uuid, g.background_uuid
								FROM games g ORDER BY g.id`)
	if err != nil {
		return nil, errors.Wrapf(utils.ErrInternal, "get game list error: %v", err)
	}
	defer rows.Close()

	games := make([]*GameModel, 0)
	for rows.Next() {
		g := &GameModel{}
		err = rows.Scan(&g.ID, &g.Slug, &g.Title, &g.Description,
			&g.Rules, &g.CodeExample, &g.BotCode, &g.LogoUUID, &g.BackgroundUUID)
		if err != nil {
			return nil, errors.Wrapf(utils.ErrInternal, "get games scan game error: %v", err)
		}
		games = append(games, g)
	}

	return games, nil
}

func (gs *AccessObject) getGameImpl(q postgresql.Queryer, field, value string) (*GameModel, error) {
	g := &GameModel{}

	//nolint: gosec уверены в том, что field корректно, так как сами его передаём
	row := q.QueryRow(`SELECT g.id, g.slug, g.title, g.description,
						g.rules, g.code_example, g.bot_code, g.logo_uuid, g.background_uuid
						FROM games g WHERE `+field+` = $1;`, value)
	if err := row.Scan(&g.ID, &g.Slug, &g.Title, &g.Description,
		&g.Rules, &g.CodeExample, &g.BotCode, &g.LogoUUID, &g.BackgroundUUID); err != nil {
		return nil, err
	}

	return g, nil
}
