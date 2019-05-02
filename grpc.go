package main

import (
	"context"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/pkg/errors"
)

type GamesManager struct{}

func (gm *GamesManager) GetGameBySlug(ctx context.Context, gameSlug *models.GameSlug) (*models.InfoGame, error) {
	game, err := getGameBySlugImpl(gameSlug.Slug)
	if err != nil {
		return nil, errors.Wrap(err, "can not get game by slug")
	}

	return &models.InfoGame{
		ID:             game.ID,
		Slug:           game.Slug,
		Title:          game.Title,
		Description:    game.Description,
		Rules:          game.Rules,
		CodeExample:    game.CodeExample,
		BotCode:        game.BotCode,
		LogoUUID:       game.LogoUUID,
		BackgroundUUID: game.BackgroundUUID,
	}, nil
}
