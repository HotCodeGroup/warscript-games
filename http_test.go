package main

import (
	"database/sql"
	"io/ioutil"
	"testing"

	"github.com/HotCodeGroup/warscript-utils/logging"
	"github.com/HotCodeGroup/warscript-utils/testutils"
	"github.com/HotCodeGroup/warscript-utils/utils"
)

func init() {
	// выключаем логгер
	logger, _ = logging.NewLogger(ioutil.Discard, "")
}

func initTests() {
	Games = &gameTest{
		games: map[string]*GameModel{
			"pong": {
				ID:             1,
				Slug:           "pong",
				Title:          "Pong",
				Description:    "Very cool game(net)",
				Rules:          "Do not cheat, please",
				CodeExample:    "const a = 5;",
				BotCode:        "const a = 5;",
				LogoUUID:       sql.NullString{String: "2eb4a823-3a6d-4cba-8767-4d4946890f4f", Valid: true},
				BackgroundUUID: sql.NullString{String: "2eb4a823-3a6d-5xyz-8767-4d4946890f4f", Valid: true},
			},
		},
	}
}

type GameTestCase struct {
	testutils.Case
	Failure error
}

func runTableAPITests(t *testing.T, cases []*GameTestCase) {
	for i, c := range cases {
		runAPITest(t, i, c)
	}
}

func runAPITest(t *testing.T, i int, c *GameTestCase) {
	if c.Failure != nil {
		Games.(*gameTest).SetNextFail(c.Failure)
	}

	testutils.RunAPITest(t, i, &c.Case)
}

func TestGetGame(t *testing.T) {
	initTests()

	cases := []*GameTestCase{
		{ // Всё ок
			Case: testutils.Case{
				ExpectedCode: 200,
				ExpectedBody: `{"slug":"pong","title":"Pong","background_uuid":"2eb4a823-3a6d-5xyz-8767-4d4946890f4f",` +
					`"description":"Very cool game(net)","rules":"Do not cheat, please",` +
					`"code_example":"const a = 5;","bot_code":"const a = 5;",` +
					`"logo_uuid":"2eb4a823-3a6d-4cba-8767-4d4946890f4f"}`,
				Method:   "GET",
				Pattern:  "/games/{game_slug}",
				Endpoint: "/games/pong",
				Function: GetGame,
			},
		},
		{ // Такой игрули нет
			Case: testutils.Case{
				ExpectedCode: 404,
				ExpectedBody: `{"message":"game not exists: not_exists"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}",
				Endpoint:     "/games/not_pong",
				Function:     GetGame,
			},
			Failure: utils.ErrNotExists,
		},
		{ // база сломалась
			Case: testutils.Case{
				ExpectedCode: 500,
				ExpectedBody: `{"message":"get game method error: internal server error"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}",
				Endpoint:     "/games/not_pong",
				Function:     GetGame,
			},
			Failure: utils.ErrInternal,
		},
	}

	runTableAPITests(t, cases)
}

func TestGetGameList(t *testing.T) {
	initTests()

	cases := []*GameTestCase{
		{ // Всё ок
			Case: testutils.Case{
				ExpectedCode: 200,
				ExpectedBody: `[{"slug":"pong","title":"Pong","background_uuid":"2eb4a823-3a6d-5xyz-8767-4d4946890f4f"}]`,
				Method:       "GET",
				Pattern:      "/games",
				Endpoint:     "/games",
				Function:     GetGameList,
			},
		},
		{ // база сломалась
			Case: testutils.Case{
				ExpectedCode: 500,
				ExpectedBody: `{"message":"get game list method error: internal server error"}`,
				Method:       "GET",
				Pattern:      "/games",
				Endpoint:     "/games",
				Function:     GetGameList,
			},
			Failure: utils.ErrInternal,
		},
	}

	runTableAPITests(t, cases)
}

func TestGetGameLeaderboard(t *testing.T) {
	initTests()

	cases := []*GameTestCase{
		{ // Всё ок
			Case: testutils.Case{
				ExpectedCode: 200,
				ExpectedBody: `[{"username":"GDVFox","photo_uuid":"2eb4a823-3a6d-4cba-8767-4d4946890f4f","id":1,"active":false,"score":1337},` +
					`{"username":"GDVFox1337","photo_uuid":"","id":2,"active":false,"score":1337}]`,
				Method:   "GET",
				Pattern:  "/games/{game_slug}/leaderboard",
				Endpoint: "/games/pong/leaderboard",
				Function: GetGameLeaderboard,
			},
		},
		{ // Такой игрули нет
			Case: testutils.Case{
				ExpectedCode: 404,
				ExpectedBody: `{"message":"game not exists or offset is large: not_exists"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}/leaderboard",
				Endpoint:     "/games/pong/leaderboard",
				Function:     GetGameLeaderboard,
			},
			Failure: utils.ErrNotExists,
		},
		{ // база сломалась
			Case: testutils.Case{
				ExpectedCode: 500,
				ExpectedBody: `{"message":"get game method error: internal server error"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}/leaderboard",
				Endpoint:     "/games/pong/leaderboard",
				Function:     GetGameLeaderboard,
			},
			Failure: utils.ErrInternal,
		},
	}

	runTableAPITests(t, cases)
}

func TestGetGameTotalPlayers(t *testing.T) {
	initTests()

	cases := []*GameTestCase{
		{ // Всё ок
			Case: testutils.Case{
				ExpectedCode: 200,
				ExpectedBody: `{"count":1}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}/leaderboard/count",
				Endpoint:     "/games/pong/leaderboard/count",
				Function:     GetGameTotalPlayers,
			},
		},
		{ // Такой игрули нет
			Case: testutils.Case{
				ExpectedCode: 404,
				ExpectedBody: `{"message":"game not exists: not_exists"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}/leaderboard/count",
				Endpoint:     "/games/pong/leaderboard/count",
				Function:     GetGameTotalPlayers,
			},
			Failure: utils.ErrNotExists,
		},
		{ // база сломалась
			Case: testutils.Case{
				ExpectedCode: 500,
				ExpectedBody: `{"message":"get game method error: internal server error"}`,
				Method:       "GET",
				Pattern:      "/games/{game_slug}/leaderboard/count",
				Endpoint:     "/games/pong/leaderboard/count",
				Function:     GetGameTotalPlayers,
			},
			Failure: utils.ErrInternal,
		},
	}

	runTableAPITests(t, cases)
}
