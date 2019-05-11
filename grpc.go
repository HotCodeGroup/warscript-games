package main

import (
	"context"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/pkg/errors"
)

// GamesManager реализация GRPC сервера
type GamesManager struct{}

// GetGameBySlug отдаёт информацию о игре по заданному slug
func (gm *GamesManager) GetGameBySlug(ctx context.Context, gameSlug *models.GameSlug) (*models.InfoGame, error) {
	game, err := getGameBySlugImpl(gameSlug.Slug)
	if err != nil {
		return nil, errors.Wrap(err, "can not get game by slug")
	}

	return &models.InfoGame{
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
