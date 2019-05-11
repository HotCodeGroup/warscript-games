package main

import (
	"database/sql"

	"github.com/HotCodeGroup/warscript-utils/testutils"
	"github.com/HotCodeGroup/warscript-utils/utils"
)

type gameTest struct {
	games map[string]*GameModel

	testutils.Failer
}

func (gt *gameTest) GetGameBySlug(slug string) (*GameModel, error) {
	if err := gt.NextFail(); err != nil {
		return nil, err
	}

	g, ok := gt.games[slug]
	if !ok {
		return nil, utils.ErrNotExists
	}

	return g, nil
}

func (gt *gameTest) GetGameTotalPlayersBySlug(slug string) (int64, error) {
	if err := gt.NextFail(); err != nil {
		return 0, err
	}

	return 1, nil
}

func (gt *gameTest) GetGameList() ([]*GameModel, error) {
	if err := gt.NextFail(); err != nil {
		return nil, err
	}

	games := make([]*GameModel, 0, len(gt.games))
	for _, game := range gt.games {
		games = append(games, game)
	}

	return games, nil
}

func (gt *gameTest) GetGameLeaderboardBySlug(slug string, limit, offset int) ([]*ScoredUserModel, error) {
	if err := gt.NextFail(); err != nil {
		return nil, err
	}

	leaderboard := []*ScoredUserModel{
		{
			ID:        1,
			Username:  "GDVFox",
			PhotoUUID: sql.NullString{String: "2eb4a823-3a6d-4cba-8767-4d4946890f4f", Valid: true},
			Score:     1337,
		},
		{
			ID:       2,
			Username: "GDVFox1337",
			Score:    1337,
		},
	}

	return leaderboard, nil
}
