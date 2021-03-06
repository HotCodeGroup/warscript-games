package main

import (
	"net/http"
	"strconv"

	"github.com/HotCodeGroup/warscript-games/jmodels"
	"github.com/HotCodeGroup/warscript-utils/utils"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// GetGame получает объект игры
func GetGame(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger(r, logger, "GetGame")
	errWriter := utils.NewErrorResponseWriter(w, logger)
	vars := mux.Vars(r)

	game, err := getGameBySlugImpl(vars["game_slug"])
	if err != nil {
		if errors.Cause(err) == utils.ErrNotExists {
			errWriter.WriteWarn(http.StatusNotFound, errors.Wrap(err, "game not exists"))
		} else {
			errWriter.WriteError(http.StatusInternalServerError, errors.Wrap(err, "get game method error"))
		}
		return
	}

	utils.WriteApplicationJSON(w, http.StatusOK, game)
}

// GetGameList gets list of games
func GetGameList(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger(r, logger, "GetGameList")
	errWriter := utils.NewErrorResponseWriter(w, logger)

	games, err := Games.GetGameList()
	if err != nil {
		errWriter.WriteError(http.StatusInternalServerError, errors.Wrap(err, "get game list method error"))

		return
	}

	respGames := make([]*jmodels.Game, len(games))
	for i, game := range games {
		respGames[i] = &jmodels.Game{
			Slug:           game.Slug,
			Title:          game.Title,
			BackgroundUUID: game.GetBackgroundUUID(), // точно 16 байт
		}
	}

	utils.WriteApplicationJSON(w, http.StatusOK, respGames)
}

// GetGameLeaderboard gets list of leaders in game
func GetGameLeaderboard(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger(r, logger, "GetGameLeaderboard")
	errWriter := utils.NewErrorResponseWriter(w, logger)
	vars := mux.Vars(r)

	query := r.URL.Query()
	limitParam, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limitParam = 5
	}
	offsetParam, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offsetParam = 0
	}

	leadersModels, err := Games.GetGameLeaderboardBySlug(vars["game_slug"], limitParam, offsetParam)
	if err != nil {
		if errors.Cause(err) == utils.ErrNotExists {
			errWriter.WriteWarn(http.StatusNotFound, errors.Wrap(err, "game not exists or offset is large"))
		} else {
			errWriter.WriteError(http.StatusInternalServerError, errors.Wrap(err, "get game method error"))
		}
		return
	}

	leaders := make([]*jmodels.ScoredUser, len(leadersModels))
	for i, leader := range leadersModels {
		leaders[i] = &jmodels.ScoredUser{
			InfoUser: jmodels.InfoUser{
				BasicUser: jmodels.BasicUser{
					Username:  leader.Username,
					PhotoUUID: leader.GetPhotoUUID(),
				},
				ID:     leader.ID,
				Active: leader.Active,
			},
			Score: leader.Score,
		}
	}

	utils.WriteApplicationJSON(w, http.StatusOK, leaders)
}

// GetGameTotalPlayers количество юзеров игравших в game_id
func GetGameTotalPlayers(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLogger(r, logger, "GetGameTotalPlayers")
	errWriter := utils.NewErrorResponseWriter(w, logger)
	vars := mux.Vars(r)

	totalCount, err := Games.GetGameTotalPlayersBySlug(vars["game_slug"])
	if err != nil {
		if errors.Cause(err) == utils.ErrNotExists {
			errWriter.WriteWarn(http.StatusNotFound, errors.Wrap(err, "game not exists"))
		} else {
			errWriter.WriteError(http.StatusInternalServerError, errors.Wrap(err, "get game method error"))
		}
		return
	}

	utils.WriteApplicationJSON(w, http.StatusOK, &struct {
		Count int64 `json:"count"`
	}{
		Count: totalCount,
	})
}
