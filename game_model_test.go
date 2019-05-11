package main

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/HotCodeGroup/warscript-utils/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HotCodeGroup/warscript-utils/utils"
	"github.com/pkg/errors"
)

func TestGetGameBySlugOK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow(1, "pong", "Pong", "very cool", "do not cheat", "a=5", "a=5", "kek", "lol"))

	pqConn = db
	Games = &AccessObject{}

	game, err := Games.GetGameBySlug("pong")
	if err != nil {
		t.Errorf("TestGetGameBySlugOK got unexpected error: %v", err)
	}

	expected := &GameModel{
		ID:             1,
		Slug:           "pong",
		Title:          "Pong",
		Description:    "very cool",
		Rules:          "do not cheat",
		BotCode:        "a=5",
		CodeExample:    "a=5",
		LogoUUID:       sql.NullString{String: "kek", Valid: true},
		BackgroundUUID: sql.NullString{String: "lol", Valid: true},
	}

	if !reflect.DeepEqual(game, expected) {
		t.Errorf("TestGetGameBySlugOK got unexpected result: %v; expected: %v",
			game, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestGetGameBySlugOK there were unfulfilled expectations: %s", err)
	}
}

func getGetGameBySlugErrors(t *testing.T, queryError, expectedError error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong").WillReturnError(queryError)

	pqConn = db
	Games = &AccessObject{}

	if _, err = Games.GetGameBySlug("pong"); err != nil {
		if errors.Cause(err) != expectedError {
			t.Errorf("TestGetGameBySlugNotExists got unexpected error: %v", err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestGetGameBySlugNotExists there were unfulfilled expectations: %s", err)
	}
}
func TestGetGameBySlugNotExists(t *testing.T) {
	getGetGameBySlugErrors(t, sql.ErrNoRows, utils.ErrNotExists)
}

func TestGetGameBySlugInternal(t *testing.T) {
	getGetGameBySlugErrors(t, sql.ErrConnDone, utils.ErrInternal)
}

func TestGetGameTotalPlayersBySlugOK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs("pong").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow(1, "pong", "Pong", "very cool", "do not cheat", "a=5", "a=5", "kek", "lol"))
	mock.ExpectQuery("SELECT").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
			AddRow(1))
	mock.ExpectCommit()

	pqConn = db
	Games = &AccessObject{}

	total, err := Games.GetGameTotalPlayersBySlug("pong")
	if err != nil {
		t.Errorf("TestGetGameTotalPlayersBySlugOK got unexpected error: %v", err)
	}

	if total != 1 {
		t.Errorf("TestGetGameTotalPlayersBySlugOK got unexpected result: %v; expected: %v",
			total, 1)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestGetGameTotalPlayersBySlugOK there were unfulfilled expectations: %s", err)
	}
}

func getGameTotalPlayersBySlugError(t *testing.T, db *sql.DB,
	mock sqlmock.Sqlmock, expectedError error) {
	pqConn = db
	Games = &AccessObject{}

	_, err := Games.GetGameTotalPlayersBySlug("pong")
	if errors.Cause(err) != expectedError {
		t.Errorf("getGameTotalPlayersBySlugError got unexpected error: %v, expected: %v", err, expectedError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("getGameTotalPlayersBySlugError there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameTotalPlayersBySlugBeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
	getGameTotalPlayersBySlugError(t, db, mock, utils.ErrInternal)
}

func TestGetGameTotalPlayersBySlugGameNotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	getGameTotalPlayersBySlugError(t, db, mock, utils.ErrNotExists)
}

func TestGetGameTotalPlayersBySlugGameInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	getGameTotalPlayersBySlugError(t, db, mock, utils.ErrInternal)
}
func TestGetGameTotalPlayersBySlugCountInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs("pong").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow(1, "pong", "Pong", "very cool", "do not cheat", "a=5", "a=5", "kek", "lol"))
	mock.ExpectQuery("SELECT").WithArgs(1).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	getGameTotalPlayersBySlugError(t, db, mock, utils.ErrInternal)
}

func TestGetGameTotalPlayersBySlugCommitInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs("pong").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow(1, "pong", "Pong", "very cool", "do not cheat", "a=5", "a=5", "kek", "lol"))
	mock.ExpectQuery("SELECT").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
			AddRow(1))
	mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

	getGameTotalPlayersBySlugError(t, db, mock, utils.ErrInternal)
}

func TestGetGameLeaderboardBySlugOK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong", 0, 6).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "score"}).
			AddRow(1, 200).
			AddRow(2, 500))

	pqConn = db
	Games = &AccessObject{}
	authGPRC = &testutils.FakeAuthClient{
		Users: map[int64]*models.InfoUser{
			1: {
				ID:       1,
				Username: "kek",
				Active:   true,
			},
			2: {
				ID:        2,
				Username:  "kek1",
				PhotoUUID: "ea04741c-68d4-4e90-814d-44ffedf7c685",
				Active:    true,
			},
		},
	}

	expected := []*ScoredUserModel{
		{
			ID:       1,
			Username: "kek",
			Active:   true,
			Score:    200,
		},
		{
			ID:        2,
			Username:  "kek1",
			Active:    true,
			PhotoUUID: sql.NullString{String: "ea04741c-68d4-4e90-814d-44ffedf7c685", Valid: true},
			Score:     500,
		},
	}

	scored, err := Games.GetGameLeaderboardBySlug("pong", 6, 0)
	if err != nil {
		t.Errorf("GetGameLeaderboardBySlug got unexpected error: %v", err)
	}

	if !reflect.DeepEqual(scored, expected) {
		t.Errorf("GetGameLeaderboardBySlug got unexpected result: %v; expected: %v",
			scored, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("GetGameLeaderboardBySlug there were unfulfilled expectations: %s", err)
	}
}

func getGameLeaderboardBySlugError(t *testing.T, db *sql.DB,
	mock sqlmock.Sqlmock, expectedError error) {
	pqConn = db
	Games = &AccessObject{}

	_, err := Games.GetGameLeaderboardBySlug("pong", 6, 0)
	if errors.Cause(err) != expectedError {
		t.Errorf("getGameLeaderboardBySlugError got unexpected error: %v, expected: %v", err, expectedError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("getGameLeaderboardBySlugError there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameLeaderboardBySlugAuthError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong", 0, 6).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "score"}).
			AddRow(1, 200).
			AddRow(2, 500))

	pqConn = db
	Games = &AccessObject{}
	authGPRC = &testutils.FakeAuthClient{}
	authGPRC.(*testutils.FakeAuthClient).SetNextFail(utils.ErrInternal)

	_, err = Games.GetGameLeaderboardBySlug("pong", 6, 0)
	if errors.Cause(err) != utils.ErrInternal {
		t.Errorf("GetGameLeaderboardBySlug got unexpected error: %v, expected: %v", err, utils.ErrInternal)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("GetGameLeaderboardBySlug there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameLeaderboardBySlugLeaderboardInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
	getGameLeaderboardBySlugError(t, db, mock, utils.ErrInternal)
}

func TestGetGameLeaderboardBySlugLeaderboardEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong", 0, 6).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "score"}))
	getGameLeaderboardBySlugError(t, db, mock, utils.ErrNotExists)
}

func TestGetGameLeaderboardBySlugLeaderboardScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WithArgs("pong", 0, 6).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "score"}).AddRow("kek", "lol"))
	getGameLeaderboardBySlugError(t, db, mock, utils.ErrInternal)
}

func TestGetGameListOK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow(1, "pong", "Pong", "very cool", "do not cheat", "a=5", "a=5", "kek", "lol"))

	pqConn = db
	Games = &AccessObject{}

	games, err := Games.GetGameList()
	if err != nil {
		t.Errorf("TestGetGameListOK got unexpected error: %v", err)
	}

	expected := []*GameModel{
		{
			ID:             1,
			Slug:           "pong",
			Title:          "Pong",
			Description:    "very cool",
			Rules:          "do not cheat",
			BotCode:        "a=5",
			CodeExample:    "a=5",
			LogoUUID:       sql.NullString{String: "kek", Valid: true},
			BackgroundUUID: sql.NullString{String: "lol", Valid: true},
		},
	}

	if !reflect.DeepEqual(games, expected) {
		t.Errorf("TestGetGameListOK got unexpected result: %v; expected: %v",
			games, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("TestGetGameListOK there were unfulfilled expectations: %s", err)
	}
}

func getGameListError(t *testing.T, db *sql.DB,
	mock sqlmock.Sqlmock, expectedError error) {
	pqConn = db
	Games = &AccessObject{}

	_, err := Games.GetGameList()
	if errors.Cause(err) != expectedError {
		t.Errorf("getGameLeaderboardBySlugError got unexpected error: %v, expected: %v", err, expectedError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("getGameLeaderboardBySlugError there were unfulfilled expectations: %s", err)
	}
}

func TestGetGameListQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
	getGameListError(t, db, mock, utils.ErrInternal)
}

func TestGetGameListScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "rules",
			"code_example", "bot_code", "logo_uuid", "background_uuid"}).
			AddRow("kek", 2, 3, 4, "do not cheat", "a=5", "a=5", "kek", "lol"))
	getGameListError(t, db, mock, utils.ErrInternal)
}
